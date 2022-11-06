package splitwise

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IfTokenIsDefinedItShouldBeInHeader(t *testing.T) {
	const testToken = "testtoken"

	client, cancel := testClientWithHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader != fmt.Sprintf("Bearer %s", testToken) {
			assert.Failf(t, "authorization header is not correct", "authorization header is [%s], should be [Bearer %s]", authHeader, testToken)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer cancel()

	ctx := context.Background()

	client.Token = testToken

	res, err := client.do(ctx, http.MethodGet, "/generic", nil, nil)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}
