package validation

import (
	lang "github.com/ng-vu/graphql-go/internal/language"
	typs "github.com/ng-vu/graphql-go/internal/types"
	util "github.com/ng-vu/graphql-go/internal/utilities"
)

type RuleVisitor struct {
	Enter func(lang.INode, lang.VisitInfo) *lang.QLError
	Leave func(lang.INode, lang.VisitInfo) *lang.QLError
}

type RuleCreator func(*Context) RuleVisitor

func Validate(
	schema typs.QLSchema,
	ast *lang.Document,
	rules []RuleCreator) []lang.QLError {
	return nil
}

type Context struct {
	schema    typs.QLSchema
	ast       lang.Document
	typeInfo  util.TypeInfo
	fragments map[string]*lang.FragmentDefinition
}

func NewContext(schema typs.QLSchema, ast lang.Document, typeInfo util.TypeInfo) *Context {
	return &Context{
		schema:   schema,
		ast:      ast,
		typeInfo: typeInfo,
	}
}

func (v Context) GetSchema() typs.QLSchema {
	return v.schema
}

func (v Context) GetDocument() lang.Document {
	return v.ast
}

func (v Context) GetFragment(name string) *lang.FragmentDefinition {
	fragments := v.fragments
	if fragments == nil {
		fragments := make(map[string]*lang.FragmentDefinition)
		for _, statement := range v.GetDocument().Definitions {
			if statement, ok := statement.(*lang.FragmentDefinition); ok {
				fragments[statement.Name.Value] = statement
			}
		}
		v.fragments = fragments
	}
	fragment, ok := fragments[name]
	if ok {
		return fragment
	}
	return nil
}

func (v Context) GetType() typs.QLOutputType {
	return v.typeInfo.GetType()
}

func (v Context) GetParentType() typs.QLCompositeType {
	return v.typeInfo.GetParentType()
}

func (v Context) GetInputType() typs.QLInputType {
	return v.typeInfo.GetInputType()
}

func (v Context) GetFieldDef() *typs.QLFieldDefinition {
	return v.typeInfo.GetFieldDef()
}

func (v Context) GetDirective() *typs.QLDirective {
	return v.typeInfo.GetDirective()
}

func (v Context) GetArgument() *typs.QLArgument {
	return v.typeInfo.GetArgument()
}
