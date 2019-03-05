package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	IterMethods = Methods{
		"Next": method0(func(this Value) Value {
			it := this.(SuIter)
			next := it.Next()
			if next == nil {
				return this
			}
			return next
		}),
	}
}
