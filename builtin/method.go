// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/core"

var _ = builtin(Method, "(object, methodName)")

func Method(th *Thread, args []Value) Value {
	ob := args[0]
	method := ToStr(args[1])

	val := ob.Lookup(th, method)
	if val == nil {
		return False
	}
	return NewSuMethod(ob, val)
}
