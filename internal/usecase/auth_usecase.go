package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/database"
	"github.com/keremdursn/hospital-case/internal/models"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/pkg/utils"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(req *RegisterRequest) (*models.Authority, error)
	Login(req *LoginRequest, cfg *config.Config) (*LoginResponse, error)
	ForgotPassword(req *ForgotPasswordRequest) (*ForgotPasswordResponse, error)
	ResetPassword(req *ResetPasswordRequest) error
}

type authUsecase struct {
	authRepo repository.AuthRepository
}

func NewAuthUsecase(r repository.AuthRepository) AuthUsecase {
	return &authUsecase{authRepo: r}
}

type RegisterRequest struct {
	HospitalName   string
	TaxNumber      string
	HospitalEmail  string
	HospitalPhone  string
	Address        string
	CityID         uint
	DistrictID     uint
	AuthorityFName string
	AuthorityLName string
	AuthorityTC    string
	AuthorityEmail string
	AuthorityPhone string
	Password       string
}

type LoginRequest struct {
	Credential string
	Password   string
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ForgotPasswordRequest struct {
	Phone string `json:"phone"`
}

type ForgotPasswordResponse struct {
	Code string `json:"code"`
}

type ResetPasswordRequest struct {
	Phone             string `json:"phone"`
	Code              string `json:"code"`
	NewPassword       string `json:"new_password"`
	RepeatNewPassword string `json:"repeat_new_password"`
}

func (u *authUsecase) Register(req *RegisterRequest) (*models.Authority, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Check uniqueness
	if exists, _ := u.authRepo.IsHospitalExists(req.TaxNumber, req.HospitalEmail, req.HospitalPhone); exists {
		return nil, errors.New("hospital already exists")
	}
	if exists, _ := u.authRepo.IsAuthorityExists(req.AuthorityTC, req.AuthorityEmail, req.AuthorityPhone); exists {
		return nil, errors.New("authority already exists")
	}

	// Create Hospital
	hospital := &models.Hospital{
		Name:       req.HospitalName,
		TaxNumber:  req.TaxNumber,
		Email:      req.HospitalEmail,
		Phone:      req.HospitalPhone,
		Address:    req.Address,
		CityID:     req.CityID,
		DistrictID: req.DistrictID,
	}
	if err := u.authRepo.CreateHospital(hospital); err != nil {
		return nil, err
	}

	// Create Authority
	authority := &models.Authority{
		FirstName:  req.AuthorityFName,
		LastName:   req.AuthorityLName,
		TC:         req.AuthorityTC,
		Email:      req.AuthorityEmail,
		Phone:      req.AuthorityPhone,
		Password:   hashedPassword,
		Role:       "yetkili",
		HospitalID: hospital.ID,
	}
	if err := u.authRepo.CreateAuthority(authority); err != nil {
		return nil, err
	}

	return authority, nil
}

func (u *authUsecase) Login(req *LoginRequest, cfg *config.Config) (*LoginResponse, error) {
	authority, err := u.authRepo.GetAuthorityByEmailOrPhone(req.Credential)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Parola kontrolü
	if !utils.CheckPasswordHash(req.Password, authority.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Token üret
	token, err := utils.GenerateToken(authority.ID, authority.HospitalID, authority.Role, cfg)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{Token: token}, nil
}

func (u *authUsecase) ForgotPassword(req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	_, err := u.authRepo.GetAuthorityByPhone(req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	code := utils.GenerateResetCode()

	ctx := context.Background()
	if err := database.RDB.Set(ctx, "reset_code:"+req.Phone, code, 5*time.Minute).Err(); err != nil {
		return nil, err
	}

	return &ForgotPasswordResponse{Code: code}, nil
}

func (u *authUsecase) ResetPassword(req *ResetPasswordRequest) error {
	if req.NewPassword != req.RepeatNewPassword {
		return errors.New("passwords do not match")
	}

	ctx := context.Background()
	storedCode, err := database.RDB.Get(ctx, "reset_code:"+req.Phone).Result()
	if err != nil || storedCode != req.Code {
		return errors.New("invalid or expired code")
	}

	authority, err := u.authRepo.GetAuthorityByPhone(req.Phone)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := u.authRepo.UpdateAuthorityPassword(authority, hashedPassword); err != nil {
		return err
	}

	_ = database.RDB.Del(ctx, "reset_code:"+req.Phone).Err()
	return nil
}
