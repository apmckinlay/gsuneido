// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/util/assert"
)

const keyBase = 1000

// createTestBtree creates a btree with the specified number of keys
func createTestBtree(treeSize int) *btree {
	bldr := Builder(heapstor(8192))
	bldr.shouldSplit = func(nd node) bool {
		return nd.noffs() >= 4 // Small split threshold for testing
	}
	for i := range treeSize {
		key := strconv.Itoa(keyBase + i)
		assert.That(bldr.Add(key, uint64(keyBase+i)))
	}
	bt := bldr.Finish().(*btree)
	bt.Check(nil)
	return bt
}

// go test ./db19/index/btree3 -run="^$" -fuzz=FuzzRandomUpdateBatches
func FuzzRandomUpdateBatches(f *testing.F) {
	f.Add(uint64(42), uint8(10), uint8(5))
	f.Add(uint64(123), uint8(20), uint8(15))
	f.Add(uint64(456), uint8(30), uint8(25))

	f.Fuzz(func(t *testing.T, seed uint64, b1, b2 uint8) {
		treeSize := int(b1) + 1
		maxBatchSize := int(b2) + 1
		// fmt.Println("seed:", seed, "treeSize:", treeSize, "maxBatchSize:", maxBatchSize)

		rng := rand.New(rand.NewPCG(seed, seed))

		bt := createTestBtree(treeSize)
		bt.shouldSplit = noSplit

		// Track expected offsets for each key
		expectedOffsets := make(map[string]uint64)
		baseKey := 1000
		for i := range treeSize {
			key := strconv.Itoa(baseKey + i)
			expectedOffsets[key] = uint64(baseKey + i)
		}

		// Do 100 MergeAndSave operations with random batches of updates
		for range 100 {
			ib := &ixbuf.T{}

			// Insert updates until we reach the desired batch size
			// ixbuf handles duplicates automatically
			batchSize := max(1, rng.IntN(min((treeSize+1)/2, maxBatchSize)))
			for ib.Len() < batchSize {
				// Pick a random existing key to update
				keyIndex := rng.IntN(treeSize)
				key := strconv.Itoa(baseKey + keyIndex)
				newOffset := uint64(rng.IntN(9000) + 1000) // random value 1000-9999
				expectedOffsets[key] = newOffset

				// fmt.Println("=" + key, "->", newOffset)
				ib.Update(key, newOffset)
			}

			bt = bt.MergeAndSave(ib.Iter()).(*btree)
			// bt.print()
			bt.Check(nil)
		}

		// Verify tree still contains all original keys with correct offsets
		iter := bt.Iterator()
		count := 0
		for iter.Next(); iter.HasCur(); iter.Next() {
			key, offset := iter.Cur()
			count++

			// Verify the offset matches our expected value
			expectedOffset, exists := expectedOffsets[key]
			assert.That(exists)
			assert.This(offset).Is(expectedOffset)
		}
		assert.This(count).Is(treeSize)

		bt.Check(nil)
	})
}

// go test ./db19/index/btree3 -run="^$" -fuzz=FuzzRandomInsertBatches
func FuzzRandomInsertBatches(f *testing.F) {
	f.Add(uint64(42), uint8(20), uint8(5))
	f.Add(uint64(123), uint8(50), uint8(10))
	f.Add(uint64(789), uint8(100), uint8(15))

	f.Fuzz(func(t *testing.T, seed uint64, totalKeys, maxBatchSize uint8) {
		// Ensure reasonable bounds
		totalSize := int(totalKeys) + 1   // 1-256 keys
		maxBatch := int(maxBatchSize) + 1 // 1-256 batch size
		if maxBatch > totalSize {
			maxBatch = totalSize
		}

		// fmt.Println("seed:", seed, "totalKeys:", totalSize, "maxBatchSize:", maxBatch)

		rng := rand.New(rand.NewPCG(seed, seed))

		// Start with empty btree
		bt := Builder(heapstor(8192)).Finish()
		bt.(*btree).shouldSplit = func(nd node) bool {
			return nd.noffs() >= 4 // Small split threshold for testing
		}

		// Generate shuffled unique keys
		perm := rng.Perm(totalSize)

		// Insert keys in random-sized batches
		i := 0
		for i < len(perm) {
			batchSize := rng.IntN(maxBatch) + 1
			if i+batchSize > len(perm) {
				batchSize = len(perm) - i
			}

			ib := &ixbuf.T{}
			for j := 0; j < batchSize; j++ {
				keyNum := keyBase + perm[i+j]
				key := strconv.Itoa(keyNum)
				offset := uint64(keyNum)
				ib.Insert(key, offset)
			}

			bt = bt.MergeAndSave(ib.Iter()).(*btree)
			bt.Check(nil)
			i += batchSize
		}

		// Verify final btree contains all keys with correct offsets in sequence
		iter := bt.Iterator()
		count := 0
		for iter.Next(); iter.HasCur(); iter.Next() {
			key, offset := iter.Cur()

			// Verify we get the expected sequential key
			keyNum, err := strconv.Atoi(key)
			if err != nil {
				t.Fatalf("Failed to parse key %s: %v", key, err)
			}
			expectedKeyNum := keyBase + count
			assert.This(keyNum).Is(expectedKeyNum)

			// Verify the offset equals the key value
			assert.This(offset).Is(uint64(keyNum))

			count++
		}
		assert.This(count).Is(totalSize)

		bt.Check(nil)
	})
}

// go test ./db19/index/btree3 -run="^$" -fuzz=FuzzRandomDeleteBatches
func FuzzRandomDeleteBatches(f *testing.F) {
	f.Add(uint64(42), uint8(10), uint8(3))
	f.Add(uint64(123), uint8(50), uint8(10))
	f.Add(uint64(456), uint8(100), uint8(25))

	f.Fuzz(func(t *testing.T, seed uint64, b1, b2 uint8) {
		treeSize := int(b1) + 1
		maxBatchSize := int(b2) + 1
		// fmt.Println("seed:", seed, "treeSize:", treeSize, "maxBatchSize:", maxBatchSize)

		rng := rand.New(rand.NewPCG(seed, seed))

		bt := createTestBtree(treeSize)

		// Delete all keys in random sized batches
		remaining := rng.Perm(treeSize)
		for len(remaining) > 0 {
			batchSize := rng.IntN(min(maxBatchSize, len(remaining))) + 1

			ib := &ixbuf.T{}
			for i := 0; i < batchSize; i++ {
				key := strconv.Itoa(keyBase + remaining[i])
				keyNum, err := strconv.Atoi(key)
				if err != nil {
					t.Fatalf("Failed to parse key %s: %v", key, err)
				}
				// fmt.Println("-" + key)
				ib.Insert(key, ixbuf.Delete|uint64(keyNum))
			}
			bt = bt.MergeAndSave(ib.Iter()).(*btree)
			// bt.print()
			bt.Check(nil)

			// Remove processed keys from remaining
			remaining = remaining[batchSize:]
		}

		// Verify tree is completely empty
		iter := bt.Iterator()
		iter.Next()
		assert.That(iter.Eof())
		assert.This(bt.treeLevels).Is(0)
	})
}
