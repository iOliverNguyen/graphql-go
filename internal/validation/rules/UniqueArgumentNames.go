package rules

import "fmt"

func duplicateArgMessage(argName interface{}) string {
	return fmt.Sprintf(`There can be only one argument named "%v".`, argName)
}

func UniqueArgumentNames() RuleVisitor {
	knownArgNames := []lang.INode{}
	return RuleVisitor{
		Enter: func(node lang.INode, info *lang.VisitInfo) *lang.GraphGLError {
			switch node.(type) {
			case Field:
				knownArgNames = []lang.INode{}
			case Directive:
				knownArgNames = []lang.INode{}
			case Argument:
				argName := node.Name.Value
				if knownArgNames[argName] {
					return lang.GraphQLError{
						duplicateArgMessage(argName),
						[]lang.INode{knownArgNames[argName], node.Name},
					}
				}

				knownArgNames[argName] = node.Name
			}
		},
	}
}
