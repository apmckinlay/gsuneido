package interp

import (
	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/util/verify"
)

const maxNamedArgs = 100

// args massages the arguments on the stack (specified by ArgSpec)
// to match what is expected by the function (specified by SuFunc)
// The stack must already have been expanded.
func (t *Thread) args(fn *SuFunc, as ArgSpec, args []Value) {
	unnamed := int(as.Unnamed)
	if unnamed == fn.Nparams {
		if len(as.Spec) > 0 {
			// remove unused named args from stack
			panic("not implemented") // TODO
		}
		return // simple fast path
	}
	if unnamed > fn.Nparams {
		panic("too many arguments")
	}
	// as.Unnamed < fn.Nparams

	atParam := fn.Nparams == 1 && fn.Strings[0][0] == '@'

	// remove after debugged
	verify.That(!atParam || fn.Nparams == 1)
	verify.That(unnamed < EACH || len(as.Spec) == 0)

	if atParam {
		if unnamed >= EACH {
			// @arg => @param
			ob := args[0].(*SuObject)
			ob = ob.Slice(unnamed - EACH)
			args[0] = ob
			return
		}
		// args => @param
		panic("not implemented") // TODO
	}

	if unnamed >= EACH {
		// @args
		panic("not implemented") // TODO
	}

	if len(as.Spec) > 0 {
		// shuffle named args to match params
		verify.That(len(as.Spec) < maxNamedArgs)
		var tmp [maxNamedArgs]Value
		nargs := as.Nargs()
		// move named arguments aside, off the stack
		copy(tmp[0:], args[unnamed:nargs])
		// initialize space for named args
		for i := unnamed; i < nargs; i++ {
			args[i] = nil
		}
		// move applicable named args back to correct position
		for si, ni := range as.Spec {
			for i := 0; i < fn.Nparams; i++ {
				if as.Names[ni] == fn.Strings[i] {
					args[i] = tmp[si]
				}
			}
		}
	}

}
