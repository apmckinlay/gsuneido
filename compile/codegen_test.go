// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"os"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func Example_parseFunction_SrcPos() {
	src := `function (x, y) {
		a = x
		b = y
		return a + b
		}`
	ast := parseFunction(src)
	fn := codegen(ast)
	DisasmMixed(os.Stdout, fn, src)
	// Output:
	// 20: a = x
	// 	0: Load x
	// 	2: Store a
	// 	4: Pop
	// 28: b = y
	// 	5: Load y
	// 	7: Store b
	// 	9: Pop
	// 36: return a + b
	// 	10: Load a
	// 	12: Load b
	// 	14: Add
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
	test("1 / a", "One, Load a, Div")

	test("a % b", "Load a, Load b, Mod")
	test("a % b % c", "Load a, Load b, Mod, Load c, Mod")

	test("a | b | c", "Load a, Load b, BitOr, Load c, BitOr")

	test("a is b", "Load a, Load b, Is")
	test("a = b", "Load b, Store a")
	test("_dyn = 123", "Int 123, Store _dyn")
	test("a = b = c", "Load c, Store b, Store a")
	test("a = b; not a", "Load b, Store a, Pop, Load a, Not")
	test("n += 5", "Int 5, LoadLock n, Add, StoreUnlock")
	test("n /= 5", "Int 5, LoadLock n, Swap, Div, StoreUnlock")
	test("++n", "LoadLock n, One, Add, StoreUnlock")
	test("n--", "LoadLock n, Dup, One, Sub, StoreUnlock, Pop")
	test("a.b", "Load a, Value 'b', Get")
	test("a[2]", "Load a, Int 2, Get")
	test("a.b = 123", "Load a, Value 'b', Int 123, Put")
	test("a[2] = false", "Load a, Int 2, False, Put")
	test("a.b += 5", "Int 5, Load a, Value 'b', GetLock, Add, PutUnlock")
	test("a.b -= 5", "Int 5, Load a, Value 'b', GetLock, Swap, Sub, PutUnlock")
	test("++a.b", "Load a, Value 'b', GetLock, One, Add, PutUnlock")
	test("a.b++", "Load a, Value 'b', GetLock, Dup, One, Add, PutUnlock, Pop")
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

	test("f()", "Load f, CallFuncNilOk ()")
	test("(f())", "Load f, CallFuncNilOk ()")
	test("f(); f()", "Load f, CallFuncDiscard (), Load f, CallFuncNilOk ()")
	test("F()", "Global F, CallFuncNilOk ()")
	test("f(a, b)", "Load a, Load b, Load f, CallFuncNilOk (?, ?)")
	test("f(1,2,3,4)", "One, Int 2, Int 3, Int 4, Load f, CallFuncNilOk (?, ?, ?, ?)")
	test("f(1,2,3,4,5)", "One, Int 2, Int 3, Int 4, Int 5, Load f, CallFuncNilOk (?, ?, ?, ?, ?)")
	test("f(a, b, c:, d: 0)", "Load a, Load b, True, Zero, Load f, CallFuncNilOk (?, ?, c:, d:)")
	test("f(@args)", "Load args, Load f, CallFuncNilOk (@)")
	test("f(@+1args)", "Load args, Load f, CallFuncNilOk (@+1)")
	test("f(a: a)", "Load a, Load f, CallFuncNilOk (a:)")
	test("f(:a)", "Load a, Load f, CallFuncNilOk (a:)")
	test("f(12, 34: 56, false:)",
		"Int 12, Int 56, True, Load f, CallFuncNilOk (?, 34:, false:)")
	test("f(1,a:2); f(3,a:4)",
		"One, Int 2, Load f, CallFuncDiscard (?, a:), Int 3, Int 4, Load f, CallFuncNilOk (?, a:)")

	test("[a: 2, :b]", "Int 2, Load b, Global Record, CallFuncNilOk (a:, b:)")
	test("[1, a: 2, :b]", "One, Int 2, Load b, Global Object, CallFuncNilOk (?, a:, b:)")

	test("char.Size()", "Load char, Value 'Size', CallMethNilOk ()")
	test("a.f(123)", "Load a, Int 123, Value 'f', CallMethNilOk (?)")
	test("a.f(1,2,3)", "Load a, One, Int 2, Int 3, Value 'f', CallMethNilOk (?, ?, ?)")
	test("a.f(1,2,3,4)", "Load a, One, Int 2, Int 3, Int 4, Value 'f', CallMethNilOk (?, ?, ?, ?)")
	test("a.f(x:)", "Load a, True, Value 'f', CallMethNilOk (x:)")
	test("a[b](123)", "Load a, Int 123, Load b, CallMethNilOk (?)")
	test("a[b $ c](123)", "Load a, Int 123, Load b, Load c, Cat, CallMethNilOk (?)")
	test("a().Add(123)", "Load a, CallFuncNoNil (), Int 123, Value 'Add', CallMethNilOk (?)")
	test("a().Add(123).Size()",
		"Load a, CallFuncNoNil (), Int 123, Value 'Add', CallMethNoNil (?), Value 'Size', CallMethNilOk ()")
	test("a.b(1).c(2)",
		"Load a, One, Value 'b', CallMethNoNil (?), Int 2, Value 'c', CallMethNilOk (?)")

	test("function () { }", "Value /* function */")

	test("new c", "Load c, Value '*new*', CallMethNilOk ()")
	test("new c()", "Load c, Value '*new*', CallMethNilOk ()")
	test("new c(1)", "Load c, One, Value '*new*', CallMethNilOk (?)")
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
		Assert(t).That(actual, Like(expected))
		// if actual != expected {
		// 	t.Errorf("\n%s\nexpect: %s\nactual: %s", src, expected, actual)
		// }
	}
	test("New(){}", "This, Value 'New', Super Foo, CallMethNilOk ()")
	test("New(){ F() }", "This, Value 'New', Super Foo, CallMethDiscard (), "+
		"Global F, CallFuncNilOk ()")

	// Super(...) => Super.New(...)
	test("New(){super(1)}", "This, One, Value 'New', Super Foo, CallMethNilOk (?)")

	test("F(){super.Bar(0,1)}", "This, Zero, One, Value 'Bar', Super Foo, CallMethNilOk (?, ?)")

	test("F() { 1.Times() { super.Push(123) } }", `One, Closure
		0: This
		1: Int 123
		4: Value 'Push'
		6: Super Foo
		9: CallMethNilOk (?), Value 'Times', CallMethNilOk (block:)`)
}

