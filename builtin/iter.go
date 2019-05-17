package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	IterMethods = Methods{
		"Dup": method0(func(this Value) Value {
			it := this.(SuIter)
			return SuIter{Iter: it.Dup()}
		}),
		"Infinite?": method0(func(this Value) Value {
			it := this.(SuIter)
			return SuBool(it.Infinite())
		}),
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
