// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bytes"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/str"
)

// fNode is a file based btree node with partial incremental encoding.
// Nodes are variable length and are packed into a sequence of bytes
// with variable length entries.
// So we can only iterate from the beginning, no random access or binary search.
//
// Entry is:
//		- 5 byte smalloffset
//		- one byte prefix length
//		- one byte key part length
//		- key part bytes (variable length)
type fNode []byte

func fAppend(fn fNode, offset uint64, npre int, diff string) fNode {
	fn = stor.AppendSmallOffset(fn, offset)
	fn = append(fn, byte(npre), byte(len(diff)))
	fn = append(fn, diff...)
	return fn
}

func fRead(fn fNode) (npre int, diff []byte, offset uint64) {
	offset = stor.ReadSmallOffset(fn)
	npre = int(fn[5])
	dn := int(fn[6])
	diff = fn[7 : 7+dn]
	return
}

func fLen(diff []byte) int {
	return 5 + 1 + 1 + len(diff)
}

func (fn fNode) next(i int) int {
	return i + 7 + int(fn[i+6])
}

type fData struct {
	key    string
	offset uint64
}

type fNodeBuilder struct {
	fe       fNode
	notFirst bool
	fi       int
	prev     string
	known    string
	offset   uint64
	fi2      int
	known2   string
	offset2  uint64
}

func (fb *fNodeBuilder) Add(key string, offset uint64, embedLen int) {
	if fb.notFirst {
		if key <= fb.prev {
			panic("fBuilder keys must be inserted in order, without duplicates")
		}
	} else {
		fb.notFirst = true
	}
	if len(fb.fe) == 0 {
		fb.fe = fAppend(fb.fe, offset, 0, "")
		fb.known = ""
	} else {
		fb.fi2 = fb.fi
		fb.fi = len(fb.fe)
		npre, diff, known := addone(key, fb.prev, fb.known, embedLen)
		fb.fe = fAppend(fb.fe, offset, npre, diff)
		fb.known2 = fb.known
		fb.known = known
		fb.offset2 = fb.offset
		fb.offset = offset
	}
	fb.prev = key
}

// Split saves all but the last two entries as the left node
// and initializes fb.fe with the last two entries
func (fb *fNodeBuilder) Split(store *stor.Stor) (leftOff uint64, splitKey string) {
	splitKey = fb.known2 // known of second last entry
	left := fb.fe[:fb.fi2]
	leftOff = left.putNode(store)
	// first entry becomes 0, ""
	right := fAppend(fb.fe[:0], fb.offset2, 0, "") // offset of second last entry
	// second entry becomes 0, known
	right = fAppend(right, fb.offset, 0, fb.known) // offset,known of last entry
	fb.fe = right
	return
}

func (fb *fNodeBuilder) Entries() fNode {
	return fb.fe
}

func addone(key, prev, known string, embedLen int) (npre int, diff string, knownNew string) {
	if key <= prev {
		print("OUT OF ORDER: prev", prev, "key", key)
	}
	assert.That(key > prev)
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
		diff = key[len(known):ints.Min(npre+embedLen, len(key))]
	}
	assert.That(len(diff) > 0)
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

func commonSlicePrefixLen(s, t []byte) int {
	for i := 0; ; i++ {
		if i >= len(s) || i >= len(t) || s[i] != t[i] {
			return i
		}
	}
}

// search returns the offset and range
// of the entry that could match the search string
func (fn fNode) search(s string) (uint64, string, string) {
	var off uint64
	var known []byte
	it := fn.iter()
	for it.next() && s >= string(it.known) {
		off = it.offset
		known = append(known[:0], it.known...)
	}
	return off, string(known), string(it.known)
}

func (fn fNode) contains(s string, get func(uint64) string) bool {
	if len(fn) == 0 {
		return false
	}

	offset, _, _ := fn.search(s)
	return s == get(offset)
}

const (
	insMiddle = iota
	insStart
	insEnd
)

