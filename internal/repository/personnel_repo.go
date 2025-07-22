package repository

import (
	"github.com/keremdursn/hospital-case/internal/database"
	"github.com/keremdursn/hospital-case/internal/models"
	"gorm.io/gorm"
)

type PersonnelRepository interface {
	GetAllJobGroups() ([]models.JobGroup, error)
	GetAllTitlesByJobGroup(jobGroupID uint) ([]models.Title, error)

	IsTCOrPhoneExists(tc, phone string) (bool, error)
	GetJobGroupByID(id uint) (*models.JobGroup, error)
	GetTitleByID(id uint) (*models.Title, error)
	CountHospitalHeads(hospitalID uint) (int64, error)
	GetHospitalPolyclinicByID(id uint) (*models.HospitalPolyclinic, error)
	GetPolyclinicByID(id uint) (*models.Polyclinic, error)
	CreateStaff(staff *models.Staff) error
}

type personnelRepository struct {
	db *gorm.DB
}

func NewPersonnelRepository(db *gorm.DB) PersonnelRepository {
	return &personnelRepository{db: db}
}

func (r *personnelRepository) GetAllJobGroups() ([]models.JobGroup, error) {
	var groups []models.JobGroup
	if err := r.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *personnelRepository) GetAllTitlesByJobGroup(jobGroupID uint) ([]models.Title, error) {
	var titles []models.Title
	if err := r.db.Where("job_group_id = ?", jobGroupID).Find(&titles).Error; err != nil {
		return nil, err
	}
	return titles, nil
}

func (r *personnelRepository) IsTCOrPhoneExists(tc, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Staff{}).Where("tc = ? OR phone = ?", tc, phone).Count(&count).Error
	return count > 0, err
}

func (r *personnelRepository) GetJobGroupByID(id uint) (*models.JobGroup, error) {
	var JobGroup models.JobGroup
	if err := r.db.First(&JobGroup, id).Error; err != nil {
		return nil, err
	}
	return &JobGroup, nil
}

func (r *personnelRepository) GetTitleByID(id uint) (*models.Title, error) {
	var title models.Title
	if err := r.db.First(&title, id).Error; err != nil {
		return nil, err
	}
	return &title, nil
}

func (r *personnelRepository) CountHospitalHeads(hospitalID uint) (int64, error) {
	var count int64
	err := database.DB.Table("staffs").
		Joins("JOIN titles ON staffs.title_id = titles.id").
		Where("titles.name = ? AND staffs.hospital_id = ?", "Başhekim", hospitalID).
		Count(&count).Error
	return count, err
}

func (r *personnelRepository) GetHospitalPolyclinicByID(id uint) (*models.HospitalPolyclinic, error) {
	var hp models.HospitalPolyclinic
	if err := r.db.First(&hp, id).Error; err != nil {
		return nil, err
	}
	return &hp, nil
}

func (r *personnelRepository) GetPolyclinicByID(id uint) (*models.Polyclinic, error) {
	var p models.Polyclinic
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r * personnelRepository) CreateStaff(staff *models.Staff) error {
	return r.db.Create(staff).Error
}