package language

import (
	"fmt"
	"regexp"
	"strconv"
)

type QLError struct {
	Message   string
	Stack     string
	Nodes     []INode
	Source    *Source
	Positions []int
	Locations []SourceLocation
}

func NewQLError(message string, nodes []INode) QLError {
	return NewQLErrorWithSource(message, nodes, "", nil, nil)
}

func NewQLErrorWithSource(message string, nodes []INode, stack string, source *Source, positions []int) QLError {
	err := QLError{
		Message:   message,
		Stack:     stack,
		Nodes:     nodes,
		Source:    source,
		Positions: positions,
	}

	if stack == "" {
		err.Stack = message
	}

	if source == nil && len(nodes) > 0 {
		loc := nodes[0].Loc()
		if loc != nil {
			err.Source = loc.Source
		}
	}

	if positions == nil {
		flag := false
		nodePositions := make([]int, len(nodes))
		for i, node := range nodes {
			loc := node.Loc()
			if loc != nil {
				nodePositions[i] = loc.Start
				flag = flag || loc.Start > 0
			}
		}
		if flag {
			err.Positions = nodePositions
		}
	}

	if len(positions) > 0 && source != nil {
		locations := make([]SourceLocation, len(positions))
		for i, pos := range positions {
			locations[i] = getLocation(*source, pos)
		}
		err.Locations = locations
	}

	return err
}

func (e QLError) Error() string {
	return e.Message
}

type QLFormattedError struct {
	Message   string
	Locations []SourceLocation
}

func FormatError(err QLError) QLFormattedError {
	return QLFormattedError{
		Message:   err.Message,
		Locations: err.Locations,
	}
}

func LocatedError(err error, nodes []INode) QLError {
	message := "An unknown error occurred"
	stack := ""

	if err != nil {
		message = err.Error()

		if err, ok := err.(*QLError); ok {
			stack = err.Stack
		}
	}

	return NewQLErrorWithSource(message, nodes, stack, nil, nil)
}

func SyntaxError(source Source, position int, description string) QLError {
	location := getLocation(source, position)
	return NewQLErrorWithSource(
		fmt.Sprintf("Syntax Error %v (%v:%v) %v\n\n%v",
			source.Name, location.Line, location.Column, description,
			highlightSourceAtLocation(source, location)),
		nil, "", &source, []int{position})
}

/**
 * return high light
 */
func highlightSourceAtLocation(source Source, location SourceLocation) string {
	line := location.Line
	prevLineNum := strconv.Itoa(line - 1)
	lineNum := strconv.Itoa(line)
	nextLineNum := strconv.Itoa(line + 1)
	padLen := len(nextLineNum)
	lines := regexp.MustCompile("\r\n|[\n\r\u2028\u2029]").Split(source.Body, -1)

	result := ""
	if line >= 2 {
		result += lpad(padLen, prevLineNum) + ": " + lines[line-2] + "\n"
	}

	result += lpad(padLen, lineNum) + ": " + lines[line-1] + "\n" +
		lpad(1+padLen+location.Column, "") + "^\n"

	if line < len(lines) {
		result += lpad(padLen, nextLineNum) + ": " + lines[line] + "\n"
	}
	return result
}

/**
 * add (leng - len(str)) ' ' before str
 */
func lpad(leng int, str string) string {
	l := leng - len(str)
	a := make([]byte, l)
	for i := 0; i < l; i++ {
		a[i] = byte(' ')
	}
	return string(a) + str
}

type SourceLocation struct {
	Line   int
	Column int
}

/**
 * return line number and column number of error char
 * position is position in source.body
 */
func getLocation(source Source, position int) SourceLocation {
	// line := 1
	// lastChar := ' '
	// column := position + 1
	// startColumn := 0

	// for i, ch := range source.Body {
	// 	if i >= position {
	// 		column = position - startColumn
	// 		break
	// 	}

	// 	switch ch {
	// 	case '\r', '\n', '\u2028', '\u2029':
	// 		if ch != '\n' || lastChar != '\r' {
	// 			line++
	// 		}
	// 		startColumn = i
	// 	}
	// 	lastChar = ch
	// }

	line := 1
	column := 1

	for i, ch := range source.Body {
		if i >= position {
			break
		}

		switch ch {
		//if ch is newline, increase line, reset position
		case '\r', '\n', '\u2028', '\u2029': //all cases: '\r\n', '\r', '\n', '\u2028', '\u2029'
			nextCh := source.Body[i+1]
			if ch == '\r' && nextCh == '\n' { //case '\r\n'
				line--
			}
			line++
			column = 1
		default:
			column++
		}
	}

	return SourceLocation{
		Line:   line,
		Column: column,
	}
}
