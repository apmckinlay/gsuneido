// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// buildTree2 builds a btree with n sequential keys/offsets.
func buildTree2(n int) *btree {
	b := Builder(stor.HeapStor(8192))
	for i := 1; i <= n; i++ {
		k := strconv.Itoa(i)
		added := b.Add(k, uint64(i))
		if !added {
			panic("failed to add key " + k)
		}
	}
	bt := b.Finish().(*btree)
	GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, off uint64) string {
		return strconv.Itoa(int(off))
	}
	return bt
}

func TestRootUInitializedForBuilderFinish(t *testing.T) {
	bt := buildTree2(9)

	// Read the root node directly to see what it contains
	rootNode := readNode(bt.stor, bt.root)

	// Iterate through the root node to see its contents
	it := rootNode.iter()
	count := 0
	for it.next() {
		count++
	}

	// rootU must be non-empty for a non-empty tree
	assert.That(len(bt.rootUnode) > 0)
	// sanity: treeLevels should be >= 0
	assert.That(bt.treeLevels >= 0)
}

func TestRootUInitializedSingleEntry(t *testing.T) {
	bt := buildTree2(1)
	assert.That(len(bt.rootUnode) > 0)
}

func TestRootUEmptyForCreateBtree(t *testing.T) {
	// CreateBtree builds an empty tree (single empty root node)
	bt := CreateBtree(stor.HeapStor(8192), nil).(*btree)
	assert.That(len(bt.rootUnode) == 0)
}
