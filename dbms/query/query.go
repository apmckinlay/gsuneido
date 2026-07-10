// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package query implements query parsing, optimization, and execution.
/*
	Query
		Table
		Nothing
		ProjectNone
		schemaTable
			Tables
			TablesLookup
			Columns
			Indexes
			History
		Query1
			Extend
			Project / Remove
			Rename
			Sort
			Summarize
			TempIndex
			Where
			View
		Query2
			Compatible
				Union
				Compatible1
					Intersect
					Minus
			joinLike
				Times
				joinBase
					Join
					LeftJoin
					SemiJoin
*/
package query

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/opt"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Query interface {
	// Columns is all the available columns, including derived
	Columns() []string

	// Transform refactors the query for more efficient execution.
	// This stage is not cost based, transforms are applied when possible.
	//
	// Transform methods MUST ensure they call Transform on their children.
	// Transform is (mostly) bottom up, partly for path copying.
	// Which means Transform methods should start by calling Transform
	// on their children.
	//
	// Any changes should build new nodes, NOT modify nodes.
	// This is partly to ensure that constructor validation is done.
	Transform() Query

	// SetTran is used for cursors
	SetTran(tran QueryTran)

	// Order is nil for everything except Sort
	Order() []string

	// Fixed returns the field values that are constant from Extend or Where
	Fixed() Fixed

	// Updateable returns the table name if the rows from the query can be updated
	// else ""
	Updateable() string

	// SingleTable is used by TempIndex.
	// It is true if the query Get returns a single record stored in the database
	SingleTable() bool

	// Indexes returns all the indexes.
	// Unlike Keys, Indexes are physical access paths.
	Indexes() [][]string

	// Keys returns sets of fields that are unique primary keys.
	// On a table this will be the key indexes, but on other operations
	// this is logical, there may not be a physical index.
	Keys() [][]string

	// Nrows returns n the number of expected result rows from the query,
	// and p the "population" it was drawn from.
	// For example, a Where on a key selects a single row (n = 1),
	// from the entire table with p rows.
	//
	// Nrows should be the same regardless of the strategy.
	// For symmetrical/reversible operations e.g. join or union
	// it should give the same result both ways.
	//
	// Nrows does *not* incorporate frac
	Nrows() (n, p int)

	// rowSize returns the average number of bytes per row
	rowSize() int

	// Rewind resets the query so Get Next gets first, or Prev gets last
	// It does *not* clear any Select.
	Rewind()

	// Get returns the next or previous row, or nil if at end.
	// It sticks at eof until Rewind.
	Get(th *Thread, dir Dir) Row

	// Lookup returns the row matching the given key value, or nil if not found.
	// It is used by Where and Compatible (Intersect, Minus, Union)
	// It is valid (although not necessarily the most efficient)
	// to implement Lookup with Select and Get
	// in which case it should leave the select cleared.
	// Lookup should rewind.
	Lookup(th *Thread, sels Sels) Row

	// Select restricts the query to records matching the given packed values.
	// It is used by Where, Join, and LeftJoin.
	// To clear the select, use Select(nil)
	// Select should rewind.
	// It is only valid to call Select for the index chosen by optimize/setApproach.
	// If index is nil, then Select should not be called.
	Select(sels Sels)

	Header() *Header
	Output(th *Thread, rec Record)

	String() string

	cacheAdd(req Require, fixcost, varcost Cost, approach any)
	cacheGet(req Require) (fixcost, varcost Cost, approach any)
	cacheClear()

	// optimize determines the minimum cost strategy based on estimates.
	//
	// index is what is required for ordering, grouping, or lookup.
	// NOTE: there is no way to specify which of these is required.
	//
	// frac is the estimated fraction of the rows that will be read.
	// It affects the variable cost.
	// frac = 0 means only Lookup, else frac < 1 means Select or first/last
	//
	// varcost should already incorporate frac
	optimize(mode Mode, req Require) (Cost, Cost, any)

	// setApproach locks in the approach chosen by optimize.
	// index and frac must match a previous optimize call
	setApproach(req Require, approach any, tran QueryTran)

	// lookupCost returns the cost of one Lookup
	lookupCost() Cost

	// fastSingle returns whether it's a fast singleton.
	// This is mostly equivalent to whether it has an empty key().
	// Join, Intersect, and Union return false because it depends on strategy.
	// It is used by optimize (below)
	fastSingle() bool

	// Simple is simple, alternate execution method for testing.
	// It should normally be used after just parsing,
	// without transform or optimize.
	// The result may be modified - do not return internal data
	Simple(th *Thread) []Row

	// ValueGet is for Suneido.ParseQuery and queryvalue.go
	// It would be Get, but that is already used in Query.
	ValueGet(key Value) Value

	Metrics() *metrics

	// knowExactNrows returns true if Nrows returns an exact count.
	// Used by Summarize for sumTbl strategy.
	knowExactNrows() bool
}

