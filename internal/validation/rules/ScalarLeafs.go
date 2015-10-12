package rules

import "fmt"

func noSubselectionAllowedMessage(field, typ interface{}) string {
	return fmt.Sprintf(`Field "%v" of type "%v" must not have a sub selection.`, field, typ)
}

func requiredSubselectionMessage(field, typ interface{}) string {
	return fmt.Sprintf(`Field "%v" of type "%v" must have a sub selection.`, field, typ)
}

func ScalarLeafs(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			typ := context.GetType()
			if typ {
				if isLeafType(typ) {
					if node.SelectionSet {
						return lang.NewGraphQLError(
							noSubselectionAllowedMessage(node.Name.Value, typ),
							[]lang.INode{node.SelectionSet},
						)
					}
				} else if !node.SecletionSet {
					return lang.GraphQLError(
						requiredSubselectionMessage(node.Name.Value, typ),
						[]lang.INode{node},
					)
				}
			}
		},
	}
}
