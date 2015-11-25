package language

import "testing"

func lexOne(str string) Token {
	source := NewSource(str, "")
	lexer := newLexer(source)
	return lexer.nextToken()
}

func TestLex_DisallowsUncommonControlCharacters(T *testing.T) {
	expectPanic(T, func() {
		lexOne(`\\u0007`)
	}, `Syntax Error GraphQL (1:1) Invalid character "\\u0007"`)
}

func TestLex_AcceptsBOMHeader(T *testing.T) {
	deepEqual(T, lexOne(`\uFEFF foo`), newToken(TOKEN_NAME, 2, 5, "foo"))
}

func TestLex_SkipsWhitespace(T *testing.T) {
	deepEqual(T, lexOne(`

    foo

`), newToken(TOKEN_NAME, 6, 9, "foo"))

	deepEqual(T, lexOne(`
    #comment
    foo#comment
`), newToken(TOKEN_NAME, 18, 21, "foo"))

	deepEqual(T, lexOne(`,,,foo,,,`), newToken(TOKEN_NAME, 3, 6, "foo"))
}

func TestLex_ErrorsRespectWhitespace(T *testing.T) {
	expectPanic(T, func() {
		lexOne(`

    ?


`)
	},
		"Syntax Error GraphQL (3:5) Unexpected character \"?\".\n"+
			"\n"+
			"2: \n"+
			"3:     ?\n"+
			"       ^\n"+
			"4: \n")
}

func TestLex_Strings(T *testing.T) {
	deepEqual(T, lexOne(`"simple"`),
		newToken(TOKEN_STRING, 0, 8, "simple"))

	deepEqual(T, lexOne(`" white space "`),
		newToken(TOKEN_STRING, 0, 15, " white space "))

	deepEqual(T, lexOne(`"quote \""`),
		newToken(TOKEN_STRING, 0, 10, `quote "`))

	deepEqual(T, lexOne(`"escaped \n\r\b\t\f"`),
		newToken(TOKEN_STRING, 0, 20, "escaped \n\r\b\t\f"))

	deepEqual(T, lexOne(`"slashes \\ \/"`),
		newToken(TOKEN_STRING, 0, 15, `slashes \ /`))

	deepEqual(T, lexOne(`"unicode \u1234\u5678\u90AB\uCDEF"`),
		newToken(TOKEN_STRING, 0, 34, "unicode \u1234\u5678\u90AB\uCDEF"))
}

func TestLex_ReportStringErrors(T *testing.T) {
	expectPanic(T, func() {
		lexOne(`"no end quote`)
	}, "Syntax Error GraphQL (1:14) Unterminated string")

	expectPanic(T, func() {
		lexOne(`"multi\nline"`)
	}, "Syntax Error GraphQL (1:7) Unterminated string")

	expectPanic(T, func() {
		lexOne(`"multi\rline"`)
	}, "Syntax Error GraphQL (1:7) Unterminated string")

	expectPanic(T, func() {
		lexOne(`"multi\u2028line"`)
	}, "Syntax Error GraphQL (1:7) Unterminated string")

	expectPanic(T, func() {
		lexOne(`"multi\u2029line"`)
	}, "Syntax Error GraphQL (1:7) Unterminated string")

	expectPanic(T, func() {
		lexOne(`"bad \\z esc"`)
	}, "Syntax Error GraphQL (1:7) Bad character escape sequence")

	expectPanic(T, func() {
		lexOne(`"bad \\x esc"`)
	}, "Syntax Error GraphQL (1:7) Bad character escape sequence")

	expectPanic(T, func() {
		lexOne(`"bad \\u1 esc"`)
	}, "Syntax Error GraphQL (1:7) Bad character escape sequence")

	expectPanic(T, func() {
		lexOne(`"bad \\u0XX1 esc"`)
	}, "Syntax Error GraphQL (1:7) Bad character escape sequence")

	expectPanic(T, func() {
		lexOne(`"bad \\uXXXX esc"`)
	}, "Syntax Error GraphQL (1:7) Bad character escape sequence")

	expectPanic(T, func() {
		lexOne(`"bad \\uFXXX esc"`)
	}, "Syntax Error GraphQL (1:7) Bad character escape sequence")

	expectPanic(T, func() {
		lexOne(`"bad \\uXXXF esc"`)
	}, "Syntax Error GraphQL (1:7) Bad character escape sequence")
}

