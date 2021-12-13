package evaluationcontext

import (
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

var (
	surveyAttributeType  = prefixEntity(entitySurvey, "type")
	surveyAttributeState = prefixEntity(entitySurvey, "state")
)

type Survey struct {
	surveyID    string
	surveyType  string
	surveyState string

	ldUser lduser.User
}

type SurveyOption func(*Survey)

func WithType(surveyType string) SurveyOption {
	return func(s *Survey) {
		s.surveyType = prefixEntity(entitySurvey, surveyType)
	}
}

func WithState(surveyState string) SurveyOption {
	return func(s *Survey) {
		s.surveyState = prefixEntity(entitySurvey, surveyState)
	}
}

func NewSurvey(surveyID string, opts ...SurveyOption) Survey {
	s := &Survey{
		surveyID: prefixEntity(entitySurvey, surveyID),
	}

	for _, opt := range opts {
		opt(s)
	}

	userBuilder := lduser.NewUserBuilder(s.surveyID)
	userBuilder.Custom(surveyAttributeType, ldvalue.String(s.surveyType))
	userBuilder.Custom(surveyAttributeState, ldvalue.String(s.surveyState))
	s.ldUser = userBuilder.Build()

	return *s
}

func (s Survey) ToLDUser() lduser.User {
	return s.ldUser
}
