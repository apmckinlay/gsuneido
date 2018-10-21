package runtime

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestArgs(t *testing.T) {
	th := &Thread{}
	setStack := func(nums ...int) {
		th.Reset()
		for _, n := range nums {
			th.Push(SuInt(n))
		}
	}
	ckStack := func(vals ...int) {
		t.Helper()
		Assert(t).That(fmt.Sprint(th.stack[:th.sp]), Equals(fmt.Sprint(vals)))
	}

	// 0 args => 0 params
	f := &ParamSpec{}
	a := ArgSpec{}
	th.Args(f, &a)

	// @arg => @param
	f = &ParamSpec{Nparams: 1, Flags: []Flag{AtParam}}
	a = ArgSpec{Unnamed: EACH}
	th.Reset()
	th.Push(makeOb())
	th.Args(f, &a)
	Assert(t).That(th.sp, Equals(1))
	Assert(t).True(th.stack[0].Equal(makeOb()))

	// @+1arg => @param
	f = &ParamSpec{Nparams: 1, Flags: []Flag{AtParam}}
	a = ArgSpec{Unnamed: EACH1}
	th.Reset()
	th.Push(makeOb())
	th.Args(f, &a)
	Assert(t).That(th.sp, Equals(1))
	Assert(t).True(th.stack[0].Equal(makeOb().Slice(1)))

	// 2 args => 2 params
	f = &ParamSpec{Nparams: 2, Flags: []Flag{0, 0}}
	a = ArgSpec{Unnamed: 2}
	setStack(11, 22)
	th.Args(f, &a)
	ckStack(11, 22)

	// 1 args => 2 params
	f = &ParamSpec{Nparams: 2, Flags: []Flag{0, 0}}
	a = ArgSpec{Unnamed: 1}
	setStack(11)
	Assert(t).That(func() { th.Args(f, &a) }, Panics("missing argument"))

	// 2 args => 1 param
	f = &ParamSpec{Nparams: 1, Flags: []Flag{0}}
	a = ArgSpec{Unnamed: 2}
	setStack(11, 22)
	Assert(t).That(func() { th.Args(f, &a) }, Panics("too many arguments"))

	// 1 arg => 2 params with 1 default
	f = &ParamSpec{Nparams: 2, Flags: []Flag{0, 0},
		Ndefaults: 1, Values: []Value{SuInt(22)}}
	a = ArgSpec{Unnamed: 1}
	setStack(11)
	th.Args(f, &a)
	ckStack(11, 22)

	// all named
	f = &ParamSpec{Nparams: 3, Flags: []Flag{0, 0, 0},
		Names: []string{"a", "b", "c"}}
	a = ArgSpec{Unnamed: 0,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{1, 0, 2, 3}} // b, c, a, d
	setStack(22, 33, 11, 44)
	th.Args(f, &a)
	ckStack(11, 22, 33)

	// mixed
	f = &ParamSpec{Nparams: 4, Flags: []Flag{0, 0, 0},
		Names: []string{"a", "b", "c", "d"}}
	a = ArgSpec{Unnamed: 2,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{3, 0}} // d, c
	setStack(22, 33, 11, 44) // fn(22, 33, d: 11, c: 44)
	th.Args(f, &a)
	ckStack(22, 33, 44, 11)

	// args => @param
	f = &ParamSpec{Nparams: 1, Flags: []Flag{AtParam}}
	a = ArgSpec{Unnamed: 2,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{1, 2}} // b, a
	setStack(11, 22, 44, 33)
	th.Args(f, &a)
	Assert(t).That(th.sp, Equals(1))
	Assert(t).That(th.stack[0].String(), Equals(makeOb().String()))

	// @args => params
	f = &ParamSpec{Nparams: 4, Flags: []Flag{0, 0, 0, 0},
		Names: []string{"d", "c", "b", "a"}}
	a = ArgSpec{Unnamed: EACH}
	th.Reset()
	th.Push(makeOb())
	th.Args(f, &a)
	ckStack(11, 22, 44, 33)

	// dynamic
	setStack(111, 123)
	th.frames[0] = Frame{
		fn: &SuFunc{ParamSpec: ParamSpec{Names: []string{"x", "_dyn"}}}}
	th.fp++
	f = &ParamSpec{Nparams: 1, Flags: []Flag{DynParam}, Names: []string{"dyn"}}
	a = ArgSpec{}
	th.Args(f, &a)
	ckStack(111, 123, 123)
}

func makeOb() *SuObject {
	var ob SuObject
	ob.Add(SuInt(11))
	ob.Add(SuInt(22))
	ob.Put(SuStr("a"), SuInt(33))
	ob.Put(SuStr("b"), SuInt(44))
	return &ob
}
