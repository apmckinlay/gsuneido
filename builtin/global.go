// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
)

var _ = builtin("Global(name)", func(t *Thread, args []Value) Value {
	s := ToStr(args[0])
	global, s := str.Cut(s, '.')
	val := Global.GetName(t, global)
	for s != "" {
		var mem string
		mem, s = str.Cut(s, '.')
		if val = val.Get(t, SuStr(mem)); val == nil {
			panic("Global: " + mem + " not found")
		}
	}
	return val
})
