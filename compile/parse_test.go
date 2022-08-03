// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestClassName(t *testing.T) {
	p := NewParser("")
	test := func(in, expected string) {
		t.Helper()
		classNum.Store(0)
		p.name = in
		assert.T(t).This(p.getClassName()).Is(expected)
	}
	test("", "Class1")
	test("?", "Class1")
	test("Foo", "Foo")
	test("Foo Bar", "Bar")
	test("Foo Bar Baz", "Baz")
	test("Foo ?", "Class1")
	test("Foo Bar ?", "Class1")
}

func TestParseExpression(t *testing.T) {
	className := ""
	rt.DefaultSingleQuotes = true
	defer func() { rt.DefaultSingleQuotes = false }()
	parseExpr := func(src string) ast.Expr {
		t.Helper()
		p := NewParser(src)
		p.className = className
		p.InitFuncInfo()
		result := p.Expression()
		assert.T(t).This(p.Token).Is(tok.Eof)
		return result
	}
	xtest := func(src string, expected string) {
		t.Helper()
		err := assert.Catch(func() { parseExpr(src) })
		if actual, ok := err.(string); ok {
			if !strings.Contains(actual, expected) {
				t.Errorf("\n%#v\nexpect: %#v\nactual: %#v", src, expected, actual)
			}
		} else {
			t.Error("unexpected:", err)
		}
	}
	xtest("1 = 2", "lvalue required")
	xtest("a = 5 = b", "lvalue required")
	xtest("++123", "lvalue required")
	xtest("123--", "lvalue required")
	xtest("++123--", "lvalue required")
	xtest("a.''", "expecting identifier")
	xtest("f(a:, b:, 'a':)", "duplicate argument name")
	xtest("f(a:, b:, :b)", "duplicate argument name")
	xtest("f(1, 2, a:, b: 3, 4", "un-named arguments must come before named arguments")

	test := func(src string, expected string) {
		t.Helper()
		if expected == "" {
			expected = src
		}
		ast := parseExpr(src)
		actual := ast.String()
		assert.T(t).This(actual).Like(expected)
	}

	test("123", "")
	test("foo", "")
	test("true", "")
	test("a", "")
	test("this", "")
	test("default", "")

	test("a is true", "Binary(Is a true)")

	test("a % b % c", "Binary(Mod Binary(Mod a b) c)")

	test("(123)", "123")
	test("a + b", "Nary(Add a b)")
	test("a - b", "Nary(Add a Unary(Sub b))")
	test("a * b", "Nary(Mul a b)")
	test("a / b", "Nary(Mul a Unary(Div b))")
	test("3 / a", "Nary(Mul Unary(Div a) 3)")
	test("1 / a", "Unary(Div a)")
	test("a + b * c", "Nary(Add a Nary(Mul b c))")
	test("(a + b) * c", "Nary(Mul Unary(LParen Nary(Add a b)) c)")
	test("a * b + c", "Nary(Add Nary(Mul a b) c)")

	test("a + b", "Nary(Add a b)")
	test("a - b", "Nary(Add a Unary(Sub b))")
	test("5 + a + b", "Nary(Add 5 a b)")

	test("a $ b", "Nary(Cat a b)")
	test("a $ b $ c", "Nary(Cat a b c)")
	test("'foo' $ a $ 'bar'", "Nary(Cat 'foo' a 'bar')")
	test("'x' $ 'y'.Repeat()", "Nary(Cat 'x' Call(Mem('y' 'Repeat')))")

	test("a | b & c", "Nary(BitOr a Nary(BitAnd b c))")
	test("a ^ b ^ c", "Nary(BitXor a b c)")

	test("a + b - c", "Nary(Add a b Unary(Sub c))")
	test("a + b * c", "Nary(Add a Nary(Mul b c))")

	test("a % b * c", "Nary(Mul Binary(Mod a b) c)")
	test("a / b % c", "Binary(Mod Nary(Mul a Unary(Div b)) c)")
	test("a * b * c", "Nary(Mul a b c)")
	test("a * b / c", "Nary(Mul a b Unary(Div c))")
	test("++a", "Unary(Inc a)")
	test("++a.b", "Unary(Inc Mem(a 'b'))")
	test("a--", "Unary(PostDec a)")
	test("a = 123", "Binary(Eq a 123)")
	test("a = b = c", "Binary(Eq a Binary(Eq b c))")
	test("a += 123", "Binary(AddEq a 123)")
	test("+ - not ~ x", "Unary(Add Unary(Sub Unary(Not Unary(BitNot x))))")
	test("+f()", "Unary(Add Call(f))")
	test("not f()", "Unary(Not Call(f))")

	test("a and b", "Nary(And a b)")
	test("a and b and c", "Nary(And a b c)")
	test("a or b", "Nary(Or a b)")
	test("a or b or c", "Nary(Or a b c)")

	test("a ? b : c", "Trinary(a b c)")
	test("a \n ? b \n : c", "Trinary(a b c)")
	test("a and b ? c + 1 : d * 2", "Trinary(Nary(And a b) Nary(Add c 1) Nary(Mul d 2))")
	test("a ? (b ? c : d) : (e ? f : g)",
		"Trinary(a Unary(LParen Trinary(b c d)) Unary(LParen Trinary(e f g)))")
	test("a ?  b ? c : d  :  e ? f : g", "Trinary(a Trinary(b c d) Trinary(e f g))")

	test("a in (1,2,3)", "In(a [1 2 3])")
	test("a not in (1,2,3)", "Unary(Not In(a [1 2 3]))")
	test("a in (1,2,3) in (true, false)", "In(In(a [1 2 3]) [true false])")

	test("a.b", "Mem(a 'b')")
	test(".a.b", "Mem(Mem(this 'a') 'b')") // not privatized
	className = "Foo"
	test(".a.b", "Mem(Mem(this 'Foo_a') 'b')") // privatized
	test("this.a.b", "Mem(Mem(this 'a') 'b')") // not privatized
	className = ""

	test("a[b]", "Mem(a b)")
	test("a[b][c]", "Mem(Mem(a b) c)")
	test("a[b + c]", "Mem(a Nary(Add b c))")
	test("a[1..]", "RangeTo(a 1 <nil>)")
	test("a[1..2]", "RangeTo(a 1 2)")
	test("a[..2]", "RangeTo(a <nil> 2)")
	test("a[1::]", "RangeLen(a 1 <nil>)")
	test("a[1::2]", "RangeLen(a 1 2)")
	test("a[::2]", "RangeLen(a <nil> 2)")
	test("a[0::1][0]", "Mem(RangeLen(a 0 1) 0)")

	test("b = { }", "Binary(Eq b Block())")
	test("b = { it }", "Binary(Eq b Block(it \n it))")
	test("b = {|a,b| }", "Binary(Eq b Block(a,b))")
	test("b = {|@a| }", "Binary(Eq b Block(@a))")

	test("f()", "Call(f)")
	test("f(a, b)", "Call(f a b)")
	test("f(@a)", "Call(f '@':a)")
	test("f(@+1 a)", "Call(f '@+1':a)")
	test("f(a:)", "Call(f a:true)")
	test("f(a: 123)", "Call(f a:123)")
	test("f(123:)", "Call(f 123:true)")
	test("f(123: 456)", "Call(f 123:456)")
	test("f(123: 456:)", "Call(f 123:true 456:true)")
	test("f('a b':)", "Call(f 'a b':true)")
	test("f('a b': 123)", "Call(f 'a b':123)")
	test("f(a: 1, b: 2)", "Call(f a:1 b:2)")
	test("f(1, a: 2)", "Call(f 1 a:2)")
	test("f(1, is: 2)", "Call(f 1 is:2)")
	test("f(a: a)", "Call(f a:a)")
	test("f(:a)", "Call(f a:a)")
	test("f(){ }", "Call(f block:Block())")
	test("f({ })", "Call(f Block())")
	test("c.m(a, b)", "Call(Mem(c 'm') a b)")
	className = "Foo"
	test(".m()", "Call(Mem(this 'Foo_m'))")
	className = ""
	test("false isnt x = F()", "Binary(Isnt false Binary(Eq x Call(F)))")
	test("0xB2.Chr()", "Call(Mem(178 'Chr'))")

	test("F { }", "/* class : F */")
	test("a.F({ })",
		"Call(Mem(a 'F') Block())")
	test("a.F(block:{ })",
		"Call(Mem(a 'F') block:Block())")
	test("a.F(){ }",
		"Call(Mem(a 'F') block:Block())")
	test("a.F { }",
		"Call(Mem(a 'F') block:Block())")

	test("super(1)", "Call(super 1)")
	test("super.Foo(1)", "Call(Mem(super 'Foo') 1)")

	test("new c", "Call(Mem(c '*new*'))")
	test("new c.m", "Call(Mem(Mem(c 'm') '*new*'))")
	test("new c(a, b)", "Call(Mem(c '*new*') a b)")
	test("new c.m(a, b)", "Call(Mem(Mem(c 'm') '*new*') a b)")

	test("[:a]", "Call(Record a:a)")

}

