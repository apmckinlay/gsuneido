// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/str"
)

var _ = builtin(global, "(name)")

func global(th *Thread, args []Value) Value {
	s := ToStr(args[0])
	global, s := str.Cut(s, '.')
	val := Global.GetName(th, global)
	for s != "" {
		var mem string
		mem, s = str.Cut(s, '.')
		if val = val.Get(th, SuStr(mem)); val == nil {
			panic("Global: " + mem + " not found")
		}
	}
	return val
}
