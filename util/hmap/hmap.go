// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package hmap implements a hash map
// Based on: https://github.com/skarupke/flat_hash_map bytell_hash_map
package hmap

import (
	"github.com/apmckinlay/gsuneido/util/verify"
)

// metaData is a single byte of information about each slot
// one bit is used for whether the slot is "direct" or not
// seven bits are used for the jump index that chains slots
// zero is used for empty so new arrays don't have to be initialized
type metaData byte

const (
	metaEmpty    = 0
	metaDirect   = 0x80
	metaJumpBits = 0x7f
	noJump       = 0x7f
)

// jump returns the jump part of the metaData
func (m metaData) jump() int {
	return int(m & metaJumpBits)
}

// withJump returns a metaData with the jump index set to the given value
func (m metaData) withJump(ji int) metaData {
	return (m & metaDirect) | metaData(ji)
}

// jumpSize is the possible relative jumps in a chain of slots
// start with sequential values for locality
// then increasingly bigger values to cover the table
var jumpSize = [...]int{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,

	21, 28, 36, 45, 55, 66, 78, 91, 105, 120, 136, 153, 171, 190, 210, 231,
	253, 276, 300, 325, 351, 378, 406, 435, 465, 496, 528, 561, 595, 630,
	666, 703, 741, 780, 820, 861, 903, 946, 990, 1035, 1081, 1128, 1176,
	1225, 1275, 1326, 1378, 1431, 1485, 1540, 1596, 1653, 1711, 1770, 1830,
	1891, 1953, 2016, 2080, 2145, 2211, 2278, 2346, 2415, 2485, 2556,

	3741, 8385, 18915, 42486, 95703, 215496, 485605, 1091503, 2456436,
	5529475, 12437578, 27986421, 62972253, 141700195, 318819126, 717314626,
}

// Val holds any value
type Val = interface{}

// Key is anything with Hash and Equal
type Key interface {
	Hash() uint32
	Equal(Val) bool
}

const (
	blockSize      = 8 // must be power of two
	blockSizeMask  = 7
	blockSizeShift = 3
)

// The table is organized into blocks with separate arrays for meta,key,val
// This is to eliminate padding while still maintaining locality
type block struct {
	meta [blockSize]metaData
	key  [blockSize]Key
	val  [blockSize]Val
}

// Hmap is the hash table
// Its zero value is a valid empty table
type Hmap struct {
	blocks   []block
	size     int32
	version  uint16
	capShift uint8
	growing  bool
}

// cap returns the raw capacity
func (h *Hmap) cap() int {
	return len(h.blocks) * blockSize
}

// Size returns the current number of elements in the table
func (h *Hmap) Size() int {
	return int(h.size)
}

// isFull returns true if adding an element would exceed 7/8 full
func (h *Hmap) isFull() bool {
	return int(h.size)+1 > 7*h.cap()/8 // .875
}

// whichBlock returns the index into the blocks array
func whichBlock(i int) int {
	return i >> blockSizeShift
}

// indexInBlock returns the index into the arrays in a block
func indexInBlock(i int) int {
	return i & blockSizeMask
}

// hashToIndex uses Fibonaci
func (h *Hmap) hashToIndex(hash uint32) int {
	const phi32 = 2654435769
	return int((hash * phi32) >> h.capShift)
	//return int(hash % h.cap()) // simpler for debugging
}

// keepInRange is used after advancing an index to make it wrap around
func (h *Hmap) keepInRange(index int) int {
	return index & (h.cap() - 1)
}

// Get returns the value for the Key or nil if not found
func (h *Hmap) Get(key Key) Val {
	if h.size == 0 {
		return nil
	}
	// for performance, don't use an iterator for the first probe
	index := h.hashToIndex(key.Hash())
	b := &h.blocks[whichBlock(index)]
	ib := indexInBlock(index)
	if (b.meta[ib] & metaDirect) != metaDirect {
		return nil
	}
	if key.Equal(b.key[ib]) {
		return b.val[ib]
	}
	iter := chainIter{h, b, index, ib}
	for iter.next() {
		if key.Equal(iter.key()) {
			return iter.val()
		}
	}
	return nil
}

