package runtime

// see also: ArgSpec

import (
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/verify"
)

func (t *Thread) Args(ps *ParamSpec, as *ArgSpec) []Value {
	if ps.Signature^as.Signature == 0xff {
		// fast path if signatures match, hopefully inlined
		return t.stack[t.sp-int(as.Nargs):]
	}
	return t.args(ps, as)
}

// args adjusts sp to shrink or grow the stack
// so the correct amount is reserved for the function
// and then calls massage
// It returns a slice of the stack containing the locals
func (t *Thread) args(ps *ParamSpec, as *ArgSpec) []Value {
	nargs := int(as.Nargs)
	base := t.sp - nargs

	// reserve stack space for params
	for expand := int(ps.Nparams) - nargs; expand > 0; expand-- {
		t.Push(nil)
	}
	locals := t.stack[base:]
	t.massage(ps, as, locals)

	// shrink stack if excess args
	t.sp = base + int(ps.Nparams)

	return locals
}

// massage adjust the arguments on the stack (described by ArgSpec)
// to match what is expected by the function (described by Func)
// The stack must already have been expanded (e.g. by args)
func (t *Thread) massage(ps *ParamSpec, as *ArgSpec, args []Value) {
	unnamed := as.Unnamed()
	atParam := ps.Nparams == 1 && ps.Flags[0] == AtParam
	if unnamed == int(ps.Nparams) && len(as.Spec) == 0 && !atParam {
		return // simple fast path
	}
	atArg := as.Each >= EACH
	if !atArg && !atParam && unnamed > int(ps.Nparams) {
		panic("too many arguments")
	}
	// as.Unnamed < fn.Nparams

	each := int(as.Each) - 1

	// remove after debugged
	verify.That(!atParam || ps.Nparams == 1)
	verify.That(!atArg || len(as.Spec) == 0)

	if atParam {
		if atArg {
			// @arg => @param
			ob := ToObject(args[0])
			ob = ob.Slice(each) // makes a copy
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
			ob.Put(as.Names[ni], args[unnamed+i])
			args[unnamed+i] = nil
		}
		args[0] = ob
		return
	}

	if atArg {
		// @args => params
		ob := args[0].(*SuObject)
		if ob.ListSize()-each > int(ps.Nparams) {
			panic("too many arguments")
		}
		for i := 0; i < ints.Min(int(ps.Nparams), ob.ListSize()-each); i++ {
			args[i] = ob.ListGet(i + each)
		}
		// named members may overwrite unnamed (same as when passed individually)
		for i := 0; i < int(ps.Nparams); i++ {
			if x := ob.Get(t, SuStr(ps.Names[i])); x != nil {
				args[i] = x
			}
		}
	}

	if len(as.Spec) > 0 {
		// shuffle named args to match params
		verify.That(len(as.Spec) < MaxArgs)
		var tmp [MaxArgs]Value
		nargs := int(as.Nargs)
		// move named arguments aside, off the stack
		copy(tmp[0:], args[unnamed:nargs])
		// initialize space for named args
		for i := unnamed; i < nargs; i++ {
			args[i] = nil
		}
		// move applicable named args back to correct position
		for si, ni := range as.Spec {
			for i := 0; i < int(ps.Nparams); i++ {
				if as.Names[ni] == SuStr(ps.Names[i]) {
					args[i] = tmp[si]
				}
			}
		}
	}

	// fill in dynamic
	for i := 0; i < int(ps.Nparams); i++ {
		if args[i] == nil && ps.Flags[i]&DynParam != 0 {
			if x := t.dyn("_" + ps.Names[i]); x != nil {
				args[i] = x
			}
		}
	}

	// fill in defaults and check for missing
	v := 0
	for i := as.Unnamed(); i < int(ps.Nparams); i++ {
		if args[i] == nil {
			if i >= int(ps.Nparams-ps.Ndefaults) {
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
