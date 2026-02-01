// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ixbuf defines an ordered list ixbuf.T
// with a mutating Insert and a immutable persistent Merge.
// It is designed for intermediate numbers of values, e.g. up to 16k or so.
// It is used mutably for per transaction index buffers
// and immutably for global index buffers.
package ixbuf

import (
	"cmp"
	"fmt"
	"log"
	"strings"

	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/slc"
)

const Update = 1 << 62
const Delete = 1 << 63
const Insert = 0

const Mask = 0xffffffffff

type T = ixbuf

type ixbuf struct {
	chunks []chunk
	size   int32
	// modCount is used by Iterator to detect modifications.
	// No locking since ixbuf is thread contained when mutable.
	// Incremented by Insert and Clear.
	modCount int32
}

type chunk []slot

type slot struct {
	key string
	off uint64
}

func (ib *ixbuf) Len() int {
	return int(ib.size)
}

func (ib *ixbuf) Clear() {
	ib.size = 0
	ib.chunks = nil
	ib.modCount++
}

// goal is the desired chunk size for a given item count.
// It is chosen so the size of the chunk list is roughly the chunk size.
func goal(n int32) int {
	// add 50% because average size is 2/3 full
	switch {
	case n < 256:
		return 16 + 16/2
	case n < 1024:
		return 32 + 32/2
	case n < 4*1024:
		return 64 + 64/2
	case n < 16*1024:
		return 128 + 128/2
	case n < 64*1024:
		return 256 + 256/2
	default:
		return 512 + 512/2
	}
}

// Insert adds an element. It mutates and is NOT thread-safe.
// off can have the update or delete bit set.
func (ib *ixbuf) Insert(key string, off uint64) (oldoff uint64) {
	if off == 0 {
		panic("ixbuf Insert: offset cannot be zero")
	}
	ib.modCount++
	if len(ib.chunks) == 0 {
		ib.size++
		ib.chunks = make([]chunk, 1, 4)   // ???
		ib.chunks[0] = make([]slot, 1, 8) // ???
		ib.chunks[0][0] = slot{key: key, off: off}
		return
	}
	ci := ib.searchChunks(key)
	c := ib.chunks[ci]
	i := search(ib.chunks[ci], key)

	if i < len(c) && c[i].key == key {
		// already exists, combine
		slot := &c[i]
		slot.off, oldoff = Combine(slot.off, off) // handles update/delete
		if slot.off == 0 {
			ib.remove(ci, i)
		}
		return
	}

	// insert in place
	ib.size++
	c = append(c, slot{})
	ib.chunks[ci] = c
	copy(c[i+1:], c[i:])
	c[i] = slot{key: key, off: off}

	if len(c) > goal(ib.size) {
		// split
		n := len(c)
		at := n / 2
		if i == 0 {
			at = n / 4
		} else if i == len(c)-1 {
			at += n / 4
		}
		left := c[:at] // re-use
		right := make([]slot, n-at)
		copy(right, c[at:])
		ib.chunks[ci] = left
		ib.chunks = append(ib.chunks, nil)
		ci++
		copy(ib.chunks[ci+1:], ib.chunks[ci:])
		ib.chunks[ci] = right
	}
	return
}

func (ib *ixbuf) remove(ci int, i int) {
	c := ib.chunks[ci]
	if len(c) == 1 {
		ib.chunks = append(ib.chunks[:ci], ib.chunks[ci+1:]...)
	} else {
		ib.chunks[ci] = append(c[:i], c[i+1:]...)
	}
	ib.size--
}

func (ib *ixbuf) search(key string) (int, chunk, int) {
	ci := ib.searchChunks(key)
	c := ib.chunks[ci]
	i := search(ib.chunks[ci], key)
	return ci, c, i
}

// searchChunks does a binary search of the first key in each chunk.
// It returns len-1 if the key is greater than all keys.
func (ib *ixbuf) searchChunks(key string) int {
	i, j := 0, len(ib.chunks)
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		c := ib.chunks[h]
		if key > c.lastKey() {
			i = h + 1
		} else {
			j = h
		}
	}
	return min(i, len(ib.chunks)-1)
}

