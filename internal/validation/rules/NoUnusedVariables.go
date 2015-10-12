package rules

import (
	lang "github.com/ng-vu/graphql-go/internal/language"
)

type noUnusedVariables struct {
}

func NoUnusedVariables() RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			switch node.(type) {
			case lang.OperationDefinition:
			case lang.VariableDefinition:
			case lang.Variable:
			case lang.FragmentSpread:
			}
			return nil
		},

		Leave: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			switch node.(type) {
			case lang.OperationDefinition:
			}
			return nil
		},
	}
}
