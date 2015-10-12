package language

import (
	"unsafe"
)

type QueryKeyMap struct {
	Name struct{}

	Document            struct{ Definitions bool }
	OperationDefinition struct{ Name, VariableDefinitions, Directives, SelectionSet bool }
	VariableDefinition  struct{ Variable, IType, DefaultValue bool }
	Variable            struct{ Name bool }
	SelectionSet        struct{ Selections bool }
	Field               struct{ Alias, Name, Arguments, Directives, SelectionSet bool }
	Argument            struct{ Name, IValue bool }

	FragmentSpread     struct{ Name, Directives bool }
	InlineFragment     struct{ TypeCondition, Directives, SelectionSet bool }
	FragmentDefinition struct{ Name, TypeCondition, Directives, SelectionSet bool }

	IntValue     struct{}
	FloatValue   struct{}
	StringValue  struct{}
	BooleanValue struct{}
	EnumValue    struct{}
	ListValue    struct{ Values bool }
	ObjectValue  struct{ Fields bool }
	ObjectField  struct{ Name, IValue bool }

	Directive struct{ Name, Arguments bool }

	NamedType   struct{ Name bool }
	ListType    struct{ IType bool }
	NonNullType struct{ IType bool }

	ObjectTypeDefinition      struct{ Name, Interfaces, Fields bool }
	FieldDefinition           struct{ Name, Arguments, IType bool }
	InputValueDefinition      struct{ Name, IType, DefaultValue bool }
	InterfaceTypeDefinition   struct{ Name, Fields bool }
	UnionTypeDefinition       struct{ Name, Types bool }
	ScalarTypeDefinition      struct{ Name bool }
	EnumTypeDefinition        struct{ Name, Values bool }
	EnumValueDefinition       struct{ Name bool }
	InputObjectTypeDefinition struct{ Name, Fields bool }
	TypeExtensionDefinition   struct{ IDefinition bool }
}

var QueryDocumentKeys QueryKeyMap

func init() {
	q := &QueryDocumentKeys
	bools := (*[unsafe.Sizeof(*q)]bool)(unsafe.Pointer(q))
	for i := range bools {
		bools[i] = true
	}
}

type _Visitable interface {
	visit(*QueryKeyMap) []_VisitNode
}

func (n Name) visit(keyMap *QueryKeyMap) []_VisitNode {
	return nil
}

func (n Document) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, len(n.Definitions))
	for i, node := range n.Definitions {
		result = append(result, vsi(node, "Definitions", i))
	}
	return result
}

func (n OperationDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 2+len(n.VariableDefinitions)+len(n.Directives))[:0]
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.VariableDefinitions {
		result = append(result, vsi(node, "VariableDefinitions", i))
	}
	for i, node := range n.Directives {
		result = append(result, vsi(node, "Directives", i))
	}
	result = append(result, vs(n.SelectionSet, "SelectionSet"))
	return result
}

func (n VariableDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{
		vs(n.Variable, "Variable"),
		vs(n.Type, "IType"),
		vs(n.DefaultValue, "DefaultValue"),
	}
}

func (n Variable) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{vs(n.Name, "Name")}
}

func (n SelectionSet) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, len(n.Selections))
	for i, node := range n.Selections {
		result[i] = vsi(node, "Selections", i)
	}
	return result
}

func (n Field) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 3+len(n.Arguments)+len(n.Directives))[:0]
	result = append(result, vs(n.Alias, "Alias"))
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Arguments {
		result = append(result, vsi(node, "Arguments", i))
	}
	for i, node := range n.Directives {
		result = append(result, vsi(node, "Directives", i))
	}
	result = append(result, vs(n.SelectionSet, "SelectionSet"))
	return result
}

func (n Argument) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{
		vs(n.Name, "Name"),
		vs(n.Value, "IValue"),
	}
}

func (n FragmentSpread) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 1+len(n.Directives))[:0]
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Directives {
		result = append(result, vsi(node, "Directives", i))
	}
	return result
}

func (n InlineFragment) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 2+len(n.Directives))[:0]
	result = append(result, vs(n.TypeCondition, "TypeCondition"))
	for i, node := range n.Directives {
		result = append(result, vsi(node, "Directives", i))
	}
	result = append(result, vs(n.SelectionSet, "SelectionSet"))
	return result
}

func (n FragmentDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 3+len(n.Directives))[:0]
	result = append(result, vs(n.Name, "Name"))
	result = append(result, vs(n.TypeCondition, "TypeCondition"))
	for i, node := range n.Directives {
		result = append(result, vsi(node, "Directives", i))
	}
	result = append(result, vs(n.SelectionSet, "SelectionSet"))
	return result
}

func (n IntValue) visit(keyMap *QueryKeyMap) []_VisitNode {
	return nil
}

func (n FloatValue) visit(keyMap *QueryKeyMap) []_VisitNode {
	return nil
}

func (n StringValue) visit(keyMap *QueryKeyMap) []_VisitNode {
	return nil
}

func (n BooleanValue) visit(keyMap *QueryKeyMap) []_VisitNode {
	return nil
}

func (n EnumValue) visit(keyMap *QueryKeyMap) []_VisitNode {
	return nil
}

func (n ListValue) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, len(n.Values))
	for i, node := range n.Values {
		result[i] = vsi(node, "Values", i)
	}
	return result
}

func (n ObjectValue) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, len(n.Fields))
	for i, node := range n.Fields {
		result[i] = vsi(node, "Fields", i)
	}
	return result
}

func (n ObjectField) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{
		vs(n.Name, "Name"),
		vs(n.Value, "IValue"),
	}
}

func (n Directive) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 1+len(n.Arguments))[:0]
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Arguments {
		result = append(result, vsi(node, "Arguments", i))
	}
	return result
}