func TestParseParams(t *testing.T) {
	test := func(src string) {
		t.Helper()
		p := NewParser(src + "{}")
		result := p.function(true) // method to allow dot params
		assert.T(t).This(p.Token).Is(tok.Eof)
		s := result.String()
		s = s[8:] // remove "Function"
		assert.T(t).This(s).Is(src)
	}
	test("()")
	test("(@a)")
	test("(a,b)")
	test("(ab=1)")
	test("(a=1)")
	test("(a,b=1)")
	test("(_a,_b=1)")
	test("(.a,._b=1)")
}

func TestParseStatements(t *testing.T) {
	rt.DefaultSingleQuotes = true
	defer func() { rt.DefaultSingleQuotes = false }()
	test := func(src string, expected string) {
		t.Helper()
		p := NewParser(src + " }")
		p.InitFuncInfo()
		stmts := p.statements()
		assert.T(t).This(p.Token).Is(tok.RCurly)
		s := ""
		sep := ""
		for _, stmt := range stmts {
			s += sep + stmt.String()
			sep = "\n"
		}
		assert.T(t).This(s).Like(expected)
	}
	test("x=123;;", "Binary(Eq x 123)\n{}")

	// return
	test("return", "Return()")
	test("return 123", "Return(123)")
	test("return \n 123", "Return()\n123")
	test("return; 123", "Return()\n123")
	test("return a + \n b", "Return(Nary(Add a b))")
	test("return \n while b \n c", "Return()\nWhile(b c)")

	test("forever\na", "Forever(a)")

	// while
	test("while (a) { b }", "While(a b)")
	test("while a { b }", "While(a b)")
	test("while (a) \n b", "While(a b)")
	test("while a \n b", "While(a b)")
	test("while a \n ;", "While(a {})")

	// if-else
	test("if (a) stmt", "If(a stmt)")
	test("if a \n stmt", "If(a stmt)")
	test("if (a) stmt else stmt2", "If(a stmt \n else stmt2)")
	test("if f() { stmt } else stmt2", "If(Call(f) stmt \n else stmt2)")
	test("if F { stmt }", "If(F stmt)")
	test("if ((a = b) is true) { stmt }",
		"If(Unary(LParen Binary(Is Unary(LParen Binary(Eq a b)) true)) stmt)")

	// switch
	test("switch { case 1: b }",
		"Switch(true \n Case(1 \n b))")
	test("switch { \n case x < 3: \n return 123 \n }",
		"Switch(true \n Case(Binary(Lt x 3) \n Return(123)))")
	test("switch a { case 1,2: b case 3: c default: d }",
		"Switch(a \n Case(1,2 \n b) \n Case(3 \n c) \n d)")
	test("switch a { case 1,2: b case 3: c default: }",
		"Switch(a \n Case(1,2 \n b) \n Case(3 \n c) \n ())")

	test("throw 'fubar'", "Throw('fubar')")
	test("throw 123 throw 456", "Throw(123) \n Throw(456)")

	test("break", "Break")

	test("continue", "Continue")

	// do-while
	test("do stmt while e", "DoWhile(stmt e)")

	// for-in
	test("for x in ob\nstmt", "ForIn(x ob\nstmt)")
	test("for x in ob { stmt }", "ForIn(x ob\nstmt)")
	test("for (x in ob) stmt", "ForIn(x ob\nstmt)")

	// for
	test("for (;;) stmt", "For(; ; \n stmt)")
	test("for (i = 0; i < 9; ++i) stmt",
		"For(Binary(Eq i 0); Binary(Lt i 9); Unary(Inc i) \n stmt)")

	// try-catch
	test("try stmt", "Try(stmt)")
	test("try stmt catch stmt2", "Try(stmt \n catch stmt2)")
	test("try stmt catch (e) stmt2", "Try(stmt \n catch (e) stmt2)")
	test("try stmt catch (e, 'err') stmt2", "Try(stmt \n catch (e,'err') stmt2)")

	// newline handling
	test("+a \n -b", "Unary(Add a)\nUnary(Sub b)")
	test("a + b \n -c", "Nary(Add a b)\nUnary(Sub c)")
	test("a + \n b + c", "Nary(Add a b c)")
	test("a = b; .F()", "Binary(Eq a b)\nCall(Mem(this 'F'))")
	test("a = b; \n .F()", "Binary(Eq a b)\nCall(Mem(this 'F'))")
	test("a = b \n .F()", "Binary(Eq a b)\nCall(Mem(this 'F'))")
	test("a = b .F()", "Binary(Eq a Call(Mem(b 'F')))")

	xtest := func(src string, expected string) {
		t.Helper()
		actual := assert.Catch(func() {
			p := NewParser(src + "}")
			p.statements()
			assert.T(t).This(p.Token).Is(tok.Eof)
		}).(string)
		if !strings.Contains(actual, expected) {
			t.Errorf("%#v expected: %#v but got: %#v", src, expected, actual)
		}
	}
	xtest("a \n * b", "syntax error @4 unexpected '*'")
	xtest("foo bar", "syntax error")
	xtest("1+2 3+4", "syntax error")
	xtest("1+2 3+4", "syntax error")
	xtest("return 1+2 3+4", "syntax error")
	xtest("throw 1+2 3+4", "syntax error")
}
