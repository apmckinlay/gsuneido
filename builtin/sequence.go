// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtinRaw("Sequence(iter)", // raw to get thread
	func(t *Thread, as *ArgSpec, args []Value) Value {
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

func (wi *wrapIter) Infinite() (result bool) {
	for i := 0; i < 1; i++ { // workaround for 1.14 bug
		defer func() {
			if e := recover(); e != nil {
				result = false
			}
		}()
	}
	return wi.call("Infinite?") == True
}

func (wi *wrapIter) Dup() Iter {
	it := wi.call("Dup")
	return &wrapIter{iter: it, t: wi.t}
}

func (wi *wrapIter) call(method string) Value {
	return wi.t.CallLookup(wi.iter, method)
}

var _ Iter = (*wrapIter)(nil)

// for SuSequence

func init() {
	SequenceMethods = Methods{
		"Copy": method0(func(this Value) Value {
			return this.(*SuSequence).Copy()
		}),
		"Infinite?": method0(func(this Value) Value {
			return SuBool(this.(*SuSequence).Infinite())
		}),
		"Instantiated?": method0(func(this Value) Value {
			return SuBool(this.(*SuSequence).Instantiated())
		}),
		"Iter": method0(func(this Value) Value {
			iter := this.(*SuSequence).Iter()
			if wi, ok := iter.(*wrapIter); ok {
				return wi.iter
			}
			return SuIter{Iter: iter}
		}),
		"Join": method1("(separator='')", func(this, arg Value) Value {
			iter := this.(*SuSequence).Iter()
			separator := ToStr(arg)
			sep := ""
			var buf strings.Builder
			for {
				val := iter.Next()
				if val == nil {
					break
				}
				buf.WriteString(sep)
				sep = separator
				if s, ok := val.ToStr(); ok {
					buf.WriteString(s)
				} else {
					buf.WriteString(val.String())
				}
			}
			return SuStr(buf.String())
		}),
	}
}