var emptyKey = [][]string{{}}

// queryBase is embedded by almost all Query types
type queryBase struct {
	// header must be set by constructors and setApproach.
	// setApproach is necessary because the sources may get reversed
	// which affects the order of Fields
	header    *Header
	keys      [][]string
	indexes   [][]string
	fixed     Fixed
	nNrows    opt.Int
	pNrows    opt.Int
	rowSiz    opt.Int
	fast1     opt.Bool
	singleTbl opt.Bool
	lookCost  opt.Int
	cache
	metrics
}

type state byte

const (
	rewound state = iota
	within
	eof
)

type metrics struct {
	fixcost  Cost
	varcost  Cost
	costself Cost
	frac     float64
	ngets    int32
	nsels    int32
	nlooks   int32
	tget     uint64
	tgetself uint64
}

func (m *metrics) String() string {
	return fmt.Sprintf("metrics{fixcost: %v varcost: %v costself: %v frac: %.2f ngets: %d nsels: %d nlooks: %d tget: %d tgetself: %d}",
		m.fixcost, m.varcost, m.costself, m.frac, m.ngets, m.nsels, m.nlooks, m.tget, m.tgetself)
}

func (m *metrics) setCost(frac float64, fixcost, varcost Cost) {
	m.frac = frac
	m.fixcost = fixcost
	m.varcost = varcost
}

func (q *queryBase) Columns() []string {
	return q.header.Columns
}

func (q *queryBase) Header() *Header {
	return q.header
}

func (q *queryBase) Keys() [][]string {
	return q.keys
}

func (q *queryBase) Indexes() [][]string {
	return q.indexes
}

func (*queryBase) Order() []string {
	return nil
}

func (q *queryBase) Fixed() Fixed {
	return q.fixed
}

func (q *queryBase) Nrows() (int, int) {
	return q.nNrows.Get(), q.pNrows.Get()
}

func (q *queryBase) setNrows(n, p int) {
	q.nNrows.Set(n)
	q.pNrows.Set(p)
}

func (q *queryBase) rowSize() int {
	return q.rowSiz.Get()
}

func (q *queryBase) fastSingle() bool {
	return q.fast1.Get()
}

func (q *queryBase) SingleTable() bool {
	return q.singleTbl.Get()
}

func (q *queryBase) lookupCost() Cost {
	return q.lookCost.Get()
}

// Updateable is overridden by Query1
func (*queryBase) Updateable() string {
	return ""
}

func (q *queryBase) Metrics() *metrics {
	return &q.metrics
}

func (*queryBase) knowExactNrows() bool {
	return false
}

// Mode is the transaction context - cursor, read, or update.
// It affects the use of temporary indexes.
type Mode int

const (
	CursorMode Mode = iota
	ReadMode
	UpdateMode
)

func (mode Mode) String() string {
	switch mode {
	case CursorMode:
		return "cursorMode"
	case ReadMode:
		return "readMode"
	case UpdateMode:
		return "updateMode"
	default:
		panic("invalid mode")
	}
}

type Cost = int

