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
	tran    oiTran
	overlay *Overlay
	rng     Range
	table   string
	// We need to keep our own curKey/Off independent of the source iterators
	// because new source iterators may be returned by the callback.
	curKey string
	iters  []iterT
	iIndex int
	curOff uint64
	// curIter is the iterator containing the current item = iters[curIter]
	curIter int
	state
	lastDir dir
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

func (oi *OverIter) Eof() bool {
	return oi.state == eof
}

func (oi *OverIter) Cur() (string, uint64) {
	assert.Msg("OverIter Cur eof").That(oi.state != eof && oi.state != rewound)
	assert.Msg("OverIter Cur deleted").That(oi.curOff&ixbuf.Delete == 0)
	return oi.curKey, oi.curOff
}

func (oi *OverIter) Range(rng Range) {
	oi.rng = rng
	oi.state = rewound
	for _, it := range oi.iters {
		it.Range(rng)
	}
}

func (oi *OverIter) Print() {
	oi.overlay.Print()
}

// Next -------------------------------------------------------------

func (oi *OverIter) Next(t oiTran) {
	// NOTE: keep this code in sync with Prev
	if oi.state == eof {
		return // stick at eof
	}
	modified := oi.update(t)
	if oi.state == rewound {
		oi.all(iterT.Next)
		oi.state = front
	} else {
		oi.modNext(modified)
	}
	lastState := oi.state
	oi.curIter, oi.curKey, oi.curOff = oi.minIter()
	if oi.curIter == -1 {
		oi.state = eof
	}
	oi.lastDir = next
	if lastState == front {
		if oi.state == eof {
			oi.tran.Read(oi.table, oi.iIndex, oi.rng.Org, oi.rng.End)
		} else {
			oi.tran.Read(oi.table, oi.iIndex, oi.rng.Org, oi.curKey)
		}
	}
}

func (oi *OverIter) update(t oiTran) bool {
	ov := t.GetIndexI(oi.table, oi.iIndex)
	if ov == oi.overlay {
		return false
	}
	oi.tran = t
	oi.newIters(ov)
	return true
}

func (oi *OverIter) newIters(ov *Overlay) {
	its := make([]iterT, 0, 2+len(ov.layers))
	its = append(its, ov.bt.Iterator())
	for _, ib := range ov.layers {
		its = append(its, ib.Iterator())
	}
	if ov.mut != nil {
		its = append(its, ov.mut.Iterator())
	}
	for _, it := range its {
		it.Range(oi.rng)
	}
	oi.iters = its
	oi.overlay = ov
}

func (oi *OverIter) all(fn func(it iterT)) {
	for _, it := range oi.iters {
		fn(it)
	}
}

func (oi *OverIter) modNext(modified bool) {
	// NOTE: keep this code in sync with modPrev
	for _, it := range oi.iters {
		if modified || it.Modified() {
			it.Seek(oi.curKey)
			if !it.Eof() {
				key, _ := it.Cur()
				if key <= oi.curKey {
					it.Next()
				}
			}
		} else if oi.lastDir != next {
			if it.Eof() {
				it.Rewind()
			}
			it.Next()
		} else if atKey(it, oi.curKey) {
			it.Next()
		}
	}
}

func atKey(it iterator.T, key string) bool {
	if it.Eof() {
		return false
	}
	itkey, _ := it.Cur()
	return itkey == key
}

// minIter finds the the minimum current key
func (oi *OverIter) minIter() (int, string, uint64) {
	// NOTE: keep this code in sync with maxIter
outer:
	for {
		itMin := -1
		var keyMin string
		var offMin uint64
		for i, it := range oi.iters {
			if it.Eof() {
				continue
			}
			key, off := it.Cur()
			if itMin == -1 || key < keyMin {
				itMin = i
				keyMin = key
				offMin = off
			} else if key == keyMin {
				off = ixbuf.Combine(offMin, off)
				if off == 0 || off&ixbuf.Delete != 0 {
					oi.ifKey(i, key, iterator.T.Next)
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

func (oi *OverIter) ifKey(i int, key string, fn func(iterator.T)) {
	for j := 0; j < i; j++ {
		itj := oi.iters[j]
		if !itj.Eof() {
			k, _ := itj.Cur()
			if k == key {
				fn(itj)
			}
		}
	}
}

// Prev -------------------------------------------------------------

func (oi *OverIter) Prev(t oiTran) {
	// NOTE: keep this code in sync with Next
	if oi.state == eof {
		return // stick at eof
	}
	modified := oi.update(t)
	if oi.state == rewound {
		oi.all(iterT.Prev)
		oi.state = back
	} else {
		oi.modPrev(modified)
	}
	lastState := oi.state
	oi.curIter, oi.curKey, oi.curOff = oi.maxIter()
	if oi.curIter == -1 {
		oi.state = eof
	}
	oi.lastDir = prev
	if lastState == back {
		if oi.state == eof {
			oi.tran.Read(oi.table, oi.iIndex, oi.rng.Org, oi.rng.End)
		} else {
			oi.tran.Read(oi.table, oi.iIndex, oi.curKey, oi.rng.End)
		}
	}
}

func (oi *OverIter) modPrev(modified bool) {
	// NOTE: keep this code in sync with modNext
	for _, it := range oi.iters {
		if modified || it.Modified() {
			it.SeekAll(oi.curKey)
			if !it.Eof() {
				key, _ := it.Cur()
				if key >= oi.curKey {
					it.Prev()
				}
			}
		} else if oi.lastDir != prev {
			if it.Eof() {
				it.Rewind()
			}
			it.Prev()
		} else if atKey(it, oi.curKey) {
			it.Prev()
		}
	}
}

// maxIter finds the maximum current key
func (oi *OverIter) maxIter() (int, string, uint64) {
	// NOTE: keep this code in sync with minIter
outer:
	for {
		itMax := -1
		var keyMax string
		var offMax uint64
		for i, it := range oi.iters {
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
				if off == 0 || off&ixbuf.Delete != 0 {
					oi.ifKey(i, key, iterator.T.Prev)
					// delete so skip
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

func (oi *OverIter) Rewind() {
	oi.all(iterT.Rewind)
	oi.state = rewound
	oi.curIter = -1
	oi.curKey = ""
	oi.curOff = 0
}
