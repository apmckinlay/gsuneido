// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package strs has miscellaneous functions for slices of strings
package intern

import (
	"strconv"
	"strings"
	"testing"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/strp"
)

func TestUnique(t *testing.T) {
	ck := func(s string) *byte {
		t.Helper()
		p := StrP(s)
		assert.This(strp.String(p)).Is(s)
		assert.This(String(s)).Is(s)
		return p
	}
	Clear()
	assert.This(Count()).Is(0)
	assert.This(Bytes()).Is(0)

	a := ck("hello")
	b := ck("world")
	c := ck("hello")
	assert.That(unsafe.Pointer(a) == unsafe.Pointer(c)) // deduplicated
	assert.That(unsafe.Pointer(a) != unsafe.Pointer(b))
	assert.This(Count()).Is(2)
	assert.This(Bytes()).Is(6 + 6) // length prefix + contents

	const n = 10000
	for range 2 {
		for i := range n {
			s := (strconv.Itoa(i) + "                      ")[:16]
			ck(s)
		}
		assert.This(Count()).Is(n + 2)
		assert.This(Bytes()).Is(12 + n*17)
	}

	Clear()
	assert.This(Count()).Is(0)
	assert.This(Bytes()).Is(0)
	// a pointer remains valid after Clear
	assert.This(strp.String(a)).Is("hello")
}

func TestStrP(t *testing.T) {
	Clear()
	// length prefix
	p := StrP("hello")
	assert.This(int(*p)).Is(5)
	assert.This(strp.String(p)).Is("hello")

	// maximum length is interned
	maxs := strings.Repeat("x", strp.MaxLen)
	p1 := StrP(maxs)
	p2 := StrP(maxs)
	assert.That(p1 != nil)
	assert.That(unsafe.Pointer(p1) == unsafe.Pointer(p2))

	// too long returns nil
	assert.That(StrP(strings.Repeat("x", strp.MaxLen+1)) == nil)

	// String clones (does not intern) strings that are too long
	long := strings.Repeat("x", strp.MaxLen+1)
	s1 := String(long)
	s2 := String(long)
	assert.This(s1).Is(long)
	assert.That(unsafe.StringData(s1) != unsafe.StringData(s2))
	Clear()
}
