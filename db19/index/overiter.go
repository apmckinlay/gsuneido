// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
)

type iterT = iterator.T

type Range = iterator.Range

// mergeCallback is a function passed into a MergeIter
// so it can determine if the underlying container (normally an Overlay)
// has been modified.
// The iterator passes its last known modCount
// and if the container's modCount has changed,
// it returns the new modCount and the new source iterators.
// If the modCount has not changed, it returns nil instead of new iterators.
type mergeCallback func(modCount int) (int, []iterT)

// MergeIter is a Suneido style iterator
// that merges several other Suneido style iterators.
//
// We need to keep our own curKey/Off independent of the source iterators
// because new source iterators may be returned by the callback.
type MergeIter struct {
	callback mergeCallback
	iters    []iterT
	modCount int
	curKey   string
	curOff   uint64
	// curIter is the iterator containing the current item = iters[curIter]
	curIter int
	state
	lastDir dir
	rng     Range
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

func NewMergeIter(callback mergeCallback) *MergeIter {
	modCount, iters := callback(-1)
	return &MergeIter{callback: callback, modCount: modCount, iters: iters}
}

func (mi *MergeIter) Eof() bool {
	return mi.state == eof
}

func (mi *MergeIter) Cur() (string, uint64) {
	return mi.curKey, mi.curOff
}

func (mi *MergeIter) Range(rng Range) {
	 mi.rng = rng
	 mi.state = rewound
	 for _,it := range mi.iters {
		it.Range(rng)
	 }
}

func (mi *MergeIter) Next() {
	if mi.state == eof {
		return // stick at eof
	}
	if mi.state == rewound {
		modCount, iters := mi.callback(mi.modCount)
		if iters != nil { // modified
			mi.modCount, mi.iters = modCount, iters
		}
		mi.all(iterT.Next)
		mi.state = within
	} else {
		mi.modNext()
	}
	mi.curIter, mi.curKey, mi.curOff = mi.minIter()
	if mi.curIter == -1 {
		mi.state = eof
	}
	mi.lastDir = next
}

func (mi *MergeIter) all(fn func(it iterT)) {
	for _, it := range mi.iters {
		fn(it)
	}
}

func (mi *MergeIter) modNext() {
	modCount, iters := mi.callback(mi.modCount)
	modified := iters != nil
	mi.modCount = modCount
	if modified {
		mi.iters = iters
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
			nextRewind(it)
		} else if i == mi.curIter {
			it.Next()
		}
	}
}

func nextRewind(it iterT) {
	if it.Eof() {
		it.Rewind()
	}
	it.Next()
}

// minIter finds the the minimum current key
func (mi *MergeIter) minIter() (int, string, uint64) {
outer:
	for {
		itMin := -1
		var keyMin string
		var offMin uint64
		for i, it := range mi.iters {
			if !it.Eof() {
				key, off := it.Cur()
				if itMin == -1 || key < keyMin {
					itMin = i
					keyMin = key
					offMin = off
				} else if key == keyMin {
					off = ixbuf.Combine(offMin, off)
					if off == 0 {
						// add,delete so skip
						// may not be the final minimum, but still need to skip
						it.Next()
						mi.iters[itMin].Next()
						continue outer
					}
				}
			}
		}
		return itMin, keyMin, offMin
	}
}

func (mi *MergeIter) Prev() {
	if mi.state == eof {
		return // stick at eof
	}
	if mi.state == rewound {
		modCount, iters := mi.callback(mi.modCount)
		if iters != nil { // modified
			mi.modCount, mi.iters = modCount, iters
		}
		mi.all(iterT.Prev)
		mi.state = within
	} else {
		mi.modPrev()
	}
	mi.curIter, mi.curKey, mi.curOff = mi.maxIter()
	if mi.curIter == -1 {
		mi.state = eof
	}
	mi.lastDir = prev
}

func (mi *MergeIter) modPrev() {
	modCount, iters := mi.callback(mi.modCount)
	modified := iters != nil
	mi.modCount = modCount
	if modified {
		mi.iters = iters
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
			prevRewind(it)
		} else if i == mi.curIter {
			it.Prev()
		}
	}
}

func prevRewind(it iterT) {
	if it.Eof() {
		it.Rewind()
	}
	it.Prev()
}

// maxIter finds the maximum current key
func (mi *MergeIter) maxIter() (int, string, uint64) {
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

func (mi *MergeIter) Rewind() {
	mi.all(iterT.Rewind)
	mi.state = rewound
	mi.curIter = -1
	mi.curKey = ""
	mi.curOff = 0
}
