package v1

// FindResult is the result of a find request.
type FindResult struct {
	Count int64 `json:"count,omitempty"`
}
