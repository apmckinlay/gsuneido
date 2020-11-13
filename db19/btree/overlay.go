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

// Overlay is an fbtree plus a base inter,
// plus overlay inters from un-merged transactions,
// plus a mutable inter within update transactions.
type Overlay struct {
	fb    *fbtree
	under []*inter.T
	// mut is the per transaction mutable top inter.T, nil if read-only
	mut *inter.T
}

func NewOverlay(store *stor.Stor, is *ixspec.T) *Overlay {
	assert.That(is != nil)
	return &Overlay{fb: CreateFbtree(store, is),
		under: []*inter.T{{}}}
}

// Mutable returns a modifiable copy of an Overlay
func (ov *Overlay) Mutable(tranNum int) *Overlay {
	assert.That(ov.mut == nil)
	under := make([]*inter.T, len(ov.under))
	copy(under, ov.under)
	assert.That(len(under) >= 1)
	return &Overlay{fb: ov.fb, under: under, mut: &inter.T{TranNum: tranNum}}
}

func (ov *Overlay) GetIxspec() *ixspec.T {
	return ov.fb.ixspec
}

func (ov *Overlay) SetIxspec(is *ixspec.T) {
	ov.fb.ixspec = is
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
	n, _, _ := ov.fb.check(fn)
	return n
}

func (ov *Overlay) QuickCheck() {
	ov.fb.quickCheck()
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
	return 5 + 1
}

func (ov *Overlay) Write(w *stor.Writer) {
	fb := ov.fb
	w.Put5(fb.root).Put1(fb.treeLevels)
}

// ReadOverlay reads an Overlay from storage BUT without ixspec
func ReadOverlay(st *stor.Stor, r *stor.Reader) *Overlay {
	root := r.Get5()
	treeLevels := r.Get1()
	return &Overlay{fb: OpenFbtree(st, root, treeLevels), under: []*inter.T{{}}}
}

//-------------------------------------------------------------------

// UpdateWith combines the overlay result of a transaction
// with the latest overlay.
// The immutable part of ov was taken at the start of the transaction
// so it will be out of date.
// The checker ensures that the updates are independent.
func (ov *Overlay) UpdateWith(latest *Overlay) {
	ov.fb = latest.fb
	// reuse the new slice and overwrite ov.under with the latest
	ov.under = append(ov.under[:0], latest.under...)
	// add mut updates
	ov.under = append(ov.under, ov.mut)
	ov.mut = nil
	assert.That(len(ov.under) >= 2)
}

//-------------------------------------------------------------------

type MergeResult = *inter.T

// Merge merges the base inter with one or more of the transaction inters
// to produce a new base inter. It does not modify the original inter's.
func (ov *Overlay) Merge(nmerge int) MergeResult {
	assert.That(ov.mut == nil)
	return inter.Merge(ov.under[:nmerge+1]...)
}

func (ov *Overlay) WithMerged(mr MergeResult, nmerged int) *Overlay {
	under := make([]*inter.T, len(ov.under)-nmerged)
	under[0] = mr
	copy(under[1:], ov.under[1+nmerged:])
	return &Overlay{fb: ov.fb, under: under}
}

//-------------------------------------------------------------------

type SaveResult = *fbtree

// Save updates the stored fbtree with the base inter
// and returns the new fbtree to later pass to WithSaved
func (ov *Overlay) Save() SaveResult {
	assert.That(ov.mut == nil)
	return ov.fb.MergeAndSave(ov.under[0].Iter(false))
}

// WithSaved returns a new Overlay,
// combining the current state (ov) with the updated fbtree (in ov2)
func (ov *Overlay) WithSaved(fb SaveResult) *Overlay {
	under := make([]*inter.T, len(ov.under))
	under[0] = &inter.T{} // new empty base inter
	copy(under[1:], ov.under[1:])
	return &Overlay{fb: fb, under: under}
}

//-------------------------------------------------------------------

func (ov *Overlay) CheckFlat() {
	assert.Msg("not flat").That(len(ov.under) == 1)
}

func (ov *Overlay) CheckTnMerged(tn int) {
	for i := 1; i < len(ov.under); i++ {
		assert.That(ov.under[i].TranNum != tn)
	}
}
