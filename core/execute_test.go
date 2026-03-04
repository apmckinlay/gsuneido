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
