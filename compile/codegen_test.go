package compile

import (
	"os"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func ExampleSrcPos() {
	src := `function () {
		a = 123
		b = 4
		return a + b
		}`
	ast := parseFunction(src)
	fn := codegen(ast)
	DisasmMixed(os.Stdout, fn, src)
	// Output:
	// 16: a = 123
	// 	0: Int 123
	// 	3: Store a
	// 	5: Pop
	// 26: b = 4
	// 	6: Int 4
	// 	9: Store b
	// 	11: Pop
	// 34: return a + b
	// 	12: Load a
	// 	14: Load b
	// 	16: Add
}

func TestCodegen(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		classNum = 0
		ast := parseFunction("function () {\n" + src + "\n}")
		fn := codegen(ast)
		actual := disasm(fn)
		if actual != expected {
			t.Errorf("\n%s\nexpect: %s\nactual: %s", src, expected, actual)
		}
	}
	test("true", "True")
	test("", "")
	test("return", "")
	test("return true", "True")
	test("true", "True")
	test("123", "Int 123")
	test("a", "Load a")
	test("_a", "Dyload _a")
	test("G", "Global G")
	test("this", "This")

	test("-a", "Load a, UnaryMinus")
	test("a + b", "Load a, Load b, Add")
	test("a - b", "Load a, Load b, Sub")
	test("a + b + c", "Load a, Load b, Add, Load c, Add")
	test("a + b - c", "Load a, Load b, Add, Load c, Sub")
	test("a - b - c", "Load a, Load b, Sub, Load c, Sub")

	test("a * b", "Load a, Load b, Mul")
	test("a / b", "Load a, Load b, Div")
	test("a * b * c", "Load a, Load b, Mul, Load c, Mul")
	test("a * b / c", "Load a, Load b, Mul, Load c, Div")
	test("a / b / c", "Load a, Load b, Load c, Mul, Div")
	test("a * b / c / d", "Load a, Load b, Mul, Load c, Load d, Mul, Div")

	test("a % b", "Load a, Load b, Mod")
	test("a % b % c", "Load a, Load b, Mod, Load c, Mod")

	test("a | b | c", "Load a, Load b, BitOr, Load c, BitOr")

	test("a is true", "Load a, True, Is")
	test("s = 'hello'", "Value 'hello', Store s")
	test("_dyn = 123", "Int 123, Store _dyn")
	test("a = b = c", "Load c, Store b, Store a")
	test("a = true; not a", "True, Store a, Pop, Load a, Not")
	test("n += 5", "Load n, Int 5, Add, Store n")
	test("++n", "Load n, One, Add, Store n")
	test("n--", "Load n, Dup, One, Sub, Store n, Pop")
	test("a.b", "Load a, Value 'b', Get")
	test("a[2]", "Load a, Int 2, Get")
	test("a.b = 123", "Load a, Value 'b', Int 123, Put")
	test("a[2] = false", "Load a, Int 2, False, Put")
	test("a.b += 5", "Load a, Value 'b', Dup2, Get, Int 5, Add, Put")
	test("++a.b", "Load a, Value 'b', Dup2, Get, One, Add, Put")
	test("a.b++", "Load a, Value 'b', Dup2, Get, Dupx2, One, Add, Put, Pop")
	test("a[..]", "Load a, Zero, MaxInt, RangeTo")
	test("a[..3]", "Load a, Zero, Int 3, RangeTo")
	test("a[2..]", "Load a, Int 2, MaxInt, RangeTo")
	test("a[2..3]", "Load a, Int 2, Int 3, RangeTo")
	test("a[::]", "Load a, Zero, MaxInt, RangeLen")
	test("a[::3]", "Load a, Zero, Int 3, RangeLen")
	test("a[2::]", "Load a, Int 2, MaxInt, RangeLen")
	test("a[2::3]", "Load a, Int 2, Int 3, RangeLen")

	test("return", "")
	test("return 123", "Int 123")

	test("throw 'fubar'", "Value 'fubar', Throw")

	test("f()", "Load f, CallFunc ()")
	test("F()", "Global F, CallFunc ()")
	test("f(a, b)", "Load a, Load b, Load f, CallFunc (?, ?)")
	test("f(1,2,3,4)", "One, Int 2, Int 3, Int 4, Load f, CallFunc (?, ?, ?, ?)")
	test("f(1,2,3,4,5)", "One, Int 2, Int 3, Int 4, Int 5, Load f, CallFunc (?, ?, ?, ?, ?)")
	test("f(a, b, c:, d: 0)", "Load a, Load b, True, Zero, Load f, CallFunc (?, ?, c:, d:)")
	test("f(@args)", "Load args, Load f, CallFunc (@)")
	test("f(@+1args)", "Load args, Load f, CallFunc (@+1)")
	test("f(a: a)", "Load a, Load f, CallFunc (a:)")
	test("f(:a)", "Load a, Load f, CallFunc (a:)")
	test("f(12, 34: 56, false:)",
		"Int 12, Int 56, True, Load f, CallFunc (?, 34:, false:)")
	test("f(1,a:2); f(3,a:4)",
		"One, Int 2, Load f, CallFunc (?, a:), Pop, Int 3, Int 4, Load f, CallFunc (?, a:)")

	test("[a: 2, :b]", "Int 2, Load b, Global Record, CallFunc (a:, b:)")
	test("[1, a: 2, :b]", "One, Int 2, Load b, Global Object, CallFunc (?, a:, b:)")

	test("char.Size()", "Load char, Value 'Size', CallMeth ()")
	test("a.f(123)", "Load a, Int 123, Value 'f', CallMeth (?)")
	test("a.f(1,2,3)", "Load a, One, Int 2, Int 3, Value 'f', CallMeth (?, ?, ?)")
	test("a.f(1,2,3,4)", "Load a, One, Int 2, Int 3, Int 4, Value 'f', CallMeth (?, ?, ?, ?)")
	test("a.f(x:)", "Load a, True, Value 'f', CallMeth (x:)")
	test("a[b](123)", "Load a, Int 123, Load b, CallMeth (?)")
	test("a[b $ c](123)", "Load a, Int 123, Load b, Load c, Cat, CallMeth (?)")
	test("a().Add(123)", "Load a, CallFunc (), Int 123, Value 'Add', CallMeth (?)")
	test("a().Add(123).Size()",
		"Load a, CallFunc (), Int 123, Value 'Add', CallMeth (?), Value 'Size', CallMeth ()")
	test("a.b(1).c(2)",
		"Load a, One, Value 'b', CallMeth (?), Int 2, Value 'c', CallMeth (?)")

	test("function () { }", "Value /* function */")

	test("new c", "Load c, Value '*new*', CallMeth ()")
	test("new c()", "Load c, Value '*new*', CallMeth ()")
	test("new c(1)", "Load c, One, Value '*new*', CallMeth (?)")
}

