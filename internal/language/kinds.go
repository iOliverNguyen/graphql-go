package language

type NodeKind string

const (
	// Name

	NAME = NodeKind("Name")

	// Document

	DOCUMENT             = NodeKind("Document")
	OPERATION_DEFINITION = NodeKind("OperationDefinition")
	VARIABLE_DEFINITION  = NodeKind("VariableDefinition")
	VARIABLE             = NodeKind("Variable")
	SELECTION_SET        = NodeKind("SelectionSet")
	FIELD                = NodeKind("Field")
	ARGUMENT             = NodeKind("Argument")

	// Fragments

	FRAGMENT_SPREAD     = NodeKind("FragmentSpread")
	INLINE_FRAGMENT     = NodeKind("InlineFragment")
	FRAGMENT_DEFINITION = NodeKind("FragmentDefinition")

	// Values

	INT_VALUE     = NodeKind("IntValue")
	FLOAT_VALUE   = NodeKind("FloatValue")
	STRING_VALUE  = NodeKind("StringValue")
	BOOLEAN_VALUE = NodeKind("BooleanValue")
	ENUM_VALUE    = NodeKind("EnumValue")
	LIST_VALUE    = NodeKind("ListValue")
	OBJECT_VALUE  = NodeKind("ObjectValue")
	OBJECT_FIELD  = NodeKind("ObjectField")

	// Directives

	DIRECTIVE = NodeKind("Directive")

	// Types

	NAMED_TYPE    = NodeKind("NamedType")
	LIST_TYPE     = NodeKind("ListType")
	NON_NULL_TYPE = NodeKind("NonNullType")

	// IType Definitions

	OBJECT_TYPE_DEFINITION       = NodeKind("ObjectTypeDefinition")
	FIELD_DEFINITION             = NodeKind("FieldDefinition")
	INPUT_VALUE_DEFINITION       = NodeKind("InputValueDefinition")
	INTERFACE_TYPE_DEFINITION    = NodeKind("InterfaceTypeDefinition")
	UNION_TYPE_DEFINITION        = NodeKind("UnionTypeDefinition")
	SCALAR_TYPE_DEFINITION       = NodeKind("ScalarTypeDefinition")
	ENUM_TYPE_DEFINITION         = NodeKind("EnumTypeDefinition")
	ENUM_VALUE_DEFINITION        = NodeKind("EnumValueDefinition")
	INPUT_OBJECT_TYPE_DEFINITION = NodeKind("InputObjectTypeDefinition")
	TYPE_EXTENSION_DEFINITION    = NodeKind("TypeExtensionDefinition")
)
