// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"

	"slices"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
)

type tran struct {
	db   *Database
	meta *meta.Meta
	asof int64
	off  uint64
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

func (t *UpdateTran) getRwInfo(table string) *meta.Info {
	if ti := t.meta.GetRwInfo(table); ti != nil {
		return ti
	}
	panic("nonexistent table: " + table)
}

func (t *tran) GetAllInfo() []*meta.Info {
	return slices.AppendSeq(make([]*meta.Info, 0, 32), t.meta.Infos())
}

func (t *tran) GetAllSchema() []*meta.Schema {
	return slices.AppendSeq(make([]*meta.Schema, 0, 32), t.meta.Tables())
}

func (t *tran) GetAllViews() []string {
	defs := make([]string, 0, 16)
	for name, def := range t.meta.Views() {
		defs = append(defs, name, def)
	}
	return defs
}

func (t *tran) GetView(name string) string {
	return t.db.GetView(name)
}

func (t *tran) GetStore() *stor.Stor {
	return t.db.Store
}

//-------------------------------------------------------------------

type ReadTran struct {
	tran
	num int
}

var nextReadTran atomic.Int32

func (db *Database) NewReadTran() *ReadTran {
	state := db.GetState()
	return &ReadTran{tran: tran{db: db, meta: state.Meta},
		num: int(nextReadTran.Add(2))} // even
}

func (t *ReadTran) String() string {
	return "rt" + strconv.Itoa(t.num)
}

func (t *ReadTran) Num() int {
	return t.num
}

func (t *ReadTran) GetIndex(table string, cols []string) *index.Overlay {
	ts := t.meta.GetRoSchema(table)
	if ts == nil {
		return nil
	}
	for i, ix := range ts.Indexes {
		if slices.Equal(cols, ix.Columns) {
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

func (t *ReadTran) GetRecord(off uint64) core.Record {
	buf := t.db.Store.Data(off)
	size := core.RecLen(buf)
	return core.Record(hacks.BStoS(buf[:size]))
}

func (t *ReadTran) ColToFld(table, col string) int {
	ts := t.meta.GetRoSchema(table)
	return slices.Index(ts.Columns, col)
}

func (t *ReadTran) RangeFrac(table string, iIndex int, org, end string) float64 {
	info := t.meta.GetRoInfo(table)
	idx := info.Indexes[iIndex]
	return idx.RangeFrac(org, end, info.Nrows)
}

// Lookup returns the DbRec for a key, or nil if not found
func (t *ReadTran) Lookup(table string, iIndex int, key string) *core.DbRec {
	idx := t.meta.GetRoInfo(table).Indexes[iIndex]
	off := idx.Lookup(key)
	if off == 0 {
		return nil
	}
	return &core.DbRec{Off: off, Record: t.GetRecord(off)}
}

func (t *ReadTran) Read(string, int, string, string) {
	// Read transactions don't need to track reads.
	// See UpdateTran Read.
}

func (t *ReadTran) Output(*core.Thread, string, core.Record) {
	panic("can't output to read-only transaction")
}

func (t *ReadTran) Delete(*core.Thread, string, uint64) {
	panic("can't delete from read-only transaction")
}

func (t *ReadTran) Update(*core.Thread, string, uint64, core.Record) uint64 {
	panic("can't update from read-only transaction")
}

func (t *ReadTran) ReadCount() int {
	return 0
}

func (t *ReadTran) WriteCount() int {
	return 0
}

func (t *ReadTran) Asof(asof int64) int64 {
	var state *DbState
	switch asof {
	case 0:
		return t.asof
	case -1:
		state = PrevState(t.db.Store, t.off)
	case 1:
		state = NextState(t.db.Store, t.off)
	default:
		state = StateAsof(t.db.Store, asof)
	}
	if state == nil {
		return 0
	}
	t.meta = state.Meta
	t.asof = state.Asof
	t.off = state.Off
	return t.asof
}

func (t *ReadTran) MakeLess(is *ixkey.Spec) func(x, y uint64) bool {
	return MakeLess(t.db.Store, is)
}

func (t *ReadTran) Complete() string {
	return ""
}

func (t *ReadTran) Abort() string {
	return ""
}

//-------------------------------------------------------------------

type UpdateTran struct {
	ct *CkTran
	ReadTran
	writeCount int
}

func (db *Database) NewUpdateTran() *UpdateTran {
	db.ckOpen()
	ct := db.ck.StartTran()
	if ct == nil {
		return nil
	}
	meta := ct.state.Meta.Mutable()
	return &UpdateTran{ct: ct,
		ReadTran: ReadTran{tran: tran{db: db, meta: meta}}}
}

func (t *UpdateTran) String() string {
	return t.ct.String()
}

func (t *UpdateTran) Num() int {
	return t.ct.start
}

func (t *UpdateTran) ReadCount() int {
	// read count is tracked by checker
	// so it knows about range consolidation
	return t.db.ck.ReadCount(t.ct)
}

func (t *UpdateTran) WriteCount() int {
	return t.writeCount
}

// Complete returns "" on success, otherwise an error
func (t *UpdateTran) Complete() string {
	if !t.db.ck.Commit(t) {
		return t.ct.failure.Load()
	}
	return ""
}

// Commit is used by tests. It panics on error.
func (t *UpdateTran) Commit() {
	t.ck(t.db.ck.Commit(t))
}

// commit is internal, called by checkco (to serialize)
func (t *UpdateTran) commit() int {
	t.db.UpdateState(func(state *DbState) {
		state.Meta = t.meta.LayeredOnto(state.Meta)
	})
	return t.num()
}

// Abort returns "" if it succeeds or if the transaction was already aborted.
func (t *UpdateTran) Abort() string {
	if !t.db.ck.Abort(t.ct, "aborted") {
		return "abort failed"
	}
	return ""
}

func (t *UpdateTran) num() int {
	return t.ct.start
}

// Lookup returns the DbRec for a key, or nil if not found
func (t *UpdateTran) Lookup(table string, iIndex int, key string) *core.DbRec {
	t.Read(table, iIndex, key, key)
	return t.ReadTran.Lookup(table, iIndex, key)
}

// Read adds a transaction read event to the checker
func (t *UpdateTran) Read(table string, iIndex int, from, to string) {
	t.ck(t.db.ck.Read(t.ct, table, iIndex, from, to))
}

const writeMax = 10000

func (t *UpdateTran) write() {
	if t.writeCount++; t.writeCount >= writeMax {
		t.Abort()
		panic("too many writes (output, update, or delete) in one transaction")
	}
}

func (t *UpdateTran) Output(th *core.Thread, table string, rec core.Record) {
	if t.db.corrupted.Load() {
		return // prevent appending to database
	}
	trace.Dbms.Println("tran Output", table)
	t.write()
	ts := t.getSchema(table)
	ti := t.tran.GetInfo(table) // readonly
	rec = rec.Truncate(len(ts.Columns))
	n := rec.Len()
	off, buf := t.db.Store.Alloc(n + cksum.Len)
	copy(buf, rec[:n])
	cksum.Update(buf)
	keys := make([]string, len(ts.Indexes))
	for i := range ts.Indexes {
		ix := ts.Indexes[i]
		keys[i] = ix.Ixspec.Key(rec)
		if ix.Mode == 'k' && len(ix.Columns) == 0 {
			if ti.Nrows > 0 {
				panic(fmt.Sprint("duplicate key: () in ", table))
			}
			t.Read(table, i, "", "")
		} else {
			t.dupOutputBlock(table, i, ix, ti.Indexes[i], rec, keys[i])
		}
		t.fkeyOutputBlock(ts, i, rec)
	}
	t.ck(t.db.ck.Output(t.ct, table, keys))
	func() {
		defer func() {
			if e := recover(); e != nil {
				t.Abort()
				panic(e)
			}
		}()
		ti = t.getRwInfo(table)
		for i := range ts.Indexes {
			ti.Indexes[i].Insert(keys[i], off)
		}
	}()
	ti.Nrows++
	ti.Size += int64(n)
	t.db.CallTrigger(th, t, table, "", rec)
}

func (t *UpdateTran) dupOutputBlock(table string, iIndex int, ix schema.Index,
	ov *index.Overlay, rec core.Record, key string) {
	if needsDupCheck(ix, rec) {
		if ov.Lookup(key) != 0 {
			panic(fmt.Sprint("duplicate key: ",
				str.Join(",", ix.Columns), " in ", table))
		}
		t.Read(table, iIndex, key, key)
	}
}

func needsDupCheck(ix schema.Index, rec core.Record) bool {
	if ix.Primary {
		return true
	}
	if ix.Mode == 'u' && !ix.ContainsKey && !uniqueIndexEmpty(rec, ix.Ixspec) {
		return true
	}
	return false
}

func uniqueIndexEmpty(rec core.Record, is ixkey.Spec) bool {
	for _, f := range is.Fields {
		if rec.GetRaw(f) != "" {
			return false
		}
	}
	return true
}

func (t *UpdateTran) fkeyOutputBlock(ts *meta.Schema, i int, rec core.Record) {
	ix := &ts.Indexes[i]
	fk := ix.Fk
	if fk.Table != "" {
		key := ix.Ixspec.Trunc(len(ix.Columns)).Key(rec)
		if key != "" && !t.fkeyOutputExists(fk.Table, fk.IIndex, key) {
			panic("output blocked by foreign key: " +
				fk.Table + " " + ix.String())
		}
	}
}

func (t *UpdateTran) fkeyOutputExists(table string, iIndex int, key string) bool {
	t.Read(table, iIndex, key, key)
	return t.ReadTran.fkeyOutputExists(table, iIndex, key)
}

// fkeyOutputExists on ReadTran is used by Database buildIndexes
func (t *ReadTran) fkeyOutputExists(table string, iIndex int, key string) bool {
	idx := t.meta.GetRoInfo(table).Indexes[iIndex]
	return idx.Lookup(key) != 0
}

func (t *UpdateTran) Delete(th *core.Thread, table string, off uint64) {
	trace.Dbms.Println("tran Delete", table, off)
	t.write()
	ts := t.getSchema(table)
	rec := t.GetRecord(off)
	n := int64(rec.Len())
	keys := make([]string, len(ts.Indexes))
	for i := range ts.Indexes {
		is := ts.Indexes[i].Ixspec
		keys[i] = is.Key(rec)
		t.fkeyDeleteBlock(ts, i, keys[i])
	}
	t.ck(t.db.ck.Delete(t.ct, table, off, keys))
	func() {
		defer func() {
			if e := recover(); e != nil {
				t.Abort()
				panic(e)
			}
		}()
		ti := t.getRwInfo(table)
		for i := range ts.Indexes {
			t.fkeyDeleteCascade(th, ts, i, keys[i])
		}
		for i := range ts.Indexes {
			prevoff := ti.Indexes[i].Delete(keys[i], off)
			if prevoff != 0 && prevoff != off {
				panic("update & delete on same record")
			}
	}
		assert.Msg("Delete Nrows").That(ti.Nrows > 0)
		ti.Nrows--
		assert.Msg("Delete Size").That(ti.Size >= n)
		ti.Size -= n
	}()
	t.db.CallTrigger(th, t, table, rec, "")
}

func (t *UpdateTran) fkeyDeleteBlock(ts *meta.Schema, i int, key string) {
	if key == "" {
		return
	}
	encoded := ts.Indexes[i].Ixspec.Encodes()
	fkToHere := ts.Indexes[i].FkToHere
	encKey := ""
	for j := range fkToHere {
		fk := &fkToHere[j]
		fkis := t.meta.GetRoSchema(fk.Table).Indexes[fk.IIndex].Ixspec
		fkey := key
		if !encoded && fkis.Encodes() {
			if encKey == "" {
				encKey = ixkey.Encode(key)
			}
			fkey = encKey
		}
		if fk.Mode == schema.Block && t.fkeyDeleteExists(fk, fkey) {
			panic("delete blocked by foreign key: " +
				fk.Table + " " + str.Join("(,)", fk.Columns))
		}
	}
}

func (t *UpdateTran) fkeyDeleteExists(fk *schema.Fkey, key string) bool {
	iter := index.NewOverIter(fk.Table, fk.IIndex)
	iter.Range(index.Range{Org: key, End: rangeEnd(key, len(fk.Columns))})
	iter.Next(t)
	return !iter.Eof()
}

// rangeEnd returns the end of the range for a key.
// Naively, you can just append a separator and Max.
// But that isn't correct if key doesn't have the correct (n) number of fields.
// WARNING: key must be encoded
func rangeEnd(key string, n int) string {
	kn := len(key)
	var sb strings.Builder
	sb.Grow(kn + len(ixkey.Sep) + len(ixkey.Max))
	for i := 0; i < kn; i++ {
		sb.WriteByte(key[i])
		if key[i] == '\x00' && i+1 < kn && key[i+1] == '\x00' {
			n--
			i++
			sb.WriteByte(key[i])
			if n == 0 {
				break
			}
		}
	}
	assert.That(n >= 0)
	for ; n > 0; n-- {
		sb.WriteByte(0)
		sb.WriteByte(0)
	}
	sb.WriteString(ixkey.Max)
	return sb.String()
}

func (t *UpdateTran) fkeyDeleteCascade(th *core.Thread, ts *meta.Schema, i int, key string) {
	if key == "" {
		return
	}
	encoded := ts.Indexes[i].Ixspec.Encodes()
	fkToHere := ts.Indexes[i].FkToHere
	for j := range fkToHere {
		fk := &fkToHere[j]
		if fk.Mode&schema.CascadeDeletes != 0 {
			iter := t.cascade(fk, encoded, key)
			for iter.Next(t); !iter.Eof(); iter.Next(t) {
				_, off := iter.Cur()
				t.Delete(th, fk.Table, off)
			}
		}
	}
}

func (t *UpdateTran) cascade(fk *schema.Fkey, encoded bool, key string) *index.OverIter {
	fkis := t.meta.GetRoSchema(fk.Table).Indexes[fk.IIndex].Ixspec
	if !encoded && fkis.Encodes() {
		key = ixkey.Encode(key)
	}
	iter := index.NewOverIter(fk.Table, fk.IIndex)
	iter.Range(index.Range{Org: key, End: rangeEnd(key, len(fk.Columns))})
	return iter
}

func (t *UpdateTran) Update(th *core.Thread, table string, oldoff uint64, newrec core.Record) uint64 {
	t.write()
	return t.update(th, table, oldoff, newrec, true)
}

func (t *UpdateTran) update(th *core.Thread, table string, oldoff uint64, newrec core.Record,
	block bool) uint64 {
	if t.db.corrupted.Load() {
		return oldoff // prevent appending to database
	}
	ts := t.getSchema(table)
	newrec = newrec.Truncate(len(ts.Columns))
	n := newrec.Len()
	newrec = newrec[:n]
	oldrec := t.GetRecord(oldoff)
	if newrec == oldrec {
		// in order to update we must already have read the record
		// so we should already have sent a read to the checker
		return oldoff
	}
	newoff, buf := t.db.Store.Alloc(n + cksum.Len)
	copy(buf, newrec)
	cksum.Update(buf)
	ti := t.tran.GetInfo(table) // read-only
	oldkeys := make([]string, len(ts.Indexes))
	newkeys := make([]string, len(ts.Indexes))
	for i := range ts.Indexes {
		ix := ts.Indexes[i]
		is := ix.Ixspec
		oldkeys[i] = is.Key(oldrec)
		newkeys[i] = is.Key(newrec)
		if oldkeys[i] != newkeys[i] {
			t.dupOutputBlock(table, i, ix, ti.Indexes[i], newrec, newkeys[i])
			t.fkeyDeleteBlock(ts, i, oldkeys[i])
			if block {
				t.fkeyOutputBlock(ts, i, newrec)
			}
		}
	}
	t.ck(t.db.ck.Update(t.ct, table, oldoff, oldkeys, newkeys))
	ti = t.getRwInfo(table)
	d := int64(len(newrec)) - int64(len(oldrec))
	assert.Msg("Update Size").That(int64(ti.Size)+d > 0)
	ti.Size = ti.Size + d
	func() {
		defer func() {
			if e := recover(); e != nil {
				t.Abort()
				panic(e)
			}
		}()
		for i := range ts.Indexes {
			if oldkeys[i] != newkeys[i] {
				t.fkeyUpdateCascade(th, ts, i, newrec, oldkeys[i])
			}
		}
		for i := range ts.Indexes {
			ix := ti.Indexes[i]
			if oldkeys[i] == newkeys[i] {
				prevoff := ix.Update(oldkeys[i], newoff)
				if prevoff != 0 && prevoff != oldoff {
					panic("update & update on same record")
				}
			} else {
				ix.Delete(oldkeys[i], oldoff)
				ix.Insert(newkeys[i], newoff)
			}
		}
	}()
	t.db.CallTrigger(th, t, table, oldrec, newrec)
	return newoff
}

func (t *UpdateTran) fkeyUpdateCascade(th *core.Thread, ts *meta.Schema, i int,
	rec core.Record, key string) { // rec is old, key is new
	ixcols := ts.Indexes[i].Columns
	encoded := ts.Indexes[i].Ixspec.Encodes()
	fkToHere := ts.Indexes[i].FkToHere
	for i := range fkToHere {
		fk := &fkToHere[i]
		if fk.Mode&schema.CascadeUpdates == 0 {
			continue
		}
		ts2 := t.GetSchema(fk.Table)
		ixcols2 := fk.Columns
		iter := t.cascade(fk, encoded, key)
		for iter.Next(t); !iter.Eof(); iter.Next(t) {
			_, off := iter.Cur()
			oldrec := t.GetRecord(off)
			rb := core.RecordBuilder{}
			for i, col := range ts2.Columns {
				if j := slices.Index(ixcols2, col); j != -1 {
					k := slices.Index(ts.Columns, ixcols[j])
					rb.AddRaw(rec.GetRaw(k))
				} else {
					rb.AddRaw(oldrec.GetRaw(i))
				}
			}
			newrec := rb.Trim().Build()
			t.update(th, fk.Table, off, newrec, false) // no output block
		}
	}
}

func (t *UpdateTran) ck(result bool) {
	if !result {
		failure := t.ct.failure.Load()
		if failure == "" {
			panic("transaction already ended")
		}
		panic("transaction aborted: " + failure)
	}
}
