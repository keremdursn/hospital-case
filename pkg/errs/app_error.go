package errs

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
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

// SendErrorResponse Fiber context'e standardize edilmiş error response gönderir
func SendErrorResponse(c *fiber.Ctx, appErr *AppError) error {
	return c.Status(appErr.StatusCode).JSON(appErr)
}

// SendErrorResponseWithDetails Detaylı error response gönderir
func SendErrorResponseWithDetails(c *fiber.Ctx, appErr *AppError, details interface{}) error {
	response := fiber.Map{
		"code":    appErr.Code,
		"message": appErr.Message,
		"status":  appErr.StatusCode,
	}

	if details != nil {
		response["details"] = details
	}

	return c.Status(appErr.StatusCode).JSON(response)
}

// HandleError Genel error'ı AppError'a çevirir ve response gönderir
func HandleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*AppError); ok {
		return SendErrorResponse(c, appErr)
	}

	// Genel error için internal server error döner
	return SendErrorResponse(c, ErrInternal)
}
