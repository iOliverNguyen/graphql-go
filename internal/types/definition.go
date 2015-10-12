package types

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	lang "github.com/ng-vu/graphql-go/internal/language"
	"github.com/ng-vu/graphql-go/ql"
)

type QLType interface {
	GetName() string
	graphqlType()
}

func (t *QLScalar) graphqlType()      {}
func (t *QLObject) graphqlType()      {}
func (t *QLInterface) graphqlType()   {}
func (t *QLUnion) graphqlType()       {}
func (t *QLEnum) graphqlType()        {}
func (t *QLInputObject) graphqlType() {}
func (t *QLList) graphqlType()        {}
func (t *QLNonNull) graphqlType()     {}

func (t *QLScalar) GetName() string      { return t.Name }
func (t *QLObject) GetName() string      { return t.Name }
func (t *QLInterface) GetName() string   { return t.Name }
func (t *QLUnion) GetName() string       { return t.Name }
func (t *QLEnum) GetName() string        { return t.Name }
func (t *QLInputObject) GetName() string { return t.Name }
func (t *QLList) GetName() string        { return "" }
func (t *QLNonNull) GetName() string     { return "" }

// NOTE(ng-vu): possible missing in graphql-js
func (t *QLArgument) graphqlType()    {}
func (t *QLArgument) GetName() string { return t.Name }

type QLInputType interface {
	QLType
	graphqlInputType()
}

func (t *QLScalar) graphqlInputType()      {}
func (t *QLEnum) graphqlInputType()        {}
func (t *QLInputObject) graphqlInputType() {}
func (t *QLList) graphqlInputType()        {}
func (t *QLNonNull) graphqlInputType()     {}

type QLOutputType interface {
	QLType
	graphqlOutputType()
}

func (t *QLScalar) graphqlOutputType()    {}
func (t *QLObject) graphqlOutputType()    {}
func (t *QLInterface) graphqlOutputType() {}
func (t *QLUnion) graphqlOutputType()     {}
func (t *QLEnum) graphqlOutputType()      {}
func (t *QLList) graphqlOutputType()      {}
func (t *QLNonNull) graphqlOutputType()   {}

type QLLeafType interface {
	QLType
	graphqlLeafType()
}

func (t *QLScalar) graphqlLeafType() {}
func (t *QLEnum) graphqlLeafType()   {}

type QLCompositeType interface {
	QLType
	graphqlCompositeType()
}

func (t *QLObject) graphqlCompositeType()    {}
func (t *QLInterface) graphqlCompositeType() {}
func (t *QLUnion) graphqlCompositeType()     {}

type QLAbstractType interface {
	QLType
	graphqlAbstractType()
	GetPossibleTypes() []*QLObject
	IsPossibleType(*QLObject) bool
}

func (t *QLInterface) graphqlAbstractType() {}
func (t *QLUnion) graphqlAbstractType()     {}

type QLNullableType interface {
	QLType
	graphqlNullableType()
}

func (t *QLScalar) graphqlNullableType()      {}
func (t *QLObject) graphqlNullableType()      {}
func (t *QLInterface) graphqlNullableType()   {}
func (t *QLUnion) graphqlNullableType()       {}
func (t *QLEnum) graphqlNullableType()        {}
func (t *QLInputObject) graphqlNullableType() {}
func (t *QLList) graphqlNullableType()        {}

type QLNamedType interface {
	QLType
	graphqlNamedType()
}

func (t *QLScalar) graphqlNamedType()      {}
func (t *QLObject) graphqlNamedType()      {}
func (t *QLInterface) graphqlNamedType()   {}
func (t *QLUnion) graphqlNamedType()       {}
func (t *QLEnum) graphqlNamedType()        {}
func (t *QLInputObject) graphqlNamedType() {}

type QLObjectInterface interface {
	GetFields() map[string]QLFieldDefinition
}

func NewQLType(config ql.Type) QLType {
	switch config := config.(type) {
	case ql.Scalar:
		return NewQLScalar(config)
	case ql.Object:
		return NewQLObject(config)
	case ql.Interface:
		return NewQLInterface(config)
	case ql.Union:
		return NewQLUnion(config)
	case ql.Enum:
		return NewQLEnum(config)
	case ql.InputObject:
		return NewQLInputObject(config)
	case ql.List:
		return NewQLList(NewQLType(config.OfType))
	case ql.NonNull:
		return NewQLNonNull(NewQLType(config.OfType))
	default:
		panic("graphql-go/types: unreachable")
	}
}

