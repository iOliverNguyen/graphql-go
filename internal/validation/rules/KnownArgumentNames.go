package rules

import "fmt"

func unknownArgMessage(argName, fieldName, typ interface{}) string {
	return fmt.Sprintf(`Unknown argument "%v" on field "%v" of type "%v".`, argName, fieldName, typ)
}

func unknownDirectiveArgMessage(argName, directiveName interface{}) string {
	return fmt.Sprintf(`Unknown argument "%v" on directive "@%v".`, argName, directiveName)
}

func KnownArgumentNames(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			if node, ok := node.(lang.Argument); ok {
				key, parent, path, ancestors := info.Key, info.Parent, info.Path, info.ancestors
				argumentOf := ancestors[len(ancestors)-1]
				if argumentOf.Kind() == lang.INodeKind("Field") {
					fieldDef := context.GetFieldDef()
					if fieldDef, ok := fieldDef.(lang.FieldDefinition); ok {
						fieldArgDef := util.Find(fieldDef.Arguments, func(item interface{}) bool {
							if item, ok := item.(lang.InputValueDefinition); ok {
								return item.Name == node.Name.Value
							}
							return false
						})
						if fieldArgDef, ok := fieldArgDef.(lang.FieldArgDef); ok {
							parentType := context.GetParentType()
							util.Invariant(parentType)
							return lang.NewGraphQLError(
								unknownArgMessage(node.Name.Value, fieldDef.Name, parentType.Name),
								[]lang.INode{node},
							)
						}
					}
				} else if argumentOf.Kind() == lang.INodeKind("Directive") {
					directive := context.GetDirectvie()
					if directive {
						directiveArgDef := util.Find(directive.Arg, func(item interface{}) bool {
							if item.Name == node.Name.Value {
								return true
							}

							return false
						})
					}

          if !directiveArgDef {
            return lang.NewGraphQLError(
              unknownDirectiveArgMessage(node.Name.Value, directive.Name)
            )
          }
				}
			}
			return nil
		},
	}
}
