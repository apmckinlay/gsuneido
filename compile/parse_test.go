package compile

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestParseExpression(t *testing.T) {
	parseExpr := func(src string) *Ast {
		p := newParser(src)
		return expression(p, astBuilder).(*Ast)
	}
	xtest := func(src string, expected string) {
		actual := Catch(func() { parseExpr(src) })
		if actual != expected {
			t.Errorf("%#v expected: %#v but got: %#v", src, expected, actual)
		}
	}
	xtest("1 = 2", "syntax error: lvalue required")
	xtest("a = 5 = b", "syntax error: lvalue required")
	xtest("++123", "syntax error: lvalue required")
	xtest("123--", "syntax error: lvalue required")
	xtest("++123--", "syntax error: lvalue required")

	test := func(src string, expected string) {
		t.Helper()
		ast := parseExpr(src)
		actual := ast.String()
		if actual != expected {
			t.Errorf("%s expected: %s but got: %s", src, expected, actual)
		}
	}
	test("123", "123")
	test("foo", "foo")
	test("true", "true")
	test("-123", "-123")
	test("a", "a")
	test("this", "this")

	test("1 + 2", "3")
	test("1 + 2 + 3", "6")
	test("1 + 2 - 3", "0")
	test("1 | 2 | 4", "7")
	test("255 & 15", "15")
	test("a or true or b", "true")
	test("a and false and b", "false")

	test("a % b % c", "(% (% a b) c)")

	test("a + b * c", "(+ a (* b c))")
	test("(a + b) * c", "(* (+ a b) c)")
	test("a * b + c", "(+ (* a b) c)")

	test("a + b", "(+ a b)")
	test("a - b", "(+ a (- b))")
	test("1 + a + b", "(+ 1 a b)")

	test("1 + a + b + 2", "(+ a b 3)")
	test("5 + a + b - 2", "(+ a b 3)")
	test("2 + a + b - 5", "(+ a b -3)")
	test("a - 2 - 1", "(+ a -3)")

	test("a $ b", "($ a b)")
	test("a $ b $ c", "($ a b c)")
	test("'foo' $ 'bar'", "'foobar'")
	test("'foo' $ a $ 'bar'", "($ 'foo' a 'bar')")
	test("'foo' $ 'bar' $ b", "($ 'foobar' b)")
	test("a $ 'foo' $ 'bar' $ b", "($ a 'foobar' b)")
	test("a $ 'foo' $ 'bar'", "($ a 'foobar')")

	test("a | b & c", "(| a (& b c))")
	test("a ^ b ^ c", "(^ a b c)")

	test("a + b - c", "(+ a b (- c))")
	test("a + b * c", "(+ a (* b c))")

	test("8 % 3", "2")
	test("2 * 4", "8")
	test("8 / 2", "4")
	test("4 * 8 / 2", "16")
	test("1 * a * b", "(* 1 a b)")
	test("3 * a * b * 2", "(* a b 6)")
	test("6 * a * b / 3", "(* a b 2)")
	test("8 * a * b / 4", "(* a b 2)")
	test("a % b * c", "(* (% a b) c)")
	test("a / b % c", "(% (* a (/ b)) c)")
	test("a * b * c", "(* a b c)")
	test("a * b / c", "(* a b (/ c))")
	test("++a", "(++ a)")
	test("++a.b", "(++ (. a b))")
	test("a--", "(post a)")
	test("a = 123", "(= a 123)")
	test("a = b = c", "(= a (= b c))")
	test("a += 123", "(+= a 123)")
	test("+ - ! ~ x", "(+ (- (! (~ x))))")

	test("a and b", "(and a b)")
	test("a and b and c", "(and a b c)")
	test("a or b", "(or a b)")
	test("a or b or c", "(or a b c)")

	test("a ? b : c", "(? a b c)")
	test("a and b ? c + 1 : d * 2", "(? (and a b) (+ c 1) (* d 2))")
	test("a ? (b ? c : d) : (e ? f : g)", "(? a (? b c d) (? e f g))")
	test("a ?  b ? c : d  :  e ? f : g", "(? a (? b c d) (? e f g))")

	test("a in (1,2,3)", "(in a 1 2 3)")
	test("a in (1,2,3) in (true, false)", "(in (in a 1 2 3) true false)")

	test("a.b", "(. a b)")
	test(".a.b", "(. (. this a) b)")
	test("this.a.b", "(. (. this a) b)")

	test("a[b]", "([ a b)")
	test("a[b][c]", "([ ([ a b) c)")
	test("a[b + c]", "([ a (+ b c))")
	test("a[1..]", "([ a (.. 1 2147483647))")
	test("a[1..2]", "([ a (.. 1 2))")
	test("a[..2]", "([ a (.. 0 2))")
	test("a[1::]", "([ a (:: 1 2147483647))")
	test("a[1::2]", "([ a (:: 1 2))")
	test("a[::2]", "([ a (:: 0 2))")

	test("b = { }", "(= b (block blockParams STMTS))")
	test("b = {|a,b| a; b }", "(= b (block (blockParams a b) (STMTS a b)))")
	test("b = {|@a| a }", "(= b (block (blockParams @a) (STMTS a)))")

	test("f()", "(call f args)")
	test("f(a, b)", "(call f (args (noKwd a) (noKwd b)))")
	test("f(@a)", "(call f (atArg a))")
	test("f(@+1 a)", "(call f (at1Arg a))")
	test("f(a:)", "(call f (args (a true)))")
	test("f(a: 1, b: 2)", "(call f (args (a 1) (b 2)))")
	test("f(1, a: 2)", "(call f (args (noKwd 1) (a 2)))")
	test("f(){ b }", "(call f (args (blockArg (block blockParams (STMTS b)))))")
	test("f({ b })", "(call f (args (noKwd (block blockParams (STMTS b)))))")
	test("c.m(a, b)", "(call (. c m) (args (noKwd a) (noKwd b)))")
	test(".m()", "(call (. this m) args)")
	test("false isnt x = F()", "(isnt false (= x (call F args)))")

	test("new c", "(new c args)")
	test("new c.m", "(new (. c m) args)")
	test("new c(a, b)", "(new c (args (noKwd a) (noKwd b)))")
	test("new c.m(a, b)", "(new (. c m) (args (noKwd a) (noKwd b)))")
}

