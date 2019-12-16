// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import "github.com/apmckinlay/gsuneido/util/verify"

// Call pushes the arguments onto the stack and calls the function
func (t *Thread) Call(fn Callable, args ...Value) Value {
	return t.CallThis(fn, nil, args...)
}

func (t *Thread) CallEach1(fn Callable, arg Value) Value {
	return t.pushCall(fn, nil, &ArgSpecEach1, arg)
}

// CallLookup calls a *named* method.
func (t *Thread) CallLookup(this Value, method string, args ...Value) Value {
	return t.CallThis(t.lookup(this, method), this, args...)
}

func (t *Thread) CallThis(fn Callable, this Value, args ...Value) Value {
	verify.That(len(args) < AsEach)
	as := &StdArgSpecs[len(args)]
	return t.pushCall(fn, this, as, args...)
}

func (t *Thread) CallLookupEach1(this Value, method string, arg Value) Value {
	return t.pushCall(t.lookup(this, method), this, &ArgSpecEach1, arg)
}

func (t *Thread) lookup(this Value, method string) Callable {
	fn := this.Lookup(t, method)
	if fn == nil {
		panic("method not found: " + ErrType(this) + "." + method)
	}
	return fn
}

// pushCall pushes the arguments onto the stack and calls the function
func (t *Thread) pushCall(fn Callable, this Value, as *ArgSpec, args ...Value) Value {
	base := t.sp
	for _, x := range args {
		t.Push(x)
	}
	result := fn.Call(t, this, as)
	t.sp = base
	return result
}
