package splitwise

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Client_IfTokenIsDefinedItShouldBeInHeader(t *testing.T) {
	const testToken = "testtoken"

	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != fmt.Sprintf("Bearer %s", testToken) {
			assert.Failf(t, "authorization header is not correct", "authorization header is [%s], should be [Bearer %s]", authHeader, testToken)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer cancel()

	ctx := context.Background()

	client.Token = testToken

	res, err := client.do(ctx, http.MethodGet, "/generic", nil, nil)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func Test_Client_IfNoClientIsDefinedUseADefaultOne(t *testing.T) {
	wasCalled := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wasCalled = true

		w.Write([]byte("success"))
	}))
	defer server.Close()

	url := server.URL

	client := Client{
		BaseUrl: url,
	}

	ctx := context.Background()
	res, err := client.do(ctx, http.MethodGet, "", nil, nil)
	require.NoError(t, err)

	rawBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, "success", string(rawBody))
	assert.True(t, wasCalled)
}

func Test_Client_BaseURL(t *testing.T) {
	client := Client{}

	assert.Equal(t, DefaultBaseUrl, client.baseUrl())

	const url = "https://fakesite.com"
	client.BaseUrl = url

	assert.Equal(t, url, client.baseUrl())
}

func Test_Client_ApiVersionPath(t *testing.T) {
	client := Client{}

	assert.Equal(t, DefaultApiVersionPath, client.apiVersionPath())

	const path = "/path"
	client.ApiVersionPath = path

	assert.Equal(t, path, client.apiVersionPath())
}

func Test_JsonMarshal(t *testing.T) {
	client := Client{}

	f1ptr := reflect.ValueOf(json.Marshal)
	f2ptr := reflect.ValueOf(client.marshal())
	assert.Equal(t, f1ptr.Pointer(), f2ptr.Pointer())

	marshal := func(_ interface{}) ([]byte, error) {
		return nil, nil
	}
	client.JsonMarshaler = marshal

	f1ptr = reflect.ValueOf(marshal)
	f2ptr = reflect.ValueOf(client.marshal())
	assert.Equal(t, f1ptr.Pointer(), f2ptr.Pointer())
}

func Test_JsonUnmarshal(t *testing.T) {
	client := Client{}

	f1ptr := reflect.ValueOf(json.Unmarshal)
	f2ptr := reflect.ValueOf(client.unmarshal())
	assert.Equal(t, f1ptr.Pointer(), f2ptr.Pointer())

	unmarshal := func(_ []byte, _ interface{}) error {
		return nil
	}
	client.JsonUnmarshaler = unmarshal

	f1ptr = reflect.ValueOf(unmarshal)
	f2ptr = reflect.ValueOf(client.unmarshal())
	assert.Equal(t, f1ptr.Pointer(), f2ptr.Pointer())
}

func Test_Logger(t *testing.T) {
	client := Client{}

	assert.Equal(t, log.Default(), client.getLogger())

	var buf bytes.Buffer
	logger := log.New(io.Writer(&buf), "", 0)

	client.Logger = logger

	const data = "test"
	client.getLogger().Printf(data)

	assert.Equal(t, len(data)+1, buf.Len())
}

func Test_Client_Do_FailsIfPathIsWrong(t *testing.T) {
	client := Client{
		BaseUrl:        "https:// ////",
		ApiVersionPath: " ",
	}

	ctx := context.Background()
	_, err := client.do(ctx, http.MethodGet, "", nil, nil)
	require.Error(t, err)

	var parseErr *url.Error
	assert.ErrorAs(t, err, &parseErr)
}

func Test_Client_Do_shouldFailIfMarshalingFails(t *testing.T) {
	expectedErr := errors.New("this is expected")
	stubMarshal := func(v interface{}) ([]byte, error) {
		return nil, expectedErr
	}

	client := Client{
		JsonMarshaler: stubMarshal,
	}

	ctx := context.Background()
	_, err := client.do(ctx, http.MethodGet, "", nil, nil)
	require.Error(t, err)

	assert.ErrorIs(t, err, expectedErr)
}

func Test_Client_Do_shouldFailIfRequestCannotBeCreated(t *testing.T) {
	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("this should not have been reached")
	}))
	defer cancel()

	_, err := client.do(nil, http.MethodGet, "", nil, nil)
	assert.Error(t, err)

	_, err = client.do(context.Background(), "INVALID\n", "", nil, nil)
	assert.Error(t, err)
}
