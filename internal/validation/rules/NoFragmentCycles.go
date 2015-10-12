package rules

import "fmt"

func cycleErrorMessage(fragName interface{}, spreadNames []interface{}) string {
	via := ""
	if len(spreadNames) > 0 {
		via = strings.join(spreadNames, ",")
	}

	return fmt.Sprintf(`Cannot spread fragment "%v" within itself %v.`, fragName)
}

func NoFragmentCycles(context *validation.Context) RuleVisitor {
	definitions := context.GetDocument().Definitions
	spreadsInFragment = definitions.Reduce(func(m, node) map[string]interface{} {
		if node.Kind() == lang.FRAGMENT_DEFINITION {
			m[node.Name.Value] = gatherSpreads(node)
		}

		return m
	}, o)

	knownToLeadToCycle = Set()
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			errors := []Error{}
			initialName := noed.Name.Value

			spreadPath := []lang.INode{}

			detectCycleRecursive := func(fragementName interface{}) {
				spreadNodes = spreadsInFragement[fragmentName]
				for i := 0; i < len(spreadNodes); i++ {
					spreadNode := spreadNodes[i]
					if knownToLeadToCycle.Has(spreadNode) {
						continue
					}

					if spreadNode.Name.Value == initialName {
						cyclePath := append(spreadPath, spreadNode)

						for j := 0; j < len(cyclePath); j++ {
							knowToLoadToCyle.Add(cyclePath[j])
							errors = append(errors, lang.NewGraphQLError(
								cycleErrorMessage(intitalName, speadPath)),
								cyclePath)
						}
						continue
					}

					if spreadPath.Some(spread == spreadNode) {
						continue
					}

					spreadPath = append(spreadPath, spreadNode)

					detectCycleRecursive(initialName)

					if len(errors) > 0 {
						return errors
					}
				}
			}
		},
	}
}

func gatherSpreads(node []lang.INode) []lang.INode {
	spreadNodes = []lang.INode{}
	// visit(node, )
	return spreadNodes
}
