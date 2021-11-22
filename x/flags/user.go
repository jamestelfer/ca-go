package flags

import (
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

const (
	anonymousUser                   = "ANONYMOUS_USER"
	userAttributeAccountAggregateID = "accountAggregateId"
)

type User struct {
	userAggregateID    string
	accountAggregateID string

	ldUser lduser.User
}

type UserOption func(*User)

func WithAccountAggregateID(id string) UserOption {
	return func(u *User) {
		u.accountAggregateID = id
	}
}

func NewAnonymousUser() User {
	return User{
		ldUser: lduser.NewAnonymousUser(anonymousUser),
	}
}

func NewUser(userAggregateID string, opts ...UserOption) User {
	u := &User{
		userAggregateID: userAggregateID,
	}

	for _, opt := range opts {
		opt(u)
	}

	u.ldUser = lduser.NewUserBuilder(userAggregateID).Custom(
		userAttributeAccountAggregateID,
		ldvalue.String(u.accountAggregateID)).Build()

	return *u
}

func (u User) RawUser() interface{} {
	return u.ldUser
}
