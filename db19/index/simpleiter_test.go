// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"strconv"
	"testing"

	btree "github.com/apmckinlay/gsuneido/db19/index/btree3"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type simpleTestTran struct {
	overlay *Overlay
	num     int
}

func (t *simpleTestTran) GetIndexI(string, int) *Overlay {
	return t.overlay
}

func (t *simpleTestTran) Read(string, int, string, string) {
	// SimpleIter doesn't track reads
}

func (t *simpleTestTran) Num() int {
	return t.num
}

func TestNewSimpleIter(t *testing.T) {
	// Test with valid overlay (single empty layer, no mut)
	store := stor.HeapStor(8192)
	bt := btree.CreateBtree(store, nil)
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)
	assert.That(!si.HasCur())

	// Test with overlay that has mut (should return nil)
	mut := &ixbuf.T{}
	mut.Insert("test", 1)
	ovMut := &Overlay{bt: bt, layers: []*ixbuf.T{{}}, mut: mut}
	siNil := NewSimpleIter(tran, ovMut)
	assert.That(siNil == nil)

	// Test with overlay that has non-empty layer (should return nil)
	ovLayer := &Overlay{bt: bt, layers: []*ixbuf.T{{}, {}}}
	ovLayer.layers[1].Insert("test", 1)
	siNil2 := NewSimpleIter(tran, ovLayer)
	assert.That(siNil2 == nil)

	// Test with multiple layers (should return nil)
	ovMulti := &Overlay{bt: bt, layers: []*ixbuf.T{{}, {}}}
	siNil3 := NewSimpleIter(tran, ovMulti)
	assert.That(siNil3 == nil)
}

func TestSimpleIterBasic(t *testing.T) {
	assert := assert.T(t)

	// Create overlay with some data
	store := stor.HeapStor(8192)
	bldr := btree.Builder(store)
	for i := 1; i <= 5; i++ {
		key := strconv.Itoa(i)
		assert.That(bldr.Add(key, uint64(i)))
	}
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)

	// Test forward iteration
	si.Next(tran)
	assert.That(!si.Eof())
	assert.That(si.HasCur())
	key, off := si.Cur()
	assert.This(key).Is("1")
	assert.This(off).Is(uint64(1))
	assert.This(si.CurOff()).Is(uint64(1))

	si.Next(tran)
	key, off = si.Cur()
	assert.This(key).Is("2")
	assert.This(off).Is(uint64(2))

	// Test backward iteration
	si.Rewind()
	si.Prev(tran)
	assert.That(!si.Eof())
	key, off = si.Cur()
	assert.This(key).Is("5")
	assert.This(off).Is(uint64(5))

	si.Prev(tran)
	key, off = si.Cur()
	assert.This(key).Is("4")
	assert.This(off).Is(uint64(4))
}

func TestSimpleIterEmpty(t *testing.T) {
	assert := assert.T(t)

	// Create empty overlay
	store := stor.HeapStor(8192)
	bt := btree.CreateBtree(store, nil)
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)

	// Test Next on empty
	si.Next(tran)
	assert.That(si.Eof())
	assert.That(!si.HasCur())

	// Test Prev on empty
	si.Rewind()
	si.Prev(tran)
	assert.That(si.Eof())
	assert.That(!si.HasCur())
}

func TestSimpleIterSingleItem(t *testing.T) {
	assert := assert.T(t)

	// Create overlay with single item
	store := stor.HeapStor(8192)
	bldr := btree.Builder(store)
	assert.That(bldr.Add("42", 42))
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)

	// Test forward
	si.Next(tran)
	assert.That(!si.Eof())
	key, off := si.Cur()
	assert.This(key).Is("42")
	assert.This(off).Is(uint64(42))

	si.Next(tran)
	assert.That(si.Eof())

	// Test backward
	si.Rewind()
	si.Prev(tran)
	assert.That(!si.Eof())
	key, off = si.Cur()
	assert.This(key).Is("42")
	assert.This(off).Is(uint64(42))

	si.Prev(tran)
	assert.That(si.Eof())
}

