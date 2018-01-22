package interp

import (
	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// args massages the arguments on the stack (specified by ArgSpec)
// to match what is expected by the function (specified by SuFunc)
func (t *Thread) args(fn *SuFunc, as ArgSpec) {
	if int(as.Unnamed) == fn.Nparams {
		if len(as.Spec) > 0 {
			// remove unused named args from stack
			panic("not implemented") // TODO
		}
		return // simple fast path
	}
	if int(as.Unnamed) > fn.Nparams {
		panic("too many arguments")
	}
	// as.Unnamed < fn.Nparams

	atParam := fn.Nparams == 1 && fn.Strings[0][0] == '@'

	// remove after debugged
	verify.That(!atParam || fn.Nparams == 1)
	verify.That(as.Unnamed < EACH || len(as.Spec) == 0)

	if atParam {
		if as.Unnamed >= EACH {
			// @arg => @param
			ob := t.Top().(*SuObject)
			ob = ob.Slice(int(as.Unnamed - EACH))
			t.SetTop(ob)
			return
		}
		// args => @param
		panic("not implemented") // TODO
	}

	if as.Unnamed >= EACH {
		// @args
		panic("not implemented") // TODO
	}

	if len(as.Spec) > 0 {
		// shuffle named args to match params
	}

}
