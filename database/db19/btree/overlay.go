// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

type treeIter = func() (string, uint64, bool)

type tree interface {
	Iter() treeIter
}

// Overlay is an immutable fbtree plus one or more mbtrees.
// with a mutable mbtree at the top to store updates.
type Overlay struct {
	// under are the underlying fbtree and mbtree's
	under []tree
	// mb is the mutable top mbtree
	mb *mbtree
}

// Insert inserts into the mutable top mbtree
func (ov *Overlay) Insert(key string, off uint64) {
	ov.mb.Insert(key, off)
}

const tombstone = 1 << 63

// Delete either deletes the key/offset from the mutable mbtree
// or inserts a tombstone into the mutable mbtree.
func (ov *Overlay) Delete(key string, off uint64) {
	if !ov.mb.Delete(key, off) {
		// key not present
		ov.mb.Insert(key, off|tombstone)
	}
}

//-------------------------------------------------------------------

type ovsrc struct {
	iter treeIter
	key  string
	off  uint64
	ok   bool
}

// Iter returns an treeIter function
func (ov *Overlay) Iter() treeIter {
	srcs := make([]ovsrc, len(ov.under)+1)
	srcs[0] = ovsrc{iter: ov.mb.Iter()}
	for i := 1; i < len(srcs); i++ {
		srcs[i] = ovsrc{iter: ov.under[i-1].Iter()}
	}
	for i := range srcs {
		srcs[i].next()
	}
	return func() (string, uint64, bool) {
		i := ovsrcNext(srcs)
		key, off, ok := srcs[i].key, srcs[i].off, srcs[i].ok
		srcs[i].next()
		return key, off >> 1, ok
	}
}

func (src *ovsrc) next() {
	src.key, src.off, src.ok = src.iter()
	// adjust offset so tombstone comes first
	src.off = (src.off << 1) | ((src.off >> 63) ^ 1)
}

// ovsrcNext returns the index of the next element
func ovsrcNext(srcs []ovsrc) int {
	min := 0
	for {
		for i := 1; i < len(srcs); i++ {
			if ovsrcLess(&srcs[i], &srcs[min]) {
				min = i
			}
		}
		if !isTombstone(srcs[min].off) {
			return min
		}
		// skip over the insert matching the tombstone
		for i := range srcs {
			if i != min &&
				srcs[i].key == srcs[min].key && srcs[i].off&^1 == srcs[min].off {
				srcs[i].next()
			}
		}
		srcs[min].next() // skip the tombstone itself
	}
}

func isTombstone(off uint64) bool {
	return (off & 1) == 0
}

func ovsrcLess(x, y *ovsrc) bool {
	if !x.ok {
		return false
	}
	return !y.ok || x.key < y.key || (x.key == y.key && x.off < y.off)
}
