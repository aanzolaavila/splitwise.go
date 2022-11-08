package splitwise

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const getFriends200Response = `
{
  "friends": [
    {
      "id": 15,
      "first_name": "Ada",
      "last_name": "Lovelace",
      "email": "ada@example.com",
      "registration_status": "confirmed",
      "picture": {
        "small": "string",
        "medium": "string",
        "large": "string"
      },
      "groups": [
        {
          "group_id": 571,
          "balance": [
            {
              "currency_code": "USD",
              "amount": "414.5"
            }
          ]
        }
      ],
      "balance": [
        {
          "currency_code": "USD",
          "amount": "414.5"
        }
      ],
      "updated_at": "2019-08-24T14:15:22Z"
    }
  ]
}
`

func Test_GetFriends(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, getFriends200Response)
	defer cancel()

	ctx := context.Background()
	fs, err := client.GetFriends(ctx)
	require.NoError(t, err)

	require.Len(t, fs, 1)
	f := fs[0]
	assert.Equal(t, resources.FriendID(15), f.ID)
	assert.Equal(t, "Ada", f.FirstName)
	assert.Equal(t, "Lovelace", f.LastName)

	require.Len(t, f.Balance, 1)
	assert.Equal(t, "414.5", f.Balance[0].Amount)

	d := time.Date(2019, 8, 24, 14, 15, 22, 0, time.UTC)
	assert.True(t, d.Equal(f.UpdatedAt))
}

func Test_GetFriends_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.GetFriends(ctx)

		return err
	}

	doBasicErrorChecks(t, f)
}

const getFriend200Response = `
{
  "friend": {
    "id": 15,
    "first_name": "Ada",
    "last_name": "Lovelace",
    "email": "ada@example.com",
    "registration_status": "confirmed",
    "picture": {
      "small": "string",
      "medium": "string",
      "large": "string"
    },
    "groups": [
      {
        "group_id": 571,
        "balance": [
          {
            "currency_code": "USD",
            "amount": "414.5"
          }
        ]
      }
    ],
    "balance": [
      {
        "currency_code": "USD",
        "amount": "414.5"
      }
    ],
    "updated_at": "2019-08-24T14:15:22Z"
  }
}
`

func Test_GetFriend(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, getFriend200Response)
	defer cancel()

	ctx := context.Background()
	f, err := client.GetFriend(ctx, 0)
	require.NoError(t, err)

	assert.Equal(t, resources.FriendID(15), f.ID)
	assert.Equal(t, "Ada", f.FirstName)
	assert.Equal(t, "Lovelace", f.LastName)

	require.Len(t, f.Balance, 1)
	assert.Equal(t, "414.5", f.Balance[0].Amount)

	d := time.Date(2019, 8, 24, 14, 15, 22, 0, time.UTC)
	assert.True(t, d.Equal(f.UpdatedAt))

}

func Test_GetFriend_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.GetFriend(ctx, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_CreateFriend(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		friendEmail     = "test@email.com"
		friendFirstname = "Test"
		friendLastname  = "Testing"
		friendId        = resources.FriendID(150)
	)

	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := struct {
			Email     *string `json:"user_email"`
			Firstname *string `json:"user_first_name"`
			Lastname  *string `json:"user_last_name"`
		}{}

		rawBody, err := io.ReadAll(r.Body)
		require.NoError(err)

		err = json.Unmarshal(rawBody, &in)
		require.NoError(err)

		require.NotNil(in.Email)
		require.NotNil(in.Firstname)
		require.NotNil(in.Lastname)

		assert.Equal(friendEmail, *in.Email)
		assert.Equal(friendFirstname, *in.Firstname)
		assert.Equal(friendLastname, *in.Lastname)

		c := struct {
			Friend resources.Friend `json:"friend"`
		}{
			Friend: resources.Friend{
				ID:        friendId,
				Email:     friendEmail,
				FirstName: friendFirstname,
				LastName:  friendLastname,
			},
		}

		res, err := json.Marshal(c)
		require.NoError(err)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(res)
	}))
	defer cancel()

	ctx := context.Background()
	f, err := client.AddFriend(ctx, friendEmail, FriendParams{
		FriendFirstname: friendFirstname,
		FriendLastname:  friendLastname,
	})
	require.NoError(err)

	assert.Equal(friendId, f.ID)
	assert.Equal(friendEmail, f.Email)
	assert.Equal(friendFirstname, f.FirstName)
	assert.Equal(friendLastname, f.LastName)
}

func Test_AddFriend_EmailCannotEmpty(t *testing.T) {
	client := testClientThatFailsTestIfHttpIsCalled(t)

	ctx := context.Background()
	_, err := client.AddFriend(ctx, "", nil)
	assert.ErrorIs(t, err, ErrInvalidParameter)
}

