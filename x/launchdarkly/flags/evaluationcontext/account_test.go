package evaluationcontext_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/x/launchdarkly/flags/evaluationcontext"
	"github.com/cultureamp/ca-go/x/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

func TestNewAccount(t *testing.T) {
	t.Run("can create an account", func(t *testing.T) {
		account := evaluationcontext.NewAccount("my-account-id")
		assertAccountAttributes(t, account, "toggle.account.my-account-id")
	})

	t.Run("can create an account from context", func(t *testing.T) {
		user := request.AuthenticatedUser{
			CustomerAccountID: "123",
			RealUserID:        "456",
			UserID:            "789",
		}
		ctx := context.Background()

		ctx = request.ContextWithAuthenticatedUser(ctx, user)

		toggleAccount, err := evaluationcontext.AccountFromContext(ctx)
		require.NoError(t, err)
		assertAccountAttributes(t, toggleAccount, "toggle.account.123")
	})
}

func assertAccountAttributes(t *testing.T, account evaluationcontext.Account, accountID string) {
	t.Helper()

	ldUser, ok := account.RawUser().(lduser.User)
	require.True(t, ok, "should be castable to a LaunchDarkly user object")

	assert.Equal(t, accountID, ldUser.GetKey())
}
