package language

import (
	"fmt"
	"unicode/utf8"
)

type Token struct {
	Kind  TokenKind
	Start int
	End   int
	Value string
}

type TokenKind int

const (
	_                   = iota
	TOKEN_EOF TokenKind = iota
	TOKEN_BANG
	TOKEN_DOLLAR
	TOKEN_PAREN_L
	TOKEN_PAREN_R
	TOKEN_SPREAD
	TOKEN_COLON
	TOKEN_EQUALS
	TOKEN_AT
	TOKEN_BRACKET_L
	TOKEN_BRACKET_R
	TOKEN_BRACE_L
	TOKEN_PIPE
	TOKEN_BRACE_R
	TOKEN_NAME
	TOKEN_VARIABLE
	TOKEN_INT
	TOKEN_FLOAT
	TOKEN_STRING
)

var tokenDescription = map[TokenKind]string{
	TOKEN_EOF:       "EOF",
	TOKEN_BANG:      "!",
	TOKEN_DOLLAR:    "$",
	TOKEN_PAREN_L:   "(",
	TOKEN_PAREN_R:   ")",
	TOKEN_SPREAD:    "...",
	TOKEN_COLON:     ":",
	TOKEN_EQUALS:    "=",
	TOKEN_AT:        "@",
	TOKEN_BRACKET_L: "[",
	TOKEN_BRACKET_R: "]",
	TOKEN_BRACE_L:   "{",
	TOKEN_PIPE:      "|",
	TOKEN_BRACE_R:   "}",
	TOKEN_NAME:      "Name",
	TOKEN_VARIABLE:  "Variable",
	TOKEN_INT:       "Int",
	TOKEN_FLOAT:     "Float",
	TOKEN_STRING:    "String",
}

const EOF = -1

func (t TokenKind) String() string {
	return tokenDescription[t]
}

func (t Token) String() string {
	if t.Value != "" {
		return tokenDescription[t.Kind] + " " + t.Value
	}

	return tokenDescription[t.Kind]
}

func newToken(kind TokenKind, start, end int, value string) Token {
	return Token{
		Kind:  kind,
		Start: start,
		End:   end,
		Value: value,
	}
}

type Lexer struct {
	source       Source
	body         string // source
	char         rune   // current character
	position     int    // position of current character
	prevPosition int
	nextPosition int // position after current character
}

/**
 * Helper function for constructing the Token object.
 */
func newLexer(source Source) *Lexer {
	l := &Lexer{
		source:       source,
		body:         source.Body,
		position:     -1,
		nextPosition: 0,
		char:         0,
	}
	l.next()
	return l
}

func (l *Lexer) String() string {
	return fmt.Sprintf("Lexer %v-%v : %c\n", l.position, l.nextPosition, l.char)
}

func (l *Lexer) nextToken() Token {
	token := l.readToken(l.prevPosition)
	l.prevPosition = token.End
	return token
}

func (l *Lexer) nextTokenFromPosition(resetPosition int) Token {
	token := l.readToken(resetPosition)
	l.prevPosition = token.End
	return token
}

/**
 * Gets the next token from the source starting at the given position.
 *
 * This skips over whitespace and comments until it finds the next lexable
 * token, then lexes punctuators immediately or calls the appropriate helper
 * function for more complicated tokens.
 */
func (l *Lexer) readToken(fromPosition int) (result Token) {
	l.resetPosition(fromPosition)
	l.skipWhitespace()

	ch := l.char
	switch ch {
	case EOF:
		return newToken(TOKEN_EOF, l.position, l.position, "")
	case '!':
		return newToken(TOKEN_BANG, l.position, l.position+1, "")
	case '$':
		return newToken(TOKEN_DOLLAR, l.position, l.position+1, "")
	case '(':
		return newToken(TOKEN_PAREN_L, l.position, l.position+1, "")
	case ')':
		return newToken(TOKEN_PAREN_R, l.position, l.position+1, "")
	case '.':
		// TODO(qv)
		return
	case ':':
		return newToken(TOKEN_COLON, l.position, l.position+1, "")
	case '=':
		return newToken(TOKEN_EQUALS, l.position, l.position+1, "")
	case '@':
		return newToken(TOKEN_AT, l.position, l.position+1, "")
	case '[':
		return newToken(TOKEN_BRACKET_L, l.position, l.position+1, "")
	case ']':
		return newToken(TOKEN_BRACKET_R, l.position, l.position+1, "")
	case '{':
		return newToken(TOKEN_BRACE_L, l.position, l.position+1, "")
	case '|':
		return newToken(TOKEN_PIPE, l.position, l.position+1, "")
	case '}':
		return newToken(TOKEN_BRACE_R, l.position, l.position+1, "")
	case '"':
		return l.readString()
	}

	switch {
	case ch == '_', ch >= 'A' && ch <= 'Z', ch >= 'a' && ch <= 'z':
		return l.readName()
	case ch == '-', ch >= '0' && ch <= '9':
		return l.readNumber()
	}

	panic(SyntaxError(l.source, l.position,
		fmt.Sprintf(`Unexpected character "%c".`, ch)))
}

func (l *Lexer) next() rune {
	if l.nextPosition >= len(l.body) {
		l.char = EOF
		return EOF
	}

	ch, size := utf8.DecodeRuneInString(l.body[l.nextPosition:])
	if ch == utf8.RuneError {
		panic("illegal character")
	}

	l.char = ch
	l.position = l.nextPosition
	l.nextPosition += size
	return ch
}

