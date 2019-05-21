package builtin

import (
	"runtime"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/dnum"
)

var _ = builtin0("MemoryArena()", func() Value {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return SuDnum{Dnum: dnum.FromInt(int64(ms.HeapSys))}
})
