package rules

import "fmt"

func unknownFragmentMessage(fragName interface{}) string {
	return fmt.Sprintf(`Unknown fragment "%v".`, fragName)
}

func KnownFragmentNames(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			fragmentName := node.Name.Value
			fragment := context.GetFragment(fragmentName)

			if !fragment {
				return lang.NewGraphQLError(
					unknownFragmentMessage(fragmentName),
					[]lang.INode{node.Name},
				)
			}
		},
	}
}
