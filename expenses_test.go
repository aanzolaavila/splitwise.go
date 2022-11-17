package splitwise

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getAndCheckIntExpensesParam(t *testing.T) {
	const (
		testField                = ExpensesGroupId
		expensesGroupIdAsStr     = "25"
		expensesGroupIdAsInt     = 25
		expensesGroupInvalidStr  = "invalid"
		expensesGroupInvalidType = 25.0
	)

	ps := ExpensesParams{
		testField: expensesGroupIdAsStr,
	}

	t.Run("StringRepresentingIntShouldPass", func(t *testing.T) {
		s, err := getAndCheckIntExpensesParam(ps, testField)
		assert.NoError(t, err)
		assert.Equal(t, expensesGroupIdAsStr, s)
	})

	t.Run("IntShouldPass", func(t *testing.T) {
		ps[ExpensesGroupId] = expensesGroupIdAsInt
		s, err := getAndCheckIntExpensesParam(ps, testField)
		assert.NoError(t, err)
		assert.Equal(t, expensesGroupIdAsStr, s)
	})

	t.Run("StringNOTRepresentingIntShouldFail", func(t *testing.T) {
		ps[ExpensesGroupId] = expensesGroupInvalidStr
		s, err := getAndCheckIntExpensesParam(ps, testField)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidParameter)
		assert.Zero(t, s)
	})

	t.Run("InvalidTypeShouldFail", func(t *testing.T) {
		ps[ExpensesGroupId] = expensesGroupInvalidType
		s, err := getAndCheckIntExpensesParam(ps, testField)
		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidParameter)
		assert.Zero(t, s)
	})

	t.Run("NonExistentParamShouldGiveZeroValueAndNotFail", func(t *testing.T) {
		s, err := getAndCheckIntExpensesParam(ps, ExpensesDatedAfter)
		assert.NoError(t, err)
		assert.Zero(t, s)
	})
}

func Test_getAndCheckDateExpensesParam(t *testing.T) {
	const (
		testField                             = ExpensesDatedBefore
		expensesDatedBeforeAsStr       string = "2022-01-01T12:00:00Z"
		expensesDatedBeforeInvalidStr         = "invalid"
		expensesDatedBeforeInvalidType        = 20.0
	)

	expensesDatedBeforeAsDate, err := time.Parse(time.RFC3339, expensesDatedBeforeAsStr)
	if err != nil {
		require.FailNowf(t, "failed to create date for testing: %s", err.Error())
	}

	ps := ExpensesParams{
		ExpensesDatedBefore: expensesDatedBeforeAsDate,
	}

	t.Run("DateTypeShouldPass", func(t *testing.T) {
		s, err := getAndCheckDateExpensesParam(ps, testField)
		assert.NoError(t, err)
		assert.Equal(t, expensesDatedBeforeAsStr, s)
	})

	t.Run("StringAsValidDateShouldPass", func(t *testing.T) {
		ps[testField] = expensesDatedBeforeAsStr
		s, err := getAndCheckDateExpensesParam(ps, testField)
		assert.NoError(t, err)
		assert.Equal(t, expensesDatedBeforeAsStr, s)
	})

	t.Run("InvalidStringShouldFail", func(t *testing.T) {
		ps[testField] = expensesDatedBeforeInvalidStr
		s, err := getAndCheckDateExpensesParam(ps, testField)
		assert.ErrorIs(t, err, ErrInvalidParameter)
		assert.Zero(t, s)
	})

	t.Run("NonExistentParamShouldGiveZeroValueAndNotFail", func(t *testing.T) {
		s, err := getAndCheckDateExpensesParam(ps, ExpensesDatedAfter)
		assert.NoError(t, err)
		assert.Zero(t, s)
	})

	t.Run("InvalidTypeShouldFail", func(t *testing.T) {
		ps[testField] = expensesDatedBeforeInvalidType
		s, err := getAndCheckDateExpensesParam(ps, testField)
		assert.ErrorIs(t, err, ErrInvalidParameter)
		assert.Zero(t, s)
	})
}

