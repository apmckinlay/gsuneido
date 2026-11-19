// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stats"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestRegenStats(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, time.Hour)
	defer db.Close()
	prevMakeSuTran := db19.MakeSuTran
	t.Cleanup(func() { db19.MakeSuTran = prevMakeSuTran })
	db19.MakeSuTran = func(*db19.UpdateTran) *core.SuTran {
		return core.NewSuTran(nil, true)
	}
	query.DoAdminTest(db, "create tbl (a, b, c) key(a)")
	query.DoAdminTest(db, "create gone (a) key(a)")
	insert := func(from, to int) {
		for i := from; i < to; i++ {
			ut := db.NewUpdateTran(core.AllPerms)
			query.DoAction(nil, ut, "insert { a: "+core.IntVal(i).String()+
				", b: 'x', c: 1 } into tbl")
			assert.This(ut.Complete()).Is("")
		}
	}
	insert(0, 10)

	assert.This(RegenStats(db).Error()).Is("no existing dbstats")

	// initial stats on columns a and b only
	tally := make(stats.StatsTally)
	tally["tbl"] = stats.StatsTable{
		"a": &stats.StatsColumn{}, "b": &stats.StatsColumn{}}
	tally["tbl"]["a"].Add("x")
	tally["tbl"]["b"].Add("x")
	tally["gone"] = stats.StatsTable{"a": &stats.StatsColumn{}}
	tally["gone"]["a"].Add("x")
	createStatsTable(db, tally)
	db.SetStats(db19.ReadStats(db))
	assert.This(db.Stats()["tbl"].Count).Is(1)
	assert.This(db.Stats()["gone"].Count).Is(1)

	query.DoAdminTest(db, "drop gone")
	insert(10, 20)
	assert.This(RegenStats(db)).Is(nil)

	assert.This(len(db.Stats())).Is(1)
	ts := db.Stats()["tbl"]
	assert.This(ts.Count).Is(20)
	assert.This(len(ts.Columns)).Is(2) // same columns, c not added
	_, ok := ts.Columns["a"]
	assert.True(ok)
	_, ok = ts.Columns["b"]
	assert.True(ok)
	assert.This(ts.Columns["b"].Cardinality).Is(1)

	ut := db.NewUpdateTran(core.AllPerms)
	assert.This(query.DoAction(nil, ut, "delete tbl")).Is(20)
	assert.This(ut.Complete()).Is("")
	assert.This(RegenStats(db)).Is(nil)
	assert.This(len(db.Stats())).Is(0)
	assert.This(len(db19.ReadStats(db))).Is(0)
}

func TestStatsTableReadOnly(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	defer db.Close()
	st, ti := buildStatsTable(db, nil)
	db.AddNewTable(st, ti)
	table := db19.StatsTableName
	const msg = "can't modify system table: " + db19.StatsTableName

	assert.This(func() { query.DoAdminTest(db, "drop "+table) }).Panics(msg)
	_, err := LoadDbTable(table, "", "", "", db)
	assert.This(err.Error()).Is(msg)

	ut := db.NewUpdateTran(core.AllPerms)
	defer ut.Abort()
	assert.This(ut.Lookup(table, 0, "")).Isnt("")
	assert.This(query.NewTable(ut, table).Updateable()).Is("")
	assert.This(func() { ut.Output(nil, table, "") }).Panics(msg)
	assert.This(func() { ut.Update(nil, table, 0, "") }).Panics(msg)
	assert.This(func() { ut.Delete(nil, table, 0) }).Panics(msg)
}
