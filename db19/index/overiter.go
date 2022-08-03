// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type iterT = iterator.T

type Range = iterator.Range

// OverIter is a Suneido style iterator
// that merges several other Suneido style iterators.
//
// OverIter also tracks read ranges for transaction conflict checking.
type OverIter struct {
	table  string
	iIndex int
	iters  []iterT
	// We need to keep our own curKey/Off independent of the source iterators
	// because new source iterators may be returned by the callback.
	curKey string
	curOff uint64
	// curIter is the iterator containing the current item = iters[curIter]
	curIter int
	state
	lastDir dir
	rng     Range
	tran    oiTran
	overlay *Overlay
}

type state byte

const (
	rewound state = iota
	front
	back
	eof
)

type dir int8

const (
	next dir = +1
	prev dir = -1
)

type oiTran interface {
	GetIndexI(table string, iIndex int) *Overlay
	Read(table string, iIndex int, from, to string)
}

func NewOverIter(table string, iIndex int) *OverIter {
	return &OverIter{table: table, iIndex: iIndex, rng: iterator.All}
}

func (mi *OverIter) Eof() bool {
	return mi.state == eof
}

func (mi *OverIter) Cur() (string, uint64) {
	assert.Msg("OverIter Cur eof").That(mi.state != eof && mi.state != rewound)
	assert.Msg("OverIter Cur deleted").That(mi.curOff&ixbuf.Delete == 0)
	return mi.curKey, mi.curOff
}

func (mi *OverIter) Range(rng Range) { //TODO org & end params
	mi.rng = rng
	mi.state = rewound
	for _, it := range mi.iters {
		it.Range(rng)
	}
}

func (mi *OverIter) Next(t oiTran) {
	// NOTE: keep this code in sync with Prev
	if mi.state == eof {
		return // stick at eof
	}
	if mi.state == rewound {
		if t != mi.tran {
			// some of the iterators may still be valid
			// but simplest is just to get fresh iterators
			ov := t.GetIndexI(mi.table, mi.iIndex)
			mi.newIters(ov)
			mi.tran = t
		}
		mi.all(iterT.Next)
		mi.state = front
	} else {
		mi.modNext(t)
	}
	lastState := mi.state
	mi.curIter, mi.curKey, mi.curOff = mi.minIter()
	if mi.curIter == -1 {
		mi.state = eof
	}
	mi.lastDir = next
	if lastState == front {
		if mi.state == eof {
			mi.tran.Read(mi.table, mi.iIndex, mi.rng.Org, mi.rng.End)
		} else {
			mi.tran.Read(mi.table, mi.iIndex, mi.rng.Org, mi.curKey)
		}
	}
}

func (mi *OverIter) newIters(ov *Overlay) {
	its := make([]iterT, 0, 2+len(ov.layers))
	its = append(its, ov.bt.Iterator())
	for _, ib := range ov.layers {
		its = append(its, ib.Iterator())
	}
	if ov.mut != nil {
		its = append(its, ov.mut.Iterator())
	}
	for _, it := range its {
		it.Range(mi.rng)
	}
	mi.iters = its
	mi.overlay = ov
}

func (mi *OverIter) all(fn func(it iterT)) {
	for _, it := range mi.iters {
		fn(it)
	}
}

func (mi *OverIter) modNext(t oiTran) {
	// NOTE: keep this code in sync with modPrev
	modified := false
	if t != mi.tran {
		ov := t.GetIndexI(mi.table, mi.iIndex)
		if !mi.keepIters(ov) {
			mi.newIters(ov)
			modified = true
		}
		mi.tran = t
	}
	for i, it := range mi.iters {
		if modified || it.Modified() {
			it.Seek(mi.curKey)
			if !it.Eof() {
				key, _ := it.Cur()
				if key <= mi.curKey {
					it.Next()
				}
			}
		} else if mi.lastDir != next {
			if it.Eof() {
				it.Rewind()
			}
			it.Next()
		} else if i == mi.curIter {
			it.Next()
		}
	}
}

