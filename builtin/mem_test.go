package builtin

import (
	"runtime"
	"runtime/metrics"
	"testing"
)

func BenchmarkMetric(b *testing.B) {
	sample := make([]metrics.Sample, 1)
	sample[0].Name = "/memory/classes/heap/objects:bytes"
	for range b.N {
		metrics.Read(sample)
	}
}

func BenchmarkMemStat(b *testing.B) {
	var ms runtime.MemStats
	for range b.N {
		runtime.ReadMemStats(&ms)
	}
}

func BenchmarkSystemMemory(b *testing.B) {
	for range b.N {
		systemMemory()
	}
}

// func TestHeapSize(*testing.T) {
// 	sample := []metrics.Sample{{Name: "/memory/classes/heap/objects:bytes"}}
// 	metrics.Read(sample)
// 	fmt.Println("metrics /memory/classes/heap/objects:bytes",
// 		int64(sample[0].Value.Uint64()))

// 	var ms runtime.MemStats
// 	runtime.ReadMemStats(&ms)
// 	fmt.Println("MemStats HeapAlloc", ms.HeapAlloc)
// }
