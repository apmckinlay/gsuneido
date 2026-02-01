// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"fmt"
	"strings"

	btree1 "github.com/apmckinlay/gsuneido/db19/index/btree"
	btree3 "github.com/apmckinlay/gsuneido/db19/index/btree3"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/slc"
)

// Overlay is the composite in-memory representation of an index
// @immutable
type Overlay struct {
	// bt is the stored base btree (immutable)
	bt iface.Btree
	// mut is the per transaction mutable top ixbuf.T, nil if read-only
	mut *ixbuf.T
	// layers is: (immutable)
	// - a base ixbuf of merged but not persisted changes,
	// - plus ixbuf's from completed but un-merged transactions
	layers []*ixbuf.T
}

func (ov *Overlay) Cksum() uint32 {
	return ov.bt.Cksum()
}

// NewOverlay is only used by big_test.go
func NewOverlay(store *stor.Stor, is *ixkey.Spec) *Overlay {
	assert.That(is != nil)
	return &Overlay{bt: btree3.CreateBtree(store, is),
		layers: []*ixbuf.T{{}}} // single base layer
}

// OverlayStub is for tests
func OverlayStub() *Overlay {
	return &Overlay{bt: &btree3.T{}}
}

func OverlayFor(bt iface.Btree) *Overlay {
	return &Overlay{bt: bt, layers: []*ixbuf.T{{}}} // single base layer
}

func OverlayForN(bt iface.Btree, nlayers int) *Overlay {
	layers := make([]*ixbuf.T, nlayers)
	for i := range layers {
		layers[i] = &ixbuf.T{}
	}
	return &Overlay{bt: bt, layers: layers}
}

func (ov *Overlay) Nlayers() int {
	return len(ov.layers)
}

// BtreeLevels includes leaf level as well as tree levels
func (ov *Overlay) BtreeLevels() int {
	return ov.bt.TreeLevels() + 1
}

// Mutable returns a modifiable copy of an Overlay
func (ov *Overlay) Mutable() *Overlay {
	assert.That(ov.mut == nil && len(ov.layers) > 0 && ov.layers[0] != nil)
	return &Overlay{bt: ov.bt, layers: ov.layers, mut: &ixbuf.T{}}
}

