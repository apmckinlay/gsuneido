// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin0("Crash!()",
	func() runtime.Value {
		// force a crash, mostly to test output capture
		go func() { panic("Crash!") }()
		return nil
	})
