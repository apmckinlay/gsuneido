package compile

import (
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestVars(t *testing.T) {
	test := func(src string, expected ...string) {
		f := ParseFunction(src)
		vars := ast.VarList(f)
		sort.Strings(vars)
		sort.Strings(expected)
		Assert(t).That(vars, Equals(expected))
	}
	test("function () {}", []string{}...)
	test("function (c,a,b) { X() }", "c", "a", "b")
	test("function () { x = y = z }", "x", "y", "z")
	test("function () { for x in [] {} }", "x")
	test("function () { try f() catch {} }", "f")
	test("function () { try f() catch(x) {} }", "f", "x")
	test("function (a) { function(x){y} }", "a")
	test("function (a) { class { F(x){y} } }", "a")
	test("function () { b = {|x| y } }", "b")
}
