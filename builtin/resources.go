// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin0("ResourceCounts()", func() Value {
	ob := NewSuObject()
	ob.Put(nil, SuStr("Callbacks"), IntVal(CallbacksCount()))
	return ob
})
