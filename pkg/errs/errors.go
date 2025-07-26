package errs

var (
	// Authentication errors
	ErrInvalidCredentials   = NewAppError("ERR_INVALID_CREDENTIALS", "Invalid credentials", 401, nil)
	ErrUserNotFound         = NewAppError("ERR_USER_NOT_FOUND", "User not found", 404, nil)
	ErrHospitalExists       = NewAppError("ERR_HOSPITAL_EXISTS", "Hospital already exists", 409, nil)
	ErrAuthorityExists      = NewAppError("ERR_AUTHORITY_EXISTS", "Authority already exists", 409, nil)
	ErrPasswordsMismatch    = NewAppError("ERR_PASSWORDS_MISMATCH", "Passwords do not match", 400, nil)
	ErrInvalidResetCode     = NewAppError("ERR_INVALID_RESET_CODE", "Invalid or expired reset code", 400, nil)
	ErrInvalidRefreshToken  = NewAppError("ERR_INVALID_REFRESH_TOKEN", "Invalid refresh token", 401, nil)
	ErrRefreshTokenRequired = NewAppError("ERR_REFRESH_TOKEN_REQUIRED", "Refresh token is required", 400, nil)

	// Validation errors
	ErrInvalidJSON       = NewAppError("ERR_INVALID_JSON", "Cannot parse JSON", 400, nil)
	ErrInvalidUserID     = NewAppError("ERR_INVALID_USER_ID", "Invalid user ID", 400, nil)
	ErrInvalidStaffID    = NewAppError("ERR_INVALID_STAFF_ID", "Invalid staff ID", 400, nil)
	ErrInvalidJobGroupID = NewAppError("ERR_INVALID_JOB_GROUP_ID", "Invalid job_group_id", 400, nil)
	ErrInvalidCityID     = NewAppError("ERR_INVALID_CITY_ID", "Invalid city_id", 400, nil)
	ErrInvalidID         = NewAppError("ERR_INVALID_ID", "Invalid id", 400, nil)
	ErrValidationFailed  = NewAppError("ERR_VALIDATION_FAILED", "Validation failed", 400, nil)

	// Authorization errors
	ErrUnauthorized = NewAppError("ERR_UNAUTHORIZED", "Unauthorized", 401, nil)
	ErrForbidden    = NewAppError("ERR_FORBIDDEN", "Forbidden", 403, nil)

	// Resource errors
	ErrNotFound = NewAppError("ERR_NOT_FOUND", "Resource not found", 404, nil)
	ErrConflict = NewAppError("ERR_CONFLICT", "Resource conflict", 409, nil)

	// Server errors
	ErrInternal = NewAppError("ERR_INTERNAL", "Internal server error", 500, nil)
	ErrDatabase = NewAppError("ERR_DATABASE", "Database error", 500, nil)

	// Rate limiting
	ErrTooManyRequests = NewAppError("ERR_TOO_MANY_REQUESTS", "Too many requests", 429, nil)
)
