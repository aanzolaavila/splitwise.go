package splitwise

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const currenciesSuccessResponse = `
{
  "currencies": [
    {
      "currency_code": "BRL",
      "unit": "R$"
    }
  ]
}
`

func Test_GetCurrencies_SanityCheck(t *testing.T) {
	client, cancel := testClient(t, http.StatusOK, http.MethodGet, currenciesSuccessResponse)
	defer cancel()

	ctx := context.Background()

	curs, err := client.GetCurrencies(ctx)
	require.NoError(t, err)

	require.Len(t, curs, 1)

	cur := curs[0]
	assert.Equal(t, "BRL", cur.CurrencyCode)
	assert.Equal(t, "R$", cur.Unit)
}

func Test_GetCurrencies_BasicErrorTests(t *testing.T) {
	f := func(client Client) error {
		ctx := context.Background()
		curs, err := client.GetCurrencies(ctx)
		assert.Len(t, curs, 0)

		return err
	}

	doBasicErrorChecks(t, f)
}
