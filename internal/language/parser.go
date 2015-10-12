package language

import (
	"errors"
	"fmt"
)

type ParseOptions struct {
	NoLocation bool
	NoSource   bool
}

func Parse(source Source, options ...ParseOptions) (result *Document, err error) {
	defer func() {
		e := recover()
		if e != nil {
			result = nil
			if e, ok := e.(error); ok {
				err = e
			} else {
				err = errors.New(fmt.Sprint("graphql/parser: ", err))
			}
		}
	}()

	var opts ParseOptions
	if len(options) > 0 {
		opts = options[0]
	}
	parser := newParser(source, opts)
	return parser.parseDocument(), nil
}

func ParseValue(source Source, options ...ParseOptions) (result IValue, err error) {
	defer func() {
		e := recover()
		if e != nil {
			result = nil
			if e, ok := e.(error); ok {
				err = e
			} else {
				err = errors.New(fmt.Sprint("graphql/parser: ", err))
			}
		}
	}()

	var opts ParseOptions
	if len(options) > 0 {
		opts = options[0]
	}
	parser := newParser(source, opts)
	return parser.parseValueLiteral(false), nil
}

type Parser struct {
	lexer   *Lexer
	token   Token
	source  Source
	options ParseOptions
	prevEnd int
}

/**
 * Returns the parser object that is used to store state throughout the
 * process of parsing.
 */
func newParser(source Source, options ParseOptions) *Parser {
	lexer := newLexer(source)
	return &Parser{
		lexer:   lexer,
		token:   lexer.nextToken(),
		source:  source,
		options: options,
		prevEnd: 0,
	}
}

/**
 * Returns a location object, used to identify the place in
 * the source that created a given parsed object.
 */
func (p *Parser) loc(start int) *Location {
	if p.options.NoLocation {
		return nil
	}

	if p.options.NoSource {
		return &Location{
			Start: start,
			End:   p.prevEnd,
		}
	}

	return &Location{
		Start:  start,
		End:    p.prevEnd,
		Source: &p.source,
	}
}

/**
 * Moves the internal parser object to the next lexed token.
 */
func (p *Parser) advance() {
	prevEnd := p.token.End
	p.prevEnd = prevEnd
	p.token = p.lexer.nextTokenFromPosition(prevEnd)
}

/**
 * Determines if the next token is of a given kind
 */
func (p *Parser) peek(kind TokenKind) bool {
	return p.token.Kind == kind
}

/**
 * If the next token is of the given kind, return true after advancing
 * the parser. Otherwise, do not change the parser state and return false.
 */
func (p *Parser) skip(kind TokenKind) bool {
	match := p.token.Kind == kind
	if match {
		p.advance()
	}
	return match
}

/**
 * If the next token is of the given kind, return that token after advancing
 * the parser. Otherwise, do not change the parser state and return false.
 */
func (p *Parser) expect(kind TokenKind) Token {
	token := p.token
	if token.Kind == kind {
		p.advance()
		return token
	}

	panic(SyntaxError(p.source, token.Start,
		fmt.Sprintf("Expected %v, found %v", kind, token)))
}

/**
 * If the next token is a keyword with the given value, return that token after
 * advancing the parser. Otherwise, do not change the parser state and return
 * false.
 */
func (p *Parser) expectKeyword(value string) Token {
	token := p.token
	if token.Kind == TOKEN_NAME && token.Value == value {
		p.advance()
		return token
	}

	panic(SyntaxError(p.source, token.Start,
		fmt.Sprintf("Expected %v, found %v", value, token)))
}

/**
 * Helper function for creating an error when an unexpected lexed token
 * is encountered.
 */
func (p *Parser) unexpected(atToken *Token) error {
	token := atToken
	if atToken == nil {
		token = &p.token
	}

	return SyntaxError(p.source, token.Start,
		fmt.Sprintf("Unexpected %v", token))
}

/**
 * Returns a possibly empty list of parse nodes, determined by
 * the parseFn. This list begins with a lex token of openKind
 * and ends with a lex token of closeKind. Advances the parser
 * to the next lex token after the closing token.
 */
