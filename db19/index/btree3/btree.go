// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package btree is a btree implementation.
Technically it is a B+ tree since only the leaf nodes point to the data.
Builder produces btrees where the righthand edge may be as small as one child.
Deletes applied by MergeAndSave do not combine siblings
and may also produce nodes as small as one child (but empty nodes are removed)
The root may be as small as two children.

btree is treated as immutable.
Two things create or modify btrees:
- bulk in-order loading from load or compact (Builder)
- batch add, update, and delete from an ixbuf (MergeAndSave)

A btree consists of two kinds of nodes: leafNode and treeNode.
There are treeLevels of treeNodes.
The root is a leafNode if treeLevels == 0, otherwise it is a treeNode.

leafNode and treeNode have the same basic representations
but for performance have separate implementations.
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
	"github.com/apmckinlay/gsuneido/util/hacks"
)

var _ iface.Btree = (*btree)(nil)

const splitCount = 100   // ???
const minSplit = 1024    // ???
const maxNodeSize = 8192 // ???

// TreeHeight is the estimated average tree height.
// It is used by Table.lookupCost
const TreeHeight = 3 // = 10,000 to 1,000,000 keys with fanout of 100

type btree struct {
	stor       *stor.Stor
	root       uint64
	treeLevels int
	// count is used by RangeFrac, set by Read, Builder, and MergeAndSave
	count       int
	shouldSplit func(node) bool // overridden by tests
}

type T = btree

func CreateBtree(st *stor.Stor, _ *ixkey.Spec) iface.Btree {
	return Builder(st).Finish()
}

func OpenBtree(st *stor.Stor, root uint64, treeLevels int, nrows int) iface.Btree {
	return &btree{stor: st, root: root, treeLevels: treeLevels,
		count: nrows, shouldSplit: shouldSplit}
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

func Read(st *stor.Stor, r *stor.Reader, nrows int) iface.Btree {
	root := uint64(r.Get5())
	treeLevels := r.Get1()
	return OpenBtree(st, root, treeLevels, nrows)
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

//-------------------------------------------------------------------

// Check verifies that the keys are in order and returns the number of keys.
// If the supplied function is not nil, it is applied to each leaf offset.
// WARNING: the key string references a reused byte slice - don't hold it
func (bt *btree) Check(fn any) (count, size, nnodes int) {
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
			size += len(nd)
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
			size += len(nd)
			var prevSuffix []byte
			first := true
			for it := nd.iter(); it.next(); count++ {
				suffix := it.suffix()
				if !first {
					assert.That(string(suffix) > string(prevSuffix))
				}
				first = false
				prevSuffix = suffix
				switch fn := fn.(type) {
				case func():
					fn()
				case func(uint64):
					fn(it.offset())
				case func(string, uint64):
					prev = append(append((prev)[:0], nd.prefix()...), suffix...)
					fn(hacks.BStoS(prev), it.offset())
				default:
				}
			}
			prev = append(append((prev)[:0], nd.prefix()...),
				nd.suffix(nd.nkeys()-1)...)
		}
	}
	check1(0, bt.root)
	assert.That(count == bt.count)
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

//-------------------------------------------------------------------

type node interface {
	noffs() int
	size() int
	offset(i int) uint64
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

// readNode reads a node at the given level and offset, level 0 is the root
func (bt *btree) readNode(level int, off uint64) node {
	if level < bt.treeLevels {
		return bt.readTree(off)
	}
	return bt.readLeaf(off)
}

// ------------------------------------------------------------------

type Stats struct {
	Levels  int
	Count   int
	Size    int
	AvgSize int
	Nleaf   int
	Ntree   int
	LeafFan int
	TreeFan int
	RootFan int
}

func (bt *btree) Stats() (stats Stats) {
	stats.Levels = bt.treeLevels + 1
	bt.stats(0, bt.root, &stats)
	stats.AvgSize = stats.Size / (1 + stats.Ntree + stats.Nleaf)
	if stats.Nleaf > 0 {
		stats.LeafFan /= stats.Nleaf
	}
	if stats.Ntree > 0 {
		stats.TreeFan /= stats.Ntree
	}
	return
}

func (bt *btree) stats(depth int, offset uint64, stats *Stats) {
	nd := bt.readNode(depth, offset)
	stats.Size += nd.size()
	n := nd.noffs()
	if depth == 0 {
		stats.RootFan = n
	} else if depth < bt.treeLevels {
		stats.Ntree++
		stats.TreeFan += n
	} else {
		stats.Nleaf++
		stats.LeafFan += n
	}
	for i := range n {
		offset := nd.offset(i)
		if depth < bt.treeLevels { // tree
			bt.stats(depth+1, offset, stats) // RECURSE
		} else { // leaf
			stats.Count++
		}
	}
}

func (stats Stats) String() string {
	return fmt.Sprint(
		"lv ", stats.Levels,
		" n ", stats.Count,
		" sz ", stats.Size,
		" as ", stats.AvgSize,
		" tn ", stats.Ntree,
		" ln ", stats.Nleaf,
		" lf ", stats.LeafFan,
		" tf ", stats.TreeFan,
		" rf ", stats.RootFan)
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