func (n NamedType) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{vs(n.Name, "Name")}
}

func (n ListType) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{vs(n.Type, "IType")}
}

func (n NonNullType) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{vs(n.Type, "IType")}
}

func (n ObjectTypeDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 1+len(n.Interfaces)+len(n.Fields))[:0]
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Interfaces {
		result = append(result, vsi(node, "Interfaces", i))
	}
	for i, node := range n.Fields {
		result = append(result, vsi(node, "Fields", i))
	}
	return result
}

func (n FieldDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 2+len(n.Arguments))
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Arguments {
		result = append(result, vsi(node, "Arguments", i))
	}
	result = append(result, vs(n.Type, "IType"))
	return result
}

func (n InputValueDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{
		vs(n.Name, "Name"),
		vs(n.Type, "IType"),
		vs(n.DefaultValue, "DefaultValue"),
	}
}

func (n InterfaceTypeDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 1+len(n.Fields))[:0]
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Fields {
		result = append(result, vsi(node, "Fields", i))
	}
	return result
}

func (n UnionTypeDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 1+len(n.Types))[:0]
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Types {
		result = append(result, vsi(node, "Types", i))
	}
	return result
}

func (n ScalarTypeDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{vs(n.Name, "Name")}
}

func (n EnumTypeDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 1+len(n.Values))[:0]
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Values {
		result = append(result, vsi(node, "Values", i))
	}
	return result
}

func (n EnumValueDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{vs(n.Name, "Name")}
}

func (n InputObjectTypeDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	result := make([]_VisitNode, 1+len(n.Fields))[:0]
	result = append(result, vs(n.Name, "Name"))
	for i, node := range n.Fields {
		result = append(result, vsi(node, "Fields", i))
	}
	return result
}

func (n TypeExtensionDefinition) visit(keyMap *QueryKeyMap) []_VisitNode {
	return []_VisitNode{vs(n.Definition, "IDefinition")}
}

type ManualVisitor interface {
	Visit(INode, func(INode))
}

type ManualVisitorFunc func(INode, func(INode))

func (m ManualVisitorFunc) Visit(node INode, visit func(INode)) {
	m(node, visit)
}

type Visitor interface {
	Enter(INode, VisitInfo) VisitAction
	Leave(INode, VisitInfo) VisitAction
}

type VisitorFunc struct {
	EnterFunc func(INode, VisitInfo) VisitAction
	LeaveFunc func(INode, VisitInfo) VisitAction
}

func (v VisitorFunc) Enter(node INode, info VisitInfo) VisitAction {
	if v.Enter != nil {
		return v.Enter(node, info)
	}
	return nil
}

func (v VisitorFunc) Leave(node INode, info VisitInfo) VisitAction {
	if v.Leave != nil {
		return v.Leave(node, info)
	}
	return nil
}

type VisitInfo struct {
	Parent    INode
	Ancestors []INode
	Key       string
	Path      []string
}

type VisitAction interface {
	newNode() INode
}

type _VisitAction struct {
	node INode
}

func (v _VisitAction) newNode() INode {
	return v.node
}

var (
	VISIT_BREAK  = &_VisitAction{}
	VISIT_SKIP   = &_VisitAction{}
	VISIT_DELETE = &_VisitAction{}
)

func VISIT_REPLACE(node INode) VisitAction {
	return &_VisitAction{node}
}

type _VisitStack struct {
	index      int
	visitNodes []_VisitNode
	edits      []_VisitEdit
	prev       *_VisitStack
}

type _VisitEdit struct {
	key  string
	node INode
}

type _VisitNode struct {
	node  INode
	name  string
	index int
}

func vs(node INode, name string) _VisitNode {
	return _VisitNode{
		node:  node,
		name:  name,
		index: 0,
	}
}

func vsi(node INode, name string, index int) _VisitNode {
	return _VisitNode{
		node:  node,
		name:  name,
		index: index,
	}
}

func ManualVisit(root INode, visitor ManualVisitor) {

}

func Visit(root INode, visitor Visitor, keyMap *QueryKeyMap) {
	if keyMap == nil {
		keyMap = &QueryDocumentKeys
	}

	var parent INode
	var path []string
	var ancestors []INode
	var stack = &_VisitStack{
		visitNodes: []_VisitNode{{node: root}},
	}
	for ; ; stack.index++ {
		var visitNode _VisitNode
		var node INode

		var isEntering = stack.index < len(stack.visitNodes)
		if isEntering {
			visitNode = stack.visitNodes[stack.index]
			node = visitNode.node

		} else {
			path = path[:len(path)-1]
			node = ancestors[len(ancestors)-1]
			ancestors = ancestors[:len(ancestors)-1]
			stack = stack.prev
		}

		info := VisitInfo{
			Parent:    parent,
			Ancestors: ancestors,
			Key:       visitNode.name,
			Path:      path,
		}

		var result VisitAction
		if isEntering {
			if visitor.Enter != nil {
				result = visitor.Enter(node, info)
			}
		} else {
			if visitor.Leave != nil {
				result = visitor.Leave(node, info)
			}
		}

		if result == VISIT_BREAK {
			break
		}
		if result == VISIT_SKIP {
			// skip the node
		}
		if result == VISIT_DELETE {
			// delete the node and stop digging deeper
		}
		if result != nil && result.newNode() != nil {
			// replace the node with a new one
		}

		// if not skip or delete
		if isEntering {
			nextNodes := visitNode.node.visit(keyMap)
			if len(nextNodes) > 0 {
				stack = &_VisitStack{
					visitNodes: nextNodes,
					prev:       stack,
				}
			}

			parent = node
			path = append(path, visitNode.name)
			ancestors = append(ancestors, node)
		}

		if !isEntering && len(ancestors) == 0 {
			break
		}
	}
}