func (p *Parser) any(openKind, closeKind TokenKind, parseFn func()) {
	p.expect(openKind)
	for !p.skip(closeKind) {
		parseFn()
	}
}

/**
 * Returns a non-empty list of parse nodes, determined by
 * the parseFn. This list begins with a lex token of openKind
 * and ends with a lex token of closeKind. Advances the parser
 * to the next lex token after the closing token.
 */
func (p *Parser) many(openKind, closeKind TokenKind, parseFn func()) {
	p.expect(openKind)
	parseFn()
	for !p.skip(closeKind) {
		parseFn()
	}
}

/**
 * Converts a name lex token into a name parse node.
 */
func (p *Parser) parseName() *Name {
	token := p.expect(TOKEN_NAME)
	return &Name{
		Value:    token.Value,
		Location: p.loc(token.Start),
	}
}

// Implements the parsing rules in the Document section.

/**
 * Document : IDefinition+
 */
func (p *Parser) parseDocument() *Document {
	start := p.token.Start
	definitions := make([]IDefinition, 4)[:0]
	definitions = append(definitions, p.parseDefinition())
	for !p.skip(TOKEN_EOF) {
		definitions = append(definitions, p.parseDefinition())
	}

	return &Document{
		Definitions: definitions,
		Location:    p.loc(start),
	}
}

/**
 * IDefinition :
 *   - OperationDefinition
 *   - FragmentDefinition
 *   - TypeDefinition
 */
func (p *Parser) parseDefinition() IDefinition {
	if p.peek(TOKEN_BRACE_L) {
		return p.parseOperationDefinition()
	}

	if p.peek(TOKEN_NAME) {
		switch p.token.Value {
		case "query", "mutation", "subscription":
			p.parseOperationDefinition()
		case "fragment":
			return p.parseFragmentDefinition()
		case "type", "interface", "union", "scalar", "enum", "input", "extend":
			return p.parseTypeDefinition()
		}
	}

	panic(p.unexpected(nil))
}

// Implements the parsing rules in the Operations section.

/**
 * OperationDefinition :
 *  - SelectionSet
 *  - OperationType Name VariableDefinitions? Directives? SelectionSet
 *
 * OperationType : one of query mutation
 */
func (p *Parser) parseOperationDefinition() *OperationDefinition {
	start := p.token.Start
	if p.peek(TOKEN_BRACE_L) {
		return &OperationDefinition{
			Operation:           OperationQuery,
			Name:                nil,
			VariableDefinitions: nil,
			Directives:          nil,
			SelectionSet:        p.parseSelectionSet(),
			Location:            p.loc(start),
		}
	}

	operationToken := p.expect(TOKEN_NAME)
	operation := operationToken.Value
	name := p.parseName()
	return &OperationDefinition{
		Operation:           OperationType(operation),
		Name:                name,
		VariableDefinitions: p.parseVariableDefinitions(),
		Directives:          p.parseDirectives(),
		SelectionSet:        p.parseSelectionSet(),
		Location:            p.loc(start),
	}
}

/**
 * VariableDefinitions : ( VariableDefinition+ )
 */
func (p *Parser) parseVariableDefinitions() []*VariableDefinition {
	if p.peek(TOKEN_PAREN_L) {
		result := make([]*VariableDefinition, 4)[:0]
		p.many(TOKEN_PAREN_L, TOKEN_PAREN_R, func() {
			result = append(result, p.parseVariableDefinition())
		})
		return result
	}

	return nil
}

/**
 * VariableDefinition : Variable : IType DefaultValue?
 */
func (p *Parser) parseVariableDefinition() *VariableDefinition {
	start := p.token.Start

	result := &VariableDefinition{}
	result.Variable = p.parseVariable()
	p.expect(TOKEN_COLON)
	result.Type = p.parseType()
	if p.skip(TOKEN_EQUALS) {
		result.DefaultValue = p.parseValueLiteral(true)
	}
	result.Location = p.loc(start)
	return result
}

/**
 * Variable : $ Name
 */
func (p *Parser) parseVariable() *Variable {
	start := p.token.Start
	p.expect(TOKEN_DOLLAR)
	return &Variable{
		Name:     p.parseName(),
		Location: p.loc(start),
	}
}

