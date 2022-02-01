package evaluationcontext

import (
	"fmt"

	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

const (
	flagContextKind   = "flag"
	toggleContextKind = "toggle"
)

// FlagContext represents a set of attributes which a flag is evaluated against. The
// only context supported now is User.
type FlagContext interface {
	// ToLDFlagUser transforms the context implementation into an LDUser object that can
	// be understood by LaunchDarkly when evaluating a flag.
	ToLDFlagUser() lduser.User
}

// ToggleContext represents a set of attributes which a product toggle is evaluated
// against. The only context supported now is Account.
type ToggleContext interface {
	ToLDToggleUser() lduser.User
}

// ldKey returns a formatted string to be used as the `key` value of the lduser.User
// sent to LaunchDarkly. The key is comprised of three components:
// 	1. The kind of evaluation context the key belongs to. A `flag` or `toggle`.
// 	2. The entity prefix, for example `user` or `account`.
//	3. The unique identifier for the entity, e.g. a user ID for the User evaluation
//     context.
func ldKey(contextType, prefix, key string) string {
	return fmt.Sprintf("%s.%s.%s", contextType, prefix, key)
}
