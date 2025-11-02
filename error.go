package elevenlabs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error returned by the ElevenLabs API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Detail     string `json:"detail,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("elevenlabs api error (status %d): %s - %s", e.StatusCode, e.Message, e.Detail)
	}
	return fmt.Sprintf("elevenlabs api error (status %d): %s", e.StatusCode, e.Message)
}

// parseErrorResponse parses an error response from the API
func parseErrorResponse(resp *http.Response) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Message:    resp.Status,
	}

	// Try to read and parse the error body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apiErr
	}

	// Try to unmarshal into a structured error
	var errorResponse struct {
		Detail  string `json:"detail"`
		Message string `json:"message"`
	}

	if err := json.Unmarshal(body, &errorResponse); err == nil {
		if errorResponse.Detail != "" {
			apiErr.Detail = errorResponse.Detail
		}
		if errorResponse.Message != "" {
			apiErr.Message = errorResponse.Message
		}
	} else {
		// If JSON parsing fails, use the raw body as the detail
		apiErr.Detail = string(body)
	}

	return apiErr
}
