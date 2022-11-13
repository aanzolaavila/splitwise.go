package splitwise

import (
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
	if body == nil {
		var err error
		body, err = io.ReadAll(res.Body)
		if err != nil {
			if res.StatusCode == http.StatusOK {
				c.getLogger().Printf("Warning: could not read from response body, but response status code is %d", res.StatusCode)
				return nil
			}

			return SplitwiseError(res.StatusCode)
		}

		defer res.Body.Close()
	}

	if err := c.checkJson(body); err != nil {
		return err
	}

	if err := c.extractErrorsFromBody(body); err != nil {
		return fmt.Errorf("%w: %s", SplitwiseError(res.StatusCode), err.Error())
	}

	sv := c.extractSuccessValue(body)
	if sv != nil && !*sv {
		return fmt.Errorf("%w: there was no error in response", SplitwiseError(res.StatusCode))
	}

	if res.StatusCode != http.StatusOK {
		return SplitwiseError(res.StatusCode)
	}

	return nil
}

func (c *Client) checkJson(body []byte) error {
	if len(body) == 0 {
		return nil
	}

	var j interface{}
	return c.unmarshal()(body, &j)
}

func (c *Client) extractSuccessValue(body []byte) *bool {
	var s struct {
		Success *bool `json:"success"`
	}

	err := c.unmarshal()(body, &s)
	if err != nil {
		return nil
	}

	return s.Success
}

func (c *Client) extractErrorsFromBody(body []byte) error {
	if len(body) == 0 {
		return nil
	}

	var s []string

	e := c.extractSigleError(body)
	if e != "" {
		s = append(s, e)
	}

	errs := c.extractErrorList(body)
	s = append(s, errs...)

	errs = c.extractErrorsBase(body)
	s = append(s, errs...)

	errs = c.extractPropertyErrorStruct(body)
	s = append(s, errs...)

	if len(s) == 1 {
		return errors.New(s[0])
	}

	if len(s) > 1 {
		es := strings.Join(s, ", ")
		return fmt.Errorf("errors: [ %s ]", es)
	}

	return nil
}

func (c *Client) extractErrorList(body []byte) []string {
	e := struct {
		Errors []string `json:"errors"`
	}{}

	_ = c.unmarshal()(body, &e)

	return e.Errors
}

func (c *Client) extractErrorsBase(body []byte) []string {
	e := struct {
		Errors struct {
			Base []string `json:"base"`
		} `json:"errors"`
	}{}

	_ = c.unmarshal()(body, &e)

	return e.Errors.Base
}

func (c *Client) extractSigleError(body []byte) string {
	e := struct {
		Error string `json:"error"`
	}{}

	_ = c.unmarshal()(body, &e)

	return e.Error
}

func (c *Client) extractPropertyErrorStruct(body []byte) (es []string) {
	e := struct {
		Errors map[string][]string `json:"errors"`
	}{}

	_ = c.unmarshal()(body, &e)

	for k, v := range e.Errors {
		errs := strings.Join(v, ", ")
		err := fmt.Sprintf("property \"%s\" got errors: [ %s ]", k, errs)
		es = append(es, err)
	}

	return
}
