// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"fmt"
	"time"

	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

type dbcheck DbState

// quick check ------------------------------------------------------

// QuickCheck is the default partial checking done at start up.
// Panics on error.
func (db *Database) QuickCheck() {
	t := time.Now()
	dc := (*dbcheck)(db.GetState())
	n := dc.forEachTable(dc.quickCheckTable)
	fmt.Println("quick checked", n, "tables in",
		time.Since(t).Round(time.Millisecond))
}

func (dc dbcheck) quickCheckTable(sc *meta.Schema) {
	info := dc.meta.GetRoInfo(sc.Table)
	for _, ix := range info.Indexes {
		ix.QuickCheck()
	}
}

// full check -------------------------------------------------------

// CheckDatabase checks the integrity of the database.
func CheckDatabase(dbfile string) (ec *ErrCorrupt) {
	defer func() {
		if e := recover(); e != nil {
			ec = NewErrCorrupt(e)
		}
	}()
	db, err := openDatabase(dbfile, stor.READ, false)
	if err != nil {
		return NewErrCorrupt(err)
	}
	defer db.Close()
	dc := (*dbcheck)(db.GetState())
	dc.forEachTable(dc.checkTable)
	return nil
}

func (dc *dbcheck) forEachTable(fn func(sc *meta.Schema)) int {
	n := 0
	dc.meta.ForEachSchema(func(sc *meta.Schema) {
		n++
		fn(sc)
	})
	return n
}

func (dc *dbcheck) checkTable(sc *meta.Schema) {
	info := dc.meta.GetRoInfo(sc.Table)
	count, sum := dc.checkFirstIndex(sc.Table, info.Indexes[0])
	if count != info.Nrows {
		panic(&ErrCorrupt{table: sc.Table})
	}
	for i := 1; i < len(info.Indexes); i++ {
		ix := info.Indexes[i]
		count, sum = checkOtherIndex(sc.Table, ix, count, sum)
	}
}

func (dc *dbcheck) checkFirstIndex(table string, ix *btree.Overlay) (int, uint64) {
	sum := uint64(0)
	count := ix.Check(func(off uint64) {
		sum += off // addition so order doesn't matter
		buf := dc.store.Data(off)
		size := runtime.RecLen(buf)
		if !cksum.Check(buf[:size+cksum.Len]) {
			// fmt.Println("data checksum")
			panic(&ErrCorrupt{table: table})
		}
	})
	return count, sum
}

func checkOtherIndex(table string, ix *btree.Overlay,
	countPrev int, sumPrev uint64) (int, uint64) {
	sum := uint64(0)
	count := ix.Check(func(off uint64) {
		sum += off // addition so order doesn't matter
	})
	if count != countPrev || sum != sumPrev {
		panic(&ErrCorrupt{table: table})
	}
	return count, sum
}

//-------------------------------------------------------------------

type ErrCorrupt struct {
	err   error
	table string
}

func (ec *ErrCorrupt) Error() string {
	if ec.err == nil {
		return "database corrupt"
	}
	return "database corrupt: " + ec.err.Error()
}
func (ec *ErrCorrupt) Unwrap() error {
	return ec.err
}
func (ec *ErrCorrupt) Table() string {
	if ec == nil {
		return ""
	}
	return ec.table
}
func NewErrCorrupt(e interface{}) *ErrCorrupt {
	if e == nil {
		return nil
	}
	if e2, ok := e.(*ErrCorrupt); ok {
		return e2
	}
	if e2, ok := e.(error); ok {
		return &ErrCorrupt{err: e2}
	}
	return &ErrCorrupt{err: errors.New(fmt.Sprint(e))}
}
