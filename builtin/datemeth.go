package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	DateMethods = Methods{
		"MinusSeconds": method1("(date)", func(this Value, val Value) Value {
			t1 := this.(SuDate)
			t2 := val.(SuDate) //TODO better conversion error
			ms := t1.MinusMs(t2)
			return fromFloat(float64(ms)/1000)
		}),
	}
}