func disasm(fn *SuFunc) string {
	da := []string{}
	var s string
	for i := 0; i < len(fn.Code); {
		i, s = Disasm1(fn, i)
		// s = str.BeforeFirst(s, "\n")
		da = append(da, s)
	}
	return strings.Join(da, ", ")
}

func TestControl(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		ast := parseFunction("function () {\n" + src + "\n}")
		fn := codegen(ast)
		buf := strings.Builder{}
		Disasm(&buf, fn)
		s := buf.String()
		Assert(t).That(s, Like(expected).Comment(src))
	}

	test(`try F()`, `
		0: Try 12 ''
        4: Global F
        7: CallFuncDiscard ()
        9: Catch 13
        12: Pop`)
	test(`try F() catch G()`, `
		0: Try 12 ''
        4: Global F
        7: CallFuncDiscard ()
        9: Catch 18
        12: Pop
        13: Global G
        16: CallFuncDiscard ()`)
	test(`try F() catch (x, "y") G()`, `
		0: Try 12 'y'
        4: Global F
        7: CallFuncDiscard ()
        9: Catch 20
        12: Store x
        14: Pop
        15: Global G
        18: CallFuncDiscard ()`)

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

	test("a ? b() : c()", `
		0: Load a
        2: QMark 12
        5: Load b
        7: CallFuncNilOk ()
        9: Jump 16
        12: Load c
        14: CallFuncNilOk ()`)

	test("a ? b() : c();;", `
		0: Load a
        2: QMark 12
        5: Load b
        7: CallFuncNilOk ()
        9: Jump 16
        12: Load c
		14: CallFuncNilOk ()
		16: Pop`)

	test("(a ? b : c)", `
		0: Load a
        2: QMark 10
        5: Load b
        7: Jump 12
        10: Load c`)

	test("a ? b : c;;", `
		0: Load a
        2: QMark 10
        5: Load b
        7: Jump 12
        10: Load c
        12: Pop`)

	test("(a ? b : c);;", `
		0: Load a
        2: QMark 10
        5: Load b
        7: Jump 12
        10: Load c
        12: Pop`)

	test("return a ? b : c", `
		0: Load a
        2: QMark 10
        5: Load b
        7: Jump 12
        10: Load c`)

	test("return a ? b : c;;", `
		0: Load a
        2: QMark 10
        5: Load b
        7: Jump 12
        10: Load c
        12: Return`)

	test("return (a ? b : c);;", `
		0: Load a
        2: QMark 10
        5: Load b
        7: Jump 12
        10: Load c
        12: Return`)

	test("a in ()", `
		0: False`)

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
	test("x=1; if (x > 1) a()",
		``)

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
	test("switch a { case 1,2,3: b default: }", `
		0: Load a
        2: One
        3: JumpIs 18
        6: Int 2
        9: JumpIs 18
        12: Int 3
        15: JumpIsnt 24
        18: Load b
        20: Pop
        21: Jump 25
        24: Pop`)

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
		9: Jump 12
		12: Load a
		14: JumpTrue 3`)

	test("for (i = 0; i < 9; ++i) body", `
		0: Zero
        1: Store i
        3: Pop
        4: Jump 16
        7: Load body
        9: Pop
        10: LoadLock i
        12: One
        13: Add
        14: StoreUnlock
        15: Pop
        16: Load i
        18: Int 9
        21: Lt
        22: JumpTrue 7`)

	test("for (i = 0; i < 9; ++i) { a; continue; b }", `
		0: Zero
        1: Store i
        3: Pop
        4: Jump 22
        7: Load a
        9: Pop
        10: Jump 16
        13: Load b
        15: Pop
        16: LoadLock i
        18: One
        19: Add
        20: StoreUnlock
        21: Pop
        22: Load i
        24: Int 9
        27: Lt
        28: JumpTrue 7`)

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
}

func TestBlock(t *testing.T) {
	ast := parseFunction("function (x) {\n b = {|a| a + x }\n}")
	fn := codegen(ast)
	block := fn.Values[0].(*SuFunc)

	Assert(t).That(fn.Names, Equals([]string{"x", "b", "a|2"}))
	Assert(t).That(block.Names, Equals([]string{"x", "b", "a"}))
	Assert(t).That(int(block.Offset), Equals(2))

	Assert(t).That(block.ParamSpec.Params(), Equals("(a)"))

	Assert(t).That(disasm(fn), Like(
		`Closure
		0: Load a
		2: Load x
		4: Add, Store b`))
}

// parseFunction parses a function and returns an AST for it
func parseFunction(src string) *ast.Function {
	p := NewParser(src)
	return p.Function()
}
