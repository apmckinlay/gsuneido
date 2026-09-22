// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core_test

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/builtin"
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBuiltinString(t *testing.T) {
	f := Global.GetName(nil, "Type")
	assert.T(t).This(f.String()).Is("Type /* builtin function */")
	f = Global.GetName(nil, "Object")
	assert.T(t).This(f.String()).Is("Object /* builtin function */")
}

func TestThrow(t *testing.T) {
	assert.TestOnlyIndividually(t)

	f := compile.Constant(`function () {
        for ..100000
               try 
                    throw 'test'
                catch (e)
                    {}
	}`)
	var th Thread
	th.Call(f)
	heap := builtin.HeapSys()
	fmt.Println(trace.Number(heap))
	assert.That(heap < 32_000_000)
}

func TestMulti(t *testing.T) {
	f := compile.Constant(`function () {
	    f = function() { return 12,34 }
		a,b = f()
		return a is 12 and b is 34
	}`)
	var th Thread
	assert.This(th.Call(f)).Is(True)

	f = compile.Constant(`function () {
		fn = function (unused) 
			{
			return 12, 34
			}
		for foo in Object(1)
			{
			a, b = fn(foo)
			}
		return a is 12 and b is 34
	}`)
	assert.This(th.Call(f)).Is(True)

	f = compile.Constant(`function () {
		f = function() { return 1,2 }
		g = function() { return /* nothing */ }
		f()
		x, y = g()
	}`)
	assert.This(func() { th.Call(f) }).Panics("multiple return/assign mismatch")

	f = compile.Constant(`function () {
		f = function() { return 1,2 }
		g = function() { return 0 }
		f()
		x, y = g()
	}`)
	assert.This(func() { th.Call(f) }).Panics("multiple return/assign mismatch")
}

func TestTrinaryDiscardMulti(t *testing.T) {
	assert := assert.T(t)
	var th Thread

	// make sure that a trinary discards ReturnMulti properly

	// with op.Return
	f := compile.Constant(`function () {
		multi = function() { return 1, 2 }
		inner = function(multi, cond) {
			cond ? multi() : 0
			return 5
		}
		a, b = inner(multi, true)
	}`)
	assert.This(func() { th.Call(f) }).Panics("multiple return/assign mismatch")

	// with op.BlockReturn
	f = compile.Constant(`function () {
		multi = function() { return 1, 2 }
		inner = function(multi, cond) {
			cond ? multi() : 0
			#(111).Eval({ return 5 })
		}
		a, b = inner(multi, true)
	}`)
	assert.This(func() { th.Call(f) }).Panics("multiple return/assign mismatch")

	// with op.Gather
	f = compile.Constant(`function () {
		multi = function() { return 1, 2 }
		inner = function(multi, cond) {
			cond ? multi() : 0
			#(111).Eval({ return })
		}
		@ob = inner(multi, true)
		return ob
	}`)
	assert.This(th.Call(f)).Is(&SuObject{})
	assert.That(len(th.ReturnMulti) == 0)
}

func TestInRange(t *testing.T) {
	options.StrictCompare = true
	defer func() {
		options.StrictCompare = false
	}()
	f := compile.Constant(`function (x) { 0 < x and x <= 9 }`)
	var th Thread
	assert.That(th.Call(f, IntVal(0)) == False)
	assert.That(th.Call(f, IntVal(1)) == True)
	assert.That(th.Call(f, IntVal(9)) == True)
	assert.That(th.Call(f, IntVal(10)) == False)
	assert.That(th.Call(f, True) == False)
	assert.That(th.Call(f, EmptyStr) == False)
}

func BenchmarkForInSeq(b *testing.B) {
	f := compile.Constant(`function () {
		for i in Seq(1000)
		    {}
	}`)
	var th Thread
	for b.Loop() {
		th.Call(f)
	}
}

func BenchmarkForInCounted(b *testing.B) {
	f := compile.Constant(`function () {
		for i in ..1000
		    {}
	}`)
	var th Thread
	for b.Loop() {
		th.Call(f)
	}
}

func BenchmarkForClassic(b *testing.B) {
	f := compile.Constant(`function () {
		for (i = 0; i < 1000; i++)
		    {}
	}`)
	var th Thread
	for b.Loop() {
		th.Call(f)
	}
}

