// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !windows portable

package builtin

import (
	"os"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Exit(code = 0)",
	func(arg Value) Value {
		code := 0
		if arg != True {
			code = IfInt(arg)
		}
		os.Exit(code)
		return nil
	})
