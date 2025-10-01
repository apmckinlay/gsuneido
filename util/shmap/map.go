// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// package shmap is a Swiss hash table, based on the Go 1.24 map
// We have our own version to allow more flexibility with key hash and equals.
// It is just a single table, not multiple like Go - simpler, faster.
// Not as good for very large tables, but that is ok for our usage.
// Unlike Go, empty is defined as zero so we don't need to initialize groups.
// And keys and values are in separate arrays to avoid padding.
// Iteration is simpler because we don't support modification during iteration.
package shmap

import (
	"math/bits"
	"slices"

	"github.com/apmckinlay/gsuneido/util/assert"
)

type Map[K any, V any, H helper[K]] struct {
	groups     []group[K, V] // size is power of 2
	count      int32
	growthLeft int32
	help       H
}

// group holds 8 keys and values, plus 8 bytes of metadata
type group[K any, V any] struct {
	// control holds low 7 bits of hashes, and marks empty and tombstone
	control uint64
	// separate key and val arrays, to avoid padding
	keys [groupSize]K
	vals [groupSize]V
}

const groupSize = 8 // to match bytes in uint64 control word

const empty = 0      // high bit = zero
const deleted = 0x7f // high bit = zero
// an occupied slot is the 7 bit hash | 0x80 (high bit = one)

const loadFactor = 7 // i.e. 7 / 8

//-------------------------------------------------------------------

// helper is used to allow custom hash and equals functions
type helper[K any] interface {
	Hash(k K) uint64
	Equal(x, y K) bool
}

// NewMapMeth returns a Map that calls Hash and Equal methods on the key
func NewMapMeth[K Key, V any]() *Map[K, V, Meth[K]] {
	return &Map[K, V, Meth[K]]{help: Meth[K]{}}
}

type Meth[K Key] struct{}

type Key interface {
	Hash() uint64
	Equal(any) bool
}

func (Meth[K]) Hash(k K) uint64 {
	return k.Hash()
}
func (Meth[K]) Equal(x, y K) bool {
	return x.Equal(y)
}

// NewMapFuncs returns a Map using the provided hash and equals functions.
// Using closures allows accessing context.
func NewMapFuncs[K any, V any](
	hfn func(k K) uint64,
	eqfn func(x, y K) bool) *Map[K, V, Funcs[K]] {
	return &Map[K, V, Funcs[K]]{help: Funcs[K]{hfn: hfn, eqfn: eqfn}}
}

type Funcs[K any] struct {
	hfn  func(k K) uint64
	eqfn func(x, y K) bool
}

func (fs Funcs[K]) Hash(k K) uint64 {
	return fs.hfn(k)
}
func (fs Funcs[K]) Equal(x, y K) bool {
	return fs.eqfn(x, y)
}

//-------------------------------------------------------------------

// Size returns the number of elements in the map
func (m *Map[K, V, H]) Size() int {
	return int(m.count)
}

// Has returns true if the key exists in the table, false otherwise
func (m *Map[K, V, H]) Has(key K) bool {
	_, ok := m.Get(key)
	return ok
}

// Get returns the value and true if the key exists
// otherwise it returns the zero value and false.
func (m *Map[K, V, H]) Get(k K) (v V, ok bool) {
	// nGets++
	if m.count == 0 {
		return
	}
	// Get does a normal search,
	// stopping when the key is found or after a group with an empty slot.
	h2, seq := m.search(k)
	for ; ; seq = seq.next() {
		// nGetProbes++
		grp := &m.groups[seq.offset]
		b := findByte(grp.control, h2|0x80)
		for b != 0 {
			i := b.first()
			if m.help.Equal(k, grp.keys[i]) {
				return grp.vals[i], true
			}
			b = b.dropFirst()
		}
		if anyEmpty(grp.control) {
			return
		}
	}
}

func (m *Map[K, V, H]) search(k K) (uint8, probeSeq) {
	h := m.help.Hash(k)
	h1 := h >> 7
	h2 := uint8(h & 0x7f)
	seq := makeProbeSeq(h1, uint64(len(m.groups)-1))
	return h2, seq
}

// var nGets = 0
// var nGetProbes = 0

//-------------------------------------------------------------------

// bitset has 8 bytes of 0 or 1
type bitset uint64

const lsbits = 0x0101010101010101
const msbits = 0x8080808080808080

