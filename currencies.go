package splitwise

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/aanzolaavila/splitwise.go/resources"
)

type currenciesContainer struct {
	Currencies []resources.Currency `json:"currencies"`
}

func (c *Client) GetCurrencies(ctx context.Context) ([]resources.Currency, error) {
	const path = "/get_currencies"

	res, err := c.do(ctx, http.MethodGet, path, nil, nil)
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

	var container currenciesContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Currencies, nil
}