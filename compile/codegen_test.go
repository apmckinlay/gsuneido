package compile

import (
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/runtime"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestCodegen(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		ast := ParseFunction("function () {\n" + src + "\n}")
		fn := codegen(ast)
		actual := disasm(fn)
		if actual != expected {
			t.Errorf("\n%s\nexpect: %s\nactual: %s", src, expected, actual)
		}
	}
	test("true", "true")
	test("", "")
	test("return", "")
	test("return true", "true")
	test("true", "true")
	test("123", "int 123")
	test("a", "load a")
	test("_a", "dyload _a")
	test("G", "global G")
	test("this", "this")

	test("-a", "load a, uminus")
	test("a + b", "load a, load b, add")
	test("a - b", "load a, load b, sub")
	test("a + b + c", "load a, load b, add, load c, add")
	test("a + b - c", "load a, load b, add, load c, sub")
	test("a - b - c", "load a, load b, sub, load c, sub")

	test("a * b", "load a, load b, mul")
	test("a / b", "load a, load b, div")
	test("a * b * c", "load a, load b, mul, load c, mul")
	test("a * b / c", "load a, load b, mul, load c, div")
	test("a / b / c", "load a, load b, div, load c, div")

	test("a % b", "load a, load b, mod")
	test("a % b % c", "load a, load b, mod, load c, mod")

	test("a | b | c", "load a, load b, bitor, load c, bitor")

	test("a is true", "load a, true, is")
	test("s = 'hello'", "value 'hello', store s")
	test("_dyn = 123", "int 123, store _dyn")
	test("a = b = c", "load c, store b, store a")
	test("a = true; not a", "true, store a, pop, load a, not")
	test("n += 5", "load n, int 5, add, store n")
	test("++n", "load n, one, add, store n")
	test("n--", "load n, dup, one, sub, store n, pop")
	test("a.b", "load a, value 'b', get")
	test("a[2]", "load a, int 2, get")
	test("a.b = 123", "load a, value 'b', int 123, put")
	test("a[2] = false", "load a, int 2, false, put")
	test("a.b += 5", "load a, value 'b', dup2, get, int 5, add, put")
	test("++a.b", "load a, value 'b', dup2, get, one, add, put")
	test("a.b++", "load a, value 'b', dup2, get, dupx2, one, add, put, pop")
	test("a[..]", "load a, zero, maxint, rangeto")
	test("a[..3]", "load a, zero, int 3, rangeto")
	test("a[2..]", "load a, int 2, maxint, rangeto")
	test("a[2..3]", "load a, int 2, int 3, rangeto")
	test("a[::]", "load a, zero, maxint, rangelen")
	test("a[::3]", "load a, zero, int 3, rangelen")
	test("a[2::]", "load a, int 2, maxint, rangelen")
	test("a[2::3]", "load a, int 2, int 3, rangelen")

	test("return", "")
	test("return 123", "int 123")

	test("throw 'fubar'", "value 'fubar', throw")

	test("f()", "load f, callfunc0")
	test("F()", "global F, callfunc0")
	test("f(a, b)", "load a, load b, load f, callfunc2")
	test("f(1,2,3,4)", "one, int 2, int 3, int 4, load f, callfunc4")
	test("f(1,2,3,4,5)", "one, int 2, int 3, int 4, int 5, load f, callfunc(?, ?, ?, ?, ?)")
	test("f(a, b, c:, d: 0)", "load a, load b, true, zero, load f, callfunc(?, ?, c:, d:)")
	test("f(@args)", "load args, load f, callfunc(@)")
	test("f(@+1args)", "load args, load f, callfunc(@+1)")
	test("f(a: a)", "load a, load f, callfunc(a:)")
	test("f(:a)", "load a, load f, callfunc(a:)")
	test("f(12, 34: 56, false:)",
		"int 12, int 56, true, load f, callfunc(?, 34:, false:)")

	test("[1, a: 2, :b]", "one, int 2, load b, global Record, callfunc(?, a:, b:)")

	test("char.Size()", "load char, value 'Size', callmeth0")
	test("a.f(123)", "load a, int 123, value 'f', callmeth1")
	test("a.f(1,2,3)", "load a, one, int 2, int 3, value 'f', callmeth3")
	test("a.f(1,2,3,4)", "load a, one, int 2, int 3, int 4, value 'f', callmeth4")
	test("a.f(x:)", "load a, true, value 'f', callmeth(x:)")
	test("a[b](123)", "load a, int 123, load b, callmeth1")
	test("a[b $ c](123)", "load a, int 123, load b, load c, cat, callmeth1")
	test("a().Add(123)", "load a, callfunc0, int 123, value 'Add', callmeth1")
	test("a().Add(123).Size()",
		"load a, callfunc0, int 123, value 'Add', callmeth1, value 'Size', callmeth0")
	test("a.b(1).c(2)",
		"load a, one, value 'b', callmeth1, int 2, value 'c', callmeth1")

	test("function () { }", "value /* function */")

	test("new c", "load c, value '*new*', callmeth0")
	test("new c()", "load c, value '*new*', callmeth0")
	test("new c(1)", "load c, one, value '*new*', callmeth1")
}

