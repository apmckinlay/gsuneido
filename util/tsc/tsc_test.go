// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tsc

import (
	"fmt"
	"testing"
	"time"
)

func TestRead(*testing.T) {
	for i := 1; i < 10; i++ {
		t := time.Now()
		tsc := Read()
		time.Sleep(time.Duration(i) * time.Millisecond)
		dt := time.Since(t) / time.Microsecond
		dtsc := Read() - tsc
		fmt.Println(dtsc / uint64(dt))
	}
}

func BenchmarkRead(b *testing.B) {
	for b.Loop() {
		Read()
	}
}
