package database

import (
	"fmt"

	"github.com/keremdursn/hospital-case/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	fmt.Println("Running database migrations...")

	// Burada modelleri ekliyoruz
	return db.AutoMigrate(
		&models.City{},
		&models.District{},
		&models.Hospital{},
		&models.Authority{},
		&models.Polyclinic{},
		&models.HospitalPolyclinic{},
		&models.JobGroup{},
		&models.Title{},
		&models.Staff{},
	)
}
