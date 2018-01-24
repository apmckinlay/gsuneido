package interp

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestArgs(t *testing.T) {
	th := &Thread{}
	f := &SuFunc{}
	a := ArgSpec{}

	// 0 args => 0 params
	th.args(f, a)

	// @arg => @param
	f.Nparams, f.Nlocals = 1, 1
	f.Flags = []Flag{AT_F}
	a.Unnamed = EACH
	th.stack = []Value{&SuObject{}}
	th.args(f, a)

	// @+1arg => @param
	f.Nparams, f.Nlocals = 1, 1
	f.Flags = []Flag{AT_F}
	a.Unnamed = EACH1
	th.stack = []Value{&SuObject{}}
	th.args(f, a)

	// 2 args => 2 params
	f.Nparams, f.Nlocals = 2, 2
	a.Unnamed = 2
	th.stack = []Value{SuInt(11), SuInt(22)}
	th.args(f, a)

	// 1 args => 2 params
	f.Nparams, f.Nlocals = 2, 2
	a.Unnamed = 1
	th.stack = []Value{SuInt(11)}
	Assert(t).That(func() { th.args(f, a) }, Panics("missing argument"))

	// 1 args => 2 params with 1 default
	a.Unnamed = 1
	th.stack = []Value{SuInt(11)}
	f.Nparams, f.Nlocals = 2, 2
	f.Ndefaults = 1
	f.Values = []Value{SuInt(22)}
	th.args(f, a)
	Assert(t).That(fmt.Sprint(th.stack), Equals("[11 22]"))

	// all named
	a.Unnamed = 0
	a.Names = []string{"c", "b", "a", "d"}
	a.Spec = []byte{1, 0, 2, 3} // b, c, a, d
	th.stack = []Value{SuInt(22), SuInt(33), SuInt(11), SuInt(44)}
	f.Nparams, f.Nlocals = 3, 3
	f.Strings = []string{"a", "b", "c"}
	th.args(f, a)
	Assert(t).That(fmt.Sprint(th.stack), Equals("[11 22 33]"))
}
