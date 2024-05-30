// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
)

var _ = builtin(Hash, "(value)")

func Hash(arg Value) Value {
	return Int64Val(int64(arg.Hash()))
}