func NewQLInputType(config ql.InputType) QLInputType {
	switch config := config.(type) {
	case ql.Scalar:
		return NewQLScalar(config)
	case ql.Enum:
		return NewQLEnum(config)
	case ql.InputObject:
		return NewQLInputObject(config)
	case ql.List:
		return NewQLList(NewQLType(config.OfType))
	case ql.NonNull:
		return NewQLNonNull(NewQLType(config.OfType))
	default:
		panic("graphql-go/types: unreachable")
	}
}

func NewQLOutputType(config ql.OutputType) QLOutputType {
	switch config := config.(type) {
	case ql.Scalar:
		return NewQLScalar(config)
	case ql.Object:
		return NewQLObject(config)
	case ql.Interface:
		return NewQLInterface(config)
	case ql.Union:
		return NewQLUnion(config)
	case ql.Enum:
		return NewQLEnum(config)
	case ql.List:
		return NewQLList(NewQLType(config.OfType))
	case ql.NonNull:
		return NewQLNonNull(NewQLType(config.OfType))
	default:
		panic("graphql-go/types: unreachable")
	}
}

func NewQLResolveFunc(fn interface{}) QLFieldResolveFunc {
	v := reflect.ValueOf(fn)
	t := v.Type()
	if t.Kind() != reflect.Func {
		panic("graphql: expect resolve function")
	}
	if t.NumIn() != 1 {
		panic("graphql: resolve function must receive one struct argument")
	}
	if t.NumOut() != 1 {
		panic("graphql: resolve function must return one result")
	}

	in0 := t.In(0)
	if in0.Kind() != reflect.Struct {
		panic("graphql: resolve function must receive struct as argument")
	}

	return func(source interface{}, args map[string]interface{}, info QLResolveInfo) interface{} {
		arg := reflect.New(in0).Elem()
		result := v.Call([]reflect.Value{arg})
		return result[0].Interface()
	}
}

/*
Scalar Type Definition

The leaf values of any request and input values to arguments are
Scalars (or Enums) and are defined with a name and a series of functions
used to parse input from ast or variables and to ensure validity.

Example:

    var OddType = new QLScalar({
      name: 'Odd',
      serialize(value) {
        return value % 2 === 1 ? value : null;
      }
    });

*/
type QLScalar struct {
	Name        string
	Description string

	config ql.Scalar
}

func NewQLScalar(config ql.Scalar) *QLScalar {
	if config.Name == "" {
		throw("Type must be named.")
	}
	assertValidName(config.Name)
	if config.Serialize == nil {
		throw(`%v must provide "Serialize" function. If this custom Scalar is `+
			`also used as an input type, ensure "ParseValue" and "ParseLiteral" `+
			`functions are also provided.`,
			config.Name)
	}
	if (config.ParseValue != nil) != (config.ParseLiteral != nil) {
		throw(`%v must provide both "ParseValue" and "ParseLiteral" functions.`,
			config.Name)
	}
	return &QLScalar{
		Name:        config.Name,
		Description: config.Description,
		config:      config,
	}
}

func (g *QLScalar) String() string {
	return g.Name
}

func (g *QLScalar) Serialize(v interface{}) interface{} {
	return g.config.Serialize(v)
}

func (g *QLScalar) ParseValue(v interface{}) interface{} {
	if parser := g.config.ParseValue; parser != nil {
		return parser(v)
	}
	return nil
}

func (g *QLScalar) ParseLiteral(valueAST lang.IValue) interface{} {
	if parser := g.config.ParseLiteral; parser != nil {
		if valueAST, ok := valueAST.(lang.IScalarValue); ok {
			return parser(string(valueAST.Kind()), valueAST.GetValue())
		}
	}
	return nil
}

