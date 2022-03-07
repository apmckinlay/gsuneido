// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

// see also: ArgSpec

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
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
	for i := 0; i < nargs; i++ {
		if locals[i] == nil {
			panic("missing argument " + strconv.Itoa(i) + " in " + as.String())
		}
	}
	t.massage(ps, as, locals)

	// shrink stack if excess args
	t.sp = base + int(ps.Nparams)

	return locals
}

// MaxArgs is the maximum number of arguments allowed
const MaxArgs = 200

// massage adjust the arguments on the stack (described by ArgSpec)
// to match what is expected by the function (described by ParamSpec)
// The stack must already have been expanded (e.g. by args)
func (t *Thread) massage(ps *ParamSpec, as *ArgSpec, args []Value) {
	unnamed := int(as.Nargs) - len(as.Spec) // only valid if !atArg
	atParam := ps.Nparams == 1 && ps.Flags[0] == AtParam
	atArg := as.Each >= EACH0
	if unnamed == int(ps.Nparams) && len(as.Spec) == 0 && !atParam && !atArg {
		return // simple fast path
	}
	if !atArg && !atParam && unnamed > int(ps.Nparams) {
		panic("too many arguments")
	}
	// as.Unnamed < fn.Nparams

	each := int(as.Each) - 1

	// remove after debugged
	assert.That(!atParam || ps.Nparams == 1)
	assert.That(!atArg || len(as.Spec) == 0)

	if atParam {
		if atArg {
			// @arg => @param
			ob := ToContainer(args[0]).ToObject()
			args[0] = ob.Slice(each) // makes a copy
			return
		}
		// args => @param
		ob := &SuObject{}
		for i := 0; i < unnamed; i++ {
			ob.Add(args[i])
			args[i] = nil
		}
		for i, ni := range as.Spec {
			ob.Set(as.Names[ni], args[unnamed+i])
			args[unnamed+i] = nil
		}
		args[0] = ob
		return
	}
	if atArg {
		// @args => params
		ob := ToContainer(args[0])
		args[0] = nil
		if ob.ListSize()-each > int(ps.Nparams) {
			panic("too many arguments")
		}
		for i := 0; i < ord.Min(int(ps.Nparams), ob.ListSize()-each); i++ {
			args[i] = ob.ListGet(i + each)
		}
		// named members may overwrite unnamed (same as when passed individually)
		for i := 0; i < int(ps.Nparams); i++ {
			if x := ob.GetIfPresent(t, SuStr(ps.ParamName(i))); x != nil {
				args[i] = x
			}
		}
	} else if len(as.Spec) > 0 {
		// shuffle named args to match params
		assert.That(len(as.Spec) < MaxArgs)
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
				if as.Names[ni] == SuStr(ps.ParamName(i)) {
					args[i] = tmp[si]
				}
			}
		}
	}

	// fill in dynamic
	for i := 0; i < int(ps.Nparams); i++ {
		if args[i] == nil && ps.Flags[i]&DynParam != 0 {
			if x := t.dyn("_" + ps.ParamName(i)); x != nil {
				args[i] = x
			}
		}
	}

	// fill in defaults
	noDefs := int(ps.Nparams - ps.Ndefaults)
	for i := noDefs; i < int(ps.Nparams); i++ {
		if args[i] == nil {
			args[i] = ps.Values[i-noDefs]
		}
	}

	// check that all parameters now have values
	for i := 0; i < noDefs; i++ {
		if args[i] == nil {
			panic("missing argument")
		}
	}
}

func (t *Thread) dyn(name string) Value {
	for i := t.fp - 1; i >= 0; i-- {
		fr := t.frames[i]
		for j, s := range fr.fn.Names {
			if s == name {
				if x := fr.locals.v[j]; x != nil {
					return x
				}
			}
		}
	}
	return nil
}
