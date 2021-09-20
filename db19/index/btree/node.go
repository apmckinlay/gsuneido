// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bytes"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/str"
)

// node is a file based btree node with partial incremental encoding.
// Nodes are variable length and are packed into a sequence of bytes
// with variable length entries.
// So we can only iterate from the beginning, no random access or binary search.
//
// Entry is:
//		- 5 byte smalloffset
//		- one byte prefix length (npre)
//		- one byte key part length (len diff)
//		- key part bytes (variable length) (diff)
type node []byte

const embedAll = 255

func (nd node) append(offset uint64, npre int, diff string) node {
	nd = stor.AppendSmallOffset(nd, offset)
	nd = append(nd, byte(npre), byte(len(diff)))
	nd = append(nd, diff...)
	return nd
}

func (nd node) read() (npre int, diff []byte, offset uint64) {
	offset = stor.ReadSmallOffset(nd)
	npre = int(nd[5])
	dn := int(nd[6])
	diff = nd[7 : 7+dn]
	return
}

func fLen(diff []byte) int {
	return 5 + 1 + 1 + len(diff)
}

func (nd node) next(i int) int {
	return i + 7 + int(nd[i+6])
}

// addone calculates the encoding for a new entry.
//
// NOTE: if key is a known (not a full value) then embedLen should be embedAll
func addone(key, prev, known string, embedLen int) (npre int, diff string, knownNew string) {
	if key <= prev {
		fmt.Printf("OUT OF ORDER: prev %q key %q\n", prev, key)
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

// search returns the offset of the entry that could match the key
func (nd node) search(key string) uint64 {
	var off uint64
	it := nd.iter()
	for it.next() && key >= string(it.known) {
		off = it.offset
	}
	return off
}

// update adds, updates, or deletes a key in a node.
// get will be nil for tree nodes.
// Used by merge.
func (nd node) update(keyNew string, offNew uint64, get func(uint64) string) node {
	if len(nd) == 0 {
		return nd.append(offNew, 0, "")
	}
	// search
	curPos := 0
	curNpre := 0
	curEof := false
	var curOffset uint64
	var curDiff, curKnown []byte
	it := nd.iter()
	for it.next() && keyNew >= string(it.known) {
		//TODO switch to nodeIter.copyFrom
		curPos = it.pos
		curNpre = it.npre
		curEof = it.eof()
		curOffset = it.offset
		curDiff = it.diff
		curKnown = append(curKnown[:0], it.known...) // copy over
	}

	curoff := curOffset
	curkey := string(curKnown)
	embedLen := embedAll
	if get != nil {
		embedLen = 1
		curkey = get(curoff)
	}

	if offNew>>62 != 0 {
		if keyNew == curkey {
			if offNew&ixbuf.Delete != 0 {
				_ = t && trace("before delete", nd.knowns())
				nd, _ = nd.delete(curOffset)
				_ = t && trace("after delete", nd.knowns())
			} else { // update
				nd.setOffset(curPos, offNew)
			}
			return nd
		}
		panic("update/delete on nonexistent")
	}

	var prev string
	ins := make(node, 0, 64)
	var npre int
	var diff string
	var knownNew string
	var i, j int
	if keyNew > curkey { // newkey after curkey
		if curEof {
			// at end
			npre, diff, _ = addone(keyNew, curkey, string(curKnown), embedLen)
			return nd.append(offNew, npre, diff)
		}
		npre, diff, knownNew = addone(keyNew, curkey, string(curKnown), embedLen)
		// print("after:", "key", keyNew, "prev", curkey, "known", curKnown,
		// 	"=>", "npre", npre, "diff", diff, "knownNew", knownNew)
		ins = ins.append(offNew, npre, diff)
		i = it.pos
		j = it.pos
		prev = knownNew
	} else { // newkey before curkey
		// first entry stays the same, just update offset
		ins = ins.append(offNew, curNpre, string(curDiff))
		// old first key becomes second entry
		npre, diff, knownNew = addone(curkey, keyNew, string(curKnown), embedLen)
		// print("before:", "key", curkey, "prev", keyNew, "known", curKnown,
		// 	"=>", "npre", npre, "diff", diff, "knownNew", knownNew)
		ins = ins.append(curoff, npre, diff)
		i = curPos
		j = it.pos
		prev = curkey
	}
	if !curEof {
		npre2, diff2, _ := addone(string(it.known), prev, knownNew, embedAll)
		if npre2 != it.npre || diff2 != string(it.diff) {
			// print("following:", "key", it.known, "prev", prev, "known", knownNew,
			// 	"=>", "npre", npre, "diff", diff)
			// adjust following entry
			ins = ins.append(it.offset, npre2, diff2)
			j += fLen(it.diff)
		}
	}
	return nd.replace(i, j, ins)
}

// replace is used by insert and delete
// to replace a portion of a node (i,j) with new content (rep)
func (nd node) replace(i, j int, rep node) node {
	nr := len(rep)
	d := nr - (j - i)
	nd = bytes.Grow(nd, d)
	copy(nd[i+nr:], nd[j:])
	copy(nd[i:], rep)
	if d < 0 {
		nd = nd[:len(nd)+d]
	}
	return nd
}

func (nd node) delete(offset uint64) (node, bool) {
	// search
	var prev []byte
	it := nd.iter()
	for {
		if !it.next() {
			return nil, false // not found
		}
		if it.offset == offset {
			break
		}
		prev = append(prev[:0], it.known...)
	}
	i := it.pos
	// print("i", i)

	j := nd.next(i)
	if j >= len(nd) {
		// delete last item, simplest case, no adjustments
		return nd[:i], true
	}
	// print("1 prev", string(prev), "pos", it.pos, "known", string(it.known))

	rep := make(node, 0, 64)
	if i == 0 {
		// deleting first entry so make following into first
		if !it.next() {
			nd = append(nd[:5], 0, 0)
			return nd, true
		}
		rep = rep.updateCopy(nd, j, 0, "")
		j = nd.next(j)

		// adjust following entry if there is one
		if it.next() {
			diff := string(it.known)
			rep = rep.updateCopy(nd, j, it.npre, diff)
			j = nd.next(j)
			// print("2 prev", string(prev), "known", string(it.known), "diff", diff)
		}
		nd = nd.replace(i, j, rep)
		return nd, true
	}
	calced := append([]byte{}, prev...) // copy
	if it.npre > len(prev) {
		calced = append(calced[:0], it.known[:it.npre]...)
	}
	if it.next() {
		// adjust the following entry
		npre := commonSlicePrefixLen(calced, it.known)
		ndif := commonSlicePrefixLen(prev, it.known)
		diff := string(it.known[ndif:])
		// print("3 prev", string(prev), "calced", calced, "pos", it.pos, "known", string(it.known))
		// print("npre", npre, "n", n, "diff", diff)
		rep = rep.updateCopy(nd, j, npre, diff)
		j = nd.next(j)
	}
	nd = nd.replace(i, j, rep)
	return nd, true
}

func (nd node) updateCopy(src node, i int, npre int, diff string) node {
	nd = append(nd, src[i:i+5]...) // copy offset
	nd = append(nd, byte(npre), byte(len(diff)))
	nd = append(nd, diff...)
	return nd
}

func (nd node) setOffset(pos int, off uint64) {
	stor.WriteSmallOffset(nd[pos:], off)
}

// iter -------------------------------------------------------------

type nodeIter struct {
	node   node
	pos    int // position in node
	npre   int
	diff   []byte
	known  []byte
	offset uint64
}

func (nd node) iter() *nodeIter {
	return &nodeIter{node: nd, pos: -7}
}

func (it *nodeIter) next() bool {
	it.pos += fLen(it.diff)
	if it.pos >= len(it.node) {
		it.known = it.known[:0] // ""
		return false
	}
	it.npre, it.diff, it.offset = it.node[it.pos:].read()

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

func (it *nodeIter) eof() bool {
	return it.pos+fLen(it.diff) >= len(it.node)
}

func (it *nodeIter) copyFrom(src *nodeIter) {
	it.pos = src.pos
	it.npre = src.npre
	it.offset = src.offset
	it.diff = src.diff
	it.known = append(it.known[:0], src.known...) // copy over
}

//-------------------------------------------------------------------

func (nd node) stats() {
	n := nd.check(nil)
	avg := float32(len(nd)-7*n) / float32(n)
	print("n", n, "len", len(nd), "avg", avg)
}

func (nd node) check(get func(uint64) string) int {
	n := 0
	var knownPrev []byte
	var keyPrev string
	it := nd.iter()
	for it.next() {
		known := string(it.known)
		if known < string(knownPrev) {
			panic("out of order: known " + known + ", prev " + string(knownPrev))
		}
		if get != nil {
			key := get(it.offset)
			npre := commonPrefixLen(keyPrev, key)
			if npre > len(known) {
				panic("insufficient known: prev key " + keyPrev +
					", key " + key + ", known" + known)
			}
			if !strings.HasPrefix(key, known) {
				panic("mismatch: known " + known + ", key " + key)
			}
			keyPrev = key
		}
		knownPrev = append(knownPrev[:0], it.known...)
		n++
	}
	return n
}

func (nd node) print() {
	it := nd.iter()
	for it.next() {
		print(it.offset, it.known)
	}
}

func (nd node) printLeafNode(get func(uint64) string) {
	it := nd.iter()
	for it.next() {
		offset := it.offset
		print(strconv.Itoa(it.pos)+":"+
			"{", strconv.Itoa(int(offset))+":", it.npre, it.diff, "}",
			it.known, "("+get(offset)+")")
	}
}

func (nd node) printTreeNode() {
	it := nd.iter()
	for it.next() {
		offset := it.offset
		print(strconv.Itoa(it.pos)+": {", offset, it.npre, it.diff, "}",
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

func (nd node) String() string {
	s := "["
	it := nd.iter()
	for it.next() {
		known := string(it.known)
		if known == "" {
			known = "''"
		}
		s += fmt.Sprint(known, "=", it.offset) + " "
	}
	return strings.TrimSpace(s) + "]"
}

func (nd node) knowns() string {
	var sb strings.Builder
	it := nd.iter()
	for it.next() {
		sb.Write(it.known)
		sb.WriteByte(' ')
		sb.WriteString(offstr(it.offset))
		sb.WriteByte(' ')
	}
	return sb.String()
}
