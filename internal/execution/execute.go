package execution

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	debug "github.com/ng-vu/graphql-go/internal/debug"
	lang "github.com/ng-vu/graphql-go/internal/language"
	typs "github.com/ng-vu/graphql-go/internal/types"
)

var LOG = debug.New("graphql/execution")

type _Context struct {
	Schema         typs.QLSchema
	Fragments      map[string]*lang.FragmentDefinition
	RootValue      interface{}
	Operation      *lang.OperationDefinition
	VariableValues map[string]interface{}
	Errors         []error
}

type Result struct {
	Data   interface{}
	Errors []error
}

type Options struct {
	RootValue      interface{}
	VariableValues map[string]interface{}
	OperationName  string
}

func Execute(schema typs.QLSchema, documentAST *lang.Document, opts Options) Result {
	context := newContext(schema, documentAST,
		opts.RootValue, opts.VariableValues, opts.OperationName)
	return context.executeOperation()
}

func newContext(
	schema typs.QLSchema,
	documentAST *lang.Document,
	rootValue interface{},
	rawVariableValues map[string]interface{},
	operationName string,
) *_Context {

	operations := make(map[string]*lang.OperationDefinition)
	fragments := make(map[string]*lang.FragmentDefinition)

	for _, statement := range documentAST.Definitions {
		switch statement := statement.(type) {
		case *lang.OperationDefinition:
			name := ""
			if statement.Name != nil {
				name = statement.Name.Value
			}
			operations[name] = statement
		case *lang.FragmentDefinition:
			fragments[statement.Name.Value] = statement
		default:
			panic(lang.NewQLError(
				fmt.Sprintf(`Cannot execute a request containing a %v.`, statement.Kind()),
				[]lang.INode{statement}))
		}
	}
	if operationName != "" && len(operations) != 1 {
		panic(lang.NewQLError(`Must provide operation name if query contains multiple operations.`, nil))
	}
	if operationName == "" {
		for name := range operations {
			operationName = name
		}
	}
	operation, ok := operations[operationName]
	if !ok {
		panic(lang.NewQLError(
			fmt.Sprintf(`Unknown operation named "%v".`, operationName), nil))
	}
	variableValues := GetVariableValues(schema, operation.VariableDefinitions, rawVariableValues)
	return &_Context{
		Schema:         schema,
		Fragments:      fragments,
		RootValue:      rootValue,
		Operation:      operation,
		VariableValues: variableValues,
		Errors:         nil,
	}
}

func (c *_Context) executeOperation() Result {
	operation := c.Operation
	typ := getOperationRootType(c.Schema, operation)
	fields := c.collectFields(typ, operation.SelectionSet,
		map[string][]*lang.Field{}, map[string]struct{}{})

	if operation.Operation == lang.OperationMutation {
		return c.executeFieldsSerially(typ, c.RootValue, fields)
	}
	return c.executeFields(typ, c.RootValue, fields)
}

func (c *_Context) executeFieldsSerially(
	parentType *typs.QLObject,
	sourceValue interface{},
	fields map[string][]*lang.Field) Result {

	results := make(map[string]interface{})
	for responseName, fieldASTs := range fields {
		result := c.resolveField(parentType, sourceValue, fieldASTs)
		if result != nil {
			results[responseName] = result
		}
	}
	return Result{results, nil}
}

func (c *_Context) executeFields(
	parentType *typs.QLObject,
	sourceValue interface{},
	fields map[string][]*lang.Field) Result {

	var m sync.Mutex
	var wg sync.WaitGroup
	results := make(map[string]interface{})
	for responseName, fieldASTs := range fields {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := c.resolveField(parentType, sourceValue, fieldASTs)
			if result != nil {
				m.Lock()
				results[responseName] = result
				m.Unlock()
			}
		}()
	}
	wg.Wait()
	return Result{results, nil}
}

func (c *_Context) collectFields(
	runtimeType *typs.QLObject,
	selectionSet *lang.SelectionSet,
	fields map[string][]*lang.Field,
	visitedFragmentNames map[string]struct{},
) map[string][]*lang.Field {

	for _, selection := range selectionSet.Selections {
		switch selection := selection.(type) {
		case *lang.Field:
			if !c.shouldIncludeNode(selection.Directives) {
				continue
			}
			name := getFieldEntryKey(selection)
			fields[name] = append(fields[name], selection)

		case *lang.InlineFragment:
			if !c.shouldIncludeNode(selection.Directives) ||
				!c.doesFragmentConditionMatch(selection, runtimeType) {
				continue
			}
			c.collectFields(runtimeType, selection.SelectionSet, fields, visitedFragmentNames)

		case *lang.FragmentSpread:
			fragName := selection.Name.Value
			_, visited := visitedFragmentNames[fragName]
			if visited || !c.shouldIncludeNode(selection.Directives) {
				continue
			}
			visitedFragmentNames[fragName] = struct{}{}
			fragment, ok := c.Fragments[fragName]
			if !ok || c.shouldIncludeNode(fragment.Directives) ||
				c.doesFragmentConditionMatch(fragment, runtimeType) {
				continue
			}
			c.collectFields(runtimeType, fragment.SelectionSet, fields, visitedFragmentNames)
		}
	}
	return fields
}

