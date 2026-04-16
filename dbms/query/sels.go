// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "iter"

type Sels []Sel

type Sel struct {
	col string
	val string
}

func NewSel(col, val string) Sel {
	return Sel{col, val}
}

func (sels Sels) MustGet(col string) string {
	for _, sel := range sels {
		if sel.col == col {
			return sel.val
		}
	}
	panic("Sels.Get can't find " + col)
}

func (sels Sels) Get(col string) (string, bool) {
	for _, sel := range sels {
		if sel.col == col {
			return sel.val, true
		}
	}
	return "", false
}

func (sels Sels) HasCol(col string) bool {
	for _, sel := range sels {
		if sel.col == col {
			return true
		}
	}
	return false
}

func (sels Sels) FindCol(col string) int {
	for i, sel := range sels {
		if sel.col == col {
			return i
		}
	}
	return -1
}

// ColsAre returns true if the columns are the same set (ignoring order).
// Does NOT handle duplicates.
func (sels Sels) ColsAre(cols []string) bool {
	if len(sels) != len(cols) {
		return false
	}
outer:
	for _, col := range cols {
		for _, sel := range sels {
			if col == sel.col {
				continue outer
			}
		}
		return false
	}
	return true
}

func (sels Sels) All() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		for _, sel := range sels {
			if !yield(sel.col, sel.val) {
				return
			}
		}
	}
}
