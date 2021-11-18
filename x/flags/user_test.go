package flags_test

import (
	"testing"

	"github.com/cultureamp/ca-go/x/flags"
)

func TestNewUser(t *testing.T) {
	t.Run("can create an anonymous user", func(t *testing.T) {
		_ = flags.NewAnonymousUser()
	})

	t.Run("can create an identified user", func(t *testing.T) {
		_ = flags.NewUser("not-a-uuid")
		_ = flags.NewUser("not-a-uuid", flags.WithAccountAggregateID("not-a-uuid"))
	})
}
