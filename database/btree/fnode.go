// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/database/stor"
	"github.com/apmckinlay/gsuneido/util/bytes"
	"github.com/apmckinlay/gsuneido/util/str"
)

// fNode is a file based btree node with partial incremental encoding.
// Nodes are variable length and are packed into a sequence of bytes
// with variable length entries.
// So we can only iterate from the beginning. No random access.
//
// Entry is:
//		- 5 byte smalloffset
//		- one byte prefix length
//		- one byte key part length
//		- key part bytes (variable length)
type fNode []byte

func fAppend(fe fNode, offset uint64, npre int, diff string) fNode {
	fe = stor.AppendSmallOffset(fe, offset)
	fe = append(fe, byte(npre), byte(len(diff)))
	fe = append(fe, diff...)
	return fe
}

func fRead(fe_ fNode) (fe fNode, offset uint64, npre int, diff string) {
	fe, offset = stor.ReadSmallOffset(fe_)
	npre = int(fe[0])
	sn := int(fe[1])
	fe = fe[2:]
	diff = string(fe[:sn])
	fe = fe[sn:]
	return
}

func fLen(diff string) int {
	return 5 + 1 + 1 + len(diff)
}

type fData struct {
	key    string
	offset uint64
}

// embedLen only needs to be one.
// Making it larger reduces looking at data for failing searches.
const embedLen = 1

type fNodeBuilder struct {
	fe    fNode
	prev  string
	known string
}

func (fb *fNodeBuilder) Add(key string, offset uint64) {
	if key <= fb.prev {
		panic("fBuilder keys must be inserted in order, without duplicates")
	}
	if len(fb.fe) == 0 {
		fb.fe = fAppend(fb.fe, offset, 0, "")
		fb.known = ""
	} else {
		npre, diff, known := addone(key, fb.prev, fb.known)
		fb.fe = fAppend(fb.fe, offset, npre, diff)
		fb.known = known
	}
	fb.prev = key
}

func (fb *fNodeBuilder) Entries() fNode {
	return fb.fe
}

func addone(key, prev, known string) (npre int, diff string, knownNew string) {
	npre = commonPrefixLen(prev, key)
	if npre > 255 {
		panic("key common prefix too long")
	}
	if npre <= len(known) {
		// normal case
		diff = str.Subn(key, npre, embedLen)
	} else {
		// prefix is longer than what's known
		// so we have to embed the missing info + embedLen
		diff = key[len(known) : npre+embedLen]
	}
	knownNew = str.Subn(key, 0, npre+embedLen)
	return
}

func commonPrefixLen(s, t string) int {
	for i := 0; ; i++ {
		if i >= len(s) || i >= len(t) || s[i] != t[i] {
			return i
		}
	}
}

// search returns the offset and range
// of the entry that could match the search string
func (fn fNode) search(s string) (uint64, string, string) {
	var offset uint64
	var known string
	it := fn.Iter()
	for it.next() && s >= it.known {
		offset = it.offset
		known = it.known
	}
	return offset, known, it.known
}

func (fn fNode) contains(s string, get func(uint64) string) bool {
	offset, _, _ := fn.search(s)
	return s == get(offset)
}

func (fn fNode) insert(keyNew string, offNew uint64, get func(uint64) string) fNode {
	if len(fn) == 0 {
		return fAppend(fn, offNew, 0, "")
	}
	// search
	var cur iter
	it := fn.Iter()
	for it.next() && keyNew >= it.known {
		cur = *it
	}

	curkey := get(cur.offset)
	var prev string
	ins := make(fNode, 0, 64)
	var npre int
	var diff string
	var knownNew string
	var i, j int
	if keyNew > curkey { // newkey after curkey
		if cur.eof() {
			npre, diff, _ = addone(keyNew, curkey, cur.known)
			return fAppend(fn, offNew, npre, diff)
		}
		npre, diff, knownNew = addone(keyNew, curkey, cur.known)
		ins = fAppend(ins, offNew, npre, diff)
		i = it.fi
		j = it.fi
		prev = knownNew
	} else { // newkey before curkey
		// first entry stays the same, just update offset
		ins = fAppend(ins, offNew, cur.npre, cur.diff)
		// old first key becomes second entry
		npre, diff, knownNew = addone(curkey, keyNew, cur.known)
		ins = fAppend(ins, cur.offset, npre, diff)
		i = cur.fi
		j = it.fi
		prev = curkey
	}
	if !cur.eof() {
		npre2, diff2, _ := addone(it.known, prev, knownNew)
		if npre2 != it.npre || diff2 != it.diff {
			// adjust following entry
			ins = fAppend(ins, it.offset, npre2, diff2)
			j += fLen(it.diff)
		}
	}
	fn = replace(fn, i, j, ins)
	return fn
}

func replace(fe fNode, i, j int, ins fNode) fNode {
	d := len(ins) - (j - i)
	fe = bytes.Grow(fe, d)
	copy(fe[i+d:], fe[i:])
	copy(fe[i:], ins)
	return fe
}

// iter -------------------------------------------------------------

type iter struct {
	fn         fNode
	fi         int // position in original fEntries
	offset     uint64
	npre       int
	diff       string
	known      string
	afterFirst bool
}

func (fn fNode) Iter() *iter {
	return &iter{fn: fn}
}

func (it *iter) next() bool {
	if it.afterFirst {
		it.fi += fLen(it.diff)
	} else {
		it.afterFirst = true
	}
	if len(it.fn) == 0 {
		it.known = ""
		return false
	}
	//TODO don't decode offset unless needed
	it.fn, it.offset, it.npre, it.diff = fRead(it.fn)
	if it.known == "" && it.npre == 0 && it.diff == "" {
		// first
	} else if it.npre <= len(it.known) {
		if len(it.diff) != 1 {
			panic("unexpected diff len")
		}
	} else {
		if len(it.diff) != it.npre-len(it.known)+1 {
			panic("unexpected diff len")
		}
	}
	//TODO use a buffer for known to reduce allocation
	if it.npre <= len(it.known) {
		it.known = it.known[:it.npre] + it.diff
	} else {
		it.known = it.known + it.diff
	}
	return true
}

func (it *iter) eof() bool {
	return len(it.fn) == 0
}
