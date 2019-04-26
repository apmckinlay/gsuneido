package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Record(@args)",
	func(_ *Thread, args ...Value) Value {
		ob := args[0].(*SuObject)
		ob.SetDefault(EmptyStr)
		return SuRecordFromObject(ob)
	})
