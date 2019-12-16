// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"fmt"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin("PrintStdout(string)",
	func(t *Thread, args []Value) Value {
		fmt.Print(ToStr(args[0]))
		return nil
	})
