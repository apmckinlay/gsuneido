// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"hash/crc32"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/blkbloom"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

// These benchmarks calibrate relative optimizer costs, not costs in ns.
// The database is memory mapped and persisted before any read transactions.
// Repeated selective scans can fit in CPU caches even when the table does not.
// recordRead and all plan benchmarks checksum every fetched record, including
// candidates rejected by residual predicates. This includes checksum CPU work,
// not just storage latency. recordView measures the zero-copy fetch instead.
//
// Run individual groups, or the whole suite, with:
//
//	go test -short -timeout=30s -run '^$' -bench 'Bitmap/Factors' -benchmem ./dbms/query
//	go test -short -timeout=30s -run '^$' -bench 'Bitmap/Plan' -benchmem ./dbms/query
//	go test -short -timeout=30s -run '^$' -bench 'Bitmap/BloomCost' -benchtime=500ms -count=3 -benchmem ./dbms/query
//	go test -short -timeout=30s -run '^$' -bench 'Bitmap' -benchtime=500ms -benchmem ./dbms/query
//
// B/op measures allocations, not retained size.

type benchTable struct {
	name    string
	rows    int
	width   int
	indexes string
}

var bloomBenchTables = []benchTable{
	{"bt", 1_000_000, 200, "index(a) index(b)"},
	{"bw", 1_000_000, 1000, "index(a) index(b)"},
}

var benchCases = []struct{ main, bloom float64 }{
	{.5, .5}, {.25, .25}, {.1, .1}, {.05, .05}, {.01, .01}, {.002, .002},
	{.5, .01}, {.01, .5},
}

var (
	benchCRC        = crc32.MakeTable(crc32.Castagnoli)
	benchRecordSink uint32
	benchBloomSink  blkbloom.Bloom
	benchHitsSink   int
)

// One parent owns the shared fixture until all its subbenchmarks finish.
func BenchmarkBitmap(b *testing.B) {
	db := openBenchDb(b, bloomBenchTables)
	b.Run("Factors", func(b *testing.B) { benchmarkBitmapFactors(b, db) })
	b.Run("Plan", func(b *testing.B) { benchmarkBitmapPlan(b, db) })
	b.Run("BloomCost", func(b *testing.B) { benchmarkBloomCost(b, db) })
}

// Calibrate Bloom operations against Table.costFor, using the same field access
// as BenchmarkTableGetCost, without adding checksum work to the scan baseline.
func benchmarkBloomCost(b *testing.B, db *db19.Database) {
	for _, bt := range bloomBenchTables {
		tran := db.NewReadTran(AllPerms)
		defer tran.Complete()
		for _, col := range []string{"k", "a"} {
			b.Run(fmt.Sprintf("%s/table_%s", bt.name, col), func(b *testing.B) {
				tbl := NewTable(tran, bt.name).(*Table)
				index := tbl.indexes[slices.IndexFunc(tbl.indexes,
					func(index []string) bool { return index[0] == col })]
				tbl.SetIndex(index, ReadMode)
				_, cost, _ := tbl.costFor(index, ReadMode, NoneReq(1))
				bytes, count := 0, 0
				for b.Loop() {
					tbl.Rewind()
					for row := tbl.Get(nil, Next); row != nil; row = tbl.Get(nil, Next) {
						bytes += len(row.GetRaw(tbl.Header(), "pad"))
						count++
					}
				}
				if count != bt.rows*b.N {
					b.Fatal("scan count changed")
				}
				benchHitsSink = bytes
				reportPer(b, bt.rows, "entry")
				b.ReportMetric(float64(cost)/float64(bt.rows), "units/entry")
			})
		}
		tbl := NewTable(tran, bt.name).(*Table)
		probes := collectOffsets(tbl, benchSelection(b, tbl, "a", bt.rows, 1))
		for _, frac := range []float64{1, .1, .01, .002} {
			offs := collectOffsets(tbl, benchSelection(b, tbl, "b", bt.rows, frac))
			b.Run(fmt.Sprintf("%s/frac%g/Add", bt.name, frac), func(b *testing.B) {
				for b.Loop() {
					bf := blkbloom.New(len(offs))
					for _, off := range offs {
						bf.Add(off)
					}
					benchBloomSink = bf
				}
				reportPer(b, len(offs), "offset")
			})
			b.Run(fmt.Sprintf("%s/frac%g/Has", bt.name, frac), func(b *testing.B) {
				bf := blkbloom.New(len(offs))
				for _, off := range offs {
					bf.Add(off)
				}
				hits := 0
				for b.Loop() {
					for _, off := range probes {
						if bf.Has(off) {
							hits++
						}
					}
				}
				benchHitsSink = hits
				reportPer(b, len(probes), "probe")
			})
		}
	}
}

