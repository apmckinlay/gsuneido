// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import "github.com/apmckinlay/gsuneido/util/assert"

// Call pushes the arguments onto the stack and calls the function
func (th *Thread) Call(fn Callable, args ...Value) Value {
	return th.CallThis(fn, nil, args...)
}

func (th *Thread) CallEach(fn Callable, arg Value) Value {
	return th.PushCall(fn, nil, &ArgSpecEach0, arg)
}

func (th *Thread) CallEach1(fn Callable, arg Value) Value {
	return th.PushCall(fn, nil, &ArgSpecEach1, arg)
}

// CallLookup calls a *named* method.
func (th *Thread) CallLookup(this Value, method string, args ...Value) Value {
	return th.CallThis(th.Lookup(this, method), this, args...)
}

func (th *Thread) CallThis(fn Callable, this Value, args ...Value) Value {
	assert.That(len(args) < AsEach)
	as := &StdArgSpecs[len(args)]
	return th.PushCall(fn, this, as, args...)
}

func (th *Thread) CallLookupEach1(this Value, method string, arg Value) Value {
	return th.PushCall(th.Lookup(this, method), this, &ArgSpecEach1, arg)
}

func (th *Thread) Lookup(this Value, method string) Callable {
	fn := this.Lookup(th, method)
	if fn == nil {
		panic("method not found: " + ErrType(this) + "." + method)
	}
	return fn
}

// PushCall pushes the arguments onto the stack and calls the function
func (th *Thread) PushCall(fn Callable, this Value, as *ArgSpec, args ...Value) Value {
	base := th.sp
	for _, x := range args {
		th.Push(x)
	}
	result := fn.Call(th, this, as)
	th.sp = base
	return result
}
