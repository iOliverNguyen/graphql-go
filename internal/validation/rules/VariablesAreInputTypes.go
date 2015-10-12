package rules

import "fmt"

func nonInputTypeOnVarMessage(variableName, typeName interface{}) string {
	return fmt.Sprintf(`Variable "$%v" cannot be non-input type "%v".`, variableName, typeName)
}

func VariablesAreInputTypes(context *validation.Context) RuleVisitor {
	return RuleVisitor{
		Enter: func(node lang.INode, info lang.VisitInfo) *lang.GraphQLError {
			typ := typeFromAST(context.GetSchema(), node.Type)

			if typ && !isInputType(typ) {
				variableName := node.Variable.Name.Value
				return lang.GraphQLError(
					nonInputTypeOnVarMessage(variableName, printer.Print(typeName)),
					[]lang.INode{node.Type},
				)
			}
		},
	}
}
