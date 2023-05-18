package v1

// FindResult is the result of a find request.
//
// Property:
// - Count: is a field of the FindResult struct that represents the number of occurrences found during a search operation. It is of type int64, which means it can hold integervalues up to 64 bits in size. The `json:"count,omitempty"` tag is used to specify the field name in JSON.
type FindResult struct {
	Count int64 `json:"count,omitempty"`
}
