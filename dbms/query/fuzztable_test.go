// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bit"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

type FT struct {
	db       *db19.Database
	rnd      *rand.Rand
	nextNum  int
	rt       *db19.ReadTran
	nextRule int
	rules    map[string]string // for debugging
	ruleRefs []string
	ruleDeps map[string]string
}

type buildFT struct {
	*FT
	maxRows    int
	maxKeys    int
	maxIndexes int
	prefix     string
	// colflds is the physical record slots: base columns plus derived (rule)
	// columns, with dropped columns replaced by "-". data, colIndex, and
	// cardinality are keyed/positioned on it.
	// Use columns() for the logical view.
	colflds  []string
	colIndex map[string]int
	emptyKey bool
	keys     [][]string
	indexes  [][]string
	// uniqueIdx records which indexes are "unique index" rather than plain
	// "index", keyed by the joined column names (indexKey) so the mapping
	// stays correct even if indexes are later split, reordered, or shared
	// (e.g. by fuzz tests that build related tables from a common buildFT).
	uniqueIdx   map[string]bool
	cardinality map[string]int
	data        [][]string
	noEmptyKey  bool
	ruleCols    []string
}

// columns returns the logical columns: base columns plus derived (rule)
// columns, with dropped ("-") columns removed and derived names uncapitalized,
// matching Query.Columns(). It is derived from colflds and must not be
// modified by callers.
func (b *buildFT) columns() []string {
	cols := slices.DeleteFunc(slices.Clone(b.colflds),
		func(c string) bool { return c == "-" })
	for i, c := range cols {
		cols[i] = str.UnCapitalize(c)
	}
	return cols
}

// fldIndex returns the physical field index for a column name
func (b *buildFT) fldIndex(col string) int {
	if i, ok := b.colIndex[col]; ok {
		return i
	}
	panic("buildFT.fldIndex: unknown column: " + col)
}

func (ft *FT) NewFuzzTable() Query {
	return ft.newFT().Build()
}

const (
	ftMaxRows    = 47
	ftMaxCols    = 31
	ftMaxRules   = 3
	ftMaxKeys    = 3
	ftMaxIndexes = 3
	ftMaxDel     = 3
)

func (ft *FT) newFT() *buildFT {
	if ft.rules == nil {
		ft.rules = make(map[string]string)
	}
	if ft.ruleDeps == nil {
		ft.ruleDeps = make(map[string]string)
	}
	return &buildFT{FT: ft, maxRows: ftMaxRows, maxKeys: ftMaxKeys, maxIndexes: ftMaxIndexes, prefix: "c"}
}

func (b *buildFT) Sizes(maxRows, maxKeys, maxIndexes int) *buildFT {
	b.maxRows = maxRows
	b.maxKeys = maxKeys
	b.maxIndexes = maxIndexes
	return b
}

func (b *buildFT) Prefix(prefix string) *buildFT {
	b.prefix = prefix
	return b
}

func (b *buildFT) NoEmptyKey() *buildFT {
	b.noEmptyKey = true
	return b
}

func (b *buildFT) Build() Query {
	b.construct()
	return b.finish()
}

func (b *buildFT) construct() *buildFT {
	b.makeColumns()
	b.makeKeys()
	b.makeIndexes()
	b.makeData()
	b.dropColumns()
	return b
}

func (b *buildFT) makeRules() {
	base := b.colflds[:len(b.colflds)-len(b.ruleCols)]
	nAvail := len(base) - slc.Count(base, "-")
	pickAvail := func() string {
		k := b.rnd.IntN(nAvail)
		for _, c := range base {
			if c != "-" {
				if k == 0 {
					return c
				}
				k--
			}
		}
		panic(assert.ShouldNotReachHere())
	}
	for _, col := range b.ruleCols {
		ruleName := "Rule_" + str.UnCapitalize(col)
		var ruleExpr string
		if nAvail > 0 && b.rnd.IntN(2) == 0 {
			existingCol := pickAvail()
			b.ruleRefs = append(b.ruleRefs, existingCol)
			b.ruleDeps[str.UnCapitalize(col)] = existingCol
			ruleExpr = "." + str.UnCapitalize(existingCol)
		} else {
			ruleExpr = strconv.Itoa(b.rnd.IntN(1000))
		}
		b.rules[ruleName] = ruleExpr
		Global.TestDef(ruleName,
			compile.Constant("function() { return "+ruleExpr+" }"))
	}
}

