package rules

import "fmt"

func unknownTypeMessage(typ interface{}) string {
	return fmt.Sprintf(`Unknown type "%v".`, typ)
}

func KnownTypeNames(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			typeName := node.Name.Value
			typ := context.GetSchema().GetType(typeName)

			if !typ {
				return lang.NewGraphQLError{
					unknownTypeMessage(typeName),
					[]lang.INode{node},
				}
			}
		},
	}
}
