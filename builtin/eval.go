// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/regex"
)

// Eval executes string containing Suneido code
// i.e. string.Eval()
func EvalString(t *Thread, s string) Value {
	s = strings.Trim(s, " \t\r\n")
	if isGlobal(s) {
		// optimize if just a global name
		return Global.GetName(t, s)
	}
	s = "function () {\n" + s + "\n}"
	fn := compile.NamedConstant("", "eval", s).(*SuFunc)
	return t.Start(fn, nil)
}

var rxGlobal = regex.Compile(`\A[A-Z][_a-zA-Z0-9]*?[!?]?\Z`)

func isGlobal(s string) bool {
	return rxGlobal.Matches(s)
}

// EvalAsMethod runs a function as if it were a method of an object
// i.e. object.Eval
func EvalAsMethod(t *Thread, as *ArgSpec, ob Value, args []Value) Value {
	// first argument is function
	k, f := NewArgsIter(as, args)()
	if k != nil || f == nil {
		panic("usage: object.Eval(callable, ...)")
	}
	if m, ok := f.(*SuMethod); ok {
		f = m.GetFn()
	}
	return f.Call(t, ob, as.DropFirst())
}