func (c *_Context) shouldIncludeNode(directives []*lang.Directive) bool {
	return true
	// skipName := typs.QLSkipDirective.Name
	// var skipAST *lang.Directive
	// for _, directive := range directives {
	// 	if directive.Name.Value == skipName {
	// 		skipAST = directive
	// 		break
	// 	}
	// }
	// if skipAST != nil {
	// 	argValues := GetArgumentValues(typs.QLSkipDirective.Args, skipAST.Arguments, c.VariableValues)
	// 	skipIf := argValues["if"]
	// 	return skipIf != nil // TODO(qv): Check trust value
	// }

	// var includeAST *lang.Directive
	// for _, directive := range directives {
	// 	if directive.Name.Value == typs.QLIncludeDirective.Name {
	// 		includeAST = directive
	// 		break
	// 	}
	// }
	// if includeAST != nil {
	// 	argValues := GetArgumentValues(typs.QLIncludeDirective.Args, includeAST.Arguments, c.VariableValues)
	// 	includeIf := argValues["if"]
	// 	return includeIf != nil // TODO(qv): Check trust value
	// }
	// return false
}

func (c *_Context) doesFragmentConditionMatch(fragment lang.ITypeCondition, typ *typs.QLObject) bool {
	return true
	// 	conditionalType := util.TypeFromAST(c.Schema, fragment.GetTypeCondition())
	// 	if conditionalType == typ {
	// 		return true
	// 	}
	// 	// TODO(qv): Check isAbstractType
	// 	if conditionalType, ok := conditionalType.(typs.QLAbstractType); ok {
	// 		return conditionalType.IsPossibleType(typ)
	// 	}
	// 	return false
}

func (c *_Context) resolveField(
	parentType *typs.QLObject,
	source interface{},
	fieldASTs []*lang.Field,
) interface{} {

	fieldAST := fieldASTs[0]
	fieldName := fieldAST.Name.Value
	fieldDef := getFieldDef(c.Schema, parentType, fieldName)
	if fieldDef == nil {
		return nil
	}

	returnType := fieldDef.Type
	resolveFn := fieldDef.Resolve
	if resolveFn == nil {
		resolveFn = defaultResolveFn
	}

	args := GetArgumentValues(fieldDef.Args, fieldAST.Arguments, c.VariableValues)
	info := typs.QLResolveInfo{
		FieldName:      fieldName,
		FieldASTs:      fieldASTs,
		ReturnType:     returnType,
		ParentType:     parentType,
		Schema:         c.Schema,
		Fragments:      c.Fragments,
		RootValue:      c.RootValue,
		Operation:      c.Operation,
		VariableValues: c.VariableValues,
	}

	defer func() {
		err := recover()
		if err == nil {
			return
		}

		var reportedError lang.QLError
		nodes := make([]lang.INode, len(fieldASTs))
		for i, ast := range fieldASTs {
			nodes[i] = ast
		}
		if err, ok := err.(error); ok {
			reportedError = lang.LocatedError(err, nodes)
		} else {
			reportedError = lang.LocatedError(errors.New(fmt.Sprint(err)), nodes)
		}

		if _, ok := returnType.(*typs.QLNonNull); ok {
			panic(reportedError)
		}
		c.Errors = append(c.Errors, reportedError)
	}()

	result := resolveFn(source, args, info)
	fields := make([]*lang.Field, len(fieldASTs))
	for i, ast := range fieldASTs {
		fields[i] = ast
	}
	return c.completeValueCatchingError(returnType, fieldASTs, info, result)
}

