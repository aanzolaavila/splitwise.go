package splitwise

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
)

const getExpensesComments200Response = `
{
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
`

func Test_GetExpenseComments_SanityChecks(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, getExpensesComments200Response)
	defer cancel()

	ctx := context.Background()
	comments, err := client.GetExpenseComments(ctx, 855870953)
	assert.NoError(t, err)

	assert.Len(t, comments, 1)

	com := comments[0]
	assert.Equal(t, resources.CommentID(79800950), com.ID)
	assert.Equal(t, "System", com.CommentType)
	assert.Equal(t, "ExpenseComment", com.RelationType)
	assert.Equal(t, resources.Identifier(855870953), com.RelationID)
}

func Test_GetExpenseComments_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		comms, err := client.GetExpenseComments(ctx, 0)
		assert.Len(t, comms, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}

const createExpenseComment200Response = `
{
  "comment": {
    "id": 79800950,
    "content": "Does this include the delivery fee?",
    "comment_type": "User",
    "relation_type": "ExpenseComment",
    "relation_id": 5123,
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
}
`

func Test_CreateExpenseComments_SanityCheck(t *testing.T) {
	const expenseId = 5123
	const content = "Does this include the delivery fee?"

	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		input := struct {
			ExpenseId int    `json:"expense_id"`
			Content   string `json:"content"`
		}{}

		rawBody, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()

		err = json.Unmarshal(rawBody, &input)
		assert.NoError(t, err)

		assert.Equal(t, expenseId, input.ExpenseId)
		assert.Equal(t, content, input.Content)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(createExpenseComment200Response))
	}))
	defer cancel()

	ctx := context.Background()
	com, err := client.CreateExpenseComment(ctx, expenseId, content)
	assert.NoError(t, err)

	assert.Equal(t, resources.Identifier(expenseId), com.RelationID)
	assert.Equal(t, content, com.Content)
}

func Test_CreateExpenseComments_ShouldFailIfEmptyContent(t *testing.T) {
	httpClient := httpClientStub{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			const msg = "the client should have never been called"
			t.Fatalf(msg)
			return nil, errors.New(msg)
		},
	}

	client := Client{
		HttpClient: httpClient,
	}

	ctx := context.Background()
	_, err := client.CreateExpenseComment(ctx, 0, "")
	assert.Error(t, err)
}

func Test_CreateExpenseComments_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.CreateExpenseComment(ctx, 0, "some content")

		return err
	}

	doBasicErrorChecks(t, f)
}

const deleteExpenseComment200Response = `
{
  "comment": {
    "id": 79800950,
    "content": "Does this include the delivery fee?",
    "comment_type": "User",
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
}
`

func Test_DeleteExpenseComment_SanityCheck(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodPost, deleteExpenseComment200Response)
	defer cancel()

	ctx := context.Background()
	const commId = 79800950
	com, err := client.DeleteExpenseComment(ctx, commId)
	assert.NoError(t, err)

	assert.Equal(t, resources.CommentID(commId), com.ID)
	assert.Equal(t, "ExpenseComment", com.RelationType)
}

func Test_DeleteExpenseComment_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		_, err := client.DeleteExpenseComment(ctx, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}
