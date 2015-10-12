package execution

import (
	"reflect"

	lang "github.com/ng-vu/graphql-go/internal/language"
	typs "github.com/ng-vu/graphql-go/internal/types"
	util "github.com/ng-vu/graphql-go/internal/utilities"
)

func GetVariableValues(
	schema typs.QLSchema,
	definitionASTs []*lang.VariableDefinition,
	inputs map[string]interface{},
) map[string]interface{} {
	values := make(map[string]interface{})
	for _, defAST := range definitionASTs {
		varName := defAST.Variable.Name.Value
		values[varName] = getVariableValue(schema, defAST, inputs[varName])
	}
	return values
}

func GetArgumentValues(
	argDefs []*typs.QLArgument,
	argASTs []*lang.Argument,
	variableValues map[string]interface{},
) map[string]interface{} {
	if len(argDefs) == 0 || len(argASTs) == 0 {
		return nil
	}
	argASTMap := make(map[string]*lang.Argument)
	for _, arg := range argASTs {
		argASTMap[arg.Name.Value] = arg
	}
	result := make(map[string]interface{})
	for _, argDef := range argDefs {
		var valueAST lang.IValue
		name := argDef.Name
		if argAST, ok := argASTMap[name]; ok {
			valueAST = argAST.Value
		}
		value := util.ValueFromAST(valueAST, argDef.Type, variableValues)
		if value == nil {
			value = argDef.DefaultValue
		}
		if value != nil {
			result[name] = value
		}
	}
	return result
}

func getVariableValue(
	schema typs.QLSchema,
	definitionAST *lang.VariableDefinition,
	input interface{}) interface{} {
	return nil
	// typ := util.TypeFromAST(schema, definitionAST.Type)
	// variable := definitionAST.Variable
	// if _, ok := typ.(typs.QLInputType); !ok {
	// 	panic(lang.NewQLError(
	// 		fmt.Sprintf(
	// 			`Variable "$%v" expected value of type "%v" which cannot be used as an input type.`,
	// 			variable.Name.Value, lang.Print(definitionAST.Type)),
	// 		[]lang.INode{definitionAST}))
	// }
	// // TODO: isValidJSValue
	// return coerceValue(typ, input)
	// if input == nil {
	// 	panic(lang.NewQLError(
	// 		fmt.Sprintf(
	// 			`Variable "$%v" of required type "$%v" was not provided.`,
	// 			variable.Name.Value, lang.Print(definitionAST.Type)),
	// 		[]lang.INode{definitionAST}))
	// }
	// panic(lang.NewQLError(
	// 	fmt.Sprintf(
	// 		`Variable "$%v" expected value of type "%v" but got: %v`,
	// 		variable.Name.Value, lang.Print(definitionAST.Type)),
	// 	[]lang.INode{definitionAST}))
}

/**
 * Given a type and any value, return a runtime value coerced to match the type.
 */
func coerceValue(typ typs.QLType, value interface{}) interface{} {
	if typ, ok := typ.(*typs.QLNonNull); ok {
		nullableType := typ.OfType
		return coerceValue(nullableType, value)
	}

	v := reflect.ValueOf(value)
	if value == nil || v.IsNil() {
		return nil
	}

	switch typ := typ.(type) {
	case *typs.QLList:
		itemType := typ.OfType
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
			result := make([]interface{}, v.Len())[:0]
			for i, n := 0, v.Len(); i < n; i++ {
				fieldValue := v.Field(i)
				result = append(result, coerceValue(itemType, fieldValue))
			}
			return result
		}
		return []interface{}{coerceValue(itemType, value)}

	case *typs.QLInputObject:
		fields := typ.GetFields()
		obj := make(map[string]interface{})
		for fieldName, field := range fields {
			runtimeField := v.FieldByName(fieldName)
			if !runtimeField.IsValid() {
				continue
			}

			fieldValue := coerceValue(field.Type, runtimeField.Interface())
			if fieldValue == nil {
				fieldValue = field.DefaultValue
			}
			if fieldValue != nil {
				obj[fieldName] = fieldValue
			}
		}
		return obj

	case *typs.QLScalar:
		return typ.ParseValue(value)

	case *typs.QLEnum:
		return typ.ParseValue(value)

	default:
		panic("Must be input type")
	}
}