type QueryTran interface {
	GetSchema(table string) *schema.Schema
	GetInfo(table string) *meta.Info
	GetAllInfo() []*meta.Info
	GetAllSchema() []*meta.Schema
	GetAllViews() []string
	GetView(string) string
	GetStore() *stor.Stor
	RangeFrac(table string, iIndex int, org, end string) float64
	Lookup(table string, iIndex int, key string) *DbRec
	Output(th *Thread, table string, rec Record)
	GetIndexI(table string, iIndex int) *index.Overlay
	GetRecord(off uint64) Record
	MakeLess(is *ixkey.Spec) func(x, y uint64) bool
	Read(string, int, string, string)
	IndexIter(table string, iIndex int) index.IndexIter
	Num() int
}

// Setup prepares a parsed query for execution.
// It calls Transform, Optimize, and SetApproach.
// The resulting Query is ready for execution.
func Setup(q Query, mode Mode, t QueryTran) (Query, Cost, Cost) {
	q = q.Transform()
	return setup(q, mode, 1, t)
}

// Setup1 is the same as Setup except it passes a frac of 1/nrows
// which will minimize fixed cost e.g. by avoiding temp indexes.
// It is used by DbmsLocal for Query1, QueryFirst, QueryLast
func Setup1(q Query, mode Mode, t QueryTran) (Query, Cost, Cost) {
	q = q.Transform()
	nrows, _ := q.Nrows()
	nrows = max(1, nrows) // avoid divide by zero
	return setup(q, mode, 1/float64(nrows), t)
}

func setup(q Query, mode Mode, frac float64, t QueryTran) (Query, Cost, Cost) {
	req := NoneReq(float32(frac))
	fixcost, varcost := Optimize(q, mode, req)
	if fixcost+varcost >= impossible {
		panic("invalid query: " + String(q))
	}
	q = SetApproach(q, req, t)
	if mode == CursorMode {
		setCursorMode(q)
	}
	return q, fixcost, varcost
}

// SetupKey is like Setup but it ensures a key index
// It is used by updateAction (action.go)
func SetupKey(q Query, mode Mode, t QueryTran) Query {
	// we use ReqGroup because it matches the search for a key
	// although we don't actually need the rows to be grouped
	q = q.Transform()
	best := newBest[[]string]()
	for _, key := range q.Keys() {
		f, v, _ := optimize(q, mode, GroupReq(key, 1, 1))
		best.update(f, v, key)
	}
	if best.none() {
		panic("invalid query: " + String(q))
	}
	q = SetApproach(q, GroupReq(best.data, 1, 1), t)
	return q
}

// SetupIdx is like Setup but specifies an index
// e.g. to test Select or Lookup
func SetupIdx(q Query, mode Mode, t QueryTran, index []string) Query {
	req := OrderReq(index, 1)
	fixcost, varcost := Optimize(q, mode, req)
	if fixcost+varcost >= impossible {
		panic("invalid query: " + String(q))
	}
	q = SetApproach(q, req, t)
	if mode == CursorMode {
		setCursorMode(q)
	}
	return q
}

// outOfOrder is added when we reverse sources.
// This discourages reversing the order without a good reason
// which makes tests and debugging easier.
const outOfOrder = 10

const impossible = Cost(math.MaxInt / 64) // allow for adding impossible's

//-------------------------------------------------------------------
// new version of Optimize using Require (not used yet)
// initially duplicates the existing one, will eventually replace it

func Optimize(q Query, mode Mode, req Require) (fixcost, varcost Cost) {
	fixcost, varcost, _ = optimize(q, mode, req)
	return fixcost, varcost
}

func optimize(q Query, mode Mode, req Require) (
	fixcost, varcost Cost, approach any) {
	assert.That(!math.IsNaN(float64(req.frac)) && !math.IsInf(float64(req.frac), 0))
	if !set.Subset(q.Columns(), req.cols) {
		return impossible, impossible, nil
	}

	// this condition must match SetApproach
	// A fastSingle node (or one whose fixed covers req.cols) trivially
	// satisfies any require, so the qualitative aspect (cols/use) is
	// irrelevant. Clear cols AND nseeks
	// frac is kept as it scales the (single) row's varcost.
	if q.fastSingle() || q.Fixed().All(req.cols) {
		req.cols = nil
		req.nseeks = 0
		req.use = ReqNone
	}
	if fixcost, varcost, app := q.cacheGet(req); varcost >= 0 {
		return fixcost, varcost, app
	}
	fixcost, varcost, app := optTempIndex(q, mode, req)
	assert.That(fixcost >= 0 && varcost >= 0)
	q.cacheAdd(req, fixcost, varcost, app)
	return fixcost, varcost, app
}

