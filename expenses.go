package splitwise

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aanzolaavila/splitwise.go/resources"
)

type expensesContainer struct {
	Expenses []resources.ExpenseResponse `json:"expenses"`
}

type expenseContainer struct {
	Expense resources.ExpenseResponse `json:"expense"`
}

type expensesParam string
type ExpensesParams map[expensesParam]interface{}

const (
	ExpensesGroupId       expensesParam = "group_id"
	ExpensesFriendId      expensesParam = "friend_id"
	ExpensesDatedAfter    expensesParam = "dated_after"
	ExpensesDatedBefore   expensesParam = "dated_before"
	ExpensesUpdatedAfter  expensesParam = "updated_after"
	ExpensesUpdatedBefore expensesParam = "updated_before"
	ExpensesLimit         expensesParam = "limit"
	ExpensesOffset        expensesParam = "offset"
)

func getAndCheckIntExpensesParam(params ExpensesParams, key expensesParam) (string, error) {
	value, ok := params[key]
	if ok {
		strValue, ok := value.(string)
		if ok {
			if _, err := strconv.Atoi(strValue); err != nil {
				return "", fmt.Errorf("%s is not convertable to int: %v", key, err)
			}
			return strValue, nil
		}

		intValue, ok := value.(int)
		if !ok {
			return "", fmt.Errorf("%s is not an int", key)
		}
		return strconv.Itoa(intValue), nil
	}

	return "", nil
}

func getAndCheckDateExpensesParam(params ExpensesParams, key expensesParam) (string, error) {
	value, ok := params[key]
	if ok {
		const timeFormat = time.RFC3339
		timeValue, ok := value.(time.Time)
		if ok {
			return timeValue.Format(timeFormat), nil
		}

		strValue, ok := value.(string)
		if ok {
			if _, err := time.Parse(timeFormat, strValue); err != nil {
				return "", fmt.Errorf("%s does not have [%s] format", key, timeFormat)
			}

			return strValue, nil
		}

		return "", fmt.Errorf("%s is not a string nor a date", key)
	}

	return "", nil
}

func expensesParamsToUrlValues(params ExpensesParams) (url.Values, error) {
	values := url.Values{}

	intParams := []expensesParam{ExpensesGroupId, ExpensesFriendId, ExpensesLimit, ExpensesOffset}
	for _, p := range intParams {
		val, err := getAndCheckIntExpensesParam(params, p)
		if err != nil {
			return nil, err
		}

		if val != "" {
			values.Set(string(p), val)
		}
	}

	dateParams := []expensesParam{ExpensesDatedAfter, ExpensesDatedBefore, ExpensesUpdatedBefore, ExpensesUpdatedAfter}
	for _, p := range dateParams {
		val, err := getAndCheckDateExpensesParam(params, p)
		if err != nil {
			return nil, err
		}

		if val != "" {
			values.Set(string(p), val)
		}
	}

	return values, nil
}

func (c *Client) GetExpenses(ctx context.Context, params ExpensesParams) ([]resources.ExpenseResponse, error) {
	const basePath = "/get_expenses"

	expensesValues, err := expensesParamsToUrlValues(params)
	if err != nil {
		return nil, err
	}

	res, err := c.do(ctx, http.MethodGet, basePath, expensesValues, nil)
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

	var container expensesContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Expenses, nil
}

func (c *Client) GetExpense(ctx context.Context, id int) (resources.ExpenseResponse, error) {
	const basePath = "/get_expense"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return resources.ExpenseResponse{}, err
	}

	if res.StatusCode != http.StatusOK {
		return resources.ExpenseResponse{}, handleResponseError(res)
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.ExpenseResponse{}, err
	}
	defer res.Body.Close()

	var container expenseContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return resources.ExpenseResponse{}, err
	}

	return container.Expense, nil
}

type createExpenseParam string
type CreateExpenseParams map[createExpenseParam]interface{}

const (
	CreateExpenseDetails        createExpenseParam = "details"
	CreateExpenseDate           createExpenseParam = "date"
	CreateExpenseRepeatInterval createExpenseParam = "repeat_interval"
	CreateExpenseCurrencyCode   createExpenseParam = "currency_code"
	CreateExpenseCategoryId     createExpenseParam = "category_id"
)

func (c *Client) CreateExpenseEqualGroupSplit(ctx context.Context, cost, description string, groupId int, splitEqually bool, params CreateExpenseParams) ([]resources.ExpenseResponse, error) {
	const basePath = "/create_expense"

	if cost == "" || description == "" {
		return nil, fmt.Errorf("cost and description must be non-empty")
	}

	m := make(map[string]interface{})
	for k, v := range params {
		m[string(k)] = v
	}

	m["cost"] = cost
	m["description"] = description
	m["group_id"] = groupId
	m["split_equally"] = splitEqually

	res, err := c.do(ctx, http.MethodPost, basePath, nil, m)
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
	if err := handleStatusOkErrorResponse(res, rawBody); err != nil {
		return nil, err
	}

	var container expensesContainer
	err = json.Unmarshal(rawBody, &container)
	if err != nil {
		return nil, err
	}

	return container.Expenses, nil
}

func (c *Client) DeleteExpense(ctx context.Context, id int) error {
	const basePath = "/delete_expense"

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

func (c *Client) RestoreExpense(ctx context.Context, id int) error {
	const basePath = "/undelete_expense"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodPost, path, nil, nil)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return handleResponseError(res)
	}

	if err := handleStatusOkErrorResponse(res, nil); err != nil {
		return err
	}

	return nil
}
