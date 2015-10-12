package utilities

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	lang "github.com/ng-vu/graphql-go/internal/language"
	typs "github.com/ng-vu/graphql-go/internal/types"
)

func ASTFromValue(value interface{}, typ typs.QLType) lang.IValue {

	if typ, ok := typ.(*typs.QLNonNull); ok {
		return ASTFromValue(value, typ.OfType)
	}

	if IsNil(value) {
		return nil
	}

	val := reflect.ValueOf(value)
	val = reflect.Indirect(val)

	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		var itemType typs.QLType
		if typ, ok := typ.(*typs.QLList); ok {
			itemType = typ.OfType
		}

		values := []lang.IValue{}
		for i := 0; i < val.Len(); i++ {
			values = append(values, ASTFromValue(val.Index(i), itemType))
		}

		return &lang.ListValue{
			Values: values,
		}
	} else if typ, ok := typ.(*typs.QLList); ok {
		return ASTFromValue(value, typ.OfType)
	}

	if val.Kind() == reflect.Bool {
		return &lang.BooleanValue{
			Value: strconv.FormatBool(val.Bool()),
		}
	}

	if _, ok := typ.(*typs.QLScalar); ok {
		if val.Kind() == reflect.Int {
			return &lang.IntValue{
				Value: fmt.Sprintf("%v", val.Int()),
			}
		}

		if val.Kind() == reflect.Float64 {
			return &lang.FloatValue{
				Value: fmt.Sprintf("%v", val.Float()),
			}
		}
	}

	if val.Kind() == reflect.String {

		if _, ok := typ.(*typs.QLEnum); ok {
			matched, _ := regexp.MatchString(`^[_a-zA-Z][_a-zA-Z0-9]*$`, val.String())
			if matched {
				return &lang.EnumValue{
					Value: val.String(),
				}
			}
		}
		return &lang.StringValue{
			Value: val.String()[1 : val.Len()-1],
		}
	}

	v := val.Elem()

	fields := []*lang.ObjectField{}
	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		typeField := val.Type().Field(i)

		var fieldTyp typs.QLType

		if typ, ok := typ.(*typs.QLInputObject); ok {
			fieldDef := typ.GetFields()[typeField.Name]
			if fieldDef != nil {
				fieldTyp = fieldDef.Type
			}
		}

		fieldValue := ASTFromValue(valueField, fieldTyp)

		if fieldValue != nil {
			fields = append(fields, &lang.ObjectField{
				Name: &lang.Name{
					Value: typeField.Name,
				},
				Value: fieldValue,
			})
		}
	}

	return &lang.ObjectValue{
		Fields: fields,
	}

	return nil
}

func BuildASTSchema() *typs.QLSchema {
	return nil
}

func BuildClientSchema() *typs.QLSchema {
	return nil
}

func TypeFromAST(schema typs.QLSchema, inputTypeAST lang.IType) typs.QLType {
	switch typ := inputTypeAST.(type) {
	case *lang.ListType:
		innerType := TypeFromAST(schema, typ.Type)
		if innerType == nil {
			return nil
		}
		return typs.NewQLList(innerType)

	case *lang.NonNullType:
		innerType := TypeFromAST(schema, typ.Type)
		if innerType == nil {
			return nil
		}
		return typs.NewQLNonNull(innerType)

	case *lang.NamedType:
		return schema.GetType(typ.Name.Value)

	default:
		throw("Must be a named type.")
		return nil
	}
}

/**
 * Produces a Go value given a GraphQL Value AST.
 *
 * A GraphQL type must be provided, which will be used to interpret different
 * GraphQL Value literals.
 *
 * | GraphQL Value        | JSON Value    |
 * | -------------------- | ------------- |
 * | Input Object         | Object        |
 * | List                 | Array         |
 * | Boolean              | Boolean       |
 * | String / Enum Value  | String        |
 * | Int / Float          | Number        |
 *
 */
func ValueFromAST(
	valueAST lang.IValue,
	typ typs.QLInputType,
	variables map[string]interface{},
) interface{} {
	if typ, ok := typ.(*typs.QLNonNull); ok {
		nullableType := typ.OfType.(typs.QLInputType)
		return ValueFromAST(valueAST, nullableType, variables)
	}

	if IsNil(valueAST) {
		return nil
	}

	if valueAST, ok := valueAST.(*lang.Variable); ok {
		variableName := valueAST.Name.Value
		return variables[variableName]
	}

	switch typ := typ.(type) {
	case *typs.QLList:
		if valueAST, ok := valueAST.(*lang.ListValue); ok {
			itemType := typ.OfType.(typs.QLInputType)
			result := make([]interface{}, len(valueAST.Values))[:0]
			for _, itemAST := range valueAST.Values {
				result = append(result, ValueFromAST(itemAST, itemType, variables))
			}
			return result
		}
		return []interface{}{ValueFromAST(valueAST, typ, variables)}

	case *typs.QLInputObject:
		valueAST, ok := valueAST.(*lang.ObjectValue)
		if !ok {
			return nil
		}

		fieldASTMap := make(map[string]*lang.ObjectField)
		for _, fieldAST := range valueAST.Fields {
			fieldASTMap[fieldAST.Name.Value] = fieldAST
		}

		result := make(map[string]interface{})
		fields := typ.GetFields()
		for fieldName, field := range fields {
			fieldAST := fieldASTMap[fieldName]
			var fieldValue interface{}
			if fieldAST != nil {
				fieldValue = ValueFromAST(fieldAST.Value, field.Type, variables)
			}
			if IsNil(fieldValue) {
				fieldValue = field.DefaultValue
			}
			if !IsNil(fieldValue) {
				result[fieldName] = fieldValue
			}
		}
		return result

	case *typs.QLScalar:
		parsed := typ.ParseLiteral(valueAST)
		if IsNil(parsed) {
			return nil
		}
		return parsed

	default:
		throw("Must be input type")
		return nil
	}
}
