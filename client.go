package splitwise

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseUrl = "https://secure.splitwise.com"
const DefaultApiVersionPath = "/api/v3.0"

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	HttpClient     HttpClient
	BaseUrl        string
	ApiVersionPath string
	Token          string
}

func defaultHttpClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

func (c *Client) baseUrl() string {
	if c.BaseUrl == "" {
		return DefaultBaseUrl
	}

	return c.BaseUrl
}

func (c *Client) apiVersionPath() string {
	if c.ApiVersionPath == "" {
		return DefaultApiVersionPath
	}

	return c.ApiVersionPath
}

func (c *Client) httpClient() HttpClient {
	if c.HttpClient == nil {
		c.HttpClient = defaultHttpClient()
	}

	return c.HttpClient
}

func (c *Client) addAuthorizationHeader(req *http.Request) {
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
}

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

	errsArray, ok := errsValue.([]string)
	if !ok {
		return nil
	}

	var errs []error
	for _, errStr := range errsArray {
		err := errors.New(errStr)
		errs = append(errs, err)
	}

	return errs
}

func handleStatusOkErrorResponse(res *http.Response) error {
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error response is not 200")
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var m map[string]interface{}
	err = json.Unmarshal(rawBody, &m)
	if err != nil {
		return err
	}

	successValue, ok := m["success"]
	if !ok {
		return fmt.Errorf("unknown error status")
	}

	successStatus, ok := successValue.(bool)
	if !ok {
		return fmt.Errorf("unexpected success response: %v", successValue)
	}

	if successStatus {
		return nil
	}

	respErrors := extractErrorsFromMap(m)

	if len(respErrors) == 1 {
		return fmt.Errorf("%v", respErrors[0])
	}

	if len(respErrors) > 1 {
		return fmt.Errorf("got multiple errors: %v", respErrors)
	}

	return fmt.Errorf("unsuccessful but unknown causes")
}

func paramsToJsonBytesReader(params map[string]interface{}) (io.Reader, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

func (c *Client) do(ctx context.Context, method string, path string, queryParams url.Values, bodyParams map[string]interface{}) (*http.Response, error) {
	path = c.baseUrl() + c.apiVersionPath() + path

	endpointUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	endpointUrl.RawQuery = queryParams.Encode()

	bodyReader, err := paramsToJsonBytesReader(bodyParams)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpointUrl.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	c.addAuthorizationHeader(req)

	req.Header.Add("Content-Type", "application/json")

	return c.httpClient().Do(req)
}
