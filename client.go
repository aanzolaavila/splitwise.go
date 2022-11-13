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

type jsonMarshaler func(interface{}) ([]byte, error)
type jsonUnmarshaler func([]byte, interface{}) error

type logger interface {
	Printf(string, ...interface{})
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	Logger          logger
	HttpClient      httpClient
	BaseUrl         string
	ApiVersionPath  string
	Token           string
	JsonMarshaler   jsonMarshaler
	JsonUnmarshaler jsonUnmarshaler
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

func (c *Client) marshal() jsonMarshaler {
	if c.JsonMarshaler == nil {
		return json.Marshal
	}

	return c.JsonMarshaler
}

func (c *Client) unmarshal() jsonUnmarshaler {
	if c.JsonUnmarshaler == nil {
		return json.Unmarshal
	}

	return c.JsonUnmarshaler
}

func (c *Client) getHttpClient() httpClient {
	if c.HttpClient == nil {
		c.HttpClient = defaultHttpClient()
	}

	return c.HttpClient
}

func (c *Client) getLogger() logger {
	if c.Logger == nil {
		return log.Default()
	}

	return c.Logger
}

func (c *Client) addRequiredHeaders(req *http.Request) {
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	req.Header.Add("Content-Type", "application/json")
}

func (c *Client) paramsToJsonBytesReader(params map[string]interface{}) (io.Reader, error) {
	body, err := c.marshal()(params)
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

	bodyReader, err := c.paramsToJsonBytesReader(bodyParams)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, endpointUrl.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	c.addRequiredHeaders(req)

	c.getLogger().Printf("%s %s\n", method, endpointUrl.String())

	return c.getHttpClient().Do(req)
}
