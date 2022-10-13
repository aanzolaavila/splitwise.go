package splitwise

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/aanzolaavila/splitwise.go/resources"
)

type groupsContainer struct {
	Groups []resources.Group `json:"groups"`
}

func (c *Client) GetGroups(ctx context.Context) ([]resources.Group, error) {
	const path = "/get_groups"

	res, err := c.do(ctx, http.MethodGet, path, url.Values{}, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, handleError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var container groupsContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Groups, nil
}
