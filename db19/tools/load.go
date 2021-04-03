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

	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/sortlist"
)

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
	nTables := 0
	for ; ; nTables++ {
		schema := readLinePrefixed(r, "====== ")
		if schema == "" {
			break
		}
		loadTable(db, r, schema)
		trace()
		assert.That(nTables < 1010)
	}
	trace("SIZE", db.Store.Size())
	db.GetState().Write(true)
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
	if _, err := os.Stat(dbfile); os.IsNotExist(err) {
		db, err = CreateDatabase(dbfile)
	} else {
		db, err = OpenDatabase(dbfile)
		db.DropTable(table)
	}
	ck(err)
	defer db.Close()
	f, r := open(table + ".su")
	defer f.Close()
	schema := table + " " + readLinePrefixed(r, "====== ")
	nrecs := loadTable(db, r, schema)
	db.GetState().Write(true)
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

func loadTable(db *Database, r *bufio.Reader, schema string) int {
	trace(schema)
	sch := query.NewRequestParser(schema).Schema()

	store := db.Store
	list := sortlist.NewUnsorted()
	before := store.Size()
	nrecs := readRecords(r, store, list)
	beforeIndexes := store.Size()
	dataSize := beforeIndexes - before
	trace("nrecs", nrecs, "data size", dataSize)
	list.Finish()
	ts := &meta.Schema{Schema: sch}
	ov := buildIndexes(ts, list, store, nrecs)
	trace("indexes size", store.Size()-beforeIndexes)
	ti := &meta.Info{Table: sch.Table, Nrows: nrecs, Size: dataSize, Indexes: ov}
	db.LoadedTable(ts, ti)
	return nrecs
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

func readRecords(in *bufio.Reader, store *stor.Stor, list *sortlist.Builder) int {
	nrecs := 0
	intbuf := make([]byte, 4)
	for { // each record
		_, err := io.ReadFull(in, intbuf)
		if err == io.EOF {
			break
		}
		ck(err)
		size := int(binary.BigEndian.Uint32(intbuf))
		if size == 0 {
			break
		}
		off, buf := store.Alloc(size + cksum.Len)
		_, err = io.ReadFull(in, buf[:size])
		ck(err)
		cksum.Update(buf)
		list.Add(off)
		nrecs++
	}
	return nrecs
}

func buildIndexes(ts *meta.Schema, list *sortlist.Builder, store *stor.Stor, nrecs int) []*index.Overlay {
	ts.Ixspecs()
	ov := make([]*index.Overlay, len(ts.Indexes))
	for i := range ts.Indexes {
		ix := ts.Indexes[i]
		trace(ix)
		if i > 0 || ix.Mode != 'k' {
			list.Sort(mkcmp(store, &ix.Ixspec))
		}
		before := store.Size()
		bldr := btree.Builder(store)
		iter := list.Iter()
		n := 0
		for off := iter(); off != 0; off = iter() {
			bldr.Add(btree.GetLeafKey(store, &ix.Ixspec, off), off)
			n++
		}
		ov[i] = index.OverlayFor(bldr.Finish())
		assert.This(n).Is(nrecs)
		trace("size", store.Size()-before)
	}
	return ov
}

func mkcmp(store *stor.Stor, is *ixkey.Spec) func(x, y uint64) int {
	return func(x, y uint64) int {
		xr := OffToRec(store, x)
		yr := OffToRec(store, y)
		return is.Compare(xr, yr)
	}
}

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func trace(...interface{}) {
	// fmt.Println(args...) // comment out to disable tracing
}
