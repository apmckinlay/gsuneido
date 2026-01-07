// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import . "github.com/apmckinlay/gsuneido/core"

// QueryMock is a test mock that satisfies the Query interface.
// By default all methods panic. Set fields to override return values.
type QueryMock struct {
	cache
	ColumnsResult     []string
	TransformResult   Query
	OrderResult       []string
	FixedResult       []Fixed
	UpdateableResult  string
	SingleTableResult bool
	IndexesResult     [][]string
	KeysResult        [][]string
	NrowsN, NrowsP    int
	RowSizeResult     int
	GetResult         Row
	LookupResult      Row
	HeaderResult      *Header
	StringResult      string
	OptimizeResult    struct {
		Fixcost, Varcost Cost
		Approach         any
	}
	LookupCostResult Cost
	FastSingleResult bool
	SimpleResult     []Row
	ValueGetResult   Value
	MetricsResult    *metrics
}

var _ Query = (*QueryMock)(nil)

func (m *QueryMock) Columns() []string {
	if m.ColumnsResult != nil {
		return m.ColumnsResult
	}
	panic("QueryMock.Columns not implemented")
}

func (m *QueryMock) Transform() Query {
	if m.TransformResult != nil {
		return m.TransformResult
	}
	panic("QueryMock.Transform not implemented")
}

func (*QueryMock) SetTran(QueryTran) {
	panic("QueryMock.SetTran not implemented")
}

func (m *QueryMock) Order() []string {
	if m.OrderResult != nil {
		return m.OrderResult
	}
	panic("QueryMock.Order not implemented")
}

func (m *QueryMock) Fixed() []Fixed {
	if m.FixedResult != nil {
		return m.FixedResult
	}
	panic("QueryMock.Fixed not implemented")
}

func (m *QueryMock) Updateable() string {
	return m.UpdateableResult
}

func (m *QueryMock) SingleTable() bool {
	return m.SingleTableResult
}

func (m *QueryMock) Indexes() [][]string {
	if m.IndexesResult != nil {
		return m.IndexesResult
	}
	panic("QueryMock.Indexes not implemented")
}

func (m *QueryMock) Keys() [][]string {
	if m.KeysResult != nil {
		return m.KeysResult
	}
	panic("QueryMock.Keys not implemented")
}

func (m *QueryMock) Nrows() (int, int) {
	return m.NrowsN, m.NrowsP
}

func (m *QueryMock) rowSize() int {
	return m.RowSizeResult
}

func (*QueryMock) Rewind() {
	panic("QueryMock.Rewind not implemented")
}

func (m *QueryMock) Get(*Thread, Dir) Row {
	if m.GetResult != nil {
		return m.GetResult
	}
	panic("QueryMock.Get not implemented")
}

func (m *QueryMock) Lookup(*Thread, []string, []string) Row {
	if m.LookupResult != nil {
		return m.LookupResult
	}
	panic("QueryMock.Lookup not implemented")
}

func (*QueryMock) Select([]string, []string) {
	panic("QueryMock.Select not implemented")
}

func (m *QueryMock) Header() *Header {
	if m.HeaderResult != nil {
		return m.HeaderResult
	}
	panic("QueryMock.Header not implemented")
}

func (*QueryMock) Output(*Thread, Record) {
	panic("QueryMock.Output not implemented")
}

func (m *QueryMock) String() string {
	return m.StringResult
}

func (m *QueryMock) optimize(Mode, []string, float64) (Cost, Cost, any) {
	return m.OptimizeResult.Fixcost, m.OptimizeResult.Varcost,
		m.OptimizeResult.Approach
}

func (*QueryMock) setApproach([]string, float64, any, QueryTran) {
	panic("QueryMock.setApproach not implemented")
}

func (m *QueryMock) lookupCost() Cost {
	return m.LookupCostResult
}

func (m *QueryMock) fastSingle() bool {
	return m.FastSingleResult
}

func (m *QueryMock) Simple(*Thread) []Row {
	if m.SimpleResult != nil {
		return m.SimpleResult
	}
	panic("QueryMock.Simple not implemented")
}

func (m *QueryMock) ValueGet(Value) Value {
	if m.ValueGetResult != nil {
		return m.ValueGetResult
	}
	panic("QueryMock.ValueGet not implemented")
}

func (m *QueryMock) Metrics() *metrics {
	if m.MetricsResult != nil {
		return m.MetricsResult
	}
	panic("QueryMock.Metrics not implemented")
}

func (*QueryMock) knowExactNrows() bool {
	return false
}
