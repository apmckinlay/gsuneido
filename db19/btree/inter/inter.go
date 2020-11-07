// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package inter defines an ordered list inter.T
// with a mutating Insert and a immutable persistent Merge.
// It is designed for intermediate numbers of values, e.g. up to 16k or so.
// It is used mutably for per transaction index buffers
// and immutably for global index buffers.
package inter

import (
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
)

type T struct {
	chunks  []chunk
	size    int
	TranNum int //TODO temporary
}

type chunk []slot

type slot struct {
	key string
	off uint64
}

func (t *T) Len() int {
	return t.size
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
func (t *T) Insert(key string, off uint64) {
	t.size++
	if len(t.chunks) == 0 {
		t.chunks = make([]chunk, 1, 4)
		t.chunks[0] = make([]slot, 1, 8)
		t.chunks[0][0] = slot{key: key, off: off}
		return
	}
	ci := t.search(key)
	// insert in place
	i := search(t.chunks[ci], key)
	t.chunks[ci] = append(t.chunks[ci], slot{})
	c := t.chunks[ci]
	copy(c[i+1:], c[i:])
	c[i] = slot{key: key, off: off}

	if len(c) > goal(t.size) {
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
		t.chunks[ci] = left
		t.chunks = append(t.chunks, nil)
		ci++
		copy(t.chunks[ci+1:], t.chunks[ci:])
		t.chunks[ci] = right
	}
}

// search does a binary search the first key in each chunk.
// It returns len-1 if the key is greater than all keys.
func (t *T) search(key string) int {
	i, j := 0, len(t.chunks)
	for i < j {
		h := int(uint(i+j) >> 1) // i ≤ h < j
		c := t.chunks[h]
		if key > c.lastKey() {
			i = h + 1
		} else {
			j = h
		}
	}
	return ints.Min(i, len(t.chunks)-1)
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

func (t *T) Delete(key string) bool {
	ci := t.search(key)
	c := t.chunks[ci]
	i := search(c, key)
	if i >= len(c) || c[i].key != key {
		return false
	}
	if len(c) == 1 {
		t.chunks = append(t.chunks[:ci], t.chunks[ci+1:]...) // remove chunk
	} else {
		t.chunks[ci] = append(c[:i], c[i+1:]...) // remove slot
	}
	t.size--
	return true
}

//-------------------------------------------------------------------

// Merge combines several T into a new one.
// It does not modify its inputs so it is thread-safe
// as long as the inputs don't change.
// It is immutable persistent and the result may share chunks of the inputs
// so again, the inputs cannot change.
func Merge(ts ...*T) *T {
	assert.That(len(ts) > 1)
	in := make([][]chunk, 0, len(ts))
	size := 0
	nc := 0
	var single *T
	for _, t := range ts {
		if t.size == 0 {
			continue
		}
		in = append(in, t.chunks)
		size += t.size
		nc += len(t.chunks)
		single = t
	}
	if len(in) == 0 {
		return &T{}
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

func (o *output) result() *T {
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
	return &T{chunks: o.out, size: o.size}
}

//-------------------------------------------------------------------

type Iter = func() (string, uint64, bool)

func (t *T) Iter(bool) Iter {
	if t.size == 0 {
		return func() (string, uint64, bool) {
				return "", 0, false
		}
	}
	ti := 0
	c := t.chunks[0]
	i := -1
	return func() (string, uint64, bool) {
		i++
		if i >= len(c) {
			if ti+1 >= len(t.chunks) {
				return "", 0, false
			}
			ti++
			c = t.chunks[ti]
			i = 0
		}
		slot := c[i]
		return slot.key, slot.off, true
	}
}

type Visitor func(key string, off uint64)

func (t *T) ForEach(fn Visitor) {
	for _, c := range t.chunks {
		for _, slot := range c {
			fn(slot.key, slot.off)
		}
	}
}
