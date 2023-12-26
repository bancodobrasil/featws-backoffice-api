package v1

import (
	"unicode"

	responses "github.com/bancodobrasil/featws-api/responses/v1"
	"github.com/go-playground/validator/v10"
)

// validatePayload validates a payload using a validator library and returns any validation errors as a
// custom error response.
func validatePayload(payload interface{}) *responses.Error {
	validate := validator.New()
	validate.RegisterValidation("doesNotStartWithDigit", validateDoesNotStartWithDigit)
	validate.RegisterValidation("mustContainLetter", validateMustContainLetter)
	if validationErr := validate.Struct(payload); validationErr != nil {
		err2, ok := validationErr.(validator.ValidationErrors)

		if !ok {
			return &responses.Error{
				Error: validationErr.Error(),
			}
		}
		errors := make([]responses.ValidationError, 0)

		for _, field := range err2 {
			errors = append(errors, responses.ValidationError{
				Field: field.Field(),
				Tag:   field.Tag(),
				Error: field.Error(),
			})
		}

		return &responses.Error{
			ValidationErrors: errors,
		}
	}

	return nil
}

func validateDoesNotStartWithDigit(fl validator.FieldLevel) bool {
	if unicode.IsDigit(rune(fl.Field().String()[0])) {
		return false
	}
	return true
}

func validateMustContainLetter(fl validator.FieldLevel) bool {
	for _, r := range fl.Field().String() {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}
