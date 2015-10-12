package rules

import "fmt"

func undefinedFieldMessage(fieldName, typ interface{}) string {
	return fmt.Sprintf(`Cannot query field "%v" on "%v".`, fieldName, typ)
}

func FieldsOnCorrectType(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			if node, ok := node.(lang.Field); ok {
				typ := context.GetParentType()
				if typ {
					fieldDef := context.GetFieldDef()
					if !fieldDef {
						return lang.NewGraphQLError(
							undefinedFieldMessage(node.Value, typ.name),
							[]lang.INode{node},
						)
					}
				}
			}
			return nil
		},
	}
}
