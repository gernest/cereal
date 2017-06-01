package messages

import "net/http"

const (
	BadJSON           = "Failed to Parse JSON"
	InvalidJSON       = "Body must be a valid json object"
	FailedValidation  = "Failed Validation"
	CodeMissing       = "missing"
	CodeMissingField  = "missingField"
	CodeInvalid       = "invalid"
	CodeAlreadyExists = "alreadyExists"
)

// Message is the error response sent to the client
type Message struct {
	Message string  `json:"message"`
	Errors  []Error `json:"errors,omitempty"`
}

// Error represents error details sent to the client
type Error struct {
	Resource string `json:"resource"`
	Field    string `json:"field"`
	Code     string `json:"code"`
}

func OK() interface{} {
	data := make(map[string]interface{})
	data["status"] = http.StatusText(http.StatusOK)
	return data
}