// optTempIndex determines if a TempIndex is a benefit
// and if it is, returns a special tempIndex approach
// that is processed by SetApproach which creates the actual TempIndex
func optTempIndex(q Query, mode Mode, req Require) (
	fixcost, varcost Cost, approach any) {
	traceQO := func(more ...any) {
		if trace.QueryOpt.On() {
			args := append([]any{req, "="}, more...)
			trace.QueryOpt.Println(mode, args...)
			trace.Println(strategy(q, 1))
		}
	}
	traceQO("optTempIndex", "----------------")

	indexedFixCost, indexedVarCost, indexedApp := q.optimize(mode, req)
	assert.That(indexedFixCost >= 0 && indexedVarCost >= 0)

	u := req.use
	if u == ReqNone || !tempIndexable(mode) {
		traceQO(indexedFixCost + indexedVarCost)
		return indexedFixCost, indexedVarCost, indexedApp
	}

	nrows, _ := q.Nrows()
	assert.That(nrows >= 0)
	best := newBest[tiApproach]()

	// with no index
	noIdxOrder := req.cols
	if u == ReqUnique {
		noIdxOrder = tempIndexKey(q, req.cols)
	}
	optTI(&best, q, mode, NoneReq(req.frac), nrows, factorNone, noIdxOrder)

	// with required index
	optTI(&best, q, mode, req, nrows, factorAll, req.cols)

	// with "best" index
	if bestIndex := tempIndexBest(q, req.cols); bestIndex != nil {
		optTI(&best, q, mode, OrderReq(bestIndex, req.frac), nrows, factorPre, req.cols)
	}

	// key-subset candidates for ReqUnique
	if u == ReqUnique {
		fixed := q.Fixed()
		nReqColsUnfixed := countUnfixed(req.cols, fixed)
		for _, key := range q.Keys() {
			if !indexCovered(key, req.cols, fixed) {
				continue
			}
			nKeyUnfixed := countUnfixed(key, fixed)
			if nKeyUnfixed == 0 || nKeyUnfixed >= nReqColsUnfixed {
				continue
			}
			keyUnfixed := fixed.RemoveFrom(key)
			optTI(&best, q, mode, UniqueReq(keyUnfixed, req.nseeks), nrows, factorAll, keyUnfixed)
		}
	}

	// for ReqUnique or ReqGroup with nseeks, add per-lookup cost on temp index
	if u == ReqUnique || (u == ReqGroup && req.nseeks > 0) {
		perLookup := Cost(400)
		if q.SingleTable() {
			perLookup = Cost(200)
		}
		best.varcost += Cost(req.nseeks) * perLookup
	}

	tempIndexCost := best.cost()
	indexedCost := indexedFixCost + indexedVarCost
	if indexedCost <= tempIndexCost {
		traceQO("indexed", indexedCost, "<=", tempIndexCost)
		return indexedFixCost, indexedVarCost, indexedApp
	}
	traceQO("tempindex", best.data.index, tempIndexCost, "<", indexedCost)
	return best.fixcost, best.varcost,
		&tempIndex{index: best.data.tiOrder, srcapp: best.data.srcapp,
			srcindex:   best.data.index,
			srcfixcost: best.data.srcfixcost, srcvarcost: best.data.srcvarcost}
}

