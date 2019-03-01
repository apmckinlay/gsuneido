package compile

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	rt "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestParseExpression(t *testing.T) {
	rt.DefaultSingleQuotes = true
	defer func() { rt.DefaultSingleQuotes = false }()
	parseExpr := func(src string) ast.Expr {
		t.Helper()
		p := newParser(src)
		result := p.expr()
		Assert(t).That(p.Token, Equals(tok.Eof))
		return result
	}
	xtest := func(src string, expected string) {
		t.Helper()
		err := Catch(func() { parseExpr(src) })
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
		if actual != expected {
			t.Errorf("%s expected: %s but got: %s", src, expected, actual)
		}
	}

	test("123", "")
	test("foo", "")
	test("true", "")
	test("a", "")
	test("this", "")
	test("default", "")

	test("a is true", "(Is a true)")

	test("a % b % c", "(Mod (Mod a b) c)")

	test("(123)", "123")
	test("a + b", "(Add a b)")
	test("a - b", "(Add a (Sub b))")
	test("a * b", "(Mul a b)")
	test("a / b", "(Mul a (Div b))")
	test("a + b * c", "(Add a (Mul b c))")
	test("(a + b) * c", "(Mul (LParen (Add a b)) c)")
	test("a * b + c", "(Add (Mul a b) c)")

	test("a + b", "(Add a b)")
	test("a - b", "(Add a (Sub b))")
	test("5 + a + b", "(Add 5 a b)")

	test("a $ b", "(Cat a b)")
	test("a $ b $ c", "(Cat a b c)")
	test("'foo' $ a $ 'bar'", "(Cat 'foo' a 'bar')")

	test("a | b & c", "(BitOr a (BitAnd b c))")
	test("a ^ b ^ c", "(BitXor a b c)")

	test("a + b - c", "(Add a b (Sub c))")
	test("a + b * c", "(Add a (Mul b c))")

	test("a % b * c", "(Mul (Mod a b) c)")
	test("a / b % c", "(Mod (Mul a (Div b)) c)")
	test("a * b * c", "(Mul a b c)")
	test("a * b / c", "(Mul a b (Div c))")
	test("++a", "(Inc a)")
	test("++a.b", "(Inc a.b)")
	test("a--", "(PostDec a)")
	test("a = 123", "(Eq a 123)")
	test("a = b = c", "(Eq a (Eq b c))")
	test("a += 123", "(AddEq a 123)")
	test("+ - not ~ x", "(Add (Sub (Not (BitNot x))))")
	test("+f()", "(Add (call f))")
	test("not f()", "(Not (call f))")

	test("a and b", "(And a b)")
	test("a and b and c", "(And a b c)")
	test("a or b", "(Or a b)")
	test("a or b or c", "(Or a b c)")

	test("a ? b : c", "(? a b c)")
	test("a \n ? b \n : c", "(? a b c)")
	test("a and b ? c + 1 : d * 2", "(? (And a b) (Add c 1) (Mul d 2))")
	test("a ? (b ? c : d) : (e ? f : g)", "(? a (LParen (? b c d)) (LParen (? e f g)))")
	test("a ?  b ? c : d  :  e ? f : g", "(? a (? b c d) (? e f g))")
	test("true ? b : c", "b")
	test("false ? b : c", "c")

	test("a in (1,2,3)", "(a in 1 2 3)")
	test("a not in (1,2,3)", "(Not (a in 1 2 3))")
	test("a in (1,2,3) in (true, false)", "((a in 1 2 3) in true false)")

	test("a.b", "")
	test(".a.b", "this._a.b")
	test("this.a.b", "")

	test("a[b]", "")
	test("a[b][c]", "")
	test("a[b + c]", "a[(Add b c)]")
	test("a[1..]", "")
	test("a[1..2]", "")
	test("a[..2]", "")
	test("a[1::]", "")
	test("a[1::2]", "")
	test("a[::2]", "")
	test("a[0::1][0]", "")

	test("b = { }", "(Eq b { })")
	test("b = {|a,b| }", "(Eq b {|a,b| })")
	test("b = {|@a| }", "(Eq b {|@a| })")

	test("f()", "(call f)")
	test("f(a, b)", "(call f a b)")
	test("f(@a)", "(call f '@': a)")
	test("f(@+1 a)", "(call f '@+1': a)")
	test("f(a:)", "(call f a: true)")
	test("f(a: 123)", "(call f a: 123)")
	test("f(123:)", "(call f 123: true)")
	test("f(123: 456)", "(call f 123: 456)")
	test("f(123: 456:)", "(call f 123: true 456: true)")
	test("f('a b':)", "(call f 'a b': true)")
	test("f('a b': 123)", "(call f 'a b': 123)")
	test("f(a: 1, b: 2)", "(call f a: 1 b: 2)")
	test("f(1, a: 2)", "(call f 1 a: 2)")
	test("f(1, is: 2)", "(call f 1 is: 2)")
	test("f(a: a)", "(call f a: a)")
	test("f(:a)", "(call f a: a)")
	test("f(){ }", "(call f block: { })")
	test("f({ })", "(call f { })")
	test("c.m(a, b)", "(call c.m a b)")
	test(".m()", "(call this._m)")
	test("false isnt x = F()", "(Isnt false (Eq x (call F)))")
	test("0xB2.Chr()", "(call 178.Chr)")

	test("F { }", "/* class : F */")
	test("a.F({ })",
		"(call a.F { })")
	test("a.F(block: { })",
		"(call a.F block: { })")
	test("a.F(){ }",
		"(call a.F block: { })")
	test("a.F { }",
		"(call a.F block: { })")

	test("super(1)", "(call super 1)")
	test("super.Foo(1)", "(call super.Foo 1)")

	test("new c", "(call c.*new*)")
	test("new c.m", "(call c.m.*new*)")
	test("new c(a, b)", "(call c.*new* a b)")
	test("new c.m(a, b)", "(call c.m.*new* a b)")

	test("[:a]", "(call Record a: a)")

	// folding ------------------------------------------------------

	// unary
	test("-123", "")
	test("not true", "false")
	test("(123)", "123")

	// binary
	test("8 % 3", "2")
	test("1 << 4", "16")
	test("'foobar' =~ 'oo'", "true")
	test("'foobar' !~ 'obo'", "true")

	// commutative
	test("a * 0 * b", "0") // short circuit
	test("a & 0 & b", "0") // short circuit
	test("1 * a * 1", "a") // skip identity
	test("1 + 2", "3")
	test("1 + 2 + 3", "6")
	test("1 + 2 - 3", "0")
	test("1 | 2 | 4", "7")
	test("255 & 15", "15")
	test("a and true and true", "a") // skip identity
	test("a or false or false", "a") // skip identity
	test("a or true or b", "true")   // short circuit
	test("a and false and b", "false")

	test("1 + a + b + 2", "(Add 3 a b)")
	test("5 + a + b - 2", "(Add 3 a b)")
	test("2 + a + b - 5", "(Add -3 a b)")
	test("a - 2 - 1", "(Add a -3)")

	test("1 * 8", "8")
	test("(1 * 8)", "8")
	test("1 / 8", ".125")
	test("2 / 8", ".25")
	test("2 * 4", "8")
	test("a / 2", "(Mul a .5)")
	test("8 / 2", "4")
	test("4 * 8 / 2", "16")
	test("2 * a * b", "(Mul 2 a b)")
	test("3 * a * b * 2", "(Mul 6 a b)")
	test("a * 6 * b / 3", "(Mul a 2 b)")
	test("a * 8 * b / 4", "(Mul a 2 b)")

	// concatenation
	test("'foo' $ 'bar'", "'foobar'")
	test("'foo' $ 'bar' $ b", "(Cat 'foobar' b)")
	test("a $ 'foo' $ 'bar' $ b", "(Cat a 'foobar' b)")
	test("a $ 'foo' $ 'bar'", "(Cat a 'foobar')")
	test(`'foo' $
		'bar'`, "'foobar'")
	test(`'foo' $
		'bar' $
		'baz'`, "'foobarbaz'")
}

