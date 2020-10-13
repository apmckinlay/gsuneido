// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestStateReadWrite(*testing.T) {
	store := stor.HeapStor(1024)
	off := writeState(store, 1234, 5678)
	offSchema, offInfo, _ := readState(store, off)
	assert.This(offSchema).Is(1234)
	assert.This(offInfo).Is(5678)
}
