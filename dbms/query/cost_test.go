// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"hash/maphash"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"

	"github.com/apmckinlay/gsuneido/db19"
)

const dbName = "costing.db"
const tableSchema = "(a, b, c, d, e, f, g, h, i, j) key(a) index(b)"

var nrecs = []int{1000, 10_000, 100_000, 1_000_000}
var rsizes = []int{200, 1000}

// go test -v -run ^TestCosting_create$ ./dbms/query

func TestCost_create(t *testing.T) {
	assert.TestOnlyIndividually(t)
	db := createDb()
	defer db.Close()
}

func createDb() *db19.Database {
	db, err := db19.CreateDatabase(dbName)
	if err != nil {
		panic(err.Error())
	}
	db19.StartConcur(db, time.Hour)
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	createTables(db)
	db.Close()
	if err := CompactDb(dbName); err != nil {
		panic(err.Error())
	}
	db, err = db19.OpenDatabase(dbName)
	if err != nil {
		panic(err.Error())
	}
	db19.StartConcur(db, time.Hour)
	return db
}

func createTables(db *db19.Database) {
	for _, nrec := range nrecs {
		for _, rsize := range rsizes {
			fmt.Println(nrec, rsize)
			doAdmin(db, "create "+tblname(nrec, rsize)+" "+tableSchema)
			createData(db, nrec, rsize)
		}
	}
}

func tblname(nrecs, rsize int) string {
	return "tbl_" + strconv.Itoa(nrecs) + "_" + strconv.Itoa(rsize)
}

func createData(db *db19.Database, nrecs, rsize int) {
	fill := SuStr(strings.Repeat("x", rsize/8))
	for i := range nrecs {
		t := db.NewUpdateTran(AllPerms)
		var rb RecordBuilder
		rb.Add(IntVal(i))            // a: ordered
		rb.Add(IntVal(rand.Int64())) // b: random
		for range 8 {
			rb.Add(fill)
		}
		t.Output(nil, tblname(nrecs, rsize), rb.Build())
		t.Commit()
	}
}

//-------------------------------------------------------------------
// go test -run=NONE -bench ^BenchmarkTableGetCost$ ./dbms/query

