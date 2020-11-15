// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/ixspec"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

// fbtree is an immutable btree designed to be stored in a file.
type fbtree struct {
	// treeLevels is how many levels of tree nodes there are (initially 0)
	// Nodes do not store whether they are leaf or tree nodes.
	// Since we always start at the root and descend,
	// the code tracks the depth and compares it to treeLevels
	// to differentiate leaf or tree nodes.
	// When the root splits, treeLevels is incremented.
	treeLevels int
	// root is the offset of the root node
	root uint64
	// store is where the btree is stored
	store *stor.Stor
	// ixspec is an opaque value passed to GetLeafKey.
	// It specifies which fields make up the key, based on the schema.
	ixspec *ixspec.T
}

const maxlevels = 8

// MaxNodeSize is the maximum node size in bytes, split if larger.
// Overridden by tests.
var MaxNodeSize = 256 //TODO tune

// GetLeafKey is used to get the key for a data offset.
// It is a dependency that must be injected
var GetLeafKey func(st *stor.Stor, is *ixspec.T, off uint64) string

func CreateFbtree(store *stor.Stor, is *ixspec.T) *fbtree {
	rootNode := fNode{}
	root := rootNode.putNode(store)
	return &fbtree{root: root, store: store, ixspec: is}
}

func OpenFbtree(store *stor.Stor, root uint64, treeLevels int) *fbtree {
	return &fbtree{root: root, treeLevels: treeLevels, store: store}
}

func (fb *fbtree) getLeafKey(off uint64) string {
	return GetLeafKey(fb.store, fb.ixspec, off)
}

func (fb *fbtree) Search(key string) uint64 {
	off := fb.root
	for i := 0; i <= fb.treeLevels; i++ {
		node := fb.getNode(off)
		off, _, _ = node.search(key)
	}
	return off
}

// putNode stores the node
func (node fNode) putNode(store *stor.Stor) uint64 {
	n := len(node)
	off, buf := store.Alloc(2 + n + cksum.Len)
	stor.NewWriter(buf).Put2(n)
	buf = buf[2:]
	copy(buf, node)
	cksum.Update(buf)
	// if len(node) > 0 && rand.Intn(500) == 42 {
	// 	// corrupt some nodes to test checking
	// 	fmt.Println("ZAP")
	// 	buf := store.Data(off)
	// 	buf[3 + rand.Intn(len(node))] = byte(rand.Intn(256))
	// }
	return off
}

// getNode returns the node for a given offset
func (fb *fbtree) getNode(off uint64) fNode {
	return readNode(fb.store, off)
}

func (fb *fbtree) getNodeCk(off uint64, check bool) fNode {
	node := readNode(fb.store, off)
	if check {
		cksum.MustCheck(node[:len(node)+cksum.Len])
	}
	return node
}

func readNode(store *stor.Stor, off uint64) fNode {
	buf := store.Data(off)
	n := stor.NewReader(buf).Get2()
	return fNode(buf[2 : 2+n])
}

//-------------------------------------------------------------------
// Quick check is used when opening a database. It should be fast.
// To be fast it should only look at the end (recent) part of the file.

// recentSize is the length of the tail of the file that we look at
const recentSize = 32 * 1024 * 1024 // ???

func (fb *fbtree) quickCheck() {
	recent := int64(fb.store.Size()) - recentSize
	fb.quickCheck1(0, fb.root, recent)
}

func (fb *fbtree) quickCheck1(depth int, offset uint64, recent int64) {
	// only look at nodes in the recent part of the file
	if int64(offset) < recent {
		return
	}
	node := fb.getNodeCk(offset, true)
	if depth < fb.treeLevels {
		// tree node
		for it := node.iter(); it.next(); {
			fb.quickCheck1(depth+1, it.offset, recent)
		}
	} else {
		// leaf node
		for it := node.iter(); it.next(); {
			// only checksum data records in the recent part of the file
			if int64(it.offset) > recent {
				buf := fb.store.Data(it.offset)
				size := runtime.RecLen(buf)
				cksum.MustCheck(buf[:size+cksum.Len])
			}
		}
	}
}

// check verifies that the keys are in order and returns the number of keys.
// The supplied fn is applied to each leaf offset.
func (fb *fbtree) check(fn func(uint64)) (count, size, nnodes int) {
	key := ""
	return fb.check1(0, fb.root, &key, fn)
}