/*
Object Type Definition

Almost all of the QL types you define will be object types. Object types
have a name, but most importantly describe their fields.

Example:

   var AddressType = new QLObject({
     name: 'Address',
     fields: {
       street: { type: QLString },
       number: { type: QLInt },
       formatted: {
         type: QLString,
         resolve(obj) {
           return obj.number + ' ' + obj.street
         }
       }
     }
   });

When two types need to refer to each other, or a type needs to refer to
itself in a field, you can use a function expression (aka a closure or a
thunk) to supply the fields lazily.

Example:

   var PersonType = new QLObject({
     name: 'Person',
     fields: () => ({
       name: { type: QLString },
       bestFriend: { type: PersonType },
     })
   });

*/
type QLObject struct {
	Name        string
	Description string
	IsTypeOf    func(v interface{}, info *QLResolveInfo) bool

	config     ql.Object
	fields     map[string]*QLFieldDefinition
	interfaces []*QLInterface
}

func NewQLObject(config ql.Object) *QLObject {
	if config.Name == "" {
		throw("Type must be named.")
	}
	assertValidName(config.Name)

	if config.InterfacesFunc != nil && config.Interfaces != nil {
		throw(`%v must provide either "Interfaces" or "InterfacesFunc", not both.`, config.Name)
	}

	if config.FieldsFunc != nil && config.Fields != nil {
		throw(`%v must provide either "Fields" or "FieldsFunc", not both.`, config.Name)
	}

	g := &QLObject{
		Name:        config.Name,
		Description: config.Description,
		// TODO(qv): IsTypeOf
		// IsTypeOf:    config.IsTypeOf,
		config: config,
	}
	addImplementationToInterfaces(g)
	return g
}

func (g *QLObject) String() string {
	return g.Name
}

func (g *QLObject) GetFields() map[string]*QLFieldDefinition {
	if g.fields == nil {
		g.fields = defineFieldMap(g, g.config.Fields, g.config.FieldsFunc)
	}
	return g.fields
}

func (g *QLObject) GetInterfaces() []*QLInterface {
	if g.interfaces == nil {
		g.interfaces = g.defineInterfaces()
	}
	return g.interfaces
}

func (g *QLObject) defineInterfaces() []*QLInterface {
	interfaces := g.config.Interfaces
	if fn := g.config.InterfacesFunc; fn != nil {
		interfaces = fn()
	}
	result := make([]*QLInterface, len(interfaces))
	for i, iface := range interfaces {
		result[i] = NewQLInterface(iface)
	}
	return result
}

func defineFieldMap(
	typ QLNamedType,
	fieldMap ql.FieldMap,
	fieldMapFunc func() ql.FieldMap,
) map[string]*QLFieldDefinition {

	if fieldMapFunc != nil {
		fieldMap = fieldMapFunc()
	}
	result := make(map[string]*QLFieldDefinition)
	for fieldName, fieldConfig := range fieldMap {
		args := make([]*QLArgument, len(fieldConfig.Args))[:0]
		for argName, argConfig := range fieldConfig.Args {
			arg := &QLArgument{
				Name:        argName,
				Description: argConfig.Description,
				// TODO(qv): convert ql.Input to QLInput
				// Type:         argConfig.Type,
				DefaultValue: argConfig.DefaultValue,
			}
			args = append(args, arg)
		}

		field := &QLFieldDefinition{
			Name:              fieldName,
			Description:       fieldConfig.Description,
			Type:              NewQLOutputType(fieldConfig.Type),
			Args:              args,
			Resolve:           NewQLResolveFunc(fieldConfig.Resolve),
			DeprecationReason: fieldConfig.DeprecationReason,
		}
		result[fieldName] = field
	}
	return result
}

func addImplementationToInterfaces(impl *QLObject) {
	// TODO(qv)
}

type QLFieldResolveFunc func(
	source interface{},
	args map[string]interface{},
	info QLResolveInfo,
) interface{}

type QLResolveInfo struct {
	FieldName      string
	FieldASTs      []*lang.Field
	ReturnType     QLOutputType
	ParentType     QLCompositeType
	Schema         QLSchema
	Fragments      map[string]*lang.FragmentDefinition
	RootValue      interface{}
	Operation      *lang.OperationDefinition
	VariableValues map[string]interface{}
}

