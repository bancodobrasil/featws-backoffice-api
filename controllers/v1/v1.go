package v1

import (
	responses "github.com/bancodobrasil/featws-api/responses/v1"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// validatePayload validates a payload using a validator library and returns any validation errors as a
// custom error response.
func validatePayload(payload interface{}) *responses.Error {
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