// Put adds or updates an entry
func (h *Hmap) Put(key Key, val Val) {
	h.version++
	if h.cap() == 0 {
		h.grow()
	}
	index := h.hashToIndex(key.Hash())
	b := &h.blocks[whichBlock(index)]
	ib := indexInBlock(index)
	if b.meta[ib] == metaEmpty {
		b.meta[ib] = metaDirect | noJump
		b.key[ib] = key
		b.val[ib] = val
		h.size++
		return
	}
	iter := chainIter{h, b, index, ib}
	if iter.meta()&metaDirect != metaDirect {
		h.putDirect(&iter, key, val)
		return
	}
	for {
		if key.Equal(iter.key()) {
			iter.valSet(val) // update value for existing key
			h.version--      // version doesn't need to change
			return
		}
		if !iter.next() {
			h.putChain(&iter, key, val)
			return
		}
	}
}

// grow initially makes the capacity one block, and then doubles after that
func (h *Hmap) grow() {
	if h.cap() == 0 {
		h.blocks = make([]block, 1)
		h.capShift = 32 - 3
		return
	}
	if h.growing {
		panic("grow while growing")
	}
	h.growing = true
	oldblocks := h.blocks
	h.blocks = make([]block, 2*len(oldblocks))
	h.capShift--
	h.size = 0
	for _, b := range oldblocks {
		for ib := 0; ib < blockSize; ib++ {
			if b.meta[ib] != metaEmpty {
				h.Put(b.key[ib], b.val[ib])
			}
		}
	}
	h.growing = false
}

// putDirect starts a new chain
// i.e. it handles when this is the first entry with a certain hash index
func (h *Hmap) putDirect(slot *chainIter, key Key, val Val) {
	if h.isFull() {
		h.grow()
		h.Put(key, val) // recursive restart
		return
	}
	if slot.meta() != metaEmpty {
		// move the chain passing through here
		oldchain := *slot
		prev := h.findPrev(slot)
		for {
			ji, free := h.findEmpty(prev.index)
			if ji == -1 {
				h.grow()
				h.Put(key, val) // recursive restart
				return
			}
			prev.jumpSet(ji)
			free.set(noJump, oldchain.key(), oldchain.val())
			jump := oldchain.jump()
			if oldchain.index != slot.index { // keep direct slot from being used
				oldchain.set(metaEmpty, nil, nil)
			}
			if !oldchain.nextWithJump(jump) {
				break
			}
			prev = free
		}
	}
	slot.set(metaDirect|noJump, key, val)
	h.size++
}

// findPrev finds the previous element in a chain
// because chains are single linked, we search from the beginning of the chain
// panics if unsuccessful
func (h *Hmap) findPrev(slot *chainIter) chainIter {
	iter := h.iterFromKey(slot.key())
	for iter.nextIndex() != slot.index {
		verify.That(iter.next())
	}
	return iter
}

// findEmpty finds a jump from a given index to an empty slot
func (h *Hmap) findEmpty(index int) (int, chainIter) {
	for ji := range jumpSize {
		free := h.keepInRange(index + jumpSize[ji])
		b := &h.blocks[whichBlock(free)]
		ib := indexInBlock(free)
		if b.meta[ib] == metaEmpty {
			return ji, chainIter{h, b, free, ib}
		}
	}
	return -1, chainIter{}
}

// putChain adds to the end of a chain
func (h *Hmap) putChain(iter *chainIter, key Key, val Val) {
	if h.isFull() {
		h.grow()
		h.Put(key, val) // recursive restart
		return
	}
	ji, free := h.findEmpty(iter.index)
	if ji == -1 {
		h.grow()
		h.Put(key, val) // recursive restart
		return
	}
	iter.jumpSet(ji)
	free.set(noJump, key, val)
	h.size++
}

