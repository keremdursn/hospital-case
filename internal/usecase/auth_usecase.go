package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/keremdursn/hospital-case/internal/config"
	"github.com/keremdursn/hospital-case/internal/dto"
	"github.com/keremdursn/hospital-case/internal/models"
	"github.com/keremdursn/hospital-case/internal/repository"
	"github.com/keremdursn/hospital-case/pkg/utils"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(req *dto.RegisterRequest) (*models.Authority, error)
	Login(req *dto.LoginRequest, cfg *config.Config) (*dto.LoginResponse, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error)
	ResetPassword(req *dto.ResetPasswordRequest) error
}

type authUsecase struct {
	authRepo repository.AuthRepository
	redis    *redis.Client
}

func NewAuthUsecase(r repository.AuthRepository, redis *redis.Client) AuthUsecase {
	return &authUsecase{
		authRepo: r,
		redis:    redis,
	}
}

func (u *authUsecase) Register(req *dto.RegisterRequest) (*models.Authority, error) {
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

func (u *authUsecase) Login(req *dto.LoginRequest, cfg *config.Config) (*dto.LoginResponse, error) {
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
	tokenPair, err := utils.GenerateTokenPair(authority.ID, authority.HospitalID, authority.Role, cfg)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

func (u *authUsecase) ForgotPassword(req *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, error) {
	_, err := u.authRepo.GetAuthorityByPhone(req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	code := utils.GenerateResetCode()

	ctx := context.Background()
	if err := u.redis.Set(ctx, "reset_code:"+req.Phone, code, 5*time.Minute).Err(); err != nil {
		return nil, err
	}

	return &dto.ForgotPasswordResponse{Code: code}, nil
}

func (u *authUsecase) ResetPassword(req *dto.ResetPasswordRequest) error {
	if req.NewPassword != req.RepeatNewPassword {
		return errors.New("passwords do not match")
	}

	ctx := context.Background()
	storedCode, err := u.redis.Get(ctx, "reset_code:"+req.Phone).Result()
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

	_ = u.redis.Del(ctx, "reset_code:"+req.Phone).Err()
	return nil
}
