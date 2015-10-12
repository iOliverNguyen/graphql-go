package ql

const (
	INT     = "IntValue"
	FLOAT   = "FloatValue"
	STRING  = "StringValue"
	BOOLEAN = "BooleanValue"
)

type Type interface {
	graphqlType()
}

func (t Scalar) graphqlType()      {}
func (t Object) graphqlType()      {}
func (t Interface) graphqlType()   {}
func (t Union) graphqlType()       {}
func (t Enum) graphqlType()        {}
func (t InputObject) graphqlType() {}
func (t List) graphqlType()        {}
func (t NonNull) graphqlType()     {}

type InputType interface {
	inputType()
}

func (t Scalar) inputType()      {}
func (t Enum) inputType()        {}
func (t InputObject) inputType() {}
func (t List) inputType()        {}
func (t NonNull) inputType()     {}

type OutputType interface {
	outputType()
}

func (t Scalar) outputType()    {}
func (t Object) outputType()    {}
func (t Interface) outputType() {}
func (t Union) outputType()     {}
func (t Enum) outputType()      {}
func (t List) outputType()      {}
func (t NonNull) outputType()   {}

type Objects []Object
type Object struct {
	Name           string
	Interfaces     Interfaces
	InterfacesFunc func() Interfaces
	Fields         FieldMap
	FieldsFunc     func() FieldMap
	// IsTypeOf     func(v interface{}, info *GraphQLResolveInfo) bool
	IsTypeOf    func(v interface{}, info interface{}) bool
	Description string
}

type FieldMap map[string]Field
type Field struct {
	Type              OutputType
	Args              ArgumentMap
	Resolve           interface{}
	DeprecationReason string
	Description       string
}

type ArgumentMap map[string]Argument
type Argument struct {
	Type         InputType
	DefaultValue interface{}
	Description  string
}

type Scalar struct {
	Name         string
	Description  string
	Serialize    func(v interface{}) interface{}
	ParseValue   func(v interface{}) interface{}
	ParseLiteral func(kind, value string) interface{}
}

type Interfaces []Interface
type Interface struct {
	Name       string
	Fields     FieldMap
	FieldsFunc func() FieldMap
	// ResolveType func(v interface{}, info *GraphQLResolveInfo) *GraphQLObjectType
	ResolveType func(v interface{}, info interface{}) interface{}
	Description string
}

type Union struct {
	Name  string
	Types Objects
	// ResolveType func(v interface{}, info *GraphQLResolveInfo) *GraphQLObjectType
	ResolveType func(v interface{}, info interface{}) interface{}
	Description string
}

type Enum struct {
	Name        string
	Values      EnumValueMap
	Description string
}

type EnumValueMap map[string]EnumValue
type EnumValue struct {
	Value             interface{}
	DeprecationReason string
	Description       string
}

type InputObject struct {
	Name        string
	Fields      InputObjectFieldMap
	FieldsFunc  func() InputObjectFieldMap
	Description string
}

type InputObjectFieldMap map[string]InputObjectField
type InputObjectField struct {
	Type         InputType
	DefaultValue interface{}
	Description  string
}

type List struct {
	OfType Type
}

type NonNull struct {
	OfType Type
}
