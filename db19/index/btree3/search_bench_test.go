// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// generateKeys creates 100 keys from 1000-1099
func generateKeys() []string {
	keys := make([]string, 100)
	for i := range 100 {
		keys[i] = strconv.Itoa(1000 + i)
	}
	return keys
}

func BenchmarkSlot1Search(b *testing.B) {
	keys := generateKeys()
	slot1 := makeSlot1(keys)
	searchKey := keys[50] // Search for middle key
	assert.This(slot1.Search(searchKey)).Is(50)

	for b.Loop() {
		slot1.Search(searchKey)
	}
}

func BenchmarkSlot2Search(b *testing.B) {
	keys := generateKeys()
	slot2 := makeSlot2(keys)
	searchKey := keys[50] // Search for middle key
	assert.This(slot2.Search(searchKey)).Is(50)

	for b.Loop() {
		slot2.Search(searchKey)
	}
}
