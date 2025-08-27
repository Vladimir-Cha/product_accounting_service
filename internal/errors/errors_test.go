package errors

import "testing"

func TestErrorWithDetails(t *testing.T) {
	err := ErrBadRequest.WithDetails(map[string]any{"field": "email"})
	if err.Details == nil {
		t.Error("Details not set")
	}
}
