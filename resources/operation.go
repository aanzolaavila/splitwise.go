package resources

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Stringer interface {
	String() string
}

type ApiClient interface {
	Get(endpoint string, queryParams url.Values, bodyParams url.Values) (*http.Response, error)
	Post(endpoint string, queryParams url.Values, bodyParams url.Values) (*http.Response, error)
}

type ErrorResponse struct {
	ErrCode int16
	Message string
}

func (e *ErrorResponse) StatusText() string {
	s := http.StatusText(int(e.ErrCode))
	if s == "" {
		return "Unknown"
	}

	return s
}

type errorMap struct {
	Error  string `json:"error,omitempty"`
	Errors struct {
		Base []string `json:"base"`
	} `json:"errors,omitempty"`
}

func (e *ErrorResponse) UnmarshalJSON(data []byte) error {
	var msgMap errorMap
	if err := json.Unmarshal(data, &msgMap); err != nil {
		return err
	}

	if msgMap.Error != "" {
		e.Message = msgMap.Error
	} else {
		if len(msgMap.Errors.Base) > 0 {
			e.Message = strings.Join(msgMap.Errors.Base, ", ")
		}
	}

	return nil
}

type Response[T any] struct {
	Result T
	Error  ErrorResponse
}

type RequiredParams interface {
	Set(key string, value string) error
}

type Params struct {
	Required RequiredParams
	Optional map[string]string
}

type Operation[T any] interface {
	Do(ApiClient, Params) (Response[T], error)
}
