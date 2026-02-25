package stadia

import (
	"encoding/json"
	"fmt"
	"io"

	"resty.dev/v3"
)

// Unfortunately, Stadia Maps is inconsistent in how it handles errors.
// We therefore have to have a function that handles all the different JSON
// error variations.
func parseError(resp *resty.Response) error {
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading all body: %w", err)
	}
	var server_error serverError
	err = json.Unmarshal(content, &server_error)
	if err == nil {
		return newAPIError(resp.StatusCode(), server_error.Error.Reason)
	}

	// At this point we've exhausted all of our options, so just pass the JSON through
	return newAPIError(resp.StatusCode(), string(content))
}

type apiError struct {
	Message string
	Status  int
}

func newAPIError(status int, msg string) apiError {
	return apiError{
		Message: msg,
		Status:  status,
	}
}
func (e apiError) Error() string {
	return e.Message
}

type Error struct {
	ErrorMessage string   `json:"error"`
	Errors       []string `json:"errors"`
}

func (e *Error) Error() string {
	return e.ErrorMessage
}

/*
Got this when I managed to bork the server

	{
	    "error": {
	        "reason": "Internal Server Error"
	    },
	    "status": 500
	}
*/
type errorWithReason struct {
	Reason string `json:"reason"`
}
type serverError struct {
	Error  errorWithReason `json:"error"`
	Status int             `json:"status"`
}

/*
	if len(result.Geocode.Errors) > 0 {
		joined := strings.Join(result.Geocode.Errors, ", ")
		return nil, fmt.Errorf("structured geocoding failure: %d '%s'", resp.StatusCode(), joined)
	} else if result.Geocode.Error != "" {
		return nil, fmt.Errorf("structured geocoding failure: %d '%s'", resp.StatusCode(), result.Geocode.Error)
	} else {
		return nil, fmt.Errorf("structured geocoding failure: %d", resp.StatusCode())
	}
*/
