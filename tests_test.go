package splitwise

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type readCloserStub struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (r readCloserStub) Read(p []byte) (n int, err error) {
	if r.ReadFunc == nil {
		panic("read func is nil")
	}

	return r.ReadFunc(p)
}

func (r readCloserStub) Close() error {
	if r.CloseFunc == nil {
		return nil
	}

	return r.CloseFunc()
}

type httpClientStub struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (c httpClientStub) Do(r *http.Request) (*http.Response, error) {
	if c.DoFunc == nil {
		panic("mocked function is nil")
	}

	return c.DoFunc(r)
}

func testClient(t *testing.T, statusCode int, method, response string) (_ Client, cancel func()) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if method != "" {
				assert.Equal(t, method, r.Method)
			}

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

func testClientThatFailsTestIfHttpIsCalled(t *testing.T) Client {
	stubClient := httpClientStub{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			t.Fatal("the stub http client should not have been called")
			return nil, fmt.Errorf("failed test")
		},
	}

	return Client{
		HttpClient: stubClient,
		BaseUrl:    "test.invalid.host",
	}
}

func testClientWithFaultyResponse() (Client, error) {
	expectedError := errors.New("this error is expected")

	stubClient := httpClientStub{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, expectedError
		},
	}

	return Client{
		HttpClient: stubClient,
		BaseUrl:    "test.invalid.host",
	}, expectedError
}

func testClientWithFaultyResponseBody(t *testing.T, statusCode int) (Client, error, func()) {
	expectedError := errors.New("this error is expected")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))

	httpClient := server.Client()
	url := server.URL

	stubHttpClient := httpClientStub{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			res, err := httpClient.Do(r)
			if err != nil {
				t.Fatalf("connection to mocked http server failed: %v", err)
			}

			reader := readCloserStub{
				ReadFunc: func(p []byte) (n int, err error) {
					return 0, expectedError
				},
			}

			res.Body = reader
			return res, nil
		},
	}

	client := Client{
		HttpClient: stubHttpClient,
		BaseUrl:    url,
	}

	return client, expectedError, func() {
		defer server.Close()
	}
}

func doBasicErrorChecks(t *testing.T, f func(Client) error) {
	if f == nil {
		panic("callback function is nil")
	}

	tests := []struct {
		Name string
		Func func(*testing.T, func(Client) error)
	}{
		{
			Name: "doFaultyClientTest",
			Func: doFaultyClientTest,
		},
		{
			Name: "doFaultyResponseBodyTest",
			Func: doFaultyResponseBodyTest,
		},
		{
			Name: "doErrorResponseTests",
			Func: doErrorResponseTests,
		},
		{
			Name: "doInvalidJsonResponseErrorTest",
			Func: doInvalidJsonResponseErrorTest,
		},
	}

	for _, tf := range tests {
		fnName := tf.Name
		testName := t.Name()
		t.Run(fmt.Sprintf("%s:%v", testName, fnName), func(t *testing.T) {
			tf.Func(t, f)
		})
	}
}

func doFaultyClientTest(t *testing.T, f func(Client) error) {
	client, expectedErr := testClientWithFaultyResponse()

	err := f(client)
	assert.ErrorIs(t, err, expectedErr)
}

func doFaultyResponseBodyTest(t *testing.T, f func(Client) error) {
	client, expectedErr, cancel := testClientWithFaultyResponseBody(t, http.StatusOK)
	defer cancel()

	err := f(client)
	assert.ErrorIs(t, err, expectedErr)
}

const (
	emptyJsonResponse  = `{}`
	successResponse    = `{ "success": true }`
	noSuccessResponse  = `{ "success": false }`
	errorResponse      = `{ "error": "one error" }`
	errorsResponse     = `{ "errors": ["err 1", "err 2"] }`
	errorsBaseResponse = `{ "errors": { "base": ["err 1", "err 2"] } }`
)

func doErrorResponseTests(t *testing.T, f func(Client) error) {
	checks := []struct {
		StatusCode    int
		Body          string
		ExpectedError error
	}{
		{http.StatusOK, successResponse, nil},
		{http.StatusOK, emptyJsonResponse, nil},
		{http.StatusOK, noSuccessResponse, ErrUnsuccessful},
		{http.StatusOK, errorResponse, ErrUnsuccessful},
		{http.StatusOK, errorsResponse, ErrUnsuccessful},
		{http.StatusOK, errorsBaseResponse, ErrUnsuccessful},
		{http.StatusNotFound, emptyJsonResponse, ErrNotFound},
		{http.StatusUnauthorized, emptyJsonResponse, ErrNotLoggedIn},
		{http.StatusForbidden, emptyJsonResponse, ErrForbidden},
		{http.StatusInternalServerError, emptyJsonResponse, ErrSplitwiseServer},
		{http.StatusBadRequest, emptyJsonResponse, ErrBadRequest},
	}

	for i, c := range checks {
		t.Run(fmt.Sprintf("%s:T%d", t.Name(), i), func(t *testing.T) {
			err := doErrorResponseTest(t, f, c.StatusCode, c.Body)
			if c.ExpectedError == nil {
				assert.NoErrorf(t, err, "was NOT expecting error on test #%d: %v", i, err)
			} else {
				assert.ErrorIsf(t, err, c.ExpectedError, "was expecting error [%v] on test #%d, got [%v] instead", c.ExpectedError, i, err)
			}
		})
	}
}

func doErrorResponseTest(t *testing.T, f func(Client) error, statusCode int, body string) error {
	client, cancel := testClient(t, statusCode, "", body)
	defer cancel()

	return f(client)
}

func doInvalidJsonResponseErrorTest(t *testing.T, f func(Client) error) {
	const invalidJson = `{ invalid }`
	client, cancel := testClient(t, http.StatusOK, "", invalidJson)
	defer cancel()

	err := f(client)

	var syntaxErr *json.SyntaxError
	assert.ErrorAs(t, err, &syntaxErr)
}
