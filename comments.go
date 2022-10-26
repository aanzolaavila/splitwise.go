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

type commentsContainer struct {
	Comments []resources.Comment `json:"comments"`
}

type commentContainer struct {
	Comment resources.Comment `json:"comment"`
}

func (c *Client) GetExpenseComments(ctx context.Context, expenseId int) ([]resources.Comment, error) {
	const path = "/get_comments"

	queryParams := url.Values{}
	queryParams.Set("expense_id", strconv.Itoa(expenseId))

	res, err := c.do(ctx, http.MethodGet, path, queryParams, nil)
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

	var container commentsContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Comments, nil
}

func (c *Client) CreateExpenseComment(ctx context.Context, expenseId int, content string) (resources.Comment, error) {
	const basePath = "/create_comment"

	if content == "" {
		return resources.Comment{}, fmt.Errorf("content cannot be empty")
	}

	bodyParams := make(map[string]interface{})
	bodyParams["expense_id"] = expenseId
	bodyParams["content"] = content

	res, err := c.do(ctx, http.MethodPost, basePath, nil, bodyParams)
	if err != nil {
		return resources.Comment{}, err
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.Comment{}, err
	}
	defer res.Body.Close()

	if err := c.getErrorFromResponse(res, rawBody); err != nil {
		return resources.Comment{}, err
	}

	var container commentContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.Comment{}, err
	}

	return container.Comment, nil
}

func (c *Client) DeleteExpenseComment(ctx context.Context, id int) (resources.Comment, error) {
	const basePath = "/delete_comment"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return resources.Comment{}, err
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.Comment{}, err
	}
	defer res.Body.Close()

	if err := c.getErrorFromResponse(res, rawBody); err != nil {
		return resources.Comment{}, err
	}

	var container commentContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.Comment{}, err
	}

	return container.Comment, nil
}
