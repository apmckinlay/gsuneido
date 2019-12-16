// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"os"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Getenv(string)",
	func(arg Value) Value {
		return SuStr(os.Getenv(ToStr(arg)))
	})
