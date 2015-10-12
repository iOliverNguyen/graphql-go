package rules

import (
	"fmt"

	lang "github.com/ng-vu/graphql-go/internal/language"
	util "github.com/ng-vu/graphql-go/internal/utilities"
	"github.com/ng-vu/graphql-go/validation"
)

func badValueMessage(argName, typ, value interface{}) string {
	return fmt.Sprintf(
		`Argument "%v" expected type "%v" but got: ${%v}.`,
		argName, typ, value)
}

func ArgumentsOfCorrectType(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			if argAST, ok := node.(lang.Argument); ok {
				argDef := context.GetArgument()
				if argDef != nil && !util.IsValidLiteralValue(argDef.Type, argAST.Value) {
					return lang.NewGraphQLError(
						badValueMessage(argAST.Name.Value, argDef.Type, lang.Print(argAST.Value)),
						[]lang.INode{argAST.Value})
				}
			}
			return nil
		},
	}
}