func (b *buildFT) finish() Query {
	b.makeRules()
	table := "table" + strconv.Itoa(b.nextNum)
	b.nextNum++
	var sb strings.Builder
	sb.WriteString("create ")
	sb.WriteString(table)
	sb.WriteString(str.Join("(,)", b.colflds))
	for _, key := range b.keys {
		sb.WriteString(" key")
		sb.WriteString(str.Join("(,)", key))
	}
	for _, index := range b.indexes {
		sb.WriteString(" index")
		if b.uniqueIdx[indexKey(index)] {
			sb.WriteString(" unique")
		}
		sb.WriteString(str.Join("(,)", index))
	}
	DoAdminTest(b.db, sb.String())
	// fmt.Println(sb.String())

	// output the data
	ut := b.db.NewUpdateTran(AllPerms)
	for _, vals := range b.data {
		ut.Output(nil, table, b.dataToRecord(vals))
	}
	b.db.CommitMerge(ut)

	// recreate read tran so it includes the new table
	b.rt = b.db.NewReadTran(AllPerms)
	return NewTable(b.rt, table)
}

func (b *buildFT) makeColumns() {
	ncols := 1 + b.rnd.IntN(ftMaxCols)
	b.colflds = make([]string, ncols)
	b.colIndex = make(map[string]int, ncols)
	b.cardinality = make(map[string]int, ncols)
	for i := range ncols {
		col := b.prefix + strconv.Itoa(i)
		b.colflds[i] = col
		b.colIndex[col] = i
		b.cardinality[col] = 1 + b.rnd.IntN(1009)
	}

	if !b.noEmptyKey && b.rnd.IntN(2) == 0 {
		nrule := 1 + b.rnd.IntN(ftMaxRules)
		b.ruleCols = make([]string, 0, nrule)
		for range nrule {
			i := b.nextRule
			b.nextRule++
			// capitalized so the column is derived (not stored), only from a rule
			ruleCol := str.Capitalize(b.prefix) + "r" + strconv.Itoa(i)
			b.ruleCols = append(b.ruleCols, ruleCol)
			b.colflds = append(b.colflds, ruleCol)
			b.colIndex[ruleCol] = len(b.colflds) - 1
		}
	}
}

func (b *buildFT) makeKeys() {
	if !b.noEmptyKey && b.rnd.IntN(11) == 5 {
		b.emptyKey = true
		b.keys = emptyKey
		return
	}
	nkeys := 1 + b.rnd.IntN(b.maxKeys)

	nonRuleCols := make([]string, 0, len(b.colflds))
	for _, col := range b.colflds {
		if !slices.Contains(b.ruleCols, col) {
			nonRuleCols = append(nonRuleCols, col)
		}
	}

	// to simplify creating unique data, keys do not overlap
	p := b.rnd.Perm(len(nonRuleCols))
	b.keys = make([][]string, 0, nkeys)
	for range nkeys {
		if len(p) == 0 {
			break
		}
		keylen := 1 + b.rnd.IntN(3) // 1=3 cols in key
		keylen = min(keylen, len(p))
		key := make([]string, keylen)
		for j := range keylen {
			key[j] = nonRuleCols[p[0]]
			p = p[1:]
		}
		b.keys = append(b.keys, key)
	}
}

func (b *buildFT) makeIndexes() {
	if b.emptyKey || len(b.colflds) < 2 {
		return
	}
	nonRuleCols := make([]string, 0, len(b.colflds))
	for _, col := range b.colflds {
		if !slices.Contains(b.ruleCols, col) {
			nonRuleCols = append(nonRuleCols, col)
		}
	}
	if len(nonRuleCols) < 2 {
		return
	}
	keyCols := make([]string, 0, 8)
	for _, key := range b.keys {
		keyCols = append(keyCols, key...)
	}
	var uniqueCols []string // columns already claimed by a unique index
	nindexes := b.rnd.IntN(b.maxIndexes)
	b.indexes = make([][]string, 0, nindexes)
	b.uniqueIdx = make(map[string]bool)
	maxcols := min(nindexes, len(nonRuleCols))
	for ncols := 1; ncols < maxcols; ncols++ {
		idx := set.RandPerm(b.rnd, nonRuleCols, ncols)
		if slc.ContainsFn(b.indexes, idx, slices.Equal) ||
			slc.ContainsFn(b.keys, idx, containsKey) {
			continue
		}
		b.indexes = append(b.indexes, idx)
		// occasionally make it a unique index rather than a plain index.
		// keep unique index columns disjoint from keys and from each other
		// so generating the data (mostly unique, sometimes empty) stays simple.
		isUnique := b.rnd.IntN(3) == 0 &&
			set.Disjoint(idx, keyCols) &&
			set.Disjoint(idx, uniqueCols)
		if isUnique {
			uniqueCols = append(uniqueCols, idx...)
			b.uniqueIdx[indexKey(idx)] = true
		}
	}
}

