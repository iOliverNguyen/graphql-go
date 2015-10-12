package rules

import (
	"fmt"

	lang "github.com/ng-vu/graphql-go/internal/language"
	util "github.com/ng-vu/graphql-go/internal/utilities"
	typs "github.com/ng-vu/graphql-go/types"
	"github.com/ng-vu/graphql-go/validation"
)

func defaultForNonNullArgMessage(varName, typ, guessType interface{}) string {
	return fmt.Sprintf(
		`Variable "$%v" of type "%v" is required and will not use the default value. Perhaps you meant to use type "%v".`,
		varName, typ, guessType)
}

func badValueForDefaultArgMessage(varName, typ, value interface{}) string {
	return fmt.Sprintf(
		`Variable "$%v" of type "%v" has invalid default value: %v.`,
		varName, typ, value)
}

/**
 * Variable default values of correct type
 *
 * A GraphQL document is only valid if all variable default values are of the
 * type expected by their definition.
 */
func DefaultValuesOfCorrectType(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			if varDefAST, ok := node.(lang.VariableDefinition); ok {
				name := varDefAST.Variable.Name.Value
				defaultValue := varDefAST.DefaultValue
				typ := context.GetInputType()
				if typ, ok := typ.(typs.GraphQLNonNull); ok {
					return lang.NewGraphQLError(
						defaultForNonNullArgMessage(name, typ, typ.OfType),
						[]lang.INode{defaultValue})
				}
				if typ != nil && !util.IsValidLiteralValue(typ, defaultValue) {
					return lang.NewGraphQLError(
						badValueForDefaultArgMessage(name, typ, lang.Print(defaultValue)),
						[]lang.INode{defaultValue})
				}
			}
			return nil
		},
	}
}