// CompactDb is wired up by an external test (package query_test).
// The query package cannot import db19/tools directly because tools imports query.
var CompactDb func(dbfile string) error

func openBenchDb(t testing.TB, tables []benchTable) *db19.Database {
	t.Helper()
	// The database must share a filesystem with the compaction temp file,
	// which db19/tools creates in the working directory.
	path := t.Name() + ".db"
	t.Cleanup(func() { os.Remove(path); os.Remove(path + ".bak") })
	db, err := db19.CreateDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	oldDbTran, oldQueryTran := db19.MakeSuTran, MakeSuTran
	t.Cleanup(func() {
		db.Close()
		db19.MakeSuTran, MakeSuTran = oldDbTran, oldQueryTran
	})
	db19.MakeSuTran = func(*db19.UpdateTran) *SuTran { return NewSuTran(nil, true) }
	MakeSuTran = func(QueryTran) *SuTran { return nil }
	db19.StartConcur(db, time.Hour)
	fillBenchDb(db, tables)
	db.Persist() // drain writes and merge all index layers before taking snapshots
	// Rewrite the database so table data and index nodes are contiguous
	// and dead records are removed, as in a compacted production database.
	db.Close()
	if err := CompactDb(path); err != nil {
		t.Fatal(err)
	}
	db, err = db19.OpenDatabase(path)
	if err != nil {
		t.Fatal(err)
	}
	db19.StartConcur(db, time.Hour)
	return db
}

func fillBenchDb(db *db19.Database, tables []benchTable) {
	for _, bt := range tables {
		doAdmin(db, "create "+bt.name+" (k, a, b, pad) key(k) "+bt.indexes)
		pad := SuStr(strings.Repeat("x", bt.width))
		r := rand.New(rand.NewPCG(1, 2))
		const batch = 5000 // writeMax is 10000 per transaction
		for start := 0; start < bt.rows; start += batch {
			ut := db.NewUpdateTran(AllPerms)
			for k := start; k < start+batch && k < bt.rows; k++ {
				var rb RecordBuilder
				rb.Add(IntVal(k))
				rb.Add(IntVal(r.IntN(bt.rows)))
				rb.Add(IntVal(r.IntN(bt.rows)))
				rb.Add(pad)
				ut.Output(nil, bt.name, rb.Build())
			}
			ut.Commit()
		}
	}
}

// Intercept fetching, not output, so residual rejects also touch the payload.
type benchReadTran struct {
	QueryTran
	checksum uint32
	records  int64
	bytes    int64
}

func (t *benchReadTran) GetRecord(off uint64) Record {
	rec := t.QueryTran.GetRecord(off)
	t.consume(rec)
	return rec
}

func (t *benchReadTran) Lookup(table string, index int, key string) Record {
	rec := t.QueryTran.Lookup(table, index, key)
	if rec != "" {
		t.consume(rec)
	}
	return rec
}

func (t *benchReadTran) consume(rec Record) {
	t.checksum = crc32.Update(t.checksum, benchCRC, hacks.Stobs(string(rec)))
	t.records++
	t.bytes += int64(len(rec))
}

func (t *benchReadTran) reset() {
	t.checksum, t.records, t.bytes = 0, 0, 0
}

func (t *benchReadTran) report(b *testing.B) {
	benchRecordSink = t.checksum
	b.ReportMetric(float64(t.records)/float64(b.N), "reads/op")
	b.ReportMetric(float64(t.bytes)/float64(b.N), "record-B/op")
	if t.bytes > 0 {
		b.SetBytes(t.bytes / int64(b.N))
	}
}

func findIdxSel(w *Where, col string) *idxSel {
	for _, is := range w.idxSels {
		if len(is.index) > 0 && is.index[0] == col {
			return is
		}
	}
	return nil
}

