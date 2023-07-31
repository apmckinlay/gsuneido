// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(Construct, "(what, suffix='')")

func Construct(th *Thread, args []Value) Value {
	what := args[0]
	suffix := ToStr(args[1])
	c, ok := what.ToContainer()
	if ok {
		if c.ListSize() < 1 {
			panic("Construct: object requires member 0")
		}
		what = c.ListGet(0)
	}
	if s, ok := what.ToStr(); ok {
		if !strings.HasSuffix(s, suffix) {
			s += suffix
		}
		what = Global.GetName(th, s)
	}
	if c == nil {
		return th.CallLookup(what, "*new*")
	}
	return th.CallLookupEach1(what, "*new*", c)
}
