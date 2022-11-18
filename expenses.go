package splitwise

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aanzolaavila/splitwise.go/resources"
)

type expensesContainer struct {
	Expenses []resources.Expense `json:"expenses"`
}

type expenseContainer struct {
	Expense resources.Expense `json:"expense"`
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
	if !ok {
		return "", nil
	}

	strValue, ok := value.(string)
	if ok {
		if _, err := strconv.Atoi(strValue); err != nil {
			return "", fmt.Errorf("%w: %s is not convertable to int", ErrInvalidParameter, key)
		}
		return strValue, nil
	}

	intValue, ok := value.(int)
	if !ok {
		return "", fmt.Errorf("%w: %s is not an int", ErrInvalidParameter, key)
	}
	return strconv.Itoa(intValue), nil
}

func getAndCheckDateExpensesParam(params ExpensesParams, key expensesParam) (string, error) {
	value, ok := params[key]
	if !ok {
		return "", nil
	}

	const timeFormat = time.RFC3339
	timeValue, ok := value.(time.Time)
	if ok {
		return timeValue.Format(timeFormat), nil
	}

	strValue, ok := value.(string)
	if ok {
		if _, err := time.Parse(timeFormat, strValue); err != nil {
			return "", fmt.Errorf("%w: %s does not have [%s] format", ErrInvalidParameter, key, timeFormat)
		}

		return strValue, nil
	}

	return "", fmt.Errorf("%w: %s is not a string nor a date", ErrInvalidParameter, key)
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

func (c *Client) GetExpenses(ctx context.Context, params ExpensesParams) ([]resources.Expense, error) {
	const basePath = "/get_expenses"

	expensesValues, err := expensesParamsToUrlValues(params)
	if err != nil {
		return nil, err
	}

	res, err := c.do(ctx, http.MethodGet, basePath, expensesValues, nil)
	if err != nil {
		return nil, err
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var container expensesContainer
	err = c.unmarshal()(rawBody, &container)
	if err != nil {
		return nil, err
	}

	if err := c.getErrorFromResponse(res, rawBody); err != nil {
		return nil, err
	}

	return container.Expenses, nil
}

func (c *Client) GetExpense(ctx context.Context, id int) (resources.Expense, error) {
	const basePath = "/get_expense"

	path := fmt.Sprintf("%s/%d", basePath, id)

	res, err := c.do(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return resources.Expense{}, err
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return resources.Expense{}, err
	}
	defer res.Body.Close()

	var container expenseContainer
	err = c.unmarshal()(rawBody, &container)
	if err != nil {
		return resources.Expense{}, err
	}

	if err := c.getErrorFromResponse(res, rawBody); err != nil {
		return resources.Expense{}, err
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

func (c *Client) CreateExpenseEqualGroupSplit(ctx context.Context, cost float64, description string, groupId int, params CreateExpenseParams) ([]resources.Expense, error) {
	const basePath = "/create_expense"

	if cost == 0.0 || description == "" {
		return nil, fmt.Errorf("%w: cost and description must be non-empty", ErrInvalidParameter)
	}

	m := make(map[string]interface{})
	for k, v := range params {
		m[string(k)] = v
	}

	m["cost"] = fmt.Sprintf("%.2f", cost)
	m["description"] = description
	m["group_id"] = groupId
	m["split_equally"] = true

	res, err := c.do(ctx, http.MethodPost, basePath, nil, m)
	if err != nil {
		return nil, err
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var container expensesContainer
	err = c.unmarshal()(rawBody, &container)
	if err != nil {
		return nil, err
	}

	if err := c.getErrorFromResponse(res, rawBody); err != nil {
		return nil, err
	}

	return container.Expenses, nil
}

type ExpenseUser struct {
	Id        resources.UserID
	Email     string
	Firstname string
	Lastname  string
	PaidShare float64
	OwedShare float64
}

func addExpenseUserParamsToMap(idx int, u ExpenseUser, m map[string]interface{}) error {
	const format = "users__%d__%s"

	if u.Id == 0 && u.Email == "" {
		return fmt.Errorf("%w: id or email is required for the user", ErrInvalidParameter)
	}

	if v := int(u.Id); u.Id != 0 {
		k := fmt.Sprintf(format, idx, "user_id")
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

	if v := u.PaidShare; true {
		k := fmt.Sprintf(format, idx, "paid_share")
		m[k] = fmt.Sprintf("%.2f", v)
	}

	if v := u.OwedShare; true {
		k := fmt.Sprintf(format, idx, "owed_share")
		m[k] = fmt.Sprintf("%.2f", v)
	}

	return nil
}

func (c *Client) CreateExpenseByShares(ctx context.Context, cost float64, description string, groupId int, params CreateExpenseParams, users []ExpenseUser) ([]resources.Expense, error) {
	const basePath = "/create_expense"

	if cost == 0.0 || description == "" {
		return nil, fmt.Errorf("%w: cost and description must be non-empty", ErrInvalidParameter)
	}

	m := make(map[string]interface{})
	for k, v := range params {
		m[string(k)] = v
	}

	m["cost"] = fmt.Sprintf("%.2f", cost)
	m["description"] = description
	m["group_id"] = groupId

	for idx, user := range users {
		err := addExpenseUserParamsToMap(idx, user, m)
		if err != nil {
			return nil, err
		}
	}

	res, err := c.do(ctx, http.MethodPost, basePath, nil, m)
	if err != nil {
		return nil, err
	}

	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var container expensesContainer
	err = c.unmarshal()(rawBody, &container)
	if err != nil {
		return nil, err
	}

	if err := c.getErrorFromResponse(res, rawBody); err != nil {
		return nil, err
	}

	return container.Expenses, nil
}

func (c *Client) UpdateExpense(ctx context.Context, id int, cost float64, description string, groupId int, params CreateExpenseParams, users []ExpenseUser) ([]resources.Expense, error) {
	const basePath = "/update_expense"

	path := fmt.Sprintf("%s/%d", basePath, id)

	if cost == 0.0 || description == "" {
		return nil, fmt.Errorf("%w: cost and description must be non-empty", ErrInvalidParameter)
	}

	m := make(map[string]interface{})
	for k, v := range params {
		m[string(k)] = v
	}

	m["cost"] = fmt.Sprintf("%.2f", cost)
	m["description"] = description
	m["group_id"] = groupId

	for idx, user := range users {
		err := addExpenseUserParamsToMap(idx, user, m)
		if err != nil {
			return nil, err
		}
	}

	res, err := c.do(ctx, http.MethodPost, path, nil, m)
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

	var container expensesContainer
	err = c.unmarshal()(rawBody, &container)
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

	if err := c.getErrorFromResponse(res, nil); err != nil {
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

	if err := c.getErrorFromResponse(res, nil); err != nil {
		return err
	}

	return nil
}