func (b *buildFT) makeData() {
	var nrows int
	if b.emptyKey {
		nrows = b.rnd.IntN(2) // 0 or 1
	} else {
		nrows = b.rnd.IntN(b.maxRows)
	}
	b.data = b.makeRowsData(nrows)
}

func (b *buildFT) makeRowsData(nrows int) [][]string {
	x := uint16(b.rnd.Int())

	keyCols := make(map[string]bool)
	for _, key := range b.keys {
		for _, col := range key {
			keyCols[col] = true
		}
	}
	// columns belonging to a plain (non-unique) index occasionally get
	// empty data too, just to exercise indexes with missing values.
	regularIndexCols := make(map[string]bool)
	for _, index := range b.indexes {
		if !b.uniqueIdx[indexKey(index)] {
			for _, col := range index {
				regularIndexCols[col] = true
			}
		}
	}

	data := make([][]string, nrows)
	for i := range nrows {
		vals := make([]string, len(b.colflds))
		// generate data for all the columns
		for j, col := range b.colflds {
			if slices.Contains(b.ruleCols, col) {
				continue
			}
			if regularIndexCols[col] && !keyCols[col] && b.rnd.IntN(10) == 0 {
				vals[j] = "" // leave some regular index data empty
				continue
			}
			vals[j] = col + "_" + strconv.Itoa(b.rnd.IntN(b.cardinality[col]))
		}
		// overwrite with unique values for keys
		for _, key := range b.keys {
			n := x
			x = bit.Shuffle16(x) // shuffle ensures unique key values
			for k, col := range key {
				// split n (a unique value) over the columns of the key
				var v uint16
				if k < len(key)-1 {
					// 4 bits = 0 - 15 gives chance of duplicates
					v = n & 0b1111
					n >>= 4
				} else {
					// last column gets the rest
					v = n
				}
				vals[b.colIndex[col]] = col + "_" + strconv.Itoa(int(v))
			}
		}
		// overwrite with unique values for unique indexes, except sometimes
		// leave them entirely empty - a unique index allows multiple records
		// with no value, unlike a key.
		for _, index := range b.indexes {
			if !b.uniqueIdx[indexKey(index)] {
				continue
			}
			if b.rnd.IntN(5) == 0 {
				for _, col := range index {
					vals[b.colIndex[col]] = ""
				}
				continue
			}
			n := x
			x = bit.Shuffle16(x)
			for k, col := range index {
				var v uint16
				if k < len(index)-1 {
					v = n & 0b1111
					n >>= 4
				} else {
					v = n
				}
				vals[b.colIndex[col]] = col + "_" + strconv.Itoa(int(v))
			}
		}
		data[i] = vals
	}
	return data
}

func (b *buildFT) dataToRecord(vals []string) Record {
	var rb RecordBuilder
	for i, val := range vals {
		if i < len(b.colflds) && slices.Contains(b.ruleCols, b.colflds[i]) {
			continue
		}
		rb.Add(SuStr(val))
	}
	return rb.Build()
}

func (b *buildFT) dropColumns() {
	nDel := b.rnd.IntN(ftMaxDel)
	for range nDel {
		i := b.rnd.IntN(len(b.colflds))
		col := b.colflds[i]
		if b.isUsed(col) {
			continue
		}
		delete(b.colIndex, col)
		b.colflds[i] = "-"
	}
}

func (b *buildFT) isUsed(col string) bool {
	for _, key := range b.keys {
		if slices.Contains(key, col) {
			return true
		}
	}
	for _, idx := range b.indexes {
		if slices.Contains(idx, col) {
			return true
		}
	}
	return slices.Contains(b.ruleCols, col)
}

//-------------------------------------------------------------------
// Tests for FuzzTable

