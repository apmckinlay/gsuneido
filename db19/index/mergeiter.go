// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import "github.com/apmckinlay/gsuneido/db19/index/iterator"

type iterT = iterator.T

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
	mi.curIter = mi.minIter()
	if mi.curIter == -1 {
		mi.state = eof
	} else {
		mi.curKey, mi.curOff = mi.iters[mi.curIter].Cur()
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
			if it.Seek(mi.curKey) {
				it.Next()
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

// minIter returns the index of the iterator with the minimum current value
func (mi *MergeIter) minIter() int {
	itMin := -1
	var keyMin string
	for i, it := range mi.iters {
		if !it.Eof() {
			key, _ := it.Cur()
			if itMin == -1 || key < keyMin {
				itMin = i
				keyMin = key
			}
		}
	}
	return itMin
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
	mi.curIter = mi.maxIter()
	if mi.curIter == -1 {
		mi.state = eof
	} else {
		mi.curKey, mi.curOff = mi.iters[mi.curIter].Cur()
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
			prevRewind(it)
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

// maxIter returns the index of the iterator with the maximum current value
func (mi *MergeIter) maxIter() int {
	itMax := -1
	var keyMax string
	for i, it := range mi.iters {
		if !it.Eof() {
			key, _ := it.Cur()
			if itMax == -1 || key > keyMax {
				itMax = i
				keyMax = key
			}
		}
	}
	return itMax
}

func (mi *MergeIter) Rewind() {
	mi.all(iterT.Rewind)
	mi.state = rewound
	mi.curIter = -1
	mi.curKey = ""
	mi.curOff = 0
}
