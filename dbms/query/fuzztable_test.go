// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"

	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/set"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/str"
)

type FT struct {
	db      *db19.Database
	rnd     *rand.Rand
	nextNum int
	rt      *db19.ReadTran
}

type buildFT struct {
	*FT
	maxRows     int
	maxKeys     int
	maxIndexes  int
	prefix      string
	columns     []string
	colIndex    map[string]int
	emptyKey    bool
	keys        [][]string
	indexes     [][]string
	cardinality map[string]int
	data        [][]string
	noEmptyKey  bool
}

func (ft *FT) NewFuzzTable() Query {
	return ft.newFT().Build()
}

func (ft *FT) newFT() *buildFT {
	return &buildFT{FT: ft, maxRows: 47, maxKeys: 3, maxIndexes: 3, prefix: "c"}
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
	return b
}

func (b *buildFT) finish() Query {
	table := "table" + strconv.Itoa(b.nextNum)
	b.nextNum++
	var sb strings.Builder
	sb.WriteString("create ")
	sb.WriteString(table)
	sb.WriteString(str.Join("(,)", b.columns))
	for _, key := range b.keys {
		sb.WriteString(" key")
		sb.WriteString(str.Join("(,)", key))
	}
	for _, index := range b.indexes {
		sb.WriteString(" index")
		sb.WriteString(str.Join("(,)", index))
	}
	DoAdmin(b.db, sb.String(), nil)
	// fmt.Println(sb.String())

	// output the data
	ut := b.db.NewUpdateTran()
	for _, vals := range b.data {
		ut.Output(nil, table, b.dataToRecord(vals))
	}
	b.db.CommitMerge(ut)

	// recreate read tran so it includes the new table
	b.rt = b.db.NewReadTran()
	return NewTable(b.rt, table)
}

func (b *buildFT) makeColumns() {
	ncols := 1 + b.rnd.IntN(31)
	b.columns = make([]string, ncols)
	b.colIndex = make(map[string]int, ncols)
	b.cardinality = make(map[string]int, ncols)
	for i := range ncols {
		col := b.prefix + strconv.Itoa(i)
		b.columns[i] = col
		b.colIndex[col] = i
		b.cardinality[col] = 1 + b.rnd.IntN(1009)
	}
}

func (b *buildFT) makeKeys() {
	if !b.noEmptyKey && b.rnd.IntN(11) == 5 {
		b.emptyKey = true
		b.keys = emptyKey
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

func (b *buildFT) makeIndexes() {
	if b.emptyKey || len(b.columns) < 2 {
		return
	}
	nindexes := b.rnd.IntN(b.maxIndexes)
	b.indexes = make([][]string, 0, nindexes)
	maxcols := min(nindexes, len(b.columns))
	for ncols := 1; ncols < maxcols; ncols++ {
		idx := set.RandPerm(b.rnd, b.columns, ncols)
		if slc.ContainsFn(b.indexes, idx, slices.Equal) ||
			slc.ContainsFn(b.keys, idx, containsKey) {
			continue
		}
		b.indexes = append(b.indexes, idx)
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
	data := make([][]string, nrows)
	for i := range nrows {
		vals := make([]string, len(b.columns))
		// generate data for all the columns
		for j, col := range b.columns {
			vals[j] = col + "_" + strconv.Itoa(b.rnd.IntN(b.cardinality[col]))
		}
		// overwrite with unique values for keys
		for _, key := range b.keys {
			n := x
			x = bits.Shuffle16(x) // shuffle ensures unique key values
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
		data[i] = vals
	}
	return data
}

func (b *buildFT) dataToRecord(vals []string) Record {
	var rb RecordBuilder
	for _, val := range vals {
		rb.Add(SuStr(val))
	}
	return rb.Build()
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
	q.SetTran(ft.db.NewReadTran())
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
	if row == nil && len(hdr.Columns) > 0 {
		t.Log("Table produced no rows (may be valid if emptyKey)")
	}
}

func TestFuzzTable_makeColumns(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	bf := ft.newFT()
	bf.makeColumns()

	if len(bf.columns) == 0 {
		t.Fatal("No columns generated")
	}
	if len(bf.columns) > 32 {
		t.Errorf("Too many columns: %d", len(bf.columns))
	}
	if len(bf.colIndex) != len(bf.columns) {
		t.Errorf("colIndex size mismatch")
	}
	for col, card := range bf.cardinality {
		if card < 1 || card > 1009 {
			t.Errorf("Invalid cardinality for %s: %d", col, card)
		}
	}
}

func TestFuzzTable_makeKeys(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	bf := ft.newFT()
	bf.columns = []string{"a", "b", "c"}
	bf.colIndex = map[string]int{"a": 0, "b": 1, "c": 2}
	bf.makeKeys()

	for _, key := range bf.keys {
		for _, col := range key {
			if !slices.Contains(bf.columns, col) {
				t.Errorf("Key column %s not in columns", col)
			}
		}
	}
}

func TestFuzzTable_makeIndexes(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	bf := ft.newFT()
	bf.columns = []string{"a", "b", "c"}
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
			if !slices.Contains(bf.columns, col) {
				t.Errorf("Index contains invalid column: %s", col)
			}
		}
	}
}

func containsKey(key, idx []string) bool {
	return set.Subset(idx, key)
}

func TestFuzzTable_makeData(t *testing.T) {
	ft := testFT()
	defer ft.db.Close()

	bf := ft.newFT()
	bf.columns = []string{"c0", "c1", "c2", "c3", "c4", "c5"}
	bf.colIndex = map[string]int{"c0": 0, "c1": 1, "c2": 2, "c3": 3, "c4": 4, "c5": 5}
	bf.keys = [][]string{{"c0"}, {"c1", "c2"}}
	bf.cardinality = map[string]int{"c0": 10, "c1": 3, "c2": 50, "c3": 200, "c4": 5, "c5": 8}
	bf.makeData()

	// single-column key c0 values should be unique
	seen := map[string]int{}
	for i, row := range bf.data {
		if len(row) != len(bf.columns) {
			t.Errorf("Row %d has %d columns, expected %d", i, len(row), len(bf.columns))
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
