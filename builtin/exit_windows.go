// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build win32

package builtin

import (
	"os"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Exit(code = 0)",
	func(arg Value) Value {
		if arg == True {
			os.Exit(0)
		}
		PostQuitMessage(uintptr(IfInt(arg)))
		return nil
	})
