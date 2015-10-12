package language

type Location struct {
	Start  int
	End    int
	Source *Source
}

func (l *Location) Loc() *Location {
	return l
}

type INode interface {
	_Visitable
	Kind() NodeKind
	Loc() *Location
	clone() INode
}

func (n *Name) Kind() NodeKind                      { return NAME }
func (n *Document) Kind() NodeKind                  { return DOCUMENT }
func (n *OperationDefinition) Kind() NodeKind       { return OPERATION_DEFINITION }
func (n *VariableDefinition) Kind() NodeKind        { return VARIABLE_DEFINITION }
func (n *Variable) Kind() NodeKind                  { return VARIABLE }
func (n *SelectionSet) Kind() NodeKind              { return SELECTION_SET }
func (n *Field) Kind() NodeKind                     { return FIELD }
func (n *Argument) Kind() NodeKind                  { return ARGUMENT }
func (n *FragmentSpread) Kind() NodeKind            { return FRAGMENT_SPREAD }
func (n *InlineFragment) Kind() NodeKind            { return INLINE_FRAGMENT }
func (n *FragmentDefinition) Kind() NodeKind        { return FRAGMENT_DEFINITION }
func (n *IntValue) Kind() NodeKind                  { return INT_VALUE }
func (n *FloatValue) Kind() NodeKind                { return FLOAT_VALUE }
func (n *StringValue) Kind() NodeKind               { return STRING_VALUE }
func (n *BooleanValue) Kind() NodeKind              { return BOOLEAN_VALUE }
func (n *EnumValue) Kind() NodeKind                 { return ENUM_VALUE }
func (n *ListValue) Kind() NodeKind                 { return LIST_VALUE }
func (n *ObjectValue) Kind() NodeKind               { return OBJECT_VALUE }
func (n *ObjectField) Kind() NodeKind               { return OBJECT_FIELD }
func (n *Directive) Kind() NodeKind                 { return DIRECTIVE }
func (n *NamedType) Kind() NodeKind                 { return NAMED_TYPE }
func (n *ListType) Kind() NodeKind                  { return LIST_TYPE }
func (n *NonNullType) Kind() NodeKind               { return NON_NULL_TYPE }
func (n *ObjectTypeDefinition) Kind() NodeKind      { return OBJECT_TYPE_DEFINITION }
func (n *FieldDefinition) Kind() NodeKind           { return FIELD_DEFINITION }
func (n *InputValueDefinition) Kind() NodeKind      { return INPUT_VALUE_DEFINITION }
func (n *InterfaceTypeDefinition) Kind() NodeKind   { return INTERFACE_TYPE_DEFINITION }
func (n *UnionTypeDefinition) Kind() NodeKind       { return UNION_TYPE_DEFINITION }
func (n *ScalarTypeDefinition) Kind() NodeKind      { return SCALAR_TYPE_DEFINITION }
func (n *EnumTypeDefinition) Kind() NodeKind        { return ENUM_TYPE_DEFINITION }
func (n *EnumValueDefinition) Kind() NodeKind       { return ENUM_VALUE_DEFINITION }
func (n *InputObjectTypeDefinition) Kind() NodeKind { return INPUT_OBJECT_TYPE_DEFINITION }
func (n *TypeExtensionDefinition) Kind() NodeKind   { return TYPE_EXTENSION_DEFINITION }

