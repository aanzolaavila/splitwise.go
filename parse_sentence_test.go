package splitwise

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const parseSentenceIntoExpense200Response = `
{
  "expense": {
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
  },
  "valid": true,
  "confidence": 0.5
}
`

func Test_ParseSentenceIntoExpenseFreeform_SanityCheck(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	client, cancel := testClient(t, http.StatusOK, http.MethodPost, parseSentenceIntoExpense200Response)
	defer cancel()

	ctx := context.Background()
	e, err := client.ParseSentenceIntoExpenseFreeform(ctx, "I owe Ada 5 bucks", false)
	require.NoError(err)

	assert.True(e.Valid)
	assert.Equal(resources.ExpenseID(51023), e.Expense.ID)
}

func Test_ParseSentenceIntoExpenseFreeform_BasicErrorChecks(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.ParseSentenceIntoExpenseFreeform(ctx, "statement", false)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_ParseSentenceIntoExpenseFreeform_ShouldFailIfInputEmpty(t *testing.T) {
	client := testClientThatFailsTestIfHttpIsCalled(t)

	ctx := context.Background()
	_, err := client.ParseSentenceIntoExpenseFreeform(ctx, "", false)
	assert.ErrorIs(t, err, ErrInvalidParameter)
}

func Test_ParseSentenceIntoExpenseWithFriend(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		expenseID = resources.ExpenseID(100)
		input     = "expected input"
		friendID  = 50
		autosave  = true
	)

	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := struct {
			Input    string `json:"input"`
			FriendID int    `json:"friend_id"`
			Autosave bool   `json:"autosave"`
		}{}

		rawBody, err := io.ReadAll(r.Body)
		require.NoError(err)

		err = json.Unmarshal(rawBody, &in)
		require.NoError(err)

		assert.Equal(input, in.Input)
		assert.Equal(friendID, in.FriendID)
		assert.Equal(autosave, in.Autosave)

		res := struct {
			Expense    resources.Expense `json:"expense"`
			Valid      bool              `json:"valid"`
			Confidence float64           `json:"confidence"`
		}{
			Expense: resources.Expense{
				ID: expenseID,
			},
			Valid:      true,
			Confidence: 1.0,
		}

		err = json.NewEncoder(w).Encode(res)
		require.NoError(err)
	}))
	defer cancel()

	ctx := context.Background()
	e, err := client.ParseSentenceIntoExpenseWithFriend(ctx, input, friendID, autosave)
	require.NoError(err)

	assert.Equal(expenseID, e.Expense.ID)
	assert.True(e.Valid)
}

func Test_ParseSentenceIntoExpenseWithFriend_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.ParseSentenceIntoExpenseWithFriend(ctx, "statement", 0, false)

		return err
	}

	doBasicErrorChecks(t, f)
}

func Test_ParseSentenceIntoExpenseWithGroup(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	const (
		expenseID = resources.ExpenseID(100)
		input     = "expected input"
		groupId   = 50
		autosave  = true
	)

	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := struct {
			Input    string `json:"input"`
			GroupID  int    `json:"group_id"`
			Autosave bool   `json:"autosave"`
		}{}

		rawBody, err := io.ReadAll(r.Body)
		require.NoError(err)

		err = json.Unmarshal(rawBody, &in)
		require.NoError(err)

		assert.Equal(input, in.Input)
		assert.Equal(groupId, in.GroupID)
		assert.Equal(autosave, in.Autosave)

		res := struct {
			Expense    resources.Expense `json:"expense"`
			Valid      bool              `json:"valid"`
			Confidence float64           `json:"confidence"`
		}{
			Expense: resources.Expense{
				ID: expenseID,
			},
			Valid:      true,
			Confidence: 1.0,
		}

		err = json.NewEncoder(w).Encode(res)
		require.NoError(err)
	}))
	defer cancel()

	ctx := context.Background()
	e, err := client.ParseSentenceIntoExpenseWithGroup(ctx, input, groupId, autosave)
	require.NoError(err)

	assert.Equal(expenseID, e.Expense.ID)
	assert.True(e.Valid)
}

func Test_ParseSentenceIntoExpenseWithGroup_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.ParseSentenceIntoExpenseWithGroup(ctx, "statement", 0, false)

		return err
	}

	doBasicErrorChecks(t, f)
}
