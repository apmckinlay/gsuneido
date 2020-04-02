// +build !portable windows

package builtin

import (
	"testing"

	heap "github.com/apmckinlay/gsuneido/builtin/heapstack"
	. "github.com/apmckinlay/gsuneido/runtime"
)

var result Value

func BenchmarkGetStr(b *testing.B) {
	const n = 100
	buf := heap.Alloc(n)
	var s SuStr
	for i := 0; i < b.N; i++ {
		s = SuStr(heap.GetStrN(buf, n))
	}
	result = s
}

func BenchmarkBufToStr(b *testing.B) {
	const n = 100
	buf := heap.Alloc(n)
	var s Value
	for i := 0; i < b.N; i++ {
		s = bufRet(buf, n)
	}
	result = s
}
