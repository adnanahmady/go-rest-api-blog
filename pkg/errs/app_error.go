package errs

import (
	"encoding/json"
	"net/http"
)

var _ error = (*AppError)(nil)

type AppError struct {
	Code    int `json:"-"`
	Message string `json:"message"`
	Errors  Errors `json:"errors"`
}

type Errors map[string]string

func (e *AppError) Error() string {
	res, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(res)
}

func NewValidationError(errors Errors) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Validation error",
		Errors:  errors,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
	}
}