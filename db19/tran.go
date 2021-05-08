// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"strconv"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
)

type tran struct {
	db   *Database
	meta *meta.Meta
}

func (t *tran) GetInfo(table string) *meta.Info {
	return t.meta.GetRoInfo(table)
}

func (t *tran) GetSchema(table string) *schema.Schema {
	return &t.getSchema(table).Schema
}

func (t *tran) getSchema(table string) *meta.Schema {
	if ts := t.meta.GetRoSchema(table); ts != nil {
		return ts
	}
	panic("table not found: " + table)
}

func (t *tran) GetAllInfo() []*meta.Info {
	infos := make([]*meta.Info, 0, 32)
	t.meta.ForEachInfo(func(info *meta.Info) { infos = append(infos, info) })
	return infos
}

func (t *tran) GetAllSchema() []*meta.Schema {
	schemas := make([]*meta.Schema, 0, 32)
	t.meta.ForEachSchema(
		func(schema *meta.Schema) { schemas = append(schemas, schema) })
	return schemas
}

//-------------------------------------------------------------------

type ReadTran struct {
	tran
	num int
}

var nextReadTran int32

func (db *Database) NewReadTran() *ReadTran {
	state := db.GetState()
	return &ReadTran{tran: tran{db: db, meta: state.Meta},
		num: int(atomic.AddInt32(&nextReadTran, 1))}
}

func (t *ReadTran) String() string {
	return "rt" + strconv.Itoa(t.num)
}

func (t *ReadTran) GetIndex(table string, cols []string) *index.Overlay {
	ts := t.meta.GetRoSchema(table)
	ti := t.meta.GetRoInfo(table)
	for i, ix := range ts.Indexes {
		if str.List(cols).Equal(ix.Columns) {
			return ti.Indexes[i]
		}
	}
	return nil
}

func (t *ReadTran) GetRecord(off uint64) rt.Record {
	buf := t.db.Store.Data(off)
	size := rt.RecLen(buf)
	return rt.Record(hacks.BStoS(buf[:size]))
}

func (t *ReadTran) ColToFld(table, col string) int {
	ts := t.meta.GetRoSchema(table)
	return str.List(ts.Columns).Index(col)
}

func (t *ReadTran) RangeFrac(table string, iIndex int, org, end string) float64 {
	idx := t.meta.GetRoInfo(table).Indexes[iIndex]
	return float64(idx.RangeFrac(org, end))
}

// Lookup returns the DbRec for a key, or nil if not found
func (t *ReadTran) Lookup(table string, iIndex int, key string) *rt.DbRec {
	idx := t.meta.GetRoInfo(table).Indexes[iIndex]
	off := idx.Lookup(key)
	if off == 0 {
		return nil
	}
	return &rt.DbRec{Off: off, Record: t.GetRecord(off)}
}

func (t *ReadTran) Output(string, rt.Record) {
	panic("can't output to read-only transaction")
}

func (t *ReadTran) Erase(uint64) {
	panic("can't delete from read-only transaction")
}

func (t *ReadTran) Update(uint64, rt.Record) uint64 {
	panic("can't update from read-only transaction")
}

func (t *ReadTran) ReadCount() int {
	return 0
}

func (t *ReadTran) WriteCount() int {
	return 0
}

func (t *ReadTran) MakeLess(is *ixkey.Spec) func(x, y uint64) bool {
	return t.db.MakeLess(is)
}

func (t *ReadTran) Complete() string {
	return ""
}

func (t *ReadTran) Abort() {
}

//-------------------------------------------------------------------

type UpdateTran struct {
	ReadTran
	ct *CkTran
}

func (db *Database) NewUpdateTran() *UpdateTran {
	state := db.GetState()
	meta := state.Meta.Mutable()
	ct := db.ck.StartTran()
	return &UpdateTran{ct: ct,
		ReadTran: ReadTran{tran: tran{db: db, meta: meta}}}
}

func (t *UpdateTran) String() string {
	return t.ct.String()
}

func (t *UpdateTran) Complete() string {
	if !t.db.ck.Commit(t) {
		return "commit failed" //TODO conflict description
	}
	return ""
}

func (t *UpdateTran) Commit() {
	// send commit to checker
	// which starts the pipeline to merger to persister
	t.ck(t.db.ck.Commit(t))
}

// commit is internal, called by checker (to serialize)
func (t *UpdateTran) commit() int {
	t.db.UpdateState(func(state *DbState) {
		state.Meta = t.meta.LayeredOnto(state.Meta)
	})
	return t.num()
}

func (t *UpdateTran) Abort() {
	t.ck(t.db.ck.Abort(t.ct))
}

func (t *UpdateTran) num() int {
	return t.ct.start
}

func (t *UpdateTran) Output(table string, rec rt.Record) {
	ts := t.getSchema(table)
	ti := t.getInfo(table)
	n := rec.Len()
	off, buf := t.db.Store.Alloc(n + cksum.Len)
	copy(buf, rec[:n])
	cksum.Update(buf)
	keys := make([]string, len(ts.Indexes))
	for i := range ts.Indexes {
		is := ts.Indexes[i].Ixspec
		keys[i] = is.Key(rec)
		ti.Indexes[i].Insert(keys[i], off)
	}
	t.ck(t.db.ck.Write(t.ct, table, keys))
	ti.Nrows++
	ti.Size += uint64(n)
}

func (t *UpdateTran) getInfo(table string) *meta.Info {
	if ti := t.meta.GetRwInfo(table); ti != nil {
		return ti
	}
	panic("table not found: " + table)
}

func (t *UpdateTran) ck(result bool) {
	if !result {
		conflict := t.ct.conflict.Load()
		if conflict == nil {
			panic("transaction already ended")
		}
		panic("transaction aborted: " + conflict.(string))
	}
}

func (t *UpdateTran) Erase(uint64) {
	panic("Erase not implemented") //TODO
}

func (t *UpdateTran) Update(uint64, rt.Record) uint64 {
	panic("Update not implemented") //TODO
}
