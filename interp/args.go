package interp

import (
	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/verify"
)

var nilValues [MaxArgs]Value

func (t *Thread) args(fn *SuFunc, as ArgSpec) int {
	nargs := as.Nargs()
	base := len(t.stack) - nargs
	// expand stack to allow for locals (including params)
	if expand := fn.Nlocals - nargs; expand > 0 {
		t.stack = append(t.stack, nilValues[:expand]...)
	}
	locals := t.stack[base:]
	// shrink stack if excess args (locals still has full args)
	if nargs > fn.Nlocals {
		t.stack = t.stack[:base+fn.Nlocals]
	}
	t.massage(fn, as, locals)
	return base
}

// args massages the arguments on the stack (specified by ArgSpec)
// to match what is expected by the function (specified by SuFunc)
// The stack must already have been expanded.
func (t *Thread) massage(fn *SuFunc, as ArgSpec, args []Value) {
	unnamed := int(as.Unnamed)
	if unnamed == fn.Nparams && len(as.Spec) == 0 {
		return // simple fast path
	}
	if unnamed < EACH && fn.Flags[0] != AT_F && unnamed > fn.Nparams {
		panic("too many arguments")
	}
	// as.Unnamed < fn.Nparams

	atParam := fn.Nparams == 1 && fn.Flags[0] == AT_F

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
		ob := &SuObject{}
		for i := 0; i < unnamed; i++ {
			ob.Add(args[i])
			args[i] = nil
		}
		for i, ni := range as.Spec {
			ob.Put(SuStr(as.Names[ni]), args[unnamed+i])
			args[unnamed+i] = nil
		}
		args[0] = ob
		return
	}

	if unnamed >= EACH {
		// @args => params
		ob := args[0].(*SuObject)
		for i := 0; i < ints.Min(fn.Nparams, ob.Vsize()); i++ {
			args[i] = ob.Vget(i)
		}
		// named members may overwrite unnamed (same as when passed individually)
		for i := 0; i < fn.Nparams; i++ {
			if x := ob.Get(SuStr(fn.Strings[i])); x != nil {
				args[i] = x
			}
		}
	}

	if len(as.Spec) > 0 {
		// shuffle named args to match params
		verify.That(len(as.Spec) < MaxArgs)
		var tmp [MaxArgs]Value
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

	// fill in dynamic
	for i := 0; i < fn.Nparams; i++ {
		if args[i] == nil && fn.Flags[i]&DYN_F != 0 {
			if x := t.dyn("_" + fn.Strings[i]); x != nil {
				args[i] = x
			}
		}
	}

	// fill in defaults and check for missing
	v := 0
	for i := int(as.Unnamed); i < fn.Nparams; i++ {
		if args[i] == nil {
			if i >= fn.Nparams-fn.Ndefaults {
				args[i] = fn.Values[v]
				v++
			} else {
				panic("missing argument")
			}
		}
	}

}

func (t *Thread) dyn(name string) Value {
	for i := len(t.frames) - 1; i >= 0; i-- {
		fr := t.frames[i]
		for j, s := range fr.fn.Strings {
			if s == name {
				return fr.locals[j]
			}
		}
	}
	return nil
}
