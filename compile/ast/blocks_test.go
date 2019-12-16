// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast_test

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestBlocks(t *testing.T) {
	test := func(src, expected string) {
		t.Helper()
		p := compile.NewParser("function () {\n" + src + "\n}")
		f := p.Function()
		ast.Blocks(f)
		s := f.String()
		Assert(t).That(s[9:len(s)-1], Like(expected))
	}
	test("", "")
	test("a=1; b=2",
		`Binary(Eq a 1)
		Binary(Eq b 2)`)
	test("b = {}",
		`Binary(Eq b Block-func())`)
	test("b = { x }",
		`Binary(Eq b Block-func(
			x))`)
	test("b = { it }",
		`Binary(Eq b Block-func(it
			it))`)
	test("a=1; b = { a }",
		`Binary(Eq a 1)
        Binary(Eq b Block(
        	a))`)
	test("b = { a }; a",
		`Binary(Eq b Block(
			a))
		a`)
	test("a=1; b = {|a| a }",
		`Binary(Eq a 1)
        Binary(Eq b Block-func(a
			a))`)
	test("b1 = { x }; b2 = { x }", // peers
		`Binary(Eq b1 Block(
        	x))
        Binary(Eq b2 Block(
			x))`)
	test("b1 = { b2 = { x } }",
		`Binary(Eq b1 Block-func(
        	Binary(Eq b2 Block-func(
        		x))))`)
	test("b1 = { b2 = { x } }; x", // child
		`Binary(Eq b1 Block(
        	Binary(Eq b2 Block(
        		x))))
		x`)
	test("x=1; b = {|x| x }", // params not sharing
		`Binary(Eq x 1)
        	Binary(Eq b Block-func(x
			x))`)
	test("b = { .x }", // "this" requires closure
		`Binary(Eq b Block(
			Mem(this "x")))`)
	test("b = { super.F() }", // "super" requires closure
		`Binary(Eq b Block(
			Call(Mem(super "F"))))`)
	test("b1 = {|p| b2 = { p }}", // inner references outer param
		`Binary(Eq b1 Block(p
        	Binary(Eq b2 Block(
			p))))`)
	test("_x = 5; b = { _x }", // dynamic
		`Binary(Eq _x 5)
        	Binary(Eq b Block(
        	_x))`)
}