// search does a binary search of one chunk for a key.
// It returns len if the key is greater than all keys in the chunk.
func search(c chunk, key string) int {
	i, j := 0, len(c)
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		if key > c[h].key {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}

// Update combines if the key exists, otherwise it adds an update entry
func (ib *ixbuf) Update(key string, off uint64) uint64 {
	return ib.Insert(key, off|Update)
}

// Delete combines if the key exists, otherwise it adds a delete tombstone.
func (ib *ixbuf) Delete(key string, off uint64) uint64 {
	return ib.Insert(key, off|Delete)
}

//-------------------------------------------------------------------

const (
	add_update    = 0b_00_01
	add_delete    = 0b_00_10
	update_update = 0b_01_01
	update_delete = 0b_01_10
	delete_add    = 0b_10_00
	// INVALID:
	delete_delete = 0b10_10
	delete_update = 0b10_01
	// add_add    = 0b00_00
	// update_add = 0b01_00
)

func Combine(off1, off2 uint64) (result uint64, oldoff uint64) {
	ops := off1>>60 | off2>>62
	switch ops {
	case add_update: // => add
		result = off2 & Mask
	case add_delete: // => <nil>
		result = 0 // => should be removed
	case update_update: // => update
		result = off2
		oldoff = off1 & Mask
		// oldoff is so tran.Update can detect update,update of same record
	case update_delete: // => delete
		// update_delete needs to keep delete (not return 0)
		// because the add is in another layer
		result = off2
		oldoff = off1 & Mask
		// oldoff is so tran.Delete can detect update,delete of same record
	case delete_add: // => update
		result = off2 | Update
	case delete_delete:
		panic("delete & delete on same record")
	case delete_update:
		panic("delete & update on same record")
	default:
		log.Println("ERROR: ixbuf invalid Combine", OffString(off1), OffString(off2))
		dbg.PrintStack()
		panic("ixbuf invalid Combine")
	}
	return
}

//-------------------------------------------------------------------

// Merge combines several ixbuf's into a new one.
// It does not modify its inputs so it is thread-safe
// as long as the inputs don't change.
// It is immutable persistent and the result may share chunks of the inputs
// so again, the inputs cannot change.
func Merge(ibs ...*ixbuf) *ixbuf {
	assert.That(len(ibs) > 1)
	in := make([][]chunk, 0, len(ibs))
	size := int32(0)
	nc := 0
	var single *ixbuf
	for _, ib := range ibs {
		if ib.size == 0 {
			continue
		}
		in = append(in, ib.chunks)
		size += ib.size
		nc += len(ib.chunks)
		single = ib
	}
	if len(in) == 0 {
		return &ixbuf{}
	} else if len(in) == 1 {
		return single
	}
	m := merge{goal: goal(size), in: in, out: make([]chunk, 0, nc)}
	return m.merge()
}

type merge struct {
	in   [][]chunk
	buf  chunk
	out  []chunk
	goal int
	size int
}

func (m *merge) merge() *ixbuf {
	in := make([]chunk, len(m.in))
	for i := range m.in {
		in[i] = m.in[i][0]
		m.in[i] = m.in[i][1:]
	}
	passthru := false
	for len(m.in) > 0 {
		// find minimum, inlined for speed
		i := 0
		key := in[0].firstKey()
		for j := 1; j < len(in); j++ {
			key2 := in[j].firstKey()
			if key2 < key {
				i = j
				key = key2
			}
		}
		// i now has the minimum key
		if !passthru || !m.passthru(in, i) {
			m.outputSlot(in[i][0])
			if len(in[i]) > 1 {
				// advance to next slot
				in[i] = in[i][1:]
				continue
			}
		}
		if len(m.in[i]) > 0 {
			// advance to next chunk
			in[i] = m.in[i][0]
			m.in[i] = m.in[i][1:]
		} else {
			// remove empty input
			in = append(in[:i], in[i+1:]...)
			m.in = append(m.in[:i], m.in[i+1:]...)
		}
		passthru = true
	}
	return m.result()
}

func (m *merge) outputSlot(s2 slot) {
	last := len(m.buf) - 1
	if last >= 0 {
		s1 := &m.buf[last]
		if s2.key == s1.key {
			s1.off, _ = Combine(s1.off, s2.off)
			if s1.off == 0 {
				m.buf = m.buf[:last]
			}
			return
		}
	}
	if len(m.buf) > m.goal {
		m.flushbuf()
	}
	m.buf = append(m.buf, s2)
}

func (m *merge) passthru(in []chunk, i int) bool {
	// i is the minimum
	lastkey := in[i].lastKey()
	for j := range len(m.in) {
		if j != i && lastkey >= in[j].firstKey() {
			return false
		}
	}
	// if the chunk updates the previous, we can't pass through
	last := len(m.buf) - 1
	if last >= 0 && in[i].firstKey() == m.buf[last].key {
		return false
	}
	m.outputChunk(in[i])
	return true
}

func (m *merge) outputChunk(c chunk) {
	if len(c) > m.goal/2 {
		m.flushbuf()
		// pass entire chunk through to output
		m.out = append(m.out, c)
		m.size += len(c)
	} else {
		m.buf = append(m.buf, c...) // could exceed goal
	}
}

func (m *merge) flushbuf() {
	if len(m.buf) == 0 {
		return
	}
	m.out = append(m.out, slc.Clone(m.buf))
	m.size += len(m.buf)
	m.buf = m.buf[:0] // reuse buf
}

func (m *merge) result() *ixbuf {
	m.flushbuf()
	return &ixbuf{chunks: m.out, size: int32(m.size)}
}

func (c chunk) firstKey() string {
	return c[0].key
}

func (c chunk) lastKey() string {
	return c[len(c)-1].key
}

//-------------------------------------------------------------------

// Lookup searches for a key returns its associated offset or 0 if not found.
// Note: It may return offsets with the Delete or Update bits set.
func (ib *ixbuf) Lookup(key string) uint64 {
	if ib.size == 0 {
		return 0
	}
	_, c, i := ib.search(key)
	if i >= len(c) || c[i].key != key {
		return 0
	}
	return c[i].off
}

//-------------------------------------------------------------------

// Iter is used with btree.MergeAndSave
func (ib *ixbuf) Iter() iface.IterFn {
	if ib.size == 0 {
		return func() (string, uint64, bool) {
			return "", 0, false
		}
	}
	ti := 0
	c := ib.chunks[0]
	i := -1
	return func() (string, uint64, bool) {
		i++
		if i >= len(c) {
			if ti+1 >= len(ib.chunks) {
				return "", 0, false
			}
			ti++
			c = ib.chunks[ti]
			i = 0
		}
		slot := c[i]
		return slot.key, slot.off, true
	}
}

//-------------------------------------------------------------------

type Range = iface.Range

// Iterator is a Suneido style iterator for an ixbuf.
type Iterator struct {
	ib *ixbuf
	// rng is the Range of the iterator
	rng Range
	c   chunk
	// cur is the current key and offset.
	// We need to keep a copy of it because the ixbuf could change.
	cur slot
	// ci, i, and c point to the current slot = ib.chunks[ci][i]
	ci       int
	i        int
	modCount int32 // gets updated by Seek(All)
	state
}

var _ iface.Iter = (*Iterator)(nil)

type state byte

const (
	rewound state = iota
	within
	eof
)

func (ib *ixbuf) Iterator() iface.Iter {
	return &Iterator{ib: ib, modCount: ib.modCount, state: rewound,
		rng: iface.All}
}

func (it *Iterator) Range(rng Range) {
	it.rng = rng
	it.Rewind()
}

func (it *Iterator) Eof() bool {
	return it.ib.size == 0 || it.state == eof
}

func (it *Iterator) Modified() bool {
	return it.modCount != it.ib.modCount
}

func (it *Iterator) Cur() (string, uint64) {
	assert.That(it.state == within)
	return it.cur.key, it.cur.off
}

func (it *Iterator) Key() string {
	return it.cur.key
}

func (it *Iterator) Offset() uint64 {
	return it.cur.off
}

func (it *Iterator) HasCur() bool {
	return it.state == within
}

func (it *Iterator) Next() {
	if it.state == eof {
		return // stick at eof
	}
	if it.state == rewound {
		it.Seek(it.rng.Org)
		return
	}
	it.i++
	if it.i >= len(it.c) {
		if it.ci+1 >= len(it.ib.chunks) {
			it.state = eof
			return
		}
		it.ci++
		it.c = it.ib.chunks[it.ci]
		it.i = 0
	}
	it.cur = it.c[it.i]
	if it.cur.key >= it.rng.End {
		it.state = eof
	}
}

func (it *Iterator) Prev() {
	if it.state == eof {
		return // stick at eof
	}
	if it.state == rewound {
		it.SeekAll(it.rng.End)
		if it.Eof() || (it.rng.Org <= it.cur.key && it.cur.key < it.rng.End) {
			return
		}
		// Seek goes to >= so fallthrough to do previous
	}
	it.i--
	if it.i < 0 {
		if it.ci <= 0 {
			it.state = eof
			return
		}
		it.ci--
		it.c = it.ib.chunks[it.ci]
		it.i = len(it.c) - 1
	}
	it.cur = it.c[it.i]
	if it.cur.key < it.rng.Org || it.rng.End <= it.cur.key {
		it.state = eof
	}
}

func (it *Iterator) Rewind() {
	it.state = rewound
}

func (it *Iterator) Seek(key string) {
	it.SeekAll(key)
	if it.cur.key < it.rng.Org || it.rng.End <= it.cur.key {
		it.state = eof
	}
}

func (it *Iterator) SeekAll(key string) {
	if len(it.ib.chunks) == 0 {
		it.state = eof
		return
	}
	it.ci, it.c, it.i = it.ib.search(key)
	it.modCount = it.ib.modCount
	if it.i >= len(it.c) {
		it.i--
	}
	it.cur = it.c[it.i]
	it.state = within
}

//-------------------------------------------------------------------

func (ib *ixbuf) String() string {
	var sb strings.Builder
	sep := ""
	for _, c := range ib.chunks {
		for _, s := range c {
			sb.WriteString(sep)
			sep = " "
			sb.WriteString(s.String())
		}
	}
	return sb.String()
}

func (ib *ixbuf) Print() {
	fmt.Println("<<<------------------------")
	for i, c := range ib.chunks {
		if i > 0 {
			fmt.Println("+++")
		}
		c.print()
	}
	fmt.Println("------------------------>>>")
}

func (c chunk) print() {
	for _, s := range c {
		fmt.Print(" " + s.String())
	}
}

// OffString is for debugging
func OffString(off uint64) string {
	var s string
	switch off &^ Mask {
	case 0:
		s = "+"
	case Update:
		s = "="
	case Delete:
		s = "-"
	default:
		s += fmt.Sprintf("BAD BITS %b", off>>60)
	}
	return s + trace.Number(off&Mask)
}

func (s slot) String() string {
	return s.key + OffString(s.off)
}

// Check verifies that the keys are in order and there are no duplicates
func (ib *ixbuf) Check() {
	n := 0
	prev := ""
	for _, c := range ib.chunks {
		assert.That(len(c) > 0)
		for _, s := range c {
			switch cmp.Compare(prev, s.key) {
			case 0:
				panic("ixbuf: invalid duplicate key: " + prev)
			case 1:
				panic("ixbuf: out of order " + prev + " " + s.key)
			}
			prev = s.key
			n++
		}
	}
	assert.That(int(ib.size) == n)
}
