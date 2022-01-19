// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/options"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/sortlist"
)

type loadJob struct {
	sch   schema.Schema
	list  *sortlist.Builder
	nrecs int
	size  uint64
	db    *Database
}

// LoadDatabase imports a dumped database from a file.
// It returns the number of tables loaded or panics on error.
func LoadDatabase(from, dbfile string) int {
	defer func() {
		if e := recover(); e != nil {
			panic("load failed: " + fmt.Sprint(e))
		}
	}()
	f, r := open(from)
	defer f.Close()
	db, tmpfile := tmpdb()
	defer func() { db.Close(); os.Remove(tmpfile) }()
	var wg sync.WaitGroup
	channel := make(chan loadJob)
	for i := 0; i < options.Nworkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range channel {
				loadTable2(job.db, job.sch, job.list, job.nrecs, job.size, false)
			}
		}()
	}
	nTables := 0
	for ; ; nTables++ {
		schema := readLinePrefixed(r, "====== ")
		if schema == "" {
			break
		}
		loadTable(db, r, schema, channel)
		trace()
	}
	close(channel)
	wg.Wait()
	trace("SIZE", db.Store.Size())
	db.GetState().Write()
	db.Close()
	ck(RenameBak(tmpfile, dbfile))
	return nTables
}

// LoadTable imports a dumped table from a file.
// It will replace an already existing table.
// It returns the number of records loaded or panics on error.
func LoadTable(table, dbfile string) int {
	defer func() {
		if e := recover(); e != nil {
			panic("load failed: " + table + " " + fmt.Sprint(e))
		}
	}()
	var db *Database
	var err error
	if _, err = os.Stat(dbfile); os.IsNotExist(err) {
		db, err = CreateDatabase(dbfile)
	} else {
		db, err = OpenDatabase(dbfile)
	}
	ck(err)
	defer db.Close()
	return LoadDbTable(table, db)
}

// LoadDbTable is use by dbms.Load
func LoadDbTable(table string, db *Database) int {
	db.AddExclusive(table)
	defer func() {
		if e := recover(); e != nil {
			db.EndExclusive(table)
			panic(e)
		}
	}()
	f, r := open(table + ".su")
	defer f.Close()
	schema := table + " " + readLinePrefixed(r, "====== ")
	nrecs := loadTable(db, r, schema, nil)
	db.Persist() // for safety, not strictly required
	return nrecs
}

func open(filename string) (*os.File, *bufio.Reader) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	readLinePrefixed(r, "Suneido dump 2")
	return f, r
}

func loadTable(db *Database, r *bufio.Reader, schema string, channel chan loadJob) int {
	trace(schema)
	if strings.HasPrefix(schema, "views") {
		return loadViews(db, r, schema)
	}
	sch := query.NewAdminParser(schema).Schema()
	store := db.Store
	list := sortlist.NewUnsorted()
	nrecs, size := readRecords(r, store, list)
	trace("nrecs", nrecs, "data size", size)
	list.Finish()
	if channel == nil { // not concurrent
		loadTable2(db, sch, list, nrecs, size, true)
	} else {
		channel <- loadJob{db: db, sch: sch, list: list, nrecs: nrecs, size: size}
	}
	return nrecs
}

// loadTable2 is multi-threaded when loading entire database.
func loadTable2(db *Database, sch schema.Schema, list *sortlist.Builder, nrecs int, size uint64, exclusive bool) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("ERROR:", sch.Table, e) //FIXME
		}
	}()
	ts := &meta.Schema{Schema: sch}
	ovs := buildIndexes(ts, list, db.Store, nrecs)
	ti := &meta.Info{Table: sch.Table, Nrows: nrecs, Size: size, Indexes: ovs}
	if exclusive {
		db.RunEndExclusive(sch.Table, func() {
			db.OverwriteTable(ts, ti)
		})
	} else {
		db.OverwriteTable(ts, ti)
	}
}

func readLinePrefixed(r *bufio.Reader, pre string) string {
	s, err := r.ReadString('\n') // file header
	if err == io.EOF {
		return ""
	}
	ck(err)
	if !strings.HasPrefix(s, pre) {
		panic("not a valid dump file")
	}
	return s[len(pre):]
}

func readRecords(in *bufio.Reader, store *stor.Stor, list *sortlist.Builder) (
	nrecs int, size uint64) {
	intbuf := make([]byte, 4)
	for { // each record
		_, err := io.ReadFull(in, intbuf)
		if err == io.EOF {
			break
		}
		ck(err)
		n := int(binary.BigEndian.Uint32(intbuf))
		if n == 0 {
			break
		}
		off, buf := store.Alloc(n + cksum.Len)
		_, err = io.ReadFull(in, buf[:n])
		ck(err)
		cksum.Update(buf)
		list.Add(off)
		nrecs++
		size += uint64(n)
	}
	return nrecs, size
}

func buildIndexes(ts *meta.Schema, list *sortlist.Builder, store *stor.Stor, nrecs int) []*index.Overlay {
	ts.Ixspecs(ts.Indexes)
	ov := make([]*index.Overlay, len(ts.Indexes))
	for i := range ts.Indexes {
		ix := ts.Indexes[i]
		trace(ix)
		if i > 0 || ix.Mode != 'k' {
			list.Sort(MakeLess(store, &ix.Ixspec))
		}
		before := store.Size()
		bldr := btree.Builder(store)
		iter := list.Iter()
		n := 0
		for off := iter(); off != 0; off = iter() {
			bldr.Add(btree.GetLeafKey(store, &ix.Ixspec, off), off)
			n++
		}
		bt := bldr.Finish()
		bt.SetIxspec(&ix.Ixspec)
		ov[i] = index.OverlayFor(bt)
		assert.This(n).Is(nrecs)
		trace("size", store.Size()-before)
	}
	return ov
}

func loadViews(db *Database, in *bufio.Reader, schema string) int {
	assert.That(strings.HasPrefix(schema, "views (view_name,view_definition)"))
	intbuf := make([]byte, 4)
	buf := make([]byte, 32768)
	nrecs := 0
	for { // each record
		_, err := io.ReadFull(in, intbuf)
		if err == io.EOF {
			break
		}
		ck(err)
		n := int(binary.BigEndian.Uint32(intbuf))
		if n == 0 {
			break
		}
		_, err = io.ReadFull(in, buf[:n])
		ck(err)

		rec := rt.Record(string(buf[:n]))
		name := rec.GetStr(0)
		def := rec.GetStr(1)
		db.AddView(name, def)

		nrecs++
	}
	return nrecs
}

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func trace(...interface{}) {
	// fmt.Println(args...) // comment out to disable tracing
}
