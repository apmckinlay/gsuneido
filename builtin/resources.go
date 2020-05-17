// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin0("ResourceCounts()", func() Value {
	ob := NewSuObject()
	add(ob, "File", nFile)
	add(ob, "RunPiped", nRunPiped)
	add(ob, "SocketClient", nSocketClient)
	add(ob, "Callbacks", CallbacksCount())
	add(ob, "WndProcs", WndProcCount())
	return ob
})

func add(ob *SuObject, name string, n int) {
	if n != 0 {
		ob.Set(SuStr(name), IntVal(n))
	}
}
