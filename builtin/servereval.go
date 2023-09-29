// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(ServerEval, "(@args)")

func ServerEval(th *Thread, args []Value) Value {
	return th.Dbms().Exec(th, args[0])
}