// tempIndexKey finds the smallest key (by unfixed count) that's covered by cols.
// Used to minimize tempindex order when source is unordered.
func tempIndexKey(q Query, cols []string) []string {
	fixed := q.Fixed()
	var bestKey []string
	bestN := -1
	for _, key := range q.Keys() {
		if !indexCovered(key, cols, fixed) {
			continue
		}
		n := countUnfixed(key, fixed)
		if n == 0 {
			continue
		}
		if bestKey == nil || n < bestN {
			bestKey = key
			bestN = n
		}
	}
	if bestKey != nil {
		return fixed.RemoveFrom(bestKey)
	}
	return cols
}

func optTI(best *best[tiApproach], q Query, mode Mode, req Require, nrows, factor int, tiOrder []string) {
	srcReq := OrderReq(req.cols, 1)
	srcfixcost, srcvarcost, srcapp := q.optimize(mode, srcReq)
	assert.That(srcfixcost >= 0 && srcvarcost >= 0)
	fixcost, varcost := ticost(srcfixcost+srcvarcost, q, req.cols, nrows, float64(req.frac), factor)
	best.update(fixcost, varcost, tiApproach{
		index:      req.cols,
		tiOrder:    tiOrder,
		srcfixcost: srcfixcost,
		srcvarcost: srcvarcost,
		srcapp:     srcapp,
	})
}

//-------------------------------------------------------------------

const factorAll = 105  // ???
const factorPre = 110  // ???
const factorNone = 256 // ???

var ticostAdj = 0 // for tests, to discourage temp indexes

func ticost(srccost int, q Query, index []string, nrows int, frac float64,
	factor int) (Cost, Cost) {
	fixcost := srccost + ticostAdj + 1000 // ???
	fixcost += 100 * len(index)           // prefer fewer fields
	if nrows > 0 {
		fnrows := float64(nrows)
		fixcost += factor * Cost(fnrows*math.Log(fnrows)) // empirical
	}
	varcost := Cost(frac * float64(nrows) * 100) // ???
	if !q.SingleTable() {
		varcost *= 2 // ???
	}
	return fixcost, varcost
}

// tempIndexBest finds the index that has the longest common prefix.
// NOTE: it assumes that all indexes have the same cost
// which is true for simple cases like a table
// but not for more complex queries.
func tempIndexBest(q Query, index []string) []string {
	// fmt.Println("tempIndexBest", index, q)
	fixed := q.Fixed()
	var bestIndex []string
	var bestOn int
	for _, ix := range q.Indexes() {
		on := orderedn(ix, index, fixed)
		// fmt.Println("\torderedn", ix, index, fixed, "=>", on)
		if on > 0 && on < len(index) && on > bestOn {
			bestOn = on
			bestIndex = ix
		}
	}
	return bestIndex
}

type tiApproach struct {
	index      []string
	tiOrder    []string
	srcfixcost Cost
	srcvarcost Cost
	srcapp     any
}

// tempIndex is a special approach that is added by optTempIndex
// to be used by SetApproach to insert a TempIndex when required
type tempIndex struct {
	index      []string
	srcapp     any
	srcindex   []string
	srcfixcost Cost
	srcvarcost Cost
}

func tempIndexable(mode Mode) bool {
	if mode == ReadMode {
		return true
	}
	if mode == CursorMode {
		return false
	}
	// else updateMode
	return true
	// BUG this matches jSuneido, but it is not correct.
	// A temp index allows reading deleted or old versions of records.
	// But there is a big performance penalty
	// especially from the key sort added by QueryApply.
}

func min3(fixcost1, varcost1 Cost, app1 any, fixcost2, varcost2 Cost, app2 any,
	fixcost3, varcost3 Cost, app3 any) (Cost, Cost, any) {
	fixcost, varcost, app := fixcost1, varcost1, app1
	if fixcost2+varcost2 < fixcost+varcost {
		fixcost, varcost, app = fixcost2, varcost2, app2
	}
	if fixcost3+varcost3 < fixcost+varcost {
		fixcost, varcost, app = fixcost3, varcost3, app3
	}
	return fixcost, varcost, app
}

var tempIndexCount atomic.Int64
var _ = AddInfo("query.tempindex", &tempIndexCount)

