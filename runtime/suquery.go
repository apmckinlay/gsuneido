package runtime

import "github.com/apmckinlay/gsuneido/runtime/types"

// SuQuery is a database query
type SuQuery struct {
	CantConvert
	tran   *SuTran
	query  string
	iquery IQuery
	eof    Dir
}

var _ Value = (*SuQuery)(nil)

func NewSuQuery(t *SuTran, query string, iquery IQuery) *SuQuery {
	return &SuQuery{tran: t, query: query, iquery: iquery}
}

func (*SuQuery) Get(*Thread, Value) Value {
	panic("query does not support get")
}

func (*SuQuery) Put(*Thread, Value, Value) {
	panic("query does not support put")
}

func (*SuQuery) RangeTo(int, int) Value {
	panic("query does not support range")
}

func (*SuQuery) RangeLen(int, int) Value {
	panic("query does not support range")
}

func (*SuQuery) Hash() uint32 {
	panic("query hash not implemented")
}

func (*SuQuery) Hash2() uint32 {
	panic("query hash not implemented")
}

func (q *SuQuery) Equal(other interface{}) bool {
	if q2, ok := other.(*SuQuery); ok {
		return q == q2
	}
	return false
}

func (*SuQuery) Compare(Value) int {
	panic("query compare not implemented")
}

func (*SuQuery) Call(*Thread, *ArgSpec) Value {
	panic("can't call query")
}

func (*SuQuery) Type() types.Type {
	return types.Query
}

func (q *SuQuery) String() string {
	return "Query('" + q.query + "')"
}

// QueryMethods is initialized by the builtin package
var QueryMethods Methods

func (*SuQuery) Lookup(_ *Thread, method string) Callable {
	return QueryMethods[method]
}

//-------------------------------------------------------------------

func (q *SuQuery) Close() {
	q.iquery.Close()
}

func (q *SuQuery) Columns() Value {
	hdr := q.iquery.Header()
	ob := &SuObject{}
	for _, col := range hdr.Columns {
		ob.Add(SuStr(col))
	}
	return ob
}

func (q *SuQuery) GetRec(dir Dir) Value {
	if dir == q.eof {
		return False
	}
	row := q.iquery.Get(dir)
	if row == nil {
		q.eof = dir
		return False
	}
	q.eof = 0
	return SuRecordFromRow(row, q.iquery.Header(), q.tran)
}

func (q *SuQuery) Keys() Value {
	return q.iquery.Keys()
}

func (q *SuQuery) Order() Value {
	return q.iquery.Order()
}

func (q *SuQuery) Output(th *Thread, ob Container) {
	rec := ob.ToRecord(th, q.iquery.Header())
	q.iquery.Output(rec)
}

func (q *SuQuery) Rewind() {
	q.iquery.Rewind()
}

func (q *SuQuery) RuleColumns() Value {
	hdr := q.iquery.Header()
	ob := &SuObject{}
	for _, col := range hdr.Rules() {
		ob.Add(SuStr(col))
	}
	return ob
}

func (q *SuQuery) Strategy() Value {
	return SuStr(q.iquery.Strategy())
}
