package types

import (
	"github.com/ng-vu/graphql-go/ql"
)

var __Schema = NewQLObject(__SchemaConfig)
var __SchemaConfig = ql.Object{
	Name:        "__Schema",
	Description: "A QL Schema defines the capabilities of a QL server. It exposes all available types and directives on the server, as well as the entry points for query and mutation operations.",
	FieldsFunc: func() ql.FieldMap {
		return ql.FieldMap{
			"types": {
				Description: "A list of all types supported by this server.",
				Type:        ql.NonNull{ql.List{ql.NonNull{__TypeConfig}}},
				Resolve: func(schema QLSchema) interface{} {
					typeMap := schema.GetTypeMap()
					result := make(map[string]QLType)
					for key, value := range typeMap {
						result[key] = value
					}
					return result
				},
			},
			"queryType": {
				Description: "The type that query operations will be rooted at.",
				Type:        ql.NonNull{__TypeConfig},
				Resolve: func(schema QLSchema) interface{} {
					return schema.GetQueryType()
				},
			},
			"mutationType": {
				Description: "If this server supports mutation, the type that mutation operations will be rooted at.",
				Type:        __TypeConfig,
				Resolve: func(schema QLSchema) interface{} {
					return schema.GetMutationType()
				},
			},
			"directives": {
				Description: "A list of all directives supported by this server.",
				Type:        ql.NonNull{ql.List{ql.NonNull{__DirectiveConfig}}},
				Resolve: func(schema QLSchema) interface{} {
					return schema.GetDirectives()
				},
			},
		}
	},
}

var __Directive = NewQLObject(__DirectiveConfig)
var __DirectiveConfig = ql.Object{
	Name: "__Directive",
	FieldsFunc: func() ql.FieldMap {
		return ql.FieldMap{
			"name":        {Type: ql.NonNull{ql.String}},
			"description": {Type: ql.String},
			"args": {
				Type: ql.NonNull{ql.List{ql.NonNull{__InputValueConfig}}},
				Resolve: func(directive QLDirective) interface{} {
					return directive.Args
				},
			},
			"onOperation": {Type: ql.NonNull{ql.Boolean}},
			"onFragment":  {Type: ql.NonNull{ql.Boolean}},
			"onField":     {Type: ql.NonNull{ql.Boolean}},
		}
	},
}

var __Field = NewQLObject(__FieldConfig)
var __FieldConfig = ql.Object{
	Name: "__Field",
	FieldsFunc: func() ql.FieldMap {
		return ql.FieldMap{
			"name":        {Type: ql.NonNull{ql.String}},
			"description": {Type: ql.String},
			"args": {
				Type: ql.NonNull{ql.List{ql.NonNull{__InputValueConfig}}},
				Resolve: func(field QLFieldDefinition) interface{} {
					return field.Args
				},
			},
			"type": {Type: ql.NonNull{__TypeConfig}},
			"isDeprecated": {
				Type: ql.NonNull{ql.Boolean},
				Resolve: func(field QLFieldDefinition) interface{} {
					return field.DeprecationReason != ""
				},
			},
			"deprecationReason": {Type: ql.String},
		}
	},
}

var __InputValue = NewQLObject(__InputValueConfig)
var __InputValueConfig = ql.Object{
	Name: "__InputValue",
	FieldsFunc: func() ql.FieldMap {
		return ql.FieldMap{
			"name":        {Type: ql.NonNull{ql.String}},
			"description": {Type: ql.String},
			"type":        {Type: ql.NonNull{__TypeConfig}},
			"defaultValue": {
				Type: ql.String,
				Resolve: func(inputVal *QLInputObject) interface{} {
					// return inputVal.DefaultValue
					// TODO
					return nil
				},
			},
		}
	},
}

var __EnumValue = NewQLObject(__EnumValueConfig)
var __EnumValueConfig = ql.Object{
	Name: "__EnumValue",
	FieldsFunc: func() ql.FieldMap {
		return ql.FieldMap{
			"name":        {Type: ql.NonNull{ql.String}},
			"description": {Type: ql.String},
			"isDeprecated": {
				Type: ql.NonNull{ql.Boolean},
				Resolve: func(enumValue QLEnumValueDefinition) interface{} {
					return enumValue.DeprecationReason
				},
			},
		}
	},
}

const (
	TYPE_SCALAR       = "SCALAR"
	TYPE_OBJECT       = "OBJECT"
	TYPE_INTERFACE    = "INTERFACE"
	TYPE_UNION        = "UNION"
	TYPE_ENUM         = "ENUM"
	TYPE_INPUT_OBJECT = "INPUT_OBJECT"
	TYPE_LIST         = "LIST"
	TYPE_NON_NULL     = "NON_NULL"
)

