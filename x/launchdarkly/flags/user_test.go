package flags_test

import (
	"context"
	"testing"

	"github.com/cultureamp/ca-go/request"
	"github.com/cultureamp/ca-go/x/launchdarkly/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

func TestNewUser(t *testing.T) {
	t.Run("can create an anonymous user", func(t *testing.T) {
		user := flags.NewAnonymousUser()
		assertUserAttributes(t, user, "ANONYMOUS_USER", "", "")
	})

	t.Run("can create an identified user", func(t *testing.T) {
		user := flags.NewUser("not-a-uuid")
		assertUserAttributes(t, user, "not-a-uuid", "", "")

		user = flags.NewUser(
			"not-a-uuid",
			flags.WithAccountAggregateID("not-a-uuid"),
			flags.WithRealUserAggregateID("not-a-uuid"))
		assertUserAttributes(t, user, "not-a-uuid", "not-a-uuid", "not-a-uuid")
	})

	t.Run("can create a user from context", func(t *testing.T) {
		user := request.NewAuthenticatedUser("123", "456", "789")
		ctx := context.Background()

		ctx = request.ContextWithAuthenticatedUser(ctx, user)

		flagsUser, err := flags.UserFromContext(ctx)
		require.NoError(t, err)
		assertUserAttributes(t, flagsUser, "789", "456", "123")
	})
}

func assertUserAttributes(t *testing.T, user flags.User, effectiveUserAggregateID, realUserAggregateID, accountAggregateID string) {
	t.Helper()

	ldUser, ok := user.RawUser().(lduser.User)
	require.True(t, ok, "should be castable to a LaunchDarkly user object")

	assert.Equal(t, effectiveUserAggregateID, ldUser.GetKey())
	assert.Equal(t, realUserAggregateID, ldUser.GetAttribute("realUserAggregateID").StringValue())
	assert.Equal(t, accountAggregateID, ldUser.GetAttribute("accountAggregateID").StringValue())
}