func BenchmarkTableGetCost(b *testing.B) {
	db, err := db19.OpenDatabase("costing.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var nrec, rsize int
	th := &Thread{}
	var q Query
	seed := maphash.MakeSeed()
	fn := func(b *testing.B) {
		var total time.Duration
		for b.Loop() {
			if nrec*rsize < cacheSize {
				// Queries normally run once on a cold cache.
				// Only needed on small tables that could be cached.
				evictCache()
			}
			q.Rewind()
			t := time.Now()
			for row := q.Get(th, Next); !row.Eof(); row = q.Get(th, Next) {
				maphash.String(seed, string(row.Get(0))) // access all data
			}
			total += time.Since(t)
		}
		b.ReportMetric(0, "ns/op")
		b.ReportMetric(
			float64(1000*int64(b.N*nrec*rsize)/total.Nanoseconds()),
			"mb/s")
	}
	for _, nrec = range nrecs {
		for _, rsize = range rsizes {
			tbl := tblname(nrec, rsize)
			tran := db.NewReadTran(AllPerms)
			for _, idx := range []string{"a", "b"} {
				q = ParseQuery(tbl, tran, nil)
				q = setupIndex(q, ReadMode, tran, []string{idx})
				b.Run(tbl+"^"+idx, fn)
			}
		}
	}
}

const cacheSize = 64 * 1024 * 1024

var sink byte
var evictBuf = make([]byte, 2*cacheSize)

func evictCache() {
	for i := 0; i < len(evictBuf); i += 64 {
		evictBuf[i]++
	}
	sink += evictBuf[len(evictBuf)/2]
}

func BenchmarkTempindex(b *testing.B) {
	db, err := db19.OpenDatabase("costing.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	th := &Thread{}
	for _, nrec := range nrecs {
		for _, rsize := range rsizes {
			tbl := tblname(nrec, rsize)
			tran := db.NewReadTran(AllPerms)
			q := ParseQuery(tbl+" sort "+tbl+"_c", tran, nil)
			q, _, _ = Setup(q, ReadMode, tran)
			ti := q.(*Sort).source.(*TempIndex)
			b.Run(tbl, func(b *testing.B) {
				for b.Loop() {
					q.Rewind()
					ti.iter = nil // force rebuilding index
					ti.source.Rewind()
					n := 0
					for {
						row := q.Get(th, Next)
						if row.Eof() {
							break
						}
						n++
					}
					assert.This(n).Is(nrec)
				}
			})
		}
	}
}

func BenchmarkTempindexScan(b *testing.B) {
	db, err := db19.OpenDatabase("costing.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	th := &Thread{}
	for _, nrec := range nrecs {
		for _, rsize := range rsizes {
			tbl := tblname(nrec, rsize)
			tran := db.NewReadTran(AllPerms)
			q := ParseQuery(tbl+" sort "+tbl+"_c", tran, nil)
			q, _, _ = Setup(q, ReadMode, tran)
			ti := q.(*Sort).source.(*TempIndex)
			// build the temp index once before measuring
			q.Rewind()
			for !q.Get(th, Next).Eof() {
			}
			b.Run(tbl, func(b *testing.B) {
				for b.Loop() {
					q.Rewind()
					// ti.iter is kept; no rebuild
					_ = ti
					n := 0
					for {
						row := q.Get(th, Next)
						if row.Eof() {
							break
						}
						n++
					}
					assert.This(n).Is(nrec)
				}
			})
		}
	}
}

// BenchmarkTempindexCreate measures tempIndex build time for three sort orderings
// to validate the factorAll/Pre/None ratios in ticost.
//   - factorAll:  sort _a    — data arrives in _a order (key index), already sorted
//   - factorPre:  sort _b,_a — data arrives in _b order (index), sort within groups by _a
//   - factorNone: sort _c    — no index on _c, full sort from arbitrary order
func BenchmarkTempindexCreate(b *testing.B) {
	db, err := db19.OpenDatabase("costing.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	th := &Thread{}
	tran := db.NewReadTran(AllPerms)
	for _, nrec := range nrecs {
		for _, rsize := range rsizes {
			tbl := tblname(nrec, rsize)
			q := ParseQuery(tbl, tran, nil)
			q, _, _ = Setup(q, ReadMode, tran)
			fmt.Println(String(q))
			// factorAll: sort _a (pre-sorted input)
			{
				q.Rewind()
				ti := NewTempIndex(q, []string{"a"}, tran)
				b.Run(tbl+"/factorAll", func(b *testing.B) { ti.makeIndex() })
			}
			// factorPre: sort _b,_a (_b has index — data arrives in _b order,
			// only need to sort within each _b group by _a)
			{ //FIXME
				q.Rewind()
				ti := NewTempIndex(q, []string{tbl + "_b", "a"}, tran)
				b.Run(tbl+"/factorPre", func(b *testing.B) { ti.makeIndex() })
			}
			// factorNone: sort _c (no index — full sort from arbitrary order)
			{
				q.Rewind()
				ti := NewTempIndex(q, []string{tbl + "_b", "a"}, tran)
				b.Run(tbl+"/factorNone", func(b *testing.B) { ti.makeIndex() })
			}
			// measure extend overhead
			{
				tran := db.NewReadTran(AllPerms)
				q := ParseQuery(tbl+" extend x="+tbl+"_a", tran, nil)
				q, _, _ = Setup(q, ReadMode, tran)
				b.Run(tbl+"/extend", func(b *testing.B) {
					for b.Loop() {
						q.Rewind()
						for {
							row := q.Get(th, Next)
							if row.Eof() {
								break
							}
						}
					}
				})
			}
		}
	}
}

func BenchmarkLookup(b *testing.B) {
	db, err := db19.OpenDatabase("costing.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	th := &Thread{}
	for _, nrec := range nrecs {
		for _, rsize := range rsizes {
			tbl := tblname(nrec, rsize)
			sels := []Sel{{col: "a"}}
			tran := db.NewReadTran(AllPerms)
			q := ParseQuery(tbl, tran, nil)
			q, _, _ = SetupReq(q, ReadMode, tran, UniqueReq([]string{"a"}, 0))
			b.Run(tbl, func(b *testing.B) {
				successful := 0
				for b.Loop() {
					r := rand.IntN(nrec * 2)
					sels[0].val = Pack(IntVal(r))
					if !q.Lookup(th, sels).Eof() {
						successful++
					}
				}
				if b.N >= 10 {
					assert.That(successful*2 >= b.N-b.N/5)
					assert.That(successful*2 <= b.N+b.N/5)
				}
			})
		}
	}
}