func TestParseParams(t *testing.T) {
	test := func(src string) {
		t.Helper()
		p := newParser(src + "{}")
		result := p.method() // method to allow dot params
		Assert(t).That(p.Token, Equals(tok.Eof))
		s := result.String()
		s = s[8:] // remove "function"
		s = s[:len(s)-4]
		Assert(t).That(s, Equals(src))
	}
	test("()")
	test("(@a)")
	test("(a,b)")
	test("(a,b=1)")
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
		p := newParser(src + " }")
		stmts := p.statements()
		Assert(t).That(p.Token, Equals(tok.RCurly))
		s := ""
		sep := ""
		for _, stmt := range stmts {
			s += sep + stmt.String()
			sep = "\n"
		}
		Assert(t).That(s, Like(expected))
	}
	test("x=123;;", "(Eq x 123) {}")
	test("return", "return")
	test("return 123", "return 123")
	test("return \n 123", "return\n123")
	test("return; 123", "return\n123")
	test("return a + \n b", "return (Add a b)")
	test("return \n while b \n c", "return\nwhile b\nc")

	test("forever\na", "forever\na")

	test("while (a) { b }", "while a \n b")
	test("while a { b }", "while a \n b")
	test("while (a) \n b", "while a \n b")
	test("while a \n b", "while a \n b")
	test("while a \n ;", "while a \n {}")

	test("if (a) b", "if a \n b")
	test("if a \n b", "if a \n b")
	test("if (a) b else c", "if a \n b \n else \n c")
	test("if f() { b } else c", "if (call f) \n b \n else \n c")
	test("if F { b }", "if F \n b")

	test("switch { case 1: b }",
		"switch true { \n case 1 \n b \n }")
	test("switch { \n case x < 3: \n return -1 \n }",
		"switch true { \n case (Lt x 3) \n return -1 \n }")
	test("switch a { case 1,2: b case 3: c default: d }",
		"switch a { \n case 1,2 \n b \n case 3 \n c \n default: \n d \n }")

	test("throw 'fubar'", "throw 'fubar'")

	test("break", "break")
	test("continue", "continue")

	test("do a while b", "do \n a \n while b")

	test("for x in ob\na", "for x in ob \n a")
	test("for x in ob { a }", "for x in ob \n a")
	test("for (x in ob) a", "for x in ob \n a")

	test("for (;;) x", "for ; ; \n x")
	test("for (i = 0; i < 9; ++i) X",
		"for (Eq i 0); (Lt i 9); (Inc i) \n X")

	test("try x", "try \n x")
	test("try x catch y", "try \n x \n catch \n y")
	test("try x catch (e) y", "try \n x \n catch (e) \n y")
	test("try x catch (e, 'err') y", "try \n x \n catch (e,'err') \n y")

	test("+a \n -b", "(Add a) (Sub b)")
	test("a + b \n -c", "(Add a b) (Sub c)")
	test("a = b; .F()", "(Eq a b) (call this.F)")
	test("a = b; \n .F()", "(Eq a b) (call this.F)")
	test("a = b \n .F()", "(Eq a b) (call this.F)")

	xtest := func(src string, expected string) {
		t.Helper()
		actual := Catch(func() {
			p := newParser(src + "}")
			p.statements()
			Assert(t).That(p.Token, Equals(tok.Eof))
		}).(string)
		if !strings.Contains(actual, expected) {
			t.Errorf("%#v expected: %#v but got: %#v", src, expected, actual)
		}
	}
	xtest("a \n * b", "syntax error: unexpected '*'")
}