func testFT() *FT {
	return newFT(rand.Uint64(), rand.Uint64())
}

func newFT(seed1, seed2 uint64) *FT {
	st := stor.HeapStor(64 * 1024)
	st.Alloc(1)
	db := db19.CreateDb(st)
	// db19.StartConcur(db, 50*time.Millisecond)
	db.CheckerSync()
	return &FT{
		db:  db,
		rnd: rand.New(rand.NewPCG(seed1, seed2)),
	}
}

func TestNewFuzzTable(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	q := ft.NewFuzzTable()
	if q == nil {
		t.Fatal("NewFuzzTable returned nil")
	}
	hdr := q.Header()
	if len(hdr.Columns) == 0 {
		t.Error("Table has no columns")
	}
}

func TestFuzzTable_BuilderMethods(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	// Test Sizes
	bf := ft.newFT().Sizes(50, 5, 4)
	if bf.maxRows != 50 || bf.maxKeys != 5 || bf.maxIndexes != 4 {
		t.Error("Sizes not set correctly")
	}

	// Test Prefix
	bf = ft.newFT().Prefix("test_")
	if bf.prefix != "test_" {
		t.Errorf("Expected prefix='test_', got '%s'", bf.prefix)
	}

	// Test NoEmptyKey
	bf = ft.newFT().NoEmptyKey()
	if !bf.noEmptyKey {
		t.Error("Expected noEmptyKey=true")
	}
}

func TestFuzzTable_Build(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	q := ft.newFT().Sizes(10, 2, 2).Prefix("col_").Build()
	q.SetTran(ft.db.NewReadTran(AllPerms))
	hdr := q.Header()
	if len(hdr.Columns) == 0 {
		t.Error("Table has no columns")
	}
	for _, col := range hdr.Columns {
		if !strings.HasPrefix(col, "col_") {
			t.Errorf("Column %s doesn't have expected prefix", col)
		}
	}
	q.Rewind()
	row := q.Get(nil, Next)
	if row.Eof() && len(hdr.Columns) > 0 {
		t.Log("Table produced no rows (may be valid if emptyKey)")
	}
}

func TestFuzzTable_dropColumns(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	// deterministic: build a table, then drop one column
	bf := ft.newFT().Sizes(10, 2, 2)
	bf.colflds = []string{"a", "b", "c", "d"}
	bf.colIndex = map[string]int{"a": 0, "b": 1, "c": 2, "d": 3}
	bf.cardinality = map[string]int{"a": 10, "b": 10, "c": 10, "d": 10}
	bf.keys = [][]string{{"a"}}
	bf.indexes = [][]string{{"b"}}
	bf.ruleCols = []string{}
	bf.dropColumns()

	// "a" and "b" are in key/index, can't be dropped
	for _, col := range []string{"a", "b"} {
		if bf.colflds[bf.colIndex[col]] != col {
			t.Errorf("column %q was dropped but it's in a key/index", col)
		}
	}
	// "c" or "d" (or both) may have been dropped
	nDel := 0
	for _, col := range bf.colflds {
		if col == "-" {
			nDel++
		}
	}
	if nDel > 2 {
		t.Errorf("dropped %d columns, expected at most 2", nDel)
	}

	// build the table and verify Fields vs Columns
	bf.data = [][]string{{"a_1", "b_2", "c_3", "d_4"}}
	q := bf.finish()
	hdr := q.Header()
	nFields := len(hdr.Fields[0])
	nCols := len(hdr.Columns)
	if nFields != 4 {
		t.Fatalf("expected 4 physical fields, got %d", nFields)
	}
	if nCols == 4 && nDel > 0 {
		t.Errorf("dropped columns should not appear in Columns: %v", hdr.Columns)
	}
	for _, col := range hdr.Columns {
		if col == "-" {
			t.Errorf("deleted column should not appear in Columns: %v", hdr.Columns)
		}
	}
}

func TestFuzzTable_makeColumns(t *testing.T) {
	for range 100 {
		ft := testFT()
		bf := ft.newFT()
		bf.makeColumns()

		if len(bf.colflds) == 0 {
			t.Fatal("No columns generated")
		}
		if len(bf.colflds) > ftMaxCols+ftMaxRules {
			t.Errorf("Too many columns: %d", len(bf.colflds))
		}
		if len(bf.colIndex) != len(bf.colflds) {
			t.Errorf("colIndex size mismatch")
		}
		for col, card := range bf.cardinality {
			if card < 1 || card > 1009 {
				t.Errorf("Invalid cardinality for %s: %d", col, card)
			}
		}
	}
}

