// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

import (
	"fmt"
	"strings"
	"testing"

	lex "github.com/apmckinlay/gsuneido/compile/lexer"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// TestRegistered runs all tests registered by Register
func TestRegistered(t *testing.T) {
	if !RunRegistered(t) {
		t.Fail()
	}
}

/*
Register adds tests to be run by TestRegistered

Each registration provides:

- a suite name
- a source string containing row-based test cases

Example:

```go
var _ = Register("basic", `
"123", 123
"1 + 2", 3
"throw 'x'" throws "x"
`)

The source string is parsed with the Suneido lexer.

- One test case per row.
- A row ends at newline unless continued by syntax (e.g. multiline strings/code).
- Commas between row values are optional and ignored where valid.
- Comments are allowed.

Supported row forms:

- `code`
- `code, expected`
- `code throws "exception substring"`
*/
func Register(name, source string) bool {
	tests = append(tests, registered{name: name, source: source})
	return true
}

type registered struct {
	name   string
	source string
}

var tests []registered

// RunRegistered runs all tests registered by Register
func RunRegistered(t *testing.T) bool {
	ok := true
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if !runSource(tc.source) {
				ok = false
				t.Fail()
			}
		})
	}
	return ok
}

func runSource(source string) bool {
	p := parser{lxr: lex.NewLexer(source)}
	p.next(true)
	n := 0
	ok := true
	for p.Token != tok.Eof {
		row := []string{}
		str := []bool{} // parallel array, whether arg was a quoted string
		for {
			str = append(str, p.Token == tok.String)
			text := p.Text
			if p.Token == tok.Sub || p.Token == tok.Add {
				p.next(false)
				text += p.Text
			}
			row = append(row, text)
			p.next(false)
			if p.Token == tok.Comma {
				p.next(true)
			}
			if p.Token == tok.Eof || p.Token == tok.Newline {
				break
			}
		}
		if !runCase(row, str) {
			ok = false
		} else {
			n++
		}
		p.next(true)
	}
	return ok
}

type parser struct {
	lxr     *lex.Lexer
	comment string
	lex.Item
}

func (p *parser) next(skip bool) {
	p.comment = ""
	nl := false
	for {
		p.Item = p.lxr.Next()
		switch p.Token {
		case tok.Newline:
			if !skip {
				return
			}
			nl = true
		case tok.Whitespace:
			continue
		case tok.Comment:
			// capture trailing comment on same line
			if !nl {
				p.comment = p.Item.Text
			}
			continue
		default:
			return
		}
	}
}

func runCase(row []string, str []bool) (ok bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("\tFAILED: ", Fmt(row, str))
			fmt.Println("\tthrew: ", err)
			//dbg.PrintStack()
			ok = false
		}
	}()
	ok = true
	if !execute(row, str) {
		fmt.Println("\tFAILED: ", Fmt(row, str))
		ok = false
	}
	return
}

func Fmt(row []string, str []bool) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	sep := ""
	for i, s := range row {
		sb.WriteString(sep)
		sep = ", "
		if str[i] {
			sb.WriteRune('`')
		}
		sb.WriteString(s)
		if str[i] {
			sb.WriteRune('`')
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func init() {
	builtin.DefDef()
}

func execute(args []string, _ []bool) bool {
	src := "function () {\n" + args[0] + "\n}"
	var th Thread
	expected := "**notfalse**"
	if len(args) > 1 {
		expected = args[1]
	}
	var success bool
	var actual Value
	if expected == "throws" {
		expected = "throws " + args[2]
		e := assert.Catch(func() {
			fn := compile.Constant(src).(*SuFunc)
			actual = th.Call(fn)
		})
		if e == nil {
			success = false
		} else if es, ok := e.(string); ok {
			actual = SuStr(es)
			success = strings.Contains(es, args[2])
		} else if ss, ok := e.(SuStr); ok {
			actual = ss
			success = strings.Contains(string(ss), args[2])
		} else if se, ok := e.(*SuExcept); ok {
			actual = se.SuStr
			success = strings.Contains(string(se.SuStr), args[2])
		} else {
			actual = SuStr(fmt.Sprintf("%#v", e))
			success = false
		}
	} else {
		fn := compile.Constant(src).(*SuFunc)
		actual = th.Call(fn)
		if actual == nil {
			success = expected == "nil"
		} else if expected == "**notfalse**" {
			success = actual != False
		} else {
			expectedValue := compile.Constant(expected)
			success = actual.Equal(expectedValue)
			expected = WithType(expectedValue)
		}
	}
	if !success {
		fmt.Printf("\tgot: %s  expected: %s\n", WithType(actual), expected)
	}
	return success
}
