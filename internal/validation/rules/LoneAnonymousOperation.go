package rules

func anonOperationNotAloneMessage() string {
	return fmt.SPrintf(`This anonymous operation must be the only defined operation.`)
}

func LoneAnonymousOperation() RuleVisitor {
	operationCount := 0
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GrahpQLError {
			switch node.(type) {
			case Document:
				operationCount = len(node.Defifinitons.Filter(func(definition interface{}) bool {
					definition.Kind() == "OperationDefinition"
				}))
			case OperationDefinition:
				if !node.Name && operationCount > 1 {
					return lang.NewGraphQLError(
						anonOperationNotAloneMessage(),
						[]lang.INode{node},
					)
				}
			}
		},
	}
}