func reportPer(b *testing.B, n int, unit string) {
	b.Helper()
	if n > 0 {
		ns := float64(b.Elapsed().Nanoseconds()) / float64(b.N)
		b.ReportMetric(ns/float64(n), "ns/"+unit)
	}
}

func benchThreshold(rows int, frac float64) int {
	return int(float64(rows)*(1-frac)) - 1
}

func benchSelection(t testing.TB, tbl *Table, col string, rows int, frac float64) *idxSel {
	t.Helper()
	q := fmt.Sprintf("%s where %s > %d", tbl.name, col, benchThreshold(rows, frac))
	w := ParseQuery(q, tbl.tran, nil).(*Where)
	w.optInit()
	sel := findIdxSel(w, col)
	if sel == nil || sel.indexFilterFrac != 1 || sel.skipLen != 0 {
		t.Fatalf("expected a simple index range for %q", q)
	}
	for _, pr := range sel.prefixRanges {
		if !pr.isRange() {
			t.Fatalf("expected a range for %q", q)
		}
	}
	return sel
}

func scanBenchOffsets(tbl *Table, sel *idxSel, visit func(string, uint64) bool) {
	for _, pr := range sel.prefixRanges {
		tbl.SelectRaw(pr.Org, pr.End)
		tbl.GetFilter(Next, visit)
	}
}

func collectOffsets(tbl *Table, sel *idxSel) []uint64 {
	tbl.iter = nil // SetIndex does not replace an existing iterator
	tbl.SetIndex(sel.index, ReadMode)
	var offs []uint64
	scanBenchOffsets(tbl, sel, func(_ string, off uint64) bool {
		offs = append(offs, off)
		return false
	})
	return offs
}

// Each factor uses the same range, Cur/key path, callback, and GetFilter count.
// Differences include GetRecord and (for recordRead) full-record checksumming.
func benchmarkBitmapFactors(b *testing.B, db *db19.Database) {
	for _, bt := range bloomBenchTables {
		tran := db.NewReadTran(AllPerms)
		tbl := NewTable(tran, bt.name).(*Table)
		for _, frac := range []float64{1, .5, .1, .01, .002} {
			sel := benchSelection(b, tbl, "a", bt.rows, frac)
			n := len(collectOffsets(tbl, sel))
			for _, mode := range []string{"offsetScan", "recordView", "recordRead"} {
				b.Run(fmt.Sprintf("%s/frac%g/%s", bt.name, frac, mode), func(b *testing.B) {
					full := &benchReadTran{QueryTran: tran}
					var bytes int64
					count := 0
					visit := func(_ string, off uint64) bool {
						count++
						if mode == "recordView" {
							bytes += int64(len(tran.GetRecord(off)))
						} else if mode == "recordRead" {
							full.GetRecord(off)
						}
						return false
					}
					for b.Loop() {
						scanBenchOffsets(tbl, sel, visit)
					}
					if count != n*b.N {
						b.Fatal("scan count changed")
					}
					benchHitsSink = int(bytes)
					reportPer(b, n, "entry")
					if mode == "recordRead" {
						full.report(b)
					}
				})
			}
		}
	}
}

func setupPlan(t testing.TB, query string, tran QueryTran, bloom bool) *Where {
	t.Helper()
	w := ParseQuery(query, tran, nil).(*Where)
	w.optInit()
	main, second := findIdxSel(w, "a"), findIdxSel(w, "b")
	if main == nil || second == nil {
		t.Fatalf("missing index for %q", query)
	}
	ap := &whereApproach{index: main.index, idxSel: main, mode: ReadMode}
	if bloom {
		ap.bloomSel = second
	}
	w.setApproach(NoneReq(1), ap, tran)
	return w
}

type benchResult struct {
	rows int
	sum  uint64
}

func drainBenchPlan(w *Where, th *Thread, limit int) benchResult {
	w.Rewind()
	var result benchResult
	for row := w.Get(th, Next); row != nil; row = w.Get(th, Next) {
		result.rows++
		result.sum += db19.RecToOff(w.t.GetStore(), row[0])
		if limit > 0 && result.rows == limit {
			break
		}
	}
	return result
}

