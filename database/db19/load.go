// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/comp"
	"github.com/apmckinlay/gsuneido/database/db19/ixspec"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/sortlist"
)

// LoadDatabase imports a dumped database from a file.
// It returns the number of tables loaded or panics on error.
func LoadDatabase() int {
	defer func() {
		if e := recover(); e != nil {
			panic("load failed: " + fmt.Sprint(e))
		}
	}()
	f, r, store := open("database.su")
	defer f.Close()
	defer store.Close()
	nTables := 0
	for ; ; nTables++ {
		schema := readLinePrefixed(r, "====== ")
		if schema == "" {
			break
		}
		loadTable(store, r, schema)
		trace()
		assert.That(nTables < 1010)
	}
	trace("SIZE", store.Size())
	return nTables
}

// LoadTable imports a dumped table from a file.
// It returns the number of records loaded or panics on error.
func LoadTable(filename string) int {
	table := filename
	if strings.HasSuffix(table, ".su") {
		table = filename[:len(filename)-3]
	}
	defer func() {
		if e := recover(); e != nil {
			panic("load failed: " + table + " " + fmt.Sprint(e))
		}
	}()
	f, r, store := open(filename)
	defer f.Close()
	defer store.Close()
	schema := table + " " + readLinePrefixed(r, "====== ")
	return loadTable(store, r, schema)
}

func open(filename string) (*os.File, *bufio.Reader, *stor.Stor) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	readLinePrefixed(r, "Suneido dump 2")
	store, err := stor.MmapStor("tmp.db", stor.CREATE)
	ckerr(err)
	store.Alloc(1) // don't use offset 0
	return f, r, store
}

func loadTable(store *stor.Stor, r *bufio.Reader, schema string) int {
	trace(schema)
	req := compile.ParseRequest("create " + schema).(*compile.Schema)

	list := sortlist.NewUnsorted()
	before := store.Size()
	nrecs := readRecords(r, store, list)
	beforeIndexes := store.Size()
	trace("nrecs", nrecs, "data size", beforeIndexes-before)
	list.Finish()
	ts := reqToSchema(req)
	for i := range ts.Indexes {
		ix := ts.Indexes[i]
		trace(ix)
		if i > 0 || ix.Mode != 'k' { // non-key order may be different
			list.Sort(mkcmp(store, &ix.Ixspec))
		}
		before := store.Size()
		bldr := btree.NewFbtreeBuilder(store)
		iter := list.Iter()
		n := 0
		for off := iter(); off != 0; off = iter() {
			bldr.Add(getLeafKey(store, &ix.Ixspec, off), off)
			n++
		}
		bldr.Finish()
		assert.This(n).Is(nrecs)
		trace("size", store.Size()-before)
	}
	trace("indexes size", store.Size()-beforeIndexes)
	return nrecs
}

func reqToSchema(req *compile.Schema) *meta.Schema {
	var ts meta.Schema
	ts.Indexes = make([]meta.IndexSchema, len(req.Indexes))
	for i := range req.Indexes {
		ri := req.Indexes[i]
		ts.Indexes[i] = meta.IndexSchema{Fields: ri.Fields, Mode: ri.Mode}
	}
	ts.Ixspecs()
	return &ts
}

func readLinePrefixed(r *bufio.Reader, pre string) string {
	s, err := r.ReadString('\n') // file header
	if err == io.EOF {
		return ""
	}
	ckerr(err)
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
		ckerr(err)
		size := int(binary.BigEndian.Uint32(intbuf))
		if size == 0 {
			break
		}
		off, buf := store.Alloc(size)
		_, err = io.ReadFull(in, buf)
		ckerr(err)
		list.Add(off)
		nrecs++
	}
	return nrecs
}

func ckerr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func getLeafKey(store *stor.Stor, ix *ixspec.T, off uint64) string {
	rec := offToRec(store, off)
	return comp.Key(rt.Record(rec), ix.Cols, ix.Cols2)
}

func mkcmp(store *stor.Stor, ix *ixspec.T) func(x, y uint64) int {
	return func(x, y uint64) int {
		xr := offToRec(store, x)
		yr := offToRec(store, y)
		return comp.Compare(xr, yr, ix.Cols, ix.Cols2)
	}
}

func offToRec(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	return rt.Record(hacks.BStoS(buf))
}
