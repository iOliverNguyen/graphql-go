package rules

import "fmt"

func duplicateOperationNameMessage(operationName interface{}) string {
	return fmt.Sprintf(`There can only be one operation named "%v".`, operationName)
}

func UniqueOperationNames() RuleVisitor {
	knownOperationNames := lang.INode{}

	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			operationName := node.Name
			if operationName {
				if knownOperationNames(operationName.Value) {
					return lang.NewGraphQLError(
						duplicateOperationNameMessage(operationName.Value),
						[]lang.INode{knownOperationNames(operationName.Value), operationName},
					)
				}

				knownOperationNames[operationName.Value] = operationName
			}
		},
	}
}
