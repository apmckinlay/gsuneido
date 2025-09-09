// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSlot1Search(t *testing.T) {
	// Test empty slot
	emptySlot := slot1([]byte{0}) // 0 entries
	assert.T(t).This(emptySlot.Search("any")).Is(0)

	// Create a test slot with 3 entries: "apple", "banana", "cherry"
	// Entry format: 1 byte size + 5 byte offset
	slot := slot1([]byte{
		3,                // 3 entries
		5, 0, 0, 0, 0, 0, // "apple" - size 5, offset 0
		6, 0, 0, 0, 0, 1, // "banana" - size 6, offset 1
		6, 0, 0, 0, 0, 2, // "cherry" - size 6, offset 2
		'a', 'p', 'p', 'l', 'e', // "apple"
		'b', 'a', 'n', 'a', 'n', 'a', // "banana"
		'c', 'h', 'e', 'r', 'r', 'y', // "cherry"
	})

	// Test exact matches
	assert.T(t).This(slot.Search("apple")).Is(0)
	assert.T(t).This(slot.Search("banana")).Is(1)
	assert.T(t).This(slot.Search("cherry")).Is(2)

	// Test insertion positions
	assert.T(t).This(slot.Search("aaa")).Is(0)       // before "apple"
	assert.T(t).This(slot.Search("avocado")).Is(1)   // between "apple" and "banana"
	assert.T(t).This(slot.Search("blueberry")).Is(2) // between "banana" and "cherry"
	assert.T(t).This(slot.Search("date")).Is(3)      // after "cherry"
}

func TestSlot1SearchSingleEntry(t *testing.T) {
	// Test slot with single entry "test"
	slot := slot1([]byte{
		1,                // 1 entry
		4, 0, 0, 0, 0, 0, // "test" - size 4, offset 0
		't', 'e', 's', 't', // "test"
	})

	assert.T(t).This(slot.Search("apple")).Is(0) // before "test"
	assert.T(t).This(slot.Search("test")).Is(0)  // exact match
	assert.T(t).This(slot.Search("zebra")).Is(1) // after "test"
}

func makeSlot1(keys []string) slot1 {
	n := len(keys)
	if n == 0 {
		return slot1([]byte{0})
	}

	// Calculate total size needed
	headerSize := 1 + n*6 // 1 byte count + n entries * 6 bytes each
	keyDataSize := 0
	for _, key := range keys {
		keyDataSize += len(key)
	}

	slot := make([]byte, headerSize+keyDataSize)
	slot[0] = byte(n) // number of entries

	// Write entry headers (size + 5 byte offset)
	pos := 1
	for _, key := range keys {
		slot[pos] = byte(len(key)) // key size
		// Fill in the 5 offset bytes with zeros
		slot[pos+1] = 0
		slot[pos+2] = 0
		slot[pos+3] = 0
		slot[pos+4] = 0
		slot[pos+5] = 0
		pos += 6
	}

	// Write key data
	keyPos := headerSize
	for _, key := range keys {
		copy(slot[keyPos:], []byte(key))
		keyPos += len(key)
	}

	return slot1(slot)
}

func TestSlot1SearchWithHelper(t *testing.T) {
	// Test with helper function
	slot := makeSlot1([]string{"alpha", "beta", "gamma"})

	assert.T(t).This(slot.Search("aardvark")).Is(0)
	assert.T(t).This(slot.Search("alpha")).Is(0)
	assert.T(t).This(slot.Search("beta")).Is(1)
	assert.T(t).This(slot.Search("gamma")).Is(2)
	assert.T(t).This(slot.Search("delta")).Is(2)
	assert.T(t).This(slot.Search("zeta")).Is(3)
}