func TestLex_Numbers(T *testing.T) {
	deepEqual(T, lexOne(`4`),
		newToken(TOKEN_INT, 0, 1, `4`))

	deepEqual(T, lexOne(`4.123`),
		newToken(TOKEN_FLOAT, 0, 5, `4.123`))

	deepEqual(T, lexOne(`-4`),
		newToken(TOKEN_INT, 0, 2, `-4`))

	deepEqual(T, lexOne(`9`),
		newToken(TOKEN_INT, 0, 1, `9`))

	deepEqual(T, lexOne(`0`),
		newToken(TOKEN_INT, 0, 1, `0`))

	deepEqual(T, lexOne(`-4.123`),
		newToken(TOKEN_FLOAT, 0, 6, `-4.123`))

	deepEqual(T, lexOne(`0.123`),
		newToken(TOKEN_FLOAT, 0, 5, `0.123`))

	deepEqual(T, lexOne(`123e4`),
		newToken(TOKEN_FLOAT, 0, 5, `123e4`))

	deepEqual(T, lexOne(`123E4`),
		newToken(TOKEN_FLOAT, 0, 5, `123E4`))

	deepEqual(T, lexOne(`123e-4`),
		newToken(TOKEN_FLOAT, 0, 6, `123e-4`))

	deepEqual(T, lexOne(`123e+4`),
		newToken(TOKEN_FLOAT, 0, 6, `123e+4`))

	deepEqual(T, lexOne(`-1.123e4`),
		newToken(TOKEN_FLOAT, 0, 8, `-1.123e4`))

	deepEqual(T, lexOne(`-1.123E4`),
		newToken(TOKEN_FLOAT, 0, 8, `-1.123E4`))

	deepEqual(T, lexOne(`-1.123e-4`),
		newToken(TOKEN_FLOAT, 0, 9, `-1.123e-4`))

	deepEqual(T, lexOne(`-1.123e+4`),
		newToken(TOKEN_FLOAT, 0, 9, `-1.123e+4`))

	deepEqual(T, lexOne(`-1.123e4567`),
		newToken(TOKEN_FLOAT, 0, 11, `-1.123e4567`))
}

func TestLex_ReportNumberErrors(T *testing.T) {
	expectPanic(T, func() {
		lexOne(`00`)
	}, `Syntax Error GraphQL (1:2) Invalid number, unexpected digit after 0: "0".`)

	expectPanic(T, func() {
		lexOne(`+1`)
	}, `Syntax Error GraphQL (1:1) Unexpected character "+"`)

	expectPanic(T, func() {
		lexOne(`1.`)
	}, `Syntax Error GraphQL (1:3) Invalid number, expected digit but got: EOF.`)

	expectPanic(T, func() {
		lexOne(`.123`)
	}, `Syntax Error GraphQL (1:1) Unexpected character "."`)

	expectPanic(T, func() {
		lexOne(`1.A`)
	}, `Syntax Error GraphQL (1:3) Invalid number, expected digit but got: "A".`)

	expectPanic(T, func() {
		lexOne(`-A`)
	}, `Syntax Error GraphQL (1:2) Invalid number, expected digit but got: "A".`)

	expectPanic(T, func() {
		lexOne(`-A`)
	}, `Syntax Error GraphQL (1:2) Invalid number, expected digit but got: "A".`)

	expectPanic(T, func() {
		lexOne(`1.0e`)
	}, `Syntax Error GraphQL (1:5) Invalid number, expected digit but got: EOF.`)

	expectPanic(T, func() {
		lexOne(`1.0eA`)
	}, `Syntax Error GraphQL (1:5) Invalid number, expected digit but got: "A".`)
}

func TestLex_Punctuation(T *testing.T) {
	deepEqual(T, lexOne(`!`),
		newToken(TOKEN_BANG, 0, 1, ""))

	deepEqual(T, lexOne(`$`),
		newToken(TOKEN_DOLLAR, 0, 1, ""))

	deepEqual(T, lexOne(`(`),
		newToken(TOKEN_PAREN_L, 0, 1, ""))

	deepEqual(T, lexOne(`)`),
		newToken(TOKEN_PAREN_R, 0, 1, ""))

	deepEqual(T, lexOne(`...`),
		newToken(TOKEN_SPREAD, 0, 3, ""))

	deepEqual(T, lexOne(`:`),
		newToken(TOKEN_COLON, 0, 1, ""))

	deepEqual(T, lexOne(`=`),
		newToken(TOKEN_EQUALS, 0, 1, ""))

	deepEqual(T, lexOne(`@`),
		newToken(TOKEN_AT, 0, 1, ""))

	deepEqual(T, lexOne(`[`),
		newToken(TOKEN_BRACKET_L, 0, 1, ""))

	deepEqual(T, lexOne(`]`),
		newToken(TOKEN_BRACKET_R, 0, 1, ""))

	deepEqual(T, lexOne(`{`),
		newToken(TOKEN_BRACE_L, 0, 1, ""))

	deepEqual(T, lexOne(`|`),
		newToken(TOKEN_PIPE, 0, 1, ""))

	deepEqual(T, lexOne(`}`),
		newToken(TOKEN_BRACE_R, 0, 1, ""))

}

func TestLex_ReportUnknownCharacterError(T *testing.T) {
	expectPanic(T, func() {
		lexOne(`..`)
	}, `Syntax Error GraphQL (1:1) Unexpected character "."`)

	expectPanic(T, func() {
		lexOne(`?`)
	}, `Syntax Error GraphQL (1:1) Unexpected character "?"`)

	expectPanic(T, func() {
		lexOne(`\u203B`)
	}, `Syntax Error GraphQL (1:1) Unexpected character "\u203B"`)
}

func TestLex_ReportDashesInNames(T *testing.T) {
	q := `a-b`
	lexer := newLexer(NewSource(q, ""))
	firstToken := lexer.nextToken()

	deepEqual(T, firstToken,
		newToken(TOKEN_NAME, 0, 1, `a`))

	expectPanic(T, func() {
		lexer.nextToken()
	}, `Syntax Error GraphQL (1:3) Invalid number, expected digit but got: "b".`)
}
