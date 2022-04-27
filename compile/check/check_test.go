// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package check_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/check"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCheckVars(t *testing.T) {
	test := func(src, initExp, usedExp string) {
		t.Helper()
		// fmt.Println(src)
		p := compile.NewParser(src)
		ast := p.Function()
		ck := check.Check{}
		init := ck.Check(ast)
		sort.Strings(init)
		assert.T(t).This(fmt.Sprint(init)).Is("[" + initExp + "]")

		var used []string
		for s := range ck.AllUsed {
			used = append(used, s)
		}
		sort.Strings(used)
		assert.T(t).This(fmt.Sprint(used)).Is("[" + usedExp + "]")
	}
	test("function (a,b,c) { }", "a b c", "")
	test("function (a,b) { c = 1; d = 2 }", "a b c d", "")
	test("function () { a; b; c }", "", "a b c")
	test("function () { a = b + c }", "a", "b c")
	test("function () { a = F(b,c) }", "a", "b c")
	test("function () { return a + b }", "", "a b")
	test("function () { throw a $ b }", "", "a b")
	test("function () { try a = b; catch (e) c = d }", "a e", "b d")
	test("function () { if a { b=c; b } else { d=e; d; b } }", "", "a b c d e")
	test("function () { (a=1) ? (b=c) : (d=e) }", "a", "c e")
	test("function (a) { if a { b=1 } else {b=2 } }", "a b", "a")
	test("function (a,b) { a < b }", "a b", "a b")
	test("function () { while (false isnt x = Next()) { } }", "x", "")
}

