package dto

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
