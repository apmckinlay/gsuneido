// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

// Note: can't be in ast package because it uses parser

import (
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFinal(t *testing.T) {
	test := func(src string, expected string) {
		t.Helper()
		p := NewParser("function (p) {\n" + src + "\n}")
		f := p.Function()
		list := []string{}
		for v := range f.Final {
			list = append(list, v)
		}
		sort.Strings(list)
		assert.T(t).This(str.Join(",", list)).Like(expected)
	}

	test("123", "")
	test("x = 5", "x")                 // normal usage
	test("x = 5; y = 6; x + y", "x,y") // normal usage
	test("f = function(){ x = F() }; f()", "f")
	test("x = #foo", "x")
	test("p = 5", "")              // parameter
	test("x = 5; ++x", "")         // modification
	test("x = 5; x += 2", "")      // modification
	test("x = 5; x = 6", "")       // modification
	test("x = 5; b = {|x| }", "")  // block parameters
	test("x = 5; b = {|@x| }", "") // block parameters
	test(`x = 0
		for (i = 0; i < 10; ++i)
			x += i
		return x`, "")
}

func TestPropFold(t *testing.T) {
	core.DefaultSingleQuotes = true
	defer func() { core.DefaultSingleQuotes = false }()
	test := func(src string, expected string) {
		t.Helper()
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
	utest := func(src string) {
		t.Helper()
		p := NewParser("function () {\n" + src + "\n}")
		f := p.Function()
		assert.T(t).This(func() { ast.PropFold(f) }).
			Panics("uninitialized variable: u")
	}

	test("b = { it = it + 1 }",
		"Binary(Eq b Block(it \n Binary(Eq it Nary(Add it 1))))")
	test("b = {|x| x = x + 1 }",
		"Binary(Eq b Block(x \n Binary(Eq x Nary(Add x 1))))")

	test("x = 'x'; x $ 'a' $ x", "'x' \n 'xax'")

	test("a = b; a", "Binary(Eq a b) \n a")

	utest("throw u = 5; u")

	utest("for (; ; u=5) { u }")

	test("f(a = 5) and g(b = 6) and h(a + b)",
		"Nary(And Call(f 5) Call(g 6) Call(h 11))")
	utest("a and (u = 5); u")
	utest("a ? u = 5 : u")

	utest("if x { u = 5 } if y { f(u) }")
	utest("if false { u = true } if u is 'hello' { hello } world")

	test("return 123",
		"Return(123)") // no change
	test("x = 5; F(x)",
		"5 \n Call(F 5)") // propagate
	test("x = 5; F(-x)",
		"5 \n Call(F -5)") // unary
	test("a = b = c = 0; a + b + c",
		"0 \n 0")
	test("x = 5; x = 6; x",
		"Binary(Eq x 5) \n Binary(Eq x 6) \n x") // not final
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
	test("if x { a = 5; f(a) } b",
		"If(x { \n 5 \n Call(f 5) \n }) \n b")

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

	// or => in
	test("a is 1 or a is 2", "In(a [1 2])")
	test("a is 1 or a is 2 or b is 3 or b is 4",
		"Nary(Or In(a [1 2]) In(b [3 4]))")
	test("x or a is 1 or a is 2 or y or b is 3 or b is 4 or z",
		"Nary(Or x In(a [1 2]) y In(b [3 4]) z)")
	test("a is 1 or b is 2", "Nary(Or Binary(Is a 1) Binary(Is b 2))")

	// trinary
	test("true ? b : c", "b")  // fold
	test("false ? b : c", "c") // fold

	// range
	test(".x >= 0 and .x < 10",
		"InRange(Mem(this 'x') Gte 0 Lt 10)")
	test("f() and 5 < x and x <= 10 and g()",
		"Nary(And Call(f) InRange(x Gt 5 Lte 10) Call(g))")
	test("x > 0 and .x < 10",
		"Nary(And Binary(Gt x 0) Binary(Lt Mem(this 'x') 10))")
	test(".x > 0 and x < 10",
		"Nary(And Binary(Gt Mem(this 'x') 0) Binary(Lt x 10))")

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
	test("a and true and true", "Nary(And a true)") // keep op
	test("a and true and b", "Nary(And a b)")       // remove identity
	test("a or false or false", "Nary(Or a false)") // keep op
	test("a or true or b", "true")                  // short circuit
	test("a and false and b", "false")              // short circuit

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

	test("a in ()", "false")
}
