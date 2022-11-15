package splitwise

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const getGroups200Response = `
{
  "groups": [
    {
      "id": 321,
      "name": "Housemates 2020",
      "group_type": "apartment",
      "updated_at": "2019-08-24T14:15:22Z",
      "simplify_by_default": true,
      "members": [
        {
          "id": 0,
          "first_name": "Ada",
          "last_name": "Lovelace",
          "email": "ada@example.com",
          "registration_status": "confirmed",
          "picture": {
            "small": "string",
            "medium": "string",
            "large": "string"
          },
          "balance": [
            {
              "currency_code": "USD",
              "amount": "-5.02"
            }
          ]
        }
      ],
      "original_debts": [
        {
          "from": 18523,
          "to": 90261,
          "amount": "414.5",
          "currency_code": "USD"
        }
      ],
      "simplified_debts": [
        {
          "from": 18523,
          "to": 90261,
          "amount": "414.5",
          "currency_code": "USD"
        }
      ],
      "avatar": {
        "original": null,
        "xxlarge": "https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-ruby2-house-1000px.png",
        "xlarge": "https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-ruby2-house-500px.png",
        "large": "https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-ruby2-house-200px.png",
        "medium": "https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-ruby2-house-100px.png",
        "small": "https://s3.amazonaws.com/splitwise/uploads/group/default_avatars/avatar-ruby2-house-50px.png"
      },
      "custom_avatar": true,
      "cover_photo": {
        "xxlarge": "https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-ruby-1000px.png",
        "xlarge": "https://s3.amazonaws.com/splitwise/uploads/group/default_cover_photos/coverphoto-ruby-500px.png"
      },
      "invite_link": "https://www.splitwise.com/join/abQwErTyuI+12"
    }
  ]
}
`

func Test_GetGroups_SanityCheck(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	client, cancel := testClient(t, http.StatusOK, http.MethodGet, getGroups200Response)
	defer cancel()

	ctx := context.Background()
	gs, err := client.GetGroups(ctx)
	require.NoError(err)

	require.Len(gs, 1)

	g := gs[0]
	assert.Equal(resources.GroupID(321), g.ID)
	assert.Equal("Housemates 2020", g.Name)
	assert.Equal("apartment", g.Type)
}

func Test_GetGroups_BasicErrorChecks(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		gs, err := client.GetGroups(ctx)
		assert.Len(t, gs, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_GetGroup(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		groupID   = 50
		groupName = "Testing"
		groupType = "apartment"
	)

	client, cancel := testClientWithHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var gId int
		path := r.URL.Path
		_, err := fmt.Sscanf(path, DefaultApiVersionPath+"/get_group/%d", &gId)
		require.NoError(err)

		require.Equal(groupID, gId)

		res := struct {
			Group resources.Group `json:"group"`
		}{
			Group: resources.Group{
				ID:   resources.GroupID(gId),
				Name: groupName,
				Type: groupType,
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}))
	defer cancel()

	ctx := context.Background()
	g, err := client.GetGroup(ctx, groupID)
	require.NoError(err)

	assert.Equal(resources.GroupID(groupID), g.ID)
	assert.Equal(groupName, g.Name)
	assert.Equal(groupType, g.Type)
}

func Test_GetGroup_BasicErrorChecks(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.GetGroup(ctx, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_CreateGroup(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		groupID        = 100
		groupName      = "Test name"
		groupType      = "trip"
		user0ID        = 200
		user0Firstname = "Testing"
		user0Lastname  = "Testing"
		user0Email     = "testing@test.com"
		user1ID        = 5823
	)

	client, cancel := testClientWithHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p string
		path := r.URL.Path
		_, err := fmt.Sscanf(path, DefaultApiVersionPath+"/%s", &p)
		require.NoError(err)
		require.Equal("create_group", p)

		in := struct {
			Name           string `json:"name"`
			Type           string `json:"group_type"`
			User0Firstname string `json:"users__0__first_name"`
			User0Lastname  string `json:"users__0__last_name"`
			User0Email     string `json:"users__0__email"`
			User1ID        string `json:"users__1__id"`
		}{}

		rawBody, err := io.ReadAll(r.Body)
		require.NoError(err)
		defer r.Body.Close()

		err = json.Unmarshal(rawBody, &in)
		require.NoError(err)

		require.Equal(groupName, in.Name)
		require.Equal(groupType, in.Type)
		require.Equal(user0Firstname, in.User0Firstname)
		require.Equal(user0Lastname, in.User0Lastname)
		require.Equal(user0Email, in.User0Email)
		require.Equal(strconv.Itoa(user1ID), in.User1ID)

		res := struct {
			Group resources.Group `json:"group"`
		}{
			Group: resources.Group{
				ID:   resources.GroupID(groupID),
				Name: groupName,
				Type: groupType,
				Members: []resources.User{
					{
						ID:        resources.UserID(user0ID),
						FirstName: user0Firstname,
						LastName:  user0Lastname,
						Email:     user0Email,
					},
					{
						ID: resources.UserID(user1ID),
					},
				},
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}))
	defer cancel()

	ctx := context.Background()
	g, err := client.CreateGroup(ctx, groupName, GroupParams{
		GroupType: groupType,
	}, []GroupUser{
		{
			Firstname: user0Firstname,
			Lastname:  user0Lastname,
			Email:     user0Email,
		},
		{
			Id: resources.UserID(user1ID),
		},
	})
	require.NoError(err)

	assert.Equal(resources.GroupID(groupID), g.ID)
	assert.Equal(groupName, g.Name)
	assert.Equal(groupType, g.Type)

	require.Len(g.Members, 2)

	mems := g.Members

	assert.Equal(resources.UserID(user0ID), mems[0].ID)
	assert.Equal(user0Firstname, mems[0].FirstName)
	assert.Equal(user0Lastname, mems[0].LastName)
	assert.Equal(user0Email, mems[0].Email)

	assert.Equal(resources.UserID(user1ID), mems[1].ID)
}

