// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/db19/btree/inter"
	"github.com/apmckinlay/gsuneido/db19/ixspec"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type treeIter = func() (string, uint64, bool)

type tree interface {
	Iter(check bool) treeIter
}

// Overlay is an immutable fbtree plus one or more mbtrees.
// Update transactions have a mutable inter.T at the top to store updates.
type Overlay struct {
	// under are the underlying fbtree and inter.T's
	under []tree
	// mut is the mutable top inter.T, nil if read-only
	mut *inter.T
}

func NewOverlay(store *stor.Stor, is *ixspec.T) *Overlay {
	assert.That(is != nil)
	return &Overlay{under: []tree{CreateFbtree(store, is)}}
}

// var Under [8]int64

// Mutable returns a modifiable copy of an Overlay
func (ov *Overlay) Mutable(tranNum int) *Overlay {
	assert.That(ov.mut == nil)
	under := append([]tree(nil), ov.under...) // copy
	// atomic.AddInt64(&Under[len(under)], 1)
	return &Overlay{under: under, mut: &inter.T{TranNum: tranNum}}
}

func (ov *Overlay) GetIxspec() *ixspec.T {
	return ov.base().ixspec
}

func (ov *Overlay) SetIxspec(is *ixspec.T) {
	ov.base().ixspec = is
}

// Insert inserts into the mutable top inter.T
func (ov *Overlay) Insert(key string, off uint64) {
	ov.mut.Insert(key, off)
}

const tombstone = 1 << 63

// Delete either deletes the key/offset from the mutable inter.T
// or inserts a tombstone into the mutable inter.T.
func (ov *Overlay) Delete(key string, off uint64) {
	if !ov.mut.Delete(key) {
		// key not present
		ov.mut.Insert(key, off|tombstone)
	}
}

func (ov *Overlay) Check(fn func(uint64)) int {
	n, _, _ := ov.base().check(fn)
	return n
}

func (ov *Overlay) QuickCheck() {
	ov.base().quickCheck()
}

func (ov *Overlay) base() *fbtree {
	return ov.under[0].(*fbtree)
}

// iter -------------------------------------------------------------

type ovsrc struct {
	iter treeIter
	key  string
	off  uint64
	ok   bool
}

// Iter returns a treeIter function
func (ov *Overlay) Iter(check bool) treeIter {
	if ov.mut == nil && len(ov.under) == 1 {
		// only fbtree, no merge needed
		return ov.under[0].Iter(check)
	}
	srcs := make([]ovsrc, 0, len(ov.under)+1)
	if ov.mut != nil {
		srcs = append(srcs, ovsrc{iter: ov.mut.Iter(check)})
	}
	for i := range ov.under {
		srcs = append(srcs, ovsrc{iter: ov.under[i].Iter(check)})
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

//-------------------------------------------------------------------

func (ov *Overlay) StorSize() int {
	return 5 + 1 + 5
}

func (ov *Overlay) Write(w *stor.Writer) {
	fb := ov.base()
	w.Put5(fb.root).Put1(fb.treeLevels).Put5(fb.redirsOff)
}

// ReadOverlay reads an Overlay from storage BUT without ixspec
func ReadOverlay(st *stor.Stor, r *stor.Reader) *Overlay {
	root := r.Get5()
	treeLevels := r.Get1()
	redirsOff := r.Get5()
	return &Overlay{under: []tree{OpenFbtree(st, root, treeLevels, redirsOff)}}
}

//-------------------------------------------------------------------

// UpdateWith takes the inter.T updates from ov2 and adds them as a new layer to ov
func (ov *Overlay) UpdateWith(latest *Overlay) {
	// reuse the new slice and overwrite ov.under with the latest
	ov.under = append(ov.under[:0], latest.under...)
	// add inter.T updates
	ov.Freeze()
}

func (ov *Overlay) Freeze() {
	ov.under = append(ov.under, ov.mut)
	ov.mut = nil
}

//-------------------------------------------------------------------

type Result = *fbtree

// Merge merges the inter.T for tranNum (if there is one) into the fbtree
func (ov *Overlay) Merge(tns []int) Result {
	assert.That(ov.mut == nil)
	for i, tn := range tns {
		mut := ov.under[1+i].(*inter.T)
		assert.That(mut.TranNum == tn)
	}
	return ov.merge(len(tns))
}

// merge combines the base fbtree with one or more of the inter.T's
// to produce a new fbtree. It does not modify the original fbtree or inter.T's.
func (ov *Overlay) merge(nmb int) *fbtree {
	return ov.base().Update(func(fb *fbtree) {
		// ??? maybe faster to merge-iterate the inter.T's
		for i := 1; i <= nmb; i++ {
			mut := ov.under[i].(*inter.T)
			mut.ForEach(func(key string, off uint64) {
				if (off & tombstone) == 0 {
					fb.Insert(key, off)
				} else {
					fb.Delete(key, off&^tombstone)
				}
			})
		}
	})
}

func (ov *Overlay) WithMerged(fb Result, nmerged int) *Overlay {
	under := make([]tree, len(ov.under)-nmerged)
	under[0] = fb
	copy(under[1:], ov.under[1+nmerged:])
	return &Overlay{under: under}
}

//-------------------------------------------------------------------

// Save writes the Overlay's base fbtree to storage
// and returns the new fbtree (in an Overlay) to later pass to With
func (ov *Overlay) Save(flatten bool) Result {
	assert.That(ov.mut == nil)
	return ov.base().Save(flatten)
}

// WithSaved returns a new Overlay,
// combining the current state (ov) with the updated fbtree (in ov2)
func (ov *Overlay) WithSaved(r Result) *Overlay {
	under := make([]tree, len(ov.under))
	under[0] = r
	copy(under[1:], ov.under[1:])
	return &Overlay{under: under}
}

//-------------------------------------------------------------------

func (ov *Overlay) CheckFlat() {
	assert.Msg("not flat").That(len(ov.under) == 1)
}

func (ov *Overlay) CheckTnMerged(tn int) {
	for i := 1; i < len(ov.under); i++ {
		assert.That(ov.under[i].(*inter.T).TranNum != tn)
	}
}