func TestCodegenSuper(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		c := Constant("Foo { " + src + " }")
		m := src[0:strings.IndexByte(src, '(')]
		fn := c.Lookup(nil, m).(*SuFunc)
		actual := disasm(fn)
		if actual != expected {
			t.Errorf("\n%s\nexpect: %s\nactual: %s", src, expected, actual)
		}
	}
	test("New(){}", "This, Value 'New', Super Foo, CallMeth ()")

	// Super(...) => Super.New(...)
	test("New(){super(1)}", "This, One, Value 'New', Super Foo, CallMeth (?)")

	test("F(){super.Bar(0,1)}", "This, Zero, One, Value 'Bar', Super Foo, CallMeth (?, ?)")
}

func disasm(fn *SuFunc) string {
	da := []string{}
	var s string
	for i := 0; i < len(fn.Code); {
		i, s = Disasm1(fn, i)
		da = append(da, s)
	}
	return strings.Join(da, ", ")
}

func TestControl(t *testing.T) {
	asBlock := false
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		ast := parseFunction("function () {\n" + src + "\n}")
		fn := codegen2(ast, asBlock)
		buf := strings.Builder{}
		Disasm(&buf, fn)
		s := buf.String()
		Assert(t).That(s, Like(expected).Comment(src))
	}

	test(`try F()`, `
		0: Try 13 ''
        4: Global F
        7: CallFunc ()
        9: Pop
        10: Catch 14
        13: Pop`)
	test(`try F() catch G()`, `
		0: Try 13 ''
        4: Global F
        7: CallFunc ()
        9: Pop
        10: Catch 20
        13: Pop
        14: Global G
        17: CallFunc ()
        19: Pop`)
	test(`try F() catch (x, "y") G()`, `
		0: Try 13 'y'
        4: Global F
        7: CallFunc ()
        9: Pop
        10: Catch 22
        13: Store x
        15: Pop
        16: Global G
        19: CallFunc ()
        21: Pop`)

	test("a and b", `
		0: Load a
		2: And 8
		5: Load b
		7: Bool`)
	test("a or b", `
		0: Load a
		2: Or 8
		5: Load b
		7: Bool`)
	test("a or b or c", `
		0: Load a
		2: Or 13
		5: Load b
		7: Or 13
		10: Load c
		12: Bool`)
	test("a is b or c < d", `
		0: Load a
        2: Load b
        4: Is
        5: Or 13
        8: Load c
        10: Load d
        12: Lt`) // no Bool needed

	test("a ? b : c", `
		0: Load a
		2: QMark 10
		5: Load b
		7: Jump 12
		10: Load c`)

	test("a in (4,5,6)", `
		0: Load a
        2: Int 4
        5: In 18
        8: Int 5
        11: In 18
        14: Int 6
        17: Is`)

	test("while (a) b", `
		0: Jump 6
		3: Load b
		5: Pop
		6: Load a
		8: JumpTrue 3`)
	test("while a\n;", `
		0: Jump 3
		3: Load a
		5: JumpTrue 3`)

	test("if (a) b", `
		0: Load a
		2: JumpFalse 8
		5: Load b
		7: Pop`)
	test("if (a) b else c", `
		0: Load a
		2: JumpFalse 11
		5: Load b
		7: Pop
		8: Jump 14
		11: Load c
		13: Pop`)

	test("switch { case 1: b }", `
		0: True
        1: One
        2: JumpIsnt 11
        5: Load b
        7: Pop
        8: Jump 15
        11: Pop
        12: Value 'unhandled switch value'
        14: Throw`)
	test("switch a { case 1,2: b case 3: c default: d }", `
		0: Load a
        2: One
        3: JumpIs 12
        6: Int 2
        9: JumpIsnt 18
        12: Load b
        14: Pop
        15: Jump 34
        18: Int 3
        21: JumpIsnt 30
        24: Load c
        26: Pop
        27: Jump 34
        30: Pop
        31: Load d
        33: Pop`)

	test("forever { break }", `
		0: Jump 6
		3: Jump 0`)

	test("for(;;) { break }", `
		0: Jump 6
		3: Jump 0`)

	test("while a { b; break; continue }", `
		0: Jump 12
		3: Load b
		5: Pop
		6: Jump 17
		9: Jump 0
		12: Load a
		14: JumpTrue 3`)

	test("do a while b", `
		0: Load a
		2: Pop
		3: Load b
		5: JumpTrue 0`)

	test("for (;a;) { b; break; continue }", `
		0: Jump 12
		3: Load b
		5: Pop
		6: Jump 17
		9: Jump 0
		12: Load a
		14: JumpTrue 3`)

	test("for (i = 0; i < 9; ++i) body", `
		0: Zero
        1: Store i
        3: Pop
        4: Jump 17
        7: Load body
        9: Pop
        10: Load i
        12: One
        13: Add
        14: Store i
        16: Pop
        17: Load i
        19: Int 9
        22: Lt
        23: JumpTrue 7`)

	test(`for (x in y) { a; break; continue }`, `
		0: Load y
        2: Iter
        3: ForIn x 19
        7: Load a
        9: Pop
        10: Jump 19
        13: Jump 3
        16: Jump 3
        19: Pop`)

	asBlock = true
	test(`break`, `
		0: BlockBreak`)
	test(`continue`, `
		0: BlockContinue`)
	test(`return`, `
		0: BlockReturnNil`)
	test(`return true`, `
		0: True
		1: BlockReturn`)
}

func TestBlock(t *testing.T) {
	ast := parseFunction("function (x) {\n b = {|a| a + x }\n}")
	fn := codegen(ast)
	block := fn.Values[0].(*SuFunc)

	Assert(t).That(fn.Names, Equals([]string{"x", "b", ""}))
	Assert(t).That(block.Names, Equals([]string{"x", "b", "a"}))
	Assert(t).That(int(block.Offset), Equals(2))

	Assert(t).That(block.ParamSpec.Params(), Equals("(a)"))

	Assert(t).That(disasm(fn), Like(
		`Block
		0: Load a
		2: Load x
		4: Add, Store b`))
}

// parseFunction parses a function and returns an AST for it
func parseFunction(src string) *ast.Function {
	p := NewParser(src)
	return p.Function()
}
