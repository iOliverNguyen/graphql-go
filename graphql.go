package graphql

import (
	"fmt"
	"reflect"

	"github.com/ng-vu/graphql-go/internal/debug"
	"github.com/ng-vu/graphql-go/internal/execution"
	"github.com/ng-vu/graphql-go/internal/language"
	"github.com/ng-vu/graphql-go/internal/types"
	"github.com/ng-vu/graphql-go/internal/validation"
	"github.com/ng-vu/graphql-go/ql"
)

var LOG = debug.New("graphql")

type RequestOpts struct {
	RootValue      interface{}
	VariableValues map[string]interface{}
	OperationName  string
}

type Schema struct {
	schema types.QLSchema
}

func NewSchema(query ql.Object, mutations ...ql.Object) (Schema, error) {
	if len(mutations) > 1 {
		panic("graphql: must provide only one mutation object")
	}
	var mutation *ql.Object
	if len(mutations) == 1 {
		mutation = new(ql.Object)
		*mutation = mutations[0]
	}
	schema := types.NewQLSchema(query, mutation)
	return Schema{schema}, nil
}

type Request struct {
	schema      types.QLSchema
	documentAST *language.Document
	opts        RequestOpts
}

func NewRequest(schema Schema, request string, opts ...RequestOpts) (*Request, Errors) {
	var _opts RequestOpts
	if len(opts) > 1 {
		panic("graphql: must provide only one Options object")
	} else if len(opts) == 1 {
		_opts = opts[0]
	}

	source := language.NewSource(request, "GraphQL request")
	documentAST, err := language.Parse(source)
	if err != nil {
		return nil, _Errors{[]error{err}}
	}

	validationErrors := validation.Validate(schema.schema, documentAST, nil)
	if validationErrors != nil {
		errs := make([]error, len(validationErrors))
		for i, e := range validationErrors {
			errs[i] = e
		}
		return nil, _Errors{errs}
	}

	return &Request{schema.schema, documentAST, _opts}, nil
}

func (r *Request) Print() string {
	return language.Print(r.documentAST)
}

func (r *Request) Run(value interface{}) error {
	opts := execution.Options{
		RootValue:      r.opts.RootValue,
		VariableValues: r.opts.VariableValues,
		OperationName:  r.opts.OperationName,
	}
	result := execution.Execute(r.schema, r.documentAST, opts)

	if len(result.Errors) > 0 {
		errors := make([]error, len(result.Errors))
		for i, err := range result.Errors {
			errors[i] = err
		}
		return _Errors{errors}
	}

	v := reflect.Indirect(reflect.ValueOf(value))
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic("Only support struct as argument")
	}

	vResult := reflect.ValueOf(result.Data)
	if vResult.Kind() != reflect.Map {
		panic("Result must be a map")
	}

	for i, n := 0, t.NumField(); i < n; i++ {
		tField := t.Field(i)
		tag := tField.Tag.Get("graphql")
		if tag == "" {
			continue
		}

		vField := v.Field(i)
		resultField := vResult.MapIndex(reflect.ValueOf(tag))
		elemField := resultField.Elem()
		if elemField.Type().AssignableTo(tField.Type) {
			vField.Set(elemField)
		} else {
			LOG.Println("Cannot assign value to field", tField.Name, resultField)
		}
	}
	return nil
}

func Run(schema Schema, request string, opts *RequestOpts, v interface{}) error {
	req, err := NewRequest(schema, request, *opts)
	if err != nil {
		return err
	}
	return req.Run(v)
}

type Errors interface {
	Error() string
	AllErrors() []error
}

type _Errors struct {
	errors []error
}

func (e _Errors) Error() string {
	if len(e.errors) == 0 {
		return ""
	}
	if len(e.errors) == 1 {
		return e.errors[0].Error()
	}
	return fmt.Sprintf("%v (%v more)", e.errors[0].Error(), len(e.errors))
}

func (e _Errors) AllErrors() []error {
	return e.errors
}
