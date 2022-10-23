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
	case 200:
		return "Request was unsuccessful"
	case 400:
		return "Bad request"
	case 401:
		return "Invalid API request: you are not logged in"
	case 403:
		return "Invalid API request: you do not have permission to perform that action"
	case 404:
		return "Invalid API Request: record not found"
	default:
		return "Unknown error"
	}
}

func (e SplitwiseError) Error() string {
	return fmt.Sprintf("[%d] %s", e, e.String())
}

func (e SplitwiseError) Is(target error) bool {
	if ne, ok := target.(SplitwiseError); ok {
		return e == ne
	}

	return false
}

const (
	// Splitwise responds 200s on some erroneous requests
	ErrUnsuccessful = SplitwiseError(200)
	ErrBadRequest   = SplitwiseError(400)
	ErrNotLoggedIn  = SplitwiseError(401)
	ErrUnauthorized = SplitwiseError(403)
	ErrNotFound     = SplitwiseError(404)
)

type errorMap struct {
	Error  string `json:"error,omitempty"`
	Errors struct {
		Base []string `json:"base"`
	} `json:"errors,omitempty"`
}

func getErrorMessage(data []byte) (string, error) {
	var msgMap errorMap
	if err := json.Unmarshal(data, &msgMap); err != nil {
		return "", err
	}

	if msgMap.Error != "" {
		return msgMap.Error, nil
	} else {
		if len(msgMap.Errors.Base) > 0 {
			return strings.Join(msgMap.Errors.Base, ", "), nil
		}
	}

	return "", nil
}

func handleResponseError(res *http.Response) error {
	statusCode := res.StatusCode
	message := "Unknown"

	rawBody, err := io.ReadAll(res.Body)
	if err == nil {
		msg, err := getErrorMessage(rawBody)
		if err == nil {
			message = msg
		}
	}
	defer res.Body.Close()

	return fmt.Errorf("[%d] %s", statusCode, message)
}

func extractErrorsFromMap(m map[string]interface{}) []error {
	errsValue, ok := m["errors"]
	if !ok {
		return nil
	}

	errsArray, ok := errsValue.([]interface{})
	if !ok {
		baseValue, ok := errsValue.(map[string]interface{})
		if !ok {
			return nil
		}

		base, ok := baseValue["base"]
		if !ok {
			return nil
		}

		errsArray, ok = base.([]interface{})
		if !ok {
			return nil
		}

	}

	var strSlice []string
	for _, e := range errsArray {
		err, ok := e.(string)
		if ok {
			strSlice = append(strSlice, err)
		}
	}

	var errs []error
	for _, errStr := range strSlice {
		err := errors.New(errStr)
		errs = append(errs, err)
	}

	return errs
}

func handleStatusOkErrorResponse(res *http.Response, body []byte) error {
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error response is not 200")
	}

	var rawBody []byte
	var err error
	if body == nil {
		rawBody, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	} else {
		rawBody = body
	}

	var m map[string]interface{}
	err = json.Unmarshal(rawBody, &m)
	if err != nil {
		return err
	}

	respErrors := extractErrorsFromMap(m)

	if len(respErrors) == 1 {
		return fmt.Errorf("%w", respErrors[0])
	}

	if len(respErrors) > 1 {
		return fmt.Errorf("got multiple errors: %+v", respErrors)
	}

	var successStatus bool = true
	successValue, ok := m["success"]
	if ok {
		successStatus, ok = successValue.(bool)
		if !ok {
			return fmt.Errorf("unexpected success response: %v", successValue)
		}
	}

	if successStatus {
		return nil
	} else {
		return fmt.Errorf("unsuccessful with unknown causes")
	}
}
