// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type iter = func() (string, uint64, bool)

// Overlay is the composite in-memory representation of an index
type Overlay struct {
	// bt is the stored base btree
	bt *btree.T
	// layers is a base ixbuf of merged but not persisted changes,
	// plus ixbuf's from completed but un-merged transactions
	layers []*ixbuf.T
	// mut is the per transaction mutable top ixbuf.T, nil if read-only
	mut *ixbuf.T
}

func NewOverlay(store *stor.Stor, is *ixkey.Spec) *Overlay {
	assert.That(is != nil)
	return &Overlay{bt: btree.CreateBtree(store, is),
		layers: []*ixbuf.T{{}}}
}

func OverlayFor(bt *btree.T) *Overlay {
	return &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
}

// Mutable returns a modifiable copy of an Overlay
func (ov *Overlay) Mutable() *Overlay {
	assert.That(ov.mut == nil)
	layers := make([]*ixbuf.T, len(ov.layers))
	copy(layers, ov.layers)
	assert.That(len(layers) >= 1)
	return &Overlay{bt: ov.bt, layers: layers, mut: &ixbuf.T{}}
}

func (ov *Overlay) GetIxspec() *ixkey.Spec {
	return ov.bt.GetIxspec()
}

func (ov *Overlay) SetIxspec(is *ixkey.Spec) {
	ov.bt.SetIxspec(is)
}

// Insert inserts into the mutable top ixbuf.T
func (ov *Overlay) Insert(key string, off uint64) {
	ov.mut.Insert(key, off)
}

// Delete either deletes the key/offset from the mutable ixbuf.T
// or inserts a tombstone into the mutable ixbuf.T.
func (ov *Overlay) Delete(key string, off uint64) {
	ov.mut.Delete(key, off)
}

func (ov *Overlay) Lookup(key string) uint64 {
	if ov.mut != nil {
		if off := ov.mut.Lookup(key); off != 0 {
			return off
		}
	}
	for i := len(ov.layers) - 1; i >= 0; i-- {
		if off := ov.layers[i].Lookup(key); off != 0 {
			return off
		}
	}
	if off := ov.bt.Lookup(key); off != 0 {
		return off
	}
	return 0
}

func (ov *Overlay) Check(fn func(uint64)) int {
	n, _, _ := ov.bt.Check(fn)
	return n
}

func (ov *Overlay) QuickCheck() {
	ov.bt.QuickCheck()
}

// Modified is used by info.Persist.
// It only looks at the base ixbuf (layer[0])
// which accumulates changes between persists
// via merging other per transaction ixbuf's.
func (ov *Overlay) Modified() bool {
	return ov.layers[0].Len() > 0
}

func (ov *Overlay) Iterator() *OverIter {
	callback := func(mc int) (int, []iterT) {
		if mc == -1 { // first time
			its := make([]iterT, 0, 2+len(ov.layers))
			its = append(its, ov.bt.Iterator())
			for _, ib := range ov.layers {
				its = append(its, ib.Iterator())
			}
			if ov.mut != nil {
				its = append(its, ov.mut.Iterator())
			}
			return 0, its
		}
		return mc, nil
	}
	return NewOverIter(callback)
}

//-------------------------------------------------------------------

func (ov *Overlay) StorSize() int {
	return ov.bt.StorSize()
}

func (ov *Overlay) Write(w *stor.Writer) {
	ov.bt.Write(w)
}

// ReadOverlay reads an Overlay from storage BUT without ixspec
func ReadOverlay(st *stor.Stor, r *stor.Reader) *Overlay {
	return &Overlay{bt: btree.Read(st, r), layers: []*ixbuf.T{{}}}
}

//-------------------------------------------------------------------

// UpdateWith combines the overlay result of a transaction
// with the latest overlay.
// The immutable part of ov was taken at the start of the transaction
// so it will be out of date.
// The checker ensures that the updates are independent.
func (ov *Overlay) UpdateWith(latest *Overlay) {
	ov.bt = latest.bt
	// reuse the new slice and overwrite ov.layers with the latest
	ov.layers = append(ov.layers[:0], latest.layers...)
	// add mut updates
	ov.layers = append(ov.layers, ov.mut)
	ov.mut = nil
	assert.That(len(ov.layers) >= 2)
}

//-------------------------------------------------------------------

type MergeResult = *ixbuf.T

// Merge merges the base ixbuf with one or more of the transaction ixbuf's
// to produce a new base ixbuf. It does not modify the original ixbuf's.
func (ov *Overlay) Merge(nmerge int) MergeResult {
	assert.That(ov.mut == nil)
	return ixbuf.Merge(ov.layers[:nmerge+1]...)
}

func (ov *Overlay) WithMerged(mr MergeResult, nmerged int) *Overlay {
	layers := make([]*ixbuf.T, len(ov.layers)-nmerged)
	layers[0] = mr
	copy(layers[1:], ov.layers[1+nmerged:])
	return &Overlay{bt: ov.bt, layers: layers}
}

//-------------------------------------------------------------------

type SaveResult = *btree.T

// Save updates the stored btree with the base ixbuf
// and returns the new btree to later pass to WithSaved
func (ov *Overlay) Save() SaveResult {
	assert.That(ov.mut == nil)
	return ov.bt.MergeAndSave(ov.layers[0].Iter())
}

// WithSaved returns a new Overlay,
// combining the current state (ov) with the updated btree (in ov2)
func (ov *Overlay) WithSaved(bt SaveResult) *Overlay {
	layers := make([]*ixbuf.T, len(ov.layers))
	layers[0] = &ixbuf.T{} // new empty base ixbuf
	copy(layers[1:], ov.layers[1:])
	return &Overlay{bt: bt, layers: layers}
}

//-------------------------------------------------------------------

func (ov *Overlay) CheckFlat() {
	assert.Msg("not flat").That(len(ov.layers) == 1)
}