func TestSimpleIterRange(t *testing.T) {
	assert := assert.T(t)

	// Create overlay with data 1..9
	store := stor.HeapStor(8192)
	bldr := btree.Builder(store)
	for i := 1; i <= 9; i++ {
		key := strconv.Itoa(i)
		assert.That(bldr.Add(key, uint64(i)))
	}
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)

	// Test range 3..6
	si.Range(Range{Org: "3", End: "7"})

	// Forward within range
	si.Next(tran)
	key, off := si.Cur()
	assert.This(key).Is("3")
	assert.This(off).Is(uint64(3))

	si.Next(tran)
	key, off = si.Cur()
	assert.This(key).Is("4")
	assert.This(off).Is(uint64(4))

	si.Next(tran)
	key, off = si.Cur()
	assert.This(key).Is("5")
	assert.This(off).Is(uint64(5))

	si.Next(tran)
	key, off = si.Cur()
	assert.This(key).Is("6")
	assert.This(off).Is(uint64(6))

	si.Next(tran)
	assert.That(si.Eof())

	// Backward within range
	si.Rewind()
	si.Prev(tran)
	key, off = si.Cur()
	assert.This(key).Is("6")
	assert.This(off).Is(uint64(6))

	si.Prev(tran)
	key, off = si.Cur()
	assert.This(key).Is("5")
	assert.This(off).Is(uint64(5))
}

func TestSimpleIterStateTransitions(t *testing.T) {
	// Create overlay with data 1..3
	store := stor.HeapStor(8192)
	bldr := btree.Builder(store)
	for i := 1; i <= 3; i++ {
		key := strconv.Itoa(i)
		assert.That(bldr.Add(key, uint64(i)))
	}
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)

	// Test state transitions
	assert.That(!si.HasCur())

	// Move to first item
	si.Next(tran)
	assert.That(!si.Eof())
	assert.That(si.HasCur())

	// Move to EOF
	si.Next(tran) // to 2
	si.Next(tran) // to 3
	si.Next(tran) // to EOF
	assert.That(si.Eof())
	assert.That(!si.HasCur())

	// Next from EOF stays at EOF
	si.Next(tran)
	assert.That(si.Eof())

	// Rewind resets state
	si.Rewind()
	assert.That(!si.Eof()) // rewound, not eof
	assert.That(!si.HasCur())

	// Prev from rewind goes to last
	si.Prev(tran)
	assert.That(!si.Eof())
	assert.That(si.HasCur())
	key, off := si.Cur()
	assert.This(key).Is("3")
	assert.This(off).Is(uint64(3))
}

func TestSimpleIterPanicConditions(t *testing.T) {
	assert := assert.T(t)

	// Create overlay with data
	store := stor.HeapStor(8192)
	bldr := btree.Builder(store)
	assert.That(bldr.Add("test", 123))
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)

	// Test Cur() panic when rewound
	assert.This(func() { si.Cur() }).Panics("SimpleIter rewound")
	assert.This(func() { si.CurOff() }).Panics("SimpleIter rewound")

	// Test Cur() panic when EOF
	si.Next(tran)
	si.Next(tran) // to EOF
	assert.This(func() { si.Cur() }).Panics("SimpleIter eof")
	assert.This(func() { si.CurOff() }).Panics("SimpleIter eof")

	// Rewind and test again
	si.Rewind()
	si.Next(tran) // to first item
	si.Next(tran) // to EOF
	assert.This(func() { si.Cur() }).Panics("SimpleIter eof")
}

func TestSimpleIterTransactionMismatch(t *testing.T) {
	// Create overlay with data
	store := stor.HeapStor(8192)
	bldr := btree.Builder(store)
	assert.That(bldr.Add("test", 123))
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran1 := &simpleTestTran{overlay: ov, num: 0}
	tran2 := &simpleTestTran{overlay: ov, num: 2}

	si := NewSimpleIter(tran1, ov)
	assert.That(si != nil)

	// Test with wrong transaction - should panic due to assert
	assert.This(func() { si.Next(tran2) }).Panics("SimpleIter tran changed")
	assert.This(func() { si.Prev(tran2) }).Panics("SimpleIter tran changed")
}

