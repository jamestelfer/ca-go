package evaluationcontext

import (
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

type Context interface {
	ToLDUser() lduser.User
}
