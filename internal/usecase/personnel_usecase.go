package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/models"
	"github.com/keremdursn/hospital-case/internal/repository"
)

type PersonnelUsecase interface {
	ListAllJobGroups() ([]dto.JobGroupLookup, error)
	ListTitleByJobGroup(jobGroupID uint) ([]dto.TitleLookup, error)

	AddStaff(req *dto.AddStaffRequest, hospitalID uint) (*dto.StaffResponse, error)
	UpdateStaff(id uint, req *dto.UpdateStaffRequest, hospitalID uint) (*dto.StaffResponse, error)
	DeleteStaff(id, hospitalID uint) error
	ListStaff(hospitalID uint, filter dto.StaffListFilter, page, size int) (*dto.StaffListResponse, error)
}

type personnelUsecase struct {
	repo  repository.PersonnelRepository
	redis *redis.Client
}

func NewPersonnelUsecase(repo repository.PersonnelRepository, redis *redis.Client) PersonnelUsecase {
	return &personnelUsecase{
		repo:  repo,
		redis: redis,
	}
}

func (u *personnelUsecase) ListAllJobGroups() ([]dto.JobGroupLookup, error) {
	ctx := context.Background()
	cacheKey := "job_groups"

	// Önce Redis'te var mı bak
	cached, err := u.redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var resp []dto.JobGroupLookup
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return resp, nil
		}
	}

	// 2. Yoksa DB'den çek
	groups, err := u.repo.GetAllJobGroups()
	if err != nil {
		return nil, err
	}

	resp := make([]dto.JobGroupLookup, 0, len(groups))
	for _, g := range groups {
		resp = append(resp, dto.JobGroupLookup{
			ID:   g.ID,
			Name: g.Name,
		})
	}

	// 3. Redis'e yaz
	if data, err := json.Marshal(resp); err == nil {
		_ = u.redis.Set(ctx, cacheKey, data, 0).Err() // Hata olursa cache'siz devam et
	}
	return resp, nil
}

func (u *personnelUsecase) ListTitleByJobGroup(jobGroupID uint) ([]dto.TitleLookup, error) {
	ctx := context.Background()
	cacheKey := "titles_by_jobgroup_" + string(rune(jobGroupID))

	// Önce Redis'te var mı bak
	cached, err := u.redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var resp []dto.TitleLookup
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return resp, nil
		}
	}

	// Yoksa DB'den çek
	titles, err := u.repo.GetAllTitlesByJobGroup(jobGroupID)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.TitleLookup, 0, len(titles))
	for _, t := range titles {
		resp = append(resp, dto.TitleLookup{
			ID:   t.ID,
			Name: t.Name,
		})
	}

	// Redis'e yaz
	if data, err := json.Marshal(resp); err == nil {
		_ = u.redis.Set(ctx, cacheKey, data, 0).Err()
	}
	return resp, nil
}

func (u *personnelUsecase) AddStaff(req *dto.AddStaffRequest, hospitalID uint) (*dto.StaffResponse, error) {
	// TC ve telefon benzersiz mi?
	exists, err := u.repo.IsTCOrPhoneExists(req.TC, req.Phone)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("staff with given TC or phone already exists")
	}

	// JobGroup ve Title kontrolü
	jobGroup, err := u.repo.GetJobGroupByID(req.JobGroupID)
	if err != nil {
		return nil, errors.New("job group not found")
	}

	title, err := u.repo.GetTitleByID(req.TitleID)
	if err != nil {
		return nil, errors.New("title not found")
	}

	if title.JobGroupID != jobGroup.ID {
		return nil, errors.New("title does not belong to the selected job group")
	}

	if title.Name == "Başhekim" {
		count, err := u.repo.CountHospitalHeads(hospitalID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("there can be only one Başhekim in a hospital")
		}
	}

	// Poliklinik kontrolü
	var polyName *string
	if req.HospitalPolyclinicID != nil {
		hp, err := u.repo.GetHospitalPolyclinicByID(*req.HospitalPolyclinicID)
		if err != nil {
			return nil, errors.New("hospital polyclinic not found")
		}
		if hp.HospitalID != hospitalID {
			return nil, errors.New("polyclinic does not belong to your hospital")
		}
		p, _ := u.repo.GetPolyclinicByID(hp.PolyclinicID)
		polyName = &p.Name
	}

	staff := models.Staff{
		FirstName:            req.FirstName,
		LastName:             req.LastName,
		TC:                   req.TC,
		Phone:                req.Phone,
		JobGroupID:           req.JobGroupID,
		TitleID:              req.TitleID,
		HospitalID:           hospitalID,
		HospitalPolyclinicID: req.HospitalPolyclinicID,
		WorkingDays:          req.WorkingDays,
	}

	if err := u.repo.CreateStaff(&staff); err != nil {
		return nil, err
	}

	return &dto.StaffResponse{
		ID:                   staff.ID,
		FirstName:            staff.FirstName,
		LastName:             staff.LastName,
		TC:                   staff.TC,
		Phone:                staff.Phone,
		JobGroupID:           jobGroup.ID,
		JobGroupName:         jobGroup.Name,
		TitleID:              title.ID,
		TitleName:            title.Name,
		HospitalPolyclinicID: staff.HospitalPolyclinicID,
		PolyclinicName:       polyName,
		WorkingDays:          staff.WorkingDays,
	}, nil
}

