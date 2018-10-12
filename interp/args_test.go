package interp

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestArgs(t *testing.T) {
	th := &Thread{}
	setStack := func(nums ...int) {
		th.sp = 0
		for _, n := range nums {
			th.Push(SuInt(n))
		}
	}
	ckStack := func(vals ...int) {
		t.Helper()
		Assert(t).That(fmt.Sprint(th.stack[:th.sp]), Equals(fmt.Sprint(vals)))
	}

	// 0 args => 0 params
	f := &Func{}
	a := ArgSpec{}
	th.args(f, a)

	// @arg => @param
	f = &Func{Nparams: 1, Flags: []Flag{AtParam}}
	a = ArgSpec{Unnamed: EACH}
	th.sp = 0
	th.Push(makeOb())
	th.args(f, a)
	Assert(t).That(th.sp, Equals(1))
	Assert(t).True(th.stack[0].Equal(makeOb()))

	// @+1arg => @param
	f = &Func{Nparams: 1, Flags: []Flag{AtParam}}
	a = ArgSpec{Unnamed: EACH1}
	th.sp = 0
	th.Push(makeOb())
	th.args(f, a)
	Assert(t).That(th.sp, Equals(1))
	Assert(t).True(th.stack[0].Equal(makeOb().Slice(1)))

	// 2 args => 2 params
	f = &Func{Nparams: 2, Flags: []Flag{0, 0}}
	a = ArgSpec{Unnamed: 2}
	setStack(11, 22)
	th.args(f, a)
	ckStack(11, 22)

	// 1 args => 2 params
	f = &Func{Nparams: 2, Flags: []Flag{0, 0}}
	a = ArgSpec{Unnamed: 1}
	setStack(11)
	Assert(t).That(func() { th.args(f, a) }, Panics("missing argument"))

	// 2 args => 1 param
	f = &Func{Nparams: 1, Flags: []Flag{0}}
	a = ArgSpec{Unnamed: 2}
	setStack(11, 22)
	Assert(t).That(func() { th.args(f, a) }, Panics("too many arguments"))

	// 1 arg => 2 params with 1 default
	f = &Func{Nparams: 2, Flags: []Flag{0, 0},
		Ndefaults: 1, Values: []Value{SuInt(22)}}
	a = ArgSpec{Unnamed: 1}
	setStack(11)
	th.args(f, a)
	ckStack(11, 22)

	// all named
	f = &Func{Nparams: 3, Flags: []Flag{0, 0, 0},
		Strings: []string{"a", "b", "c"}}
	a = ArgSpec{Unnamed: 0,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{1, 0, 2, 3}} // b, c, a, d
	setStack(22, 33, 11, 44)
	th.args(f, a)
	ckStack(11, 22, 33)

	// mixed
	f = &Func{Nparams: 4, Flags: []Flag{0, 0, 0},
		Strings: []string{"a", "b", "c", "d"}}
	a = ArgSpec{Unnamed: 2,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{3, 0}} // d, c
	setStack(22, 33, 11, 44)  // fn(22, 33, d: 11, c: 44)
	th.args(f, a)
	ckStack(22, 33, 44, 11)

	// args => @param
	f = &Func{Nparams: 1, Flags: []Flag{AtParam}}
	a = ArgSpec{Unnamed: 2,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{1, 2}} // b, a
	setStack(11, 22, 44, 33)
	th.args(f, a)
	Assert(t).That(th.sp, Equals(1))
	Assert(t).That(th.stack[0].String(), Equals(makeOb().String()))

	// @args => params
	f = &Func{Nparams: 4, Flags: []Flag{0, 0, 0, 0},
		Strings: []string{"d", "c", "b", "a"}}
	a = ArgSpec{Unnamed: EACH}
	th.sp = 0
	th.Push(makeOb())
	th.args(f, a)
	ckStack(11, 22, 44, 33)

	// dynamic
	th.frames = append(th.frames, Frame{
		fn:     &SuFunc{Func: Func{Strings: []string{"x", "_dyn"}}},
		locals: []Value{SuInt(111), SuInt(123)},
	})
	f = &Func{Nparams: 1, Flags: []Flag{DynParam}, Strings: []string{"dyn"}}
	a = ArgSpec{}
	th.sp = 0
	th.args(f, a)
	ckStack(123)
}

func makeOb() *SuObject {
	var ob SuObject
	ob.Add(SuInt(11))
	ob.Add(SuInt(22))
	ob.Put(SuStr("a"), SuInt(33))
	ob.Put(SuStr("b"), SuInt(44))
	return &ob
}
