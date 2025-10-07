// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
db19/index/btree3 is a new, optimized btree implementation.
It is not used yet.

Two things create or modify btrees:
- bulk in-order loading from load or compact (Builder)
- add, update, and delete from an ixbuf (MergeAndSave)

A btree consists of two kinds of nodes: leafNode and treeNode.
There are treeLevels of treeNodes.
The root is a leafNode if treeLevels == 0, otherwise it is a treeNode.

leafNode and treeNode have the same basic representations.
leafNode has prefix compression.
treeNode has an extra offset since the fields are separators.
*/
package btree

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

var _ iface.Btree = (*btree)(nil)

// MaxNodeSize is the maximum node size in bytes, split if larger.
// var rather than const because it is overridden by tests.
const minSplit = 1024 // ???
const maxSplit = 8192 // ???
const splitSize = 100 // ???

// avgFanout is the estimated average number of children per node
const avgFanout = splitSize * 95 / 100 // ???

// TreeHeight is the estimated average tree height.
// It is used by Table.lookupCost
const TreeHeight = 3 // = 10,000 to 1,000,000 keys with fanout of 100

type btree struct {
	stor        *stor.Stor
	root        uint64
	treeLevels  int
	shouldSplit func(node) bool
}

type T = btree

func CreateBtree(st *stor.Stor, _ *ixkey.Spec) iface.Btree {
	return Builder(st).Finish()
}

func OpenBtree(st *stor.Stor, root uint64, treeLevels int) iface.Btree {
	return &btree{stor: st, root: root, treeLevels: treeLevels,
		shouldSplit: shouldSplit}
}

// SetSplit is for tests
func (bt *btree) SetSplit(ndsize int) {
	bt.shouldSplit = func(nd node) bool {
		return nd.size() >= ndsize
	}
}

func (bt *btree) Cksum() uint32 {
	return uint32(bt.treeLevels) + uint32(bt.root)
}

func (bt *btree) TreeLevels() int {
	return bt.treeLevels
}

func (bt *btree) SetIxspec(is *ixkey.Spec) {
	// temporary for transition
}

func (bt *btree) Write(w *stor.Writer) {
	w.Put5(int64(bt.root)).Put1(bt.treeLevels)
}

func Read(st *stor.Stor, r *stor.Reader) iface.Btree {
	root := uint64(r.Get5())
	treeLevels := r.Get1()
	return OpenBtree(st, root, treeLevels)
}

// Lookup returns the offset for a key, or 0 if not found.
func (bt *btree) Lookup(key string) uint64 {
	off := bt.root
	for range bt.treeLevels {
		nd := bt.readTree(off)
		_, off = nd.search(key)
	}
	nd := bt.readLeaf(off)
	i, found := nd.search(key)
	if !found {
		return 0 // not found
	}
	return nd.offset(i)
}

func (bt *btree) readTree(off uint64) treeNode {
	return readTree(bt.stor, off)
}

func (bt *btree) readLeaf(off uint64) leafNode {
	return readLeaf(bt.stor, off)
}

// Check verifies that the keys are in order and returns the number of keys.
// If the supplied function is not nil, it is applied to each leaf offset.
func (bt *btree) Check(fn func(uint64)) (count, size, nnodes int) {
	var prev []byte // updated by leaf
	var check1 func(int, uint64)
	check1 = func(depth int, offset uint64) {
		nnodes++
		if depth < bt.treeLevels {
			// tree
			nd := bt.readTreeCk(offset)
			if depth == 0 {
				assert.That(nd.nkeys() >= 1)
			} else {
				assert.That(nd.noffs() >= 1)
			}
			size = len(nd)
			for i := 0; i < nd.nkeys(); i++ {
				sep := nd.key(i)
				assert.That(string(sep) > string(prev))
				check1(depth+1, nd.offset(i)) // RECURSE
			}
			check1(depth+1, nd.offset(nd.nkeys())) // RECURSE
		} else {
			// leaf
			nd := bt.readLeafCk(offset)
			if nd.nkeys() == 0 {
				assert.That(bt.treeLevels == 0)
				return
			}
			k := nd.key(0)
			assert.That(count == 0 || k > string(prev))
			size = len(nd)
			var prevSuffix []byte
			first := true
			for it := nd.iter(); it.next(); count++ {
				k := it.suffix()
				if !first {
					assert.That(string(k) > string(prevSuffix))
				}
				first = false
				prevSuffix = k
				if fn != nil {
					fn(it.offset())
				}
			}
			prev = append(append((prev)[:0], nd.prefix()...),
				nd.suffix(nd.nkeys()-1)...)
		}
	}
	check1(0, bt.root)
	return
}

// Quick check is used when opening a database. It should be fast.
// To be fast it should only look at the end (recent) part of the file.
func (bt *btree) QuickCheck() {
	const recentSize = 32 * 1024 * 1024 // ???
	recent := int64(bt.stor.Size()) - recentSize
	bt.quickCheck1(0, bt.root, recent)
}

