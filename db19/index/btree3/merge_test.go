// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/testdata"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

var noSplit = func(nd splitable) bool {
	return false
}

func heapstor(chunksize int) *stor.Stor {
	st := stor.HeapStor(chunksize)
	st.Alloc(1) // avoid offset 0
	return st
}

// update -----------------------------------------------------------

func TestMergeUpdate(t *testing.T) {
	bt := createTestBtree(10)
	bt.shouldSplit = noSplit
	// bt.print()

	bt = testUpdate(bt, "1000", 9999)
	// bt.print()
	assert.This(bt.String()).Is("1000 9999 1001 1001 1002 1002 1003 1003 " +
		"1004 1004 1005 1005 1006 1006 1007 1007 1008 1008 1009 1009")

	bt = testUpdate(bt, "1009", 8888)
	// bt.print()
	assert.This(bt.String()).Is("1000 9999 1001 1001 1002 1002 1003 1003 " +
		"1004 1004 1005 1005 1006 1006 1007 1007 1008 1008 1009 8888")

	ib := &ixbuf.T{}
	ib.Insert("1001", ixbuf.Update|1) // first leaf
	ib.Insert("1002", ixbuf.Update|2) // same leaf
	ib.Insert("1004", ixbuf.Update|4) // second leaf
	ib.Insert("1007", ixbuf.Update|7) // same leaf
	ib.Insert("1009", ixbuf.Update|9)
	bt = bt.MergeAndSave(ib.Iter())
	// bt.print()
	assert.This(bt.String()).Is("1000 9999 1001 1 1002 2 1003 1003 " +
		"1004 4 1005 1005 1006 1006 1007 7 1008 1008 1009 9")

}

func testUpdate(bt *btree, key string, off uint64) *btree {
	ib := &ixbuf.T{}
	ib.Insert(key, ixbuf.Update|off)
	return bt.MergeAndSave(ib.Iter())
}

func TestMergeRootLeaf(t *testing.T) {
	// Create a simple btree with initial entries (single root leaf node)
	bldr := Builder(heapstor(8192))
	bldr.shouldSplit = func(nd splitable) bool {
		return false // Never split to ensure single root leaf node
	}
	assert.That(bldr.Add("apple", 1))
	assert.That(bldr.Add("cherry", 3))
	bt := bldr.Finish()
	assert.This(bt.String()).Is("apple 1 cherry 3")

	var ib *ixbuf.T
	test := func(expected string) {
		bt = bt.MergeAndSave(ib.Iter())
		bt.Check(nil)
		assert.This(bt.String()).Is(expected)
	}

	// add
	ib = &ixbuf.T{}
	ib.Insert("banana", 2)
	test("apple 1 banana 2 cherry 3")

	// update
	ib = &ixbuf.T{}
	ib.Insert("banana", ixbuf.Update|222)
	test("apple 1 banana 222 cherry 3")

	// delete
	ib = &ixbuf.T{}
	ib.Insert("banana", ixbuf.Delete|222)
	test("apple 1 cherry 3")
}

func (bt *btree) String() string {
	iter := bt.Iterator()
	var sb strings.Builder
	sep := ""
	for iter.Next(); iter.HasCur(); iter.Next() {
		fmt.Fprintf(&sb, "%s%s %d", sep, iter.Key(), iter.Offset())
		sep = " "
	}
	return sb.String()
}

func TestMergeOneTreeLevel(t *testing.T) {
	// Create a btree with one tree level (tree root with two leaf nodes)
	bldr := Builder(heapstor(8192))
	bldr.shouldSplit = func(nd splitable) bool {
		return nd.nkeys() >= 2
	}

	// Add entries to create tree structure: root -> [leaf1, leaf2]
	assert.That(bldr.Add("apple", 1))
	assert.That(bldr.Add("banana", 2)) // This should trigger split
	assert.That(bldr.Add("cherry", 3)) // This goes to new leaf
	bt := bldr.Finish()
	bt.shouldSplit = noSplit
	assert.This(bt.treeLevels).Is(1)
	assert.This(bt.String()).Is("apple 1 banana 2 cherry 3")

	var ib *ixbuf.T
	test := func(expected string) {
		bt = bt.MergeAndSave(ib.Iter())
		bt.Check(nil)
		// bt.print()
		// fmt.Println("-------------------------------")
		assert.This(bt.String()).Is(expected)
		assert.This(bt.treeLevels).Is(1) // Should maintain tree structure
	}

	// add entry at end
	ib = &ixbuf.T{}
	ib.Insert("date", 4)
	test("apple 1 banana 2 cherry 3 date 4")

	// update entry
	ib = &ixbuf.T{}
	ib.Insert("banana", ixbuf.Update|222)
	test("apple 1 banana 222 cherry 3 date 4")

	// delete entry
	ib = &ixbuf.T{}
	ib.Insert("banana", ixbuf.Delete|222)
	test("apple 1 cherry 3 date 4")

	// add entry at beginning
	ib = &ixbuf.T{}
	ib.Insert("aaa", 11)
	test("aaa 11 apple 1 cherry 3 date 4")

	// multiple actions in one ixbuf
	ib = &ixbuf.T{}
	ib.Insert("elephant", 5)             // add new entry
	ib.Insert("apple", ixbuf.Update|111) // update existing entry
	ib.Insert("date", ixbuf.Delete|4)    // delete existing entry
	ib.Insert("bear", 22)                // add another new entry
	test("aaa 11 apple 111 bear 22 cherry 3 elephant 5")
}

