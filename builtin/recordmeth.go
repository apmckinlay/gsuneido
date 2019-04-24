package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	RecordMethods = Methods{
		"Copy": method0(func(this Value) Value {
			return this.(*SuRecord).Copy()
		}),
		"Observer": method1("(observer)", func(this Value, arg Value) Value {
			this.(*SuRecord).Observer(arg)
			return nil
		}),
		"RemoveObserver": method1("(observer)", func(this Value, arg Value) Value {
			return SuBool(this.(*SuRecord).RemoveObserver(arg))
		}),
	}
}