// findByte returns a one bit in the LSB of each byte that matches
// see: Hacker's Delight, Chapter 6 Searching Words 6.2
// Find First String of 1-bits of a Given length
func findByte(w uint64, b byte) bitset {
	x := uint64(b) * lsbits
	x = w ^ ^x
	x = x & (x >> 1)
	x = x & (x >> 2)
	x = x & (x >> 4)
	x = x & lsbits
	return bitset(x)
}

func (b bitset) first() int {
	return bits.TrailingZeros64(uint64(b)) / 8
}

func (b bitset) dropFirst() bitset {
	return b & (b - 1)
}

func emptyOrDeleted(w uint64) bitset {
	return bitset(^(w >> 7) & lsbits)
}

// anyEmpty returns true if any bit in the control word is empty
func anyEmpty(c uint64) bool {
	c = ^c // flip all bits, empty is now 0xff
	c = (c & (c << 1)) & msbits
	return c != 0
}

//-------------------------------------------------------------------

// Del deletes a key from the map
// and returns the old value and whether the key existed.
// It does a normal search, stopping probing at a group with an empty slot.
// If the group where the key is found has any empty slots,
// then we set the found slot to empty.
// Otherwise we set the found slot to deleted (tombstone).
func (m *Map[K, V, H]) Del(k K) (V, bool) {
	var zeroVal V
	var zeroKey K
	if m.count == 0 {
		return zeroVal, false
	}
	h2, seq := m.search(k)
	for ; ; seq = seq.next() {
		grp := &m.groups[seq.offset]
		ctrls := grp.control
		hasEmpty := false
		for i := range groupSize {
			c := uint8(ctrls)
			ctrls >>= 8
			if c == h2|0x80 {
				if m.help.Equal(k, grp.keys[i]) {
					// found, delete it
					ae := anyEmpty(grp.control)
					grp.control &= ^(0xff << (i * 8)) // set to empty
					if ae {
						m.growthLeft++
					} else {
						grp.control |= deleted << (i * 8) // set to deleted
					}
					v := grp.vals[i]
					// zero to help gc
					grp.keys[i] = zeroKey
					grp.vals[i] = zeroVal
					m.count--
					return v, true
				}
			} else if c == empty {
				hasEmpty = true
			}
		}
		if hasEmpty {
			return zeroVal, false
		}
	}
}

//-------------------------------------------------------------------

// Put adds or updates a key-value entry in the map.
func (m *Map[K, V, H]) Put(k K, v V) {
	m.put(k, v, true)
}

// GetInit creates a new entry in the map (with zero value) if it doesn't exist,
// It returns the (original existing) key and whether it existed.
func (m *Map[K, V, H]) GetInit(k K) (K, bool) {
	var zero V
	return m.put(k, zero, false)
}

// put does a normal search, stopping probing at a group with an empty slot.
// If the key is found, and update is true it updates the value.
// Else if a deleted slot is found during the search, it stores there.
// Else it fills in the empty slot.
func (m *Map[K, V, H]) put(k K, v V, update bool) (k2 K, existed bool) {
	// nPuts++
	if len(m.groups) == 0 {
		m.init()
	}
	h := m.help.Hash(k)
	h1 := h >> 7
	h2 := uint8(h & 0x7f)
	grown := false
outer:
	for {
		seq := makeProbeSeq(h1, uint64(len(m.groups)-1))
		probes := 1
		var del_grp *group[K, V]
		var del_i int
		for ; ; seq = seq.next() {
			// nPutProbes++
			grp := &m.groups[seq.offset]
			ctrls := grp.control
			ie := -1
			b := findByte(ctrls, h2|0x80) | emptyOrDeleted((ctrls))
			for b != 0 {
				i := b.first()
				c := uint8(ctrls >> (i * 8))
				if c == deleted {
					del_grp = grp
					del_i = i
				} else if c == empty {
					if ie < 0 {
						ie = i
					}
				} else if m.help.Equal(k, grp.keys[i]) {
					if update {
						grp.vals[i] = v
					}
					return grp.keys[i], true // found it
				}
				b = b.dropFirst()
			}
			if ie >= 0 {
				// if we found a deleted slot during the probes, insert there,
				// otherwise insert in the empty slot
				if del_grp != nil {
					grp = del_grp
					ie = del_i
					grp.control &= ^(0xff << (ie * 8)) // set to empty
				} else {
					if m.growthLeft < 1 {
						assert.That(!grown)
						grown = true
						m.grow()
						continue outer
					}
					m.growthLeft--
				}
				grp.keys[ie] = k
				grp.vals[ie] = v
				grp.control |= uint64(h2|0x80) << (ie * 8)
				m.count++
				return k, false
			}
			if probes++; probes > len(m.groups) {
				panic("too many probes")
			}
		}
	}
}