// delete -----------------------------------------------------------

func TestMergeDelete(t *testing.T) {
	test := func(bt *btree) {
		bt.Check(nil)
		iter := bt.Iterator()
		iter.Next()
		assert.That(iter.Eof())
		assert.Msg("treeLevels").This(bt.treeLevels).Is(0)
	}
	for n := 1; n < 100; n++ {
		orig := testBtree(n, 4)

		// all in one batch
		ib := &ixbuf.T{}
		for i := 0; i < n; i++ {
			key := strconv.Itoa(base + i)
			ib.Delete(key, uint64(base+i))
		}
		bt := orig.MergeAndSave(ib.Iter())
		test(bt)

		// one at a time, forward
		bt = orig
		for i := 0; i < n; i++ {
			ib := &ixbuf.T{}
			key := strconv.Itoa(base + i)
			ib.Delete(key, uint64(base+i))
			bt = bt.MergeAndSave(ib.Iter())
		}
		test(bt)

		// one at a time, reverse
		bt = orig
		for i := n - 1; i >= 0; i-- {
			ib := &ixbuf.T{}
			key := strconv.Itoa(base + i)
			ib.Delete(key, uint64(base+i))
			bt = bt.MergeAndSave(ib.Iter())
		}
		test(bt)

		for batch := 2; batch <= 7; batch++ {
			perm := rand.Perm(n)

			// one at a time, random
			bt = orig
			for i := 0; i < n; i++ {
				ib := &ixbuf.T{}
				key := strconv.Itoa(base + perm[i])
				ib.Delete(key, uint64(base+perm[i]))
				bt = bt.MergeAndSave(ib.Iter())
			}
			test(bt)

			// batches, random
			bt = orig
			for i := 0; i < n; i += batch {
				ib := &ixbuf.T{}
				for j := 0; j < batch && i+j < n; j++ {
					key := strconv.Itoa(base + perm[i+j])
					ib.Delete(key, uint64(base+perm[i+j]))
				}
				bt = bt.MergeAndSave(ib.Iter())
			}
			test(bt)
		}
	}
}

func TestMergeDeleteLast(t *testing.T) {
	bldr := Builder(heapstor(8192))
	bldr.shouldSplit = func(nd splitable) bool {
		return nd.nkeys() >= 3 // Small split threshold to trigger split easily
	}
	for i := 1000; i <= 1009; i++ {
		assert.That(bldr.Add(strconv.Itoa(i), uint64(i)))
	}
	bt := bldr.Finish()
	// bt.print()
	ib := &ixbuf.T{}
	ib.Insert("1009", ixbuf.Delete|1009)
	bt = bt.MergeAndSave(ib.Iter())
	// bt.print()
	bt.Check(nil)
}

/*
deleting the last leaf when it's on a branch

	r -	a2 - a
*/
func TestMergeDelete1(t *testing.T) {
	hs := heapstor(8192)
	a := makeLeaf("a", 1)
	ao := a.write(hs)
	a2 := makeTree(ao)
	a2o := a2.write(hs)

	r := makeTree(a2o)
	ro := r.write(hs)

	bt := &btree{stor: hs, root: ro, treeLevels: 2}
	// bt.print()

	ib := &ixbuf.T{}
	ib.Insert("a", ixbuf.Delete|1)
	bt = bt.MergeAndSave(ib.Iter())
	// bt.print()

	assert.This(bt.treeLevels).Is(0)
	assert.This(bt.String()).Is("")
}

/*
deleting a leaf from one branch
means pulling up a leaf from another branch

		a2 - a
	  /
	r
	  \
		b2 - b
*/
func TestMergeDelete2(t *testing.T) {
	hs := heapstor(8192)
	a := makeLeaf("a", 1)
	ao := a.write(hs)
	a2 := makeTree(ao)
	a2o := a2.write(hs)

	b := makeLeaf("b", 2)
	bo := b.write(hs)
	b2 := makeTree(bo)
	b2o := b2.write(hs)

	r := makeTree(a2o, "b", b2o)
	ro := r.write(hs)

	bt := &btree{stor: hs, root: ro, treeLevels: 2}
	// bt.print()

	ib := &ixbuf.T{}
	ib.Insert("a", ixbuf.Delete|1)
	bt = bt.MergeAndSave(ib.Iter())
	// bt.print()

	assert.This(bt.treeLevels).Is(0)
	assert.This(bt.String()).Is("b 2")
}

// insert -----------------------------------------------------------