func TestSimpleIterDirectionChanges(t *testing.T) {
	assert := assert.T(t)

	// Create overlay with data 1..5
	store := stor.HeapStor(8192)
	bldr := btree.Builder(store)
	for i := 1; i <= 5; i++ {
		key := strconv.Itoa(i)
		assert.That(bldr.Add(key, uint64(i)))
	}
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)

	// Test direction changes
	si.Next(tran) // to 1
	key, off := si.Cur()
	assert.This(key).Is("1")
	assert.This(off).Is(uint64(1))

	si.Prev(tran) // should stay at EOF (since we were at front)
	assert.That(si.Eof())

	si.Rewind()
	si.Prev(tran) // to 5
	key, off = si.Cur()
	assert.This(key).Is("5")
	assert.This(off).Is(uint64(5))

	si.Prev(tran) // to 4
	key, off = si.Cur()
	assert.This(key).Is("4")
	assert.This(off).Is(uint64(4))

	si.Next(tran) // back to 5
	key, off = si.Cur()
	assert.This(key).Is("5")
	assert.This(off).Is(uint64(5))
}

func TestSimpleIterLargeDataset(t *testing.T) {
	assert := assert.T(t)

	// Create overlay with larger dataset
	store := stor.HeapStor(8192)
	bldr := btree.Builder(store)
	const n = 1000
	for i := 1; i <= n; i++ {
		// Zero-pad to ensure lexicographic order
		key := strconv.Itoa(i + 1000000)
		assert.That(bldr.Add(key, uint64(i)))
	}
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran := &simpleTestTran{overlay: ov}

	si := NewSimpleIter(tran, ov)
	assert.That(si != nil)

	// Test forward iteration through all items
	count := 0
	for si.Next(tran); !si.Eof(); si.Next(tran) {
		count++
	}
	assert.This(count).Is(n)

	// Test backward iteration through all items
	si.Rewind()
	count = 0
	for si.Prev(tran); !si.Eof(); si.Prev(tran) {
		count++
	}
	assert.This(count).Is(n)
}

func TestSimpleIterOptimizationConditions(t *testing.T) {
	assert := assert.T(t)

	store := stor.HeapStor(8192)
	bt := btree.CreateBtree(store, nil)

	// Test case 1: No layers, no mut - should work
	ov1 := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran1 := &simpleTestTran{overlay: ov1}
	si1 := NewSimpleIter(tran1, ov1)
	assert.That(si1 != nil)

	// Test case 2: No layers, empty mut - should work
	ov2 := &Overlay{bt: bt, layers: []*ixbuf.T{{}}, mut: &ixbuf.T{}}
	tran2 := &simpleTestTran{overlay: ov2}
	si2 := NewSimpleIter(tran2, ov2)
	assert.That(si2 != nil)

	// Test case 3: Single empty layer, no mut - should work
	ov3 := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	tran3 := &simpleTestTran{overlay: ov3}
	si3 := NewSimpleIter(tran3, ov3)
	assert.That(si3 != nil)

	// Test case 4: Non-empty mut - should return nil
	ov4 := &Overlay{bt: bt, layers: []*ixbuf.T{{}}, mut: &ixbuf.T{}}
	ov4.mut.Insert("test", 1)
	tran4 := &simpleTestTran{overlay: ov4}
	si4 := NewSimpleIter(tran4, ov4)
	assert.That(si4 == nil)

	// Test case 5: Non-empty layer - should return nil
	ov5 := &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
	ov5.layers[0].Insert("test", 1)
	tran5 := &simpleTestTran{overlay: ov5}
	si5 := NewSimpleIter(tran5, ov5)
	assert.That(si5 == nil)

	// Test case 6: Multiple layers - should return nil
	ov6 := &Overlay{bt: bt, layers: []*ixbuf.T{{}, {}}}
	tran6 := &simpleTestTran{overlay: ov6}
	si6 := NewSimpleIter(tran6, ov6)
	assert.That(si6 == nil)
}
