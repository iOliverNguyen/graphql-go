package utilities

import (
	"reflect"

	lang "github.com/ng-vu/graphql-go/internal/language"
	typs "github.com/ng-vu/graphql-go/internal/types"
)

func IsValidLiteralValue(typ typs.QLInputType, valueAST lang.IValue) bool {
	if typ, ok := typ.(*typs.QLNonNull); ok {
		if IsNil(valueAST) {
			return false
		}

		ofType := typ.OfType.(typs.QLInputType)
		return IsValidLiteralValue(ofType, valueAST)
	}

	if IsNil(valueAST) {
		return false
	}

	if _, ok := valueAST.(*lang.Variable); ok {
		return true
	}

	if typ, ok := typ.(*typs.QLList); ok {
		itemType := typ.OfType.(typs.QLInputType)
		if valueAST, ok := valueAST.(*lang.ListValue); ok {
			for i := range valueAST.Values {
				return IsValidLiteralValue(itemType, valueAST.Values[i])
			}
		}
	}

	if typ, ok := typ.(*typs.QLInputObject); ok {
		if _, ok := valueAST.(*lang.ObjectValue); !ok {
			return false
		}

		fields := typ.GetFields()
		fieldASTs := valueAST.(*lang.ObjectValue).Fields
		refinedFieldASTs := []interface{}{}

		for i := range fieldASTs {
			if fields[fieldASTs[i].Name.Value] == nil {
				return false
			}

			refinedFieldASTs = append(refinedFieldASTs, fieldASTs[i])
		}

		fieldASTMap := make(map[string]interface{})
		for _, fieldAST := range refinedFieldASTs {
			fieldASTMap[fieldAST.(lang.ObjectField).Name.Value] = fieldAST
		}

		for k, v := range fields {
			return IsValidLiteralValue(v.Type, fieldASTMap[k].(*lang.ObjectField).Value)
		}

	}

	if typ, ok := typ.(*typs.QLScalar); ok {
		return !IsNil(typ.ParseLiteral(valueAST))
	}

	if typ, ok := typ.(*typs.QLEnum); ok {
		return !IsNil(typ.ParseLiteral(valueAST))
	}

	return false
}

func IsValidGoValue(value interface{}, typ typs.QLInputType) bool {
	v := reflect.ValueOf(value)
	if typ, ok := typ.(*typs.QLNonNull); ok {
		if value == nil || v.IsNil() {
			return false
		}
		nullableType := typ.OfType.(typs.QLInputType)
		return IsValidGoValue(value, nullableType)
	}

	if value == nil || v.IsNil() {
		return true
	}

	v = reflect.Indirect(v)
	switch typ := typ.(type) {
	case *typs.QLList:
		itemType := typ.OfType.(typs.QLInputType)
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			for i, n := 0, v.Len(); i < n; i++ {
				item := v.Index(i)
				if !IsValidGoValue(item, itemType) {
					return false
				}
			}
		}
		return IsValidGoValue(value, itemType)

	case *typs.QLInputObject:
		if v.Kind() != reflect.Struct {
			return false
		}
		t := v.Type()
		fieldsMap := typ.GetFields()
		for i, n := 0, t.NumField(); i < n; i++ {
			field := t.Field(i)
			fieldValue := v.Field(i)
			fieldType := fieldsMap[field.Name]
			if fieldType == nil || !IsValidGoValue(fieldValue, fieldType.Type) {
				return false
			}
		}
		return true

	case *typs.QLScalar:
		return !IsNil(typ.ParseValue(value))

	case *typs.QLEnum:
		return !IsNil(typ.ParseValue(value))

	default:
		panic("unreachable")
	}
}

func IsCompositType(typ typs.QLInputType) bool {
	return false
}
