package splitwise

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aanzolaavila/splitwise.go/resources"
)

type categoriesContainer struct {
	Categories []resources.MainCategory `json:"categories"`
}

func (c *Client) GetCategories(ctx context.Context) ([]resources.MainCategory, error) {
	const path = "/get_categories"

	res, err := c.do(ctx, http.MethodGet, path, nil, nil)
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

	var container categoriesContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Categories, nil
}