func TestParseFunction(t *testing.T) {
	test := func(src, expected string) {
		result := ParseFunction(src)
		Assert(t).That(result.String(), Equals(expected))
	}
	test("function () { }", "(function params STMTS)")
	test("function (@a) { }", "(function (params @a) STMTS)")
	//test("function (@+1 a) { }", "(function (params @+1a) STMTS)")
	test("function (a, b) { }", "(function (params a b) STMTS)")
	test("function (a, b = 1) { }", "(function (params a (b 1)) STMTS)")
	test("function (a = 1) { }", "(function (params (a 1)) STMTS)")
	test("function (a, b = 1) { }", "(function (params a (b 1)) STMTS)")
	test("function (.a, .B, _c) { }", "(function (params .a .B _c) STMTS)")
}

func TestParseStatements(t *testing.T) {
	test := func(src string, expected string) {
		p := newParser(src + " }")
		ast := p.statements()
		ast = ast.first()
		s := ast.String()
		Assert(t).That(s, Like(expected))
	}
	test("return", "return")
	test("return a + b", "(return (+ a b))")
	test("forever\na", "(forever a)")
	test("while (a) { b }", "(while a (STMTS b))")
	test("while a { b }", "(while a (STMTS b))")
	test("while (a)\nb", "(while a b)")
	test("while a\nb", "(while a b)")
	test("while a\n;", "(while a STMTS)")

	test("if (a) b", "(if a b)")
	test("if (a) b else c", "(if a b c)")

	test("switch { case 1: b }",
		"(switch true (cases ( (vals 1) (STMTS b))))")
	test(`switch {
		case x < 3: return -1
		}`,
		"(switch true (cases ( (vals (< x 3)) (STMTS (return -1)))))")
	test("switch a { case 1,2: b case 3: c default: d }", `
		(switch a
		    (cases
		    	( (vals 1 2) (STMTS b))
				( (vals 3) (STMTS c)))
		    (STMTS d))`)
	test("throw 'fubar'", "(throw 'fubar')")

	test("break", "break")
	test("continue", "continue")

	test("do a while b", "(do a b)")

	test("for x in ob\na", "(for-in x ob a)")
	test("for x in ob { a }", "(for-in x ob (STMTS a))")
	test("for (x in ob) a", "(for-in x ob a)")

	test("for (i = 0; i < 9; ++i) X",
		"(for (exprs (= i 0)) (< i 9) (exprs (++ i)) X)")

	test("try x", "(try x)")
	test("try x catch y", "(try x (catch y))")
	test("try x catch (e) y", "(try x (catch e y))")
	test("try x catch (e, 'err') y", "(try x (catch e 'err' y))")
}

func BenchmarkParseExpr(b *testing.B) {
	for n := 0; n < b.N; n++ {
		p := newParser("a = b + c")
		expression(p, astBuilder)
	}
}
