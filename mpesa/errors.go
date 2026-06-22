package mpesa

import "fmt"

// APIError is returned for non-2xx HTTP responses and known non-success M-Pesa
// response codes.
type APIError struct {
	StatusCode   int
	ResponseCode string
	Description  string
	Body         string
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.ResponseCode != "" || e.Description != "" {
		return fmt.Sprintf("mpesa api error: http=%d code=%s desc=%s", e.StatusCode, e.ResponseCode, e.Description)
	}
	return fmt.Sprintf("mpesa api error: http=%d body=%s", e.StatusCode, e.Body)
}
