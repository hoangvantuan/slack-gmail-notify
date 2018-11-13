package app

import (
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

// DefaultValidator is default validator
type DefaultValidator struct {
	validator *validator.Validate
}

// NewDefaultValidator return new default validator
func NewDefaultValidator() echo.Validator {
	return &DefaultValidator{validator: validator.New()}
}

// Validate validate
func (cv *DefaultValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