// SetApproach finalizes the chosen approach.
// It also adds temp indexes where required.
func SetApproach(q Query, req Require, tran QueryTran) Query {
	// must match optimize's guard (see comment there)
	if q.fastSingle() || q.Fixed().All(req.cols) {
		req.cols = nil
		req.nseeks = 0
		req.use = ReqNone
	}
	fixcost, varcost, approach := q.cacheGet(req)
	q.cacheClear()
	if fixcost == -1 {
		panic("SetApproach: not found in cache")
	}
	assert.That(fixcost >= 0 && varcost >= 0)
	if app, ok := approach.(*tempIndex); ok {
		q.Metrics().setCost(1, app.srcfixcost, app.srcvarcost)
		q.setApproach(OrderReq(app.srcindex, 1), app.srcapp, tran)
		ti := NewTempIndex(q, app.index, tran)
		ti.setCost(float64(req.frac), fixcost, varcost)
		tempIndexCount.Add(1)
		return ti
	}
	q.Metrics().setCost(float64(req.frac), fixcost, varcost)
	q.setApproach(req, approach, tran)
	return q
}

// execution --------------------------------------------------------

// GetNext1 returns the next row from q if it matches sels, else nil.
// Used when Lookup is implemented with Select+Get —
// Select only restricts by the physical index prefix,
// so GetNext1 verifies the row matches all of sels.
func GetNext1(q Query, th *Thread, sels Sels) Row {
	// this should *not* have to loop because the index should be unique
	row := q.Get(th, Next)
	if row != nil {
		debug.assert(q.Get(th, Next) == nil)
		if singletonFilter(q.Header(), row, sels) {
			return row
		}
	}
	return nil
}

// lookupViaSelectGet implements Lookup via Select+Get,
// verifying the row matches all sels (since Select only restricts
// by the physical index prefix) and clearing the select afterwards.
func lookupViaSelectGet(q Query, th *Thread, sels Sels) Row {
	q.Select(sels)
	defer q.Select(nil)
	return GetNext1(q, th, sels)
}

// Query1 -----------------------------------------------------------

type Query1 struct {
	queryBase
	source Query
}

func (q1 *Query1) Updateable() string {
	return q1.source.Updateable()
}

func (q1 *Query1) SetTran(t QueryTran) {
	q1.source.SetTran(t)
}

func (q1 *Query1) Output(th *Thread, rec Record) {
	q1.source.Output(th, rec)
}

func (q1 *Query1) Rewind() {
	q1.source.Rewind()
}

type q1i interface {
	Source() Query
	String() string
}

func (q1 *Query1) Source() Query {
	return q1.source
}

// Query2 -----------------------------------------------------------

type Query2 struct {
	queryBase
	source1 Query
	source2 Query
}

func (q2 *Query2) SetTran(t QueryTran) {
	q2.source1.SetTran(t)
	q2.source2.SetTran(t)
}

func (q2 *Query2) SingleTable() bool {
	return false // not single
}

func (*Query2) Output(*Thread, Record) {
	panic("can't output to this query")
}

func (q2 *Query2) keypairs() [][]string {
	var keys [][]string
	for _, k1 := range q2.source1.Keys() {
		for _, k2 := range q2.source2.Keys() {
			keys = set.AddUniqueFn(keys, set.Union(k1, k2), set.Equal[string])
		}
	}
	assert.That(len(keys) != 0)
	return keys
}

type q2i interface {
	q1i
	Source2() Query
}

func (q2 *Query2) Source() Query {
	return q2.source1
}

func (q2 *Query2) Source2() Query {
	return q2.source2
}

// ------------------------------------------------------------------

// String prints the full query, including child sources
// whereas query.String only shows that operation
func String(q Query) string {
	switch qi := q.(type) {
	case q2i:
		return paren2(qi.Source()) + " " + q.String() + " " + paren1(qi.Source2())
	case *Sort:
		return String(qi.Source()) + str.Opt(" ", q.String()) // no parens
	case *View:
		return q.String()
	case q1i:
		return paren2(qi.Source()) + str.Opt(" ", q.String())
	default:
		return q.String()
	}
}

