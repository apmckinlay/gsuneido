// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

//TODO add some kind of registry of name->func

var _ = builtin0("ResourceCounts()", func() Value {
	ob := NewSuObject()
	ob.Put(nil, SuStr("Callbacks"), IntVal(CallbacksCount()))
	return ob
})
