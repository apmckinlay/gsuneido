// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type iterT = iface.Iter

type Range = iface.Range

// OverIter is a Suneido style iterator
// that merges several other Suneido style iterators.
//
// OverIter also tracks read ranges for transaction conflict checking.
type OverIter struct {
	overlay *Overlay
	rng     Range
	skipRng Range // suffix range for skip-scan mode
	table   string
	// We need to keep our own curKey/Off independent of the source iterators
	// because new source iterators may be returned by the callback.
	curKey string
	iters  []iterT
	iIndex int
	curOff uint64
	state
	lastDir dir
	// fastIdx is the index of the winning iterator from the last minIter/maxIter.
	// -1 means no fast path is available.
	// secondMin is the minimum key of all non-winning iterators (ixkey.Max if none),
	// used by the Next fast path to detect when re-merging is needed.
	fastIdx   int
	secondMin string
	secondMax string
	// number of leading fields treated as prefix in skip-scan mode
	skipStart int
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
	Num() int
}

func NewOverIter(table string, iIndex int) *OverIter {
	return &OverIter{table: table, iIndex: iIndex, rng: iface.All, fastIdx: -1}
}

func (oi *OverIter) Eof() bool {
	return oi.state == eof
}

func (oi *OverIter) HasCur() bool {
	return oi.state != eof && oi.state != rewound
}

func (oi *OverIter) Cur() (string, uint64) {
	oi.checkHasCur()
	return oi.curKey, oi.curOff
}

func (oi *OverIter) CurOff() uint64 {
	oi.checkHasCur()
	return oi.curOff
}

func (oi *OverIter) checkHasCur() {
	if oi.state == eof {
		panic("OverIter Cur eof")
	}
	if oi.state == rewound {
		panic("OverIter Cur rewound")
	}
	if oi.curOff&ixbuf.Delete != 0 {
		panic("OverIter Cur deleted")
	}
}

func (oi *OverIter) Range(rng Range) {
	oi.rng = rng
	oi.skipStart = 0
	oi.state = rewound
	for _, it := range oi.iters {
		it.Range(rng)
	}
}

// SkipScan enables skip-scan mode.
// prefixRng restricts visited prefix groups; iface.All means unrestricted.
// suffixRng applies to suffix fields (excluding prefix fields).
func (oi *OverIter) SkipScan(prefixRng Range, suffixRng Range, skipStart int) {
	assert.That(skipStart > 0)
	oi.rng = prefixRng
	oi.skipRng = suffixRng
	oi.skipStart = skipStart
	oi.state = rewound
	for _, it := range oi.iters {
		it.SkipScan(prefixRng, suffixRng, skipStart)
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
		oi.fastIdx = -1
	} else if oi.canFast(modified, next) && oi.fastNext(t, prevKey) {
		return
	} else {
		oi.fastIdx = -1
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
		if oi.skipStart != 0 {
			it.SkipScan(oi.rng, oi.skipRng, oi.skipStart)
		} else {
			it.Range(oi.rng)
		}
	}
	oi.iters = its
	oi.overlay = ov
}

func (oi *OverIter) canFast(modified bool, d dir) bool {
	return !modified &&
		oi.lastDir == d &&
		oi.fastIdx >= 0 &&
		!oi.iters[len(oi.iters)-1].Modified()
}

func (oi *OverIter) fastNext(t oiTran, prevKey string) bool {
	it := oi.iters[oi.fastIdx]
	it.Next()
	if !it.Eof() && it.Key() < oi.secondMin {
		// still the sole winner, skip full scan.
		// Safe to use Cur() without checking Delete: any tombstone in this iter
		// has a companion in another iter, creating a tie in minIter which sets
		// secondMin <= tombstone key, so the fast path is never taken past a tombstone.
		oi.curKey, oi.curOff = it.Cur()
		t.Read(oi.table, oi.iIndex, prevKey, oi.curKey)
		oi.lastDir = next
		return true
	}
	// winner advanced past secondMin; fall back to full merge.
	// The winner's key is now > curKey, so the "it.Key() == curKey" guard in
	// modNext will skip it — no double-advance.
	oi.fastIdx = -1
	oi.modNext(false)
	return false
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
			if it.Key() <= oi.curKey {
				it.Next()
			}
		} else if oi.lastDir != next {
			if it.Eof() {
				it.Rewind()
			}
			it.Next()
		} else if it.Key() == oi.curKey {
			it.Next()
		}
	}
}

