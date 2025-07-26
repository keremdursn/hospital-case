package dto

import "time"

type RegisterRequest struct {
	HospitalName   string `json:"hospital_name" validate:"required,min=2,max=100"`
	TaxNumber      string `json:"tax_number" validate:"required,len=10"`
	HospitalEmail  string `json:"hospital_email" validate:"required,email"`
	HospitalPhone  string `json:"hospital_phone" validate:"required,phone"`
	Address        string `json:"address" validate:"required,min=10,max=200"`
	CityID         uint   `json:"city_id" validate:"required,gt=0"`
	DistrictID     uint   `json:"district_id" validate:"required,gt=0"`
	AuthorityFName string `json:"authority_fname" validate:"required,min=2,max=50"`
	AuthorityLName string `json:"authority_lname" validate:"required,min=2,max=50"`
	AuthorityTC    string `json:"authority_tc" validate:"required,tc"`
	AuthorityEmail string `json:"authority_email" validate:"required,email"`
	AuthorityPhone string `json:"authority_phone" validate:"required,phone"`
	Password       string `json:"password" validate:"required,password"`
}

type LoginRequest struct {
	Credential string
	Password   string
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
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

type AuthorityResponse struct {
	ID         uint       `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	TC         string     `json:"tc"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	Role       string     `json:"role"`
	HospitalID uint       `json:"hospital_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}