func (fb *fbtree) check1(depth int, offset uint64, key *string,
	fn func(uint64)) (count, size, nnodes int) {
	node := fb.getNodeCk(offset, true)
	size += len(node)
	nnodes++
	for it := node.iter(); it.next(); {
		offset := it.offset
		if depth < fb.treeLevels {
			// tree
			if it.fi > 0 && *key > string(it.known) {
				panic("keys out of order")
			}
			*key = string(it.known)
			c, s, n := fb.check1(depth+1, offset, key, fn) // RECURSE
			count += c
			size += s
			nnodes += n
		} else {
			// leaf
			count++
			if fn != nil {
				fn(offset)
			}
			itkey := fb.getLeafKey(offset)
			if !strings.HasPrefix(itkey, string(it.known)) {
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

// iter -------------------------------------------------------------

type fbIter = func() (string, uint64, bool)

// Iter returns a function that can be called to return consecutive entries.
// NOTE: The returned key is only the known prefix.
// (unlike inter.Iter which returns the actual key)
func (fb *fbtree) Iter(check bool) fbIter {
	var stack [maxlevels]*fnIter

	// traverse down the tree to the leftmost leaf, making a stack of iterators
	nodeOff := fb.root
	for i := 0; i < fb.treeLevels; i++ {
		stack[i] = fb.getNodeCk(nodeOff, check).iter()
		stack[i].next()
		nodeOff = stack[i].offset
	}
	iter := fb.getNodeCk(nodeOff, check).iter()

	return func() (string, uint64, bool) {
		for {
			if iter.next() {
				return string(iter.known), iter.offset, true // most common path
			}
			// end of leaf, go up the tree
			i := fb.treeLevels - 1
			for ; i >= 0; i-- {
				if stack[i].next() {
					nodeOff = stack[i].offset
					break
				}
			}
			if i == -1 {
				return "", 0, false // eof
			}
			// and then back down to the next leaf
			for i++; i < fb.treeLevels; i++ {
				stack[i] = fb.getNodeCk(nodeOff, check).iter()
				stack[i].next()
				nodeOff = stack[i].offset
			}
			iter = fb.getNodeCk(nodeOff, check).iter()
		}
	}
}

// print ------------------------------------------------------------

func (fb *fbtree) print() {
	fmt.Println("<<<------------------------------")
	fb.print1(0, fb.root)
	fmt.Println("------------------------------>>>")
}

func (fb *fbtree) print1(depth int, offset uint64) {
	explan := ""
	if depth >= fb.treeLevels {
		explan += " LEAF"
	}
	print(strings.Repeat(" . ", depth)+"offset", offset, explan)
	node := fb.getNode(offset)
	var sb strings.Builder
	sep := ""
	for it := node.iter(); it.next(); {
		offset := it.offset
		if depth < fb.treeLevels {
			// tree
			print(strings.Repeat(" . ", depth)+strconv.Itoa(it.fi)+":",
				it.npre, it.diff, "=", it.known)
			fb.print1(depth+1, offset) // recurse
		} else {
			// leaf
			// print(strings.Repeat(" . ", depth)+strconv.Itoa(it.fi)+":",
			// 	strconv.Itoa(int(offset))+",", it.npre, it.diff, "=", it.known,
			// 	"("+fb.getLeafKey(offset)+")")
			sb.WriteString(sep)
			sep = " "
			if len(it.known) == 0 {
				sb.WriteString("''")
			} else {
				sb.Write(it.known)
			}
			// sb.WriteString(" = " + fb.getLeafKey(offset))
		}
	}
	if depth == fb.treeLevels {
		print(strings.Repeat(" . ", depth) + sb.String())
	}
}

// builder ----------------------------------------------------------

// fbtreeBuilder is used to bulk load an fbtree.
// Keys must be added in order.
// The fbtree is built bottom up with no splitting or inserting.
// All nodes will be "full" except for the right hand edge.
type fbtreeBuilder struct {
	levels []*level // leaf is [0]
	prev   string
	store  *stor.Stor
	count  int
}

type level struct {
	splitKey string
	builder  fNodeBuilder
}

func NewFbtreeBuilder(store *stor.Stor) *fbtreeBuilder {
	return &fbtreeBuilder{store: store, levels: []*level{{}}}
}

func (fb *fbtreeBuilder) Add(key string, off uint64) {
	if fb.count > 0 {
		if key == fb.prev {
			panic("fbtreeBuilder keys must not have duplicates")
		}
		if key < fb.prev {
			panic("fbtreeBuilder keys must be inserted in order")
		}
	}
	fb.add(0, key, off)
	fb.prev = key
	fb.count++
}

func (fb *fbtreeBuilder) add(li int, key string, off uint64) {
	if li >= len(fb.levels) {
		fb.levels = append(fb.levels, &level{})
	}
	lev := fb.levels[li]
	if len(lev.builder.fe) > (MaxNodeSize * 3 / 4) {
		// split full node to stor
		offNode, splitKey := lev.builder.Split(fb.store)
		fb.add(li+1, lev.splitKey, offNode) // RECURSE
		lev.splitKey = splitKey
	}
	embedLen := 1
	if li > 0 /*|| fb.count == 1*/ {
		embedLen = 255
	}
	lev.builder.Add(key, off, embedLen)
}

func (fb *fbtreeBuilder) Finish() *Overlay {
	var key string
	var off uint64
	for li := 0; li < len(fb.levels); li++ {
		if li > 0 {
			// allow node to slightly exceed max size
			fb.levels[li].builder.Add(key, off, 255)
		}
		key = fb.levels[li].splitKey
		off = fb.levels[li].builder.fe.putNode(fb.store)
	}
	treeLevels := len(fb.levels) - 1
	bt := OpenFbtree(fb.store, off, treeLevels)
	return &Overlay{fb: bt}
}

// trace ------------------------------------------------------------

const T = false // set to true to enable tracing

func trace(args ...interface{}) bool {
	fmt.Println(args...)
	return true
}

func traced(depth int, args ...interface{}) bool {
	fmt.Print(strings.Repeat("    ", depth))
	fmt.Println(args...)
	return true
}
