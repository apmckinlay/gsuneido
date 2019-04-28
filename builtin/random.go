package builtin

import (
	"math/rand"

	. "github.com/apmckinlay/gsuneido/runtime"
)

var _ = builtin1("Random(limit)", func(arg Value) Value {
	limit := IfInt(arg)
	return SuInt(rand.Intn(limit))
})
