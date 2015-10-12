package language

import (
	"bytes"
	"encoding/json"
	"reflect"
)

const INDENT_CHAR = "    "

func Print(ast INode) string {
	visitor := &printASTVisitor{}
	visitor.visit(ast)
	return visitor.buf.String()
}

type printASTVisitor struct {
	buf bytes.Buffer

	indentLevel string
	wrapStart   []string
	wrapEnd     []string
	pendingEnd  string
}

func (p *printASTVisitor) visit(node INode) {
	if node == nil || reflect.ValueOf(node).IsNil() {
		return
	}

	switch node := node.(type) {
	case *Name:
		p.write(node.Value)

	case *Variable:
		p.write("$")
		p.visit(node.Name)

	// Document

	case *Document:
		for i, def := range node.Definitions {
			p.writeIp(i, "\n\n")
			p.visit(def)
		}
		p.write("\n")

	case *OperationDefinition:
		if node.Name == nil || node.Name.Value == "" {
			p.visit(node.SelectionSet)
			break
		}
		p.write(string(node.Operation))
		p.write(" ")
		p.visit(node.Name)
		p.wrapOpen("(", ")")
		for i, def := range node.VariableDefinitions {
			p.writeIp(i, ", ")
			p.visit(def)
		}
		p.wrapClose()
		for _, directive := range node.Directives {
			p.write(" ")
			p.visit(directive)
		}
		p.write(" ")
		p.visit(node.SelectionSet)

	case *VariableDefinition:
		p.visit(node.Variable)
		p.write(": ")
		p.visit(node.Type)

		p.wrapOpen(" = ", "")
		p.visit(node.DefaultValue)
		p.wrapClose()

	case *SelectionSet:
		p.blockOpen()
		for i, selection := range node.Selections {
			p.blockLine(i)
			p.visit(selection)
		}
		p.blockClose()

	case *Field:
		p.visit(node.Alias)
		p.writeIf(node.Alias != nil, ": ")
		p.visit(node.Name)

		p.writeIp(len(node.Arguments), "(")
		for i, arg := range node.Arguments {
			p.writeIp(i, ",")
			p.visit(arg)
		}
		p.writeIp(len(node.Arguments), ")")

		for _, directive := range node.Directives {
			p.write(" ")
			p.visit(directive)
		}
		p.write(" ")
		p.visit(node.SelectionSet)

	case *Argument:
		p.visit(node.Name)
		p.write(": ")
		p.visit(node.Value)

	case *FragmentSpread:
		p.write("...")
		p.visit(node.Name)
		for _, directive := range node.Directives {
			p.write(" ")
			p.visit(directive)
		}

	// Fragments

	case *InlineFragment:
		p.write("... on ")
		p.visit(node.TypeCondition)
		for _, directive := range node.Directives {
			p.write(" ")
			p.visit(directive)
		}
		p.write(" ")
		p.visit(node.SelectionSet)

	case *FragmentDefinition:
		p.write("fragment ")
		p.visit(node.Name)
		p.write(" on ")
		p.visit(node.TypeCondition)
		for _, directive := range node.Directives {
			p.write(" ")
			p.visit(directive)
		}
		p.write(" ")
		p.visit(node.SelectionSet)

	// IValue

	case *IntValue:
		p.write(node.Value)

	case *FloatValue:
		p.write(node.Value)

	case *StringValue:
		data, err := json.Marshal(node.Value)
		if err != nil {
			panic(err)
		}
		p.write(string(data))

	case *BooleanValue:
		p.write(node.Value)

	case *EnumValue:
		p.write(node.Value)

	case *ListValue:
		p.write("[")
		for i, value := range node.Values {
			p.writeIp(i, ", ")
			p.visit(value)
		}
		p.write("]")

	case *ObjectValue:
		p.write("{")
		for i, field := range node.Fields {
			p.writeIp(i, ", ")
			p.visit(field)
		}
		p.write("}")

	case *ObjectField:
		p.visit(node.Name)
		p.write(": ")
		p.visit(node.Value)

	// Directive

	case *Directive:
		p.write("@")
		p.visit(node.Name)
		p.writeIp(len(node.Arguments), "(")
		for i, arg := range node.Arguments {
			p.writeIp(i, ", ")
			p.visit(arg)
		}
		p.write(")")

	// IType

	case *NamedType:
		p.visit(node.Name)

	case *ListType:
		p.write("[")
		p.visit(node.Type)
		p.write("]")

	case *NonNullType:
		p.visit(node.Type)
		p.write("!")

	// IType Definitions

	case *ObjectTypeDefinition:
		p.write("type ")
		p.visit(node.Name)
		p.write(" ")
		p.writeIp(len(node.Interfaces), "implements ")
		for i, iface := range node.Interfaces {
			p.writeIp(i, ", ")
			p.visit(iface)
		}
		p.write(" ")
		p.blockOpen()
		for i, field := range node.Fields {
			p.blockLine(i)
			p.visit(field)
		}
		p.blockClose()

	case *FieldDefinition:
		p.visit(node.Name)
		p.writeIp(len(node.Arguments), "(")
		for i, arg := range node.Arguments {
			p.writeIp(i, ", ")
			p.visit(arg)
		}
		p.writeIp(len(node.Arguments), ")")
		p.write(": ")
		p.visit(node.Type)

	case *InputValueDefinition:
		p.visit(node.Name)
		p.write(": ")
		p.visit(node.Type)

		p.wrapOpen(" = ", "")
		p.visit(node.DefaultValue)
		p.wrapClose()

	case *InterfaceTypeDefinition:
		p.write("interface ")
		p.visit(node.Name)
		p.write(" ")
		p.blockOpen()
		for i, field := range node.Fields {
			p.blockLine(i)
			p.visit(field)
		}
		p.blockClose()

	case *UnionTypeDefinition:
		p.write("union ")
		p.visit(node.Name)
		p.write(" = ")
		for i, typ := range node.Types {
			p.writeIp(i, " | ")
			p.visit(typ)
		}

	case *ScalarTypeDefinition:
		p.write("scalar ")
		p.visit(node.Name)

	case *EnumTypeDefinition:
		p.write("enum ")
		p.visit(node.Name)
		p.write(" ")
		p.blockOpen()
		for i, value := range node.Values {
			p.blockLine(i)
			p.visit(value)
		}
		p.blockClose()

	case *EnumValueDefinition:
		p.visit(node.Name)

	case *InputObjectTypeDefinition:
		p.write("input ")
		p.visit(node.Name)
		p.write(" ")
		p.blockOpen()
		for i, field := range node.Fields {
			p.blockLine(i)
			p.visit(field)
		}
		p.blockClose()

	case *TypeExtensionDefinition:
		p.write("extend ")
		p.visit(node.Definition)
	}
}

