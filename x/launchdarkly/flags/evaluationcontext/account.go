package evaluationcontext

import (
	"context"
	"errors"

	"github.com/cultureamp/ca-go/x/request"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

const (
	accountEntityPrefix = "account"
)

// Account is an evaluation context used for toggles.
type Account struct {
	key string

	ldUser lduser.User
}

func (a Account) ToLDToggleUser() lduser.User {
	return a.ldUser
}

// NewAccount returns an evaluation context representing a customer account.
// accountID is the ID of a customer account, and will generally be an
// "account_aggregate_id".
// This is used for toggles only.
func NewAccount(accountID string) Account {
	a := &Account{
		key: ldKey(toggleContextKind, accountEntityPrefix, accountID),
	}

	userBuilder := lduser.NewUserBuilder(a.key)
	a.ldUser = userBuilder.Build()

	return *a
}

// RawUser returns the wrapped LaunchDarkly user object. The return value should
// be casted to an lduser.User object.
func (a Account) RawUser() interface{} {
	return a.ldUser
}

// AccountFromContext extracts the account aggregate ID from the context. This
// value is used to create a new Account object. An error is returned if
// the account identifier is not present in the context.
func AccountFromContext(ctx context.Context) (Account, error) {
	authenticatedUser, ok := request.AuthenticatedUserFromContext(ctx)
	if !ok {
		return Account{}, errors.New("no AuthenticatedUser in supplied context")
	}

	return NewAccount(authenticatedUser.CustomerAccountID), nil
}
