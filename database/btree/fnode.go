// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"strconv"

	"github.com/apmckinlay/gsuneido/database/stor"
	"github.com/apmckinlay/gsuneido/util/bytes"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/verify"
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

func fRead(fe_ fNode) (fe fNode, npre int, diff string) {
	fe = fe_[stor.SmallOffsetLen:]
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
		npre, diff, known := addone(key, fb.prev, fb.known, 1)
		fb.fe = fAppend(fb.fe, offset, npre, diff)
		fb.known = known
	}
	fb.prev = key
}

func (fb *fNodeBuilder) Entries() fNode {
	return fb.fe
}

func addone(key, prev, known string, embedLen int) (npre int, diff string, knownNew string) {
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
	var ofi int
	var known string
	it := fn.Iter()
	for it.next() && s >= it.known {
		ofi = it.fi
		known = it.known
	}
	return fn.offset(ofi), known, it.known
}

func (fn fNode) contains(s string, get func(uint64) string) bool {
	offset, _, _ := fn.search(s)
	return s == get(offset)
}

// insert adds a new key to a node. get will be nil for tree nodes.
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

	curoff := fn.offset(cur.fi)
	curkey := cur.known
	embedLen := 255
	if get != nil {
		embedLen = 1
		curkey = get(curoff)
	}
	var prev string
	ins := make(fNode, 0, 64)
	var npre int
	var diff string
	var knownNew string
	var i, j int
	if keyNew > curkey { // newkey after curkey
		if cur.eof() {
			npre, diff, _ = addone(keyNew, curkey, cur.known, embedLen)
			return fAppend(fn, offNew, npre, diff)
		}
		npre, diff, knownNew = addone(keyNew, curkey, cur.known, embedLen)
		ins = fAppend(ins, offNew, npre, diff)
		i = it.fi
		j = it.fi
		prev = knownNew
	} else { // newkey before curkey
		// first entry stays the same, just update offset
		ins = fAppend(ins, offNew, cur.npre, cur.diff)
		// old first key becomes second entry
		npre, diff, knownNew = addone(curkey, keyNew, cur.known, embedLen)
		ins = fAppend(ins, curoff, npre, diff)
		i = cur.fi
		j = it.fi
		prev = curkey
	}
	if !cur.eof() {
		npre2, diff2, _ := addone(it.known, prev, knownNew, embedLen)
		if npre2 != it.npre || diff2 != it.diff {
			// adjust following entry
			ins = fAppend(ins, fn.offset(it.fi), npre2, diff2)
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

func (fn fNode) split(fe fNode, fi int, newkey string, newoff uint64) {

}

func (fn fNode) offset(fi int) uint64 {
	_, offset := stor.ReadSmallOffset(fn[fi:])
	return offset
}

// iter -------------------------------------------------------------

type iter struct {
	fn         fNode
	fi         int // position in original fEntries
	npre       int
	diff       string
	known      string
}

func (fn fNode) Iter() *iter {
	return &iter{fn: fn, fi: -7}
}

func (it *iter) next() bool {
	it.fi += fLen(it.diff)
	if len(it.fn) == 0 {
		it.known = ""
		return false
	}
	it.fn, it.npre, it.diff = fRead(it.fn)
	if it.known == "" && it.npre == 0 && it.diff == "" {
		// first
	} else if it.npre <= len(it.known) {
		if len(it.diff) < 1 {
			panic("unexpected diff len " + it.diff)
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

//-------------------------------------------------------------------

func (fn fNode) stats() {
	n := fn.check()
	avg := float32(len(fn)-7*n) / float32(n)
	print("    n", n, "len", len(fn), "avg", avg)
}

func (fn fNode) checkData(data []string, get func(uint64) string) {
	n := len(data)
	fn.checkUpTo(n-1, data, get)
}

// checkUpTo is used during inserting.
// It checks that inserted keys are present
// and uninserted keys are not present.
func (fn fNode) checkUpTo(i int, data []string, get func(uint64) string) {
	verify.That(fn.check() == i+1)
	for j, d := range data {
		if j <= i != fn.contains(d, get) {
			panic("can't find " + d)
		}
	}
}

func (fn fNode) check() int {
	n := 0
	prev := ""
	it := fn.Iter()
	for it.next() {
		if it.known < prev {
			panic("fEntries out of order")
		}
		prev = it.known
		n++
	}
	return n
}

func (fn fNode) print() {
	it := fn.Iter()
	for it.next() {
		print(fn.offset(it.fi), it.known)
	}
}

func (fn fNode) printLeafNode(get func(uint64) string) {
	it := fn.Iter()
	for it.next() {
		offset := fn.offset(it.fi)
		print(strconv.Itoa(it.fi)+": {", offset, it.npre, it.diff, "}",
			it.known, "("+get(offset)+")")
	}
}

func (fn fNode) printTreeNode() {
	it := fn.Iter()
	for it.next() {
		offset := fn.offset(it.fi)
		print(strconv.Itoa(it.fi)+": {", offset, it.npre, it.diff, "}",
			it.known)
	}
}

func print(args ...interface{}) {
	for i, x := range args {
		switch x := x.(type) {
		case string:
			if x == "" {
				args[i] = "'" + x + "'"
			}
		}
	}
	fmt.Println(args...)
}
