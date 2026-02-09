package client

import (
	"fmt"
	"net/http"
	"strings"
)

// APIError normalizes AWX API and transport failures into a consistent shape.
type APIError struct {
	Method     string
	URL        string
	StatusCode int
	Detail     string
	Body       string
	Retryable  bool
	Err        error
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}

	parts := []string{fmt.Sprintf("%s %s", strings.ToUpper(e.Method), e.URL)}
	if e.StatusCode > 0 {
		parts = append(parts, fmt.Sprintf("status=%d", e.StatusCode))
	}
	if e.Detail != "" {
		parts = append(parts, fmt.Sprintf("detail=%s", e.Detail))
	}
	if e.Err != nil {
		parts = append(parts, fmt.Sprintf("error=%s", e.Err.Error()))
	}
	return strings.Join(parts, " ")
}

func (e *APIError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *APIError) IsNotFound() bool {
	if e == nil {
		return false
	}
	return e.StatusCode == http.StatusNotFound
}
