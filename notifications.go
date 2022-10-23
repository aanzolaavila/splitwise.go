package splitwise

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aanzolaavila/splitwise.go/resources"
)

type notificationsContainer struct {
	Notifications []resources.Notification `json:"notifications"`
}

type notificationsParam string
type NotificationsParams map[notificationsParam]interface{}

const (
	NotificationsUpdatedAfter notificationsParam = "updated_after"
	NotificationsLimit        notificationsParam = "limit"
)

func (c *Client) GetNotifications(ctx context.Context, params NotificationsParams) ([]resources.Notification, error) {
	const path = "/get_notifications"

	bodyParams := make(map[string]interface{})
	for k, v := range params {
		bodyParams[string(k)] = v
	}

	res, err := c.do(ctx, http.MethodGet, path, nil, bodyParams)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var container notificationsContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Notifications, nil
}
