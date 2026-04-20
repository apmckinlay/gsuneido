// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"

	"github.com/apmckinlay/gsuneido/db19"
)

var nrecs = []int{1000, 10_000, 100_000, 1_000_000}
var rsizes = []int{200}

// go test -v -run ^TestCosting_create$ ./dbms/query

func TestCosting_create(t *testing.T) {
	assert.TestOnlyIndividually(t)
	db := createDb()
	defer db.Close()
}

func createDb() *db19.Database {
	db, err := db19.CreateDatabase("costing.db")
	if err != nil {
		panic(err.Error())
	}
	db19.StartConcur(db, 10*time.Second)
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	createTables(db)
	return db
}

func createTables(db *db19.Database) {
	for _, nrec := range nrecs {
		for _, rsize := range rsizes {
			fmt.Println(nrec, rsize)
			doAdmin(db, "create "+tblname(nrec, rsize)+
				" (a, b, c, d, e, f, g, h, i, j) key(a) key(b) index(c)")
			createData(db, nrec, rsize)
		}
	}
}

func tblname(nrecs, rsize int) string {
	return "tbl_" + strconv.Itoa(nrecs) + "_" + strconv.Itoa(rsize)
}

const cGroupSize = 16 // cluster size for c index

func createData(db *db19.Database, nrecs, rsize int) {
	tbl := tblname(nrecs, rsize)
	ngroups := max(1, nrecs/cGroupSize)
	var cgroup int
	for i := range nrecs {
		if i%cGroupSize == 0 {
			cgroup = rand.Intn(ngroups)
		}
		t := db.NewUpdateTran()
		var rb RecordBuilder
		rb.Add(IntVal(i))                              // a: in order
		rb.Add(IntVal(int(bits.Shuffle32(uint32(i))))) // b: random
		rb.Add(IntVal(cgroup))                         // c: clustered
		for range 7 {
			rb.Add(SuStr(strings.Repeat("x", rsize/8)))
		}
		t.Output(nil, tbl, rb.Build())
		t.Commit()
	}
}

//-------------------------------------------------------------------

func BenchmarkTableGetCost(b *testing.B) {
	db, err := db19.OpenDatabase("costing.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	th := &Thread{}
	var q Query
	fn := func(b *testing.B) {
		for b.Loop() {
			q.Rewind()
			hdr := q.Header()
			for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
				row.GetRaw(hdr, "j")
			}
		}
	}
	for _, nrec := range nrecs {
		for _, rsize := range rsizes {
			tbl := tblname(nrec, rsize)
			tran := db.NewReadTran()
			q = ParseQuery(tbl, tran, nil)
			q = setupIndex(q, ReadMode, []string{"a"}, tran) // physical order
			b.Run(tbl+"^a", fn)
			q = ParseQuery(tbl, tran, nil)
			q = setupIndex(q, ReadMode, []string{"b"}, tran) // random order
			b.Run(tbl+"^b", fn)
			q = ParseQuery(tbl, tran, nil)
			q = setupIndex(q, ReadMode, []string{"c"}, tran) // clustered
			b.Run(tbl+"^c", fn)
		}
	}
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
			tran := db.NewReadTran()
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
						if row == nil {
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
			tran := db.NewReadTran()
			q := ParseQuery(tbl+" sort "+tbl+"_c", tran, nil)
			q, _, _ = Setup(q, ReadMode, tran)
			ti := q.(*Sort).source.(*TempIndex)
			// build the temp index once before measuring
			q.Rewind()
			for q.Get(th, Next) != nil {
			}
			b.Run(tbl, func(b *testing.B) {
				for b.Loop() {
					q.Rewind()
					// ti.iter is kept; no rebuild
					_ = ti
					n := 0
					for {
						row := q.Get(th, Next)
						if row == nil {
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
	tran := db.NewReadTran()
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
				tran := db.NewReadTran()
				q := ParseQuery(tbl+" extend x="+tbl+"_a", tran, nil)
				q, _, _ = Setup(q, ReadMode, tran)
				b.Run(tbl+"/extend", func(b *testing.B) {
					for b.Loop() {
						q.Rewind()
						for {
							row := q.Get(th, Next)
							if row == nil {
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
			tran := db.NewReadTran()
			q := ParseQuery(tbl, tran, nil)
			q, _, _ = Setup(q, ReadMode, tran)
			b.Run(tbl, func(b *testing.B) {
				successful := 0
				for b.Loop() {
					r := rand.Intn(nrec * 2)
					sels[0].val = Pack(IntVal(r))
					if q.Lookup(th, sels) != nil {
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