type QLFieldDefinition struct {
	Name              string
	Description       string
	Type              QLOutputType
	Args              []*QLArgument
	Resolve           QLFieldResolveFunc
	DeprecationReason string
}

type QLArgument struct {
	Name         string
	Type         QLInputType
	DefaultValue interface{}
	Description  string
}

// type QLFieldDefinitionMap map[string]QLFieldDefinition

/*
Interface Type Definition

When a field can return one of a heterogeneous set of types, a Interface type
is used to describe what types are possible, what fields are in common across
all types, as well as a function to determine which type is actually used
when the field is resolved.

Example:

    var EntityType = new QLInterface({
      name: 'Entity',
      fields: {
        name: { type: QLString }
      }
    });

*/
type QLInterface struct {
	Name        string
	Description string
	ResolveType func(v interface{}, info *QLResolveInfo) *QLObject

	config          ql.Interface
	fields          map[string]*QLFieldDefinition
	implementations []*QLObject
	positionTypes   map[string]*QLObject
}

func NewQLInterface(config ql.Interface) *QLInterface {
	if config.Name == "" {
		throw("Type must be named.")
	}
	assertValidName(config.Name)
	return &QLInterface{
		Name:        config.Name,
		Description: config.Description,
		// TODO(qv): resolve type
		// ResolveType:     config.ResolveType,
		config:          config,
		implementations: make([]*QLObject, 4)[:0],
	}
}

func (g *QLInterface) String() string {
	return g.Name
}

func (g *QLInterface) GetFields() map[string]*QLFieldDefinition {
	if g.fields == nil {
		g.fields = defineFieldMap(g, g.config.Fields, g.config.FieldsFunc)
	}
	return g.fields
}

func (g *QLInterface) GetPossibleTypes() []*QLObject {
	return g.implementations
}

func (g *QLInterface) IsPossibleType(typ *QLObject) bool {
	if g.positionTypes == nil {
		possibleTypes := make(map[string]*QLObject)
		for _, objType := range g.GetPossibleTypes() {
			possibleTypes[objType.Name] = objType
		}
		g.positionTypes = possibleTypes
	}
	_, ok := g.positionTypes[typ.Name]
	return ok
}

func (g *QLInterface) GetObjectType(v interface{}, info *QLResolveInfo) *QLObject {
	resolver := g.ResolveType
	if resolver != nil {
		return resolver(v, info)
	}
	return getTypeOf(v, info, g)
}

func getTypeOf(v interface{}, info *QLResolveInfo, abstractType QLAbstractType) *QLObject {
	possibleTypes := abstractType.GetPossibleTypes()
	for _, typ := range possibleTypes {
		if typ.IsTypeOf != nil && typ.IsTypeOf(v, info) {
			return typ
		}
	}
	return nil
}

/*
Union Type Definition

When a field can return one of a heterogeneous set of types, a Union type
is used to describe what types are possible as well as providing a function
to determine which type is actually used when the field is resolved.

Example:

    var PetType = new QLUnion({
      name: 'Pet',
      types: [ DogType, CatType ],
      resolveType(value) {
        if (value instanceof Dog) {
          return DogType;
        }
        if (value instanceof Cat) {
          return CatType;
        }
      }
    });

*/
type QLUnion struct {
	Name        string
	Description string
	ResolveType func(v interface{}, info *QLResolveInfo) *QLObject

	config            ql.Union
	types             []*QLObject
	possibleTypeNames map[string]struct{}
}

func NewQLUnion(config ql.Union) *QLUnion {
	if config.Name == "" {
		throw("Type must be named.")
	}
	assertValidName(config.Name)
	if len(config.Types) == 0 {
		throw("Must provide Array of types for Union %v", config.Name)
	}
	types := make([]*QLObject, len(config.Types))
	if config.ResolveType == nil {
		for i, typ := range config.Types {
			if typ.IsTypeOf == nil {
				throw(`Union Type %v does not provide a "ResolveType" function and possible Type %v does not provide a "IsTypeOf" function. There is no way to resolve this possible type during execution.`, config.Name, typ.Name)
			}
			types[i] = NewQLObject(config.Types[i])
		}
	}
	return &QLUnion{
		Name:        config.Name,
		Description: config.Description,
		// TODO(qv): resolveType
		// ResolveType: config.ResolveType,
		config: config,
		types:  types,
	}
}

