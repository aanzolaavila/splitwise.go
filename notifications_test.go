package splitwise

import (
	"context"
	"net/http"
	"testing"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
)

const notificationsSuccessResponse = `
{
  "notifications": [
    {
      "id": 32514315,
      "type": 0,
      "created_at": "2019-08-24T14:15:22Z",
      "created_by": 2,
      "source": {
        "type": "Expense",
        "id": 865077,
        "url": "string"
      },
      "image_url": "https://s3.amazonaws.com/splitwise/uploads/notifications/v2/0-venmo.png",
      "image_shape": "square",
      "content": "<strong>You</strong> paid <strong>Jon H.</strong>.<br><font color=\\\"#5bc5a7\\\">You paid $23.45</font>"
    }
  ]
}
`

func Test_GetNotifications_SanityCheck(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, notificationsSuccessResponse)
	defer cancel()

	ctx := context.Background()

	nots, err := client.GetNotifications(ctx, NotificationsParams{
		NotificationsLimit: 100,
	})
	assert.NoError(t, err)

	assert.Len(t, nots, 1)

	not := nots[0]
	assert.Equal(t, resources.NotificationID(32514315), not.ID)
	assert.Equal(t, 0, not.Type)
	assert.Equal(t, "Expense", not.Source.Type)
	assert.Equal(t, resources.Identifier(865077), not.Source.ID)
}

func Test_GetNotifications_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		nots, err := client.GetNotifications(ctx, nil)
		assert.Len(t, nots, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}
