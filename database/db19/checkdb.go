// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"fmt"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

type dbcheck DbState

// quick check ------------------------------------------------------

// QuickCheck is the default partial checking done at start up.
// Panics on error.
func (db *Database) QuickCheck() {
	dc := (*dbcheck)(db.GetState())
	dc.forEachTable(dc.quickCheckTable)
}

func (dc dbcheck) quickCheckTable(sc *meta.Schema) {
	info := dc.meta.GetRoInfo(sc.Table)
	for _, ix := range info.Indexes {
		ix.QuickCheck(func(off uint64) {
			buf := dc.store.Data(off)
			size := rt.RecLen(buf)
			cksum.MustCheck(buf[:size+cksum.Len])
		})
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

func (dc *dbcheck) forEachTable(fn func(sc *meta.Schema)) {
	n := 0
	dc.meta.ForEachSchema(func(sc *meta.Schema) {
		n++
		fn(sc)
	})
	fmt.Println("processed", n, "tables")
}

func (dc *dbcheck) checkTable(sc *meta.Schema) {
	info := dc.meta.GetRoInfo(sc.Table)
	sumPrev := uint64(0)
	for i, ix := range info.Indexes {
		count, sum := dc.checkIndex(sc, i, ix)
		if count != info.Nrows || (i > 0 && sum != sumPrev) {
			// fmt.Println("i", i, "nrows", info.Nrows, "count", count, "sumPrev", sumPrev, "sum", sum)
			panic(&ErrCorrupt{table: sc.Table})
		}
		sumPrev = sum
	}
}

func (dc *dbcheck) checkIndex(sc *meta.Schema, i int, ix *btree.Overlay) (int, uint64) {
	sum := uint64(0)
	var fn func(uint64)
	if i == 0 {
		fn = func(off uint64) {
			sum += off // addition so order doesn't matter
			buf := dc.store.Data(off)
			size := rt.RecLen(buf)
			if !cksum.Check(buf[:size+cksum.Len]) {
				// fmt.Println("data checksum")
				panic(&ErrCorrupt{table: sc.Table})
			}
		}
	} else {
		fn = func(off uint64) {
			sum += off // addition so order doesn't matter
		}
	}
	n := ix.Check(fn)
	return n, sum
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
