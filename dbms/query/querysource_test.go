// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

type dataSource struct {
	rows []Row
	pos  dsState
}

type dsState int

const (
	dsRewound dsState = math.MaxInt - 1
	dsAtEnd           = math.MaxInt
	dsAtOrg           = math.MinInt
)

func NewDataSource(rows []Row) *dataSource {
	return &dataSource{rows: rows, pos: dsRewound}
}

func (ds dsState) String() string {
	switch ds {
	case dsRewound:
		return "Rewound"
	case dsAtEnd:
		return "AtEnd"
	case dsAtOrg:
		return "AtOrg"
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
	case dsRewound:
		if dir == Next {
			ds.pos = 0
		} else { // Prev
			ds.pos = dsState(len(ds.rows) - 1)
		}
	case dsAtEnd:
		if dir == Next {
			panic("QuerySource.Get: Next at end")
		}
		ds.pos = dsState(len(ds.rows) - 1)
	case dsAtOrg:
		if dir == Prev {
			panic("QuerySource.Get: Prev at beginning")
		}
		ds.pos = 0
	default: // within
		if dir == Next {
			ds.pos++
		} else { // Prev
			ds.pos--
		}
	}
	if ds.pos < 0 {
		ds.pos = dsAtOrg
		return nil
	}
	if int(ds.pos) >= len(ds.rows) {
		ds.pos = dsAtEnd
		return nil
	}
	return ds.rows[ds.pos]
}

//-------------------------------------------------------------------

// QuerySource is a mock Query source for fuzzing.
// It generates rows of data that can be filtered with Select or Lookup.
type QuerySource struct {
	QueryMock
	index []string
	dataSource
	selCols []string
	selVals []string
}

type buildQS struct {
	rnd         *rand.Rand
	maxRows     int
	maxKeys     int
	maxIndexes  int
	columns     []string
	colIndex    map[string]int
	emptyKey    bool
	keys        [][]string
	indexes     [][]string
	cardinality map[string]int
	rows        []Row
	x           uint16
}

func NewQuerySource(rnd *rand.Rand) *QuerySource {
	return newQuerySource(rnd, 101, 3, 3)
}

func newQuerySource(rnd *rand.Rand, maxRows, maxKeys, maxIndexes int) *QuerySource {
	b := buildQS{rnd: rnd, maxRows: maxRows, maxKeys: maxKeys, maxIndexes: maxIndexes}
	b.makeColumns()
	b.makeKeys()
	b.makeIndexes()
	b.makeRows()
	nrows := len(b.rows) // actual number of rows generated
	qs := &QuerySource{}
	qs.rows = b.rows
	qs.pos = dsRewound
	qs.HeaderResult = SimpleHeader(b.columns)
	qs.ColumnsResult = b.columns
	qs.IndexesResult = b.indexes
	qs.KeysResult = b.keys
	qs.NrowsN = nrows
	qs.NrowsP = nrows
	qs.MetricsResult = &metrics{}
	qs.FixedResult = []Fixed{}
	qs.FastSingleResult = b.emptyKey
	qs.LookupLevels = 2 // ???
	qs.KnowExactNrowsResult = true
	return qs
}

func (qs *QuerySource) String() string {
	var sb strings.Builder
	sb.WriteString("QuerySource ")
	sb.WriteString(str.Join("(,)", qs.ColumnsResult))
	for _, k := range qs.KeysResult {
		sb.WriteString(" key(")
		sb.WriteString(str.Join(",", k))
		sb.WriteString(")")
	}
	for _, i := range qs.IndexesResult {
		sb.WriteString(" index(")
		sb.WriteString(str.Join(",", i))
		sb.WriteString(")")
	}
	return sb.String()
}

func (b *buildQS) makeColumns() {
	ncols := 1 + b.rnd.IntN(31)
	b.columns = make([]string, ncols)
	b.colIndex = make(map[string]int, ncols)
	b.cardinality = make(map[string]int, ncols)
	for i := range ncols {
		col := "c" + strconv.Itoa(i)
		b.columns[i] = col
		b.colIndex[col] = i
		b.cardinality[col] = 1 + b.rnd.IntN(1009)
	}
}

