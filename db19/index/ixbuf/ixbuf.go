// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ixbuf defines an ordered list ixbuf.T
// with a mutating Insert and a immutable persistent Merge.
// It is designed for intermediate numbers of values, e.g. up to 16k or so.
// It is used mutably for per transaction index buffers
// and immutably for global index buffers.
package ixbuf

import (
	"fmt"
	"log"

	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"golang.org/x/exp/slices"
)

type T = ixbuf

type ixbuf struct {
	chunks []chunk
	size   int32
	// modCount is used by Iterator to detect modifications.
	// No locking since ixbuf is thread contained when mutable.
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
func (ib *ixbuf) Insert(key string, off uint64) {
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
		slot.off = Combine(slot.off, off) // handles update/delete
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
	return ord.Min(i, len(ib.chunks)-1)
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
func (ib *ixbuf) Update(key string, off uint64) {
	ib.Insert(key, off|Update)
}

// Delete combines if the key exists, otherwise it adds a delete tombstone.
func (ib *ixbuf) Delete(key string, off uint64) {
	ib.Insert(key, off|Delete)
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
	goal int
	in   [][]chunk
	buf  chunk
	out  []chunk
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

const Update = 1 << 62
const Delete = 1 << 63

// OffString is for debugging
func OffString(off uint64) string {
	s := fmt.Sprintf("%d", off&0xffffffffff)
	off &^= 0xffffffffff
	if off&Update != 0 {
		s += " update"
		off &^= Update
	}
	if off&Delete != 0 {
		s += " delete"
		off &^= Delete
	}
	if off != 0 {
		s += fmt.Sprintf("BAD BITS %x", off>>40)
	}
	return s
}

const (
	add_update    = 0b_00_01 //lint:ignore ST1003 for clarity
	add_delete    = 0b_00_10 //lint:ignore ST1003 for clarity
	update_update = 0b_01_01 //lint:ignore ST1003 for clarity
	update_delete = 0b_01_10 //lint:ignore ST1003 for clarity
	delete_add    = 0b_10_00 //lint:ignore ST1003 for clarity
)

func Combine(off1, off2 uint64) uint64 {
	ops := off1>>60 | off2>>62
	switch ops {
	case add_update:
		return off2 &^ Update // => add
	case add_delete:
		return 0 // => should be removed
	case update_update, update_delete:
		// update_delete needs to keep delete (not return 0)
		// because the add is in another layer
		return off2
	case delete_add:
		return off2 | Update
	default:
		log.Printf("ixbuf invalid Combine %b\n", ops)
		panic("ixbuf invalid Combine")
	}
}

func (m *merge) outputSlot(s2 slot) {
	last := len(m.buf) - 1
	if last >= 0 {
		s1 := &m.buf[last]
		if s2.key == s1.key {
			s1.off = Combine(s1.off, s2.off)
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
	for j := 0; j < len(m.in); j++ {
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
	m.out = append(m.out, slices.Clone(m.buf))
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

func (ib *ixbuf) PrefixLookup(key string) uint64 {
	if ib.size == 0 {
		return 0
	}
	_, c, i := ib.search(key)
	if i >= len(c) || !ixkey.HasPrefix(c[i].key, key) {
		return 0
	}
	return c[i].off
}

//-------------------------------------------------------------------

type Iter = func() (key string, off uint64, ok bool)

// Iter is used with btree.MergeAndSave
func (ib *ixbuf) Iter() Iter {
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

type Range = iterator.Range

// Iterator is a Suneido style iterator for an ixbuf.
type Iterator struct {
	ib       *ixbuf
	modCount int32
	state
	// ci, i, and c point to the current slot = ib.chunks[ci][i]
	ci int
	i  int
	c  chunk
	// cur is the current key and offset.
	// We need to keep a copy of it because the ixbuf could change.
	cur slot
	// rng is the Range of the iterator
	rng Range
}

var _ iterator.T = (*Iterator)(nil)

type state byte

const (
	rewound state = iota
	within
	eof
)

func (ib *ixbuf) Iterator() *Iterator {
	return &Iterator{ib: ib, modCount: ib.modCount, state: rewound,
		rng: iterator.All}
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
	assert.Msg("Cur when not within").That(it.state == within)
	return it.cur.key, it.cur.off
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
		off := s.off
		d := ""
		if (off & Delete) == Delete {
			off &^= Delete
			d = "d"
		}
		if (off & Update) == Update {
			off &^= Update
			d += "u"
		}
		fmt.Printf("%q %d %s\n", s.key, off, d)
	}
}