func (g *QLUnion) GetPossibleTypes() []*QLObject {
	return g.types
}

func (g *QLUnion) IsPossibleType(typ *QLObject) bool {
	possibleTypeNames := g.possibleTypeNames
	if possibleTypeNames == nil {
		possibleTypeNames = make(map[string]struct{})
		for _, possibleType := range g.GetPossibleTypes() {
			possibleTypeNames[possibleType.Name] = struct{}{}
		}
		g.possibleTypeNames = possibleTypeNames
	}
	_, ok := possibleTypeNames[typ.Name]
	return ok
}

func (g *QLUnion) GetObjectType(v interface{}, info *QLResolveInfo) *QLObject {
	if resolver := g.ResolveType; resolver != nil {
		return resolver(v, info)
	}
	return getTypeOf(v, info, g)
}

func (g *QLUnion) String() string {
	return g.Name
}

/*
Enum Type Definition

Some leaf values of requests and input values are Enums. QL serializes
Enum values as strings, however internally Enums can be represented by any
kind of type, often integers.

Example:

    var RGBType = new QLEnum({
      name: 'RGB',
      values: {
        RED: { value: 0 },
        GREEN: { value: 1 },
        BLUE: { value: 2 }
      }
    });

Note: If a value is not provided in a definition, the name of the enum value
will be used as it's internal value.
*/
type QLEnum struct {
	Name        string
	Description string

	config      ql.Enum
	values      []*QLEnumValueDefinition
	valueLookup map[interface{}]*QLEnumValueDefinition
	nameLookup  map[string]*QLEnumValueDefinition
}

type QLEnumValueDefinition struct {
	Name              string
	Value             interface{}
	DeprecationReason string
	Description       string
}

func NewQLEnum(config ql.Enum) *QLEnum {
	assertValidName(config.Name)
	g := &QLEnum{
		Name:        config.Name,
		Description: config.Description,
		config:      config,
	}
	g.values = g.defineEnumValues(config.Values)
	return g
}

func (g *QLEnum) String() string {
	return g.Name
}

func (g *QLEnum) GetValues() []*QLEnumValueDefinition {
	return g.values
}

func (g *QLEnum) Serialize(v interface{}) string {
	enumValue, ok := g.getValueLookup()[v]
	if ok {
		return enumValue.Name
	}
	return ""
}

func (g *QLEnum) ParseValue(v interface{}) interface{} {
	if v, ok := v.(string); ok {
		if enumValue, ok := g.getNameLookup()[v]; ok {
			return enumValue.Value
		}
	}
	return nil
}

func (g *QLEnum) ParseLiteral(valueAST lang.IValue) interface{} {
	if value, ok := valueAST.(*lang.EnumValue); ok {
		enumValue, ok := g.getNameLookup()[value.Value]
		if ok {
			return enumValue.Value
		}
	}
	return nil
}

func (g *QLEnum) getValueLookup() map[interface{}]*QLEnumValueDefinition {
	valueLookup := g.valueLookup
	if valueLookup == nil {
		valueLookup = make(map[interface{}]*QLEnumValueDefinition)
		for _, value := range g.GetValues() {
			valueLookup[value.Value] = value
		}
		g.valueLookup = valueLookup
	}
	return valueLookup
}

func (g *QLEnum) getNameLookup() map[string]*QLEnumValueDefinition {
	nameLookup := g.nameLookup
	if nameLookup == nil {
		nameLookup = make(map[string]*QLEnumValueDefinition)
		for _, value := range g.GetValues() {
			nameLookup[value.Name] = value
		}
		g.nameLookup = nameLookup
	}
	return nameLookup
}

func (g *QLEnum) defineEnumValues(
	valueMap ql.EnumValueMap,
) []*QLEnumValueDefinition {
	values := make([]*QLEnumValueDefinition, len(valueMap))[:0]
	for name, v := range valueMap {
		assertValidName(name)
		value := &QLEnumValueDefinition{
			Name:              name,
			Value:             v.Value,
			Description:       v.Description,
			DeprecationReason: v.DeprecationReason,
		}
		if value.Value == nil {
			value.Value = name
		}
		values = append(values, value)
	}
	return values
}

