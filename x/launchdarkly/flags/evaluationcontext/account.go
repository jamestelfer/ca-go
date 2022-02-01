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

// Account is a type of ToggleContext, representing the identifier of a customer
// account to evaluate a toggle against.
type Account struct {
	key string

	ldUser lduser.User
}

func (a Account) ToLDToggleUser() lduser.User {
	return a.ldUser
}

// NewAccount returns a new account object with the given account ID.
// accountID is the ID of a customer account, and will generally
// be an "account_aggregate_id".
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