func (l *Lexer) resetPosition(position int) {
	l.nextPosition = position
	l.next()
}

/**
 * Reads from body starting at startPosition until it finds a non-whitespace
 * or commented character, then returns the position of that character for
 * lexing.
 */
func (l *Lexer) skipWhitespace() {
	for {
		switch l.char {
		case ' ', ',', 9, 10, 11, 12, 13, 0xa0, 0x2028, 0x2029:
			l.next()

		case '#':
			// skip comments
			l.next()

		SKIP_COMMENT:
			for {
				switch l.char {
				case EOF, 10, 13, 0x2028, 0x2029:
					break SKIP_COMMENT
				default:
					l.next()
				}
			}

		default:
			return
		}
	}
}

/**
 * Reads a number token from the source file, either a float
 * or an int depending on whether a decimal point appears.
 *
 * Int:   -?(0|[1-9][0-9]*)
 * Float: -?(0|[1-9][0-9]*)(\.[0-9]+)?((E|e)(+|-)?[0-9]+)?
 */
func (l *Lexer) readNumber() Token {
	kind := TOKEN_INT
	start := l.position

	if l.char == '-' {
		l.next()
	}

	if l.char == '0' {
		l.next()
		if l.char >= '0' && l.char <= '9' {
			panic(SyntaxError(l.source, l.position,
				fmt.Sprintf("Invalid number, unexpected digit after 0: %c", l.char)))
		}
	} else {
		l.readDigits()
	}

	if l.char == '.' {
		kind = TOKEN_FLOAT
		l.next()
		l.readDigits()
	}

	if l.char == 'E' || l.char == 'e' {
		kind = TOKEN_FLOAT
		l.next()
		if l.char == '+' || l.char == '-' {
			l.next()
		}
		l.readDigits()
	}

	return newToken(kind, start, l.position, l.body[start:l.position])
}

/**
 * Returns the new position in the source after reading digits.
 */
func (l *Lexer) readDigits() {
	if l.char < '0' || l.char > '9' {
		c := "EOF"
		if l.char != EOF {
			c = string(l.char)
		}

		panic(SyntaxError(l.source, l.position,
			fmt.Sprintf(`Invalid number, expected digit but got: "%s".`, c)))
	}

	l.next()
	for l.char >= '0' && l.char <= '9' {
		l.next()
	}
}

/**
 * Reads a string token from the source file.
 *
 * "([^"\\\u000A\u000D\u2028\u2029]|(\\(u[0-9a-fA-F]{4}|["\\/bfnrt])))*"
 */
func (l *Lexer) readString() Token {
	start := l.position
	value := ""

	l.next()
	chunkStart := l.position
	for l.char != EOF &&
		l.char != '"' &&
		l.char != '\r' && l.char != '\n' &&
		l.char != 0x2028 && l.char != 0x2029 {
		if l.char == '\\' {
			value += l.source.Body[chunkStart:l.position]
			l.next()
			switch l.char {
			case '"':
				value += "\""
			case '/':
				value += "/"
			case '\\':
				value += "\\"
			case 'b':
				value += "\b"
			case 'f':
				value += "\f"
			case 'n':
				value += "\n"
			case 'r':
				value += "\r"
			case 't':
				value += "\t"

			case 'u':
				a := l.next()
				b := l.next()
				c := l.next()
				d := l.next()
				uniChar := uniCharCode(a, b, c, d)
				if uniChar < 0 {
					panic(SyntaxError(l.source, l.position, "Bad character escape sequence."))
				}
				value += string(uniChar)

			default:
				panic(SyntaxError(l.source, l.position, "Bad character escape sequence."))
			}
			chunkStart = l.position + 1
		}
		l.next()
	}

	if l.char != '"' {
		panic(SyntaxError(l.source, l.position, "Unterminated string."))
	}

	value += l.source.Body[chunkStart:l.position]
	return newToken(TOKEN_STRING, start, l.position+1, value)
}

/**
 * Converts four hexidecimal chars to the integer that the
 * string represents. For example, uniCharCode('0','0','0','f')
 * will return 15, and uniCharCode('0','0','f','f') returns 255.
 *
 * Returns a negative number on error, if a char was invalid.
 *
 * This is implemented by noting that char2hex() returns -1 on error,
 * which means the result of ORing the char2hex() will also be negative.
 */
func uniCharCode(a, b, c, d rune) rune {
	return char2hex(a)<<12 | char2hex(b)<<8 | char2hex(c)<<4 | char2hex(d)
}

/**
 * Converts a hex character to its integer value.
 * '0' becomes 0, '9' becomes 9
 * 'A' becomes 10, 'F' becomes 15
 * 'a' becomes 10, 'f' becomes 15
 *
 * Returns -1 on error.
 */
func char2hex(a rune) rune {
	switch {
	case a >= '0' && a <= '9':
		return a - '0'
	case a >= 'A' && a <= 'F':
		return a - 'A' + 10
	case a >= 'a' && a <= 'f':
		return a - 'a' + 10
	default:
		return -1
	}
}

/**
 * Reads an alphanumeric + underscore name from the source.
 *
 * [_A-Za-z][_0-9A-Za-z]*
 */
func (l *Lexer) readName() Token {
	start := l.position
	ch := l.next()
	for ch == '_' ||
		ch >= '0' && ch <= '9' ||
		ch >= 'A' && ch <= 'Z' ||
		ch >= 'a' && ch <= 'z' {
		ch = l.next()
	}

	end := l.position
	return newToken(TOKEN_NAME, start, end, l.source.Body[start:end])
}
