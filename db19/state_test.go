// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestStateReadWrite(*testing.T) {
	offsets := [4]uint64{1, 2, 3, 4}
	store := stor.HeapStor(1024)
	off := writeState(store, offsets)
	offsets2, _ := readState(store, off)
	assert.This(offsets2).Is(offsets)
}
