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

func Test_GetCategories_SanityCheck(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, categoriesSuccessResponse)
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

func Test_GetCategories_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		cats, err := client.GetCategories(ctx)
		assert.Len(t, cats, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}
