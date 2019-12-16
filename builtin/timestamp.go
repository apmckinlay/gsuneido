// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var prevTimestamp SuDate

var _ = builtin("Timestamp()", func(t *Thread, args []Value) Value {
	return t.Dbms().Timestamp()
})
