// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
)

type iterT = iterator.T

type Range = iterator.Range

// OverIter is a Suneido style iterator
// that merges several other Suneido style iterators.
//
// OverIter also tracks read ranges for transaction conflict checking.
type OverIter struct {
	overlay *Overlay
	rng     Range
	table   string
	// We need to keep our own curKey/Off independent of the source iterators
	// because new source iterators may be returned by the callback.
	curKey string
	iters  []iterT
	iIndex int
	curOff uint64
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

func (s state) String() string {
	switch s {
	case rewound:
		return "rewound"
	case front:
		return "front"
	case back:
		return "back"
	case eof:
		return "eof"
	default:
		panic("unknown state")
	}
}

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

func (oi *OverIter) HasCur() bool {
	return oi.state != eof && oi.state != rewound
}

func (oi *OverIter) Cur() (string, uint64) {
	if oi.state == eof {
		panic("OverIter Cur eof")
	}
	if oi.state == rewound {
		panic("OverIter Cur rewound")
	}
	if oi.curOff&ixbuf.Delete != 0 {
		panic("OverIter Cur deleted")
	}
	return oi.curKey, oi.curOff
}

func (oi *OverIter) Range(rng Range) {
	oi.rng = rng
	oi.state = rewound
	for _, it := range oi.iters {
		it.Range(rng)
	}
}

// Next -------------------------------------------------------------

func (oi *OverIter) Next(t oiTran) {
	// NOTE: keep this code in sync with Prev
	if oi.state == eof {
		return // stick at eof
	}

	prevKey := oi.curKey

	modified := oi.update(t)
	if oi.state == rewound {
		oi.all(iterT.Next)
		oi.state = front
		prevKey = oi.rng.Org
	} else {
		oi.modNext(modified)
	}
	var found bool
	found, oi.curKey, oi.curOff = oi.minIter()
	if found {
		t.Read(oi.table, oi.iIndex, prevKey, oi.curKey)
	} else {
		oi.state = eof
		t.Read(oi.table, oi.iIndex, prevKey, oi.rng.End)
	}
	oi.lastDir = next
}

func (oi *OverIter) update(t oiTran) bool {
	ov := t.GetIndexI(oi.table, oi.iIndex)
	if ov == oi.overlay {
		return false
	}
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
	for i, it := range oi.iters {
		if it.Modified() && i != len(oi.iters)-1 {
			panic("OverIter modNext modified not last")
		}
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
func (oi *OverIter) minIter() (bool, string, uint64) {
	// NOTE: keep this code in sync with maxIter
	for {
		var keyMin string
		var result uint64
		found := false

		// find minimum key and its most recent operation
		for _, it := range oi.iters {
			if it.Eof() {
				continue
			}
			key, off := it.Cur()
			if !found || key < keyMin { // new minimum found
				keyMin = key
				found = true
			}
			if key == keyMin { // track the most recent operation for keyMin
				if off&ixbuf.Delete != 0 {
					result = 0
				} else {
					result = off
				}
			}
		}
		if !found {
			return false, "", 0
		}
		if result != 0 {
			return true, keyMin, result &^ ixbuf.Update
		}

		// Skip this key - advance all iterators that have it
		for _, it := range oi.iters {
			if !it.Eof() {
				key, _ := it.Cur()
				if key == keyMin {
					it.Next()
				}
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

	prevKey := oi.curKey

	modified := oi.update(t)
	if oi.state == rewound {
		oi.all(iterT.Prev)
		oi.state = back
		prevKey = oi.rng.End
	} else {
		oi.modPrev(modified)
	}
	var found bool
	found, oi.curKey, oi.curOff = oi.maxIter()
	if found {
		t.Read(oi.table, oi.iIndex, oi.curKey, prevKey)
	} else {
		oi.state = eof
		t.Read(oi.table, oi.iIndex, oi.rng.Org, prevKey)
	}
	oi.lastDir = prev
}

func (oi *OverIter) modPrev(modified bool) {
	// NOTE: keep this code in sync with modNext
	for i, it := range oi.iters {
		if it.Modified() && i != len(oi.iters)-1 {
			panic("OverIter modPrev modified not last")
		}
		if modified || it.Modified() {
			it.Seek(oi.curKey) // finds first key >= curKey
			// Seek is upper bound but Prev needs lower bound
			if it.Eof() {
				// Seek hit EOF, all keys < curKey
				// For reverse iteration, position at the last key
				it.Rewind()
				it.Prev() // go to last key in range
			} else {
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
func (oi *OverIter) maxIter() (bool, string, uint64) {
	// NOTE: keep this code in sync with minIter
	for {
		var keyMax string
		var result uint64
		found := false

		// find maximum key and its most recent operation
		for _, it := range oi.iters {
			if it.Eof() {
				continue
			}
			key, off := it.Cur()
			if !found || key > keyMax { // new maximum found
				keyMax = key
				found = true
			}
			if key == keyMax { // track the most recent operation for keyMax
				if off&ixbuf.Delete != 0 {
					result = 0
				} else {
					result = off
				}
			}
		}
		if !found {
			return false, "", 0
		}
		if result != 0 {
			return true, keyMax, result &^ ixbuf.Update
		}

		// Skip this key - advance all iterators that have it
		for _, it := range oi.iters {
			if !it.Eof() {
				key, _ := it.Cur()
				if key == keyMax {
					it.Prev()
				}
			}
		}
	}
}

func (oi *OverIter) Rewind() {
	oi.all(iterT.Rewind)
	oi.state = rewound
	oi.curKey = ""
	oi.curOff = 0
}

func (oi *OverIter) String() string {
	var sb strings.Builder
	sb.WriteString("OverIter ")
	sb.WriteString(oi.state.String())
	sb.WriteString(" ")
	sb.WriteString(oi.curKey)
	sb.WriteString(ixbuf.OffString(oi.curOff))
	sb.WriteString(" :")
	for _, it := range oi.iters {
		if it.Eof() {
			sb.WriteString(" eof")
		} else {
			key, off := it.Cur()
			sb.WriteString(" ")
			sb.WriteString(key)
			sb.WriteString(ixbuf.OffString(off))
		}
	}
	return sb.String()
}

func (oi *OverIter) Check() {
	if !oi.HasCur() {
		return
	}
	curKey, curOff := oi.Cur()
	off := uint64(0)
	for _, it := range oi.iters {
		if it.HasCur() {
			itkey, itoff := it.Cur()
			if itkey == curKey {
				if off == 0 {
					off = itoff
				} else {
					off, _ = ixbuf.Combine(off, itoff)
				}
			}
		}
	}
	if off != curOff {
		panic("OverIter Check: offset mismatch")
	}
}
