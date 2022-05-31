// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

var _ = builtin0("Crash!()",
	func() runtime.Value {
		// force a crash, mostly to test output capture
		go func() { panic("Crash!") }()
		return nil
	})

var _ = builtin0("AssertFail()",
	func() runtime.Value {
		assert.That(false)
		return nil
	})

var _ = builtin0("BoundsFail()",
	func() runtime.Value {
		return []runtime.Value{}[1]
	})
