package types

import (
	"github.com/ng-vu/graphql-go/ql"
)

type QLSchema struct {
	queryType    *QLObject
	mutationType *QLObject
	directives   []*QLDirective

	typeMap map[string]QLType
}

func NewQLSchema(query ql.Object, mutation *ql.Object) QLSchema {
	queryType := NewQLObject(query)
	var mutationType *QLObject
	if mutation != nil {
		mutationType = NewQLObject(*mutation)
	}

	typeMap := make(map[string]QLType)
	if mutationType == nil {
		typeMapReducer(typeMap, queryType)
	} else {
		typeMapReducer(typeMap, queryType, mutationType)
	}

	for _, typ := range typeMap {
		if typ, ok := typ.(*QLObject); ok {
			for _, iface := range typ.GetInterfaces() {
				assertObjectImplementsInterface(typ, iface)
			}
		}
	}

	directives := []*QLDirective{
		QLIncludeDirective,
		QLSkipDirective,
	}

	return QLSchema{
		queryType:    queryType,
		mutationType: mutationType,
		directives:   directives,
		typeMap:      typeMap,
	}
}

func (g QLSchema) GetQueryType() *QLObject {
	return g.queryType
}

func (g QLSchema) GetMutationType() *QLObject {
	return g.mutationType
}

func (g QLSchema) GetTypeMap() map[string]QLType {
	return g.typeMap
}

func (g QLSchema) GetType(name string) QLType {
	typ, ok := g.typeMap[name]
	if ok {
		return typ
	}
	return nil
}

func (g QLSchema) GetDirectives() []*QLDirective {
	return g.directives
}

func (g QLSchema) GetDirective(name string) *QLDirective {
	for _, directive := range g.directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func typeMapReducer(typeMap map[string]QLType, types ...QLType) {
	for _, typ := range types {
		if typ == nil {
			continue
		}
		switch typ := typ.(type) {
		case *QLList:
			typeMapReducer(typeMap, typ.OfType)
			continue
		case *QLNonNull:
			typeMapReducer(typeMap, typ.OfType)
			continue
		}
		name := typ.GetName()
		if _typ, ok := typeMap[name]; ok {
			if _typ != typ {
				throw(`Schema must contain unique named types but contains multiple types named %v`, name)
			}
			continue
		}
		typeMap[name] = typ

		switch typ := typ.(type) {
		case QLAbstractType:
			possibleTypes := typ.GetPossibleTypes()
			innerTypes := make([]QLType, len(possibleTypes))
			for i, t := range possibleTypes {
				innerTypes[i] = t
			}
			typeMapReducer(typeMap, innerTypes...)

		case *QLObject:
			possibleTypes := typ.GetInterfaces()
			innerTypes := make([]QLType, len(possibleTypes))
			for i, t := range possibleTypes {
				innerTypes[i] = t
			}
			typeMapReducer(typeMap, innerTypes...)
		}

		switch typ := typ.(type) {
		case *QLObject:
			for _, field := range typ.GetFields() {
				argTypes := field.Args
				innerTypes := make([]QLType, len(argTypes))
				for i, t := range argTypes {
					innerTypes[i] = t
				}
				typeMapReducer(typeMap, innerTypes...)
				typeMapReducer(typeMap, field.Type)
			}

		case *QLInterface:
			for _, field := range typ.GetFields() {
				argTypes := field.Args
				innerTypes := make([]QLType, len(argTypes))
				for i, t := range argTypes {
					innerTypes[i] = t
				}
				typeMapReducer(typeMap, innerTypes...)
				typeMapReducer(typeMap, field.Type)
			}

		case *QLInputObject:
			for _, field := range typ.GetFields() {
				typeMapReducer(typeMap, field.Type)
			}
		}
	}
}

func assertObjectImplementsInterface(object *QLObject, iface *QLInterface) {
	objectFieldMap := object.GetFields()
	ifaceFieldMap := iface.GetFields()

	for fieldName, ifaceField := range ifaceFieldMap {
		objectField, ok := objectFieldMap[fieldName]

		if !ok {
			throw(`"%v" expect field "%v" but "%v" does not provide it.`,
				iface, fieldName, object)
		}

		if !isEqualType(ifaceField.Type, objectField.Type) {
			throw(`%v.%v expects type "%v" but %v.%v provides type "%v"`,
				iface, fieldName, ifaceField.Type,
				object, fieldName, objectField.Type)
		}

		for _, ifaceArg := range ifaceField.Args {
			argName := ifaceArg.Name
			ok := false
			for _, objectArg := range objectField.Args {
				if objectArg.Name == argName {
					if !isEqualType(ifaceArg.Type, objectArg.Type) {
						throw(`%v.%v(%v:) expects type "%v" but %v.%v(%v:) provides type "%v"`,
							iface, fieldName, argName, ifaceArg.Type,
							object, fieldName, argName, objectArg.Type)
					}
					ok = true
					break
				}
			}
			if !ok {
				throw(`%v.%v expects argument "%v" but %v.%v does not provide it.`,
					iface, fieldName, argName, object, fieldName)
			}
		}

		for _, objectArg := range objectField.Args {
			argName := objectArg.Name
			ok := false
			for _, ifaceArg := range ifaceField.Args {
				if ifaceArg.Name == argName {
					ok = true
					break
				}
			}
			if !ok {
				throw(`%v.%v does not define argument "%v" but %v.%v provide it.`,
					iface, fieldName, argName, object, fieldName)
			}
		}
	}
}

func isEqualType(typeA, typeB QLType) bool {
	{
		tA, okA := typeA.(*QLNonNull)
		tB, okB := typeB.(*QLNonNull)
		if okA && okB {
			return isEqualType(tA.OfType, tB.OfType)
		}
	}

	{
		tA, okA := typeA.(*QLList)
		tB, okB := typeB.(*QLList)
		if okA && okB {
			return isEqualType(tA.OfType, tB.OfType)
		}
	}

	return typeA == typeB
}
