package compile

import (
	"testing"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestParseExpression(t *testing.T) {
	test := func(src string, expected string) {
		p := newParser(src)
		ast := expression(p, astBuilder).(Ast)
		Assert(t).That(ast.String(), Equals(expected))
	}
	test("123", "123")
	test("foo", "foo")
	test("true", "true")
	test("-123", "-123")
	test("1 + 2", "3")
	test("1 + 2 + 3", "6")
	test("1 + 2 - 3", "0")
	test("1 + a + b", "(+ 1 a b)")
	test("1 + a + b + 2", "(+ a b 3)")
	test("5 + a + b - 2", "(+ a b 3)")
	test("2 + a + b - 5", "(+ a b -3)")
	test("a | b & c", "(| a (& b c))")
	test("a ^ b ^ c", "(^ (^ a b) c)")
	test("a + b - c", "(+ a b (- c))")
	test("a * b / c", "(/ (* a b) c)")
	test("a + b * c", "(+ a (* b c))")
	test("++a", "(++ a)")
	test("a--", "(POSTDEC a)")
	test("a = 123", "(= a 123)")
	test("a += 123", "(+= a 123)")
	test("+ - ! ~ x", "(+ (- (! (~ x))))")
}
