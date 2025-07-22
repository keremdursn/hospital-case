package usecase

import (
	"errors"

	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/models"
	"github.com/keremdursn/hospital-case/internal/repository"
)

type PolyclinicUsecase interface {
	ListAllPolyclinics() ([]dto.PolyclinicLookup, error)
	AddPolyclinicToHospital(req *dto.AddHospitalPolyclinicRequest, hospitalID uint) (*dto.HospitalPolyclinicResponse, error)
}

type polyclinicUsecase struct {
	repo repository.PolyclinicRepository
}

func NewPolyclinicUsecase(repo repository.PolyclinicRepository) PolyclinicUsecase {
	return &polyclinicUsecase{repo: repo}
}

func (u *polyclinicUsecase) ListAllPolyclinics() ([]dto.PolyclinicLookup, error) {
	polys, err := u.repo.GetAllPolyclinic()
	if err != nil {
		return nil, err
	}

	resp := make([]dto.PolyclinicLookup, 0, len(polys))
	for _, p := range polys {
		resp = append(resp, dto.PolyclinicLookup{
			ID:   p.ID,
			Name: p.Name,
		})
	}
	return resp, nil
}

func (u *polyclinicUsecase) AddPolyclinicToHospital(req *dto.AddHospitalPolyclinicRequest, hospitalID uint) (*dto.HospitalPolyclinicResponse, error) {
	alreadyExists, err := u.repo.IsPolyclinicAlreadyAdded(hospitalID, req.PolyclinicID)
	if err != nil {
		return nil, err
	}

	if alreadyExists {
		return nil, errors.New("this polyclinic is already added to the hospital")
	}

	poly, err := u.repo.GetPolyclinicByID(req.PolyclinicID)
	if err != nil {
		return nil, err
	}

	hp := &models.HospitalPolyclinic{
		HospitalID:   hospitalID,
		PolyclinicID: req.PolyclinicID,
	}

	if err := u.repo.CreateHospitalPolyclinic(hp); err != nil {
		return nil, err
	}

	return &dto.HospitalPolyclinicResponse{
		ID: hp.ID,
		PolyclinicID: poly.ID,
		PolyclinicName: poly.Name,
	}, nil

}