// minIter finds the minimum current key.
// Eof iterators return ixkey.Max so they are naturally excluded.
// As a side effect it sets oi.fastIdx and oi.secondMin for the Next fast path.
func (oi *OverIter) minIter() (bool, string, uint64) {
	// NOTE: keep this code in sync with maxIter
	for {
		keyMin := ixkey.Max
		second := ixkey.Max // second-smallest key (across non-winning iters)
		winIdx := -1
		var result uint64

		// find minimum key and its most recent operation;
		// simultaneously track second-smallest key for the fast path.
		// eof iterators return ixkey.Max so they are naturally excluded.
		for i, it := range oi.iters {
			key, off := it.Cur()
			if key < keyMin { // new minimum found
				// old keyMin (if any) becomes a candidate for second
				if keyMin < second {
					second = keyMin
				}
				keyMin = key
				winIdx = i
			} else if key == keyMin {
				// tie: this iter is also at keyMin; it's a non-winner at keyMin
				if keyMin < second {
					second = keyMin
				}
			} else { // key > keyMin
				if key < second {
					second = key
				}
			}
			if key == keyMin { // track the most recent operation for keyMin
				if off&ixbuf.Delete != 0 {
					result = 0
				} else {
					result = off
				}
			}
		}
		if keyMin == ixkey.Max {
			oi.fastIdx = -1
			return false, "", 0
		}
		if result != 0 {
			oi.fastIdx = winIdx
			// secondMin <= any tombstone/update key in this iter: such an entry must
			// have a companion in another iter (invariant), creating a tie here, which
			// causes second to be <= that key. ixkey.Max when len(iters)==1.
			oi.secondMin = second
			return true, keyMin, result &^ ixbuf.Update
		}

		// Skip this deleted key - lose fast path since multiple iters may advance
		oi.fastIdx = -1
		for _, it := range oi.iters {
			if it.Key() == keyMin {
				it.Next()
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
		oi.fastIdx = -1
	} else if oi.canFast(modified, prev) && oi.fastPrev(t, prevKey) {
		return
	} else {
		oi.fastIdx = -1
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

func (oi *OverIter) fastPrev(t oiTran, prevKey string) bool {
	it := oi.iters[oi.fastIdx]
	it.Prev()
	if !it.Eof() && it.Key() > oi.secondMax {
		// still the sole winner, skip full scan.
		// Safe to use Cur() without checking Delete/Update: any tombstone or update
		// in this iter has a companion in another iter, creating a tie in maxIter
		// which sets secondMax >= that key, so the fast path is never taken past one.
		oi.curKey, oi.curOff = it.Cur()
		t.Read(oi.table, oi.iIndex, oi.curKey, prevKey)
		oi.lastDir = prev
		return true
	}
	// winner retreated past secondMax; fall back to full merge.
	// The winner's key is now < curKey, so the "it.Key() == curKey" guard in
	// modPrev will skip it — no double-advance.
	oi.fastIdx = -1
	oi.modPrev(false)
	return false
}

func (oi *OverIter) modPrev(modified bool) {
	// NOTE: keep this code in sync with modNext
	for i, it := range oi.iters {
		if it.Modified() && i != len(oi.iters)-1 {
			panic("OverIter modPrev modified not last")
		}
		if modified || it.Modified() {
			it.Seek(oi.curKey) // finds first key >= curKey
			// Seek eof means all keys < curKey; rewind to position at last key
			if it.Eof() {
				it.Rewind()
			}
			if it.Key() >= oi.curKey {
				it.Prev()
			}
		} else if oi.lastDir != prev {
			if it.Eof() {
				it.Rewind()
			}
			it.Prev()
		} else if it.Key() == oi.curKey {
			it.Prev()
		}
	}
}

// maxIter finds the maximum current key.
// As a side effect it sets oi.fastIdx and oi.secondMax for the Prev fast path.
func (oi *OverIter) maxIter() (bool, string, uint64) {
	// NOTE: keep this code in sync with minIter
	for {
		var keyMax string
		second := ixkey.Min // second-largest key (across non-winning iters)
		winIdx := -1
		var result uint64
		found := false

		// find maximum key and its most recent operation;
		// simultaneously track second-largest key for the fast path.
		for i, it := range oi.iters {
			if it.Eof() {
				continue
			}
			key, off := it.Cur()
			if !found || key > keyMax { // new maximum found
				// old keyMax (if any) becomes a candidate for second
				if found && keyMax > second {
					second = keyMax
				}
				keyMax = key
				winIdx = i
				found = true
			} else if key == keyMax {
				// tie: this iter is also at keyMax; it's a non-winner at keyMax
				if keyMax > second {
					second = keyMax
				}
			} else { // key < keyMax
				if key > second {
					second = key
				}
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
			oi.fastIdx = -1
			return false, "", 0
		}
		if result != 0 {
			oi.fastIdx = winIdx
			// secondMax >= any tombstone/update key in this iter: such an entry must
			// have a companion in another iter (invariant), creating a tie here, which
			// causes second to be >= that key. ixkey.Min when len(iters)==1.
			oi.secondMax = second
			return true, keyMax, result &^ ixbuf.Update
		}

		// Skip this deleted key - lose fast path since multiple iters may advance
		oi.fastIdx = -1
		for _, it := range oi.iters {
			if !it.Eof() && it.Key() == keyMax {
				it.Prev()
			}
		}
	}
}

func (oi *OverIter) Rewind() {
	oi.all(iterT.Rewind)
	oi.state = rewound
	oi.curKey = ""
	oi.curOff = 0
	oi.fastIdx = -1
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