func (u *personnelUsecase) UpdateStaff(id uint, req *dto.UpdateStaffRequest, hospitalID uint) (*dto.StaffResponse, error) {
	staff, err := u.repo.GetStaffByID(id)
	if err != nil {
		return nil, errors.New("staff not found")
	}
	if staff.HospitalID != hospitalID {
		return nil, errors.New("forbidden: cannot update staff from another hospital")
	}

	// TC ve telefon benzersiz mi? (kendisi hariç)
	exists, err := u.repo.IsTCOrPhoneExistsExcludeID(id, req.TC, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("another staff with given TC or phone already exists")
	}

	// Meslek grubu ve unvan var mı, ilişkili mi?
	jobGroup, err := u.repo.GetJobGroupByID(req.JobGroupID)
	if err != nil {
		return nil, errors.New("job group not found")
	}

	title, err := u.repo.GetTitleByID(req.TitleID)
	if err != nil {
		return nil, errors.New("title not found")
	}
	if title.JobGroupID != jobGroup.ID {
		return nil, errors.New("title does not belong to the selected job group")
	}

	if title.Name == "Başhekim" {
		count, err := u.repo.CountHospitalHeads(hospitalID)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, errors.New("there can be only one Başhekim in a hospital")
		}
	}

	var polyName *string
	if req.HospitalPolyclinicID != nil {
		hp, err := u.repo.GetHospitalPolyclinicByID(*req.HospitalPolyclinicID)
		if err != nil {
			return nil, errors.New("hospital polyclinic not found")
		}
		if hp.HospitalID != hospitalID {
			return nil, errors.New("polyclinic does not belong to your hospital")
		}
		p, _ := u.repo.GetPolyclinicByID(hp.PolyclinicID)
		polyName = &p.Name
	}

	staff.FirstName = req.FirstName
	staff.LastName = req.LastName
	staff.TC = req.TC
	staff.Phone = req.Phone
	staff.JobGroupID = req.JobGroupID
	staff.TitleID = req.TitleID
	staff.HospitalPolyclinicID = req.HospitalPolyclinicID
	staff.WorkingDays = req.WorkingDays

	if err := u.repo.UpdateStaff(staff); err != nil {
		return nil, err
	}

	return &dto.StaffResponse{
		ID:                   staff.ID,
		FirstName:            staff.FirstName,
		LastName:             staff.LastName,
		TC:                   staff.TC,
		Phone:                staff.Phone,
		JobGroupID:           jobGroup.ID,
		JobGroupName:         jobGroup.Name,
		TitleID:              title.ID,
		TitleName:            title.Name,
		HospitalPolyclinicID: staff.HospitalPolyclinicID,
		PolyclinicName:       polyName,
		WorkingDays:          staff.WorkingDays,
	}, nil
}

func (u *personnelUsecase) DeleteStaff(id, hospitalID uint) error {
	staff, err := u.repo.GetStaffByID(id)
	if err != nil {
		return errors.New("staff not found")
	}
	if staff.HospitalID != hospitalID {
		return errors.New("forbidden: cannot delete staff from another hospital")
	}
	return u.repo.DeleteStaff(staff)
}

func (u *personnelUsecase) ListStaff(hospitalID uint, filter dto.StaffListFilter, page, size int) (*dto.StaffListResponse, error) {
	// 1. Personelleri filtreyle listele
	staffs, err := u.repo.ListStaffWithFilter(hospitalID, filter, page, size)
	if err != nil {
		return nil, err
	}

	// 2. Toplam kayıt sayısını al
	totalCount, err := u.repo.CountStaffWithFilter(hospitalID, filter)
	if err != nil {
		return nil, err
	}

	// 3. Dönüştürülecek response dizisi hazırlanıyor
	resp := make([]dto.StaffResponse, 0, len(staffs))
	for _, s := range staffs {
		// Her staff için jobGroup ve title bilgilerini getir
		jobGroup, err := u.repo.GetJobGroupByID(s.JobGroupID)
		if err != nil {
			return nil, errors.New("job group not found for staff")
		}

		title, err := u.repo.GetTitleByID(s.TitleID)
		if err != nil {
			return nil, errors.New("title not found for staff")
		}

		// Eğer poliklinik ID varsa, poliklinik adını al
		var polyName *string
		if s.HospitalPolyclinicID != nil {
			hp, err := u.repo.GetHospitalPolyclinicByID(*s.HospitalPolyclinicID)
			if err != nil {
				return nil, err
			}
			p, _ := u.repo.GetPolyclinicByID(hp.PolyclinicID)
			polyName = &p.Name
		}

		resp = append(resp, dto.StaffResponse{
			ID:                   s.ID,
			FirstName:            s.FirstName,
			LastName:             s.LastName,
			TC:                   s.TC,
			Phone:                s.Phone,
			JobGroupID:           jobGroup.ID,
			JobGroupName:         jobGroup.Name,
			TitleID:              title.ID,
			TitleName:            title.Name,
			HospitalPolyclinicID: s.HospitalPolyclinicID,
			PolyclinicName:       polyName,
			WorkingDays:          s.WorkingDays,
		})
	}

	return &dto.StaffListResponse{
		Staff: resp,
		Total: int(totalCount),
		Page:  page,
		Size:  size,
	}, nil
}
