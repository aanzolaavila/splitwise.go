package splitwise

import (
	"net/http"
	"net/http/httptest"
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