func paren1(q Query) string {
	switch q.(type) {
	case *Table, *Tables, *TablesLookup, *Columns, *Indexes, *Views,
		*Nothing, *ProjectNone:
		return String(q)
	}
	return "(" + String(q) + ")"
}

func paren2(q Query) string {
	if _, ok := q.(q2i); ok {
		return "(" + String(q) + ")"
	}
	return String(q)
}

// ------------------------------------------------------------------

func Strategy(q Query) string {
	return strategy(q, 0)
}

const indent1 = "    "

func strategy(q Query, indent int) string { // recursive
	in := strings.Repeat(indent1, indent)
	nrows, pop := q.Nrows()
	m := q.Metrics()
	cost := "{"
	if m.frac != 1 {
		cost += fmt.Sprintf("%.3fx ", m.frac)
	}
	cost += trace.Number(nrows)
	if nrows != pop {
		cost += "/" + trace.Number(pop)
	}
	cost += " " + trace.Number(m.fixcost) + "+" + trace.Number(m.varcost)
	cost += "} "
	switch q := q.(type) {
	case *Sort:
		if q.String() == "" {
			return strategy(q.Source(), indent)
		} else {
			return strategy(q.Source(), indent) + "\n" +
				in + cost + q.String()
		}
	case q2i:
		return strategy(q.Source(), indent+1) + "\n" +
			in + cost + q.String() + "\n" +
			strategy(q.Source2(), indent+1)
	case q1i:
		return strategy(q.Source(), indent) + "\n" +
			in + cost + q.String()
	default:
		return in + cost + q.String()
	}
}

// Strategy2 is like Strategy but without the cost/size estimates
// so it is more stable for tests
func Strategy2(q Query) string {
	return strategy2(q, 0)
}

func strategy2(q Query, indent int) string { // recursive
	in := strings.Repeat(indent1, indent)
	switch q := q.(type) {
	case *Sort:
		if q.String() == "" {
			return strategy2(q.Source(), indent)
		} else {
			return strategy2(q.Source(), indent) + "\n" +
				in + q.String()
		}
	case q2i:
		return strategy2(q.Source(), indent+1) + "\n" +
			in + q.String() + "\n" +
			strategy2(q.Source2(), indent+1)
	case q1i:
		return strategy2(q.Source(), indent) + "\n" +
			in + q.String()
	default:
		return in + q.String()
	}
}

func CalcSelf(q0 Query) { // recursive
	m := q0.Metrics()
	if m.tgetself != 0 {
		return // already calculated
	}
	switch q := q0.(type) {
	case q2i:
		m1 := q.Source().Metrics()
		m2 := q.Source2().Metrics()
		m.tgetself = m.tget - (m1.tget + m2.tget)
		m.costself = (m.fixcost + m.varcost) -
			(m1.fixcost + m1.varcost + m2.fixcost + m2.varcost)
		CalcSelf(q.Source())
		CalcSelf(q.Source2())
	case q1i:
		sm := q.Source().Metrics()
		m.tgetself = m.tget - sm.tget
		m.costself = (m.fixcost + m.varcost) - (sm.fixcost + sm.varcost)
		CalcSelf(q.Source())
	default:
		m.tgetself = q0.Metrics().tget
		m.costself = q0.Metrics().fixcost + q0.Metrics().varcost
	}
}

func setCursorMode(q Query) {
	switch q := q.(type) {
	case q2i:
		setCursorMode(q.Source())
		setCursorMode(q.Source2())
	case q1i:
		setCursorMode(q.Source())
	case *Table:
		q.cursorMode = true
	}
}

// func unpack(packed []string) []Value {
// 	vals := make([]Value, len(packed))
// 	for i, p := range packed {
// 		if p == ixkey.Max {
// 			vals[i] = SuStr("<max>")
// 		} else {
// 			vals[i] = Unpack(p)
// 		}
// 	}
// 	return vals
// }

//-------------------------------------------------------------------

var debug debugT

type debugT struct{}

func (debugT) assert(cond bool) {
	assert.That(cond)
}

// func (debugT) assert(cond bool) {
// }
