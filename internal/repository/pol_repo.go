package repository

import (
	"errors"

	"github.com/keremdursn/hospital-case/internal/models"
	"gorm.io/gorm"
)

type PolyclinicRepository interface {
	GetAllPolyclinic() ([]models.Polyclinic, error)
	IsPolyclinicAlreadyAdded(hospitalID, PolyclinicID uint) (bool, error)
	GetPolyclinicByID(PolyclinicID uint) (*models.Polyclinic, error)
	CreateHospitalPolyclinic(hp *models.HospitalPolyclinic) error
}

type polyclinicRepository struct {
	db *gorm.DB
}

func NewPolyclinicRepository(db *gorm.DB) PolyclinicRepository {
	return &polyclinicRepository{db: db}
}

func (r *polyclinicRepository) GetAllPolyclinic() ([]models.Polyclinic, error) {
	var polys []models.Polyclinic
	if err := r.db.Find(&polys).Error; err != nil {
		return nil, err
	}
	return polys, nil
}

func (r *polyclinicRepository) IsPolyclinicAlreadyAdded(hospitalID, PolyclinicID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.HospitalPolyclinic{}).
		Where("hospital_id = ? AND polyclinic_id = ?", hospitalID, PolyclinicID).
		Count(&count).Error
	return count > 0, err

}

func (r *polyclinicRepository) GetPolyclinicByID(PolyclinicID uint) (*models.Polyclinic, error) {
	var poly models.Polyclinic
	if err := r.db.First(&poly, PolyclinicID).Error; err != nil {
		return nil, errors.New("polyclinic not found")
	}
	return &poly, nil
}

func (r *polyclinicRepository) CreateHospitalPolyclinic(hp *models.HospitalPolyclinic) error {
	return r.db.Create(hp).Error
}
