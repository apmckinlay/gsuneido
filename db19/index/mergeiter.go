// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

type iterator interface {
	Eof() bool
	Modified() bool
	Cur() (key string, off uint64)
	Next()
	Prev()
	Rewind()
	// Seek returns true if the key was found
	Seek(key string) bool
}

type mergeCallback func(modCount int) (int, []iterator)

type MergeIter struct {
	callback mergeCallback
	iters    []iterator
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
	if mi.state != within {
		return "", 0
	}
	return mi.curKey, mi.curOff
}

func (mi *MergeIter) Next() {
	if mi.state == eof {
		return // stick at eof
	}
	if mi.state == rewound {
		mi.all(iterator.Next)
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

func (mi *MergeIter) all(fn func(it iterator)) {
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

func nextRewind(it iterator) {
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
		mi.all(iterator.Prev)
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

func prevRewind(it iterator) {
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
	mi.all(iterator.Rewind)
	mi.state = rewound
	mi.curIter = -1
	mi.curKey = ""
	mi.curOff = 0
}