func TestNaming(t *testing.T) {
	var th Thread
	builtin.DefDef()
	test := func(src, expected string) {
		t.Helper()
		f := compile.Constant("function () {\n" + src + "\n}").(*SuFunc)
		result := th.Call(f)
		assert.T(t).This(result).Is(SuStr(expected))
	}
	test(`foo = function(){}; Name(foo)`, "foo")
	test(`foo = class{}; Name(foo)`, "foo")
	test(`foo = bar = class{}; Name(bar)`, "bar")
	test(`Def('Tmp', 'function(){}'); Name(Tmp)`, "Tmp")
	test(`Def('Tmp', 'function(){ return function(){} }'); Name(Tmp())`, "Tmp")
	test(`Def('Tmp', 'function(){ return {} }'); Name(Tmp())`, "Tmp")
	test(`Def('Tmp', 'function(){ fn = function(){} }'); Name(Tmp())`, "Tmp fn")
	test(`Def('Tmp', 'function(){ b = {} }'); Name(Tmp())`, "Tmp b")
	test(`Def('Tmp', 'class { F(){} }'); Name(Tmp.F)`, "Tmp.F")
	test(`Def('Tmp', 'class { Inner: class { F(){} } }');
		Name(Tmp.Inner.F)`, "Tmp.Inner.F")
	test(`Def('Tmp', 'function(){ myclass = class { F(){} } }');
		Name(Tmp().F)`, "Tmp myclass.F")
	test(`Def('Tmp', 'function() { Object(class{}) }'); Name(Tmp()[0])`,
		"Tmp")
	test(`Def('Tmp', 'class { A() { class { B(){} } } }'); Name(Tmp.A().B)`,
		"Tmp.A.B")
}

func BenchmarkCat2(b *testing.B) {
	f := compile.Constant(
		`function ()
			{
			s = ''
			for (i = 0; i < 1000; ++i)
				s $= "abc"
			}`).(*SuFunc)
	var th Thread
	for b.Loop() {
		th.Call(f)
	}
}

func BenchmarkJoin2(b *testing.B) {
	f := compile.Constant(
		`function ()
			{
			ob = Object()
			for (i = 0; i < 1000; ++i)
				ob.Add("abc")
			ob.Join()
			}`).(*SuFunc)
	var th Thread
	for b.Loop() {
		th.Call(f)
	}
}

func BenchmarkBase(b *testing.B) {
	f := compile.Constant(
		`function ()
			{
			for (i = 0; i < 1000; ++i)
				;
			}`).(*SuFunc)
	var th Thread
	for b.Loop() {
		th.Call(f)
	}
}

// compare to BenchmarkJit in interp_test.go
func BenchmarkInterp2(b *testing.B) {
	src := `function (x,y) { x + y }`
	if !Global.Exists("ADD") {
		Global.Add("ADD", compile.Constant(src).(*SuFunc))
	}
	src = `function () {
		sum = 0
		for (i = 0; i < 100; ++i)
			sum = ADD(sum, i)
		return sum
	}`
	fn := compile.Constant(src).(*SuFunc)
	var th Thread
	for b.Loop() {
		result := th.Call(fn)
		if !result.Equal(SuInt(4950)) {
			panic("wrong result " + result.String())
		}
	}
}

func BenchmarkCall(b *testing.B) {
	f := Global.GetName(nil, "Type")
	as := &ArgSpec1
	th := &Thread{}
	th.Push(SuInt(123))
	for b.Loop() {
		f.Call(th, nil, as)
	}
}

func TestCoverage(t *testing.T) {
	options.Coverage.Store(true)
	fn := compile.Constant(`function()
		{
		x = 0
		for (i = 0; i < 10; ++i)
			x += i
		return x
		}`).(*SuFunc)
	fn.StartCoverage(true)
	var th Thread
	th.Call(fn)
	cover := fn.StopCoverage()
	assert.T(t).This(cover).
		Is(compile.Constant("#(17: 1, 25: 1, 53: 10, 62: 1)").(*SuObject))
}

func TestSuClassDefaultGet(t *testing.T) {
	f := compile.Constant(`function() {
		c = class {
			Default() { return 123 }
		}
		c.X
	}`)
	th := &Thread{}
	assert.This(th.Call(f).String()).Is("Default(X /* method */")
}

func TestAtAssign(t *testing.T) {
	var th Thread

	// multiple return values
	f := compile.Constant(`function () {
		f = function() { return 12, 34 }
		@ob = f()
		return ob
	}`)
	result := th.Call(f)
	assert.T(t).This(result).Is(SuObjectOf(SuInt(12), SuInt(34)))

	// single return value
	f = compile.Constant(`function () {
		f = function() { return 42 }
		@ob = f()
		return ob
	}`)
	result = th.Call(f)
	assert.T(t).This(result).Is(SuObjectOf(SuInt(42)))

	// no return value
	f = compile.Constant(`function () {
		f = function() { }
		@ob = f()
		return ob
	}`)
	result = th.Call(f)
	assert.T(t).This(result).Is(&SuObject{})

	// method call
	f = compile.Constant(`function () {
		obj = class { F() { return 1, 2, 3 } }
		instance = obj()
		@ob = instance.F()
		return ob
	}`)
	result = th.Call(f)
	assert.T(t).This(result).Is(SuObjectOf(SuInt(1), SuInt(2), SuInt(3)))
}

