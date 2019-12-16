// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("ServerEval(@args)", func(t *Thread, args []Value) Value {
	return t.Dbms().Exec(t, args[0])
})