/**
 * SelectionSet : { ISelection+ }
 */
func (p *Parser) parseSelectionSet() *SelectionSet {
	start := p.token.Start
	selections := make([]ISelection, 4)[:0]
	p.many(TOKEN_BRACE_L, TOKEN_BRACE_R, func() {
		selections = append(selections, p.parseSelection())
	})
	return &SelectionSet{
		Selections: selections,
		Location:   p.loc(start),
	}
}

/**
 * ISelection :
 *   - Field
 *   - FragmentSpread
 *   - InlineFragment
 */
func (p *Parser) parseSelection() ISelection {
	if p.peek(TOKEN_SPREAD) {
		return p.parseFragment()
	}

	return p.parseField()
}

/**
 * Field : Alias? Name Arguments? Directives? SelectionSet?
 *
 * Alias : Name :
 */
func (p *Parser) parseField() *Field {
	start := p.token.Start
	nameOrAlias := p.parseName()
	var alias *Name
	var name *Name
	if p.skip(TOKEN_COLON) {
		alias = nameOrAlias
		name = p.parseName()
	} else {
		name = nameOrAlias
	}

	result := &Field{}
	result.Alias = alias
	result.Name = name
	result.Arguments = p.parseArguments()
	result.Directives = p.parseDirectives()
	if p.peek(TOKEN_BRACE_L) {
		selectionSet := p.parseSelectionSet()
		result.SelectionSet = selectionSet
	}
	result.Location = p.loc(start)
	return result
}

/**
 * Arguments : ( Argument+ )
 */
func (p *Parser) parseArguments() []*Argument {
	if p.peek(TOKEN_PAREN_L) {
		result := make([]*Argument, 4)[:0]
		p.many(TOKEN_PAREN_L, TOKEN_PAREN_R, func() {
			result = append(result, p.parseArgument())
		})
		return result
	}

	return nil
}

/**
 * Argument : Name : IValue
 */
func (p *Parser) parseArgument() *Argument {
	start := p.token.Start

	result := &Argument{}
	result.Name = p.parseName()
	p.expect(TOKEN_COLON)
	result.Value = p.parseValueLiteral(false)
	result.Location = p.loc(start)
	return result
}

// Implements the parsing rules in the Fragments section.

/**
 * Corresponds to both FragmentSpread and InlineFragment in the spec.
 *
 * FragmentSpread : ... FragmentName Directives?
 *
 * InlineFragment : ... on TypeCondition Directives? SelectionSet
 */
func (p *Parser) parseFragment() IFragment {
	start := p.token.Start
	p.expect(TOKEN_SPREAD)
	if p.token.Value == "on" {
		p.advance()
		return &InlineFragment{
			TypeCondition: p.parseNamedType(),
			Directives:    p.parseDirectives(),
			SelectionSet:  p.parseSelectionSet(),
			Location:      p.loc(start),
		}
	}

	return &FragmentSpread{
		Name:       p.parseFragmentName(),
		Directives: p.parseDirectives(),
		Location:   p.loc(start),
	}
}

/**
 * FragmentDefinition :
 *   - fragment FragmentName on TypeCondition Directives? SelectionSet
 *
 * TypeCondition : NamedType
 */
func (p *Parser) parseFragmentDefinition() *FragmentDefinition {
	start := p.token.Start

	result := &FragmentDefinition{}
	p.expectKeyword("fragment")
	result.Name = p.parseFragmentName()
	p.expectKeyword("on")
	result.TypeCondition = p.parseNamedType()
	result.Directives = p.parseDirectives()
	result.SelectionSet = p.parseSelectionSet()
	result.Location = p.loc(start)
	return result
}

/**
 * FragmentName : Name but not `on`
 */
func (p *Parser) parseFragmentName() *Name {
	if p.token.Value == "on" {
		panic(p.unexpected(nil))
	}
	return p.parseName()
}

// Implements the parsing rules in the Values section.

