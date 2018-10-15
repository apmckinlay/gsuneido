package interp

// see also: ArgSpec

import (
	. "github.com/apmckinlay/gsuneido/base"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/verify"
)

func (t *Thread) args(fn *Func, as *ArgSpec) []Value {
	nargs := as.Nargs()
	base := t.sp - nargs
	// allow stack space for params
	expand := fn.Nparams - nargs
	for ; expand > 0; expand-- {
		t.Push(nil)
	}
	locals := t.stack[base:]
	// shrink stack if excess args (locals still has full args)
	if nargs > fn.Nparams {
		t.sp = base + fn.Nparams
	}
	t.massage(fn, as, locals)
	return t.stack[base:]
}

// args massages args on the stack (described by ArgSpec)
// to match what is expected by the function (described by Func)
// The stack must already have been expanded (e.g. by args)
func (t *Thread) massage(fn *Func, as *ArgSpec, args []Value) {
	unnamed := int(as.Unnamed)
	if unnamed == fn.Nparams && len(as.Spec) == 0 {
		return // simple fast path
	}
	if unnamed < EACH && fn.Flags[0] != AtParam && unnamed > fn.Nparams {
		panic("too many arguments")
	}
	// as.Unnamed < fn.Nparams

	atArg := unnamed >= EACH
	atParam := fn.Nparams == 1 && fn.Flags[0] == AtParam

	// remove after debugged
	verify.That(!atParam || fn.Nparams == 1)
	verify.That(unnamed < EACH || len(as.Spec) == 0)

	if atParam {
		if atArg {
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

	if atArg {
		// @args => params
		ob := args[0].(*SuObject)
		for i := 0; i < ints.Min(fn.Nparams, ob.ListSize()); i++ {
			args[i] = ob.ListGet(i)
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
		if args[i] == nil && fn.Flags[i]&DynParam != 0 {
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
	for i := t.fp - 1; i >= 0; i-- {
		fr := t.frames[i]
		for j, s := range fr.fn.Strings {
			if s == name {
				return t.stack[fr.bp + j]
			}
		}
	}
	return nil
}
