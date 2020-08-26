// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"github.com/apmckinlay/gsuneido/database/db19/comp"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
)

// ck must be injected
var ck Checker

type tran struct {
	meta  *meta.Overlay
	store *stor.Stor
}

type ReadTran struct {
	tran
}

func NewReadTran() *ReadTran {
	state := GetState()
	return &ReadTran{tran: tran{meta: state.meta, store: state.store}}
}

type UpdateTran struct {
	tran
	ct *CkTran
}

func NewUpdateTran() *UpdateTran {
	state := GetState()
	meta := state.meta.NewOverlay()
	ct := ck.StartTran()
	return &UpdateTran{ct: ct, tran: tran{meta: meta, store: state.store}}
}

func (t *UpdateTran) Commit() {
	// send commit request to checker
	// which starts the pipeline to merger to persister
	t.ck(ck.Commit(t))
}

// commit is internal, called by checker (to serialize)
func (t *UpdateTran) commit() int {
	UpdateState(func(state *DbState) {
		state.meta = t.meta.LayeredOnto(state.meta)
	})
	return t.num()
}

func (t *UpdateTran) num() int {
	return t.ct.start
}

func (t *UpdateTran) Output(table string, rec rt.Record) {
	ts := t.getSchema(table)
	ti := t.getInfo(table)
	off, buf := t.store.AllocSized(len(rec))
	copy(buf, []byte(rec))
	keys := make([]string, len(ts.Indexes))
	for i := range ts.Indexes {
		is := ts.Indexes[i].Ixspec
		keys[i] = comp.Key(rec, is.Cols, is.Cols2)
		ti.Indexes[i].Insert(keys[i], off)
	}
	t.ck(ck.Write(t.ct, table, keys))
	ti.Nrows++
	ti.Size += uint64(len(rec))
}

func (t *UpdateTran) getInfo(table string) *meta.Info {
	if ti := t.meta.GetRwInfo(table, t.num()); ti != nil {
		return ti
	}
	panic("table not found: " + table)
}

func (t *UpdateTran) getSchema(table string) *meta.Schema {
	if ts := t.meta.GetRoSchema(table); ts != nil {
		return ts
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
