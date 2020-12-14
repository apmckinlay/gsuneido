// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

type iterator interface {
	Eof() bool
	Cur() (key string, off uint64)
	Next()
	Prev()
	Rewind()
}

type MergeIter struct {
	iters []iterator
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

func NewMergeIter(iters []iterator) *MergeIter {
	its := make([]iterator, len(iters))
	copy(its, iters)
	return &MergeIter{iters: its}
}

func (mi *MergeIter) Eof() bool {
	return mi.state == eof
}

func (mi *MergeIter) Cur() (string, uint64) {
	if mi.state != within {
		return "", 0
	}
	return mi.iters[mi.curIter].Cur()
}

func (mi *MergeIter) Next() {
	if mi.state == eof {
		return // stick at eof
	}
	if mi.state == rewound {
		mi.all(iterator.Next)
		mi.state = within
	} else if mi.lastDir == next {
		mi.iters[mi.curIter].Next()
	} else { // switch direction
		mi.all(nextRewind)
	}
	mi.curIter = mi.minIter()
	if mi.curIter == -1 {
		mi.state = eof
	}
	mi.lastDir = next
}

func (mi *MergeIter) all(fn func(it iterator)) {
	for _, it := range mi.iters {
		fn(it)
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
	} else if mi.lastDir == prev {
		mi.iters[mi.curIter].Prev()
	} else { // switch direction
		mi.all(prevRewind)
	}
	mi.curIter = mi.maxIter()
	if mi.curIter == -1 {
		mi.state = eof
	}
	mi.lastDir = prev
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
}