/*
Input Object Type Definition

An input object defines a structured collection of fields which may be
supplied to a field argument.

Using `NonNull` will ensure that a value must be provided by the query

Example:

    var GeoPoint = new QLInputObject({
      name: 'GeoPoint',
      fields: {
        lat: { type: new QLNonNull(QLFloat) },
        lon: { type: new QLNonNull(QLFloat) },
        alt: { type: QLFloat, defaultValue: 0 },
      }
    });

*/
type QLInputObject struct {
	Name        string
	Description string

	config ql.InputObject
	fields map[string]*InputObjectField
}

type InputObjectField struct {
	Name         string
	Type         QLInputType
	DefaultValue interface{}
	Description  string
}

func NewQLInputObject(config ql.InputObject) *QLInputObject {
	if config.Name == "" {
		throw("Type must be named.")
	}
	assertValidName(config.Name)
	return &QLInputObject{
		Name:        config.Name,
		Description: config.Description,
		config:      config,
	}
}

func (g *QLInputObject) String() string {
	return g.Name
}

func (g *QLInputObject) GetFields() map[string]*InputObjectField {
	if g.fields == nil {
		g.fields = g.defineFieldMap()
	}
	return g.fields
}

func (g *QLInputObject) defineFieldMap() map[string]*InputObjectField {
	config := g.config
	fieldsConfig := config.Fields
	if config.FieldsFunc != nil {
		if fieldsConfig != nil {
			throw(`%v must provide "Fields" or "FieldsFn", not both.`, config.Name)
		}
		fieldsConfig = config.FieldsFunc()
	}
	if len(fieldsConfig) == 0 {
		throw(`%v must provide "Fields" or "FieldsFn".`, config.Name)
	}
	result := make(map[string]*InputObjectField)
	for name, fieldConfig := range fieldsConfig {
		assertValidName(name)
		field := &InputObjectField{
			Name:         name,
			Type:         NewQLInputType(fieldConfig.Type),
			Description:  fieldConfig.Description,
			DefaultValue: fieldConfig.DefaultValue,
		}
		result[name] = field
	}
	return result
}

/*
List Modifier

A list is a kind of type marker, a wrapping type which points to another
type. Lists are often created within the context of defining the fields of
an object type.

Example:

    var PersonType = new QLObject({
      name: 'Person',
      fields: () => ({
        parents: { type: new QLList(Person) },
        children: { type: new QLList(Person) },
      })
    })

*/
type QLList struct {
	OfType QLType
}

func NewQLList(ofType QLType) *QLList {
	if ofType == nil {
		throw(`Can only create List of a QLType but got: %v.`, ofType)
	}
	return &QLList{OfType: ofType}
}

func (g *QLList) String() string {
	return fmt.Sprint("[", g.OfType, "]")
}

/*
Non-Null Modifier

A non-null is a kind of type marker, a wrapping type which points to another
type. Non-null types enforce that their values are never null and can ensure
an error is raised if this ever occurs during a request. It is useful for
fields which you can make a strong guarantee on non-nullability, for example
usually the id field of a database row will never be null.

Example:

    var RowType = new QLObject({
      name: 'Row',
      fields: () => ({
        id: { type: new QLNonNull(QLString) },
      })
    })

Note: the enforcement of non-nullability occurs within the executor.
*/
type QLNonNull struct {
	OfType QLType
}

func NewQLNonNull(typ QLType) *QLNonNull {
	if _, ok := typ.(*QLNonNull); ok {
		throw(`Can only create NonNull of a Nullable QLType but got %v`, typ)
	}
	return &QLNonNull{OfType: typ}
}

func (g *QLNonNull) String() string {
	return fmt.Sprint(g.OfType, "!")
}

// Helpers

var NAME_RX = regexp.MustCompile(`^[_a-zA-Z][_a-zA-Z0-9]*$`)

func assertValidName(name string) {
	if !NAME_RX.MatchString(name) {
		throw(`Names must match /^[_a-zA-Z][_a-zA-Z0-9]*$/ but %v does not.`, name)
	}
}

func throw(format string, args ...interface{}) {
	panic(errors.New(fmt.Sprintf(format, args...)))
}
