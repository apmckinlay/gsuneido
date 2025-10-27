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

	"github.com/apmckinlay/gsuneido/db19"
)

func TestCosting_create(t *testing.T) {
	if testing.Short() {
		return
	}
	db := createDb()
	defer db.Close()
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
	for nrecs := 10; nrecs <= 1_000_000; nrecs *= 10 {
		for rsize := 200; rsize <= 200; rsize *= 10 {
			tbl := tblname(nrecs, rsize)
			tran := db.NewReadTran()
			q := ParseQuery(tbl+" sort "+tbl+"_b", tran, nil)
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
					assert.This(n).Is(nrecs)
				}
			})
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
	for nrecs := 10; nrecs <= 1_000_000; nrecs *= 10 {
		for rsize := 200; rsize <= 200; rsize *= 10 {
			tbl := tblname(nrecs, rsize)
			cols := []string{tbl + "_a"}
			tran := db.NewReadTran()
			q := ParseQuery(tbl, tran, nil)
			q, _, _ = Setup(q, ReadMode, tran)
			b.Run(tbl, func(b *testing.B) {
				successful := 0
				for b.Loop() {
					r := rand.Intn(nrecs * 2)
					vals := []string{Pack(IntVal(r))}
					if q.Lookup(th, cols, vals) != nil {
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

func BenchmarkCosting(b *testing.B) {
	db, err := db19.OpenDatabase("costing.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	th := &Thread{}
	for nrecs := 10; nrecs <= 1_000_000; nrecs *= 10 {
		for rsize := 200; rsize <= 200; rsize *= 10 {
			tbl := tblname(nrecs, rsize)
			tran := db.NewReadTran()
			q := ParseQuery(tbl, tran, nil)
			q, _, _ = Setup(q, ReadMode, tran)
			b.Run(tbl, func(b *testing.B) {
				for b.Loop() {
					q.Rewind()
					n := 0
					for {
						row := q.Get(th, Next)
						if row == nil {
							break
						}
						n++
					}
					assert.This(n).Is(nrecs)
				}
			})
		}
	}
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
	for nrecs := 10; nrecs <= 1_000_000; nrecs *= 10 {
		for rsize := 200; rsize <= 200; rsize *= 10 {
			fmt.Println(nrecs, rsize)
			doAdmin(db, "create "+mkschema(nrecs, rsize))
			createData(db, nrecs, rsize)
		}
	}
}

func tblname(nrecs, rsize int) string {
	return "tbl_" + strconv.Itoa(nrecs) + "_" + strconv.Itoa(rsize)
}

func mkschema(nrecs, rsize int) string {
	name := tblname(nrecs, rsize)
	s := " (_a, _b, _c, _d, _e, _f, _g, _h, _i, _j) key(_a)"
	s = name + strings.ReplaceAll(s, "_", name+"_")
	return s
}

// func fldname(nrecs, rsize, i int) string {
// 	fields := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
// 	return tblname(nrecs, rsize) + "_" + fields[i]
// }

func createData(db *db19.Database, nrecs, rsize int) {
	tbl := tblname(nrecs, rsize)
	for i := range nrecs {
		t := db.NewUpdateTran()
		var rb RecordBuilder
		rb.Add(IntVal(i * 2))
		rb.Add(IntVal(rand.Intn(nrecs)))
		for range 10 {
			rb.Add(SuStr(strings.Repeat("x", rsize/10)))
		}
		t.Output(nil, tbl, rb.Build())
		t.Commit()
	}
}
