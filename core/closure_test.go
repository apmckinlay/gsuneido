// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core_test

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"

	"github.com/apmckinlay/gsuneido/compile"
)

func TestClosure_rule(t *testing.T) {
	src := `function () {
        n = 0
        r = Record()
        r.AttachRule('foo', { n++ })
        r.foo
        n
    }`

	f := compile.Constant(src)
	var th Thread
	result := th.Call(f)
	assert.This(result).Is(SuInt(1))
}

func TestClosure_return(t *testing.T) {
	src := `function () {
		f = function(){}
		b = { return f() }
		b()
		123
    }`

	f := compile.Constant(src)
	var th Thread
	result := th.Call(f)
	assert.This(result).Is(nil)
}

func TestClosure_nested(t *testing.T) {
	src := `function () {
		f = function (x) { return {|a| x * a } }
		b = f(2);
		b(3)
    }`

	f := compile.Constant(src)
	// fmt.Println(DisasmOps(f.(*SuFunc)))
	var th Thread
	result := th.Call(f)
	assert.This(result).Is(SuInt(6))
}

func TestClosure_observer1(t *testing.T) {
	src := `function () {
		r = Record()
		r.Observer({|member| o = member })
		r.foo = 123
		o
    }`

	f := compile.Constant(src)
	// fmt.Println(DisasmOps(f.(*SuFunc)))
	var th Thread
	result := th.Call(f)
	assert.This(result).Is(SuStr("foo"))
}

func TestClosure_observer2(t *testing.T) {
	src := `function () {
		r = Record()
		r.Observer({|member| o = member; .bar = 456 })
		r.foo = 123
		o
    }`

	f := compile.Constant(src)
	fmt.Println(DisasmOps(f.(*SuFunc)))
	var th Thread
	result := th.Call(f)
	assert.This(result).Is(SuStr("bar"))
}

func TestClosure_ExceptionLocals(t *testing.T) {
	src := `function () {
		local = "outer"
		f = function(x) {
			inner = function() {
				throw "error"
			}
			inner()
			return x
		}
		try {
			f(local)
		} catch(e, "error") {
			return local
		}
    }`

	f := compile.Constant(src)
	var th Thread
	result := th.Call(f)
	assert.This(result).Is(SuStr("outer"))
}

func TestClosure_ExceptionShadowing(t *testing.T) {
	src := `function () {
		local = "outer"
		f = function(local) {
			inner = function() {
				throw "error"
			}
			inner()
			return local
		}
		try {
			f("shadowed")
		} catch(e) {
			return local
		}
    }`

	f := compile.Constant(src)
	var th Thread
	result := th.Call(f)
	assert.This(result).Is(SuStr("outer"))
}

func TestClosure_ExceptionNested(t *testing.T) {
	src := `function () {
		local1 = "first"
		local2 = "second"
		f = function() {
			g = function() {
				h = function() {
					throw "error"
				}
				h()
				return local2
			}
			g()
			return local1
		}
		try {
			f()
		} catch(e) {
			return [local1, local2]
		}
    }`

	f := compile.Constant(src)
	var th Thread
	result := th.Call(f)
	expected := &SuObject{}
	expected.Add(SuStr("first"))
	expected.Add(SuStr("second"))
	assert.This(result).Is(expected)
}

func TestClosure_CallstackLocals1(t *testing.T) {
	src := `function (x, y) {
		f = {|y| Display(x); throw "error" }
		try {
			f("shadowed")
		} catch(e, "error") {
			cs = e.Callstack()
			return cs[0].locals
		}
    }`

	f := compile.Constant(src)
	var th Thread
	result := th.Call(f, Zero, One)
	ob := &SuObject{}
	ob.Set(SuStr("x"), Zero)
	ob.Set(SuStr("y"), SuStr("shadowed"))
	assert.This(result).Is(ob)
}

func TestClosure_CallstackLocals2(t *testing.T) {
	src := `function (v1) {
		b = {
			v2 = v1 + 1
			throw "error"
		}
		try {
			b()
		} catch(e, "error") {
			cs = e.Callstack()
			return cs[0].locals
		}
    }`

	f := compile.Constant(src)
	var th Thread
	result := th.Call(f, One)
	assert.T(t).This(result.String()).Is("#(v2: 2, v1: 1)")
}

func TestClosure_CallstackLocals3(t *testing.T) {
	src := `function (x) {
		b1 = {|x|
			b2 = {
				b3 = {|x|
					b4 = {
						x
						throw "error"
					}
					b4()
				}
				b3(3)
			}
			b2()
		}
		try {
			b1(1)
		} catch(e, "error") {
			return e.Callstack()
		}
    }`

	f := compile.Constant(src)
	var th Thread
	cs := th.Call(f, Zero)
	assert.T(t).This(cs.Get(nil, SuInt(0)).Get(nil, SuStr("locals")).String()).
		Is("#(x: 3)")
	assert.T(t).This(cs.Get(nil, SuInt(1)).Get(nil, SuStr("locals")).String()).
		Is("#(x: 3, b4: /* closure */)")
	assert.T(t).This(cs.Get(nil, SuInt(2)).Get(nil, SuStr("locals")).String()).
		Is("#(b3: /* closure */)")
	assert.T(t).This(cs.Get(nil, SuInt(3)).Get(nil, SuStr("locals")).String()).
		Is("#(x: 1, b2: /* closure */)")
	assert.T(t).This(cs.Get(nil, SuInt(4)).Get(nil, SuStr("locals")).String()).
		Is("#(x: 0, b1: /* closure */)")
}

func TestDynamicBug(t *testing.T) {
	src := `function () {
		_p = 123
		c = class { New(._P) { } A() { .P } }
		new c()
		}`
	f := compile.Constant(src)
	var th Thread
	th.Call(f)
}

func TestCompileNamesOverwriteRepro(t *testing.T) {
	src := `function () {
		x = 1
		#(1).Each({ x })
		x++
		RetryTransaction({|t|
			#(1).Each({|it| t })
			part = 0
			part++
			})
		}`
	compile.Constant(src)
}
