package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	RecordMethods = Methods{
		"Copy": method0(func(this Value) Value {
			return this.(*SuRecord).Copy()
		}),
		"Observer": method1("(observer)", func(this, arg Value) Value {
			this.(*SuRecord).Observer(arg)
			return nil
		}),
		"PreSet": method2("(key,value)", func(this, arg1, arg2 Value) Value {
			this.(*SuRecord).PreSet(arg1, arg2)
			return nil
		}),
		"RemoveObserver": method1("(observer)", func(this, arg Value) Value {
			return SuBool(this.(*SuRecord).RemoveObserver(arg))
		}),
	}
}
