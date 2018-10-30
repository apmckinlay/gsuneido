package runtime

// see also: ArgSpec

import (
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// args adjusts sp to shrink or grow the stack
// so the correct amount is reserved for the function
// and then calls massage
// It returns the
func (t *Thread) Args(ps *ParamSpec, as *ArgSpec) []Value {
	nargs := as.Nargs()
	base := t.sp - nargs

	// reserve stack space for params
	for expand := ps.Nparams - nargs; expand > 0; expand-- {
		t.Push(nil)
	}
	locals := t.stack[base:]
	t.massage(ps, as, locals)

	// shrink stack if excess args
	t.sp = base + ps.Nparams

	return locals
}

// massage adjust the arguments on the stack (described by ArgSpec)
// to match what is expected by the function (described by Func)
// The stack must already have been expanded (e.g. by args)
func (t *Thread) massage(ps *ParamSpec, as *ArgSpec, args []Value) {
	unnamed := int(as.Unnamed)
	atParam := ps.Nparams == 1 && ps.Flags[0] == AtParam
	if unnamed == ps.Nparams && len(as.Spec) == 0 && !atParam {
		return // simple fast path
	}
	if unnamed < EACH && ps.Flags[0] != AtParam && unnamed > ps.Nparams {
		panic("too many arguments")
	}
	// as.Unnamed < fn.Nparams

	atArg := unnamed >= EACH

	// remove after debugged
	verify.That(!atParam || ps.Nparams == 1)
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
		for i := 0; i < ints.Min(ps.Nparams, ob.ListSize()); i++ {
			args[i] = ob.ListGet(i)
		}
		// named members may overwrite unnamed (same as when passed individually)
		for i := 0; i < ps.Nparams; i++ {
			if x := ob.Get(SuStr(ps.Names[i])); x != nil {
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
			for i := 0; i < ps.Nparams; i++ {
				if as.Names[ni] == ps.Names[i] {
					args[i] = tmp[si]
				}
			}
		}
	}

	// fill in dynamic
	for i := 0; i < ps.Nparams; i++ {
		if args[i] == nil && ps.Flags[i]&DynParam != 0 {
			if x := t.dyn("_" + ps.Names[i]); x != nil {
				args[i] = x
			}
		}
	}

	// fill in defaults and check for missing
	v := 0
	for i := int(as.Unnamed); i < ps.Nparams; i++ {
		if args[i] == nil {
			if i >= ps.Nparams-ps.Ndefaults {
				args[i] = ps.Values[v]
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
		for j, s := range fr.fn.Names {
			if s == name {
				return t.stack[fr.bp+j]
			}
		}
	}
	return nil
}
