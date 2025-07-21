package repository

import (
	"github.com/keremdursn/hospital-case/internal/models"
	"gorm.io/gorm"
)

type SubUserRepository interface {
	IsAuthorityExists(tc, email, phone string) (bool, error)
	CreateAuthority(authority *models.Authority) error
}

type subUserRepository struct {
	db *gorm.DB
}

func NewSubUserRepository(db *gorm.DB) SubUserRepository {
	return &subUserRepository{db: db}
}

func (r *subUserRepository) IsAuthorityExists(tc, email, phone string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Authority{}).Where("tc = ? OR email = ? OR phone = ?", tc, email, phone).Count(&count).Error
	return count > 0, err
}

func (r *subUserRepository) CreateAuthority(authority *models.Authority) error {
	return r.db.Create(authority).Error
}