func (c *_Context) completeValueCatchingError(
	returnType typs.QLType,
	fieldASTs []*lang.Field,
	info typs.QLResolveInfo,
	result interface{},
) interface{} {

	// If the field type is non-nullable, then it is resolved without any
	// protection from errors.
	if _, ok := returnType.(*typs.QLNonNull); ok {
		return c.completeValue(returnType, fieldASTs, info, result)
	}

	// Otherwise, error protection is applied, logging the error and resolving
	// a null value for this field if one is encountered.
	defer func() {
		err := recover()
		if err != nil {
			LOG.Println("Catched error:", err)
			if err, ok := err.(error); ok {
				c.Errors = append(c.Errors, err)
				return
			}
			c.Errors = append(c.Errors, errors.New(fmt.Sprint(err)))
		}
	}()

	completed := c.completeValue(returnType, fieldASTs, info, result)
	return completed
}

func (c *_Context) completeValue(
	returnType typs.QLType,
	fieldASTs []*lang.Field,
	info typs.QLResolveInfo,
	result interface{},
) interface{} {

	if returnType, ok := returnType.(*typs.QLNonNull); ok {
		completed := c.completeValue(returnType.OfType, fieldASTs, info, result)
		if completed == nil {
			nodes := make([]lang.INode, len(fieldASTs))
			for i, node := range fieldASTs {
				nodes[i] = node
			}
			panic(lang.NewQLError(
				fmt.Sprintf(
					`Cannot return null for non-nullable field %v.%v.`,
					info.ParentType, info.FieldName),
				nodes))
		}
		return completed
	}

	v := reflect.ValueOf(result)
	if result == nil {
		return nil
	}

	v = reflect.Indirect(v)
	t := v.Type()
	switch returnType := returnType.(type) {
	case *typs.QLList:
		if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
			itemType := returnType.OfType
			list := make([]interface{}, t.NumField())
			for i, n := 0, t.NumField(); i < n; i++ {
				list[i] = c.completeValueCatchingError(itemType, fieldASTs, info, v.Field(i).Interface())
			}
			return list
		}
		panic("User Error: expected array of slice, but did not find one.")

	case *typs.QLScalar:
		if returnType.Serialize == nil {
			panic("Missing serialize method on type")
		}
		serializedResult := returnType.Serialize(result)
		if serializedResult == nil {
			return nil
		}
		return serializedResult

	case *typs.QLEnum:
		return returnType.Serialize(result)

	case *typs.QLObject:
		runtimeType := returnType
		subFieldASTs := map[string][]*lang.Field{}
		visitedFragmentNames := map[string]struct{}{}
		for _, fieldAST := range fieldASTs {
			selectionSet := fieldAST.SelectionSet
			if selectionSet != nil {
				subFieldASTs = c.collectFields(runtimeType, selectionSet, subFieldASTs, visitedFragmentNames)
			}
		}
		return c.executeFields(runtimeType, result, subFieldASTs)

	case typs.QLAbstractType:
		panic("not implemented")

	default:
		panic("unreachable")
	}
}

func getOperationRootType(
	schema typs.QLSchema,
	operation *lang.OperationDefinition,
) *typs.QLObject {
	switch operation.Operation {
	case lang.OperationQuery:
		return schema.GetQueryType()
	case lang.OperationMutation:
		murationType := schema.GetMutationType()
		if murationType == nil {
			panic(lang.NewQLError(
				"Schema is not configured for mutations",
				[]lang.INode{operation}))
		}
		return murationType
	default:
		panic(lang.NewQLError(
			"Can only execute queries and mutation",
			[]lang.INode{operation}))
	}
}

func getFieldEntryKey(node *lang.Field) string {
	if node.Alias != nil {
		return node.Alias.Value
	}
	return node.Name.Value
}

/**
 * If a resolve function is not given, then a default resolve behavior is used
 * which takes the property of the source object of the same name as the field
 * and returns it as the result, or if it's a function, returns the result
 * of calling that function.
 */
func defaultResolveFn(
	source interface{},
	args map[string]interface{},
	info typs.QLResolveInfo,
) interface{} {
	v := reflect.Indirect(reflect.ValueOf(source))
	if v.Kind() != reflect.Struct {
		return nil
	}
	field := v.FieldByName(info.FieldName)
	if !field.IsValid() {
		return nil
	}
	return field.Interface()
}

func getFieldDef(
	schema typs.QLSchema,
	parentType *typs.QLObject,
	fieldName string,
) *typs.QLFieldDefinition {
	if fieldName == typs.SchemaMetaFieldDef.Name &&
		schema.GetQueryType() == parentType {
		return typs.SchemaMetaFieldDef
	}
	if fieldName == typs.TypeMetaFieldDef.Name &&
		schema.GetQueryType() == parentType {
		return typs.TypeMetaFieldDef
	}
	if fieldName == typs.TypeNameMetaFieldDef.Name {
		return typs.TypeNameMetaFieldDef
	}
	return parentType.GetFields()[fieldName]
}
