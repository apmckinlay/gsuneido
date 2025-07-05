// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package compile

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestCodegen(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		classNum.Store(0)
		ast := parseFunction("function () {\n" + src + "\n}")
		fn := codegen("", "", ast, nil).(*SuFunc)
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
	test("a + b", "LoadLoad a b, Add")
	test("a - b", "LoadLoad a b, Sub")
	test("a + b + c", "LoadLoad a b, Add, Load c, Add")
	test("a + b - c", "LoadLoad a b, Add, Load c, Sub")
	test("a - b - c", "LoadLoad a b, Sub, Load c, Sub")

	test("a * b", "LoadLoad a b, Mul")
	test("a / b", "LoadLoad a b, Div")
	test("a * b * c", "LoadLoad a b, Mul, Load c, Mul")
	test("a * b / c", "LoadLoad a b, Mul, Load c, Div")
	test("a / b / c", "LoadLoad a b, Load c, Mul, Div")
	test("a * b / c / d", "LoadLoad a b, Mul, LoadLoad c d, Mul, Div")
	test("1 / a", "One, Load a, Div")

	test("a % b", "LoadLoad a b, Mod")
	test("a % b % c", "LoadLoad a b, Mod, Load c, Mod")

	test("a | b | c", "LoadLoad a b, BitOr, Load c, BitOr")

	test("a is b", "LoadLoad a b, Is")
	test("a = b", "Load b, Store a")
	test("a,b = f()",
		"Load f, CallFuncNilOk (), PushReturn 2, StorePop a, StorePop b")
	test("_dyn = 123", "Int 123, Store _dyn")
	test("a = b = c", "Load c, Store b, Store a")
	test("a = b; not a", "Load b, StorePop a, Load a, Not")
	test("n += 5", "Int 5, LoadStore n AddEq")
	test("n /= 5", "Int 5, LoadStore n DivEq")
	test("s =~ '^a'", "LoadValue s regex, Match")
	test("++n", "One, LoadStore n AddEq")
	test("n--", "One, LoadStore n SubEq retOrig")
	test("a.b", "LoadValue a 'b', Get")
	test("a[2]", "Load a, Int 2, Get")
	test("a.b = 123", "LoadValue a 'b', Int 123, Put")
	test("a[2] = false", "Load a, Int 2, False, Put")
	test("a.b += 5", "LoadValue a 'b', Int 5, GetPut AddEq")
	test("a[b] -= 5", "LoadLoad a b, Int 5, GetPut SubEq")
	test("++a.b", "LoadValue a 'b', One, GetPut AddEq")
	test("a[b]--", "LoadLoad a b, One, GetPut SubEq retOrig")
	test("a[..]", "Load a, Zero, MaxInt, RangeTo")
	test("a[..3]", "Load a, Zero, Int 3, RangeTo")
	test("a[2..]", "Load a, Int 2, MaxInt, RangeTo")
	test("a[2..3]", "Load a, Int 2, Int 3, RangeTo")
	test("a[::]", "Load a, Zero, MaxInt, RangeLen")
	test("a[::3]", "Load a, Zero, Int 3, RangeLen")
	test("a[2::]", "Load a, Int 2, MaxInt, RangeLen")
	test("a[2::3]", "Load a, Int 2, Int 3, RangeLen")
	test("0 < a and a < 5", "Load a, InRange Gt 0 Lt 5")
	test("A['b']", "Global A, ValueGet 'b'")
	test("this[a]", "ThisLoad a, Get")
	test("f(a.b, 'c')",
		"LoadValue a 'b', GetValue 'c', Load f, CallFuncNilOk (?, ?)")
	test("a; b", "Load a, PopLoad b")
	test("a = F()", "GlobalCallFuncNoNil F (), Store a")

	test("return", "")
	test("return 123", "Int 123")

	test("return throw 123; 123", "Int 123, ReturnThrow, Int 123")
	test("return throw 123", "Int 123, ReturnThrow")

	test("return 1, 2, 3", "One, Int 2, Int 3, ReturnMulti 3")

	test("throw 'fubar'", "Value 'fubar', Throw")

	test("f()", "Load f, CallFuncNilOk ()")
	test("(f())", "Load f, CallFuncNilOk ()")
	test("f(); f()", "Load f, CallFuncDiscard (), Load f, CallFuncNilOk ()")
	test("F()", "Global F, CallFuncNilOk ()")
	test("f(a, b)", "LoadLoad a b, Load f, CallFuncNilOk (?, ?)")
	test("f(1,2,3,4)", "One, Int 2, Int 3, Int 4, Load f, CallFuncNilOk (?, ?, ?, ?)")
	test("f(1,2,3,4,5)", "One, Int 2, Int 3, Int 4, Int 5, Load f, CallFuncNilOk (?, ?, ?, ?, ?)")
	test("f(a, b, c:, d: 0)", "LoadLoad a b, True, Zero, Load f, CallFuncNilOk (?, ?, c:, d:)")
	test("f(@args)", "LoadLoad args f, CallFuncNilOk (@)")
	test("f(@+1args)", "LoadLoad args f, CallFuncNilOk (@+1)")
	test("f(a: a)", "LoadLoad a f, CallFuncNilOk (a:)")
	test("f(:a)", "LoadLoad a f, CallFuncNilOk (a:)")
	test("f(12, 34: 56, false:)",
		"Int 12, Int 56, True, Load f, CallFuncNilOk (?, 34:, false:)")
	test("f(1,a:2); f(3,a:4)",
		"One, Int 2, Load f, CallFuncDiscard (?, a:), Int 3, Int 4, Load f, CallFuncNilOk (?, a:)")

	test("[a: 2, :b]", "Int 2, Load b, Global Record, CallFuncNilOk (a:, b:)")
	test("[1, a: 2, :b]", "One, Int 2, Load b, Global Object, CallFuncNilOk (?, a:, b:)")

	test("char.Size()", "LoadValue char 'Size', CallMethNilOk ()")
	test("a.f(123)", "Load a, Int 123, Value 'f', CallMethNilOk (?)")
	test("a.f(1,2,3)", "Load a, One, Int 2, Int 3, Value 'f', CallMethNilOk (?, ?, ?)")
	test("a.f(1,2,3,4)", "Load a, One, Int 2, Int 3, Int 4, Value 'f', CallMethNilOk (?, ?, ?, ?)")
	test("a.f(x:)", "Load a, True, Value 'f', CallMethNilOk (x:)")
	test("a[b](123)", "Load a, Int 123, Load b, CallMethNilOk (?)")
	test("a[b $ c](123)", "Load a, Int 123, LoadLoad b c, Cat, CallMethNilOk (?)")
	test("a().Add(123)", "Load a, CallFuncNoNil (), Int 123, Value 'Add', CallMethNilOk (?)")
	test("a().Add(123).Size()",
		"Load a, CallFuncNoNil (), Int 123, ValueCallMethNoNil 'Add' (?), Value 'Size', CallMethNilOk ()")
	test("a.b(1).c(2)",
		"Load a, One, ValueCallMethNoNil 'b' (?), Int 2, Value 'c', CallMethNilOk (?)")

	test("function () { }", "Value /* function */")

	test("new c", "LoadValue c '*new*', CallMethNilOk ()")
	test("new c()", "LoadValue c '*new*', CallMethNilOk ()")
	test("new c(1)", "Load c, One, Value '*new*', CallMethNilOk (?)")

	xtest := func(src, expected string) {
		t.Helper()
		classNum.Store(0)
		ast := parseFunction("function () {\n" + src + "\n}")
		assert.This(func() { codegen("", "", ast, nil) }).Panics(expected)
	}
	xtest("b = { return 1, 2, 3}", "not allowed")
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
		assert.T(t).This(actual).Like(expected)
		// if actual != expected {
		// 	t.Errorf("\n%s\nexpect: %s\nactual: %s", src, expected, actual)
		// }
	}
	test("New(){}", "ThisValue 'New', Super Foo, CallMethNilOk ()")
	test("New(){ F() }", "ThisValue 'New', Super Foo, CallMethDiscard (), "+
		"Global F, CallFuncNilOk ()")

	// Super(...) => Super.New(...)
	test("New(){super(1)}", "This, One, Value 'New', Super Foo, CallMethNilOk (?)")

	test("F(){super.Bar(0,1)}", "This, Zero, One, Value 'Bar', Super Foo, CallMethNilOk (?, ?)")

	test("F() { 1.Times() { super.Push(123) } }",
		"One, Closure, "+
			"{This, Int 123, Value 'Push', Super Foo, CallMethNilOk (?), Value 'Times'}, "+
			"CallMethNilOk (block:)")
}

