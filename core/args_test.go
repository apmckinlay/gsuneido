// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

var atParamSpec = &ParamSpec{Nparams: 1, Flags: []Flag{AtParam}}

func TestArgs(t *testing.T) {
	assert := assert.T(t).This
	th := &Thread{}
	setStack := func(nums ...int) {
		th.Reset()
		for _, n := range nums {
			th.Push(SuInt(n))
		}
	}
	ckStack := func(vals ...int) {
		t.Helper()
		assert(fmt.Sprint(th.stack[:th.sp])).Is(fmt.Sprint(vals))
	}

	// 0 args => 0 params
	f := &ParamSpec{}
	as := &ArgSpec0
	th.Args(f, as)

	// @arg => @param
	f = atParamSpec
	as = &ArgSpecEach0
	th.Reset()
	th.Push(makeOb())
	th.Args(f, as)
	assert(th.sp).Is(1)
	assert(th.stack[0]).Is(makeOb())

	// @+1arg => @param
	f = atParamSpec
	as = &ArgSpecEach1
	th.Reset()
	th.Push(makeOb())
	th.Args(f, as)
	assert(th.sp).Is(1)
	assert(th.stack[0]).Is(makeOb().Slice(1))

	// 2 args => 2 params
	f = &ParamSpec{Nparams: 2, Flags: []Flag{0, 0}}
	as = &ArgSpec2
	setStack(11, 22)
	th.Args(f, as)
	ckStack(11, 22)

	// 1 args => 2 params
	f = &ParamSpec{Nparams: 2, Flags: []Flag{0, 0}}
	as = &ArgSpec1
	setStack(11)
	assert(func() { th.Args(f, as) }).Panics("missing argument")

	// 2 args => 1 param
	f = &ParamSpec{Nparams: 1, Flags: []Flag{0}}
	as = &ArgSpec2
	setStack(11, 22)
	assert(func() { th.Args(f, as) }).Panics("too many arguments")

	// 1 arg => 2 params with 1 default
	f = &ParamSpec{Nparams: 2, Flags: []Flag{0, 0},
		Ndefaults: 1, Values: []Value{SuInt(22)}}
	as = &ArgSpec1
	setStack(11)
	th.Args(f, as)
	ckStack(11, 22)

	// 2 arg => 3 params with 2 defaults
	f = &ParamSpec{Nparams: 3, Flags: []Flag{0, 0, 0},
		Ndefaults: 2, Values: []Value{False, One}}
	as = &ArgSpec2
	setStack(2, 5)
	th.Args(f, as)
	ckStack(2, 5, 1)

	// all named
	f = &ParamSpec{Nparams: 3, Flags: []Flag{0, 0, 0},
		Names: []string{"a", "b", "c"}}
	as = &ArgSpec{Nargs: 4,
		Names: vals("c", "b", "a", "d"), Spec: []byte{1, 0, 2, 3}} // b, c, a, d
	setStack(22, 33, 11, 44)
	th.Args(f, as)
	ckStack(11, 22, 33)

	// mixed
	f = &ParamSpec{Nparams: 4, Flags: []Flag{0, 0, 0},
		Names: []string{"a", "b", "c", "d"}}
	as = &ArgSpec{Nargs: 4,
		Names: vals("c", "b", "a", "d"), Spec: []byte{3, 0}} // d, c
	setStack(22, 33, 11, 44) // fn(22, 33, d: 11, c: 44)
	th.Args(f, as)
	ckStack(22, 33, 44, 11)

	// args => @param
	f = atParamSpec
	as = &ArgSpec{Nargs: 4,
		Names: vals("c", "b", "a", "d"), Spec: []byte{1, 2}} // b, a
	setStack(11, 22, 44, 33)
	th.Args(f, as)
	assert(th.sp).Is(1)
	assert(th.stack[0]).Is(makeOb())

	// @mixed => params
	f = &ParamSpec{Nparams: 4, Flags: []Flag{0, 0, 0, 0},
		Names: []string{"d", "c", "b", "a"}}
	as = &ArgSpecEach0
	th.Reset()
	th.Push(makeOb())
	th.Args(f, as)
	ckStack(11, 22, 44, 33)

	// @list
	th.Reset()
	th.Push(SuObjectOf(SuInt(1), SuInt(2), SuInt(3), SuInt(4)))
	th.Args(f, as)
	ckStack(1, 2, 3, 4)

	// @+1 list
	th.Reset()
	th.Push(SuObjectOf(SuInt(1), SuInt(2), SuInt(3), SuInt(4), SuInt(5)))
	th.Args(f, &ArgSpecEach1)
	ckStack(2, 3, 4, 5)

	// @args => one param
	f = &ParamSpec{Nparams: 1, Flags: []Flag{0},
		Names: []string{"a"}}
	as = &ArgSpecEach0
	th.Reset()
	th.Push(SuObjectOf(SuInt(123)))
	th.Args(f, as)
	ckStack(123)

	// dynamic
	setStack(111, 123)
	th.frames[0] = Frame{locals: locals{v: th.stack[0:]},
		fn: &SuFunc{ParamSpec: ParamSpec{Names: []string{"x", "_dyn"}}}}
	th.fp++
	f = &ParamSpec{Nparams: 1, Flags: []Flag{DynParam}, Names: []string{"dyn"}}
	as = &ArgSpec0
	th.Args(f, as)
	ckStack(111, 123, 123)
}

func makeOb() *SuObject {
	var ob SuObject
	ob.Add(SuInt(11))
	ob.Add(SuInt(22))
	ob.Set(SuStr("a"), SuInt(33))
	ob.Set(SuStr("b"), SuInt(44))
	return &ob
}

func vals(names ...string) []Value {
	vals := make([]Value, len(names))
	for i, s := range names {
		vals[i] = SuStr(s)
	}
	return vals
}