func (bt *btree) quickCheck1(depth int, offset uint64, recent int64) {
	// only look at nodes in the recent part of the file
	if int64(offset) < recent {
		return
	}
	if depth < bt.treeLevels {
		// tree node
		nd := bt.readTreeCk(offset)
		for it := nd.iter(); it.next(); {
			bt.quickCheck1(depth+1, it.offset(), recent) // RECURSE
		}
	} else {
		// leaf node
		nd := bt.readLeafCk(offset)
		for it := nd.iter(); it.next(); {
			// only checksum data records in the recent part of the file
			if int64(it.offset()) > recent {
				buf := bt.stor.Data(it.offset())
				size := core.RecLen(buf)
				cksum.MustCheck(buf[:size+cksum.Len])
			}
		}
	}
}

func (bt *btree) readLeafCk(offset uint64) leafNode {
	nd := readLeaf(bt.stor, offset)
	cksum.MustCheck(nd[:len(nd)+cksum.Len])
	return nd
}

func (bt *btree) readTreeCk(offset uint64) treeNode {
	nd := readTree(bt.stor, offset)
	cksum.MustCheck(nd[:len(nd)+cksum.Len])
	return nd
}

func (bt *btree) RangeFrac(org, end string, nrecs int) float64 {
	if bt.empty() || nrecs == 0 {
		// don't know if table is empty or if there are records in the ixbufs
		// fraction is between 0 and 1 so just return half
		return .5
	}

	// count the records (up to iterLimit) to get an exact result
	const iterLimit = 100 // ???
	it := bt.Iterator()
	it.Range(Range{Org: org, End: end})
	n := 0
	for it.Next(); n < iterLimit; it.Next() {
		if it.Eof() {
			return float64(n) / float64(nrecs)
		}
		n++
	}
	minResult := iterLimit / float64(nrecs)

	frac := bt.fracPos(end) - bt.fracPos(org)
	return max(frac, minResult)
}

func (bt *btree) empty() bool {
	if bt.treeLevels > 0 {
		return false
	}
	nd := bt.readLeaf(bt.root)
	return nd.nkeys() == 0
}

const smallRoot = 8  // ???
const largeRoot = 50 // ???

func (bt *btree) fracPos(key string) float64 {
	_ = t && trace("=== fracPos", key)
	if key == ixkey.Min {
		return 0
	}
	if key == ixkey.Max {
		return 1
	}
	root := bt.readNode(0, bt.root)
	n := root.noffs()
	i, off := search(root, key)
	if bt.treeLevels == 0 || n >= largeRoot {
		// only use the root node
		return float64(i) / float64(n)
	} else if n > smallRoot {
		// read the search node
		nd := bt.readNode(1, off)
		j, _ := search(nd, key)
		m := nd.noffs()
		// get the size of the rightmost node (unless it's the search node j,m)
		last := avgFanout
		if j < n-1 {
			rn := bt.readNode(1, root.(treeNode).offset(n-1))
			last = rn.noffs()
		}
		levelSize := (n-2)*avgFanout + m + last
		levelPos := i*avgFanout + j
		return float64(levelPos) / float64(levelSize)
	} else { // n <= smallRoot
		// read the whole level
		var ni, m, j int
		for it := root.(treeNode).iter(); it.next(); ni++ {
			node := bt.readNode(1, it.offset())
			no := node.noffs()
			m += no
			if ni < i {
				j += no
			} else if ni == i {
				k, _ := search(node, key)
				j += k
			}
		}
		return float64(j) / float64(m)
	}
}

type node interface {
	noffs() int
	size() int
}

func search[T node](nd T, key string) (int, uint64) {
	switch nd := any(nd).(type) {
	case leafNode:
		i, _ := nd.search(key)
		return i, 0
	case treeNode:
		return nd.search(key)
	}
	panic("unreachable")
}

func (bt *btree) readNode(level int, off uint64) node {
	if level < bt.treeLevels {
		return bt.readTree(off)
	}
	return bt.readLeaf(off)
}

// ------------------------------------------------------------------

func (bt *btree) Print() {
	fmt.Println("-----------------------------")
	bt.print1(0, bt.root)
}

func (bt *btree) print1(depth int, offset uint64) {
	indent := strings.Repeat(" .", depth)
	if depth < bt.treeLevels {
		nd := readTree(bt.stor, offset)
		fmt.Println(indent, offset, "->", nd)
		for i := 0; i < nd.nkeys(); i++ {
			bt.print1(depth+1, nd.offset(i)) // RECURSE
			fmt.Println(indent, "<"+string(nd.key(i))+">")
		}
		bt.print1(depth+1, nd.offset(nd.nkeys())) // RECURSE
	} else {
		nd := readLeaf(bt.stor, offset)
		fmt.Println(indent, offset, "->", nd)
	}
}
