// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin(Timestamp, "()")

func Timestamp(th *Thread, args []Value) Value {
	return th.Timestamp()
}
