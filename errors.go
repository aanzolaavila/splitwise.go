package splitwise

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type SplitwiseError uint16

func (e SplitwiseError) String() string {
	switch e {
	case ErrInvalidParameter:
		return "invalid parameter"
	case ErrUnsuccessful:
		return "request was unsuccessful"
	case ErrBadRequest:
		return "bad request"
	case ErrNotLoggedIn:
		return "invalid API request: you are not logged in"
	case ErrForbidden:
		return "invalid API request: you do not have permission to perform that action"
	case ErrNotFound:
		return "invalid API Request: record not found"
	case ErrSplitwiseServer:
		return "internal Server Error"
	default:
		return "unknown error"
	}
}

func (e SplitwiseError) Error() string {
	if e != ErrInvalidParameter {
		return fmt.Sprintf("status code %d: %s", e, e.String())
	}
	return e.String()
}

func (e SplitwiseError) Is(target error) bool {
	if ne, ok := target.(SplitwiseError); ok {
		return e == ne
	}

	return false
}

const (
	ErrInvalidParameter SplitwiseError = 0
	// Splitwise responds 200s on some erroneous requests
	ErrUnsuccessful    SplitwiseError = http.StatusOK
	ErrBadRequest      SplitwiseError = http.StatusBadRequest
	ErrNotLoggedIn     SplitwiseError = http.StatusUnauthorized
	ErrForbidden       SplitwiseError = http.StatusForbidden
	ErrNotFound        SplitwiseError = http.StatusNotFound
	ErrSplitwiseServer SplitwiseError = http.StatusInternalServerError
)

func (c *Client) getErrorFromResponse(res *http.Response, body []byte) error {
	var rawBody []byte = body
	if body == nil {
		var err error
		rawBody, err = io.ReadAll(res.Body)
		if err != nil {
			if res.StatusCode == http.StatusOK {
				c.getLogger().Printf("Warning: could not read from response body, but response status code is %d", res.StatusCode)
				return nil
			}

			return SplitwiseError(res.StatusCode)
		}

		defer res.Body.Close()
	}

	err := extractErrorsFromBody(rawBody)

	if err != nil {
		return fmt.Errorf("%w: %s", SplitwiseError(res.StatusCode), err.Error())
	}

	sv := extractSuccessValue(rawBody)
	if sv != nil && !*sv {
		return fmt.Errorf("%w: there was no error in response", SplitwiseError(res.StatusCode))
	}

	if res.StatusCode != http.StatusOK {
		return SplitwiseError(res.StatusCode)
	}

	return nil
}

type successMap struct {
	Success *bool `json:"success"`
}

func extractSuccessValue(body []byte) *bool {
	var s successMap
	err := json.Unmarshal(body, &s)
	if err != nil {
		return nil
	}

	return s.Success
}

type errorMap struct {
	Error  string `json:"error"`
	Errors struct {
		Base []string `json:"base"`
	} `json:"errors"`
}

type errorsListMap struct {
	Errors []string `json:"errors"`
}

func extractErrorsFromBody(body []byte) error {
	var errSlice errorsListMap
	var errMap errorMap

	_ = json.Unmarshal(body, &errMap)
	_ = json.Unmarshal(body, &errSlice)

	if errMap.Error != "" {
		return errors.New(errMap.Error)
	}

	s := errSlice.Errors
	s = append(s, errMap.Errors.Base...)

	if len(s) > 0 {
		errs := strings.Join(s, ", ")
		return fmt.Errorf("multiple errors: %s", errs)
	}

	return nil
}
