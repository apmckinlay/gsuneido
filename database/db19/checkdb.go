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

const errLimit = 20

// CheckDatabase checks the integrity of the database.
// It returns "" if no errors are found, otherwise an error message.
func CheckDatabase(dbfile string) string {
	db := OpenDatabaseRead(dbfile)
	defer db.Close()
	state := db.GetState()
	dc := dbcheck{store: db.store, state: state}
	dc.checkTables()
	return strings.Join(dc.errs, "\n")
}

type dbcheck struct {
	store *stor.Stor
	state *DbState
	errs  []string
}

func (dc *dbcheck) checkTables() {
	defer func() {
		if e := recover(); e != nil {
			dc.errs = append(dc.errs, fmt.Sprint(e))
		}
	}()
	dc.state.meta.ForEachSchema(func(sc *meta.Schema) {
		dc.checkTable(sc)
		if len(dc.errs) > errLimit {
			panic("too many errors")
		}
	})
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
