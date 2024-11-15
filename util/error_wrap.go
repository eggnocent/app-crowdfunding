package util

import (
	"context"
	"fmt"
	"net/http"
)

// ErrorWrap is a custom error struct to wrap detailed error information
type ErrorWrap struct {
	OriginalError error
	Source        string
	Action        string
	Message       string
	Code          int
}

// NewErrorWrap creates a new wrapped error
func NewErrorWrap(originalError error, source, action string, ctx context.Context, message string, code int) *ErrorWrap {
	return &ErrorWrap{
		OriginalError: originalError,
		Source:        source,
		Action:        action,
		Message:       message,
		Code:          code,
	}
}

func (e *ErrorWrap) Error() string {
	return fmt.Sprintf("[%s:%s] %s - %v", e.Source, e.Action, e.Message, e.OriginalError)
}

// WriteErrorResponse writes the error response as JSON
func WriteErrorResponse(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}
