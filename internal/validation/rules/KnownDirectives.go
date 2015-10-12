package rules

import "fmt"

func unknownDirectiveMessage(directiveName interface{}) string {
	return fmt.Sprintf(`Unknown directive "%v".`, directiveName)
}

func misplacedDirectiveMessage(directiveName, placement interface{}) {
	return fmt.Sprintf(`Directive "%v" may not be used on "%v".`, directiveName, placement)
}

func KnownDirectives(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			directiveDef := util.Find(
				context.GetSchema().GetDirectives(),
				func(def interface{}) bool {
					return def.Name == node.Name.Value
				},
			)

			if !directiveDef {
				return lang.NewGraphQLError(
					unknownDirectiveMessage(node.Name.Value),
					[]lang.INode{node},
				)
			}

			key, parent, path, ancestors := info.Key, info.Parent, info.Path, info.ancestors

			appliedTo := ancestors[len(ancestors)-1]

			if appliedTo.Kind() == lang.OPERATION_DEFINITION && !directiveDef.OnOperation {
				return lang.NewGraphQLError(
					misplacedDirectiveMessage(node.Name.Value, "operation"),
					[]lang.INode{node},
				)
			}

			if appliedTo.Kind == lang.FIELD && directiveDef.OnField {
				return lang.NewGraphQLError(
					misplacedDirectiveMessage(node.Name.Value, "field"),
					[]lang.INode{node},
				)
			}

			if (apploed.Kind() == lang.FRAGMENT_SPREAD || applied.Kind() == lang.INLINE_FRAGMENT || appliedTo.Kind() == lang.FRAGMENT_DEFINITION) && !directiveDef.OnFragment {
				return lang.NewGraphQLError(
					misplacedDirectiveMessage(node.Name.Value, "fragment"),
					[]lang.INode{node},
				)
			}
		},
	}
}
