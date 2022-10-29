package splitwise

import (
	"net/http"
	"net/http/httptest"
)

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
