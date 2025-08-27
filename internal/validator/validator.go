package validator

import "github.com/go-playground/validator"

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (v *CustomValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
