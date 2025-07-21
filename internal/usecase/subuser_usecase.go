package usecase

import (
	"errors"

	"github.com/keremdursn/hospital-case/internal/models"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/pkg/utils"
)

type CreateSubUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TC        string `json:"tc"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type SubUserResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	TC        string `json:"tc"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

type SubUserUsecase interface {
	CreateSubUser(req *CreateSubUserRequest, hospitalID uint) (*SubUserResponse, error)
}

type subUserUsecase struct {
	repo repository.SubUserRepository
}

func NewSubUserUsecase(repo repository.SubUserRepository) SubUserUsecase {
	return &subUserUsecase{repo: repo}
}

func (u *subUserUsecase) CreateSubUser(req *CreateSubUserRequest, hospitalID uint) (*SubUserResponse, error) {
	if req.Role != "yetkili" && req.Role != "calisan" {
		return nil, errors.New("role must be 'yetkili' or 'calisan'")
	}

	exists, err := u.repo.IsAuthorityExists(req.TC, req.Email, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with given TC, email, or phone already exists")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	auth := &models.Authority{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		TC:         req.TC,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   hashedPassword,
		Role:       req.Role,
		HospitalID: hospitalID,
	}

	if err := u.repo.CreateAuthority(auth); err != nil {
		return nil, err
	}

	return &SubUserResponse{
		ID:        auth.ID,
		FirstName: auth.FirstName,
		LastName:  auth.LastName,
		TC:        auth.TC,
		Email:     auth.Email,
		Phone:     auth.Phone,
		Role:      auth.Role,
	}, nil
}
