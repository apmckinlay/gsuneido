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
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

// QuerySource is a mock Query source for fuzzing.
// It generates rows of data that can be filtered with Select or Lookup.
type QuerySource struct {
	QueryMock
	index []string
	dataSource
	selCols []string
	selVals []string
}

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

func NewDataSource(rows []Row) *dataSource {
	return &dataSource{rows: rows, pos: dsRewound}
}

// NewQuerySource creates a QuerySource with random columns, rows, keys, and indexes.
// The first key column has unique values.
// Non-key columns have duplicates (values cycle with period 10).
// Sometimes (1/10 chance) returns a QuerySource with a single empty key (no columns).
func NewQuerySource(rnd *rand.Rand) *QuerySource {
	ncols := 1 + rnd.IntN(30) // 1-30 columns
	columns := make([]string, ncols)
	for i := range columns {
		columns[i] = "c" + strconv.Itoa(i)
	}

	// 1/10 chance to use empty key instead of normal keys
	useEmptyKey := rnd.IntN(10) == 0

	var nrows int
	var keys, indexes [][]string
	if useEmptyKey {
		// Create a single empty key (no columns) and single empty index
		keys = [][]string{{}}
		indexes = [][]string{{}}
		nrows = rnd.IntN(2) // 0 or 1 rows for empty key
	} else {
		nindexes := 1 + rnd.IntN(8) // 1-8 indexes
		// Limit indexes based on number of columns to avoid infinite loop
		maxIndexes := ncols * ncols // rough upper bound
		nindexes = min(nindexes, maxIndexes)
		indexes = randomIndexes(rnd, columns, nindexes)
		nkeys := 1 + rnd.IntN(3) // 1-3 keys
		nkeys = min(nkeys, len(indexes))
		perm := rnd.Perm(len(indexes))
		keys = make([][]string, nkeys)
		for i := range nkeys {
			keys[i] = indexes[perm[i]]
		}
		nrows = rnd.IntN(101) // 0-100 rows
	}

	// Build column index map
	colIndex := make(map[string]int)
	for i, col := range columns {
		colIndex[col] = i
	}

	// Compute min key length for each column (use min because single-col keys need high cardinality)
	minKeyLen := make([]int, ncols)
	for j := range ncols {
		minKeyLen[j] = 999 // sentinel for "not in any key"
	}
	for _, key := range keys {
		for _, col := range key {
			j := colIndex[col]
			minKeyLen[j] = min(minKeyLen[j], len(key))
		}
	}

	// Compute cardinality for each column based on min key length
	// Use higher cardinality to reduce rejection rate
	cardinality := make([]int, ncols)
	for j := range ncols {
		switch minKeyLen[j] {
		case 999:
			cardinality[j] = 100 // non-key column
		case 1:
			cardinality[j] = max(2*nrows, 100) // single-column key needs high cardinality
		default:
			// Use 2x nrows to ensure enough unique combinations
			cardinality[j] = max(3, intRoot(2*nrows, minKeyLen[j]))
		}
	}

	// Track seen key combinations for rejection
	seenKeys := make([]map[string]bool, len(keys))
	for k := range keys {
		seenKeys[k] = make(map[string]bool)
	}

	// Generate row data with random values, rejecting duplicate key combinations
	fields := slices.Clone(columns)
	header := NewHeader([][]string{fields}, columns)
	rows := make([]Row, 0, nrows)
	maxAttempts := nrows * 10
	attempts := 0
	for len(rows) < nrows && attempts < maxAttempts {
		attempts++
		vals := make([]int, ncols)
		for j := range ncols {
			vals[j] = rnd.IntN(cardinality[j])
		}
		// Check if any key combination is a duplicate
		duplicate := false
		for k, key := range keys {
			combo := keyCombo(key, colIndex, vals)
			if seenKeys[k][combo] {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}
		// Mark all key combinations as seen
		for k, key := range keys {
			combo := keyCombo(key, colIndex, vals)
			seenKeys[k][combo] = true
		}
		// Build the row
		var rb RecordBuilder
		for j, col := range columns {
			rb.Add(SuStr(col + "_" + strconv.Itoa(vals[j])))
		}
		rows = append(rows, Row{DbRec{Record: rb.Build(), Off: uint64(len(rows))}})
	}
	nrows = len(rows) // actual number of rows generated
	qs := &QuerySource{}
	qs.rows = rows
	qs.pos = dsRewound
	qs.HeaderResult = header
	qs.ColumnsResult = columns
	qs.IndexesResult = indexes
	qs.KeysResult = keys
	qs.NrowsN = nrows
	qs.NrowsP = nrows
	qs.MetricsResult = &metrics{}
	qs.FixedResult = []Fixed{}
	return qs
}

// keyCombo returns a string key for the combination of values in the given columns.
func keyCombo(key []string, colIndex map[string]int, vals []int) string {
	var sb strings.Builder
	for i, col := range key {
		if i > 0 {
			sb.WriteByte('|')
		}
		sb.WriteString(strconv.Itoa(vals[colIndex[col]]))
	}
	return sb.String()
}

// intRoot returns the smallest base such that base^exp >= n
func intRoot(n, exp int) int {
	if n <= 1 {
		return 2
	}
	base := 2
	for pow(base, exp) < n {
		base++
	}
	return base
}

func pow(base, exp int) int {
	result := 1
	for range exp {
		result *= base
	}
	return result
}

// randomIndexes generates n unique random indexes from the given columns.
func randomIndexes(rnd *rand.Rand, columns []string, n int) [][]string {
	seen := make(map[string]bool)
	var result [][]string
	attempts := 0
	for len(result) < n && attempts < n*100 {
		attempts++
		size := 1 + rnd.IntN(min(len(columns), 5)) // 1-5 columns per index
		perm := rnd.Perm(len(columns))
		idx := make([]string, size)
		for j := range size {
			idx[j] = columns[perm[j]]
		}
		key := strings.Join(idx, ",")
		if !seen[key] {
			seen[key] = true
			result = append(result, idx)
		}
	}
	return result
}

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
func (qs *QuerySource) matches(row Row) bool {
	if qs.selCols == nil {
		return true
	}
	for i, col := range qs.selCols {
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
	for range 100 {
		qs := NewQuerySource(rnd)
		checkNoDuplicateIndexes(t, qs.IndexesResult)
		checkKeyDataUnique(t, qs)
	}
}

func checkNoDuplicateIndexes(t *testing.T, indexes [][]string) {
	t.Helper()
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
	t.Helper()
	for _, key := range qs.KeysResult {
		seen := make(map[string]bool)
		for _, row := range qs.rows {
			var vals []string
			for _, col := range key {
				vals = append(vals, row.GetRaw(qs.HeaderResult, col))
			}
			combo := strings.Join(vals, "|")
			if seen[combo] {
				t.Fatalf("duplicate key combination for key %v: %s", key, combo)
			}
			seen[combo] = true
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
	t.Helper()
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
