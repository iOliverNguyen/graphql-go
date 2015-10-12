package utilities

import (
	lang "github.com/ng-vu/graphql-go/internal/language"
	typs "github.com/ng-vu/graphql-go/internal/types"
)

type TypeInfo struct {
	schema          typs.QLSchema
	typeStack       []typs.QLOutputType
	parentTypeStack []typs.QLCompositeType
	inputTypeStack  []typs.QLInputType
	fieldDefStack   []*typs.QLFieldDefinition
	directive       *typs.QLDirective
	argument        *typs.QLArgument
}

func NewTypeInfo(schema typs.QLSchema) TypeInfo {
	return TypeInfo{
		schema:          schema,
		typeStack:       make([]typs.QLOutputType, 0),
		parentTypeStack: make([]typs.QLCompositeType, 0),
		inputTypeStack:  make([]typs.QLInputType, 0),
		fieldDefStack:   make([]*typs.QLFieldDefinition, 0),
		directive:       nil,
		argument:        nil,
	}
}

func (t TypeInfo) GetType() typs.QLOutputType {
	s := t.typeStack
	l := len(s)
	if l > 0 {
		return s[l-1]
	}
	return nil
}

func (t TypeInfo) GetParentType() typs.QLCompositeType {
	s := t.parentTypeStack
	l := len(s)
	if l > 0 {
		return s[l-1]
	}
	return nil
}

func (t TypeInfo) GetInputType() typs.QLInputType {
	s := t.inputTypeStack
	l := len(s)
	if l > 0 {
		return s[l-1]
	}
	return nil
}

func (t TypeInfo) GetFieldDef() *typs.QLFieldDefinition {
	s := t.fieldDefStack
	l := len(s)
	if l > 0 {
		return s[l-1]
	}
	return nil
}

func (t TypeInfo) GetDirective() *typs.QLDirective {
	return t.directive
}

func (t TypeInfo) GetArgument() *typs.QLArgument {
	return t.argument
}

func (t TypeInfo) Enter(node lang.INode) {

}

func (t TypeInfo) Leave(node lang.INode) {

}