// var nPuts = 0
// var nPutProbes = 0

// init sets up a single group
func (m *Map[K, V, H]) init() {
	m.groups = make([]group[K, V], 1)
	m.growthLeft = loadFactor
}

// grow copies to a new table twice the size
func (m *Map[K, V, H]) grow() {
	groups := make([]group[K, V], 2*len(m.groups))
	// copy the contents
	for gi := range m.groups {
		grp := &m.groups[gi]
		ctrls := grp.control
		for i := range groupSize {
			c := uint8(ctrls)
			ctrls >>= 8
			if c&0x80 != 0 {
				m.growPut(groups, grp.keys[i], grp.vals[i])
			}
		}
	}
	m.groups = groups
	m.growthLeft = int32(len(groups)*loadFactor) - m.count
}

// growPut copies to a new bigger group slice.
// It is similar to Put but doesn't need to worry about duplicate keys
func (m *Map[K, V, H]) growPut(groups []group[K, V], k K, v V) {
	h := m.help.Hash(k)
	h1 := h >> 7
	h2 := uint8(h & 0x7f)
	seq := makeProbeSeq(h1, uint64(len(groups)-1))
	for ; ; seq = seq.next() {
		grp := &groups[seq.offset]
		ctrls := grp.control
		for i := range groupSize {
			c := uint8(ctrls)
			ctrls >>= 8
			if c == empty {
				grp.keys[i] = k
				grp.vals[i] = v
				grp.control |= uint64(h2|0x80) << (i * 8)
				return
			}
		}
	}
}

//-------------------------------------------------------------------

// Iter returns a function that returns the next key and value from the map.
// WARNING: It does not handle or detect modification during iteration.
// It is up to the caller to prevent this.
func (m *Map[K, V, H]) Iter() func() (K, V, bool) {
	if m.count == 0 {
		return func() (k K, v V, ok bool) { return }
	}
	gi := 0
	grp := &m.groups[0]
	i := -1
	return func() (k K, v V, ok bool) {
		if gi >= len(m.groups) {
			return
		}
		for {
			if i++; i >= groupSize {
				i = 0
				if gi++; gi >= len(m.groups) {
					return
				}
				grp = &m.groups[gi]
			}
			c := uint8(grp.control >> (i * 8))
			if c&0x80 != 0 {
				return grp.keys[i], grp.vals[i], true
			}
		}
	}
}

// Copy makes a shallow copy of the map
func (m *Map[K, V, H]) Copy() *Map[K, V, H] {
	newMap := *m
	newMap.groups = slices.Clone(m.groups)
	return &newMap
}

// Clear removes all the entries but keeps the allocated memory
func (m *Map[K, V, H]) Clear() {
	clear(m.groups)
	m.count = 0
	m.growthLeft = int32(len(m.groups) * loadFactor)
}

//-------------------------------------------------------------------

// probeSeq maintains the state for a probe sequence that iterates through the
// groups in a table. The sequence is a triangular progression of the form
//
//	p(i) := (i^2 + i)/2 + hash (mod mask+1)
//
// The sequence effectively outputs the indexes of *groups*. The group
// machinery allows us to check an entire group with minimal branching.
//
// It turns out that this probe sequence visits every group exactly once if
// the number of groups is a power of two, since (i^2+i)/2 is a bijection in
// Z/(2^m). See https://en.wikipedia.org/wiki/Quadratic_probing
type probeSeq struct {
	mask   uint64
	offset uint64
	index  uint64
}

func makeProbeSeq(hash uint64, mask uint64) probeSeq {
	return probeSeq{
		mask:   mask,
		offset: hash & mask,
		index:  0,
	}
}

func (s probeSeq) next() probeSeq {
	s.index++
	s.offset = (s.offset + s.index) & s.mask
	return s
}
