// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"errors"
	"fmt"
	"slices"

	"github.com/apmckinlay/gsuneido/core"
	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stats"
)

// RegenStats recomputes the stats for the same tables and columns
// that are currently in the stats table and replaces the stats table.
// It returns an error if there are no existing stats.
// Tables or columns that no longer exist are dropped.
func RegenStats(db *Database) (err error) {
	if db.IsCorrupted() {
		return errors.New("not allowed when database is locked")
	}
	cur := db.Stats()
	if len(cur) == 0 {
		return errors.New("no existing dbstats")
	}
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("regenerate stats failed: %v", e)
		}
	}()
	tally := make(stats.StatsTally)
	rt := db.NewReadTran(core.AllPerms)
	defer rt.Abort()
	for table, ts := range cur {
		regenTable(rt, tally, table, ts)
	}
	st, ti := buildStatsTable(db, tally)
	db.OverwriteTable(st, ti)
	// concurrent RegenStats could interleave here, but acceptable for this
	db.SetStats(ReadStats(db))
	return nil
}

func regenTable(rt *ReadTran, tally stats.StatsTally, table string,
	ts stats.TableStats) {
	if !rt.HasTable(table) {
		return
	}
	schema := rt.GetSchema(table)
	type colIdx struct {
		idx int
		sc  *stats.StatsColumn
	}
	st := make(stats.StatsTable)
	var cols []colIdx
	for col := range ts.Columns {
		if i := slices.Index(schema.Columns, col); i >= 0 && col != "-" {
			sc := &stats.StatsColumn{}
			st[col] = sc
			cols = append(cols, colIdx{idx: i, sc: sc})
		}
	}
	if len(cols) == 0 {
		return
	}
	iter := rt.IndexIter(table, 0)
	for iter.Next(rt); !iter.Eof(); iter.Next(rt) {
		rec := iter.GetRecord(rt)
		for _, c := range cols {
			c.sc.Add(rec.GetRaw(c.idx))
		}
	}
	tally[table] = st
}
