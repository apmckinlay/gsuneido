// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"fmt"
	"os"
	"slices"
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
// (similar to a copying garbage collector)
// In the process it concurrently does a full check of the database.
// To keep table data and indexes contiguous (not interleaved)
// it runs primarily single threaded, only index checking is concurrent.
// It also uses information collected in busycols.txt during operation
// to collects stats about "busy" columns
// and to select the most scanned index for physical order (faster)
// and move it to first (indexes[0])
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
	busy, scans := stats.LoadBusy()
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
		compactTable(state, src, sc.sc, dst, channel, busy, scans, stats)
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

func compactTable(state *DbState, src *Database, ts *meta.Schema, dst *Database,
	channel chan<- indexJob, busy stats.Busy, scans stats.Scans, stats stats.StatsTally) {
	defer func() {
		if e := recover(); e != nil {
			core.Fatal(ts.Table+":", e)
		}
	}()
	ts.Check(state.Meta.GetRoSchema)
	hasdel := ts.HasDeleted()
	info := state.Meta.GetRoInfo(ts.Table)
	ixi := compactIndex(ts.Indexes, info.Nrows, scans[ts.Table])
	statsCols := buildStatsCols(ts.Table, ts.Columns, busy, stats)
	sum := uint64(0)
	size := 0
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
		size += n
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
	// Keep the source schema (also used by the check workers) unchanged.
	cp := *ts
	ts = &cp
	ts.Indexes = slices.Clone(ts.Indexes)
	if hasdel {
		ts.Columns = slc.Without(ts.Columns, "-")
	}
	// Index zero remembers the physical order
	if ixi != 0 {
		ix := ts.Indexes[ixi]
		copy(ts.Indexes[1:ixi+1], ts.Indexes[:ixi])
		ts.Indexes[0] = ix
	}
	indexes := buildIndexes(ts, list, dst, nrows)
	ti := meta.NewInfo(ts.Table, indexes, nrows, size)
	dst.AddNewTable(ts, ti)
}

func compactIndex(indexes []schema.Index, nrows int, scans map[string]int) int {
	if len(indexes) <= 1 {
		return 0
	}
	best := 0
	bestCount := scans[indexes[0].ColStr()]
	for i, ix := range indexes {
		ixs := ix.ColStr()
		count := scans[ixs]
		if count > bestCount {
			best = i
			bestCount = count
		}
	}
	// we use nrows as a threshold, because if a table isn't scanned even once
	// during the usual week between compacts
	// then it probably isn't worth switching physical indexes
	// since that forces compact to read by a non-physical index which is slower
	if bestCount < max(100, nrows) ||
		float64(bestCount) <= 1.15*float64(scans[indexes[0].ColStr()]) {
		return 0 // keep the current
	}
	return best
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
	if stats.IsEmpty() {
		return
	}
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

	ti := meta.NewInfo(ts.Table, []*index.Overlay{ov}, 1, bi.Size())
	dst.AddNewTable(ts, ti)
}