func disasm(fn *SuFunc) string {
	var ops []string
	nestPrev := 0
	Disasm(fn, func(_ *SuFunc, nest, i int, s string, _ int) {
		if nest > nestPrev {
			s = "{" + s
		} else if nest < nestPrev {
			s += "}"
		}
		nestPrev = nest
		ops = append(ops, s)
	})
	return str.Join(", ", ops)
}

func TestControl(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		ast := parseFunction("function () {\n" + src + "\n}")
		fn := codegen("", "", ast, nil).(*SuFunc)
		s := DisasmOps(fn)
		assert.T(t).Msg(src).This(s).Like(expected)
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
        9: Catch 19
        12: StorePop x
        14: Global G
        17: CallFuncDiscard ()`)

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
		0: LoadLoad a b
        3: Is
        4: Or 11
        7: LoadLoad c d
        10: Lt`) // no Bool needed

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
        15: Jump 33
        18: Int 3
        21: JumpIsnt 30
        24: Load c
        26: Pop
        27: Jump 33
        30: PopLoad d
        32: Pop`)
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

	test("do { a; continue; b } while c", `
		0: Load a
		2: Pop
		3: Jump 9
		6: Load b
		8: Pop
		9: Load c
		11: JumpTrue 0`)

	test("do { a; break; b } while c", `
		0: Load a
		2: Pop
		3: Jump 14
		6: Load b
		8: Pop
		9: Load c
		11: JumpTrue 0`)

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
		1: StorePop i
		3: Jump 14
		6: Load body
		8: Pop
		9: One
		10: LoadStore i AddEq
		13: Pop
		14: Load i
		16: Int 9
		19: Lt
		20: JumpTrue 6`)

	test("for (i = 0; i < 9; ++i) { a; continue; b }", `
		0: Zero
		1: StorePop i
		3: Jump 20
		6: Load a
		8: Pop
		9: Jump 15
		12: Load b
		14: Pop
		15: One
		16: LoadStore i AddEq
		19: Pop
		20: Load i
		22: Int 9
		25: Lt
		26: JumpTrue 6`)

	test(`for (x in y) { a; break; b; continue; c }`, `
		0: Load y
		2: Iter
		3: Jump 21
		6: Load a
		8: Pop
		9: Jump 25
		12: Load b
		14: Pop
		15: Jump 21
		18: Load c
		20: Pop
		21: ForIn x 6
		25: Pop`)

	test("for m,v in ob \n { a; break; b; continue; c }", `
		0: Load ob
		2: Iter2
		3: Jump 21
		6: Load a
		8: Pop
		9: Jump 26
		12: Load b
		14: Pop
		15: Jump 21
		18: Load c
		20: Pop
		21: ForIn2 m v 6
		26: Pop`)

	test("for m,v in ob \n Print(m, v)", `
		0: Load ob
		2: Iter2
		3: Jump 14
		6: LoadLoad m v
		9: Global Print
		12: CallFuncDiscard (?, ?)
		14: ForIn2 m v 6
		19: Pop`)

	test(`for ..10 { a; break; b; continue; c }`, `
		0: Int 10
		3: MinusOne
		4: Jump 22
		7: Load a
		9: Pop
		10: Jump 25
		13: Load b
		15: Pop
		16: Jump 22
		19: Load c
		21: Pop
		22: ForRange 7
		25: Pop
		26: Pop`)

	// same as for (i = 0; i < 10; ++i)
	test(`for i in 0..10 { a; break; b; continue; c }`, `
		0: Int 10
		3: MinusOne
		4: Jump 22
		7: Load a
		9: Pop
		10: Jump 26
		13: Load b
		15: Pop
		16: Jump 22
		19: Load c
		21: Pop
		22: ForRangeVar i 7
		26: Pop
		27: Pop`)
}

func TestBlock(t *testing.T) {
	assert := assert.T(t).This
	ast := parseFunction("function (x) {\n b = {|a| a + x }\n}")
	fn := codegen("", "", ast, nil).(*SuFunc)
	block := fn.Values[0].(*SuFunc)

	assert(fn.Names).Is([]string{"x", "b", "a|2"})
	assert(block.Names).Is([]string{"x", "b", "a"})
	assert(int(block.Offset)).Is(2)

	assert(block.ParamSpec.Params()).Is("(a)")

	assert(disasm(fn)).Is("Closure, {LoadLoad a x, Add, Store b}")
}

// parseFunction parses a function and returns an AST for it
func parseFunction(src string) *ast.Function {
	p := NewParser(src)
	return p.Function()
}

// func TestBug(t *testing.T) {
// 	assert.T(t).
// 		This(func() { parseFunction("function (x) { a = function() { _ } }") }).
// 		Panics("syntax error @32 invalid identifier: '_'")
// }