func (n *Name) clone() INode                    { var v = &Name{}; *v = *n; return v }
func (n *Document) clone() INode                { var v = &Document{}; *v = *n; return v }
func (n *OperationDefinition) clone() INode     { var v = &OperationDefinition{}; *v = *n; return v }
func (n *VariableDefinition) clone() INode      { var v = &VariableDefinition{}; *v = *n; return v }
func (n *Variable) clone() INode                { var v = &Variable{}; *v = *n; return v }
func (n *SelectionSet) clone() INode            { var v = &SelectionSet{}; *v = *n; return v }
func (n *Field) clone() INode                   { var v = &Field{}; *v = *n; return v }
func (n *Argument) clone() INode                { var v = &Argument{}; *v = *n; return v }
func (n *FragmentSpread) clone() INode          { var v = &FragmentSpread{}; *v = *n; return v }
func (n *InlineFragment) clone() INode          { var v = &InlineFragment{}; *v = *n; return v }
func (n *FragmentDefinition) clone() INode      { var v = &FragmentDefinition{}; *v = *n; return v }
func (n *IntValue) clone() INode                { var v = &IntValue{}; *v = *n; return v }
func (n *FloatValue) clone() INode              { var v = &FloatValue{}; *v = *n; return v }
func (n *StringValue) clone() INode             { var v = &StringValue{}; *v = *n; return v }
func (n *BooleanValue) clone() INode            { var v = &BooleanValue{}; *v = *n; return v }
func (n *EnumValue) clone() INode               { var v = &EnumValue{}; *v = *n; return v }
func (n *ListValue) clone() INode               { var v = &ListValue{}; *v = *n; return v }
func (n *ObjectValue) clone() INode             { var v = &ObjectValue{}; *v = *n; return v }
func (n *ObjectField) clone() INode             { var v = &ObjectField{}; *v = *n; return v }
func (n *Directive) clone() INode               { var v = &Directive{}; *v = *n; return v }
func (n *NamedType) clone() INode               { var v = &NamedType{}; *v = *n; return v }
func (n *ListType) clone() INode                { var v = &ListType{}; *v = *n; return v }
func (n *NonNullType) clone() INode             { var v = &NonNullType{}; *v = *n; return v }
func (n *ObjectTypeDefinition) clone() INode    { var v = &ObjectTypeDefinition{}; *v = *n; return v }
func (n *FieldDefinition) clone() INode         { var v = &FieldDefinition{}; *v = *n; return v }
func (n *InputValueDefinition) clone() INode    { var v = &InputValueDefinition{}; *v = *n; return v }
func (n *InterfaceTypeDefinition) clone() INode { var v = &InterfaceTypeDefinition{}; *v = *n; return v }
func (n *UnionTypeDefinition) clone() INode     { var v = &UnionTypeDefinition{}; *v = *n; return v }
func (n *ScalarTypeDefinition) clone() INode    { var v = &ScalarTypeDefinition{}; *v = *n; return v }
func (n *EnumTypeDefinition) clone() INode      { var v = &EnumTypeDefinition{}; *v = *n; return v }
func (n *EnumValueDefinition) clone() INode     { var v = &EnumValueDefinition{}; *v = *n; return v }
func (n *InputObjectTypeDefinition) clone() INode {
	var v = &InputObjectTypeDefinition{}
	*v = *n
	return v
}
func (n *TypeExtensionDefinition) clone() INode { var v = &TypeExtensionDefinition{}; *v = *n; return v }

// Name

type Name struct {
	*Location
	Value string
}

// Document

type Document struct {
	*Location
	Definitions []IDefinition
}

type IDefinition interface {
	INode
	definitionNode()
}

func (n *OperationDefinition) definitionNode() {}
func (n *FragmentDefinition) definitionNode()  {}

type OperationType string

const (
	OperationQuery        = OperationType("query")
	OperationMutation     = OperationType("mutation")
	OperationSubscription = OperationType("subscription")
)

type OperationDefinition struct {
	*Location
	Operation           OperationType
	Name                *Name
	VariableDefinitions []*VariableDefinition
	Directives          []*Directive
	SelectionSet        *SelectionSet
}

type VariableDefinition struct {
	*Location
	Variable     *Variable
	Type         IType
	DefaultValue IValue
}

type Variable struct {
	*Location
	Name *Name
}

type SelectionSet struct {
	*Location
	Selections []ISelection
}

type ISelection interface {
	INode
	selectionNode()
}

func (n *Field) selectionNode()          {}
func (n *FragmentSpread) selectionNode() {}
func (n *InlineFragment) selectionNode() {}

type IFragment interface {
	ISelection
	fragmentNode()
}

func (n *FragmentSpread) fragmentNode() {}
func (n *InlineFragment) fragmentNode() {}

type Field struct {
	*Location
	Alias        *Name
	Name         *Name
	Arguments    []*Argument
	Directives   []*Directive
	SelectionSet *SelectionSet
}

type Argument struct {
	*Location
	Name  *Name
	Value IValue
}

// Fragments

type FragmentSpread struct {
	*Location
	Name       *Name
	Directives []*Directive
}

type ITypeCondition interface {
	GetTypeCondition() *NamedType
}

func (n *InlineFragment) GetTypeCondition() *NamedType {
	return n.TypeCondition
}

func (n *FragmentDefinition) GetTypeCondition() *NamedType {
	return n.TypeCondition
}

