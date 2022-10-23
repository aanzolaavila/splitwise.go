package splitwise

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const DefaultBaseUrl = "https://secure.splitwise.com"
const DefaultApiVersionPath = "/api/v3.0"

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	Logger         *log.Logger
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

func (c *Client) logger() *log.Logger {
	if c.Logger == nil {
		return log.Default()
	}

	return c.Logger
}

func (c *Client) addAuthorizationHeader(req *http.Request) {
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
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

	c.logger().Printf("%s %s\n", method, endpointUrl.String())

	return c.httpClient().Do(req)
}
