package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin0("Libraries()", func() Value {
	ob := &SuObject{}
	ob.Add(SuStr("stdlib")) //TODO dbms.Libraries
	return ob
})
