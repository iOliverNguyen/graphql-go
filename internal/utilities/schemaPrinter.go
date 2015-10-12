package utilities

import (
	"fmt"
	"strings"

	typs "github.com/ng-vu/graphql-go/internal/types"
)

func PrintSchema(schema typs.QLSchema) string {
	return ""
}

func IsDefinedType(typename string) bool {
	return !IsIntrospectionType(typename) && !IsBuiltInScalar(typename)
}

func IsIntrospectionType(typename string) bool {
	return strings.Index(typename, "__") == 0
}

func IsBuiltInScalar(typename string) bool {
	return typename == "String" || typename == "Boolean" || typename == "Int" || typename == "Float" || typename == "ID"
}

func PrintFilteredSchema(
	schema typs.QLSchema,
	typeFilter func(string) bool,
) string {
	// typeMap := schema.GetTypeMap()

	// for k, v := range typeMap {
	// }

	return ""
}

func PrintType(typ typs.QLNamedType) string {
	return ""
}

func PrintScalar(typ typs.QLScalar) string {
	return fmt.Sprintf(`scalar %v`, typ)
}

func PrintObject(typ typs.QLObject) string {
	// interfaces := typ.GetInterfaces()
	// implementedInterfaces := ""
	// if len(interfaces) > 0 {
	// 	interfaceNames := func(ins []*QLInterface) []string {
	// 		result := []string{}
	// 		for i := range ins {
	// 			result = append(result, ins[i].Name)
	// 		}

	// 		return result
	// 	}(interfaces)

	// 	implementedInterfaces := " implements " + strings.Join(interfaceNames, ", ")
	// }
	return ""
}

func PrintInterface(typ typs.QLInterface) string {
	return ""
}

func PrintUnion(typ typs.QLUnion) string {
	return ""
}

func PrintEnum(typ typs.QLEnum) string {
	return ""
}

func PrintInputObject(typ typs.QLInputObject) string {
	return ""
}

func PrintFields() {

}

func PrintArgs() {

}

func PrintInputValue() {

}
