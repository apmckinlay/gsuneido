// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core_test

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"

	"github.com/apmckinlay/gsuneido/compile"
)

// makeWarningsThrow makes any Warning panic so it can be asserted with Panics.
// It returns a func that restores the previous setting.
// Normal usage: defer makeWarningsThrow()()
func makeWarningsThrow() func() {
	prev := options.WarningsThrow.Load()
	options.WarningsThrow.Store(options.AllWarningsThrow)
	return func() { options.WarningsThrow.Store(prev) }
}

func TestClosure_LoopSharedModifiedWarning(t *testing.T) {
	// a block created in a loop, its instances sharing a variable that is
	// reassigned, and dispatched to threads
	src := `function () {
		fs = []
		n = 1
		for (i = 0; i < 2; ++i)
			fs.Add({ n })
		mod = {|x| n = x }
		return [fs, mod]
	}`
	fn := compile.Constant(src)
	var th Thread
	v := th.Call(fn)
	v.SetConcurrent()
	ob := v.(interface{ ListGet(int) Value })
	fs := ob.ListGet(0).(interface{ ListGet(int) Value })

	defer makeWarningsThrow()()

	// reassign the shared variable while concurrent
	th.Call(ob.ListGet(1), SuInt(2))

	// dispatching the same block a second time warns
	NoteThreadClosure(fs.ListGet(0))
	assert.T(t).This(func() { NoteThreadClosure(fs.ListGet(1)) }).
		Panics("thread closure")

	// but the warning is only given once
	NoteThreadClosure(fs.ListGet(1))
}

func TestClosure_SeparateSharedNoWarning(t *testing.T) {
	// separate blocks sharing a variable are not the loop case
	src := `function () {
		n = 0
		a = { n += 3 }
		b = { n += 2 }
		return [a, b]
	}`
	fn := compile.Constant(src)
	var th Thread
	v := th.Call(fn)
	v.SetConcurrent()
	ob := v.(interface{ ListGet(int) Value })

	defer makeWarningsThrow()()

	th.Call(ob.ListGet(0)) // modifies the shared variable while concurrent

	// each block is dispatched, but they are different blocks, must not warn
	NoteThreadClosure(ob.ListGet(0))
	NoteThreadClosure(ob.ListGet(1))
}

func TestClosure_LoopSharedUnmodifiedNoWarning(t *testing.T) {
	// QueryFuzz-like: a block created in a loop dispatched to threads
	// but without reassigning the shared variables
	src := `function () {
		fs = []
		for (i = 0; i < 2; ++i)
			fs.Add({ i })
		fs
	}`
	fn := compile.Constant(src)
	var th Thread
	v := th.Call(fn)
	v.SetConcurrent()
	fs := v.(interface{ ListGet(int) Value })

	defer makeWarningsThrow()()

	// the same block dispatched twice, but nothing is reassigned, must not warn
	NoteThreadClosure(fs.ListGet(0))
	NoteThreadClosure(fs.ListGet(1))
}

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
	// fmt.Println(DisasmOps(f.(*SuFunc)))
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
