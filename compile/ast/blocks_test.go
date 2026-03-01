// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ast_test

import (
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

const SharedSlotStart = core.SharedSlotStart

func TestBlocks(t *testing.T) {
	test := func(src, expected string) {
		t.Helper()
		p := compile.NewParser("function () {\n" + src + "\n}")
		f := p.Function()
		ast.Blocks(f)
		s := f.String()
		assert.T(t).This(s[9 : len(s)-1]).Like(expected)
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
	test("b1 = { x }; b2 = { x }", // siblings - no longer detected as sharing
		`Binary(Eq b1 Block-func(
        	x))
        Binary(Eq b2 Block-func(
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
			Mem(this 'x')))`)
	test("b = { super.F() }", // "super" requires closure
		`Binary(Eq b Block(
			Call(Mem(super 'F'))))`)
	test("b1 = {|p| b2 = { p }}", // inner references outer param
		`Binary(Eq b1 Block(p
        	Binary(Eq b2 Block(
			p))))`)
	test("_x = 5; b = { _x }", // dynamic not shared
		`Binary(Eq _x 5)
	        	Binary(Eq b Block-func(
	        	_x))`)

	// block return requires closure
	test("b = { return }",
		`Binary(Eq b Block(
        	Return()))`)
	test("b1 = { b2 = { return } }",
		`Binary(Eq b1 Block(
        	Binary(Eq b2 Block(
        	Return()))))`)

	// bug from continue after finding "a" for first block
	// NOTE: sibling sharing no longer detected (minor breaking change)
	test(`a = F(); F({ x = F(a) }); F({ F(x) })`,
		`Binary(Eq a Call(F))
		Call(F Block(
			Binary(Eq x Call(F a))))
		Call(F Block-func(
			Call(F x)))`)

	test("b1 = {|x| x }; b2 = {|x| x }", // siblings
		`Binary(Eq b1 Block-func(x
			x))
		Binary(Eq b2 Block-func(x
			x))`)

	// nested block sharing with grandparent param (not immediate parent)
	test("b1 = {|x| b2 = { b3 = { x } }}",
		`Binary(Eq b1 Block(x
			Binary(Eq b2 Block(
				Binary(Eq b3 Block(
					x))))))`)
}

func TestSlotAssignments(t *testing.T) {
	parse := func(src string) *ast.Function {
		t.Helper()
		p := compile.NewParser("function () {\n" + src + "\n}")
		f := p.Function()
		ast.Blocks(f)
		return f
	}

	t.Run("no sharing", func(t *testing.T) {
		// function () { x = 1; b = { y = 2 } }
		f := parse("x = 1; b = { y = 2 }")
		// Outer: x=0, b=1 (no shared vars)
		x := f.Vars["x"]
		b := f.Vars["b"]
		assert.T(t).This(x >= 0 && x < SharedSlotStart).Is(true)
		assert.T(t).This(b >= 0 && b < SharedSlotStart).Is(true)
		assert.T(t).This(x != b).Is(true)
		// Block: y=0 (no shared vars, CompileAsFunction=true)
		block := f.Body[1].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(block.Vars["y"]).Is(int16(0))
		assert.T(t).This(block.CompileAsFunction).Is(true)
	})

	t.Run("simple sharing", func(t *testing.T) {
		// function () { x = 1; b = { x } }
		f := parse("x = 1; b = { x }")
		// Outer: x=SharedSlotStart (shared)
		assert.T(t).This(f.Vars["x"]).Is(int16(SharedSlotStart))
		assert.T(t).This(f.Vars["b"]).Is(int16(0))
		// Block: x=SharedSlotStart (shared), CompileAsFunction=false
		block := f.Body[1].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(block.Vars["x"]).Is(int16(SharedSlotStart))
		assert.T(t).This(block.CompileAsFunction).Is(false)
	})

	t.Run("parameter shadowing", func(t *testing.T) {
		// function (x) { b = { |x| x } }
		p := compile.NewParser("function (x) { b = { |x| x } }")
		f := p.Function()
		ast.Blocks(f)
		// Outer: x=0 (param, not shared)
		assert.T(t).This(f.Vars["x"]).Is(int16(0))
		assert.T(t).This(f.Vars["b"]).Is(int16(1))
		// Block: x=0 (param, shadows outer), CompileAsFunction=true
		block := f.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(block.Vars["x"]).Is(int16(0))
		assert.T(t).This(block.CompileAsFunction).Is(true)
	})

	t.Run("nested blocks sharing with outer", func(t *testing.T) {
		// function () { x = 1; b1 = { b2 = { x } } }
		f := parse("x = 1; b1 = { b2 = { x } }")
		// Outer: x=SharedSlotStart (shared)
		assert.T(t).This(f.Vars["x"]).Is(int16(SharedSlotStart))
		// b1: CompileAsFunction=false (contains closure)
		b1 := f.Body[1].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(b1.CompileAsFunction).Is(false)
		// b2: x=SharedSlotStart (shared), CompileAsFunction=false
		b2 := b1.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(b2.Vars["x"]).Is(int16(SharedSlotStart))
		assert.T(t).This(b2.CompileAsFunction).Is(false)
	})

	t.Run("mixed shared and local", func(t *testing.T) {
		// function () { x = 1; y = 2; b = { x; z = 3 } }
		f := parse("x = 1; y = 2; b = { x; z = 3 }")
		// Outer: x=SharedSlotStart (shared), y=0, b=1 (locals in order encountered after shared)
		assert.T(t).This(f.Vars["x"]).Is(int16(SharedSlotStart))
		y := f.Vars["y"]
		b := f.Vars["b"]
		assert.T(t).This(y >= 0 && y < SharedSlotStart).Is(true)
		assert.T(t).This(b >= 0 && b < SharedSlotStart).Is(true)
		assert.T(t).This(y != b).Is(true)
		// Block: x=SharedSlotStart (shared), z=0 (local)
		block := f.Body[2].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(block.Vars["x"]).Is(int16(SharedSlotStart))
		assert.T(t).This(block.Vars["z"]).Is(int16(0))
		assert.T(t).This(block.CompileAsFunction).Is(false)
	})

	t.Run("multiple shared variables", func(t *testing.T) {
		// function () { x = 1; y = 2; b = { x; y } }
		f := parse("x = 1; y = 2; b = { x; y }")
		// Both x and y should be shared (>= SharedSlotStart)
		assert.T(t).This(f.Vars["x"] >= SharedSlotStart).Is(true)
		assert.T(t).This(f.Vars["y"] >= SharedSlotStart).Is(true)
		assert.T(t).This(f.Vars["x"] != f.Vars["y"]).Is(true) // different slots
		// Block: same shared slots as outer
		block := f.Body[2].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(block.Vars["x"]).Is(f.Vars["x"])
		assert.T(t).This(block.Vars["y"]).Is(f.Vars["y"])
		assert.T(t).This(block.CompileAsFunction).Is(false)
	})

	t.Run("block param sharing with outer var", func(t *testing.T) {
		// function () { x = 1; b = { |x| x } }
		f := parse("x = 1; b = { |x| x }")
		// Outer: x is local (not shared, shadowed by block param)
		assert.T(t).This(f.Vars["x"] >= 0 && f.Vars["x"] < SharedSlotStart).Is(true)
		// Block: x=0 (param, shadows outer), CompileAsFunction=true
		block := f.Body[1].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(block.Vars["x"]).Is(int16(0))
		assert.T(t).This(block.CompileAsFunction).Is(true)
	})

	t.Run("inner block references outer param", func(t *testing.T) {
		// function (p) { b1 = { b2 = { p } } }
		p := compile.NewParser("function (p) { b1 = { b2 = { p } } }")
		f := p.Function()
		ast.Blocks(f)
		// Outer: p=SharedSlotStart (shared)
		assert.T(t).This(f.Vars["p"]).Is(int16(SharedSlotStart))
		// b1: CompileAsFunction=false
		b1 := f.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(b1.CompileAsFunction).Is(false)
		// b2: p=SharedSlotStart (shared)
		b2 := b1.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		assert.T(t).This(b2.Vars["p"]).Is(int16(SharedSlotStart))
		assert.T(t).This(b2.CompileAsFunction).Is(false)
	})

	t.Run("inner block references ancestor block param", func(t *testing.T) {
		// function (x) { b1 = {|x| b2 = { b3 = { x } } } }
		p := compile.NewParser(`function (x) {
			b1 = {|x|
				b2 = {
					b3 = { x }
					b3()
				}
				b2()
			}
			b1(1)
		}`)
		f := p.Function()
		ast.Blocks(f)

		// outer x should remain local (shadowed by b1 param x)
		assert.T(t).This(f.Vars["x"]).Is(int16(0))

		b1 := f.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		b2 := b1.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		b3 := b2.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)

		// b1 parameter x should be shared with b3
		assert.T(t).This(b1.Vars["x"] >= SharedSlotStart).Is(true)
		assert.T(t).This(b3.Vars["x"]).Is(b1.Vars["x"])
		assert.T(t).This(b3.CompileAsFunction).Is(false)
	})

	t.Run("this, super, and dynamic vars are not shared", func(t *testing.T) {
		f := parse("b = { .x }")
		b := f.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		thisSlot := b.Vars["this"]
		assert.T(t).This(thisSlot >= 0 && thisSlot < SharedSlotStart).Is(true)
		_, hasOuterThis := f.Vars["this"]
		assert.T(t).This(hasOuterThis).Is(false)

		f = parse("b = { super.F() }")
		b = f.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		superSlot := b.Vars["super"]
		assert.T(t).This(superSlot >= 0 && superSlot < SharedSlotStart).Is(true)
		_, hasOuterSuper := f.Vars["super"]
		assert.T(t).This(hasOuterSuper).Is(false)

		f = parse("_name = 1; b = { _name }")
		outerDynamic := f.Vars["_name"]
		assert.T(t).This(outerDynamic >= 0 && outerDynamic < SharedSlotStart).Is(true)
		b = f.Body[1].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		innerDynamic := b.Vars["_name"]
		assert.T(t).This(innerDynamic >= 0 && innerDynamic < SharedSlotStart).Is(true)
		assert.T(t).This(b.CompileAsFunction).Is(true)
	})

	t.Run("bug", func(t *testing.T) {
		src := `function (x) {
			b1 = {|x|
				b2 = { x }
				b2()
			}
			b1(1)
		}`
		p := compile.NewParser(src)
		f := p.Function()
		ast.Blocks(f)
		assert.This(f.Vars).Is(map[string]int16{"x": 0, "b1": 1})
		// outer `x` should not be shared
	})

	t.Run("bug2", func(t *testing.T) {
		src := `function (x) {
			b1 = {|x|
				b2 = { b3 = { b4 = { x } } }
			}
		}`
		p := compile.NewParser(src)
		f := p.Function()
		ast.Blocks(f)
		assert.This(f.Vars).Is(map[string]int16{"x": 0, "b1": 1})
		// outer `x` should not be shared
	})

	t.Run("nested capture keeps nearest binding", func(t *testing.T) {
		src := `function (x) {
			b1 = {|x|
				b2 = {
					b3 = { x }
					x
					b3()
				}
				b2()
			}
			b1(1)
		}`
		p := compile.NewParser(src)
		f := p.Function()
		ast.Blocks(f)

		// outer x is shadowed by b1 param x, so outer should remain local
		assert.T(t).This(f.Vars["x"]).Is(int16(0))

		b1 := f.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		b2 := b1.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)
		b3 := b2.Body[0].(*ast.ExprStmt).E.(*ast.Binary).Rhs.(*ast.Block)

		// b1 param x is the nearest binding and should be shared with b2 and b3
		assert.T(t).This(b1.Vars["x"] >= SharedSlotStart).Is(true)
		assert.T(t).This(b2.Vars["x"]).Is(b1.Vars["x"])
		assert.T(t).This(b3.Vars["x"]).Is(b1.Vars["x"])
	})
}
