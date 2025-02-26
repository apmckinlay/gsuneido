// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// package shmap is a Swiss hash table, based on the Go 1.24 map
// It uses the same terminology - map, table, group.
// Unlike Go, empty is defined as zero so we don't need to initialize groups.
// And keys and values are in separate arrays to avoid padding.
// Iteration is simpler because we don't allow modification during iteration.
// We have our own version to allow more flexibility with hash and equals.
package shmap

import (
	"hash/maphash"
	"math/bits"
	"slices"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// Map has a number of tables
type Map[K any, V any, H helper[K]] struct {
	// dir size is power of 2, indexed by first bits of hash
	dir []*table[K, V]
	// count is the number of entries in the map
	count int32
	// depth is the log2 of the dir size (number of tables)
	// e.g. 0 = 1 table, 1 = 2 tables, 2 = 4 tables, etc.
	depth int8

	// helper allows different methods of supplying Hash and Equal methods
	help H
	
	firstTable table[K,V]
	firstDir [1]*table[K, V]
}

// table has a number of groups
type table[K any, V any] struct {
	groups     []group[K, V] // size is power of 2
	count      int16
	growthLeft int16
	// depth is the depth of this table.
	// If < map.depth then this table will occupy multiple dir slots
	depth int8
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

const loadFactor = 7     // i.e. 7 / 8
const maxTableSize = 128 // groups = 1024 entries

//-------------------------------------------------------------------

type helper[K any] interface {
	Hash(k K) uint64
	Equal(x, y K) bool
}

// NewMapCmpable returns a map for comparable keys
func NewMapCmpable[K comparable, V any]() *Map[K, V, Cmpable[K]] {
	return &Map[K, V, Cmpable[K]]{help: Cmpable[K]{maphash.MakeSeed()}}
}

type Cmpable[K comparable] struct{ seed maphash.Seed }

func (c Cmpable[K]) Hash(k K) uint64 {
	return maphash.Comparable(c.seed, k)
}
func (Cmpable[K]) Equal(x, y K) bool {
	return x == y
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

// NewMapFuncs returns a Map using the provided hash and equals functions
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

func (m *Map[K, V, H]) Len() int {
	return int(m.count)
}

//-------------------------------------------------------------------

// Get does a normal search,
// stopping when the key is found or after a group with an empty slot.
func (m *Map[K, V, H]) Get(k K) (v V, ok bool) {
	nGets++
	if len(m.dir) == 0 {
		return
	}
	h2, tbl, seq := m.search(k)
	for ; ; seq = seq.next() {
		nGetProbes++
		grp := &tbl.groups[seq.offset]
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

		// ctrls := grp.control
		// hasEmpty := false
		// for i := range groupSize {
		// 	c := uint8(ctrls)
		// 	ctrls >>= 8
		// 	if c == h2|0x80 {
		// 		if m.help.Equal(k, grp.keys[i]) {
		// 			return grp.vals[i], true
		// 		}
		// 	} else if c == empty {
		// 		hasEmpty = true
		// 	}
		// }
		// assert.This(anyEmpty(grp.control)).Is(hasEmpty)
		// if hasEmpty {
		// 	return
		// }
	}
}

func (m *Map[K, V, H]) search(k K) (h2 uint8, tbl *table[K, V], seq probeSeq) {
	h := m.help.Hash(k)
	h1 := h >> 7
	h2 = uint8(h & 0x7f)
	ti := m.tableIndex(h)
	tbl = m.dir[ti]
	seq = makeProbeSeq(h1, uint64(len(tbl.groups)-1))
	return
}

var nGets = 0
var nGetProbes = 0

func (m *Map[K, V, H]) tableIndex(h uint64) int {
	shift := 64 - uint64(m.depth)
	return int((h >> shift) & uint64(len(m.dir)-1))
}

type bitset uint64

// findByte returns a one bit in the LSB of each byte that matches
func findByte(w uint64, b byte) bitset {
	x := uint64(b) * 0x0101010101010101
	x = w ^ ^x
	x = x & (x >> 1)
	x = x & (x >> 2)
	x = x & (x >> 4)
	x = x & 0x0101010101010101
	return bitset(x)
}

func (b bitset) first() int {
	return bits.TrailingZeros64(uint64(b)) / 8
}

func (b bitset) dropFirst() bitset {
	return b & (b - 1)
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
	if len(m.dir) == 0 {
		return zeroVal, false
	}
	h2, tbl, seq := m.search(k)
	for ; ; seq = seq.next() {
		grp := &tbl.groups[seq.offset]
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
						tbl.growthLeft++
					} else {
						grp.control |= deleted << (i * 8) // set to deleted
					}
					v := grp.vals[i]
					// zero to help gc
					grp.keys[i] = zeroKey
					grp.vals[i] = zeroVal
					tbl.count--
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

const hibits = 0x80808080_80808080

// anyEmpty returns true if any bit in the control word is empty
func anyEmpty(c uint64) bool {
	c = ^c // flip all bits, empty is now 0xff
	c = (c & (c << 1)) & hibits
	return c != 0
}

//-------------------------------------------------------------------

func (m *Map[K, V, H]) Put(k K, v V) {
	m.put(k, v, true)
}

func (m *Map[K, V, H]) GetPut(k K, v V) (K, V, bool) {
	return m.put(k, v, false)
}

// Put does a normal search, stopping probing at a group with an empty slot.
// If the key is found, and update is true it updates the value.
// Else if a deleted slot is found during the search, it stores there.
// Else it fills in the empty slot.
func (m *Map[K, V, H]) put(k K, v V, update bool) (k2 K, v2 V, existed bool) {
	nPuts++
	if len(m.dir) == 0 {
		m.init()
	}
	h := m.help.Hash(k)
	h1 := h >> 7
	h2 := uint8(h & 0x7f)
	grown := false
outer:
	for {
		ti := m.tableIndex(h)
		tbl := m.dir[ti]
		seq := makeProbeSeq(h1, uint64(len(tbl.groups)-1))
		probes := 1
		var del_grp *group[K, V]
		var del_i int
		for ; ; seq = seq.next() {
			nPutProbes++
			grp := &tbl.groups[seq.offset]
			ctrls := grp.control
			ie := -1
			for i := range groupSize {
				c := uint8(ctrls)
				ctrls >>= 8
				if c == h2|0x80 {
					if m.help.Equal(k, grp.keys[i]) {
						if update {
							grp.vals[i] = v
						}
						return grp.keys[i], grp.vals[i], true
					}
				} else if c == deleted {
					del_grp = grp
					del_i = i
				} else if c == empty && ie < 0 {
					ie = i
				}
			}
			if ie >= 0 {
				if tbl.growthLeft < 1 {
					assert.Msg("already grown").That(!grown)
					grown = true
					m.grow(ti, tbl)
					continue outer
				}
				// if we found a deleted slot during the search, insert there,
				// otherwise insert in the empty slot
				if del_grp != nil {
					grp = del_grp
					ie = del_i
					grp.control &= ^(0xff << (ie * 8)) // set to empty
				} else {
					tbl.growthLeft--
				}
				grp.keys[ie] = k
				grp.vals[ie] = v
				grp.control |= uint64(h2|0x80) << (ie * 8)
				m.count++
				tbl.count++
				return
			}
			if probes++; probes > len(tbl.groups) {
				panic("too many probes")
			}
		}
	}
}

var nPuts = 0
var nPutProbes = 0

// init sets up a single table with a single group
func (m *Map[K, V, H]) init() {
	// tbl := m.newTable(1)
	// m.dir = []*table[K, V]{tbl}
	m.firstTable = *m.newTable(1)
    m.firstDir[0] = &m.firstTable
	m.dir = m.firstDir[:]
}

//-------------------------------------------------------------------

// grow returns a new table twice the size of the current table
func (m *Map[K, V, H]) grow(ti int, tbl *table[K, V]) {
	if len(tbl.groups)*2 <= maxTableSize {
		// fmt.Println("double table", ti)
		m.dir[ti] = m.double(tbl)
	} else {
		// fmt.Println("split table", ti)
		tbl1, tbl2 := m.split(tbl)
		if tbl.depth == m.depth {
			// fmt.Println("doubleDir")
			m.doubleDir()
			m.dir[2*ti] = tbl1
			m.dir[2*ti+1] = tbl2
		} else {
			// fmt.Println("updateDir")
			m.updateDir(ti, tbl, tbl1, tbl2)
		}
	}
}

func (*Map[K, V, H]) newTable(ngroups int) *table[K, V] {
	return &table[K, V]{groups: make([]group[K, V], ngroups),
		growthLeft: int16(ngroups * loadFactor)}
}

// double copies to a new table twice the size
func (m *Map[K, V, H]) double(tbl *table[K, V]) *table[K, V] {
	old := *tbl // copy
	tbl2 := tbl // reuse
	*tbl2 = *m.newTable(len(old.groups) * 2)
	tbl2.depth = old.depth
	// copy the contents
	for gi := range old.groups {
		grp := &old.groups[gi]
		ctrls := grp.control
		for i := range groupSize {
			c := uint8(ctrls)
			ctrls >>= 8
			if c&0x80 != 0 {
				m.doublePut(tbl2, grp.keys[i], grp.vals[i])
			}
		}
	}
	assert.This(tbl2.count).Is(old.count)
	return tbl2
}

// doublePut is used by double.
// It is similar to Put but doesn't need to worry about existing keys
func (m *Map[K, V, H]) doublePut(tbl *table[K, V], k K, v V) {
	tbl.count++
	tbl.growthLeft--
	h := m.help.Hash(k)
	h1 := h >> 7
	h2 := uint8(h & 0x7f)
	seq := makeProbeSeq(h1, uint64(len(tbl.groups)-1))
	for ; ; seq = seq.next() {
		grp := &tbl.groups[seq.offset]
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

// split copies to two new tables
func (m *Map[K, V, H]) split(tbl *table[K, V]) (*table[K, V], *table[K, V]) {
	tbl1 := m.newTable(len(tbl.groups))
	tbl1.depth = tbl.depth + 1
	tbl2 := m.newTable(len(tbl.groups))
	tbl2.depth = tbl.depth + 1
	// copy the contents
	for gi := range tbl.groups {
		grp := &tbl.groups[gi]
		ctrls := grp.control
		for i := range groupSize {
			c := uint8(ctrls)
			ctrls >>= 8
			if c&0x80 != 0 {
				m.splitPut(tbl1, tbl2, grp.keys[i], grp.vals[i])
			}
		}
	}
	assert.This(tbl1.count + tbl2.count).Is(tbl.count)
	return tbl1, tbl2
}

// splitPut is used by split.
// It is similar to Put but doesn't need to worry about existing keys
func (m *Map[K, V, H]) splitPut(tbl1, tbl2 *table[K, V], k K, v V) {
	h := m.help.Hash(k)
	h1 := h >> 7
	h2 := uint8(h & 0x7f)
	tbl := tbl1
	shift := 64 - uint64(tbl1.depth)
	if (h>>shift)&1 == 1 {
		tbl = tbl2
	}
	tbl.count++
	tbl.growthLeft--
	seq := makeProbeSeq(h1, uint64(len(tbl.groups)-1))
	for ; ; seq = seq.next() {
		grp := &tbl.groups[seq.offset]
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

func (m *Map[K, V, H]) doubleDir() {
	dir := make([]*table[K, V], len(m.dir)*2)
	for i, tbl := range m.dir {
		dir[2*i] = tbl
		dir[2*i+1] = tbl
	}
	m.dir = dir
	m.depth++
	// m.print()
}

const ones64 uint64 = 0xffff_ffff_ffff_ffff

// updateDir updates the directory for a split
func (m *Map[K, V, H]) updateDir(ti int, tbl, tbl1, tbl2 *table[K, V]) {
	delta := m.depth - tbl.depth
	assert.That(delta >= 1)
	ti &= int(ones64 << delta) // the first "copy" of the table
	ncopy := 1 << delta
	for i := range ncopy {
		if i < ncopy/2 {
			m.dir[ti+i] = tbl1
		} else {
			m.dir[ti+i] = tbl2
		}
	}
}

//-------------------------------------------------------------------

// Iter returns a function that returns the next key and value from the map.
// WARNING: It does not handle or catch modification during iteration.
// It is up to the caller to prevent this.
func (m *Map[K, V, H]) Iter() func() (K, V, bool) {
	if len(m.dir) == 0 {
		return func() (k K, v V, ok bool) { return }
	}
	ti := 0
	tbl := m.dir[0]
	gi := 0
	grp := &tbl.groups[0]
	i := -1
	return func() (k K, v V, ok bool) {
		if ti >= len(m.dir) {
			return
		}
		for {
			if i++; i >= groupSize {
				i = 0
				if gi++; gi >= len(tbl.groups) {
					gi = 0
					if ti++; ti >= len(m.dir) {
						return
					}
					tbl = m.dir[ti]
				}
				grp = &tbl.groups[gi]
			}
			c := uint8(grp.control >> (i * 8))
			if c&0x80 != 0 {
				return grp.keys[i], grp.vals[i], true
			}
		}
	}
}

//-------------------------------------------------------------------

// Copy makes a shallow copy of the map
func (m *Map[K, V, H]) Copy() *Map[K, V, H] {
	newMap := *m
	newMap.dir = slices.Clone(m.dir)
	for ti, tbl := range m.dir {
		t := *tbl
		t.groups = slices.Clone(tbl.groups)
		newMap.dir[ti] = &t
	}
	return &newMap
}

//-------------------------------------------------------------------

// Clear removes all the entries but keeps the allocated memory
func (m *Map[K, V, H]) Clear() {
	for ti := range m.dir {
		tbl := m.dir[ti]
		tbl.count = 0
		tbl.growthLeft = int16(len(tbl.groups) * loadFactor)
		clear(tbl.groups)
	}
	m.count = 0
}
