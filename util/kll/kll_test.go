// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package kll

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/generic/ints"
)

func TestInsertLevel0(t *testing.T) {
	sk := New[int]()
	sk.Insert(1)
	sk.Insert(2)
	assert.This(sk.Count()).Is(2)
	assert.This(sk.levels[0]).Is([]int{1, 2})
	// sk.print()
}

func TestInsertCompact(t *testing.T) {
	sk := New[int]()
	for i := range 202 {
		sk.Insert(i)
	}
	// sk.print()
	assert.This(len(sk.levels)).Is(2)
	assert.This(sk.levels[1]).Is([]int{201})
}

func TestQuery(t *testing.T) {
	sk := New[int]()
	for i := range 5000 {
		sk.Insert(int(bits.Shuffle16(uint16(i))))
	}
	// sk.print()
	m := sk.Query(0.5)
	assert.Msg(m - 32767).That(ints.Abs(m-32767) < 800)

	sk = New[int]()
	for range 5000 {
		sk.Insert(rand.IntN(1000))
	}
	// sk.print()
	m = sk.Query(0.5)
	assert.Msg(m - 500).That(ints.Abs(m-500) < 25) // 2.5%
}

func (sk *Sketch[T]) print() {
	fmt.Println("Sketch count:", sk.count, "levels:", len(sk.levels))
	for h, level := range sk.levels {
		fmt.Println("  level", h, "capacity:", caps[h], 
			"cap:", cap(level), "len:", len(level))
	}
}
