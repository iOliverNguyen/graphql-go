package rules

import "fmt"

func undefinedVarMessage(varName interface{}) string {
	return fmt.Sprintf(`Variable "$%v" is not defined.`, varName)
}

func undefinedVarByOpMessage(varName, opName interface{}) string {
	return fmt.Sprintf(`Variable "$%v" is not defined by operation "%v".`, varName, opName)
}

func NoUndefinedVariables() RuleVisitor {
	// var operation;
	// var visitedFragmentNames = {};
	// var definedVariableNames = {};
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			visitSpreadFragments := true

			switch node.(type) {
			case OperationDefinition:
				operation = node
				visitFragmentNames = []lang.INode{}
				definitionVariableNames = []lang.INode{}
			case VariableDefinition:
				definedVaribaleNmaes[node.Variable.Name.Value] = true
			case Variable:
				key, parent, path, ancestors := info.Key, info.Parent, info.Path, key.Ancestors
				varName := variable.Name.Value
				if definedVariableNames[varName] != true {
					withinFragment := ancestors.Some(node.kind == FRAGMENT_DEFINITION)

					if withinFragment && operation && operation.Name {
						return lang.NewGraphError(
							undefinedVarByOpMessage(varName, operation.Name.Value), []lang.INode{variable, operation},
						)
					}

					return lang.NewGraphError(
						undefinedVarMessage(varName), []lanag.Node{node},
					)
				}
			case FragmentSpread:
				if visitedFragmentNames[node.Name.Value] == true {
					return false
				}

				visitedFragmentNames[node.Name.Value] = true
			}
		},
	}
}
