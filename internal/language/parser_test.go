package language

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParse_AcceptsOptionToNotIncludeSource(T *testing.T) {
	deepEqual(T, Parse(NewSource(`{ field }`, ""), ParseOptions{
		NoSource: false,
	}), &Document{
		Location: &Location{
			Start: 0,
			End:   9,
		},
		Definitions: []Definition{
			OperationDefinition{
				Location: &Location{
					Start: 0,
					End:   9,
				},
				Operation:           OperationQuery,
				Name:                nil,
				VariableDefinitions: nil,
				Directives:          []Directive{},
				SelectionSet: SelectionSet{
					Location: &Location{
						Start: 0,
						End:   9,
					},
					Selections: []Selection{
						Field{
							Location: &Location{
								Start: 2,
								End:   7,
							},
							Alias: nil,
							Name: Name{
								Location: &Location{
									Start: 2,
									End:   7,
								},
								Value: "field",
							},
							Arguments:    []Argument{},
							Directives:   []Directive{},
							SelectionSet: nil,
						},
					},
				},
			},
		},
	})
}

func TestParse_ParseProvidesUsefulErrors(T *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(QLError); ok {
				deepEqual(T, err.Error(), `Syntax Error GraphQL (1:2) Expected Name, found EOF

  1: {
      ^
  `)
				deepEqual(T, err.Positions, []int{1})

				deepEqual(T, err.Locations, []SourceLocation{
					SourceLocation{
						Line:   1,
						Column: 2,
					},
				})
			}

		}
	}()
	Parse(NewSource(`{`, ""), ParseOptions{})

	expectPanic(T, func() {
		Parse(NewSource(`{ ...MissingOn }
fragment MissingOn Type
`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (2:20) Expected "on", found Name "Type"`)

	expectPanic(T, func() {
		Parse(NewSource(`{ field: {} }`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (1:10) Expected Name, found {`)

	expectPanic(T, func() {
		Parse(NewSource(`notanoperation Foo { field }`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (1:1) Unexpected Name "notanoperation"`)

	expectPanic(T, func() {
		Parse(NewSource(`...`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (1:1) Unexpected ...`)
}

func TestParse_ParseProvidesUsefulErrorWhenUsingSource(T *testing.T) {
	expectPanic(T, func() {
		Parse(NewSource(`MyQuery.graphql`, "query"), ParseOptions{})
	}, `Syntax Error MyQuery.graphql (1:6) Expected Name, found EOF`)
}

func TestParse_ParsesVariableInlineValues(T *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			T.Error("Expect not panic", err)
		}
	}()
	Parse(NewSource(`{ field(complex: { a: { b: [ $var ] } }) }`, ""), ParseOptions{})
}

func TestParse_ParseConstantDefaultValues(T *testing.T) {
	expectPanic(T, func() {
		Parse(NewSource(`query Foo($x: Complex = { a: { b: [ $var ] } }) { field }`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (1:37) Unexpected $`)
}

func TestParse_DuplicateKeysInInputObjectIsSyntaxError(T *testing.T) {
	expectPanic(T, func() {
		Parse(NewSource(`{ field(arg: { a: 1, a: 2 }) }`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (1:22) Duplicate input object field a.`)
}

func TestParse_DoesNotAcceptFragmentsNamedOn(T *testing.T) {
	expectPanic(T, func() {
		Parse(NewSource(`fragment on on on { on }`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (1:10) Unexpected Name "on"`)
}

func TestParse_DoesNotAcceptFragmentsSpreadOfOn(T *testing.T) {
	expectPanic(T, func() {
		Parse(NewSource(`{ ...on }`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (1:9) Expected Name, found }`)
}

func TestParse_DoesNotAllowNullAsValue(T *testing.T) {
	expectPanic(T, func() {
		Parse(NewSource(`{ fieldWithNullableStringInput(input: null) }`, ""), ParseOptions{})
	}, `Syntax Error GraphQL (1:39) Unexpected Name "null"`)
}

func TestParse_ParsesKitchenSink(T *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			T.Error("Expect not panic", err)
		}
	}()

	kitchenSink, err := ioutil.ReadFile("kitchen-sink.graphql")

	if err != nil {
		panic(err)
	}

	Parse(NewSource(string(kitchenSink), ""), ParseOptions{})
}

func TestParse_AllowsNonKeywordsAnywhereANameIsAllowed(T *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			T.Error("Expect not panic", err)
		}
	}()

	nonKeywords := []string{
		"on",
		"fragment",
		"query",
		"mutation",
		"true",
		"false",
	}

	for i := range nonKeywords {
		keyword := nonKeywords[i]
		fragmentName := keyword

		if keyword == "on" {
			fragmentName = "a"
		}

		Parse(NewSource(fmt.Sprintf(`query %v {
  ... %v
  ... on %v { field }
}
fragment %v on Type {
  %v(%v: $%v) @%v(%v: %v)
}`, keyword, fragmentName, keyword, fragmentName, keyword, keyword, keyword, keyword, keyword, keyword), ""), ParseOptions{})
	}
}

func TestParse_ParsesExperimentalSubscriptionFeature(T *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			T.Error("Expect not panic", err)
		}
	}()

	Parse(NewSource(`
      subscription Foo {
        subscriptionField
      }
    `, ""), ParseOptions{})

}

func TestParse_ParseCreatesAst(T *testing.T) {
	source := NewSource(`{
  node(id: 4) {
    id,
    name
  }
}
`, "")
	deepEqual(T, Parse(source, ParseOptions{}), &Document{
		Location: &Location{
			Start:  0,
			End:    41,
			Source: &source,
		},
		Definitions: []Definition{
			OperationDefinition{
				Location: &Location{
					Start:  0,
					End:    40,
					Source: &source,
				},

				Operation:           OperationQuery,
				Name:                nil,
				VariableDefinitions: nil,
				Directives:          []Directive{},
				SelectionSet: SelectionSet{
					Location: &Location{
						Start:  0,
						End:    40,
						Source: &source,
					},
					Selections: []Selection{
						Field{
							Location: &Location{
								Start:  4,
								End:    38,
								Source: &source,
							},
							Alias: nil,
							Name: Name{
								Location: &Location{
									Start:  4,
									End:    8,
									Source: &source,
								},
								Value: "node",
							},
							Arguments: []Argument{
								Argument{
									Name: Name{
										Location: &Location{
											Start:  9,
											End:    11,
											Source: &source,
										},
										Value: "id",
									},
									Value: IntValue{
										Location: &Location{
											Start:  13,
											End:    14,
											Source: &source,
										},
										Value: "4",
									},
									Location: &Location{
										Start:  9,
										End:    14,
										Source: &source,
									},
								},
							},
							Directives: []Directive{},
							SelectionSet: &SelectionSet{
								Location: &Location{
									Start:  16,
									End:    38,
									Source: &source,
								},
								Selections: []Selection{
									Field{
										Location: &Location{
											Start:  22,
											End:    24,
											Source: &source,
										},
										Alias: nil,
										Name: Name{
											Location: &Location{
												Start:  22,
												End:    24,
												Source: &source,
											},
											Value: "id",
										},
										Arguments:    []Argument{},
										Directives:   []Directive{},
										SelectionSet: nil,
									},
									Field{
										Location: &Location{
											Start:  30,
											End:    34,
											Source: &source,
										},
										Alias: nil,
										Name: Name{
											Location: &Location{
												Start:  30,
												End:    34,
												Source: &source,
											},
											Value: "name",
										},
										Arguments:    []Argument{},
										Directives:   []Directive{},
										SelectionSet: nil,
									},
								},
							},
						},
					},
				},
			},
		},
	})
}
