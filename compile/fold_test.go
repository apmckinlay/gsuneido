// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

// Note: can't be in ast package because it uses parser

import (
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFinal(t *testing.T) {
	test := func(src string, expected string) {
		t.Helper()
		p := NewParser("function (p) {\n" + src + "\n}")
		f := p.Function()
		list := []string{}
		for v, lev := range f.Final {
			if lev != disqualified {
				list = append(list, v)
			}
		}
		sort.Strings(list)
		assert.T(t).This(str.Join(",", list)).Like(expected)
	}
	test("123", "")
	test("x = 5", "x")                 // normal usage
	test("x = 5; y = 6; x + y", "x,y") // normal usage
	test("p = 5", "")                  // parameter
	test("x = 5; ++x", "")             // modification
	test("x = 5; x += 2", "")          // modification
	test("x = 5; x = 6", "")           // modification
	test("x; x = 5", "")               // assignment after usage
	test("x = 5; b = {|x| }", "b")     // block parameters
	test("x = 5; b = {|@x| }", "b")    // block parameters
	test("x = 5; forever { x }", "x")  // usage in inner nesting level
	test("forever\n{ x = 5 }\n x", "") // usage at outer nesting level
}

func TestPropFold(t *testing.T) {
	test := func(src string, expected string) {
		t.Helper()
		rt.DefaultSingleQuotes = true
		defer func() { rt.DefaultSingleQuotes = false }()
		p := NewParser("function () {\n" + src + "\n}")
		f := p.Function()
		f = ast.PropFold(f) // this is what we're testing
		s := ""
		sep := ""
		for _, stmt := range f.Body {
			if stmt != nil {
				s += sep + stmt.String()
				sep = "\n"
			}
		}
		assert.T(t).This(s).Like(expected)
	}

	test("return 123",
		"Return(123)") // no change
	test("x = 5; F(x)",
		"5 \n Call(F 5)") // propagate
	test("x = 5; F(-x)",
		"5 \n Call(F -5)") // unary
	test("x = 5; x = 6; x",
		"Binary(Eq x 5) \n Binary(Eq x 6) \n x") // propagate
	test("x = 5; ++x",
		"Binary(Eq x 5) \n Unary(Inc x)") // don't inline lvalue
	test("x = 5; x--; x",
		"Binary(Eq x 5) \n Unary(PostDec x) \n x") // don't inline after update
	test("F(1 + 2)",
		"Call(F 3)") // simple fold
	test("F(1 + 2 * 3 / 6)",
		"Call(F 2)") // simple fold
	test("x = 2 * 4; F(x)",
		"8 \n Call(F 8)") // fold & propagate
	test("x = 2; y = 1 << x; F(y)",
		"2 \n 4 \n Call(F 4)") // binary
	test("c = true; F(c ? 4 : 5)",
		"true \n Call(F 4)") // trinary
	test("x = 3; x in (2,3,4)",
		"3 \n true") // in
	test("(123)",
		"123")
	test("i = true; while (i) { i = false }",
		"Binary(Eq i true) \n While(i Binary(Eq i false))")
	test("n = 5; b = { n }",
		"5 \n Binary(Eq b Block( \n 5))")

	// folding ------------------------------------------------------

	// unary
	test("-123", "-123")
	test("not true", "false")
	test("(123)", "123")

	// binary
	test("8 % 3", "2")
	test("1 << 4", "16")
	test("'foobar' =~ 'oo'", "true")
	test("'foobar' !~ 'obo'", "true")
	test("s =~ 'x'", "Binary(Match s 'x')")
	test("'hello' =~ 'lo'", "true")
	test("s = 'hello'; s =~ 'lo'", "'hello'\ntrue")

	// trinary
	test("true ? b : c", "b")  // fold
	test("false ? b : c", "c") // fold

	// if
	test("if (true) T() else F()",
		"Call(T)")
	test("if (false) T() else F()",
		"Call(F)")
	test("if (true) T()",
		"Call(T)")
	test("if (false) T()",
		"")
	test("if (false) if (false) if (false) T()",
		"")
	test("x=1; if (x > 1) F()",
		"1")

	// commutative
	test("a * 0 * b", "Nary(Mul a b 0)") // short circuit
	test("a & 0 & b", "0")               // short circuit
	test("1 * a * 1", "Nary(Mul a 1)")
	test("1 + 2", "3")
	test("1 + 2 + 3", "6")
	test("1 + 2 - 3", "0")
	test("1 | 2 | 4", "7")
	test("255 & 15", "15")
	test("a and true and true", "Nary(And a true)")
	test("a or false or false", "Nary(Or a false)")
	test("a or true or b", "true")     // short circuit
	test("a and false and b", "false") // short circuit

	test("1 + a + b + 2", "Nary(Add 3 a b)")
	test("5 + a + b - 2", "Nary(Add 3 a b)")
	test("2 + a + b - 5", "Nary(Add -3 a b)")
	test("a - 2 - 1", "Nary(Add a -3)")

	test("1 * 8", "8")
	test("(1 * 8)", "8")
	test("1 / 8", ".125")
	test("2 / 8", ".25")
	test("2 * 4", "8")
	test("a / 2", "Nary(Mul a Unary(Div 2))")
	test("8 / 2", "4")
	test("4 * 8 / 2", "16")
	test("2 * a * b", "Nary(Mul a b 2)")
	test("3 * a * b * 2", "Nary(Mul a b 6)")
	test("a * 6 * b / 3", "Nary(Mul a b 2)")
	test("a * 8 * b / 4", "Nary(Mul a b 2)")

	// concatenation
	test("'foo' $ 'bar'", "'foobar'")
	test("'foo' $ 'bar' $ b", "Nary(Cat 'foobar' b)")
	test("a $ 'foo' $ 'bar' $ b", "Nary(Cat a 'foobar' b)")
	test("a $ 'foo' $ 'bar'", "Nary(Cat a 'foobar')")
	test(`'foo' $
		'bar'`, "'foobar'")
	test(`'foo' $
		'bar' $
		'baz'`, "'foobarbaz'")
}