func TestCheckResults(t *testing.T) {
	test := func(src string, expected ...string) {
		t.Helper()
		// fmt.Println(src)
		_, results := compile.Checked(nil, src)
		assert.T(t).This(results).Is(expected)
	}

	test("function () { return { it } }")

	test("function (a) { }",
		"WARNING: initialized but not used: a @10")
	test("function (_a) { }",
		"WARNING: initialized but not used: a @10")
	test("function (@a) { }",
		"WARNING: initialized but not used: a @10")
	test("function (unused) { }")
	test("function (a/*unused*/) { }")
	test("function (_a/*unused*/) { }")
	test("function (@a/*unused*/) { }")
	test("function (a/*unused*/) { a }",
		"ERROR: used but not initialized: a @25")
	test("function (@a/*unused*/) { a }",
		"ERROR: used but not initialized: a @26")
	test("function () { a=1 }",
		"WARNING: initialized but not used: a @14")
	test("function () { a=b; a }",
		"ERROR: used but not initialized: b @16")
	test("function () { a++ }",
		"ERROR: used but not initialized: a @14")
	test("function () { a += 1 }",
		"ERROR: used but not initialized: a @14")
	test("function () { a=1+a; a=2 }",
		"ERROR: used but not initialized: a @18")
	test("function () { a + a }",
		"ERROR: used but not initialized: a @14",
		"ERROR: used but not initialized: a @18")
	test("function () { a }",
		"ERROR: used but not initialized: a @14")

	// if and ?:
	test("function (a) { if a { b() } }",
		"ERROR: used but not initialized: b @22")
	test("function (a) { if a { b=5 } }",
		"WARNING: initialized but not used: b @22")
	test("function (a) { if a { b=a } if a { b() } }",
		"WARNING: used but possibly not initialized: b @35")
	test("function (a) { if (a) { b=1 } else { b=2 } b }")
	test("function (a) { a ? b=1 : b=2; b }")
	test("function (a) { a ? b=a : 2; b }",
		"WARNING: used but possibly not initialized: b @28")

	// switch
	test("function (f) { switch (a=f()) { case 1: a(); default: a() } }")
	test("function (f) { switch (a=f()) { case 1: a(); default: a() }; a() }")
	test("function (f) { switch (a=f()) { case 1: b=1; default: b=2 }; a + b }")
	test("function (f) { switch (f()) { case 1: a=f; case 2: b=2;b() }; a }",
		"WARNING: used but possibly not initialized: a @62")

	// while
	test("function () { while (false isnt x = 1) { }; x }")

	// for(;;)
	test("function () { for (i=0; i < 5; j++) { } }",
		"ERROR: used but not initialized: j @31")
	test("function () { for (i=0; i < 5; i++, j++) { j=0 } }")

	// for-in
	test("function () { for x in #() { } }",
		"WARNING: initialized but not used: x @18")

	// try-catch
	test("function () { try {} catch (e) {} }",
		"WARNING: initialized but not used: e @28")
	test("function (f) { try f() catch (e) f() }",
		"WARNING: initialized but not used: e @30")
	test("function (f) { try f() catch (e, 'x') f() }",
		"WARNING: initialized but not used: e @30")
	test("function (f) { try f() catch (unused) f() }")
	test("function (f) { try f() catch (e /*unused*/) f() }")
	test("function (f) { try f() catch (e /*unused*/, 'x') f() }")

	// class
	test("class { F(a){} G(){a} }",
		"WARNING: initialized but not used: a @10",
		"ERROR: used but not initialized: a @19")
	test("class { New(.X){} }")

	// blocks
	test("function () { return { it } }")
	test("function (f) { f({ x=1 }); x }")
	test("function (f) { f({|unused| }) }")
	test("function (f) { f({|x/*unused*/| }) }")
	test("function (f) { f({|@args| args }) }")
	test("function (f) { f({|x/*unused*/| x }) }",
		"ERROR: used but not initialized: x @32")

	test("function (f) { f(a=1).x = a }",
		"ERROR: used but not initialized: a @26")

	// shadowing
	test("function (f) { f({|x| x }); x }",
		"ERROR: used but not initialized: x @28")
	test("function (f) { f({|x| x++ }); x }",
		"ERROR: used but not initialized: x @30")
	test("function (f) { f({|x| x = x + 1 }); x }",
		"ERROR: used but not initialized: x @36")
	test("function (f) { x=1; f({|x| x }); }",
		"WARNING: initialized but not used: x @15")
	test("function (f, x) { f({|x| x }); }",
		"WARNING: initialized but not used: x @13")
	test("function (f) { a=1; f({|c,d| a }); c }",
		"WARNING: initialized but not used: c @24",
		"WARNING: initialized but not used: d @26",
		"ERROR: used but not initialized: c @35")

	// and/or conditions
	test("function (f) { (f and (b=f)) ? b : 0 }")
	test("function (f) { (f or (b=f)) ? 0 : b }")
	test("function (f) { (f and (b=f)) ? 0 : b }",
		"WARNING: used but possibly not initialized: b @35")
	test("function (f) { (f and (b=f)) ? f() : f(); b }",
		"WARNING: used but possibly not initialized: b @42")
	test("function (f) { if (f and (b=f)) { b() } }")
	test("function (f) { if (f or (b=f)) {} else { b() } }")
	test("function (f) { while (f and (b=f)) { b() } }")
	test("function (f) { while (f and (b=f)) { } b }",
		"WARNING: used but possibly not initialized: b @39")

	// unreachable
	test("function () { return }")
	test("function () { return; 123 }",
		"ERROR: unreachable code @22")
	test("function (x) { if (x) { return; x() } }",
		"ERROR: unreachable code @32")
	test("function (x) { if (x) { return } else { return } x }",
		"ERROR: unreachable code @49")
	test("function (f) { forever { break; f() } }",
		"ERROR: unreachable code @32")
	test("function (f) { forever { f(); continue; f() } }",
		"ERROR: unreachable code @40")
	test("function (f) { switch { } f() }",
		"ERROR: unreachable code @26")
	test("function () { switch { default: } 123 }")
	test("function () { switch { case 0: return } 123 }",
		"ERROR: unreachable code @40")

	// useless expression (no side effects)
	test("function () { 123; return }",
		"ERROR: useless expression @14")
	test("function () { class{}; return }",
		"ERROR: useless expression @14")
	test("function (f) { if (f()) return \n 123 \n return 456 }",
		"ERROR: useless expression @33")
	test("function (x) { try x+0 catch x=0 }")

	// guard clause
	test("function (f) { if (f()) { x=f } else { return } x  }")
	test("function (f) { if (f()) { return } else { x=f } x  }")

	// copy on write
	test("function (a) { if (a) b = 1 else a(c=1, b=c) return b }")
	test("function (a) { a ? b=1 : a(c=1, b=c); b }")

	// missing comma
	runtime.Global.TestDef("Object", runtime.False)
	test(`function (f, x) { f(x[0]) }`)
	test(`function (f, x) { f(x(0)) }`)
	test(`function (f, x) { f(x{}) }`)
	test(`function (f, x) { f(x.y) }`)
	test(`function (f, x) { f(x+1) }`)
	test(`function (f, x) { f(x-1) }`)

	test(`function (f, x) { f(x,
		[0]) }`)
	test(`function (f, x) { f(x,
		(0)) }`)
	test(`function (f, x) { f(x,
		{}) }`)
	test(`function (f, x) { f(x,
		.y) }`)
	test(`function (f, x) { f(x,
		+1) }`)
	test(`function (f, x) { f(x,
		-1) }`)

	test(`function (f, x) { f(x
		[0]) }`, "ERROR: missing comma @24")
	test(`function (f, x) { f(x
		(0)) }`, "ERROR: missing comma @24")
	test(`function (f, x) { f(x
		{}) }`, "ERROR: missing comma @24")
	test(`function (f, x) { f(x
		.y) }`, "ERROR: missing comma @24")
	test(`function (f, x) { f(x
		+1) }`, "ERROR: missing comma @24")
	test(`function (f, x) { f(x
		-1) }`, "ERROR: missing comma @24")
}