// insert adds a new key to a node. get will be nil for tree nodes.
func (fn fNode) insert(keyNew string, offNew uint64, get func(uint64) string) (fNode, int) {
	where := insMiddle
	if len(fn) == 0 {
		return fAppend(fn, offNew, 0, ""), where
	}
	// search
	curFi := 0
	curNpre := 0
	curEof := false
	var curOffset uint64
	var curDiff, curKnown []byte
	it := fn.iter()
	for it.next() && keyNew >= string(it.known) {
		curFi = it.fi
		curNpre = it.npre
		curEof = it.eof()
		curOffset = it.offset
		curDiff = append(curDiff[:0], it.diff...)
		curKnown = append(curKnown[:0], it.known...)
	}

	curoff := curOffset
	curkey := string(curKnown)
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
		if curEof {
			// at end
			npre, diff, _ = addone(keyNew, curkey, string(curKnown), embedLen)
			return fAppend(fn, offNew, npre, diff), insEnd
		}
		npre, diff, knownNew = addone(keyNew, curkey, string(curKnown), embedLen)
		ins = fAppend(ins, offNew, npre, diff)
		i = it.fi
		j = it.fi
		prev = knownNew
	} else { // newkey before curkey
		if curFi == 0 {
			where = insStart
		}
		// first entry stays the same, just update offset
		ins = fAppend(ins, offNew, curNpre, string(curDiff))
		// old first key becomes second entry
		npre, diff, knownNew = addone(curkey, keyNew, string(curKnown), embedLen)
		ins = fAppend(ins, curoff, npre, diff)
		i = curFi
		j = it.fi
		prev = curkey
	}
	if !curEof {
		npre2, diff2, _ := addone(string(it.known), prev, knownNew, embedLen)
		if npre2 != it.npre || diff2 != string(it.diff) {
			// adjust following entry
			ins = fAppend(ins, it.offset, npre2, diff2)
			j += fLen(it.diff)
		}
	}
	fn = fn.replace(i, j, ins)
	return fn, where
}

// replace is used by insert and delete
// to replace a portion of a node (i,j) with new content (rep)
func (fn fNode) replace(i, j int, rep fNode) fNode {
	nr := len(rep)
	d := nr - (j - i)
	fn = bytes.Grow(fn, d)
	copy(fn[i+nr:], fn[j:])
	copy(fn[i:], rep)
	if d < 0 {
		fn = fn[:len(fn)+d]
	}
	return fn
}

func (fn fNode) delete(offset uint64) (fNode, bool) {
	// search
	var prevKnown []byte
	it := fn.iter()
	for {
		if !it.next() {
			return nil, false // not found
		}
		if stor.EqualSmallOffset(fn[it.fi:], offset) {
			break
		}
		prevKnown = append(prevKnown[:0], it.known...)
	}
	i := it.fi

	j := fn.next(i)
	if j >= len(fn) {
		// delete last item, simplest case, no adjustments
		return fn[:i], true
	}

	rep := make(fNode, 0, 64)
	if i == 0 {
		// deleting first entry so make following into first
		rep = rep.updateCopy(fn, j, 0, "")
		j = fn.next(j)
		prevKnown = it.known
		it.next()
		// then adjust following entry if there is one
	}
	if it.next() {
		npre := commonSlicePrefixLen(prevKnown, it.known)
		diff := it.known[npre:]
		rep = rep.updateCopy(fn, j, npre, string(diff))
		j = fn.next(j)
	}
	fn = fn.replace(i, j, rep)
	return fn, true
}

func (fn fNode) updateCopy(src fNode, i int, npre int, diff string) fNode {
	fn = append(fn, src[i:i+5]...) // offset
	fn = append(fn, byte(npre), byte(len(diff)))
	fn = append(fn, diff...)
	return fn
}

func (fn fNode) setOffset(fi int, off uint64) {
	stor.WriteSmallOffset(fn[fi:], off)
}

// iter -------------------------------------------------------------

