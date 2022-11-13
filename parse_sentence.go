package splitwise

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/aanzolaavila/splitwise.go/resources"
)

func (c *Client) parseSentenceIntoExpense(ctx context.Context, input string, autosave bool, params map[string]interface{}) (resources.ParsedExpense, error) {
	const path = "/parse_sentence"

	if input == "" {
		return resources.ParsedExpense{}, fmt.Errorf("input must not be empty")
	}

	bodyParams := make(map[string]interface{})
	for k, v := range params {
		bodyParams[k] = v
	}

	bodyParams["input"] = input
	bodyParams["autosave"] = autosave

	res, err := c.do(ctx, http.MethodPost, path, nil, bodyParams)
	if err != nil {
		return resources.ParsedExpense{}, err
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.ParsedExpense{}, err
	}
	defer res.Body.Close()

	err = c.getErrorFromResponse(res, rawBody)
	if err != nil {
		return resources.ParsedExpense{}, err
	}

	var parsedExpense resources.ParsedExpense
	err = c.unmarshal()(rawBody, &parsedExpense)
	if err != nil {
		return resources.ParsedExpense{}, err
	}

	return parsedExpense, nil
}

func (c *Client) ParseSentenceIntoExpenseFreeform(ctx context.Context, input string, autosave bool) (resources.ParsedExpense, error) {
	return c.parseSentenceIntoExpense(ctx, input, autosave, nil)
}

func (c *Client) ParseSentenceIntoExpenseWithFriend(ctx context.Context, input string, friendId int, autosave bool) (resources.ParsedExpense, error) {
	return c.parseSentenceIntoExpense(ctx, input, autosave, map[string]interface{}{
		"friend_id": friendId,
	})
}

func (c *Client) ParseSentenceIntoExpenseWithGroup(ctx context.Context, input string, groupId int, autosave bool) (resources.ParsedExpense, error) {
	return c.parseSentenceIntoExpense(ctx, input, autosave, map[string]interface{}{
		"group_id": groupId,
	})
}
