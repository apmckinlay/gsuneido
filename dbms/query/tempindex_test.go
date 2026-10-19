// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
)

// orderKind is the order in which mockSource generates rows (a, b).
type orderKind int

const (
	randomOrder  orderKind = iota // both fields random
	groupedOrder                  // a sorted in groups, b random within
	sortedOrder                   // (a, b) already sorted
)

const orderGroupSize = 100

// mockSource is a QueryMock that generates rows on the fly,
// so TempIndex can be benchmarked without a real database.
type mockSource struct {
	*QueryMock
	nrows  int
	kind   orderKind
	single bool
	i      int
}

func (s *mockSource) Rewind() {
	s.i = 0
}

func (s *mockSource) Get(*Thread, Dir) Row {
	if s.i >= s.nrows {
		return nil
	}
	i := s.i
	s.i++
	a, b := 0, 0
	switch s.kind {
	case randomOrder:
		a = int(bits.Shuffle32(uint32(i)))
		b = int(bits.Shuffle32(uint32(i) ^ 0x9e3779b9))
	case groupedOrder:
		a = i / orderGroupSize
		b = int(bits.Shuffle32(uint32(i)))
	case sortedOrder:
		a = i / orderGroupSize
		b = i % orderGroupSize
	}
	if s.single {
		var rb RecordBuilder
		rb.Add(IntVal(a))
		rb.Add(IntVal(b))
		return Row{DbRec{Record: rb.Build(), Off: 1}}
	}
	// multi: a and b are in separate records (like a join)
	var rb1, rb2 RecordBuilder
	rb1.Add(IntVal(a))
	rb2.Add(IntVal(b))
	return Row{DbRec{Record: rb1.Build(), Off: 1},
		DbRec{Record: rb2.Build(), Off: 1}}
}

func newMockSource(cols []string, nrows int, kind orderKind, single bool) *mockSource {
	var hdr *Header
	if single {
		hdr = SimpleHeader(cols)
	} else {
		hdr = NewHeader([][]string{{cols[0]}, {cols[1]}}, cols)
	}
	return &mockSource{
		QueryMock: &QueryMock{
			ColumnsResult:     cols,
			HeaderResult:      hdr,
			KeysResult:        [][]string{{cols[0]}},
			FixedResult:       Fixed{},
			NrowsN:            nrows,
			NrowsP:            nrows,
			RowSizeResult:     2 * 8,
			SingleTableResult: single,
		},
		nrows:  nrows,
		kind:   kind,
		single: single,
	}
}

func BenchmarkTempIndex(b *testing.B) {
	MakeSuTran = func(QueryTran) *SuTran { return nil }
	th := &Thread{}
	tran := testTran{}
	kinds := []struct {
		name string
		kind orderKind
	}{
		{"random", randomOrder},
		{"grouped", groupedOrder},
		{"sorted", sortedOrder},
	}
	var order = []string{"a", "b"}
	cols := []string{"a", "b"}
	modes := []struct {
		name   string
		single bool
	}{
		{"single", true},
		{"multi", false},
	}
	for _, m := range modes {
		for _, nrows := range []int{1_000, 10_000, 100_000} {
			for _, k := range kinds {
				name := fmt.Sprintf("%s/nrows=%d", k.name, nrows)
				// build measures the first Get, which builds the sortlist.
				b.Run(m.name+"/"+name+"/build", func(b *testing.B) {
					src := newMockSource(cols, nrows, k.kind, m.single)
					ti := NewTempIndex(src, order, tran)
					for b.Loop() {
						b.StopTimer()
						ti.iter = nil // force rebuilding the sortlist
						src.Rewind()
						ti.Rewind()
						b.StartTimer()
						if ti.Get(th, Next) == nil {
							b.Fatal("no rows")
						}
					}
				})
				// read measures the remaining Gets, with the sortlist built once.
				b.Run(m.name+"/"+name+"/read", func(b *testing.B) {
					src := newMockSource(cols, nrows, k.kind, m.single)
					ti := NewTempIndex(src, order, tran)
					ti.Get(th, Next) // build the sortlist up front
					n := 1
					for ti.Get(th, Next) != nil {
						n++
					}
					assert.This(n).Is(nrows)
					for b.Loop() {
						b.StopTimer()
						ti.Rewind()
						b.StartTimer()
						for ti.Get(th, Next) != nil {
						}
					}
				})
			}
		}
	}
}
