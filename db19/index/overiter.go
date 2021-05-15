// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
)

type iterT = iterator.T

type Range = iterator.Range

// OverIter is a Suneido style iterator
// that merges several other Suneido style iterators.
//
// We need to keep our own curKey/Off independent of the source iterators
// because new source iterators may be returned by the callback.
type OverIter struct {
	table    string
	ixcols   []string
	iters    []iterT
	modCount int
	curKey   string
	curOff   uint64
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
	within
	eof
)

type dir int8

const (
	next dir = +1
	prev dir = -1
)

type oiTran interface {
	GetIndex(table string, columns []string) *Overlay
}

func NewOverIter(table string, ixcols []string) *OverIter {
	return &OverIter{table: table, ixcols: ixcols, rng: iterator.All}
}

func (mi *OverIter) Eof() bool {
	return mi.state == eof
}

func (mi *OverIter) Cur() (string, uint64) {
	return mi.curKey, mi.curOff
}

func (mi *OverIter) Range(rng Range) {
	mi.rng = rng
	mi.state = rewound
	for _, it := range mi.iters {
		it.Range(rng)
	}
}

func (mi *OverIter) Next(t oiTran) {
	if mi.state == eof {
		return // stick at eof
	}
	if mi.state == rewound {
		if t != mi.tran {
			// some of the iterators may still be valid
			// but simplest is just to get fresh iterators
			ov := t.GetIndex(mi.table, mi.ixcols)
			mi.newIters(ov)
			mi.tran = t
		}
		mi.all(iterT.Next)
		mi.state = within
	} else {
		mi.modNext(t)
	}
	mi.curIter, mi.curKey, mi.curOff = mi.minIter()
	if mi.curIter == -1 {
		mi.state = eof
	}
	mi.lastDir = next
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
	modified := false
	if t != mi.tran {
		ov := t.GetIndex(mi.table, mi.ixcols)
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
			} else if key == keyMin {
				off = ixbuf.Combine(offMin, off)
				mi.iters[itMin].Next()
				if off == 0 {
					// add,delete so skip
					// may not be the final minimum, but still need to skip
					it.Next()
					continue outer
				}
				itMin = i
				offMin = off
			}
		}
		return itMin, keyMin, offMin
	}
}

func (mi *OverIter) Prev(t oiTran) {
	if mi.state == eof {
		return // stick at eof
	}
	if mi.state == rewound {
		if t != mi.tran {
			// some of the iterators may still be valid
			// but simplest just to get fresh iterators
			ov := t.GetIndex(mi.table, mi.ixcols)
			mi.newIters(ov)
			mi.tran = t
		}
		mi.all(iterT.Prev)
		mi.state = within
	} else {
		mi.modPrev(t)
	}
	mi.curIter, mi.curKey, mi.curOff = mi.maxIter()
	if mi.curIter == -1 {
		mi.state = eof
	}
	mi.lastDir = prev
}

func (mi *OverIter) modPrev(t oiTran) {
	modified := false
	if t != mi.tran {
		ov := t.GetIndex(mi.table, mi.ixcols)
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
outer:
	for {
		itMax := -1
		var keyMax string
		var offMax uint64
		for i, it := range mi.iters {
			if !it.Eof() {
				key, off := it.Cur()
				if itMax == -1 || key > keyMax {
					itMax = i
					keyMax = key
					offMax = off
				} else if key == keyMax {
					off = ixbuf.Combine(offMax, off)
					if off == 0 {
						// add,delete so skip
						// may not be the final minimum, but still need to skip
						it.Next()
						mi.iters[itMax].Next()
						continue outer
					}
				}
			}
		}
		return itMax, keyMax, offMax
	}
}

func (mi *OverIter) Rewind() {
	mi.all(iterT.Rewind)
	mi.state = rewound
	mi.curIter = -1
	mi.curKey = ""
	mi.curOff = 0
}
