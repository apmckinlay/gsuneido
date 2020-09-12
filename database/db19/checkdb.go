// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

type dbcheck struct {
	store *stor.Stor
	state *DbState
	errs  []string
}

func (dc *dbcheck) addError(s string) {
	dc.errs = append(dc.errs, s)
	if len(dc.errs) > errLimit {
		panic("too many errors")
	}
}

// quick check ------------------------------------------------------

// QuickCheck is the default partial checking done at start up.
func (db *Database) QuickCheck() string {
	dc := dbcheck{store: db.store, state: db.GetState()}
	dc.forEachTable(dc.quickCheckTable)
	return strings.Join(dc.errs, "\n")
}

func (dc *dbcheck) quickCheckTable(sc *meta.Schema) {
	info := dc.state.meta.GetRoInfo(sc.Table)
	for _, ix := range info.Indexes {
		ix.QuickCheck(func(off uint64) {
			buf := dc.store.Data(off)
			size := rt.RecLen(buf)
			if !cksum.Check(buf[:size+cksum.Len]) {
			}
		})
		if len(dc.errs) > 0 {
			break
		}
	}
}

// full check -------------------------------------------------------

const errLimit = 20

// CheckDatabase checks the integrity of the database.
// It returns "" if no errors are found, otherwise an error message.
func CheckDatabase(dbfile string) string {
	db := openDatabase(dbfile, stor.READ, noCheck)
	defer db.Close()
	return db.Check()
}

func (db *Database) Check() string {
	dc := dbcheck{store: db.store, state: db.GetState()}
	dc.forEachTable(dc.checkTable)
	return strings.Join(dc.errs, "\n")
}

func (dc *dbcheck) forEachTable(fn func(sc *meta.Schema)) {
	defer func() {
		if e := recover(); e != nil {
			dc.errs = append(dc.errs, fmt.Sprint(e))
		}
	}()
	dc.state.meta.ForEachSchema(fn)
}

func (dc *dbcheck) checkTable(sc *meta.Schema) {
	info := dc.state.meta.GetRoInfo(sc.Table)
	sums := make([]uint64, len(info.Indexes))
	counts := make([]int, len(info.Indexes))
	for i, ix := range info.Indexes {
		counts[i], sums[i] = dc.checkIndex(sc, i, ix)
	}
	for i := 0; i < len(sums); i++ {
		if counts[i] != info.Nrows || sums[i] != sums[0] {
			err := info.Table + ": index mismatch, count " + strconv.Itoa(info.Nrows) + "\n"
			for i := 0; i < len(counts); i++ {
				err += fmt.Sprintln("   ", "count", counts[i], "sum", sums[i],
					sc.Indexes[i].String(sc.Columns))
			}
			dc.errs = append(dc.errs, strings.TrimRight(err, "\n"))
			break
		}
	}
}

func (dc *dbcheck) checkIndex(sc *meta.Schema, i int, ix *btree.Overlay) (int, uint64) {
	defer func() {
		if e := recover(); e != nil {
			dc.errs = append(dc.errs,
				fmt.Sprint(sc.Table, ": ", sc.Indexes[i].String(sc.Columns), " ", e))
		}
	}()
	sum := uint64(0)
	cksumerrs := 0
	var fn func(uint64) bool
	if i == 0 {
		fn = func(off uint64) bool {
			sum += off // addition so order doesn't matter
			buf := dc.store.Data(off)
			size := rt.RecLen(buf)
			if !cksum.Check(buf[:size+cksum.Len]) {
				cksumerrs++
				return false
			}
			return true
		}
	} else {
		fn = func(off uint64) bool {
			sum += off // addition so order doesn't matter
			return true
		}
	}
	n := ix.Check(fn)
	if cksumerrs > 0 {
		dc.errs = append(dc.errs,
			fmt.Sprint(sc.Table, ": data checksum errors (", cksumerrs, ")"))
	}
	return n, sum
}
