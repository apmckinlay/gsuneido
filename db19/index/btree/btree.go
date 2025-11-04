// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"log"
	"math/bits"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

var _ iface.Btree = (*btree)(nil)

type T = btree

// btree is an immutable btree designed to be stored in a file.
//
// To update a btree, the changes are added to an ixbuf
// which is then merged to create a new btree.
type btree struct {
	// stor is where the btree is stored
	stor *stor.Stor
	// ixspec is an opaque value passed to GetLeafKey.
	// It specifies which fields make up the key, based on the schema.
	ixspec *ixkey.Spec
	// treeLevels is how many levels of tree nodes there are (initially 0)
	// Nodes do not store whether they are leaf or tree nodes.
	// Since we always start at the root and descend,
	// the code tracks the depth and compares it to treeLevels
	// to differentiate leaf or tree nodes.
	// When the root splits, treeLevels is incremented.
	treeLevels int
	// root is the offset of the root node
	root uint64
	// rootUnode is an uncompressed copy for faster access
	rootUnode unode
}

func (bt *btree) Cksum() uint32 {
	return uint32(bt.treeLevels) + uint32(bt.root)
}

const maxlevels = 8

// MaxNodeSize is the maximum node size in bytes, split if larger.
// var rather than const because it is overridden by tests.
var MaxNodeSize = 1024

const MinSplitSize = 6 // for builder that will be split 4 and 2

// EntrySize is the estimated average entry size
const EntrySize = 10

// TreeHeight is the estimated average tree height.
// It is used by Table.lookupCost
const TreeHeight = 3

var Fanout = MaxNodeSize / EntrySize // estimate ~100

// GetLeafKey is used to get the key for a data offset.
// It is a dependency that must be injected
var GetLeafKey func(st *stor.Stor, is *ixkey.Spec, off uint64) string

func CreateBtree(st *stor.Stor, is *ixkey.Spec) iface.Btree {
	rootNode := node{}
	rootOff := rootNode.putNode(st)
	return &btree{root: rootOff, stor: st, ixspec: is}
}

func OpenBtree(st *stor.Stor, root uint64, treeLevels int) iface.Btree {
	ru := readNode(st, root).toUnode()
	return &btree{root: root, treeLevels: treeLevels, stor: st, rootUnode: ru}
}

func (bt *btree) SetSplit(ndsize int) {
	panic("not implemented")
}

func (bt *btree) GetIxspec() *ixkey.Spec {
	return bt.ixspec
}

func (bt *btree) SetIxspec(is *ixkey.Spec) {
	bt.ixspec = is
}

func (bt *btree) TreeLevels() int {
	return bt.treeLevels
}

func (bt *btree) getLeafKey(off uint64) string {
	return GetLeafKey(bt.stor, bt.ixspec, off)
}

// Lookup returns the offset for a key, or 0 if not found.
func (bt *btree) Lookup(key string) uint64 {
	off := bt.rootUnode.search(key)
	for range bt.treeLevels {
		nd := bt.getNode(off)
		off = nd.search(key)
	}
	if off == 0 || bt.getLeafKey(off) != key {
		return 0
	}
	return off
}

// putNode stores the node
func (nd node) putNode(st *stor.Stor) uint64 {
	n := len(nd)
	if n > 8192 {
		log.Println("ERROR: btree node too large")
	}
	off, buf := st.Alloc(2 + n + cksum.Len)
	stor.NewWriter(buf).Put2(n)
	buf = buf[2:]
	copy(buf, nd)
	cksum.Update(buf)
	// if len(nd) > 0 && rand.Intn(500) == 42 {
	// 	// corrupt some nodes to test checking
	// 	fmt.Println("ZAP")
	// 	buf := st.Data(off)
	// 	buf[3 + rand.Intn(len(nd))] = byte(rand.Intn(256))
	// }
	return off
}

// PutEmptyNode is for tests
func PutEmptyNode(st *stor.Stor) {
	var nd node
	off := nd.putNode(st)
	assert.That(off == 0)
}

// getNode returns the node for a given offset
func (bt *btree) getNode(off uint64) node {
	return readNode(bt.stor, off)
}

func (bt *btree) getNodeCk(off uint64, check bool) node {
	nd := readNode(bt.stor, off)
	if check {
		cksum.MustCheck(nd[:len(nd)+cksum.Len])
	}
	return nd
}