func Test_expensesParamsToUrlValues(t *testing.T) {
	var (
		require = require.New(t)
		assert  = assert.New(t)
	)

	const timeFormat = time.RFC3339
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	ps := ExpensesParams{
		// ints
		ExpensesGroupId:  1,
		ExpensesFriendId: "2",
		ExpensesLimit:    3,
		ExpensesOffset:   "4",
		// dates
		ExpensesDatedBefore:   now,
		ExpensesDatedAfter:    now.Format(timeFormat),
		ExpensesUpdatedBefore: now,
		ExpensesUpdatedAfter:  yesterday,
	}

	vals, err := expensesParamsToUrlValues(ps)
	require.NoError(err)
	require.NotNil(vals)
	assert.Len(vals, len(ps))

	test := func(field expensesParam, expected string) {
		if k := string(field); assert.Contains(vals, k) {
			vs := vals[k]
			require.Len(vs, 1)
			v := vs[0]
			assert.Equal(expected, v)
		}
	}

	test(ExpensesGroupId, "1")
	test(ExpensesFriendId, "2")
	test(ExpensesLimit, "3")
	test(ExpensesOffset, "4")
	test(ExpensesDatedBefore, now.Format(timeFormat))
	test(ExpensesDatedAfter, now.Format(timeFormat))
	test(ExpensesUpdatedBefore, now.Format(timeFormat))
	test(ExpensesUpdatedAfter, yesterday.Format(timeFormat))
}

func Test_expensesParamsToUrlValues_ErrorCases(t *testing.T) {
	ps := ExpensesParams{
		// ints
		ExpensesGroupId: 1.0,
	}

	vals, err := expensesParamsToUrlValues(ps)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidParameter)
	assert.Zero(t, vals)

	ps = ExpensesParams{
		// dates
		ExpensesDatedBefore: 2.0,
	}

	vals, err = expensesParamsToUrlValues(ps)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidParameter)
	assert.Zero(t, vals)
}

const getExpenses200Response = `
{
  "expenses": [
    {
      "cost": "25.0",
      "description": "Brunch",
      "details": "string",
      "date": "2012-05-02T13:00:00Z",
      "repeat_interval": "never",
      "currency_code": "USD",
      "category_id": 15,
      "id": 51023,
      "group_id": 391,
      "friendship_id": 4818,
      "expense_bundle_id": 491030,
      "repeats": true,
      "email_reminder": true,
      "email_reminder_in_advance": null,
      "next_repeat": "string",
      "comments_count": 0,
      "payment": true,
      "transaction_confirmed": true,
      "repayments": [
        {
          "from": 6788709,
          "to": 270896089,
          "amount": "25.0"
        }
      ],
      "created_at": "2012-07-27T06:17:09Z",
      "created_by": {
        "id": 0,
        "first_name": "Ada",
        "last_name": "Lovelace",
        "email": "ada@example.com",
        "registration_status": "confirmed",
        "picture": {
          "small": "string",
          "medium": "string",
          "large": "string"
        }
      },
      "updated_at": "2012-12-23T05:47:02Z",
      "updated_by": {
        "id": 0,
        "first_name": "Ada",
        "last_name": "Lovelace",
        "email": "ada@example.com",
        "registration_status": "confirmed",
        "picture": {
          "small": "string",
          "medium": "string",
          "large": "string"
        }
      },
      "deleted_at": "2012-12-23T05:47:02Z",
      "deleted_by": {
        "id": 0,
        "first_name": "Ada",
        "last_name": "Lovelace",
        "email": "ada@example.com",
        "registration_status": "confirmed",
        "picture": {
          "small": "string",
          "medium": "string",
          "large": "string"
        }
      },
      "category": {
        "id": 5,
        "name": "Electricity"
      },
      "receipt": {
        "large": "https://splitwise.s3.amazonaws.com/uploads/expense/receipt/3678899/large_95f8ecd1-536b-44ce-ad9b-0a9498bb7cf0.png",
        "original": "https://splitwise.s3.amazonaws.com/uploads/expense/receipt/3678899/95f8ecd1-536b-44ce-ad9b-0a9498bb7cf0.png"
      },
      "users": [
        {
          "user": {
            "id": 491923,
            "first_name": "Jane",
            "last_name": "Doe",
            "picture": {
              "medium": "image_url"
            }
          },
          "user_id": 491923,
          "paid_share": "8.99",
          "owed_share": "4.5",
          "net_balance": "4.49"
        }
      ],
      "comments": [
        {
          "id": 79800950,
          "content": "John D. updated this transaction: - The cost changed from $6.99 to $8.99",
          "comment_type": "System",
          "relation_type": "ExpenseComment",
          "relation_id": 855870953,
          "created_at": "2019-08-24T14:15:22Z",
          "deleted_at": "2019-08-24T14:15:22Z",
          "user": {
            "id": 491923,
            "first_name": "Jane",
            "last_name": "Doe",
            "picture": {
              "medium": "image_url"
            }
          }
        }
      ]
    }
  ]
}
`

