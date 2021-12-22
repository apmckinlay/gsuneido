// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/hash"
)

type Info struct {
	Table     string
	Nrows     int
	Size      uint64
	origNrows int
	origSize  uint64
	Indexes   []*index.Overlay
	// lastmod is used for persist meta chaining/flattening
	lastmod int
}

//go:generate genny -in ../../genny/hamt/hamt.go -out infohamt.go -pkg meta gen "Item=*Info KeyType=string"

func InfoKey(ti *Info) string {
	return ti.Table
}

func InfoHash(key string) uint32 {
	return hash.HashString(key)
}

func (ti *Info) storSize() int {
	size := 2 + len(ti.Table) + 4 + 5 + 1
	for i := range ti.Indexes {
		size += ti.Indexes[i].StorSize()
	}
	return size
}

func (ti *Info) Write(w *stor.Writer) {
	w.PutStr(ti.Table).
		Put4(ti.Nrows).
		Put5(ti.Size).
		Put1(len(ti.Indexes))
	for i := range ti.Indexes {
		ti.Indexes[i].Write(w)
	}
}

func ReadInfo(st *stor.Stor, r *stor.Reader) *Info {
	var ti Info
	ti.Table = r.GetStr()
	ti.Nrows = r.Get4()
	ti.Size = r.Get5()
	if ni := r.Get1(); ni > 0 {
		ti.Indexes = make([]*index.Overlay, ni)
		for i := 0; i < ni; i++ {
			ti.Indexes[i] = index.ReadOverlay(st, r)
		}
	}
	return &ti
}

func (m *Meta) newInfoTomb(table string) *Info {
	return &Info{Table: table}
}

func (ti *Info) isTomb() bool {
	return ti.Indexes == nil
}

func (ht InfoHamt) MustGet(key string) *Info {
	it, ok := ht.Get(key)
	if !ok || it.isTomb() {
		panic("info MustGet failed for " + key)
	}
	return it
}

// GetCopy returns a copy of the Info for a table, or nil if not found
func (ht InfoHamt) GetCopy(table string) *Info {
	ti, ok := ht.Get(table)
	if !ok || ti.isTomb() {
		return nil
	}
	cp := *ti
	return &cp
}

//-------------------------------------------------------------------

type btOver = *index.Overlay
type MergeResult = index.MergeResult

type MergeUpdate struct {
	table   string
	idTran  int
	nmerged int
	results []MergeResult // per index
}

// Merge collects the updates which are then applied by ApplyMerge.
// It is called by db Merge which is called by concur merger.
// WARNING: must not modify meta.
func (m *Meta) Merge(metaWas *Meta, table string, nmerge int) MergeUpdate {
	was := metaWas.schema.MustGet(table)
	cur, ok := m.schema.Get(table)
	if !ok || cur.isTomb() || cur.Id != was.Id {
		return MergeUpdate{} // table dropped or recreated
	}
	ti := m.info.MustGet(table)
	results := make([]MergeResult, len(ti.Indexes))
	for i, ov := range ti.Indexes {
		if !skipIndex(was, cur, i) {
			results[i] = ov.Merge(nmerge)
		}
	}
	return MergeUpdate{table: table, idTran: was.Id,
		nmerged: nmerge, results: results}
}

func skipIndex(was, cur *Schema, i int) bool {
	if was == cur {
		return false
	}
	cols := was.Indexes[i].Columns
	curIdx := cur.FindIndex(cols)
	if curIdx == nil {
		return true // index dropped
	}
	wasIdx := was.FindIndex(cols)
	return curIdx != wasIdx // index modified
}

func (mu *MergeUpdate) Skip() bool {
	return mu.table == ""
}

// ApplyMerge applies the updates collected by Merge.
// It is called by state.go Database.Merge, inside UpdateState.
func (m *Meta) ApplyMerge(updates []MergeUpdate) {
	info := m.info.Mutable()
	for _, up := range updates {
		if ts, ok := m.schema.Get(up.table); ok && !ts.isTomb() && ts.Id != up.idTran {
			continue // table recreated
		}
		if ti := info.GetCopy(up.table); ti != nil { // not dropped
			ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
			for i, ov := range ti.Indexes {
				ti.Indexes[i] = ov.WithMerged(up.results[i], up.nmerged)
			}
			info.Put(ti)
		}
	}
	m.info = info.Freeze()
}

//-------------------------------------------------------------------

type SaveResult = index.SaveResult

type PersistUpdate struct {
	table   string
	idTran  int
	results []SaveResult // per index
}

// Persist is called by state.Persist to write the state to the database.
// It collects the new btree roots which are then applied by ApplyPersist.
// WARNING: must not modify meta.
func (m *Meta) Persist(exec func(func() PersistUpdate)) {
	m.info.ForEach(func(ti *Info) {
		if len(ti.Indexes) >= 1 && ti.Indexes[0].Modified() {
			exec(func() PersistUpdate {
				results := make([]SaveResult, len(ti.Indexes))
				for i, ov := range ti.Indexes {
					results[i] = ov.Save()
				}
				id := m.schema.MustGet(ti.Table).Id
				return PersistUpdate{table: ti.Table, idTran: id, results: results}
			})
		}
	})
}

// ApplyPersist takes the new btree roots from Persist
// and updates the state with them.
func (m *Meta) ApplyPersist(updates []PersistUpdate) {
	info := m.info.Mutable()
	for _, up := range updates {
		if ts, ok := m.schema.Get(up.table); ok && !ts.isTomb() && ts.Id != up.idTran {
			continue // table recreated
		}
		if ti := info.GetCopy(up.table); ti != nil { // not dropped
			ti.Indexes = append(ti.Indexes[:0:0], ti.Indexes...) // copy
			for i, ov := range ti.Indexes {
				if up.results[i] != nil {
					ti.Indexes[i] = ov.WithSaved(up.results[i])
				}
			}
			info.Put(ti)
		}
	}
	m.info = info.Freeze()
}
