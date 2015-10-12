package rules

import "fmt"

func duplicateFragmentNameMessage(fragName interface{}) string {
	return fmt.Sprintf(`There can only be one fragment named "%v".`, fragName)
}

func UniqueFragmentNames() RuleVisitor {
	knownFragmentNames := []lang.INode{}
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			fragmentName := node.Name.Value
			if knownFragmentNames[fragmentName] {
				return lang.NewGraphQLError{
					duplicateFragmentNameMessage(fragName),
					[]lang.INode{knownFragmentNames[fragmentName], node.Name},
				}
			}

			knownFragmentNames[fragmentName] = node.Name
		},
	}
}
