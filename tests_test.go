package splitwise

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockReadCloser struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (r mockReadCloser) Read(p []byte) (n int, err error) {
	if r.ReadFunc == nil {
		panic("read func is nil")
	}

	return r.ReadFunc(p)
}

func (r mockReadCloser) Close() error {
	if r.CloseFunc == nil {
		return nil
	}

	return r.CloseFunc()
}

type mockHttpClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (c mockHttpClient) Do(r *http.Request) (*http.Response, error) {
	if c.DoFunc == nil {
		panic("mocked function is nil")
	}

	return c.DoFunc(r)
}

func testClient(statusCode int, response string) (_ Client, cancel func()) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte(response))
		},
		))

	client := Client{
		HttpClient: server.Client(),
		BaseUrl:    server.URL,
	}

	return client, func() {
		defer server.Close()
	}
}

func testClientWithHandler(handler http.HandlerFunc) (_ Client, cancel func()) {
	server := httptest.NewServer(handler)

	client := Client{
		HttpClient: server.Client(),
		BaseUrl:    server.URL,
	}

	return client, func() {
		defer server.Close()
	}
}

func testClientWithFaultyResponse() (Client, error) {
	expectedError := errors.New("this error is expected")

	mockClient := mockHttpClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, expectedError
		},
	}

	return Client{
		HttpClient: mockClient,
	}, expectedError
}

func testClientWithFaultyResponseBody(t *testing.T, statusCode int) (Client, error, func()) {
	expectedError := errors.New("this error is expected")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))

	httpClient := server.Client()
	url := server.URL

	mockHttpClient := mockHttpClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			res, err := httpClient.Do(r)
			if err != nil {
				t.Fatalf("connection to mocked http server failed: %v", err)
			}

			reader := mockReadCloser{
				ReadFunc: func(p []byte) (n int, err error) {
					return 0, expectedError
				},
			}

			res.Body = reader
			return res, nil
		},
	}

	client := Client{
		HttpClient: mockHttpClient,
		BaseUrl:    url,
	}

	return client, expectedError, func() {
		defer server.Close()
	}
}
