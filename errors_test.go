package splitwise

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Error400Response(t *testing.T) {
	const error400Response = `
{
  "errors": {
    "base": [
      "You cannot add unknown users to a group by user_id"
    ]
  }
}
`
	client, cancel := testClient(t, http.StatusBadRequest, http.MethodGet, error400Response)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func Test_Error401Response(t *testing.T) {
	const error401Response = `
{
  "error": "Invalid API request: you are not logged in"
}
`
	client, cancel := testClient(t, http.StatusUnauthorized, http.MethodGet, error401Response)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrNotLoggedIn)
}

func Test_Error403Response(t *testing.T) {
	const error403Response = `
{
  "errors": {
    "base": [
      "Invalid API request: you do not have permission to perform that action"
    ]
  }
}
`
	client, cancel := testClient(t, http.StatusForbidden, http.MethodGet, error403Response)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrForbidden)
}

func Test_Error404Response(t *testing.T) {
	const error404Response = `
{
  "errors": {
    "base": [
      "Invalid API Request: record not found"
    ]
  }
}
`
	client, cancel := testClient(t, http.StatusNotFound, http.MethodGet, error404Response)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrNotFound)
}

func Test_Error200NoSuccessResponse(t *testing.T) {
	const error200UnsuccessfulResponse = `
{
  "success": false,
  "errors": []
}
`
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, error200UnsuccessfulResponse)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrUnsuccessful)
}

func Test_Error200ErrorsSliceResponse(t *testing.T) {
	const e = "This is an error"
	const error200UnsuccessfulResponse = `
{
  "errors": ["This is an error"]
}
`
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, error200UnsuccessfulResponse)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrUnsuccessful)

	assert.True(t, strings.Contains(err.Error(), e))
}

func Test_Error200SingleErrorResponse(t *testing.T) {
	const e = "This is an error"
	const error200UnsuccessfulResponse = `
{
  "error": "This is an error"
}
`
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, error200UnsuccessfulResponse)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrUnsuccessful)
	assert.True(t, strings.Contains(err.Error(), e))
}

func Test_Error200SuccessResponse_ShouldNotFail(t *testing.T) {
	const error200SuccessfulResponse = `
{
  "success": true
}
`
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, error200SuccessfulResponse)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.NoError(t, err)
}

func Test_200Response_InvalidJsonShouldNotFail(t *testing.T) {
	const error200ErroneousSuccessfulResponse = `
{
  "success": true
`
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, error200ErroneousSuccessfulResponse)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.NoError(t, err)
}

func Test_200Response_ShouldNotFailIfInvalidBody(t *testing.T) {
	const successResponse = `
{
  "success": true
}
`
	client, _, cancel := testClientWithFaultyResponseBody(t, http.StatusOK)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.NoError(t, err)
}

func Test_4XXResponse_ShouldFailIfInvalidBody(t *testing.T) {
	const successResponse = `
{
  "success": true
}
`
	client, _, cancel := testClientWithFaultyResponseBody(t, http.StatusNotFound)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.Error(t, err)
	assert.ErrorIs(t, ErrNotFound, err)
}

func Test_200ResponseWithPropertyNamedStruct(t *testing.T) {
	const response = `
{
"errors": {
    "property1": [
      "string"
    ],
    "property2": [
      "string",
      "string"
    ]
  }
}
`
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, response)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrUnsuccessful)
}

func Test_IfResponseIsNot2XXButHasErrors_ItShouldAlsoIncludeThem(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		e1 = "Error1"
		e2 = "Error2"
	)
	const response = `
{
  "errors": ["Error1", "Error2"]
}
`
	client, cancel := testClient(t, http.StatusBadRequest, http.MethodGet, response)
	defer cancel()

	ctx := context.Background()
	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	require.NoError(err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(err, ErrBadRequest)
	assert.True(strings.Contains(err.Error(), e1))
	assert.True(strings.Contains(err.Error(), e2))
}
