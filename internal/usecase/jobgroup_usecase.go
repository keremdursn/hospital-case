package usecase

import (
	"errors"

	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/models"
	"github.com/keremdursn/hospital-case/internal/repository"
)

type PersonnelUsecase interface {
	ListAllJobGroups() ([]dto.JobGroupLookup, error)
	ListTitleByJobGroup(jobGroupID uint) ([]dto.TitleLookup, error)

	AddStaff(req *dto.AddStaffRequest, hospitalID uint) (*dto.StaffResponse, error)
}

type personnelUsecase struct {
	repo repository.PersonnelRepository
}

func NewPersonnelUsecase(repo repository.PersonnelRepository) PersonnelUsecase {
	return &personnelUsecase{repo: repo}
}

func (u *personnelUsecase) ListAllJobGroups() ([]dto.JobGroupLookup, error) {
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
	return resp, nil
}

func (u *personnelUsecase) ListTitleByJobGroup(jobGroupID uint) ([]dto.TitleLookup, error) {
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

	// Başhekim kontrolü
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
