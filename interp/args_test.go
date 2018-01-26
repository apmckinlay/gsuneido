package interp

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestArgs(t *testing.T) {
	th := &Thread{}

	// 0 args => 0 params
	f := &SuFunc{}
	a := ArgSpec{}
	th.args(f, a)

	// @arg => @param
	f = &SuFunc{Nparams: 1, Nlocals: 1, Flags: []Flag{AT_F}}
	a = ArgSpec{Unnamed: EACH}
	th.stack = []Value{makeOb()}
	th.args(f, a)
	Assert(t).True(th.stack[0].Equals(makeOb()))

	// @+1arg => @param
	f = &SuFunc{Nparams: 1, Nlocals: 1, Flags: []Flag{AT_F}}
	a = ArgSpec{Unnamed: EACH1}
	th.stack = []Value{makeOb()}
	th.args(f, a)
	Assert(t).True(th.stack[0].Equals(makeOb().Slice(1)))

	// 2 args => 2 params
	f = &SuFunc{Nparams: 2, Nlocals: 2, Flags: []Flag{0, 0}}
	a = ArgSpec{Unnamed: 2}
	th.stack = []Value{SuInt(11), SuInt(22)}
	th.args(f, a)
	Assert(t).That(fmt.Sprint(th.stack), Equals("[11 22]"))

	// 1 args => 2 params
	f = &SuFunc{Nparams: 2, Nlocals: 2, Flags: []Flag{0, 0}}
	a = ArgSpec{Unnamed: 1}
	th.stack = []Value{SuInt(11)}
	Assert(t).That(func() { th.args(f, a) }, Panics("missing argument"))

	// 2 args => 1 param
	f = &SuFunc{Nparams: 1, Nlocals: 1, Flags: []Flag{0}}
	a = ArgSpec{Unnamed: 2}
	th.stack = []Value{SuInt(11), SuInt(22)}
	Assert(t).That(func() { th.args(f, a) }, Panics("too many arguments"))

	// 1 arg => 2 params with 1 default
	f = &SuFunc{Nparams: 2, Nlocals: 2, Flags: []Flag{0, 0},
		Ndefaults: 1, Values: []Value{SuInt(22)}}
	a = ArgSpec{Unnamed: 1}
	th.stack = []Value{SuInt(11)}
	th.args(f, a)
	Assert(t).That(fmt.Sprint(th.stack), Equals("[11 22]"))

	// all named
	f = &SuFunc{Nparams: 3, Nlocals: 3, Flags: []Flag{0, 0, 0},
		Strings: []string{"a", "b", "c"}}
	a = ArgSpec{Unnamed: 0,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{1, 0, 2, 3}} // b, c, a, d
	th.stack = []Value{SuInt(22), SuInt(33), SuInt(11), SuInt(44)}
	th.args(f, a)
	Assert(t).That(fmt.Sprint(th.stack), Equals("[11 22 33]"))

	// mixed
	f = &SuFunc{Nparams: 4, Nlocals: 4, Flags: []Flag{0, 0, 0},
		Strings: []string{"a", "b", "c", "d"}}
	a = ArgSpec{Unnamed: 2,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{3, 0}} // d, c
	th.stack = []Value{SuInt(11), SuInt(22), SuInt(44), SuInt(33)}
	th.args(f, a)
	Assert(t).That(fmt.Sprint(th.stack), Equals("[11 22 33 44]"))

	// args => @param
	f = &SuFunc{Nparams: 1, Nlocals: 1, Flags: []Flag{AT_F}}
	a = ArgSpec{Unnamed: 2,
		Names: []string{"c", "b", "a", "d"}, Spec: []byte{1, 2}} // b, a
	th.stack = []Value{SuInt(11), SuInt(22), SuStr("bb"), SuStr("aa")}
	th.args(f, a)
	Assert(t).That(len(th.stack), Equals(1))
	Assert(t).True(th.stack[0].Equals(makeOb()))

	// @args => params
	f = &SuFunc{Nparams: 4, Nlocals: 4, Flags: []Flag{0, 0, 0, 0},
		Strings: []string{"d", "c", "b", "a"}}
	a = ArgSpec{Unnamed: EACH}
	th.stack = []Value{makeOb()}
	th.args(f, a)
	Assert(t).That(fmt.Sprint(th.stack), Equals("[11 22 'bb' 'aa']"))

	// dynamic
	th.frames = append(th.frames, Frame{
		fn:     &SuFunc{Strings: []string{"x", "_dyn"}},
		locals: []Value{SuInt(111), SuInt(123)},
	})
	f = &SuFunc{Nparams: 1, Nlocals: 1, Flags: []Flag{DYN_F}, Strings: []string{"dyn"}}
	a = ArgSpec{}
	th.stack = []Value{}
	th.args(f, a)
	Assert(t).That(fmt.Sprint(th.stack), Equals("[123]"))
}

func makeOb() *SuObject {
	var ob SuObject
	ob.Add(SuInt(11))
	ob.Add(SuInt(22))
	ob.Put(SuStr("a"), SuStr("aa"))
	ob.Put(SuStr("b"), SuStr("bb"))
	return &ob
}
