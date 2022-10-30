package splitwise

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
)

const currentUser200Resp = `
{
  "user": {
    "id": 1,
    "first_name": "Ada",
    "last_name": "Lovelace",
    "email": "ada@example.com",
    "registration_status": "confirmed",
    "picture": {
      "small": "string",
      "medium": "string",
      "large": "string"
    },
    "notifications_read": "2017-06-02T20:21:57Z",
    "notifications_count": 12,
    "notifications": {
      "added_as_friend": true
    },
    "default_currency": "USD",
    "locale": "en"
  }
}
`

func Test_GetCurrentUser_200Response(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, currentUser200Resp)
	defer cancel()

	ctx := context.Background()

	user, err := client.GetCurrentUser(ctx)
	assert.NoError(t, err)

	assert.Equal(t, resources.UserID(1), user.ID)
	assert.Equal(t, "Ada", user.FirstName)
	assert.Equal(t, "Lovelace", user.LastName)
	assert.Equal(t, "Lovelace", user.LastName)
	assert.Equal(t, "ada@example.com", user.Email)
}

const currentUser401Resp = `
{
  "error": "Invalid API request: you are not logged in"
}
`

func Test_GetCurrentUser_401Response(t *testing.T) {
	client, cancel := testClient(t, http.StatusUnauthorized, http.MethodGet, currentUser401Resp)
	defer cancel()

	ctx := context.Background()

	_, err := client.GetCurrentUser(ctx)
	assert.ErrorIs(t, err, ErrNotLoggedIn)
}

func Test_GetCurrentUser_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.GetCurrentUser(ctx)

		return err
	}

	doBasicErrorChecks(t, f)
}

// ---

const getUser200Resp = `
{
  "user": {
    "id": 0,
    "first_name": "Ada",
    "last_name": "Lovelace",
    "email": "ada@example.com",
    "registration_status": "confirmed",
    "picture": {
      "small": "string",
      "medium": "string",
      "large": "string"
    }
  }
}
`

func Test_GetUser_200Response(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, getUser200Resp)
	defer cancel()

	ctx := context.Background()

	user, err := client.GetUser(ctx, 0)
	assert.NoError(t, err)

	assert.Equal(t, resources.UserID(0), user.ID)
	assert.Equal(t, "Ada", user.FirstName)
	assert.Equal(t, "Lovelace", user.LastName)
	assert.Equal(t, "ada@example.com", user.Email)
}

func Test_GetUser_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.GetUser(ctx, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}

// ---

func Test_UpdateUser(t *testing.T) {
	const userID = resources.UserID(15)

	client, cancel := testClientWithHandler(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		input := struct {
			Firstname       string `json:"first_name"`
			Lastname        string `json:"last_name"`
			Email           string `json:"email"`
			Password        string `json:"password"`
			Locale          string `json:"locale"`
			DefaultCurrency string `json:"default_currency"`
		}{}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &input)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userRes := resources.User{
			ID:        userID,
			FirstName: input.Firstname,
			LastName:  input.Lastname,
			Email:     input.Email,
			Locale:    input.Locale,
		}

		cont := struct {
			User resources.User `json:"user"`
		}{
			User: userRes,
		}

		res, err := json.Marshal(cont)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(res)
	})
	defer cancel()

	ctx := context.Background()

	const firstname = "Test"
	const lastname = "Test Lastname"
	user, err := client.UpdateUser(ctx, int(userID), UserParams{
		UserFirstname: firstname,
		UserLastname:  lastname,
	})

	assert.NoError(t, err)

	assert.Equal(t, userID, user.ID)
	assert.Equal(t, firstname, user.FirstName)
	assert.Equal(t, lastname, user.LastName)
}

func Test_UpdateUser_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.UpdateUser(ctx, 0, nil)

		return err
	}

	doBasicErrorChecks(t, f)
}
