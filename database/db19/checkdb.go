// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
)

const errLimit = 20

// CheckDatabase checks the integrity of the database.
// It returns "" if no errors are found, otherwise an error message.
func CheckDatabase(dbfile string) string {
	db := OpenDatabaseRead(dbfile)
	defer db.Close()
	state := db.GetState()
	dc := dbcheck{state: state}
	dc.checkTables()
	return strings.Join(dc.errs, "\n")
}

type dbcheck struct {
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
	for i, ix := range info.Indexes {
		n := dc.checkIndex(sc, i, ix)
		if n != info.Nrows {
			dc.errs = append(dc.errs,
				fmt.Sprint(sc.Table, " ", sc.Indexes[i].String(sc.Columns),
					" expected count of ", info.Nrows,
					" but got ", n))
		}
	}
}

func (dc *dbcheck) checkIndex(sc *meta.Schema, i int, ix *btree.Overlay) int {
	defer func() {
		if e := recover(); e != nil {
			dc.errs = append(dc.errs,
				fmt.Sprint(sc.Table, " ", sc.Indexes[i].String(sc.Columns), e))
		}
	}()
	return ix.Check()
}
