package splitwise

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	client, cancel := testClient(400, error400Response)
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
	client, cancel := testClient(401, error401Response)
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
	client, cancel := testClient(403, error403Response)
	defer cancel()

	ctx := context.Background()

	res, err := client.do(ctx, http.MethodGet, "/", nil, nil)
	assert.NoError(t, err)

	err = client.getErrorFromResponse(res, nil)
	assert.ErrorIs(t, err, ErrUnauthorized)
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
	client, cancel := testClient(404, error404Response)
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
	client, cancel := testClient(200, error200UnsuccessfulResponse)
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
	client, cancel := testClient(200, error200UnsuccessfulResponse)
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
	client, cancel := testClient(200, error200UnsuccessfulResponse)
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
	client, cancel := testClient(200, error200SuccessfulResponse)
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
	client, cancel := testClient(200, error200ErroneousSuccessfulResponse)
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
