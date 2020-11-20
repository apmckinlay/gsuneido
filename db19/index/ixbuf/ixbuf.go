// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ixbuf defines an ordered list ixbuf.T
// with a mutating Insert and a immutable persistent Merge.
// It is designed for intermediate numbers of values, e.g. up to 16k or so.
// It is used mutably for per transaction index buffers
// and immutably for global index buffers.
package ixbuf

import (
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
)

type T = ixbuf

type ixbuf struct {
	chunks  []chunk
	size    int
}

type chunk []slot

type slot struct {
	key string
	off uint64
}

func (ib *ixbuf) Len() int {
	return ib.size
}

// goal is the desired chunk size for a given item count.
// It is chosen so the size of the chunk list is roughly the chunk size.
func goal(n int) int {
	// add 50% because average size is 2/3 full
	switch {
	case n < 256:
		return 16 + 16/2
	case n < 1024:
		return 32 + 32/2
	case n < 4096:
		return 64 + 64/2
	default:
		return 128 + 128/2
	}
}

// Insert adds an element. It mutates and is NOT thread-safe.
func (ib *ixbuf) Insert(key string, off uint64) {
	ib.size++
	if len(ib.chunks) == 0 {
		ib.chunks = make([]chunk, 1, 4)
		ib.chunks[0] = make([]slot, 1, 8)
		ib.chunks[0][0] = slot{key: key, off: off}
		return
	}
	ci := ib.search(key)
	// insert in place
	i := search(ib.chunks[ci], key)
	ib.chunks[ci] = append(ib.chunks[ci], slot{})
	c := ib.chunks[ci]
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

// search does a binary search the first key in each chunk.
// It returns len-1 if the key is greater than all keys.
func (ib *ixbuf) search(key string) int {
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
	return ints.Min(i, len(ib.chunks)-1)
}

// search does a binary search of one chunk for a key
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

//-------------------------------------------------------------------

func (ib *ixbuf) Delete(key string) bool {
	ci := ib.search(key)
	c := ib.chunks[ci]
	i := search(c, key)
	if i >= len(c) || c[i].key != key {
		return false
	}
	if len(c) == 1 {
		ib.chunks = append(ib.chunks[:ci], ib.chunks[ci+1:]...) // remove chunk
	} else {
		ib.chunks[ci] = append(c[:i], c[i+1:]...) // remove slot
	}
	ib.size--
	return true
}

//-------------------------------------------------------------------

// Merge combines several T into a new one.
// It does not modify its inputs so it is thread-safe
// as long as the inputs don't change.
// It is immutable persistent and the result may share chunks of the inputs
// so again, the inputs cannot change.
func Merge(ibs ...*ixbuf) *ixbuf {
	assert.That(len(ibs) > 1)
	in := make([][]chunk, 0, len(ibs))
	size := 0
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
	out := output{goal: goal(size), out: make([]chunk, 0, nc)}
	// merge at the chunk level
	for len(in) > 0 {
		// find minimum, inlined for speed
		i := 0
		key := in[0][0][0].key
		for j := 1; j < len(in); j++ {
			key2 := in[j][0][0].key
			if key2 < key {
				i = j
				key = key2
			}
		}
		out.add(in[i][0])
		if len(in[i]) == 1 {
			in = append(in[:i], in[i+1:]...) // remove empty input
		} else {
			in[i] = in[i][1:] // pop chunk
		}
	}
	return out.result()
}

type output struct {
	goal int
	in   []chunk
	buf  chunk // buf is between in and out
	out  []chunk
	size int
}

func (o *output) add(a chunk) {
	// fmt.Println("add chunk", chunkstr(a))
	if len(o.in) == 1 && o.in[0].lastKey() < a.firstKey() {
		o.separate(o.in[0])
		o.in[0] = a
		return
	}
	// There are other possible scenarios for separate (non-overlapping)
	// but they are too rare to slow down for.
	o.in = append(o.in, a)
	if len(o.in) == 1 { // first time
		return
	}
	o.mergeUpto(a[0])
}

func (c chunk) firstKey() string {
	return c[0].key
}

func (c chunk) lastKey() string {
	return c[len(c)-1].key
}

func (o *output) separate(c chunk) {
	// fmt.Println("separate", chunkstr(c))
	// fmt.Println("goal", o.goal, "len buf", len(o.buf), "len chunk", len(c))
	if len(c) > o.goal/2 {
		o.flushbuf()
		// fmt.Println("output chunk direct")
		o.out = append(o.out, c)
		o.size += len(c)
	} else {
		if len(o.buf) > o.goal {
			o.flushbuf()
		}
		// fmt.Println("append to buf")
		o.buf = append(o.buf, c...)
	}
}

func (o *output) flushbuf() {
	if len(o.buf) == 0 {
		return
	}
	// fmt.Println("flushbuf", chunkstr(o.buf))
	c := make([]slot, len(o.buf))
	copy(c, o.buf)
	o.out = append(o.out, c)
	o.size += len(o.buf)
	o.buf = o.buf[:0] // reuse buf
}

// mergeUpto merges at the individual key level.
// After merging, the inputs may not be in order.
// The last (limit) will always end up the minimum
// because we merge everything less.
func (o *output) mergeUpto(limit slot) {
	// fmt.Println("mergeUpto", limit)
	// limit off == 0 means no limit, merge everything (used for flush)
	for len(o.in) > 1 {
		// // fmt.Println("buf", chunkstr(o.buf))
		// find minimum, inlined for speed (~ 5%)
		i := 0
		key := o.in[0][0].key
		for j := 1; j < len(o.in); j++ {
			key2 := o.in[j][0].key
			if key2 < key {
				i = j
				key = key2
			}
		}
		slot := o.in[i][0]
		// // fmt.Println("next", slot.key)
		if limit.off != 0 && slot.key >= limit.key {
			// // fmt.Println("limit")
			return
		}
		// assert.That(len(o.buf) == 0 || slot.key > o.buf.lastKey())
		o.buf = append(o.buf, slot)
		if len(o.buf) > o.goal {
			o.flushbuf()
		}
		if len(o.in[i]) == 1 {
			o.in = append(o.in[:i], o.in[i+1:]...) // remove empty input
		} else {
			o.in[i] = o.in[i][1:] // advance to next slot
		}
	}
}

func (o *output) result() *ixbuf {
	if len(o.in) >= 1 {
		if len(o.in) > 1 {
			o.mergeUpto(slot{})
		}
		if len(o.buf) == 0 {
			o.out = append(o.out, o.in[0])
			o.size += len(o.in[0])
		} else {
			o.separate(o.in[0])
		}
	}
	if len(o.buf) > 0 {
		o.flushbuf()
	}
	return &ixbuf{chunks: o.out, size: o.size}
}

//-------------------------------------------------------------------

type Iter = func() (key string, off uint64, ok bool)

func (ib *ixbuf) Iter(bool) Iter {
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

type Visitor func(key string, off uint64)

func (ib *ixbuf) ForEach(fn Visitor) {
	for _, c := range ib.chunks {
		for _, slot := range c {
			fn(slot.key, slot.off)
		}
	}
}