func (p *printASTVisitor) write(s string) {
	if s != "" && len(p.wrapStart) > 0 {
		if len(p.wrapStart) != len(p.wrapEnd) {
			panic("unexpected different length")
		}
		for _, start := range p.wrapStart {
			_, err := p.buf.WriteString(start)
			if err != nil {
				panic(err)
			}
		}
		for i := len(p.wrapEnd) - 1; i >= 0; i-- {
			p.pendingEnd += p.wrapEnd[i]
		}
		p.wrapStart = p.wrapStart[:0]
	}
	_, err := p.buf.WriteString(s)
	if err != nil {
		panic(err)
	}
}

func (p *printASTVisitor) writeIf(cond bool, s string) {
	if cond {
		p.write(s)
	}
}

func (p *printASTVisitor) writeIp(value int, s string) {
	if value > 0 {
		p.write(s)
	}
}

func (p *printASTVisitor) newline() {
	p.write("\n")
	p.write(p.indentLevel)
}

func (p *printASTVisitor) indent() {
	p.indentLevel += INDENT_CHAR
	p.newline()
}

func (p *printASTVisitor) outdent() {
	l := len(p.indentLevel)
	if l == 0 {
		panic("unexpected outdent")
	}
	p.indentLevel = p.indentLevel[:l-len(INDENT_CHAR)]
	p.newline()
}

func (p *printASTVisitor) blockOpen() {
	p.write("{")
	p.indent()
}

func (p *printASTVisitor) blockClose() {
	p.outdent()
	p.write("}")
}

func (p *printASTVisitor) blockLine(i int) {
	if i > 0 {
		p.newline()
	}
}

func (p *printASTVisitor) wrapOpen(start, end string) {
	p.wrapStart = append(p.wrapStart, start)
	p.wrapEnd = append(p.wrapEnd, end)
}

func (p *printASTVisitor) wrapClose() {
	l := len(p.wrapEnd)
	if l > 0 {
		if len(p.wrapStart) != l {
			panic("unexpected different length")
		}
		p.wrapEnd = p.wrapEnd[:l-1]
		p.wrapStart = p.wrapStart[:l-1]
	}
}