func TestCodegenSuper(t *testing.T) {
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		c := Constant("Foo { " + src + " }")
		m := src[0:strings.IndexByte(src, '(')]
		fn := c.Lookup(m).(*SuFunc)
		actual := disasm(fn)
		if actual != expected {
			t.Errorf("\n%s\nexpect: %s\nactual: %s", src, expected, actual)
		}
	}
	test("New(){}", "this, value 'New', super Foo, callmeth0")

	// super(...) => super.New(...)
	test("New(){super(1)}", "this, one, value 'New', super Foo, callmeth1")

	test("F(){super.Bar(0,1)}", "this, zero, one, value 'Bar', super Foo, callmeth2")
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
	DefaultSingleQuotes = true
	defer func() { DefaultSingleQuotes = false }()
	test := func(src, expected string) {
		t.Helper()
		ast := ParseFunction("function () {\n" + src + "\n}")
		fn := codegen(ast)
		buf := strings.Builder{}
		Disasm(&buf, fn)
		s := buf.String()
		Assert(t).That(s, Like(expected).Comment(src))
	}

	test(`try F()`, `
		0: try 13 ''
        4: global F
        7: callfunc0
        9: pop
        10: catch 14
        13: pop
        14:`)
	test(`try F() catch G()`, `
		0: try 13 ''
        4: global F
        7: callfunc0
        9: pop
        10: catch 20
        13: pop
        14: global G
        17: callfunc0
        19: pop
        20:`)
	test(`try F() catch (x, "y") G()`, `
		0: try 13 'y'
        4: global F
        7: callfunc0
        9: pop
        10: catch 22
        13: store x
        15: pop
        16: global G
        19: callfunc0
        21: pop
        22:`)

	test("a and b", `
		0: load a
		2: and 8
		5: load b
		7: bool
		8:`)
	test("a or b", `
		0: load a
		2: or 8
		5: load b
		7: bool
		8:`)
	test("a or b or c", `
		0: load a
		2: or 13
		5: load b
		7: or 13
		10: load c
		12: bool
		13:`)

	test("a ? b : c", `
		0: load a
		2: qmark 10
		5: load b
		7: jump 12
		10: load c
		12:`)

	test("a in (4,5,6)", `
		0: load a
        2: int 4
        5: in 18
        8: int 5
        11: in 18
        14: int 6
        17: is
        18:`)

	test("while (a) b", `
		0: jump 6
		3: load b
		5: pop
		6: load a
		8: tjump 3
		11:`)
	test("while a\n;", `
		0: jump 3
		3: load a
		5: tjump 3
		8:`)

	test("if (a) b", `
		0: load a
		2: fjump 8
		5: load b
		7: pop
		8:`)
	test("if (a) b else c", `
		0: load a
		2: fjump 11
		5: load b
		7: pop
		8: jump 14
		11: load c
		13: pop
		14:`)

	test("switch { case 1: b }", `
		0: true
        1: one
        2: nejump 11
        5: load b
        7: pop
        8: jump 15
        11: pop
        12: value 'unhandled switch value'
        14: throw
        15:`)
	test("switch a { case 1,2: b case 3: c default: d }", `
		0: load a
        2: one
        3: eqjump 12
        6: int 2
        9: nejump 18
        12: load b
        14: pop
        15: jump 34
        18: int 3
        21: nejump 30
        24: load c
        26: pop
        27: jump 34
        30: pop
        31: load d
        33: pop
        34:`)

	test("forever { break }", `
		0: jump 6
		3: jump 0
		6:`)

	test("for(;;) { break }", `
		0: jump 6
		3: jump 0
		6:`)

	test("while a { b; break; continue }", `
		0: jump 12
		3: load b
		5: pop
		6: jump 17
		9: jump 0
		12: load a
		14: tjump 3
		17:`)

	test("do a while b", `
		0: load a
		2: pop
		3: load b
		5: tjump 0
		8:`)

	test("for (;a;) { b; break; continue }", `
		0: jump 12
		3: load b
		5: pop
		6: jump 17
		9: jump 0
		12: load a
		14: tjump 3
		17:`)

	test("for (i = 0; i < 9; ++i) body", `
		0: zero
        1: store i
        3: pop
        4: jump 17
        7: load body
        9: pop
        10: load i
        12: one
        13: add
        14: store i
        16: pop
        17: load i
        19: int 9
        22: lt
        23: tjump 7
		26:`)

	test(`for (x in y) { a; break; continue }`, `
		0: load y
        2: iter
        3: forin x 19
        7: load a
        9: pop
        10: jump 19
        13: jump 3
        16: jump 3
        19: pop
        20:`)
}

func TestBlock(t *testing.T) {
	ast := ParseFunction("function (x) {\n b = {|a| a + x }\n}")
	fn := codegen(ast)
	block := fn.Values[0].(*SuFunc)

	Assert(t).That(fn.Names, Equals([]string{"x", "b", ""}))
	Assert(t).That(block.Names, Equals([]string{"x", "b", "a"}))
	Assert(t).That(int(block.Offset), Equals(2))

	Assert(t).That(block.ParamSpec.Params(), Equals("(a)"))

	Assert(t).That(disasm(fn), Equals("block, store b"))
	Assert(t).That(disasm(block), Equals("load a, load x, add"))
}
