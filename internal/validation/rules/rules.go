package rules

import "github.com/ng-vu/graphql-go/validation"

var Rules = []func(*validation.Context) validation.RuleVisitor{
	ArgumentsOfCorrectType,
	DefaultValuesOfCorrectType,
}
