package rules

import "fmt"

func unusedFragMessage(fragName interface{}) string {
	return fmt.Sprintf(`Fragment "%v" is never used.`, fragName)
}

func NoUnusedFragments() RuleVisitor {
	// var fragmentDefs = [];
	// var spreadsWithinOperation = [];
	// var fragAdjacencies = {};
	// var spreadNames = {};

	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			switch node.(type) {
			case OperationDefinition:
			case FragmentDefinition:
			case FragmentSpread:
			}
		},
		Leave: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			switch node.(type) {
			case Document:
			}
		},
	}
}
