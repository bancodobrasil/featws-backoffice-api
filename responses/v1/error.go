package v1

// Error contains an error message and optional validation errors.
//
// Property:
//
//   - Error: an interface type that can hold any value, such as error messages or objects, to describe the occurred error. It allows for providing detailed information about the error to users or developers.
//
//   - ValidationErrors: is a slice of ValidationError structs that holds validation errors encountered during request processing. The omitempty tag indicates that it will only appear in the JSON response if there are actual validation errors.
type Error struct {
	Error            interface{}       `json:"error"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// ValidationError represents an error that occurred during validation, including the field,
// tag, and error message.
//
// Property:
//   - Field: is a string property that represents the name of the field that has a validation error.
//   - Tag: a string property representing the violated validation rule. For instance, if a required field is missing, the Tag property would be "required".
//   - Error: a string property that describes the error message related to the validation error. It provides details about what went wrong during the validation process.
type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Error string `json:"error"`
}