func TestReturnSpread(t *testing.T) {
	var th Thread

	// empty object - bare return
	f := compile.Constant(`function () {
		return @Object()
	}`)
	assert.This(th.Call(f)).Is(nil)
	assert.That(len(th.ReturnMulti) == 0)

	// single value
	f = compile.Constant(`function () {
		return @Object(42)
	}`)
	assert.This(th.Call(f)).Is(SuInt(42))
	assert.That(len(th.ReturnMulti) == 0)

	// multiple values - direct call returns nil (like return 1,2,3)
	f = compile.Constant(`function () {
		return @Object(0, 1, "")
	}`)
	assert.This(th.Call(f)).Is(nil)
	assert.This(th.ReturnMulti).Is([]Value{EmptyStr, One, Zero}) // reverse

	// multiple values with assignment
	f = compile.Constant(`function () {
		fn = function() { return @Object(10, 20, 30) }
		a, b, c = fn()
		return a is 10 and b is 20 and c is 30
	}`)
	assert.This(th.Call(f)).Is(True)

	// error: named members
	f = compile.Constant(`function () {
		return @Object(a: 1)
	}`)
	assert.This(func() { th.Call(f) }).Panics("return @ cannot include named members")

	// error: not an object
	f = compile.Constant(`function () {
		return @123
	}`)
	assert.This(func() { th.Call(f) }).Panics("return @ requires an object")
}

func TestBlockReturnMulti(t *testing.T) {
	var th Thread

	// multiple return values from a directly called block
	f := compile.Constant(`function () {
		inner = function() { blk = { return 12, 34 }; blk(); 123 }
		a, b = inner()
		return a is 12 and b is 34
	}`)
	assert.This(th.Call(f)).Is(True)

	// multiple return values from a block called by a builtin
	f = compile.Constant(`function () {
		inner = function() { #(1).Eval({ return 12, 34 }) }
		a, b = inner()
		return a is 12 and b is 34
	}`)
	assert.This(th.Call(f)).Is(True)

	// multiple return values from a block called from another function
	// with try/catch (like stdlib Each) which must not catch block returns
	f = compile.Constant(`function () {
		inner = function() {
			each = function(ob, blk) {
				for x in ob
					try
						blk(x)
					catch (e, "block:")
						if e is "block:break"
							break
				}
			each(Object(1), { |x| return 12, 34 })
			}
		a, b = inner()
		return a is 12 and b is 34
	}`)
	assert.This(th.Call(f)).Is(True)

	// direct call returns nil (like return 1,2,3)
	f = compile.Constant(`function () {
		inner = function() { blk = { return 0, 1, "" }; blk(); 123 }
		return inner()
	}`)
	assert.This(th.Call(f)).Is(nil)
	assert.This(th.ReturnMulti).Is([]Value{EmptyStr, One, Zero}) // reverse

	// gathered with @ob =
	f = compile.Constant(`function () {
		inner = function() { blk = { return 12, 34 }; blk(); 123 }
		@ob = inner()
		return ob
	}`)
	assert.This(th.Call(f)).Is(SuObjectOf(SuInt(12), SuInt(34)))

	// multiple return values from a nested block
	f = compile.Constant(`function () {
		inner = function() {
			blk = { b2 = { return 12, 34 }; b2(); 123 }
			blk()
			}
		a, b = inner()
		return a is 12 and b is 34
	}`)
	assert.This(th.Call(f)).Is(True)
}

func TestBlockReturnSpread(t *testing.T) {
	var th Thread

	// multiple values
	f := compile.Constant(`function () {
		inner = function() { blk = { return @Object(10, 20) }; blk(); 123 }
		a, b = inner()
		return a is 10 and b is 20
	}`)
	assert.This(th.Call(f)).Is(True)

	// single value
	f = compile.Constant(`function () {
		inner = function() { blk = { return @Object(42) }; blk(); 123 }
		return inner()
	}`)
	assert.This(th.Call(f)).Is(SuInt(42))
	assert.That(len(th.ReturnMulti) == 0)

	// empty object - bare return
	f = compile.Constant(`function () {
		inner = function() { blk = { return @Object() }; blk(); 123 }
		return inner()
	}`)
	assert.This(th.Call(f)).Is(nil)

	// error: named members
	f = compile.Constant(`function () {
		inner = function() { blk = { return @Object(a: 1) }; blk(); 123 }
		inner()
	}`)
	assert.This(func() { th.Call(f) }).Panics("return @ cannot include named members")

	// error: not an object
	f = compile.Constant(`function () {
		inner = function() { blk = { return @123 }; blk(); 123 }
		inner()
	}`)
	assert.This(func() { th.Call(f) }).Panics("return @ requires an object")
}
