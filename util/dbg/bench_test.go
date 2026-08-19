// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbg

import "testing"

// assertOn mimics the dbg-on behavior: it calls the closure and panics on false.
// It is forced non-inlinable to match the real dbg.Assert which is too complex
// to inline because it calls log.Println and PrintStack.
//
//go:noinline
func assertOn(f func() bool) {
	if !f() {
		panic("ASSERT FAILED")
	}
}

// assertOff mimics the dbg-off behavior: an empty body that is trivially
// inlined and then dead-code eliminated along with the uncalled closure.
func assertOff(f func() bool) {
}

var sink bool

func BenchmarkAssertOnTrue(b *testing.B) {
	for b.Loop() {
		assertOn(func() bool { return true })
	}
}

func BenchmarkAssertOffTrue(b *testing.B) {
	for b.Loop() {
		assertOff(func() bool { return true })
	}
}

func BenchmarkAssertOnClosureCapture(b *testing.B) {
	x, y := 42, 100
	for b.Loop() {
		assertOn(func() bool { return x+y == 142 })
	}
}

func BenchmarkAssertOffClosureCapture(b *testing.B) {
	x, y := 42, 100
	for b.Loop() {
		assertOff(func() bool { return x+y == 142 })
	}
}

func BenchmarkAssertOnExpensiveBody(b *testing.B) {
	for b.Loop() {
		assertOn(func() bool {
			sum := 0
			for j := range 1000 {
				sum += j
			}
			return sum > 0
		})
	}
}

func BenchmarkAssertOffExpensiveBody(b *testing.B) {
	for b.Loop() {
		assertOff(func() bool {
			sum := 0
			for j := range 1000 {
				sum += j
			}
			return sum > 0
		})
	}
}

func BenchmarkAssertOnSideEffect(b *testing.B) {
	n := 0
	for b.Loop() {
		assertOn(func() bool { n++; return true })
	}
	sink = n > 0
}

func BenchmarkAssertOffSideEffect(b *testing.B) {
	n := 0
	for b.Loop() {
		assertOff(func() bool { n++; return true })
	}
	sink = n > 0
}
