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
	test("a + b", "(+ a b)")
	test("a - b", "(+ a (- b))")
	test("1 + a + b", "(+ 1 a b)")
	test("1 + a + b + 2", "(+ a b 3)")
	test("5 + a + b - 2", "(+ a b 3)")
	test("2 + a + b - 5", "(+ a b -3)")

	test("a $ b", "($ a b)")
	test("a $ b $ c", "($ a b c)")
	test("'foo' $ 'bar'", "'foobar'")
	test("'foo' $ a $ 'bar'", "($ 'foo' a 'bar')")
	test("'foo' $ 'bar' $ b", "($ 'foobar' b)")
	test("a $ 'foo' $ 'bar' $ b", "($ a 'foobar' b)")
	test("a $ 'foo' $ 'bar'", "($ a 'foobar')")

	test("a | b & c", "(| a (& b c))")
	test("a ^ b ^ c", "(^ (^ a b) c)")

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
	test("a--", "(POSTDEC a)")
	test("a = 123", "(= a 123)")
	test("a += 123", "(+= a 123)")
	test("+ - ! ~ x", "(+ (- (! (~ x))))")

	test("a and b", "(and a b)")
	test("a and b and c", "(and a b c)")
	test("a or b", "(or a b)")
	test("a or b or c", "(or a b c)")

	test("a ? b : c", "(? a b c)")

	test("a in (1,2,3)", "(in a 1 2 3)")
}

func TestParseStatement(t *testing.T) {
	test := func(src string, expected string) {
		ast := ParseFunction("function () {\n" + src + "\n}")
		ast = ast.second() // function
		ast = ast.first()  // statements
		s := ast.String()
		//fmt.Println(s)
		Assert(t).That(s, Equals(expected))
	}
	test("while (a) { b }", "(while a (STMTS (EXPR b)))")
	test("while a { b }", "(while a (STMTS (EXPR b)))")
	test("while (a)\nb", "(while a (EXPR b))")
	test("while a\nb", "(while a (EXPR b))")
	test("while a\n;", "(while a NIL)")

	test("if (a) b", "(if a (EXPR b))")
	test("if (a) b else c", "(if a (EXPR b) (EXPR c))")

	test("switch { case 1: b }",
		"(switch true (cases ( (vals 1) (STMTS (EXPR b)))))")
	test(`switch { 
		case x < 3: return -1
		}`,
		"(switch true (cases ( (vals (< x 3)) (STMTS (return -1)))))")
	test("switch a { case 1,2: b case 3: c default: d }", `(switch
    a
    (cases
        ( (vals 1 2) (STMTS (EXPR b)))
        ( (vals 3) (STMTS (EXPR c))))
    (STMTS (EXPR d)))`)
}
