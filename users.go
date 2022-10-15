package splitwise

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aanzolaavila/splitwise.go/resources"
)

type userContainer struct {
	User resources.User `json:"user"`
}

func (c *Client) GetCurrentUser(ctx context.Context) (resources.User, error) {
	const path = "/get_current_user"

	res, err := c.do(ctx, http.MethodGet, path, url.Values{}, nil)
	if err != nil {
		return resources.User{}, err
	}

	if res.StatusCode != http.StatusOK {
		return resources.User{}, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.User{}, err
	}
	defer res.Body.Close()

	var container userContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.User{}, err
	}

	return container.User, nil
}

func (c *Client) GetUser(ctx context.Context, id int) (resources.User, error) {
	const basePath = "/get_user"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodGet, path, url.Values{}, nil)
	if err != nil {
		return resources.User{}, err
	}

	if res.StatusCode != http.StatusOK {
		return resources.User{}, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.User{}, err
	}
	defer res.Body.Close()

	var container userContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.User{}, err
	}

	return container.User, nil
}

type userParam string
type UserParams map[userParam]string

const (
	UserFirstname       userParam = "first_name"
	UserLastname        userParam = "last_name"
	UserEmail           userParam = "email"
	UserPassword        userParam = "password"
	UserLocale          userParam = "locale"
	UserDefaultCurrency userParam = "default_currency"
)

func (c *Client) UpdateUser(ctx context.Context, id int, params UserParams) (resources.User, error) {
	const basePath = "/update_user"

	path := fmt.Sprintf("%s/%d", basePath, id)

	p := map[string]interface{}{}
	for k, v := range params {
		p[string(k)] = v
	}

	res, err := c.do(ctx, http.MethodPost, path, url.Values{}, p)
	if err != nil {
		return resources.User{}, err
	}

	if res.StatusCode != http.StatusOK {
		return resources.User{}, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.User{}, err
	}
	defer res.Body.Close()

	var container userContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.User{}, err
	}

	return container.User, nil
}
