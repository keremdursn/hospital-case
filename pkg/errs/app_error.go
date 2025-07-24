package errs

import (
	"encoding/json"
	"fmt"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
	Err        error  `json:"-"` // orijinal hata, response'a yazılmaz
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s | cause: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) MarshalJSON() ([]byte, error) {
	// JSON response'ta Err alanı gösterilmez
	type Alias AppError
	return json.Marshal(&struct{ *Alias }{Alias: (*Alias)(e)})
}

// NewAppError yeni bir AppError oluşturur
func NewAppError(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: status,
		Err:        err,
	}
}
