// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Timestamp, "() :date")

func Timestamp(th *Thread, args []Value) Value {
	return th.Timestamp()
}