func benchPlan(b *testing.B, w *Where, tran *benchReadTran, reuse bool, limit int, expected benchResult) {
	b.Helper()
	th := &Thread{}
	if got := drainBenchPlan(w, th, limit); got != expected {
		b.Fatalf("plan result %v, expected %v", got, expected)
	}
	tran.reset()
	rows := 0
	for b.Loop() {
		if !reuse {
			w.bloom = nil
		}
		rows += drainBenchPlan(w, th, limit).rows
	}
	if rows != expected.rows*b.N {
		b.Fatal("plan row count changed")
	}
	b.ReportMetric(float64(rows)/float64(b.N), "rows/op")
	tran.report(b)
}

// Plan compares production Where execution, reuse, optimizer choice, and first
// 10 rows.
func benchmarkBitmapPlan(b *testing.B, db *db19.Database) {
	for _, bt := range bloomBenchTables {
		for _, tc := range benchCases {
			q := fmt.Sprintf("%s where a > %d and b > %d", bt.name,
				benchThreshold(bt.rows, tc.main), benchThreshold(bt.rows, tc.bloom))
			b.Run(fmt.Sprintf("%s/a%g_b%g", bt.name, tc.main, tc.bloom), func(b *testing.B) {
				tran := &benchReadTran{QueryTran: db.NewReadTran(AllPerms)}
				baseline := setupPlan(b, q, tran, false)
				th := &Thread{}
				expected := drainBenchPlan(baseline, th, 0)
				first := drainBenchPlan(baseline, th, 10)
				for _, mode := range []string{"index", "bloom", "bloomReuse", "optimizer", "first10/index", "first10/bloom"} {
					b.Run(mode, func(b *testing.B) {
						bloom := mode == "bloom" || mode == "bloomReuse" || mode == "first10/bloom"
						w := setupPlan(b, q, tran, bloom)
						var fixed, variable Cost
						if mode == "optimizer" {
							query, f, v := Setup(ParseQuery(q, tran, nil), ReadMode, tran)
							w, fixed, variable = query.(*Where), f, v
							b.Logf("index %v, bloom %v", w.tbl.index, w.bloomSel != nil)
						} else {
							main, second := findIdxSel(w, "a"), findIdxSel(w, "b")
							_, c1, _ := w.tbl.costFor(second.index, ReadMode, NoneReq(1))
							_, c2, _ := w.tbl.costFor(main.index, ReadMode, NoneReq(1))
							if bloom {
								fixed, variable = WhereCost2(bt.rows, float64(c1), float64(c2), 1,
									second.indexRangeFrac, second.indexFilterFrac, main.indexRangeFrac,
									main.indexFilterFrac, main.hasDataFilter && second.hasDataFilter)
							} else {
								variable = WhereCost(float64(c2), 1, main.indexRangeFrac, main.indexFilterFrac, main.hasDataFilter)
							}
						}
						limit, result := 0, expected
						if strings.HasPrefix(mode, "first10/") {
							limit, result = 10, first
						}
						benchPlan(b, w, tran, mode == "bloomReuse", limit, result)
						if limit == 0 {
							if mode == "bloomReuse" {
								fixed = 0
							}
							b.ReportMetric(float64(fixed+variable), "cost-units")
						}
						estimated, _ := w.Nrows()
						b.ReportMetric(float64(estimated), "estimated-rows")
					})
				}
			})
		}
	}
}

func TestBitmapBench(t *testing.T) {
	a := assert.T(t)
	db := openBenchDb(t, []benchTable{{"small", 257, 0, "index(a) index(b)"}})
	tran := &benchReadTran{QueryTran: db.NewReadTran(AllPerms)}
	q := "small where a > 20 and b > 120"
	w := setupPlan(t, q, tran, false)
	th := &Thread{}
	expected := drainBenchPlan(w, th, 0)
	a.This(expected.rows > 0).Is(true)
	a.This(tran.records > int64(expected.rows)).Is(true) // includes residual rejects
	a.This(tran.bytes > tran.records).Is(true)
	a.This(drainBenchPlan(setupPlan(t, q, tran, true), th, 0)).Is(expected)
	tran.reset()
	tran.consume(Record("complete record"))
	a.This(tran.checksum).Is(crc32.Checksum([]byte("complete record"), benchCRC))
	a.This(tran.records).Is(int64(1))
	a.This(tran.bytes).Is(int64(len("complete record")))
}
