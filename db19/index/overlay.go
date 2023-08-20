// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"fmt"

	"slices"

	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Overlay is the composite in-memory representation of an index
type Overlay struct {
	// bt is the stored base btree
	bt *btree.T
	// mut is the per transaction mutable top ixbuf.T, nil if read-only
	mut *ixbuf.T
	// layers is a base ixbuf of merged but not persisted changes,
	// plus ixbuf's from completed but un-merged transactions
	layers []*ixbuf.T
}

func (ov *Overlay) Cksum() uint32 {
	return ov.bt.Cksum()
}

func NewOverlay(store *stor.Stor, is *ixkey.Spec) *Overlay {
	assert.That(is != nil)
	return &Overlay{bt: btree.CreateBtree(store, is),
		layers: []*ixbuf.T{{}}} // single base layer
}

// OverlayStub is for tests
func OverlayStub() *Overlay {
	return &Overlay{bt: &btree.T{}}
}

func OverlayFor(bt *btree.T) *Overlay {
	return &Overlay{bt: bt, layers: []*ixbuf.T{{}}} // single base layer
}

func OverlayForN(bt *btree.T, nlayers int) *Overlay {
	layers := make([]*ixbuf.T, nlayers)
	for i := range layers {
		layers[i] = &ixbuf.T{}
	}
	return &Overlay{bt: bt, layers: layers}
}

func (ov *Overlay) Nlayers() int {
	return len(ov.layers)
}

func (ov *Overlay) BtreeLevels() int {
	return ov.bt.TreeLevels() + 1
}

// Mutable returns a modifiable copy of an Overlay
func (ov *Overlay) Mutable() *Overlay {
	assert.That(ov.mut == nil)
	assert.That(len(ov.layers) > 0)
	assert.That(ov.layers[0] != nil)
	layers := slices.Clone(ov.layers)
	assert.That(len(layers) >= 1)
	return &Overlay{bt: ov.bt, layers: layers, mut: &ixbuf.T{}}
}

// Copy is for debugging
func (ov *Overlay) Copy() *Overlay {
	assert.That(ov.mut == nil)
	layers := slices.Clone(ov.layers)
	assert.That(len(layers) >= 1)
	return &Overlay{bt: ov.bt, layers: layers}
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

func (ov *Overlay) Update(key string, off uint64) {
	ov.mut.Update(key, off)
}

// Lookup returns the offset of the record specified by the key
// or 0 if it's not found.
// It handles the Delete bit and removes the Update bit.
func (ov *Overlay) Lookup(key string) uint64 {
	if ov.mut != nil {
		if off := ov.mut.Lookup(key); off != 0 {
			if off&ixbuf.Delete != 0 {
				return 0 // deleted
			}
			return off &^ ixbuf.Update
		}
	}
	for i := len(ov.layers) - 1; i >= 0; i-- {
		if off := ov.layers[i].Lookup(key); off != 0 {
			if off&ixbuf.Delete != 0 {
				return 0 // deleted
			}
			return off &^ ixbuf.Update
		}
	}
	if off := ov.bt.Lookup(key); off != 0 {
		return off
	}
	return 0
}

func (ov *Overlay) RangeFrac(org, end string) float32 {
	return ov.bt.RangeFrac(org, end)
}

func (ov *Overlay) Check(fn func(uint64)) int {
	n, _, _ := ov.bt.Check(fn)
	return n
}

func (ov *Overlay) QuickCheck() {
	ov.bt.QuickCheck()
}

func (ov *Overlay) Stats() btree.Stats {
	return ov.bt.Stats()
}

// Modified is used by info.Persist.
// It only looks at the base ixbuf (layer[0])
// which accumulates changes between persists
// via merging other per transaction ixbuf's.
func (ov *Overlay) Modified() bool {
	return ov.layers[0].Len() > 0
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
// with the latest overlay. It is called by Meta.LayeredOnto.
// The immutable part of ov was taken at the start of the transaction
// so it will be out of date.
// The checker ensures that the updates are independent.
func (ov *Overlay) UpdateWith(latest *Overlay) {
	ov.bt = latest.bt
	// reuse the new slice and overwrite ov.layers with the latest
	ov.layers = append(ov.layers[:0], latest.layers...) // copy
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

// WithMerged is called by Meta.ApplyMerge
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

func (ov *Overlay) CheckMerged() {
	if len(ov.layers) != 1 && ov.layers[0].Len() != 0 {
		panic("index not merged")
	}
}

func (ov *Overlay) Print() {
	fmt.Println("btree")
	ov.bt.Print()
	for i, ixb := range ov.layers {
		fmt.Println("layer", i)
		ixb.Print()
	}
	if ov.mut != nil {
		fmt.Println("mut")
		ov.mut.Print()
	}
}