// Copy is for debugging
func (ov *Overlay) Copy() *Overlay {
	assert.That(ov.mut == nil)
	layers := slc.Clone(ov.layers)
	assert.That(len(layers) >= 1)
	return &Overlay{bt: ov.bt, layers: layers}
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
func (ov *Overlay) Delete(key string, off uint64) uint64 {
	return ov.mut.Delete(key, off)
}

func (ov *Overlay) Update(key string, off uint64) uint64 {
	return ov.mut.Update(key, off)
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

func (ov *Overlay) RangeFrac(org, end string, nrecs int) float64 {
	return ov.bt.RangeFrac(org, end, nrecs)
}

// CheckBtree applies a function to each entry in the btree.
// WARNING: it ignores other layers.
func (ov *Overlay) CheckBtree(fn any) int {
	n, _, _ := ov.bt.Check(fn)
	return n
}

func (ov *Overlay) BtreeIter() iface.Iter {
	return ov.bt.Iterator()
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

//-------------------------------------------------------------------

func (ov *Overlay) StorSize() int {
	return 6 // 5 for root offset, 1 for tree levels
}

func (ov *Overlay) Write(w *stor.Writer) {
	ov.bt.Write(w)
}

// ReadOverlay reads an Overlay from storage BUT without ixspec
func ReadOverlay(st *stor.Stor, r *stor.Reader, nrows int) *Overlay {
	var bt iface.Btree
	if st.OldVer {
		bt = btree1.Read(st, r)
	} else {
		bt = btree3.Read(st, r, nrows)
	}
	return &Overlay{bt: bt, layers: []*ixbuf.T{{}}}
}

//-------------------------------------------------------------------

// UpdateWith combines the overlay result of a transaction
// with the latest overlay. It is called by Meta.LayeredOnto.
// The immutable part of ov was taken at the start of the transaction
// so it will be out of date.
// The checker ensures that the updates are independent.
func (ov *Overlay) UpdateWith(latest *Overlay) {
	ov.bt = latest.bt                           // @allow-mutate
	ov.layers = slc.With(latest.layers, ov.mut) // @allow-mutate
	ov.mut = nil                                // @allow-mutate
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

// WithMerged replaces nmerge layers with the merge result
func (ov *Overlay) WithMerged(mr MergeResult, nmerged int) *Overlay {
	layers := make([]*ixbuf.T, len(ov.layers)-nmerged)
	layers[0] = mr
	copy(layers[1:], ov.layers[1+nmerged:])
	return &Overlay{bt: ov.bt, layers: layers}
}

//-------------------------------------------------------------------

// Save updates the stored btree with the base ixbuf
// and returns the new btree to later pass to WithSaved
func (ov *Overlay) Save() iface.Btree {
	assert.That(ov.mut == nil)
	return ov.bt.MergeAndSave(ov.layers[0].Iter())
}

// WithSaved returns a new Overlay,
// combining the current state (ov) with the updated btree (in ov2)
func (ov *Overlay) WithSaved(bt iface.Btree) *Overlay {
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

// func (ov *Overlay) Print() {
// 	fmt.Println("btree")
// 	ov.bt.Print()
// 	for i, ixb := range ov.layers {
// 		fmt.Println("layer", i)
// 		ixb.Print()
// 	}
// 	if ov.mut != nil {
// 		fmt.Println("mut")
// 		ov.mut.Print()
// 	}
// }

// String returns a string representation of the overlay
// WITHOUT the btree
func (ov *Overlay) String() string {
	var sb strings.Builder
	for i, ixb := range ov.layers {
		fmt.Fprintln(&sb, "layer", i, ":", ixb)
	}
	if ov.mut != nil {
		sb.WriteString("mut: ")
		sb.WriteString(ov.mut.String())
		ov.mut.Print()
	}
	return sb.String()
}

// Check verifies an Overlay by:
// - using ixbuf.Check to check the ixbuf layers
// - checking that layer entries for each key are in sequence: add, possible updates, possible delete
func (ov *Overlay) Check() map[string]uint64 {
	// Check each ixbuf layer individually
	for _, layer := range ov.layers {
		layer.Check()
	}
	if ov.mut != nil {
		ov.mut.Check()
	}

	// Track current offset for each key (0 = doesn't exist)
	keyOffsets := make(map[string]uint64)

	// Process layers in order: layers[0] to layers[n], then mut
	for layerIdx, layer := range ov.layers {
		checkLayer(layer, layerIdx, keyOffsets)
	}
	if ov.mut != nil {
		checkLayer(ov.mut, len(ov.layers), keyOffsets)
	}
	return keyOffsets
}

func checkLayer(layer *ixbuf.T, layerIdx int, keyOffsets map[string]uint64) {
	iter := layer.Iter()
	for {
		key, off, ok := iter()
		if !ok {
			break
		}
		actualOffset := off & ixbuf.Mask
		currentOffset := keyOffsets[key]
		if off&ixbuf.Delete != 0 {
			if currentOffset == 0 {
				panic(fmt.Sprintf("delete of non-existent key %q in layer %d", key, layerIdx))
			}
			if actualOffset < currentOffset {
				panic(fmt.Sprintf("delete offset mismatch for key %q in layer %d: expected >= %d, got %d", key, layerIdx, currentOffset, actualOffset))
			}
			keyOffsets[key] = 0
		} else if off&ixbuf.Update != 0 {
			if currentOffset == 0 {
				panic(fmt.Sprintf("update of non-existent key %q in layer %d", key, layerIdx))
			}
			keyOffsets[key] = actualOffset
		} else {
			if currentOffset != 0 {
				panic(fmt.Sprintf("add of existing key %q in layer %d", key, layerIdx))
			}
			keyOffsets[key] = actualOffset
		}
	}
}
