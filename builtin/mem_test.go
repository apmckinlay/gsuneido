package builtin

import (
	"runtime"
	"runtime/metrics"
	"testing"
)

func BenchmarkMetric(b *testing.B) {
	sample := make([]metrics.Sample, 1)
	sample[0].Name = "/memory/classes/heap/objects:bytes"
	for b.Loop() {
		metrics.Read(sample)
	}
}

func BenchmarkMemStat(b *testing.B) {
	var ms runtime.MemStats
	for b.Loop() {
		runtime.ReadMemStats(&ms)
	}
}

func BenchmarkSystemMemory(b *testing.B) {
	for b.Loop() {
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