func TestFuzzTable_makeKeys(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	bf := ft.newFT()
	bf.colflds = []string{"a", "b", "c"}
	bf.colIndex = map[string]int{"a": 0, "b": 1, "c": 2}
	bf.makeKeys()

	for _, key := range bf.keys {
		for _, col := range key {
			if !slices.Contains(bf.colflds, col) {
				t.Errorf("Key column %s not in columns", col)
			}
		}
	}
}

func TestFuzzTable_makeIndexes(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	bf := ft.newFT()
	bf.colflds = []string{"a", "b", "c"}
	bf.colIndex = map[string]int{"a": 0, "b": 1, "c": 2}
	bf.keys = [][]string{{"a"}}
	bf.makeIndexes()

	// Verify indexes don't overlap with keys
	for _, idx := range bf.indexes {
		if slc.ContainsFn(bf.keys, idx, containsKey) {
			t.Errorf("Index %v overlaps with a key", idx)
		}
	}

	// Verify no duplicate indexes
	for i, idx1 := range bf.indexes {
		for j, idx2 := range bf.indexes {
			if i < j && slices.Equal(idx1, idx2) {
				t.Errorf("Duplicate indexes: %v", idx1)
			}
		}
	}

	// Verify all indexes contain only valid columns
	for _, idx := range bf.indexes {
		for _, col := range idx {
			if !slices.Contains(bf.colflds, col) {
				t.Errorf("Index contains invalid column: %s", col)
			}
		}
	}
}

func containsKey(key, idx []string) bool {
	return set.HasSubset(idx, key)
}

// indexKey returns a stable map key for an index's columns, used to track
// which indexes are unique independent of their position/order in a slice.
func indexKey(idx []string) string {
	return strings.Join(idx, ",")
}

func TestFuzzTable_makeData(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	bf := ft.newFT()
	bf.colflds = []string{"c0", "c1", "c2", "c3", "c4", "c5"}
	bf.colIndex = map[string]int{"c0": 0, "c1": 1, "c2": 2, "c3": 3, "c4": 4, "c5": 5}
	bf.keys = [][]string{{"c0"}, {"c1", "c2"}}
	bf.cardinality = map[string]int{"c0": 10, "c1": 3, "c2": 50, "c3": 200, "c4": 5, "c5": 8}
	bf.makeData()

	// single-column key c0 values should be unique
	seen := map[string]int{}
	for i, row := range bf.data {
		if len(row) != len(bf.colflds) {
			t.Errorf("Row %d has %d columns, expected %d", i, len(row), len(bf.colflds))
		}
		v := row[bf.colIndex["c0"]]
		if prev, ok := seen[v]; ok {
			t.Errorf("single-col key c0 duplicate value %q at rows %d and %d", v, prev, i)
		}
		seen[v] = i
	}
	// multi-column key (c1,c2) combinations should be unique
	seenCombo := map[string]int{}
	for i, row := range bf.data {
		v1 := row[bf.colIndex["c1"]]
		v2 := row[bf.colIndex["c2"]]
		combo := v1 + "|" + v2
		if prev, ok := seenCombo[combo]; ok {
			t.Errorf("multi-col key (c1,c2) duplicate combo %q at rows %d and %d", combo, prev, i)
		}
		seenCombo[combo] = i
	}
	// non-key column values should be within cardinality
	nonKeys := []string{"c3", "c4", "c5"}
	for _, col := range nonKeys {
		for i, row := range bf.data {
			v := row[bf.colIndex[col]]
			prefix := col + "_"
			suffix := strings.TrimPrefix(v, prefix)
			if suffix == v || suffix == "" {
				t.Errorf("row %d: %s value %q missing expected prefix", i, col, v)
				continue
			}
			n, err := strconv.Atoi(suffix)
			if err != nil {
				t.Errorf("row %d: %s value %q suffix not an integer", i, col, v)
				continue
			}
			limit := bf.cardinality[col]
			if n < 0 || n >= limit {
				t.Errorf("row %d: %s value %q suffix %d out of range [0,%d)", i, col, v, n, limit)
			}
		}
	}
}
