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

type friendsContainer struct {
	Friends []resources.Friend `json:"friends"`
}

type friendContainer struct {
	Friend resources.Friend `json:"friend"`
}

type friendUsersContainer struct {
	Users []resources.User `json:"users"`
}

func (c *Client) GetFriends(ctx context.Context) ([]resources.Friend, error) {
	const path = "/get_friends"

	res, err := c.do(ctx, http.MethodGet, path, url.Values{}, nil)
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

	var container friendsContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Friends, nil
}

func (c *Client) GetFriend(ctx context.Context, id int) (resources.Friend, error) {
	const basePath = "/get_friend"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodGet, path, url.Values{}, nil)
	if err != nil {
		return resources.Friend{}, err
	}

	if res.StatusCode != http.StatusOK {
		return resources.Friend{}, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.Friend{}, err
	}
	defer res.Body.Close()

	var container friendContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.Friend{}, err
	}

	return container.Friend, nil
}

type friendParam string
type FriendParams map[friendParam]interface{}

const (
	FriendFirstname friendParam = "user_first_name"
	FriendLastname  friendParam = "user_last_name"
)

func (c *Client) AddFriend(ctx context.Context, email string, params FriendParams) (resources.Friend, error) {
	const basePath = "/create_friend"

	if email == "" {
		return resources.Friend{}, fmt.Errorf("email cannot be empty")
	}

	bodyParams := make(map[string]interface{})
	for k, v := range params {
		bodyParams[string(k)] = v
	}

	bodyParams["user_email"] = email

	res, err := c.do(ctx, http.MethodPost, basePath, nil, bodyParams)
	if err != nil {
		return resources.Friend{}, err
	}

	if res.StatusCode != http.StatusOK {
		return resources.Friend{}, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.Friend{}, err
	}
	defer res.Body.Close()

	var container friendContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.Friend{}, err
	}

	return container.Friend, nil
}

type FriendUser struct {
	Firstname string
	Lastname  string
	Email     string
}

func addFriendUserParamsToMap(idx int, f FriendUser, m map[string]interface{}) error {
	const format = "friends__%d__%s"

	if f.Email == "" {
		return fmt.Errorf("email is required for the friend: %+v", f)
	}

	if v := f.Email; v != "" {
		k := fmt.Sprintf(format, idx, "email")
		m[k] = v
	}

	if v := f.Firstname; v != "" {
		k := fmt.Sprintf(format, idx, "first_name")
		m[k] = v
	}

	if v := f.Lastname; v != "" {
		k := fmt.Sprintf(format, idx, "last_name")
		m[k] = v
	}

	return nil
}

func (c *Client) AddFriends(ctx context.Context, friends []FriendUser) ([]resources.User, error) {
	const basePath = "/create_friends"

	bodyParams := make(map[string]interface{})

	for idx, f := range friends {
		if err := addFriendUserParamsToMap(idx, f, bodyParams); err != nil {
			return nil, err
		}
	}

	res, err := c.do(ctx, http.MethodPost, basePath, nil, bodyParams)
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

	var container friendUsersContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Users, nil
}

func (c *Client) DeleteFriend(ctx context.Context, id int) error {
	const basePath = "/delete_friend"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return err
	}

	if err := handleStatusOkErrorResponse(res, nil); err != nil {
		return err
	}

	return nil
}