func Test_CreateGroup_BasicErrorChecks(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.CreateGroup(ctx, "Test", nil, nil)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_CreateGroup_EmptyNameProducesError(t *testing.T) {
	client := testClientThatFailsTestIfHttpIsCalled(t)

	ctx := context.Background()
	_, err := client.CreateGroup(ctx, "", nil, nil)
	assert.ErrorIs(t, err, ErrInvalidParameter)
}

func Test_CreateGroup_ShouldFailIfUserHasMissingParams(t *testing.T) {
	client := testClientThatFailsTestIfHttpIsCalled(t)

	ctx := context.Background()
	_, err := client.CreateGroup(ctx, "Test", nil, []GroupUser{
		{},
	})
	assert.ErrorIs(t, err, ErrInvalidParameter)
}

func Test_DeleteGroup(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		groupID = 500
	)

	client, cancel := testClientWithHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var gID int
		path := r.URL.Path
		_, err := fmt.Sscanf(path, DefaultApiVersionPath+"/delete_group/%d", &gID)
		require.NoError(err)
		require.Equal(groupID, gID)

		w.WriteHeader(http.StatusOK)
	}))
	defer cancel()

	ctx := context.Background()
	err := client.DeleteGroup(ctx, groupID)
	assert.NoError(err)
}

func Test_DeleteGroup_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		err := client.DeleteGroup(ctx, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_RestoreGroup(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const groupID = 500

	client, cancel := testClientWithHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var gID int
		path := r.URL.Path
		_, err := fmt.Sscanf(path, DefaultApiVersionPath+"/undelete_group/%d", &gID)
		require.NoError(err)
		require.Equal(groupID, gID)

		w.WriteHeader(http.StatusOK)
	}))
	defer cancel()

	ctx := context.Background()
	err := client.RestoreGroup(ctx, groupID)
	assert.NoError(err)
}

func Test_RestoreGroup_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		err := client.RestoreGroup(ctx, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_AddUserToGroupFromUserId(t *testing.T) {
	require := require.New(t)

	const (
		groupID = 100
		userID  = 200
	)

	client, cancel := testClientWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		in := struct {
			GroupID int `json:"group_id"`
			UserID  int `json:"user_id"`
		}{}

		rawBody, err := io.ReadAll(r.Body)
		require.NoError(err)

		err = json.Unmarshal(rawBody, &in)
		require.NoError(err)

		require.Equal(groupID, in.GroupID)
		require.Equal(userID, in.UserID)

		const res = `
{
  "success": true,
  "user": {},
  "errors": {}
}`

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res))
	})
	defer cancel()

	ctx := context.Background()
	err := client.AddUserToGroupFromUserId(ctx, groupID, userID)
	assert.NoError(t, err)
}

func Test_AddUserToGroupFromUserId_ShouldFailIfInvalidParameters(t *testing.T) {
	client := testClientThatFailsTestIfHttpIsCalled(t)

	ctx := context.Background()
	err := client.AddUserToGroupFromUserId(ctx, 15, 0)
	assert.ErrorIs(t, err, ErrInvalidParameter)

	err = client.AddUserToGroupFromUserId(ctx, 0, 15)
	assert.ErrorIs(t, err, ErrInvalidParameter)

	err = client.AddUserToGroupFromUserId(ctx, 0, 0)
	assert.ErrorIs(t, err, ErrInvalidParameter)
}

func Test_AddUserToGroupFromUserId_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		err := client.AddUserToGroupFromUserId(ctx, 15, 20)

		return err
	}

	doBasicErrorChecks(t, f)
}
