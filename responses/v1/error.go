package v1

// Error ...
type Error struct {
	Error            interface{}       `json:"error"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

// ValidationError ...
type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Error string `json:"error"`
}
