package usecase

import (
	"errors"

	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/repository"
)

type HospitalUsecase interface {
	GetHospitalByID(hospitalID uint) (*dto.HospitalResponse, error)
	UpdateHospital(hospitalID uint, req *dto.UpdateHospitalRequest) (*dto.HospitalResponse, error)
}

type hospitalUsecase struct {
	repo repository.HospitalRepository
}

func NewHospitalUsecase(repo repository.HospitalRepository) HospitalUsecase {
	return &hospitalUsecase{repo: repo}
}

func (u *hospitalUsecase) GetHospitalByID(hospitalID uint) (*dto.HospitalResponse, error) {
	hospital, err := u.repo.GetByID(hospitalID)
	if err != nil {
		return nil, err
	}
	city, _ := u.repo.GetCityByID(hospital.CityID)
	district, _ := u.repo.GetDistrictByID(hospital.DistrictID)

	return &dto.HospitalResponse{
		ID:           hospital.ID,
		Name:         hospital.Name,
		TaxNumber:    hospital.TaxNumber,
		Email:        hospital.Email,
		Phone:        hospital.Phone,
		Address:      hospital.Address,
		CityID:       hospital.CityID,
		CityName:     city.Name,
		DistrictID:   hospital.DistrictID,
		DistrictName: district.Name,
	}, nil
}

func (u *hospitalUsecase) UpdateHospital(hospitalID uint, req *dto.UpdateHospitalRequest) (*dto.HospitalResponse, error) {
	hospital, err := u.repo.GetByID(hospitalID)
	if err != nil {
		return nil, err
	}
	conflict, _ := u.repo.IsUniqueFieldsConflict(hospitalID, req.TaxNumber, req.Email, req.Phone)
	if conflict {
		return nil, errors.New("another hospital with given tax number, email, or phone already exists")
	}

	city, err := u.repo.GetCityByID(req.CityID)
	if err != nil {
		return nil, err
	}
	district, err := u.repo.GetDistrictByID(req.DistrictID)
	if err != nil {
		return nil, err
	}

	hospital.Name = req.Name
	hospital.TaxNumber = req.TaxNumber
	hospital.Email = req.Email
	hospital.Phone = req.Phone
	hospital.Address = req.Address
	hospital.CityID = req.CityID
	hospital.DistrictID = req.DistrictID

	if err := u.repo.Update(hospital); err != nil {
		return nil, err
	}

	return &dto.HospitalResponse{
		ID:           hospital.ID,
		Name:         hospital.Name,
		TaxNumber:    hospital.TaxNumber,
		Email:        hospital.Email,
		Phone:        hospital.Phone,
		Address:      hospital.Address,
		CityID:       hospital.CityID,
		CityName:     city.Name,
		DistrictID:   hospital.DistrictID,
		DistrictName: district.Name,
	}, nil
}
