// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSlot2Search(t *testing.T) {
	// Test empty slot
	emptySlot := slot2([]byte{0}) // 0 entries
	assert.T(t).This(emptySlot.Search("any")).Is(0)

	// Create a test slot with 3 entries: "apple", "banana", "cherry"
	// Entry format: 2 byte offset + 5 byte database offset
	const data = 1 + 3*7 + 2 // 1 byte count + 3 entries * 7 bytes + 2 byte node size
	slot := slot2([]byte{
		3, // 3 entries
		// "apple" - offset to start of data
		0, data, 0, 0, 0, 0, 0,
		// "banana" - offset after "apple"
		0, data + 5, 0, 0, 0, 0, 0,
		// "cherry" - offset after "banana"
		0, data + 11, 0, 0, 0, 0, 0,
		// Node size (end of data)
		0, data + 17,
		// Key data
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

func TestSlot2SearchSingleEntry(t *testing.T) {
	// Test slot with single entry "test"
	dataStart := 1 + 1*7 + 2 // 1 byte count + 1 entry * 7 bytes + 2 byte node size
	slot := slot2([]byte{
		1, // 1 entry
		// "test" - offset to start of data
		byte(dataStart >> 8), byte(dataStart), 0, 0, 0, 0, 0,
		// Node size (end of data)
		byte((dataStart + 4) >> 8), byte(dataStart + 4),
		// Key data
		't', 'e', 's', 't', // "test"
	})

	assert.T(t).This(slot.Search("apple")).Is(0) // before "test"
	assert.T(t).This(slot.Search("test")).Is(0)  // exact match
	assert.T(t).This(slot.Search("zebra")).Is(1) // after "test"
}

func makeSlot2(keys []string) slot2 {
	n := len(keys)
	if n == 0 {
		return slot2([]byte{0})
	}

	// Calculate total size needed
	headerSize := 1 + n*7 + 2 // 1 byte count + n entries * 7 bytes + 2 byte node size
	keyDataSize := 0
	for _, key := range keys {
		keyDataSize += len(key)
	}

	slot := make([]byte, headerSize+keyDataSize)
	slot[0] = byte(n) // number of entries

	// Write entry headers (2 byte offset + 5 byte database offset)
	pos := 1
	keyOffset := headerSize
	for _, key := range keys {
		// Write 2-byte offset in big-endian format
		slot[pos] = byte(keyOffset >> 8)
		slot[pos+1] = byte(keyOffset)
		// Fill in the 5 database offset bytes with zeros
		slot[pos+2] = 0
		slot[pos+3] = 0
		slot[pos+4] = 0
		slot[pos+5] = 0
		slot[pos+6] = 0
		pos += 7
		keyOffset += len(key)
	}

	// Write node size (end of data) in big-endian format
	endOffset := headerSize + keyDataSize
	slot[pos] = byte(endOffset >> 8)
	slot[pos+1] = byte(endOffset)

	// Write key data
	keyPos := headerSize
	for _, key := range keys {
		copy(slot[keyPos:], []byte(key))
		keyPos += len(key)
	}

	return slot2(slot)
}

func TestSlot2SearchWithHelper(t *testing.T) {
	// Test with helper function
	slot := makeSlot2([]string{"alpha", "beta", "gamma"})

	assert.T(t).This(slot.Search("aardvark")).Is(0)
	assert.T(t).This(slot.Search("alpha")).Is(0)
	assert.T(t).This(slot.Search("beta")).Is(1)
	assert.T(t).This(slot.Search("gamma")).Is(2)
	assert.T(t).This(slot.Search("delta")).Is(2)
	assert.T(t).This(slot.Search("zeta")).Is(3)
}
