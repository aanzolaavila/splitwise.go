package splitwise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	ClientVersion = "0.1"
	BaseEndpoint  = "/api/v3.0"
)

type ApiClient interface {
	Do(*http.Request) (*http.Response, error)
}

type SplitwiseClient struct {
	Client ApiClient
}

func NewClient(client ApiClient) (*SplitwiseClient, error) {
	if client == nil {
		return nil, fmt.Errorf("client cannot be nil")
	}

	return &SplitwiseClient{
		Client: client,
	}, nil
}

func (c SplitwiseClient) do(method string, endpoint string, queryParams url.Values, bodyParams url.Values) (*http.Response, error) {
	endpoint = fmt.Sprintf("%s/%s", BaseEndpoint, endpoint)

	body, err := json.Marshal(bodyParams)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(method, endpoint, bodyReader)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = queryParams.Encode()

	return c.Client.Do(req)
}

func (c SplitwiseClient) Get(endpoint string, queryParams url.Values, bodyParams url.Values) (*http.Response, error) {
	return c.do(http.MethodGet, endpoint, queryParams, bodyParams)
}

func (c SplitwiseClient) Post(endpoint string, queryParams url.Values, bodyParams url.Values) (*http.Response, error) {
	return c.do(http.MethodPost, endpoint, queryParams, bodyParams)
}

// func main() {
// client, _ := NewTokenClient("SoyX4zHwLjZkWvFGYS6OsZhokaIMg6WQm1bh8hJ8", nil)
//
// req, err := http.NewRequest(http.MethodGet, DefaultBaseUrl+"/api/v3.0/get_current_user", nil)
// if err != nil {
// log.Fatalf("error: %v", err)
// }
//
// res, err := client.Do(req)
// if err != nil {
// log.Fatalf("error: %v", err)
// }
//
// body, err := io.ReadAll(res.Body)
// if err != nil {
// log.Fatalf("error: %v", err)
// }
// defer res.Body.Close()
//
// fmt.Printf("%s", string(body))
// }