func Test_GetExpenses(t *testing.T) {
	const timeFormat = time.RFC3339

	var (
		require = require.New(t)
		assert  = assert.New(t)
	)

	intParams := []expensesParam{ExpensesGroupId, ExpensesFriendId, ExpensesLimit, ExpensesOffset}
	dateParams := []expensesParam{ExpensesDatedAfter, ExpensesDatedBefore, ExpensesUpdatedBefore, ExpensesUpdatedAfter}
	expectedDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	client, cancel := testClientWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		vals := r.URL.Query()
		for _, p := range intParams {
			k := string(p)
			require.Contains(vals, k)
			vs := vals[k]
			require.Len(vs, 1)
			_, err := strconv.Atoi(vs[0])
			assert.NoError(err)
		}

		for _, p := range dateParams {
			k := string(p)
			require.Contains(vals, k)
			vs := vals[k]
			require.Len(vs, 1)
			d, err := time.Parse(timeFormat, vs[0])
			require.NoError(err, "date conversion failed on %s: %s", k, vs[0])
			assert.True(d.Equal(expectedDate), "dates are not equal [%v] != [%v]", expectedDate, d)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(getExpenses200Response))
	})
	defer cancel()

	ctx := context.Background()
	ex, err := client.GetExpenses(ctx, ExpensesParams{
		// ints
		ExpensesGroupId:  1,
		ExpensesFriendId: "2",
		ExpensesLimit:    3,
		ExpensesOffset:   "4",
		// dates
		ExpensesDatedBefore:   expectedDate,
		ExpensesDatedAfter:    expectedDate.Format(timeFormat),
		ExpensesUpdatedBefore: expectedDate,
		ExpensesUpdatedAfter:  expectedDate.Format(timeFormat),
	})

	require.NoError(err)
	require.Len(ex, 1)

	e := ex[0]
	assert.Equal(resources.ExpenseID(51023), e.ID)
}

func Test_GetExpenses_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		ex, err := client.GetExpenses(ctx, nil)
		assert.Len(t, ex, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_GetExpenses_ShouldFailOnInvalidInput(t *testing.T) {
	client := testClientThatFailsTestIfHttpIsCalled(t)

	ctx := context.Background()
	ex, err := client.GetExpenses(ctx, ExpensesParams{
		ExpensesGroupId: 1.0,
	})
	assert.ErrorIs(t, err, ErrInvalidParameter)
	assert.Len(t, ex, 0)
}

func Test_GetExpense(t *testing.T) {
	const expenseID = 100

	client, cancel := testClientWithHandler(t, func(w http.ResponseWriter, r *http.Request) {
		var eID int
		path := r.URL.Path
		_, err := fmt.Sscanf(path, DefaultApiVersionPath+"/get_expense/%d", &eID)
		require.NoError(t, err)

		assert.Equal(t, expenseID, eID)

		res := struct {
			Expense resources.Expense `json:"expense"`
		}{
			Expense: resources.Expense{
				ID: resources.ExpenseID(expenseID),
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	})
	defer cancel()

	ctx := context.Background()
	e, err := client.GetExpense(ctx, expenseID)
	require.NoError(t, err)

	assert.Equal(t, resources.ExpenseID(expenseID), e.ID)
}

func Test_GetExpense_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		e, err := client.GetExpense(ctx, 0)
		assert.Zero(t, e)

		return err
	}

	doBasicErrorChecks(t, f)
}
