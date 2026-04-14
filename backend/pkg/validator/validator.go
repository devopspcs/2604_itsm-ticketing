package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validate = validator.New()

func Validate(s interface{}) []FieldError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	var errs []FieldError
	for _, e := range err.(validator.ValidationErrors) {
		errs = append(errs, FieldError{
			Field:   e.Field(),
			Message: fmt.Sprintf("failed on '%s' validation", e.Tag()),
		})
	}
	return errs
}
