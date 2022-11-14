package splitwise

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		json.NewEncoder(w).Encode(res)
		w.WriteHeader(http.StatusOK)
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
