package builtin

import . "github.com/apmckinlay/gsuneido/runtime"

var _ = builtin("Record(@args)",
	func(_ *Thread, args ...Value) Value {
		return &SuRecord{SuObject: *args[0].(*SuObject)}
	})
