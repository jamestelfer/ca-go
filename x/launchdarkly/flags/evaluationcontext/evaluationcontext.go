package evaluationcontext

import (
	"fmt"

	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
)

type entity string

const (
	entityUser entity = "user"
)

const attributeEntityType = "entityType"

type Context interface {
	ToLDUser() lduser.User
}

func prefixEntity(entity entity, key string) string {
	return fmt.Sprintf("%s.%s", entity, key)
}
