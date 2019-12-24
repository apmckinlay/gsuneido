// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package check_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/check"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
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
		Assert(t).That(fmt.Sprint(init), Equals("["+initExp+"]"))

		var used []string
		for s := range ck.AllUsed {
			used = append(used, s)
		}
		sort.Strings(used)
		Assert(t).That(fmt.Sprint(used), Equals("["+usedExp+"]"))
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
		Assert(t).That(results, Equals(expected))
	}
	test("function (a) { }",
		"WARNING: initialized but not used: a @10")
	test("function (a /*unused*/, b) { }",
		"WARNING: initialized but not used: b @24")
	test("function (a/*unused*/) { a }",
		"ERROR: used but not initialized: a @25")
	test("function () { a=1 }",
		"WARNING: initialized but not used: a @14")
	test("function () { a=b; a }",
		"ERROR: used but not initialized: b @16")
	test("function () { a=1+a }",
		"ERROR: used but not initialized: a @18")
	test("function () { a + a }",
		"ERROR: used but not initialized: a @14",
		"ERROR: used but not initialized: a @18")
	test("function () { a }",
		"ERROR: used but not initialized: a @14")
	test("function (a) { if a { b } }",
		"ERROR: used but not initialized: b @22")
	test("function (a) { if a { b=5 } }",
		"WARNING: initialized but not used: b @22")
	test("function (a) { if a { b=5 } if a { b } }",
		"WARNING: used but possibly not initialized: b @35")
	test("function () { a=1; b={|c,d| a }; c }",
		"ERROR: used but not initialized: c @33",
		"WARNING: initialized but not used: b @19")
	test("function () { while (false isnt x = 1) { }; x }")
	test("function () { for (i=0; i < 5; j++) { } }",
		"ERROR: used but not initialized: j @31")
	test("function () { for x in #() { } }",
		"WARNING: initialized but not used: x @18")
	test("function () { try {} catch (e) {} }",
		"WARNING: initialized but not used: e @28")

	test("class { F(){} G(a){} }",
		"WARNING: initialized but not used: a @16")
	test("class { New(.X){} }")

	test("function () { try true catch (e) false }",
		"WARNING: initialized but not used: e @30")
	test("function () { try true catch (e, 'x') false }",
		"WARNING: initialized but not used: e @30")
	test("function () { try true catch (e /*unused*/) false }")
	test("function () { try true catch (e /*unused*/, 'x') false }")
}
