package splitwise

import (
	"bytes"
	"context"
	"encoding/json"
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

func handleError(resp *http.Response) error {
	statusCode := resp.StatusCode
	message := "Unknown"

	rawBody, err := io.ReadAll(resp.Body)
	if err == nil {
		msg, err := getErrorMessage(rawBody)
		if err == nil {
			message = msg
		}
	}

	return fmt.Errorf("[%d] %s", statusCode, message)
}

func paramsToJsonBytesReader(params map[string]string) (io.Reader, error) {
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

func (c *Client) do(ctx context.Context, method string, path string, queryParams url.Values, bodyParams map[string]string) (*http.Response, error) {
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