var __TypeKind = NewQLEnum(__TypeKindConfig)
var __TypeKindConfig = ql.Enum{
	Name:        "__TypeKind",
	Description: "An enum describing what kind of type a given __Type is",
	Values: ql.EnumValueMap{
		"SCALAR": {
			Value:       TYPE_SCALAR,
			Description: "Indicates this type is a scalar.",
		},
		"OBJECT": {
			Value:       TYPE_OBJECT,
			Description: "Indicates this type is an object. `fields` and `interfaces` are valid fields.",
		},
		"INTERFACE": {
			Value:       TYPE_INTERFACE,
			Description: "Indicates this type is an interface. `fields` and `possibleTypes` are valid fields.",
		},
		"UNION": {
			Value:       TYPE_UNION,
			Description: "Indicates this type is a union. `possibleTypes` is a valid field.",
		},
		"ENUM": {
			Value:       TYPE_ENUM,
			Description: "Indicates this type is an enum. `enumValues` is a valid field.",
		},
		"INPUT_OBJECT": {
			Value:       TYPE_INPUT_OBJECT,
			Description: "Indicates this type is an input object. `inputFields` is a valid field.",
		},
		"LIST": {
			Value:       TYPE_LIST,
			Description: "Indicates this type is a list. `ofType` is a valid field.",
		},
		"NON_NULL": {
			Value:       TYPE_NON_NULL,
			Description: "Indicates this type is a non-null. `ofType` is a valid field.",
		},
	},
}

var __Type *QLObject
var __TypeConfig ql.Object
var SchemaMetaFieldDef, TypeMetaFieldDef, TypeNameMetaFieldDef *QLFieldDefinition

func init() {
	__TypeConfig = ql.Object{
		Name: "__Type",
		FieldsFunc: func() ql.FieldMap {
			return ql.FieldMap{
				"kind": {
					Type: ql.NonNull{__TypeKindConfig},
				},
				"name":        {Type: ql.String},
				"description": {Type: ql.String},
				"fields": {
					Type: ql.List{ql.NonNull{__FieldConfig}},
					Args: ql.ArgumentMap{
						"includeDeprecated": {Type: ql.Boolean, DefaultValue: false},
					},
					Resolve: func(typ QLType, args struct{ IncludeDeprecated bool }) interface{} {
						if typ, ok := typ.(QLObjectInterface); ok {
							fieldMap := typ.GetFields()
							fields := make([]QLFieldDefinition, len(fieldMap))[:0]
							for _, field := range fieldMap {
								if args.IncludeDeprecated || field.DeprecationReason == "" {
									fields = append(fields, field)
								}
							}
							return fields
						}
						return nil
					},
				},
				"interfaces": {
					Type: ql.List{ql.NonNull{__TypeConfig}},
					Resolve: func(typ QLObject) interface{} {
						return typ.GetInterfaces()
					},
				},
				"possibleTypes": {
					Type: ql.List{ql.NonNull{__TypeConfig}},
					Resolve: func(typ QLAbstractType) interface{} {
						return typ.GetPossibleTypes()
					},
				},
				"enumValues": {
					Type: ql.List{ql.NonNull{__EnumValueConfig}},
					Args: ql.ArgumentMap{
						"includeDeprecated": {Type: ql.Boolean, DefaultValue: false},
					},
					// TODO(qv): Resolve
					Resolve: func(typ QLEnum, args struct{ IncludeDeprecated bool }) interface{} {
						vs := typ.GetValues()
						values := make([]*QLEnumValueDefinition, len(vs))[:0]
						for _, v := range vs {
							if args.IncludeDeprecated || v.DeprecationReason == "" {
								values = append(values, v)
							}
						}
						return values
					},
				},
				"inputFields": {
					Type: ql.List{ql.NonNull{__InputValueConfig}},
					Resolve: func(typ *QLInputObject) interface{} {
						fieldMap := typ.GetFields()
						result := make([]*InputObjectField, len(fieldMap))[:0]
						for _, field := range fieldMap {
							result = append(result, field)
						}
						return result
					},
				},
				"ofType": {Type: __TypeConfig},
			}
		},
	}

	__Type = NewQLObject(__TypeConfig)

	SchemaMetaFieldDef = &QLFieldDefinition{
		Name:        "__schema",
		Type:        NewQLOutputType(ql.NonNull{__SchemaConfig}),
		Description: "Access the current type schema of this server.",
		Resolve: func(source interface{}, args map[string]interface{}, info QLResolveInfo) interface{} {
			return info.Schema
		},
	}

	TypeMetaFieldDef = &QLFieldDefinition{
		Name:        "__type",
		Type:        NewQLOutputType(__TypeConfig),
		Description: "Request the type information of a single type.",
		Args:        []*QLArgument{},
		Resolve: func(source interface{}, args map[string]interface{}, info QLResolveInfo) interface{} {
			if name, ok := args["name"].(string); ok {
				return info.Schema.GetType(name)
			}
			return nil
		},
	}

	TypeNameMetaFieldDef = &QLFieldDefinition{
		Name:        "__typename",
		Type:        NewQLOutputType(ql.NonNull{ql.String}),
		Description: "The name of the current Object type at runtime.",
		Resolve: func(source interface{}, args map[string]interface{}, info QLResolveInfo) interface{} {
			return info.ParentType.GetName()
		},
	}
}
