// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

// cause a Go runtime error (for testing)
var _ = builtin0("RuntimeError()",
	func() Value {
		var x []Value
		return x[123]
	})
