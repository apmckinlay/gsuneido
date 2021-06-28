// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"math"
	"strconv"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/strs"
)

type tran struct {
	db   *Database
	meta *meta.Meta
}

// GetInfo returns read-only Info for the table or nil if not found
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
	panic("nonexistent table: " + table)
}

func (t *UpdateTran) getInfo(table string) *meta.Info {
	if ti := t.meta.GetRwInfo(table); ti != nil {
		return ti
	}
	panic("nonexistent table: " + table)
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

func (t *tran) GetAllViews() []string {
	defs := make([]string, 0, 16)
	t.meta.ForEachView(func(name, def string) {
		defs = append(defs, name, def)
	})
	return defs
}

func (t *tran) GetView(name string) string {
	return t.db.GetView(name)
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
	if ts == nil {
		return nil
	}
	for i, ix := range ts.Indexes {
		if strs.Equal(cols, ix.Columns) {
			return t.GetIndexI(table, i)
		}
	}
	return nil
}

func (t *ReadTran) GetIndexI(table string, iIndex int) *index.Overlay {
	ti := t.meta.GetRoInfo(table)
	if ti == nil {
		return nil
	}
	return ti.Indexes[iIndex]
}

func (t *ReadTran) GetRecord(off uint64) rt.Record {
	buf := t.db.Store.Data(off)
	size := rt.RecLen(buf)
	return rt.Record(hacks.BStoS(buf[:size]))
}

func (t *ReadTran) ColToFld(table, col string) int {
	ts := t.meta.GetRoSchema(table)
	return strs.Index(ts.Columns, col)
}

func (t *ReadTran) RangeFrac(table string, iIndex int, org, end string) float64 {
	info := t.meta.GetRoInfo(table)
	idx := info.Indexes[iIndex]
	f := float64(idx.RangeFrac(org, end))
	if info.Nrows > 0 {
		f = math.Max(f, 1/float64(info.Nrows))
	}
	return f
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

func (t *ReadTran) Read(string, int, string, string) {
	// Read transactions don't need to track reads.
	// See UpdateTran Read.
}

func (t *ReadTran) Output(string, rt.Record) {
	panic("can't output to read-only transaction")
}

func (t *ReadTran) Delete(string, uint64) {
	panic("can't delete from read-only transaction")
}

func (t *ReadTran) Update(string, uint64, rt.Record) uint64 {
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

// Complete returns "" on success, otherwise an error
func (t *UpdateTran) Complete() string {
	if !t.db.ck.Commit(t) {
		conflict := t.ct.conflict.Load()
		if conflict == nil {
			return "transaction already ended"
		}
		return conflict.(string)
	}
	return ""
}

// Commit is used by tests. It panics on error.
func (t *UpdateTran) Commit() {
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
	t.ck(t.db.ck.Abort(t.ct, "aborted"))
}

func (t *UpdateTran) num() int {
	return t.ct.start
}

// Lookup returns the DbRec for a key, or nil if not found
func (t *UpdateTran) Lookup(table string, iIndex int, key string) *rt.DbRec {
	t.Read(table, iIndex, key, key)
	return t.ReadTran.Lookup(table, iIndex, key)
}

func (t *UpdateTran) Read(table string, iIndex int, from, to string) {
	t.ck(t.db.ck.Read(t.ct, table, iIndex, from, to))
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
		ix := ti.Indexes[i]
		is := ts.Indexes[i].Ixspec
		keys[i] = is.Key(rec)
		if ix.Lookup(keys[i]) != 0 {
			panic(fmt.Sprint("duplicate key: ",
				strs.Join(",", ts.Indexes[i].Columns), " in ", table))
		}
	}
	for i := range ts.Indexes {
		ix := ti.Indexes[i]
		ix.Insert(keys[i], off)
	}
	t.ck(t.db.ck.Write(t.ct, table, keys))
	ti.Nrows++
	ti.Size += uint64(n)
}

func (t *UpdateTran) Delete(table string, off uint64) {
	ts := t.getSchema(table)
	ti := t.getInfo(table)
	rec := t.GetRecord(off)
	n := rec.Len()
	keys := make([]string, len(ts.Indexes))
	for i := range ts.Indexes {
		ix := ti.Indexes[i]
		is := ts.Indexes[i].Ixspec
		keys[i] = is.Key(rec)
		ix.Delete(keys[i], off)
	}
	t.ck(t.db.ck.Write(t.ct, table, keys))
	assert.Msg("Delete Nrows").That(ti.Nrows > 0)
	ti.Nrows--
	assert.Msg("Delete Size").That(ti.Size >= uint64(n))
	ti.Size -= uint64(n)
}

func (t *UpdateTran) Update(table string, oldoff uint64, newrec rt.Record) uint64 {
	ts := t.getSchema(table)
	ti := t.getInfo(table)
	n := newrec.Len()
	newrec = newrec[:n]
	oldrec := t.GetRecord(oldoff)
	newoff := oldoff
	if newrec != oldrec {
		off, buf := t.db.Store.Alloc(n + cksum.Len)
		copy(buf, newrec)
		cksum.Update(buf)
		newoff = off
	}
	oldkeys := make([]string, len(ts.Indexes))
	newkeys := make([]string, len(ts.Indexes))
	for i := range ts.Indexes {
		is := ts.Indexes[i].Ixspec
		oldkeys[i] = is.Key(oldrec)
		if newoff != oldoff {
			newkeys[i] = is.Key(newrec)
			if oldkeys[i] != newkeys[i] {
				ix := ti.Indexes[i]
				if ix.Lookup(newkeys[i]) != 0 {
					panic(fmt.Sprint("duplicate key: ",
						strs.Join(",", ts.Indexes[i].Columns), " in ", table))
				}
			}
		}
	}
	if newoff != oldoff {
		for i := range ts.Indexes {
			ix := ti.Indexes[i]
			if oldkeys[i] == newkeys[i] {
				ix.Update(oldkeys[i], newoff)
			} else {
				ix.Delete(oldkeys[i], oldoff)
				ix.Insert(newkeys[i], newoff)
			}
		}
	}
	t.ck(t.db.ck.Write(t.ct, table, oldkeys))
	if newoff != oldoff {
		t.ck(t.db.ck.Write(t.ct, table, newkeys))
		d := int64(len(newrec) - len(oldrec))
		assert.Msg("Update Size").That(int64(ti.Size)+d > 0)
		ti.Size = uint64(int64(ti.Size) + d)
	}
	return newoff
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