type InlineFragment struct {
	*Location
	TypeCondition *NamedType
	Directives    []*Directive
	SelectionSet  *SelectionSet
}

type FragmentDefinition struct {
	*Location
	Name          *Name
	TypeCondition *NamedType
	Directives    []*Directive
	SelectionSet  *SelectionSet
}

// Values

type IValue interface {
	INode
	valueNode()
}

func (n *Variable) valueNode()     {}
func (n *IntValue) valueNode()     {}
func (n *FloatValue) valueNode()   {}
func (n *StringValue) valueNode()  {}
func (n *BooleanValue) valueNode() {}
func (n *EnumValue) valueNode()    {}
func (n *ListValue) valueNode()    {}
func (n *ObjectValue) valueNode()  {}

type IScalarValue interface {
	IValue
	GetValue() string
}

func (n *IntValue) GetValue() string     { return n.Value }
func (n *FloatValue) GetValue() string   { return n.Value }
func (n *StringValue) GetValue() string  { return n.Value }
func (n *BooleanValue) GetValue() string { return n.Value }

type IntValue struct {
	*Location
	Value string
}

type FloatValue struct {
	*Location
	Value string
}

type StringValue struct {
	*Location
	Value string
}

type BooleanValue struct {
	*Location
	Value string
}

type EnumValue struct {
	*Location
	Value string
}

type ListValue struct {
	*Location
	Values []IValue
}

type ObjectValue struct {
	*Location
	Fields []*ObjectField
}

type ObjectField struct {
	*Location
	Name  *Name
	Value IValue
}

// Directives

type Directive struct {
	*Location
	Name      *Name
	Arguments []*Argument
}

// IType Reference

type IType interface {
	INode
	typeNode()
}

func (n *NamedType) typeNode()   {}
func (n *ListType) typeNode()    {}
func (n *NonNullType) typeNode() {}

type INonNullType interface {
	IType
	nonNullTypeNode()
}

func (n *NamedType) nonNullTypeNode() {}
func (n *ListType) nonNullTypeNode()  {}

type NamedType struct {
	*Location
	Name *Name
}

type ListType struct {
	*Location
	Type IType
}

type NonNullType struct {
	*Location
	Type INonNullType
}

// IType IDefinition

type ITypeDefinition interface {
	IDefinition
	typeDefinitionNode()
}

func (n *ObjectTypeDefinition) definitionNode()      {}
func (n *InterfaceTypeDefinition) definitionNode()   {}
func (n *UnionTypeDefinition) definitionNode()       {}
func (n *ScalarTypeDefinition) definitionNode()      {}
func (n *EnumTypeDefinition) definitionNode()        {}
func (n *InputObjectTypeDefinition) definitionNode() {}
func (n *TypeExtensionDefinition) definitionNode()   {}

func (n *ObjectTypeDefinition) typeDefinitionNode()      {}
func (n *InterfaceTypeDefinition) typeDefinitionNode()   {}
func (n *UnionTypeDefinition) typeDefinitionNode()       {}
func (n *ScalarTypeDefinition) typeDefinitionNode()      {}
func (n *EnumTypeDefinition) typeDefinitionNode()        {}
func (n *InputObjectTypeDefinition) typeDefinitionNode() {}
func (n *TypeExtensionDefinition) typeDefinitionNode()   {}

type ObjectTypeDefinition struct {
	*Location
	Name       *Name
	Interfaces []*NamedType
	Fields     []*FieldDefinition
}

type FieldDefinition struct {
	*Location
	Name      *Name
	Arguments []*InputValueDefinition
	Type      IType
}

type InputValueDefinition struct {
	*Location
	Name         *Name
	Type         IType
	DefaultValue IValue
}

type InterfaceTypeDefinition struct {
	*Location
	Name   *Name
	Fields []*FieldDefinition
}

type UnionTypeDefinition struct {
	*Location
	Name  *Name
	Types []*NamedType
}

type ScalarTypeDefinition struct {
	*Location
	Name *Name
}

type EnumTypeDefinition struct {
	*Location
	Name   *Name
	Values []*EnumValueDefinition
}

type EnumValueDefinition struct {
	*Location
	Name *Name
}

type InputObjectTypeDefinition struct {
	*Location
	Name   *Name
	Fields []*InputValueDefinition
}

type TypeExtensionDefinition struct {
	*Location
	Definition *ObjectTypeDefinition
}