type fnIter struct {
	fn     fNode
	fi     int // position in original fEntries
	npre   int
	diff   []byte
	known  []byte
	offset uint64
}

func (fn fNode) iter() *fnIter {
	return &fnIter{fn: fn, fi: -7}
}

func (it *fnIter) next() bool {
	it.fi += fLen(it.diff)
	if it.fi >= len(it.fn) {
		it.known = it.known[:0] // ""
		return false
	}
	it.npre, it.diff, it.offset = fRead(it.fn[it.fi:])

	// maybe remove this validation in production?
	// if it.known == "" && it.npre == 0 && it.diff == "" {
	// 	// first
	// } else if it.npre <= len(it.known) {
	// 	if len(it.diff) < 1 {
	// 		// print("bad diff len, npre", it.npre, "diff", it.diff, "known", it.known)
	// 		panic("bad diff len")
	// 	}
	// } else {
	// 	if len(it.diff) != it.npre-len(it.known)+1 {
	// 		// print("bad diff len, npre", it.npre, "diff", it.diff, "known", it.known)
	// 		panic("bad diff len")
	// 	}
	// }

	if it.npre <= len(it.known) {
		it.known = append(it.known[:it.npre], it.diff...)
	} else {
		it.known = append(it.known, it.diff...)
	}
	return true
}

func (it *fnIter) eof() bool {
	return it.fi+fLen(it.diff) >= len(it.fn)
}

//-------------------------------------------------------------------

func (fn fNode) stats() {
	n := fn.check()
	avg := float32(len(fn)-7*n) / float32(n)
	print("    n", n, "len", len(fn), "avg", avg)
}

func (fn fNode) checkData(data []string, get func(uint64) string) {
	fn.checkUpTo(len(data)-1, data, get)
}

// checkUpTo is used during inserting.
// It checks that inserted keys are present
// and uninserted keys are not present.
func (fn fNode) checkUpTo(i int, data []string, get func(uint64) string) {
	n := fn.check()
	nn := 0
	for j, d := range data {
		if (d != "" && j <= i) != fn.contains(d, get) {
			panic("can't find " + d)
		}
		if d != "" && j <= i {
			nn++
		}
	}
	if nn != n {
		panic("check count expected " + strconv.Itoa(n) +
			" got " + strconv.Itoa(nn))
	}
}

func (fn fNode) check() int {
	n := 0
	var prev []byte
	it := fn.iter()
	for it.next() {
		if string(it.known) < string(prev) {
			print("known", it.known, "prev", prev)
			panic("fEntries out of order")
		}
		if it.fi > 7 && it.npre > len(prev)+(len(it.diff)-1) {
			panic("npre > len(prev.known)")
		}
		prev = append(prev[:0], it.known...)
		n++
	}
	return n
}

func (fn fNode) print() {
	it := fn.iter()
	for it.next() {
		print(it.offset, it.known)
	}
}

func (fn fNode) printLeafNode(get func(uint64) string) {
	it := fn.iter()
	for it.next() {
		offset := it.offset
		print(strconv.Itoa(it.fi)+": {", offset, it.npre, it.diff, "}",
			it.known, "("+get(offset)+")")
	}
}

func (fn fNode) printTreeNode() {
	it := fn.iter()
	for it.next() {
		offset := it.offset
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
		case []byte:
			args[i] = string(x)
			if len(x) == 0 {
				args[i] = "''"
			}
		}
	}
	fmt.Println(args...)
}

func (fn fNode) String() string {
	s := "["
	it := fn.iter()
	for it.next() {
		known := string(it.known)
		if known == "" {
			known = "''"
		}
		s += fmt.Sprint(known, "=", it.offset) + " "
	}
	return strings.TrimSpace(s) + "]"
}

func (fn fNode) knowns() string {
	var sb strings.Builder
	it := fn.iter()
	for it.next() {
		sb.Write(it.known)
		sb.WriteByte(' ')
	}
	return sb.String()
}
