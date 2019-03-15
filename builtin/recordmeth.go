package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	RecordMethods = Methods{
		"Copy": method0(func(this Value) Value {
			return this.(*SuRecord).Copy()
		}),
	}
}
