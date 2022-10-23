package splitwise

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aanzolaavila/splitwise.go/resources"
)

type groupsContainer struct {
	Groups []resources.Group `json:"groups"`
}

type groupContainer struct {
	Group resources.Group `json:"group"`
}

func (c *Client) GetGroups(ctx context.Context) ([]resources.Group, error) {
	const path = "/get_groups"

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

	var container groupsContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Groups, nil
}

func (c *Client) GetGroup(ctx context.Context, id int) (resources.Group, error) {
	const basePath = "/get_group"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodGet, path, url.Values{}, nil)
	if err != nil {
		return resources.Group{}, err
	}

	if res.StatusCode != http.StatusOK {
		return resources.Group{}, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.Group{}, err
	}
	defer res.Body.Close()

	var container groupContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.Group{}, err
	}

	return container.Group, nil
}

type groupParam string
type GroupParams map[groupParam]interface{}

const (
	GroupType              groupParam = "group_type"
	GroupSimplifyByDefault groupParam = "simplify_by_default"
)

type GroupUser struct {
	Id        resources.UserID
	Email     string
	Firstname string
	Lastname  string
}

func addGroupUserParamsToMap(idx int, u GroupUser, m map[string]interface{}) error {
	const format = "users__%d__%s"

	if u.Id == 0 && u.Email == "" {
		return fmt.Errorf("id or email is required for the user: %+v", u)
	}

	if v := strconv.Itoa(int(u.Id)); u.Id != 0 {
		k := fmt.Sprintf(format, idx, "id")
		m[k] = v
	}

	if v := u.Email; v != "" {
		k := fmt.Sprintf(format, idx, "email")
		m[k] = v
	}

	if v := u.Firstname; v != "" {
		k := fmt.Sprintf(format, idx, "first_name")
		m[k] = v
	}

	if v := u.Lastname; v != "" {
		k := fmt.Sprintf(format, idx, "last_name")
		m[k] = v
	}

	return nil
}

func (c *Client) CreateGroup(ctx context.Context, name string, params GroupParams, users []GroupUser) (resources.Group, error) {
	const basePath = "/create_group"

	if name == "" {
		return resources.Group{}, fmt.Errorf("name cannot be empty")
	}

	bodyParams := make(map[string]interface{})
	for k, v := range params {
		bodyParams[string(k)] = v
	}

	bodyParams["name"] = name

	for idx, user := range users {
		err := addGroupUserParamsToMap(idx, user, bodyParams)
		if err != nil {
			return resources.Group{}, err
		}
	}

	res, err := c.do(ctx, http.MethodPost, basePath, nil, bodyParams)
	if err != nil {
		return resources.Group{}, err
	}

	if res.StatusCode != http.StatusOK {
		return resources.Group{}, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.Group{}, err
	}
	defer res.Body.Close()

	var container groupContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.Group{}, err
	}

	return container.Group, nil
}

func (c *Client) DeleteGroup(ctx context.Context, id int) error {
	const basePath = "/delete_group"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return handleResponseError(res)
	}

	return nil
}

func (c *Client) RestoreGroup(ctx context.Context, id int) error {
	const basePath = "/undelete_group"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return handleResponseError(res)
	}

	return nil
}

func (c *Client) AddUserToGroupFromUserId(ctx context.Context, groupId, userId int) error {
	const path = "/add_user_to_group"

	if groupId == 0 || userId == 0 {
		return fmt.Errorf("group id and user id must be defined")
	}

	bodyParams := map[string]interface{}{}
	bodyParams["group_id"] = groupId
	bodyParams["user_id"] = userId

	res, err := c.do(ctx, http.MethodPost, path, nil, bodyParams)
	if err != nil {
		return err
	}

	if err := handleStatusOkErrorResponse(res, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) AddUserToGroupFromUserInfo(ctx context.Context, groupId int, firstName, lastName, email string) error {
	const path = "/add_user_to_group"

	if groupId == 0 || firstName == "" || lastName == "" || email == "" {
		return fmt.Errorf("group id, firstname, lastname and email must be defined")
	}

	bodyParams := map[string]interface{}{}
	bodyParams["group_id"] = groupId
	bodyParams["first_name"] = firstName
	bodyParams["last_name"] = lastName
	bodyParams["email"] = email

	res, err := c.do(ctx, http.MethodPost, path, nil, bodyParams)
	if err != nil {
		return err
	}

	if err := handleStatusOkErrorResponse(res, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveUserFromGroup(ctx context.Context, groupId, userId int) error {
	const path = "/remove_user_from_group"

	if groupId == 0 || userId == 0 {
		return fmt.Errorf("group id and user id must be defined")
	}

	bodyParams := map[string]interface{}{}
	bodyParams["group_id"] = groupId
	bodyParams["user_id"] = userId

	res, err := c.do(ctx, http.MethodPost, path, nil, bodyParams)
	if err != nil {
		return err
	}

	if err := handleStatusOkErrorResponse(res, nil); err != nil {
		return err
	}

	return nil
}
