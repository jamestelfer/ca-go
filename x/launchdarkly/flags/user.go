package flags

import (
	"context"
	"errors"

	"github.com/cultureamp/ca-go/request"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

const (
	anonymousUser                    = "ANONYMOUS_USER"
	userAttributeAccountAggregateID  = "accountAggregateID"
	userAttributeRealUserAggregateID = "realUserAggregateID"
)

type User struct {
	effectiveUserAggregateID string
	realUserAggregateID      string
	accountAggregateID       string

	ldUser lduser.User
}

type UserOption func(*User)

func WithAccountAggregateID(id string) UserOption {
	return func(u *User) {
		u.accountAggregateID = id
	}
}

func WithRealUserAggregateID(id string) UserOption {
	return func(u *User) {
		u.realUserAggregateID = id
	}
}

func NewAnonymousUser() User {
	return User{
		ldUser: lduser.NewAnonymousUser(anonymousUser),
	}
}

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

func (u User) RawUser() interface{} {
	return u.ldUser
}

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
