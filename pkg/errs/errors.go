package errs

var (
	ErrInvalidCredentials = NewAppError("ERR_INVALID_CREDENTIALS", "Invalid credentials", 401, nil)
	ErrUserNotFound       = NewAppError("ERR_USER_NOT_FOUND", "User not found", 404, nil)
	ErrHospitalExists     = NewAppError("ERR_HOSPITAL_EXISTS", "Hospital already exists", 409, nil)
	ErrAuthorityExists    = NewAppError("ERR_AUTHORITY_EXISTS", "Authority already exists", 409, nil)
	ErrPasswordsMismatch  = NewAppError("ERR_PASSWORDS_MISMATCH", "Passwords do not match", 400, nil)
	ErrInvalidResetCode   = NewAppError("ERR_INVALID_RESET_CODE", "Invalid or expired reset code", 400, nil)
	ErrInternal           = NewAppError("ERR_INTERNAL", "Internal server error", 500, nil)
)