func readNode(st *stor.Stor, off uint64) node {
	buf := st.Data(off)
	n := stor.NewReader(buf).Get2()
	return node(buf[2 : 2+n])
}

//-------------------------------------------------------------------
// Quick check is used when opening a database. It should be fast.
// To be fast it should only look at the end (recent) part of the file.

// recentSize is the length of the tail of the file that we look at
const recentSize = 32 * 1024 * 1024 // ???

func (bt *btree) QuickCheck() {
	recent := int64(bt.stor.Size()) - recentSize
	bt.quickCheck1(0, bt.root, recent)
}

func (bt *btree) quickCheck1(depth int, offset uint64, recent int64) {
	// only look at nodes in the recent part of the file
	if int64(offset) < recent {
		return
	}
	nd := bt.getNodeCk(offset, true)
	if depth < bt.treeLevels {
		// tree node
		for it := nd.iter(); it.next(); {
			bt.quickCheck1(depth+1, it.offset, recent)
		}
	} else {
		// leaf node
		for it := nd.iter(); it.next(); {
			// only checksum data records in the recent part of the file
			if int64(it.offset) > recent {
				buf := bt.stor.Data(it.offset)
				size := core.RecLen(buf)
				cksum.MustCheck(buf[:size+cksum.Len])
			}
		}
	}
}

// Check verifies that the keys are in order and returns the number of keys.
// If the supplied function is not nil, it is applied to each leaf offset.
func (bt *btree) Check(fn any) (count, size, nnodes int) {
	key := ""
	return bt.check1(0, bt.root, &key, fn)
}

func (bt *btree) check1(depth int, offset uint64, key *string, fn any) (count, size, nnodes int) {
	nd := bt.getNodeCk(offset, true)
	if len(nd) == 0 && (bt.treeLevels > 0 || depth > 0) {
		panic("empty node in non-empty btree")
	}
	size += len(nd)
	nnodes++
	for it := nd.iter(); it.next(); {
		off := it.offset
		if depth < bt.treeLevels {
			// tree
			if it.pos > 0 && *key > string(it.known) {
				panic("keys out of order")
			}
			*key = string(it.known)
			c, s, n := bt.check1(depth+1, off, key, fn) // RECURSE
			count += c
			size += s
			nnodes += n
		} else {
			// leaf
			count++
			itkey := bt.getLeafKey(off)
			switch fn := fn.(type) {
			case func(uint64):
				fn(off)
			case func(string, uint64):
				fn(itkey, off)
			default:
			}
			if !strings.HasPrefix(itkey, string(it.known)) {
				// fmt.Printf("known %q index %q\nvalues %v\n",
				// 	string(it.known), itkey, ixkey.DecodeValues(itkey))
				panic("index key does not match data")
			}
			if *key > itkey {
				panic("keys out of order")
			}
			*key = itkey
		}
	}
	return
}

type Stats struct {
	Levels  int
	Count   int
	Size    int
	Nnodes  int
	RootN   int
	Fan     int
	NodeFan [8]int
}

func (bt *btree) Stats() (stats Stats) {
	stats.Levels = bt.treeLevels + 1
	bt.stats(0, bt.root, &stats)
	stats.Fan /= stats.Nnodes
	return
}

func (bt *btree) stats(depth int, offset uint64, stats *Stats) {
	nd := bt.getNode(offset)
	stats.Nnodes++
	stats.Size += len(nd)
	n := uint16(0)
	for it := nd.iter(); it.next(); n++ {
		if depth == 0 {
			stats.RootN++
		} else {
			stats.Fan++
		}
		offset := it.offset
		if depth < bt.treeLevels { // tree
			bt.stats(depth+1, offset, stats) // RECURSE
		} else { // leaf
			stats.Count++
		}
	}
	stats.NodeFan[16-bits.LeadingZeros16(n)]++
}

func (stats Stats) String() string {
	s := fmt.Sprintln(
		"lv", stats.Levels,
		" n ", stats.Count,
		" sz ", stats.Size,
		" nn ", stats.Nnodes,
		" rn ", stats.RootN,
		" f ", stats.Fan) + "    >= "
	for i, n := range stats.NodeFan {
		if n > 0 {
			s += fmt.Sprintf("%d: %d ", (1<<i)/2, n)
		}
	}
	return s
}

