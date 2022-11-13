package splitwise

import (
	"context"
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

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := c.getErrorFromResponse(res, rawBody); err != nil {
		return nil, err
	}

	var container notificationsContainer
	err = c.unmarshal()(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Notifications, nil
}