/**
 * IValue[Const] :
 *   - [~Const] Variable
 *   - IntValue
 *   - FloatValue
 *   - StringValue
 *   - BooleanValue
 *   - EnumValue
 *   - ListValue[?Const]
 *   - ObjectValue[?Const]
 *
 * BooleanValue : one of `true` `false`
 *
 * EnumValue : Name but not `true`, `false` or `null`
 */
func (p *Parser) parseValueLiteral(isConst bool) IValue {
	token := p.token
	switch token.Kind {
	case TOKEN_BRACKET_L:
		return p.parseList(isConst)
	case TOKEN_BRACE_L:
		return p.parseObject(isConst)
	case TOKEN_INT:
		p.advance()
		return &IntValue{
			Value:    token.Value,
			Location: p.loc(token.Start),
		}
	case TOKEN_FLOAT:
		p.advance()
		return &FloatValue{
			Value:    token.Value,
			Location: p.loc(token.Start),
		}
	case TOKEN_STRING:
		p.advance()
		return &StringValue{
			Value:    token.Value,
			Location: p.loc(token.Start),
		}
	case TOKEN_NAME:
		if token.Value == "true" || token.Value == "false" {
			p.advance()
			return &BooleanValue{
				Value:    token.Value,
				Location: p.loc(token.Start),
			}
		} else if token.Value != "null" {
			p.advance()
			return &EnumValue{
				Value:    token.Value,
				Location: p.loc(token.Start),
			}
		}
	case TOKEN_DOLLAR:
		if !isConst {
			return p.parseVariable()
		}
	}
	panic(p.unexpected(nil))
}

func (p *Parser) parseConstValue() IValue {
	return p.parseValueLiteral(true)
}

func (p *Parser) parseValueValue() IValue {
	return p.parseValueLiteral(false)
}

/**
 * ListValue[Const] :
 *   - [ ]
 *   - [ IValue[?Const]+ ]
 */
func (p *Parser) parseList(isConst bool) *ListValue {
	start := p.token.Start
	var itemFn func() IValue
	if isConst {
		itemFn = p.parseConstValue
	} else {
		itemFn = p.parseValueValue
	}

	values := make([]IValue, 4)[:0]
	p.any(TOKEN_BRACKET_L, TOKEN_BRACKET_R, func() {
		values = append(values, itemFn())
	})
	return &ListValue{
		Values:   values,
		Location: p.loc(start),
	}
}

/**
 * ObjectValue[Const] :
 *   - { }
 *   - { ObjectField[?Const]+ }
 */
func (p *Parser) parseObject(isConst bool) *ObjectValue {
	start := p.token.Start
	p.expect(TOKEN_BRACE_L)
	fieldNames := make(map[string]struct{})
	fields := make([]*ObjectField, 4)[:0]
	for !p.skip(TOKEN_BRACE_R) {
		fields = append(fields, p.parseObjectField(isConst, fieldNames))
	}
	return &ObjectValue{
		Fields:   fields,
		Location: p.loc(start),
	}
}

/**
 * ObjectField[Const] : Name : IValue[?Const]
 */
func (p *Parser) parseObjectField(isConst bool, fieldNames map[string]struct{}) *ObjectField {
	start := p.token.Start
	name := p.parseName()
	if _, ok := fieldNames[name.Value]; ok {
		panic(SyntaxError(p.source, start,
			fmt.Sprintf("Duplicate input object field %v", name.Value)))
	}
	fieldNames[name.Value] = struct{}{}

	p.expect(TOKEN_COLON)
	return &ObjectField{
		Name:     name,
		Value:    p.parseValueLiteral(isConst),
		Location: p.loc(start),
	}
}

// Implements the parsing rules in the Directives section.

/**
 * Directives : Directive+
 */
func (p *Parser) parseDirectives() []*Directive {
	directives := make([]*Directive, 4)[:0]
	for p.peek(TOKEN_AT) {
		directives = append(directives, p.parseDirective())
	}
	return directives
}

/**
 * Directive : @ Name Arguments?
 */
func (p *Parser) parseDirective() *Directive {
	start := p.token.Start
	p.expect(TOKEN_AT)
	return &Directive{
		Name:      p.parseName(),
		Arguments: p.parseArguments(),
		Location:  p.loc(start),
	}
}

