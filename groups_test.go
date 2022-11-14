package splitwise

import (
	"context"
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
