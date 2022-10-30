package splitwise

import (
	"context"
	"net/http"
	"testing"

	"github.com/aanzolaavila/splitwise.go/resources"
	"github.com/stretchr/testify/assert"
)

const categoriesSuccessResponse = `
{
  "categories": [
    {
      "id": 1,
      "name": "Utilities",
      "icon": "https://s3.amazonaws.com/splitwise/uploads/category/icon/square/utilities/cleaning.png",
      "icon_types": {
        "slim": {
          "small": "http://example.com",
          "large": "http://example.com"
        },
        "square": {
          "large": "http://example.com",
          "xlarge": "http://example.com"
        }
      },
      "subcategories": [
        {
          "id": 48,
          "name": "Cleaning",
          "icon": "https://s3.amazonaws.com/splitwise/uploads/category/icon/square/utilities/cleaning.png",
          "icon_types": {
            "slim": {
              "small": "http://example.com",
              "large": "http://example.com"
            },
            "square": {
              "large": "http://example.com",
              "xlarge": "http://example.com"
            }
          }
        }
      ]
    }
  ]
}`

func Test_Categories_SanityCheck(t *testing.T) {
	client, cancel := testClient(200, categoriesSuccessResponse)
	defer cancel()

	ctx := context.Background()

	categories, err := client.GetCategories(ctx)
	assert.NoError(t, err)

	assert.Len(t, categories, 1)

	cat := categories[0]
	assert.Equal(t, resources.CategoryID(1), cat.ID)

	assert.Len(t, cat.Subcategories, 1)

	subcat := cat.Subcategories[0]
	assert.Equal(t, resources.CategoryID(48), subcat.ID)
}

func Test_Categories_FaultyClient(t *testing.T) {
	client, expectedErr := testClientWithFaultyResponse()

	ctx := context.Background()

	categories, err := client.GetCategories(ctx)
	assert.ErrorIs(t, err, expectedErr)
	assert.Len(t, categories, 0)
}

func Test_Categories_FaultyBodyShouldFail(t *testing.T) {
	client, expectedErr, cancel := testClientWithFaultyResponseBody(t, http.StatusOK)
	defer cancel()

	ctx := context.Background()

	cats, err := client.GetCategories(ctx)
	assert.ErrorIs(t, err, expectedErr)
	assert.Len(t, cats, 0)
}