// Del deletes a key and returns its old value, or nil if it didn't exist
func (h *Hmap) Del(key Key) Val {
	if h.size == 0 {
		return nil
	}
	iter := h.iterFromKey(key)
	if (iter.meta() & metaDirect) != metaDirect {
		return nil // hash index does not exist in the table
	}
	var prev chainIter
	for !key.Equal(iter.key()) {
		prev = iter
		if !iter.next() {
			return nil // end of chain
		}
	}
	h.version++
	val := iter.val()
	if iter.jump() != noJump {
		// delete from within chain - move element at end of chain
		slot := iter
		for iter.jump() != noJump {
			prev = iter
			iter.next()
		}
		slot.keySet(iter.key())
		slot.valSet(iter.val())
	}
	iter.set(metaEmpty, nil, nil)
	if prev.b != nil {
		prev.jumpSet(noJump)
	}
	h.size--
	return val
}

// Copy returns a shallow copy of the Hmap
func (h *Hmap) Copy() *Hmap {
	hnew := Hmap{size: h.size, capShift: h.capShift}
	hnew.blocks = make([]block, len(h.blocks))
	copy(hnew.blocks, h.blocks)
	return &hnew
}

// Iter returns a function (closure) that is called to get the next item.
// It returns nil,nil at the end.
func (h *Hmap) Iter() func() (Key, Val) {
	i := len(h.blocks) * blockSize
	ver := h.version
	return func() (Key, Val) {
		if ver != h.version {
			panic("hmap modified during iteration")
		}
		for i > 0 {
			i--
			b := &h.blocks[whichBlock(i)]
			ib := indexInBlock(i)
			if b.meta[ib] != metaEmpty {
				return b.key[ib], b.val[ib]
			}
		}
		return nil, nil // end
	}
}

//-------------------------------------------------------------------

// chainIter points to a slot in the table
type chainIter struct {
	h     *Hmap
	b     *block
	index int
	ib    int
}

// iterFromKey creates an iterator from a key
func (h *Hmap) iterFromKey(key Key) chainIter {
	index := h.hashToIndex(key.Hash())
	b := &h.blocks[whichBlock(index)]
	ib := indexInBlock(index)
	return chainIter{h, b, index, ib}
}

// next advances to the next slot in the chain
// it returns true if it was able to advance, false if at the end of the chain
func (it *chainIter) next() bool {
	return it.nextWithJump(it.jump())
}

// nextWithJump advances by a given jump
// it returns false if the jump == noJump, otherwise true
func (it *chainIter) nextWithJump(jump int) bool {
	if jump == noJump {
		return false // end of chain
	}
	it.index = it.h.keepInRange(it.index + jumpSize[jump])
	it.b = &it.h.blocks[whichBlock(it.index)]
	it.ib = indexInBlock(it.index)
	return true
}

// meta returns the metaData for an iterator's slot
func (it *chainIter) meta() metaData {
	return it.b.meta[it.ib]
}

// key returns the key for an iterator's slot
func (it *chainIter) key() Key {
	return it.b.key[it.ib]
}

// keySet sets the key for an iterator's slot
func (it *chainIter) keySet(key Key) {
	it.b.key[it.ib] = key
}

// val returns the value for an iterator's slot
func (it *chainIter) val() Val {
	return it.b.val[it.ib]
}

// valSet sets the value for an iterator's slot
func (it *chainIter) valSet(val Val) {
	it.b.val[it.ib] = val
}

// jump returns the jump part of the metaData for an iterator's slot
func (it *chainIter) jump() int {
	return it.b.meta[it.ib].jump()
}

// jumpSet sets the jump part of the metaData for an iterator's slot
func (it *chainIter) jumpSet(jump int) {
	meta := &it.b.meta[it.ib]
	*meta = meta.withJump(jump)
}

// set updates the contents of a slot
func (it *chainIter) set(meta metaData, key Key, val Val) {
	it.b.meta[it.ib] = meta
	it.b.key[it.ib] = key
	it.b.val[it.ib] = val
}

// nextIndex returns the index of the slot after this one in the chain
func (it *chainIter) nextIndex() int {
	return it.h.keepInRange(it.index + jumpSize[it.jump()])
}
