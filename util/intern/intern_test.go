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
)

func TestUnique(t *testing.T) {
	assert.This(unsafe.Sizeof(entry{})).Is(uintptr(4))
	ck := func(s string) {
		t.Helper()
		assert.T(t).This(String(s)).Is(s)
	}
	assert.T(t).This(Count()).Is(0)
	assert.T(t).This(Bytes()).Is(0)
	ck("hello")
	ck("world")
	ck("hello")
	ck("world")
	ck(strings.Repeat("helloworld", 30)) // too big, separate Clone
	assert.This(Count()).Is(2)
	assert.This(Bytes()).Is(10)
	Clear()
	assert.T(t).This(Count()).Is(0)
	assert.T(t).This(Bytes()).Is(0)
	ck("hello")
	ck("world")
	ck("hello")
	ck("world")
	ck(strings.Repeat("helloworld", 30)) // too big, separate Clone
	assert.This(Count()).Is(2)
	assert.This(Bytes()).Is(10)
	Clear()
	assert.T(t).This(Count()).Is(0)
	assert.T(t).This(Bytes()).Is(0)
	const n = 10000
	for range 2 {
		for i := range n {
			s := (strconv.Itoa(i) + "                      ")[:16]
			ck(s)
		}
		assert.T(t).This(Count()).Is(n)
		assert.T(t).This(Bytes()).Is(n * 16)
	}
}