func TestMergeSplitLeaf(t *testing.T) {
	// Create a leaf node with multiple entries to test splitting
	var b leafBuilder
	b.add("apple", 1)
	b.add("banana", 2)
	b.add("cherry", 3)
	b.add("date", 4)
	b.add("elderberry", 5)
	originalLeaf := b.finish()

	// Verify we have the expected number of keys
	assert.This(originalLeaf.nkeys()).Is(5)
	assert.This(originalLeaf.String()).Is("leaf{apple 1 banana 2 cherry 3 date 4 elderberry 5}")

	// Create a state with the leaf node for testing split
	lm := leafMerge{
		leaf:     originalLeaf,
		modified: true, // Mark as modified so split will work
	}

	// Call splitLeaf
	left, right, splitKey := lm.split()

	// Verify the split key is correct (should be the first key of the right node)
	assert.This(splitKey).Is("cherry")

	// Verify left node contains first half of entries
	assert.This(left.nkeys()).Is(2)
	assert.This(left.String()).Is("leaf{apple 1 banana 2}")

	// Verify right node contains second half of entries
	assert.This(right.nkeys()).Is(3)
	assert.This(right.String()).Is("leaf{cherry 3 date 4 elderberry 5}")
}

func TestMergeSplitTree(t *testing.T) {
	// Create a tree node with multiple entries to test splitting
	var b treeBuilder
	b.add(100, "apple")
	b.add(200, "banana")
	b.add(300, "cherry")
	originalTree := b.finish(400)
	assert.This(originalTree.String()).Is("tree{100 <apple> 200 <banana> 300 <cherry> 400}")

	tm := &treeMerge{
		tree:     originalTree,
		modified: true,
	}
	left, right, splitKey := tm.split()

	assert.This(splitKey).Is("banana")
	assert.This(left.String()).Is("tree{100 <apple> 200}")
	assert.This(right.String()).Is("tree{300 <cherry> 400}")
}

func TestMergeInsert(t *testing.T) {
	rng := rand.New(rand.NewPCG(123, 456))
	empty := func() *btree {
		bt := Builder(heapstor(8192)).Finish() // Create empty btree
		bt.shouldSplit = func(nd splitable) bool {
			return nd.nkeys() >= 4 // Small split threshold for testing
		}
		return bt
	}
	var bt *btree
	insert := func(i int) {
		ib := &ixbuf.T{}
		ib.Insert(strconv.Itoa(i), uint64(i))
		bt = bt.MergeAndSave(ib.Iter())
	}
	for n := 1; n < 100; n++ {
		// add at end 1 at a time
		bt = empty()
		for i := 0; i <= n; i++ {
			insert(i + base)
			// bt.print()
			bt.Check(nil)
		}

		// add at end single batch
		bt = empty()
		ib := &ixbuf.T{}
		for i := 0; i <= n; i++ {
			ib.Insert(strconv.Itoa(i+base), uint64(i+base))
		}
		bt = bt.MergeAndSave(ib.Iter())
		// bt.print()
		bt.Check(nil)

		// add at beginning
		bt = empty()
		for i := n; i >= 1; i-- {
			insert(i + base)
			// bt.print()
			bt.Check(nil)
		}

		// add randomly, 1 at a time
		bt = empty()
		for _, i := range rand.Perm(n) {
			insert(i + base)
			// bt.print()
			bt.Check(nil)
		}

		// add randomly, in batches
		bt = empty()
		p := rng.Perm(n)
		batch := 5
		for j := 0; j < n; j += batch {
			ib := &ixbuf.T{}
			for k := 0; k < batch && j+k < n; k++ {
				i := p[j+k] + base
				ib.Insert(strconv.Itoa(i), uint64(i))
			}
			bt = bt.MergeAndSave(ib.Iter())
			// bt.print()
			bt.Check(nil)
		}
	}
}

//-------------------------------------------------------------------

func TestMergeMix(*testing.T) {
	rng := rand.New(rand.NewPCG(123, 456))
	nMerges := 2000
	opsPerMerge := 5
	if testing.Short() {
		nMerges = 200
		opsPerMerge = 200
	}
	d := testdata.New()

	bt := Builder(heapstor(8192)).Finish() // Create empty btree
	bt.shouldSplit = func(nd splitable) bool {
		return nd.nkeys() >= 4 // Small split threshold for testing
	}

	for range nMerges {
		trace("---")
		x := &ixbuf.T{}
		for range opsPerMerge {
			k := rng.IntN(4)
			switch {
			case k == 0 || k == 1 || d.Len() == 0:
				x.Insert(d.Gen())
			case k == 2:
				_, key, _ := d.Rand()
				off := d.NextOff()
				x.Update(key, off)
				d.Update(key, off)
			case k == 3:
				i, key, off := d.Rand()
				x.Delete(key, off)
				d.Delete(i)
			}
		}
		bt = bt.MergeAndSave(x.Iter())
	}
	bt.Check(nil)
	d.Check(bt)
	// d.CheckIter(bt.Iterator()) // TODO: Iterator doesn't implement iterator.T interface yet
}