func (mi *OverIter) keepIters(ov *Overlay) bool {
	if ov.bt != mi.overlay.bt ||
		len(ov.layers) != len(mi.overlay.layers) ||
		(ov.mut == nil) != (mi.overlay.mut == nil) {
		return false
	}
	for i, layer := range mi.overlay.layers {
		if layer != ov.layers[i] {
			return false
		}
	}
	return true
}

// minIter finds the the minimum current key
func (mi *OverIter) minIter() (int, string, uint64) {
	// NOTE: keep this code in sync with maxIter
outer:
	for {
		itMin := -1
		var keyMin string
		var offMin uint64
		for i, it := range mi.iters {
			if it.Eof() {
				continue
			}
			key, off := it.Cur()
			if itMin == -1 || key < keyMin {
				itMin = i
				keyMin = key
				offMin = off
				// may get update first
				// e.g. add, update => add, but then restart gets update
			} else if key == keyMin {
				off = ixbuf.Combine(offMin, off)
				mi.iters[itMin].Next()
				if off == 0 || off&ixbuf.Delete != 0 {
					// delete so skip
					// may not be the final minimum, but still need to skip
					it.Next()
					continue outer // restart
				}
				itMin = i
				offMin = off
			}
		}
		return itMin, keyMin, offMin &^ ixbuf.Update
	}
}

func (mi *OverIter) Prev(t oiTran) {
	// NOTE: keep this code in sync with Next
	if mi.state == eof {
		return // stick at eof
	}
	if mi.state == rewound {
		if t != mi.tran {
			// some of the iterators may still be valid
			// but simplest just to get fresh iterators
			ov := t.GetIndexI(mi.table, mi.iIndex)
			mi.newIters(ov)
			mi.tran = t
		}
		mi.all(iterT.Prev)
		mi.state = back
	} else {
		mi.modPrev(t)
	}
	lastState := mi.state
	mi.curIter, mi.curKey, mi.curOff = mi.maxIter()
	if mi.curIter == -1 {
		mi.state = eof
	}
	mi.lastDir = prev
	if lastState == back {
		if mi.state == eof {
			mi.tran.Read(mi.table, mi.iIndex, mi.rng.Org, mi.rng.End)
		} else {
			mi.tran.Read(mi.table, mi.iIndex, mi.curKey, mi.rng.End)
		}
	}
}

func (mi *OverIter) modPrev(t oiTran) {
	// NOTE: keep this code in sync with modNext
	modified := false
	if t != mi.tran {
		ov := t.GetIndexI(mi.table, mi.iIndex)
		if !mi.keepIters(ov) {
			mi.newIters(ov)
			modified = true
		}
		mi.tran = t
	}
	for i, it := range mi.iters {
		if modified || it.Modified() {
			it.SeekAll(mi.curKey)
			if !it.Eof() {
				key, _ := it.Cur()
				if key >= mi.curKey {
					it.Prev()
				}
			}
		} else if mi.lastDir != prev {
			if it.Eof() {
				it.Rewind()
			}
			it.Prev()
		} else if i == mi.curIter {
			it.Prev()
		}
	}
}

// maxIter finds the maximum current key
func (mi *OverIter) maxIter() (int, string, uint64) {
	// NOTE: keep this code in sync with minIter
outer:
	for {
		itMax := -1
		var keyMax string
		var offMax uint64
		for i, it := range mi.iters {
			if it.Eof() {
				continue
			}
			key, off := it.Cur()
			if itMax == -1 || key > keyMax {
				itMax = i
				keyMax = key
				offMax = off
			} else if key == keyMax {
				off = ixbuf.Combine(offMax, off)
				mi.iters[itMax].Prev()
				if off == 0 || off&ixbuf.Delete != 0 {
					// add,delete so skip
					// may not be the final minimum, but still need to skip
					it.Prev()
					continue outer
				}
				itMax = i
				offMax = off
			}
		}
		return itMax, keyMax, offMax &^ ixbuf.Update
	}
}

func (mi *OverIter) Rewind() {
	mi.all(iterT.Rewind)
	mi.state = rewound
	mi.curIter = -1
	mi.curKey = ""
	mi.curOff = 0
}
