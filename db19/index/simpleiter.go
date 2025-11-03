// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"github.com/apmckinlay/gsuneido/db19/index/iface"
)

// IndexIter is the interface for index iterators.
// It is implemented by OverIter and SimpleIter.
type IndexIter interface {
	Next(t oiTran)
	Prev(t oiTran)
	Cur() (string, uint64)
	CurOff() uint64
	Eof() bool
	HasCur() bool
	Rewind()
	Range(rng Range)
}

// SimpleIter is an optimized iterator for read-only use
// with a single btree (no layers/mut).
// It does not track reads and avoids unnecessary key allocations.
type SimpleIter struct {
	t      oiTran
	it     iface.Iter
	rng    Range
	state
}

// NewSimpleIter creates an optimized iterator for read-only btree iteration.
// - the overlay must not change during iteration
// - no read tracking (e.g. read-only transaction)
func NewSimpleIter(t oiTran, ov *Overlay) IndexIter {
	bt := ov.bt
	it := bt.Iterator()
	if (ov.mut == nil || ov.mut.Len() == 0) &&
		(len(ov.layers) == 0 ||
			(len(ov.layers) == 1 && ov.layers[0].Len() == 0)) {
		return &SimpleIter{
			t:     t,
			it:    it,
			rng:   iface.All,
			state: rewound,
		}
	}
	return nil
}

func (si *SimpleIter) Eof() bool {
	return si.state == eof
}

func (si *SimpleIter) HasCur() bool {
	return si.state != eof && si.state != rewound
}

func (si *SimpleIter) Cur() (string, uint64) {
	si.checkHasCur()
	return si.it.Cur()
}

func (si *SimpleIter) CurOff() uint64 {
	si.checkHasCur()
	return si.it.Offset()
}

func (si *SimpleIter) checkHasCur() {
	if si.state == eof {
		panic("SimpleIter eof")
	}
	if si.state == rewound {
		panic("SimpleIter rewound")
	}
}

func (si *SimpleIter) Range(rng Range) {
	si.rng = rng
	si.state = rewound
	si.it.Range(rng)
}

func (si *SimpleIter) Next(t oiTran) {
	if t.Num() != si.t.Num() {
		panic("SimpleIter tran changed")
	}
	if si.state == eof {
		return // stick at eof
	}
	if si.state == rewound {
		si.state = front
	}
	si.it.Next()
	if si.it.Eof() {
		si.state = eof
	}
}

func (si *SimpleIter) Prev(t oiTran) {
	if t.Num() != si.t.Num() {
		panic("SimpleIter tran changed")
	}
	if si.state == eof {
		return // stick at eof
	}
	if si.state == rewound {
		si.state = back
	}
	si.it.Prev()
	if si.it.Eof() {
		si.state = eof
	}
}

func (si *SimpleIter) Rewind() {
	si.it.Rewind()
	si.state = rewound
}