func Test_AddFriend_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.AddFriend(ctx, "test@test.com", nil)

		return err
	}

	doBasicErrorChecks(t, f)
}

const addFriends200Response = `
{
  "users": [
    {
      "id": 26,
      "first_name": "Ada",
      "last_name": "Lovelace",
      "email": "ada@example.com",
      "registration_status": "confirmed",
      "picture": {
        "small": "string",
        "medium": "string",
        "large": "string"
      },
      "groups": [
        {
          "group_id": 571,
          "balance": [
            {
              "currency_code": "USD",
              "amount": "414.5"
            }
          ]
        }
      ],
      "balance": [
        {
          "currency_code": "USD",
          "amount": "414.5"
        }
      ],
      "updated_at": "2019-08-24T14:15:22Z"
    }
  ]
}
`

func Test_AddFriends_SanityCheck(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		friendFirstname = "Ada"
		friendLastname  = "Lovelace"
		friendEmail     = "ada@example.com"
	)

	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := struct {
			Firstname *string `json:"friends__0__first_name"`
			Lastname  *string `json:"friends__0__last_name"`
			Email     *string `json:"friends__0__email"`
		}{}

		rawBody, err := io.ReadAll(r.Body)
		require.NoError(err)

		err = json.Unmarshal(rawBody, &in)
		require.NoError(err)

		require.NotNil(in.Firstname)
		require.NotNil(in.Lastname)
		require.NotNil(in.Email)

		require.Equal(friendFirstname, *in.Firstname)
		require.Equal(friendLastname, *in.Lastname)
		require.Equal(friendEmail, *in.Email)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(addFriends200Response))
	}))
	defer cancel()

	ctx := context.Background()

	users := []FriendUser{
		{
			Email:     friendEmail,
			Firstname: friendFirstname,
			Lastname:  friendLastname,
		},
	}

	fs, err := client.AddFriends(ctx, users)
	assert.NoError(err)

	require.Len(fs, 1)

	f := fs[0]
	assert.Equal(resources.UserID(26), f.ID)
	assert.Equal(friendEmail, f.Email)
	assert.Equal(friendFirstname, f.FirstName)
	assert.Equal(friendLastname, f.LastName)
}

func Test_AddFriends_ShouldFailIfNoFriends(t *testing.T) {
	ctx := context.Background()
	client := testClientThatFailsTestIfHttpIsCalled(t)

	_, err := client.AddFriends(ctx, nil)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidParameter)

	_, err = client.AddFriends(ctx, []FriendUser{})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidParameter)
}

const addFriends200ResponseWithErrors = `
{
  "users": [
    {
      "id": 30,
      "first_name": "Ada",
      "last_name": "Lovelace",
      "email": "ada@example.com",
      "registration_status": "confirmed",
      "picture": {
        "small": "string",
        "medium": "string",
        "large": "string"
      },
      "groups": [
        {
          "group_id": 571,
          "balance": [
            {
              "currency_code": "USD",
              "amount": "414.5"
            }
          ]
        }
      ],
      "balance": [
        {
          "currency_code": "USD",
          "amount": "414.5"
        }
      ],
      "updated_at": "2019-08-24T14:15:22Z"
    }
  ],
  "errors": {
    "property1": [
      "string"
    ],
    "property2": [
      "string"
    ]
  }
}
`

func Test_AddFriends_200ResponseWithErrorsShouldFail(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		friendFirstname = "Ada"
		friendLastname  = "Lovelace"
		friendEmail     = "ada@example.com"
	)

	client, cancel := testClient(t, http.StatusOK, http.MethodPost, addFriends200ResponseWithErrors)
	defer cancel()

	ctx := context.Background()

	users := []FriendUser{
		{Email: friendEmail},
	}

	fs, err := client.AddFriends(ctx, users)

	assert.Error(err)
	assert.ErrorIs(err, ErrUnsuccessful)

	require.Len(fs, 1)
	f := fs[0]
	assert.Equal(resources.UserID(30), f.ID)
	assert.Equal(friendFirstname, f.FirstName)
	assert.Equal(friendLastname, f.LastName)
	assert.Equal(friendEmail, f.Email)
}

func Test_AddFriends_ShouldFailIfEmailIsMissing(t *testing.T) {
	assert := assert.New(t)

	client := testClientThatFailsTestIfHttpIsCalled(t)

	ctx := context.Background()

	users := []FriendUser{
		{Firstname: "Test"},
	}

	fs, err := client.AddFriends(ctx, users)
	assert.ErrorIs(err, ErrInvalidParameter)
	assert.Empty(fs)
}

func Test_AddFriends_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		users := []FriendUser{
			{Email: "test@mail.com"},
			{Email: "test@mail.com"},
			{Email: "test@mail.com"},
		}

		_, err := client.AddFriends(ctx, users)

		return err
	}

	doBasicErrorChecks(t, f)
}
