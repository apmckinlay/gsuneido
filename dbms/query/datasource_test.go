// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math"
	"strconv"

	. "github.com/apmckinlay/gsuneido/core"
)

type dataSource struct {
	rows []Row
	pos  dsState
}

type dsState int

const (
	dsRewound dsState = math.MinInt
	dsEof             = math.MaxInt
)

func NewDataSource(rows []Row) *dataSource {
	return &dataSource{rows: rows, pos: dsRewound}
}

func (ds dsState) String() string {
	switch ds {
	case dsRewound:
		return "Rewound"
	case dsEof:
		return "Eof"
	default:
		return strconv.Itoa(int(ds))
	}
}

// Rewind resets the position so Next gets first or Prev gets last.
func (ds *dataSource) rewind() {
	ds.pos = dsRewound
}

func (ds *dataSource) get(dir Dir) Row {
	switch ds.pos {
	case dsEof:
		return nil
	case dsRewound:
		if dir == Next {
			ds.pos = 0
		} else { // Prev
			ds.pos = dsState(len(ds.rows) - 1)
		}
	default: // within
		if dir == Next {
			ds.pos++
		} else { // Prev
			ds.pos--
		}
	}
	if ds.pos < 0 || int(ds.pos) >= len(ds.rows) {
		ds.pos = dsEof
		return nil
	}
	return ds.rows[ds.pos]
}
