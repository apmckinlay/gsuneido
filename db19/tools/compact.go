// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/apmckinlay/gsuneido/core"
	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stats"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/dbg"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/slc"
	"github.com/apmckinlay/gsuneido/util/sortlist"
	"github.com/apmckinlay/gsuneido/util/system"
)

// Compact cleans up old records and index nodes that are no longer in use.
// It does this by copying live data to a new database file.
// In the process it concurrently does a full check of the database.
// To keep table data and indexes contiguous (not interleaved)
// it runs primarily single threaded, only index checking is concurrent.
func Compact(dbfile string) (nTables, nViews int, oldSize, newSize uint64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("compact failed: %v", e)
			dbg.PrintStack()
		}
	}()
	src, err := OpenDb(dbfile, stor.Read, false)
	ck(err)
	defer src.Close()
	oldSize = src.Store.Size()
	dst, tmpfile := tmpdb()
	defer func() { dst.Close(); os.Remove(tmpfile) }()

	state := src.GetState()
	busy := stats.LoadBusy()
	stats := make(stats.StatsTally)

	type schemaSize struct {
		sc    *meta.Schema
		nrows int
	}

	nViews = copyViews(state, dst)

	schemas := make([]schemaSize, 0, 128)
	for sc := range state.Meta.Tables() {
		ti := state.Meta.GetRoInfo(sc.Table)
		ss := schemaSize{sc: sc, nrows: ti.Nrows}
		schemas = append(schemas, ss)
		nTables++
	}
	// sort largest first to minimize load tail (like dump)
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].nrows > schemas[j].nrows
	})
	var wg sync.WaitGroup
	channel := make(chan indexJob, 8)
	for range options.Nworkers {
		wg.Go(func() {
			for job := range channel {
				CheckOtherIndex(job.store, job.tsi, job.ix, job.nrows, job.sum, false, nil)
			}
		})
	}
	for _, sc := range schemas {
		if sc.sc.Table == StatsTableName {
			continue
		}
		compactTable(state, src, sc.sc, dst, channel, busy, stats)
	}
	close(channel)
	wg.Wait()
	createStatsTable(dst, stats)
	dst.GetState().Write()
	newSize = dst.Store.Size()
	dst.Close()
	src.Close()
	ck(system.RenameBak(tmpfile, dbfile))
	return nTables, nViews, oldSize, newSize, nil
}

func tmpdb() (*Database, string) {
	dst, err := os.CreateTemp(".", "gs*.tmp")
	ck(err)
	tmpfile := dst.Name()
	dst.Close()
	db, err := CreateDatabase(tmpfile)
	ck(err)
	return db, tmpfile
}

func copyViews(state *DbState, dst *Database) int {
	n := 0
	for name, def := range state.Meta.Views() {
		dst.AddView(name, def)
		n++
	}
	return n
}

type indexJob struct {
	store *stor.Stor
	tsi   *schema.Index
	ix    *index.Overlay
	nrows int
	sum   uint64
}

func compactTable(state *DbState, src *Database, ts *meta.Schema, dst *Database, channel chan<- indexJob, busy stats.Busy, stats stats.StatsTally) {
	defer func() {
		if e := recover(); e != nil {
			core.Fatal(ts.Table+":", e)
		}
	}()
	ts.Check(state.Meta.GetRoSchema)
	hasdel := ts.HasDeleted()
	info := state.Meta.GetRoInfo(ts.Table)
	ixi := info.SmallestKeyIndex(ts.Indexes)
	statsCols := buildStatsCols(ts.Table, ts.Columns, busy, stats)
	sum := uint64(0)
	size := int64(0)
	list := sortlist.NewUnsorted(func(x uint64) bool { return x == 0 })
	var off2 uint64
	var dstbuf []byte
	nrows := info.Indexes[ixi].CheckBtree(func(off uint64) {
		sum += off // addition so order doesn't matter
		buf := src.Store.Data(off)
		n := core.RecLen(buf)
		buf = buf[:n+cksum.Len]
		cksum.MustCheck(buf)
		rec := core.Record(hacks.BStoS(buf[:n]))
		addStats(rec, statsCols)
		if hasdel {
			rec = squeeze(rec, ts.Columns)
			n = len(rec)
			off2, dstbuf = dst.Store.Alloc(n + cksum.Len)
			copy(dstbuf, rec)
			cksum.Update(dstbuf)
		} else {
			off2, dstbuf = dst.Store.Alloc(len(buf))
			copy(dstbuf, buf)
		}
		list.Add(off2)
		size += int64(n)
	})
	list.Finish()
	assert.That(nrows == info.Nrows)
	for i := range info.Indexes {
		if i == ixi {
			continue
		}
		channel <- indexJob{store: src.Store,
			tsi: &ts.Indexes[i], ix: info.Indexes[i], nrows: nrows, sum: sum}
	}
	if hasdel {
		ts.Columns = slc.Without(ts.Columns, "-")
	}
	indexes := buildIndexes(ts, list, dst, nrows, ixi)
	ti := meta.NewInfo(ts.Table, indexes, nrows, size)
	dst.AddNewTable(ts, ti)
}

type statsCol struct {
	idx int
	sc  *stats.StatsColumn
}

func buildStatsCols(table string, cols []string,
	busy stats.Busy, tally stats.StatsTally) []statsCol {
	busyCols := busy[table]
	if len(busyCols) == 0 {
		return nil
	}
	ht, ok := tally[table]
	if !ok {
		ht = make(stats.StatsTable)
		tally[table] = ht
	}
	for _, col := range busyCols {
		if _, ok := ht[col]; !ok {
			ht[col] = &stats.StatsColumn{}
		}
	}
	colIdx := make(map[string]int, len(cols))
	for i, col := range cols {
		if col != "-" {
			colIdx[col] = i
		}
	}
	var statsCols []statsCol
	for _, col := range busyCols {
		if i, ok := colIdx[col]; ok {
			statsCols = append(statsCols, statsCol{idx: i, sc: ht[col]})
		}
	}
	if len(statsCols) == 0 {
		return nil
	}
	sort.Slice(statsCols,
		func(i, j int) bool { return statsCols[i].idx < statsCols[j].idx })
	return statsCols
}

func addStats(rec core.Record, stats []statsCol) {
	for _, e := range stats {
		val := rec.GetRaw(e.idx)
		e.sc.Add(val)
	}
}

func createStatsTable(dst *Database, stats stats.StatsTally) {
	stats.Complete()
	var rb core.RecordBuilder
	rb.Add(stats)
	bi := rb.PreBuild()
	n := bi.Size() + cksum.Len
	off, buf := dst.Store.Alloc(n)
	rec := rb.BuildInto(buf[:0], bi)
	cksum.Update(buf)

	ts := &meta.Schema{
		Table:   StatsTableName,
		Columns: []string{"data"},
		Indexes: []schema.Index{{Mode: 'k'}}}
	ts.SetupIndexes()

	ix := &ts.Indexes[0]
	key := ix.Ixspec.Key(rec)
	bldr := dst.BtreeBuilder()
	if !bldr.Add(key, off) {
		panic("cannot build " + StatsTableName + " index")
	}
	bt := bldr.Finish()
	ov := index.OverlayFor(bt)

	ti := meta.NewInfo(ts.Table, []*index.Overlay{ov}, 1, int64(bi.Size()))
	dst.AddNewTable(ts, ti)
}