// Implements the parsing rules in the Types section.

/**
 * IType :
 *   - NamedType
 *   - ListType
 *   - NonNullType
 */
func (p *Parser) parseType() IType {
	start := p.token.Start
	var typ IType
	if p.skip(TOKEN_BRACKET_L) {
		typ := p.parseType()
		p.expect(TOKEN_BRACKET_R)
		typ = &ListType{
			Type:     typ,
			Location: p.loc(start),
		}
	} else {
		typ = p.parseNamedType()
	}
	if p.skip(TOKEN_BANG) {
		return &NonNullType{
			Type:     typ.(INonNullType),
			Location: p.loc(start),
		}
	}
	return typ
}

/**
 * NamedType : Name
 */
func (p *Parser) parseNamedType() *NamedType {
	start := p.token.Start
	return &NamedType{
		Name:     p.parseName(),
		Location: p.loc(start),
	}
}

// Implements the parsing rules in the IType IDefinition section.

/**
 * TypeDefinition :
 *   - ObjectTypeDefinition
 *   - InterfaceTypeDefinition
 *   - UnionTypeDefinition
 *   - ScalarTypeDefinition
 *   - EnumTypeDefinition
 *   - InputObjectTypeDefinition
 *   - TypeExtensionDefinition
 */
func (p *Parser) parseTypeDefinition() ITypeDefinition {
	if !p.peek(TOKEN_NAME) {
		panic(p.unexpected(nil))
	}
	switch p.token.Value {
	case "type":
		return p.parseObjectTypeDefinition()
	case "interface":
		return p.parseInterfaceTypeDefinition()
	case "union":
		return p.parseUnionTypeDefinition()
	case "scalar":
		return p.parseScalarTypeDefinition()
	case "enum":
		return p.parseEnumTypeDefinition()
	case "input":
		return p.parseInputObjectTypeDefinition()
	case "extend":
		return p.parseTypeExtensionDefinition()
	default:
		panic(p.unexpected(nil))
	}
}

/**
 * ObjectTypeDefinition : type Name ImplementsInterfaces? { FieldDefinition+ }
 */
func (p *Parser) parseObjectTypeDefinition() *ObjectTypeDefinition {
	start := p.token.Start
	p.expectKeyword("type")
	name := p.parseName()
	interfaces := p.parseImplementsInterfaces()
	fields := make([]*FieldDefinition, 4)[:0]
	p.any(TOKEN_BRACE_L, TOKEN_BRACE_R, func() {
		fields = append(fields, p.parseFieldDefinition())
	})
	return &ObjectTypeDefinition{
		Name:       name,
		Interfaces: interfaces,
		Fields:     fields,
		Location:   p.loc(start),
	}
}

/**
 * ImplementsInterfaces : implements NamedType+
 */
func (p *Parser) parseImplementsInterfaces() []*NamedType {
	types := make([]*NamedType, 4)[:0]
	if p.token.Value == "implements" {
		p.advance()
		types = append(types, p.parseNamedType())
		for !p.peek(TOKEN_BRACE_L) {
			types = append(types, p.parseNamedType())
		}
	}
	return types
}

/**
 * FieldDefinition : Name ArgumentsDefinition? : IType
 */
func (p *Parser) parseFieldDefinition() *FieldDefinition {
	start := p.token.Start
	result := &FieldDefinition{}
	result.Name = p.parseName()
	result.Arguments = p.parseArgumentDefs()
	p.expect(TOKEN_COLON)
	result.Type = p.parseType()
	result.Location = p.loc(start)
	return result
}

/**
 * ArgumentsDefinition : ( InputValueDefinition+ )
 */
func (p *Parser) parseArgumentDefs() []*InputValueDefinition {
	if !p.peek(TOKEN_PAREN_L) {
		return nil
	}

	result := make([]*InputValueDefinition, 4)[:0]
	p.many(TOKEN_PAREN_L, TOKEN_PAREN_R, func() {
		result = append(result, p.parseInputValueDef())
	})
	return result
}

/**
 * InputValueDefinition : Name : IType DefaultValue?
 */