// print ------------------------------------------------------------

func (bt *btree) Print() {
	fmt.Println("<<<------------------------------")
	bt.print1(0, bt.root)
	fmt.Println("------------------------------>>>")
}

func (bt *btree) print1(depth int, offset uint64) {
	explan := ""
	if depth >= bt.treeLevels {
		explan += " LEAF"
	}
	print(strings.Repeat(" . ", depth)+"offset", offset, explan)
	nd := bt.getNode(offset)
	var sb strings.Builder
	sep := ""
	for it := nd.iter(); it.next(); {
		offset := it.offset
		if depth < bt.treeLevels {
			// tree
			print(strings.Repeat(" . ", depth)+strconv.Itoa(it.pos)+":",
				it.npre, it.diff, "=", it.known)
			bt.print1(depth+1, offset) // recurse
		} else {
			// leaf
			// print(strings.Repeat(" . ", depth)+strconv.Itoa(it.pos)+":",
			// 	strconv.Itoa(int(offset))+",", it.npre, it.diff, "=", it.known,
			// 	"("+bt.getLeafKey(offset)+")")
			sb.WriteString(sep)
			sep = ", "
			// if len(it.known) == 0 {
			// 	sb.WriteString("''")
			// } else {
			// 	sb.Write(it.known)
			// }
			// fmt.Fprintf(&sb, /*" = " +*/"%q %d", bt.getLeafKey(offset), offset)
			fmt.Fprintf(&sb, "%d", offset)
		}
	}
	if depth == bt.treeLevels {
		print(strings.Repeat(" . ", depth) + sb.String())
	}
}

func (bt *btree) NodeSizes() {
	bt.nodeSizes(0, bt.root)
}

func (bt *btree) nodeSizes(depth int, offset uint64) int {
	nd := bt.getNode(offset)
	n := 0
	for it := nd.iter(); it.next(); {
		if depth < bt.treeLevels {
			n += bt.nodeSizes(depth+1, it.offset)
		} else {
			n++
		}
	}
	fmt.Println(depth, n)
	return n
}

//-------------------------------------------------------------------

func (bt *btree) Write(w *stor.Writer) {
	w.Put5(int64(bt.root)).Put1(bt.treeLevels)
}

func Read(st *stor.Stor, r *stor.Reader) iface.Btree {
	root := uint64(r.Get5())
	treeLevels := r.Get1()
	return OpenBtree(st, root, treeLevels)
}

//-------------------------------------------------------------------

// RangeFrac returns the fraction of the btree (0 to 1) in the range org to end
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
	return len(bt.rootUnode) == 0
}

func (bt *btree) fracPos(key string) float64 {
	if key == ixkey.Min {
		return 0
	}
	if key == ixkey.Max {
		return 1
	}
	frac := float64(0)
	div := float64(1)
	off := bt.root
	exact := true
	for level := 0; level <= bt.treeLevels; level++ {
		node := bt.getNode(off)
		i := 0
		n := 0
		for it := node.iter(); it.next(); n++ {
			k := string(it.known)
			if key >= k {
				i = n
				off = it.offset
			}
		}
		const smallRoot = 10 // ???
		if level == 0 && n < smallRoot {
			i, n, off = bt.rootChildren(node, key)
			level++
		}
		if n == 0 {
			return frac
		}
		frac += float64(i) / float64(n) / div
		if exact {
			exact = false
			div = float64(n)
		} else {
			div *= float64(Fanout) // ???
		}
	}
	return frac
}

// rootChildren helps when the root is small
// by scanning the children as if they were one bigger node
func (bt *btree) rootChildren(root node, key string) (i, n int, off uint64) {
	for rit := root.iter(); rit.next(); {
		node := bt.getNode(rit.offset)
		for it := node.iter(); it.next(); n++ {
			k := string(it.known)
			if k == "" {
				k = string(rit.known)
			}
			if key >= k {
				i = n
				off = it.offset
			}
		}
	}
	return
}

// trace ------------------------------------------------------------

const t = false // set to true to enable tracing

func trace(args ...any) bool {
	fmt.Println(args...)
	return true
}