func (b *buildQS) makeKeys() {
	if b.rnd.IntN(11) == 5 {
		b.emptyKey = true
		b.keys = [][]string{{}}
		b.indexes = [][]string{{}}
		return
	}
	nkeys := 1 + b.rnd.IntN(b.maxKeys)
	// to simplify creating unique data, keys do not overlap
	p := b.rnd.Perm(len(b.columns))
	b.keys = make([][]string, 0, nkeys)
	for range nkeys {
		if len(p) == 0 {
			break
		}
		keylen := 1 + b.rnd.IntN(3) // 1=3 cols in key
		keylen = min(keylen, len(p))
		key := make([]string, keylen)
		for j := range keylen {
			key[j] = b.columns[p[0]]
			p = p[1:]
		}
		b.keys = append(b.keys, key)
	}
}

func (b *buildQS) makeIndexes() {
	if b.emptyKey || len(b.columns) < 2 {
		b.indexes = b.keys
		return
	}
	nindexes := b.rnd.IntN(b.maxIndexes)
	b.indexes = make([][]string, 0, len(b.keys)+nindexes)
	b.indexes = append(b.indexes, b.keys...) // keys are indexes
	maxcols := min(nindexes, len(b.columns))
	for ncols := 1; ncols < maxcols; ncols++  {
		idx := set.RandPerm(b.rnd, b.columns, ncols)
		if !slices.ContainsFunc(b.indexes,
			func(x []string) bool { return slices.Equal(x, idx) }) {
			b.indexes = append(b.indexes, idx)
		}
	}
}

func (b *buildQS) makeRows() {
	var nrows int
	if b.emptyKey {
		nrows = b.rnd.IntN(1) // 0 or 1
	} else {
		nrows = b.rnd.IntN(b.maxRows)
	}
	b.rows = make([]Row, nrows)
	b.x = uint16(b.rnd.Int())
	for i := range nrows {
		b.rows[i] = Row{DbRec{Record: b.makeRow(), Off: uint64(i)}}
	}
}

func (b *buildQS) makeRow() Record {
	vals := make([]string, len(b.columns))
	// generate data for all the columns
	for i, col := range b.columns {
		vals[i] = col + "_" + strconv.Itoa(b.rnd.IntN(b.cardinality[col]))
	}
	// overwrite with unique values for keys
	for _, key := range b.keys {
		n := b.x
		b.x = bits.Shuffle16(b.x) // shuffle ensures unique key values
		for i, col := range key {
			// split n (a unique value) over the columns of the key
			var v uint16
			if i < len(key)-1 {
				// 4 bits = 0 - 15 gives chance of duplicates
				v = n & 0b1111
				n >>= 4
			} else {
				// last column gets the rest
				v = n
			}
			vals[b.colIndex[col]] = col + "_" + strconv.Itoa(int(v))
			// fmt.Printf("Key %v Col %s n=%d v=%d\n", key, col, n, v)
		}
	}
	var rb RecordBuilder
	for _, val := range vals {
		rb.Add(SuStr(val))
	}
	return rb.Build()
}

//-------------------------------------------------------------------

func (qs *QuerySource) optimize(_ Mode, index []string, frac float64) (Cost, Cost, any) {
	if !qs.validIndex(index) {
		return impossible, impossible, nil
	}
	varcost := Cost(float64(10000) * frac)
	if varcost < 1 {
		varcost = 1
	}
	return 0, varcost, nil
}

func (qs *QuerySource) setApproach(index []string, _ float64, _ any, _ QueryTran) {
	if index == nil {
		qs.index = qs.IndexesResult[0]
		return
	}
	if !qs.validIndex(index) {
		panic("QuerySource.setApproach: invalid index")
	}
	qs.index = index
	slices.SortFunc(qs.rows, func(a, b Row) int {
		for _, col := range index {
			av := a.GetRaw(qs.HeaderResult, col)
			bv := b.GetRaw(qs.HeaderResult, col)
			if cmp := strings.Compare(av, bv); cmp != 0 {
				return cmp
			}
		}
		return 0
	})
	qs.Rewind()
}

func (qs *QuerySource) validIndex(index []string) bool {
	// If source has empty key, accept any index
	if isEmptyKey(qs.KeysResult) {
		return true
	}
	// Normal case: check if index is a prefix of any existing index
	return slices.ContainsFunc(qs.IndexesResult, func(ix []string) bool {
		return slc.HasPrefix(ix, index)
	})
}

func (qs *QuerySource) knowExactNrows() bool {
	return true
}

// Rewind resets the position so Next gets first or Prev gets last.
func (qs *QuerySource) Rewind() {
	qs.rewind()
}

