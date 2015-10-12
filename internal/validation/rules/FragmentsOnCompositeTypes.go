package rules

import "fmt"

func inlineFragmentOnNonCompositeErrorMessage(typ interface{}) string {
	return fmt.Sprintf(`Fragment cannot condition on non composite type "%v".`, typ)
}

func fragmentOnNonCompositeErrorMessage(fragName, typ interface{}) string {
	return fmt.Sprintf(`Fragment "${fragName}" cannot condition on non composite type "%v".`, typ)
}

func FragmentsOnCompositeTypes(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			switch node.(type) {
			case lang.InlineFragment:
				typ := context.GetType()
				if typ && !util.IsCompositeType(typ) {
					return lang.NewGraphQLError(
						inlineFragmentOnNonCompositeErrorMessage(printer.Print(node.TypeCondition)),
						[]lang.INode{node.TypeCondition},
					)
				}
			case lang.FragmentDefinition:
				typ := context.GetType()
				if typ && !util.IsCompositeType(typ) {
					return lang.NewGraphQLError(
						fragmentOnNonCompositeErrorMessage(node.Name, printer.Print(node.TypeCondition)),
						[]lang.INode{node.TypeCondition},
					)
				}
			}

			return nil
		},
	}
}
