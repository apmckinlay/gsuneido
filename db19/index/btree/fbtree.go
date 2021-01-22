// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

type T = fbtree

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
	ixspec *ixkey.Spec
}

const maxlevels = 8

// MaxNodeSize is the maximum node size in bytes, split if larger.
// Overridden by tests.
var MaxNodeSize = 256 //TODO tune

// GetLeafKey is used to get the key for a data offset.
// It is a dependency that must be injected
var GetLeafKey func(st *stor.Stor, is *ixkey.Spec, off uint64) string

func CreateFbtree(store *stor.Stor, is *ixkey.Spec) *fbtree {
	rootNode := fnode{}
	root := rootNode.putNode(store)
	return &fbtree{root: root, store: store, ixspec: is}
}

func OpenFbtree(store *stor.Stor, root uint64, treeLevels int) *fbtree {
	return &fbtree{root: root, treeLevels: treeLevels, store: store}
}

func (fb *fbtree) GetIxspec() *ixkey.Spec {
	return fb.ixspec
}

func (fb *fbtree) SetIxspec(is *ixkey.Spec) {
	fb.ixspec = is
}

func (fb *fbtree) getLeafKey(off uint64) string {
	return GetLeafKey(fb.store, fb.ixspec, off)
}

// Lookup returns the offset for a key, or 0 if not found.
func (fb *fbtree) Lookup(key string) uint64 {
	off := fb.root
	for i := 0; i <= fb.treeLevels; i++ {
		node := fb.getNode(off)
		off = node.search(key)
	}
	if fb.getLeafKey(off) != key {
		return 0
	}
	return off
}

// putNode stores the node
func (node fnode) putNode(store *stor.Stor) uint64 {
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
func (fb *fbtree) getNode(off uint64) fnode {
	return readNode(fb.store, off)
}

func (fb *fbtree) getNodeCk(off uint64, check bool) fnode {
	node := readNode(fb.store, off)
	if check {
		cksum.MustCheck(node[:len(node)+cksum.Len])
	}
	return node
}

func readNode(store *stor.Stor, off uint64) fnode {
	buf := store.Data(off)
	n := stor.NewReader(buf).Get2()
	return fnode(buf[2 : 2+n])
}

//-------------------------------------------------------------------
// Quick check is used when opening a database. It should be fast.
// To be fast it should only look at the end (recent) part of the file.

// recentSize is the length of the tail of the file that we look at
const recentSize = 32 * 1024 * 1024 // ???

func (fb *fbtree) QuickCheck() {
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

// Check verifies that the keys are in order and returns the number of keys.
// The supplied fn is applied to each leaf offset.
func (fb *fbtree) Check(fn func(uint64)) (count, size, nnodes int) {
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

//-------------------------------------------------------------------

func (fb *fbtree) StorSize() int {
	return 5 + 1
}

func (fb *fbtree) Write(w *stor.Writer) {
	w.Put5(fb.root).Put1(fb.treeLevels)
}

// ReadOverlay reads an Overlay from storage BUT without ixspec
func Read(st *stor.Stor, r *stor.Reader) *fbtree {
	root := r.Get5()
	treeLevels := r.Get1()
	return OpenFbtree(st, root, treeLevels)
}

// trace ------------------------------------------------------------

const t = false // set to true to enable tracing

func trace(args ...interface{}) bool {
	fmt.Println(args...)
	return true
}