// Get returns the next or previous row, respecting any active Select.
func (qs *QuerySource) Get(_ *Thread, dir Dir) Row {
	for {
		row := qs.get(dir)
		if row == nil {
			return nil
		}
		if qs.matches(row) {
			return row
		}
	}
}

// matches checks if a row matches the current Select criteria.
// It only checks the index prefix columns (stopping at first missing column),
// matching the behavior of Table and TempIndex.
func (qs *QuerySource) matches(row Row) bool {
	if qs.selCols == nil {
		return true
	}

	// Check if source has empty keys (singleton)
	hasEmptyKey := false
	for _, key := range qs.KeysResult {
		if len(key) == 0 {
			hasEmptyKey = true
			break
		}
	}

	if hasEmptyKey {
		// Use singletonFilter for empty key case
		return singletonFilter(qs.HeaderResult, row, qs.selCols, qs.selVals)
	}

	// Normal case: check index prefix columns
	for _, col := range qs.index {
		i := slices.Index(qs.selCols, col)
		if i == -1 {
			break // stop at first missing column
		}
		raw := row.GetRaw(qs.HeaderResult, col)
		if raw != qs.selVals[i] {
			return false
		}
	}
	return true
}

// Select restricts subsequent Get calls to rows matching the given columns/values.
// cols and vals are packed values (output of Pack).
func (qs *QuerySource) Select(cols, vals []string) {
	if cols != nil {
		assert.That(!selConflict(qs.ColumnsResult, cols, vals))
		assert.That(selPrefix(qs.index, cols))
	}
	qs.selCols = cols
	qs.selVals = vals
	qs.Rewind()
}

// selPrefix does the same validation as Table
func selPrefix(index, selCols []string) bool {
	if len(index) == 0 {
		return true
	}
	data := false
	for _, col := range index {
		if slices.Index(selCols, col) == -1 {
			break
		}
		data = true
	}
	return data
}

// Lookup finds and returns the first row matching the given columns/values.
func (qs *QuerySource) Lookup(_ *Thread, cols, vals []string) Row {
	qs.Select(cols, vals)
	row := qs.Get(nil, Next)
	assert.That(row == nil || qs.Get(nil, Next) == nil)
	qs.Select(nil, nil) // clear select
	return row
}

// Simple returns all rows.
func (qs *QuerySource) Simple(*Thread) []Row {
	// need to copy because it may be modified
	return slices.Clone(qs.rows)
}

//-------------------------------------------------------------------

func TestQuerySource(t *testing.T) {
	rnd := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	for range 1000 {
		qs := NewQuerySource(rnd)
		checkNoDuplicateIndexes(t, qs.IndexesResult)
		checkKeyDataUnique(t, qs)
	}
}

func checkNoDuplicateIndexes(t *testing.T, indexes [][]string) {
	seen := make(map[string]bool)
	for _, idx := range indexes {
		key := strings.Join(idx, ",")
		if seen[key] {
			t.Fatalf("duplicate index: %v", idx)
		}
		seen[key] = true
	}
}

func checkKeyDataUnique(t *testing.T, qs *QuerySource) {
	for _, key := range qs.KeysResult {
		seen := make(map[string]int)
		for i, row := range qs.rows {
			var vals []string
			for _, col := range key {
				vals = append(vals, row.GetRaw(qs.HeaderResult, col))
			}
			combo := strings.Join(vals, "|")
			if prev, ok := seen[combo]; ok {
				t.Logf("Duplicate found for key %v", key)
				t.Logf("Row %d: %s", prev, combo)
				t.Logf("Row %d: %s", i, combo)
				t.Fatalf("duplicate key combination for key %v: %s", key, combo)
			}
			seen[combo] = i
		}
	}
}

func TestQuerySourceSetApproach(t *testing.T) {
	rnd := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	for range 100 {
		qs := NewQuerySource(rnd)
		for _, index := range qs.IndexesResult {
			qs.setApproach(index, 0, nil, nil)
			checkSorted(t, qs, index)
		}
	}
}

func checkSorted(t *testing.T, qs *QuerySource, index []string) {
	for i := 1; i < len(qs.rows); i++ {
		prev := qs.rows[i-1]
		curr := qs.rows[i]
		cmp := 0
		for _, col := range index {
			pv := prev.GetRaw(qs.HeaderResult, col)
			cv := curr.GetRaw(qs.HeaderResult, col)
			cmp = strings.Compare(pv, cv)
			if cmp != 0 {
				break
			}
		}
		if cmp > 0 {
			t.Fatalf("rows not sorted by %v at position %d", index, i)
		}
	}
}
