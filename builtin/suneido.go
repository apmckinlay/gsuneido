// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	SuneidoObjectMethods = Methods{
		"Compile": method("(source, errob = false)",
			func(t *Thread, _ Value, args []Value) Value {
				src := ToStr(args[0])
				if args[1] == False {
					return compile.Constant(src)
				}
				ob := ToContainer(args[1])
				val, checks := compile.Checked(t, src)
				for _, w := range checks {
					ob.Add(SuStr(w))
				}
				return val
			}),
		"Parse": method("(source)",
			func(t *Thread, _ Value, args []Value) Value {
				src := ToStr(args[0])
				p := compile.AstParser(src)
				return p.Const()
			}),
	}
}
