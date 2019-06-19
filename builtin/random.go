package builtin

import (
	"math/rand"
	"time"

	. "github.com/apmckinlay/gsuneido/runtime"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = builtin1("Random(limit)", func(arg Value) Value {
	limit := IfInt(arg)
	return SuInt(rand.Intn(limit))
})
