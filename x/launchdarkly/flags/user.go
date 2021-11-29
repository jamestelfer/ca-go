package flags

import (
	"context"
	"errors"

	"github.com/cultureamp/ca-go/x/request"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

const (
	anonymousUser                    = "ANONYMOUS_USER"
	userAttributeAccountAggregateID  = "accountAggregateID"
	userAttributeRealUserAggregateID = "realUserAggregateID"
)

// User wraps the LaunchDarkly user object.
type User struct {
	effectiveUserAggregateID string
	realUserAggregateID      string
	accountAggregateID       string

	ldUser lduser.User
}

// UserOption are functions that can be supplied to configure a new user with
// additional attributes.
type UserOption func(*User)

// WithAccountAggregateID configures the user with the given account aggregate ID.
func WithAccountAggregateID(id string) UserOption {
	return func(u *User) {
		u.accountAggregateID = id
	}
}

// WithRealUserAggregateID configures the user with the given real user aggregate ID.
func WithRealUserAggregateID(id string) UserOption {
	return func(u *User) {
		u.realUserAggregateID = id
	}
}

// NewAnonymousUser returns a user object suitable for use in unauthenticated
// requests or requests with no access to user identifiers.
func NewAnonymousUser() User {
	return User{
		ldUser: lduser.NewAnonymousUser(anonymousUser),
	}
}

// NewUser returns a new user object with the given effective user aggregate ID
// and options.
func NewUser(effectiveUserAggregateID string, opts ...UserOption) User {
	u := &User{
		effectiveUserAggregateID: effectiveUserAggregateID,
	}

	for _, opt := range opts {
		opt(u)
	}

	userBuilder := lduser.NewUserBuilder(effectiveUserAggregateID)
	userBuilder.Custom(
		userAttributeAccountAggregateID,
		ldvalue.String(u.accountAggregateID))
	userBuilder.Custom(
		userAttributeRealUserAggregateID,
		ldvalue.String(u.realUserAggregateID))
	u.ldUser = userBuilder.Build()

	return *u
}

// RawUser returns the wrapped LaunchDarkly user object. The return value should
// be casted to an lduser.User object.
func (u User) RawUser() interface{} {
	return u.ldUser
}

// UserFromContext extracts the effective user aggregate ID, real user aggregate
// ID, and account aggregate ID from the context. These values are used to
// create a new User object. An error is returned if user identifiers are not
// present in the context.
func UserFromContext(ctx context.Context) (User, error) {
	authenticatedUser, ok := request.AuthenticatedUserFromContext(ctx)
	if !ok {
		return User{}, errors.New("no AuthenticatedUser in supplied context")
	}

	return NewUser(
		authenticatedUser.EffectiveUserAggregateID,
		WithAccountAggregateID(authenticatedUser.AccountAggregateID),
		WithRealUserAggregateID(authenticatedUser.RealUserAggregateID)), nil
}