func (p *Parser) parseInputValueDef() *InputValueDefinition {
	start := p.token.Start
	result := &InputValueDefinition{}
	result.Name = p.parseName()
	p.expect(TOKEN_COLON)
	result.Type = p.parseType()
	if p.skip(TOKEN_EQUALS) {
		result.DefaultValue = p.parseConstValue()
	}
	result.Location = p.loc(start)
	return result
}

/**
 * InterfaceTypeDefinition : interface Name { FieldDefinition+ }
 */
func (p *Parser) parseInterfaceTypeDefinition() *InterfaceTypeDefinition {
	start := p.token.Start
	p.expectKeyword("interface")
	name := p.parseName()
	fields := make([]*FieldDefinition, 4)[:0]
	p.any(TOKEN_BRACE_L, TOKEN_BRACE_R, func() {
		fields = append(fields, p.parseFieldDefinition())
	})
	return &InterfaceTypeDefinition{
		Name:     name,
		Fields:   fields,
		Location: p.loc(start),
	}
}

/**
 * UnionTypeDefinition : union Name = UnionMembers
 */
func (p *Parser) parseUnionTypeDefinition() *UnionTypeDefinition {
	start := p.token.Start
	p.expectKeyword("union")
	name := p.parseName()
	types := p.parseUnionMembers()
	return &UnionTypeDefinition{
		Name:     name,
		Types:    types,
		Location: p.loc(start),
	}
}

/**
 * UnionMembers :
 *   - NamedType
 *   - UnionMembers | NamedType
 */
func (p *Parser) parseUnionMembers() []*NamedType {
	result := make([]*NamedType, 4)[:0]
	result = append(result, p.parseNamedType())
	for p.skip(TOKEN_PIPE) {
		result = append(result, p.parseNamedType())
	}
	return result
}

/**
 * ScalarTypeDefinition : scalar Name
 */
func (p *Parser) parseScalarTypeDefinition() *ScalarTypeDefinition {
	start := p.token.Start
	p.expectKeyword("scalar")
	name := p.parseName()
	return &ScalarTypeDefinition{
		Name:     name,
		Location: p.loc(start),
	}
}

/**
 * EnumTypeDefinition : enum Name { EnumValueDefinition+ }
 */
func (p *Parser) parseEnumTypeDefinition() *EnumTypeDefinition {
	start := p.token.Start
	p.expectKeyword("enum")
	name := p.parseName()
	values := make([]*EnumValueDefinition, 4)[:0]
	p.many(TOKEN_BRACE_L, TOKEN_BRACE_R, func() {
		values = append(values, p.parseEnumValueDefinition())
	})
	return &EnumTypeDefinition{
		Name:     name,
		Values:   values,
		Location: p.loc(start),
	}
}

/**
 * EnumValueDefinition : EnumValue
 *
 * EnumValue : Name
 */
func (p *Parser) parseEnumValueDefinition() *EnumValueDefinition {
	start := p.token.Start
	name := p.parseName()
	return &EnumValueDefinition{
		Name:     name,
		Location: p.loc(start),
	}
}

/**
 * InputObjectTypeDefinition : input Name { InputValueDefinition+ }
 */
func (p *Parser) parseInputObjectTypeDefinition() *InputObjectTypeDefinition {
	start := p.token.Start
	p.expectKeyword("input")
	name := p.parseName()
	fields := make([]*InputValueDefinition, 4)[:0]
	p.any(TOKEN_BRACE_L, TOKEN_BRACE_R, func() {
		fields = append(fields, p.parseInputValueDef())
	})
	return &InputObjectTypeDefinition{
		Name:     name,
		Fields:   fields,
		Location: p.loc(start),
	}
}

/**
 * TypeExtensionDefinition : extend ObjectTypeDefinition
 */
func (p *Parser) parseTypeExtensionDefinition() *TypeExtensionDefinition {
	start := p.token.Start
	p.expectKeyword("extend")
	definition := p.parseObjectTypeDefinition()
	return &TypeExtensionDefinition{
		Definition: definition,
		Location:   p.loc(start),
	}
}
