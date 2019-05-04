package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtinRaw("Sequence(iter)", // raw to get thread
	func(t *Thread, as *ArgSpec, args ...Value) Value {
		args = t.Args(&ParamSpec1, as)
		return NewSuSequence(&wrapIter{iter: args[0], t: t})
	})

// wrapIter adapts a Suneido iterator (a class with Next,Dup,Infinite)
// to the runtime.Iter interface
// for the reverse see runtime.SuIter
type wrapIter struct {
	iter Value
	t    *Thread
}

func (wi *wrapIter) Next() Value {
	x := wi.call("Next")
	if x == wi.iter {
		return nil
	}
	return x
}

func (wi *wrapIter) Infinite() bool {
	if wi.iter.Lookup(wi.t, "Infinite?") == nil {
		return false
	}
	return wi.call("Infinite?") == True
}

func (wi *wrapIter) Dup() Iter {
	it := wi.call("Dup")
	return &wrapIter{iter: it, t: wi.t}
}

func (wi *wrapIter) call(method string) Value {
	wi.t.Push(wi.iter)
	return wi.t.CallMethod(method, ArgSpec0)
}

var _ Iter = (*wrapIter)(nil)
