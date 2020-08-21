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
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/sortlist"
)

//TODO handle capitalized rule fields, and _lower! fields

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
	key := firstShortestKey(req)
	for i, ix := range req.Indexes {
		ixcols := ix.Fields
		if len(ixcols) > 0 && strings.HasSuffix(req.Columns[ixcols[0]], "!") {
			continue //TODO
		}
		if req.Indexes[i].Mode != 'k' {
			ixcols = append(ixcols, key...)
		}
		trace(ix, "cols", ixcols)
		//		if i > 0 {
		list.Sort(mkcmp(store, ixcols))
		//		}
		before := store.Size()
		bldr := btree.NewFbtreeBuilder(store)
		iter := list.Iter()
		n := 0
		for off := iter(); off != 0; off = iter() {
			bldr.Add(getLeafKey(store, ixcols, off), off)
			n++
		}
		assert.This(n).Is(nrecs)
		trace("size", store.Size()-before)
	}
	trace("indexes size", store.Size()-beforeIndexes)
	return nrecs
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

func firstShortestKey(req *compile.Schema) []int {
	var key []int
	for _, ix := range req.Indexes {
		if usableKey(req, ix) &&
			(key == nil || len(ix.Fields) < len(key)) {
			key = ix.Fields
		}
	}
	return key
}

func usableKey(req *compile.Schema, ix *compile.Index) bool {
	return ix.Mode == 'k' && len(ix.Fields) > 0 &&
		!strings.HasSuffix(req.Columns[ix.Fields[0]], "!")
}

func getLeafKey(store *stor.Stor, ixspec interface{}, off uint64) string {
	rec := offToRec(store, off)
	ixcols := ixspec.([]int)
	return comp.Key(rt.Record(rec), ixcols)
}

func mkcmp(store *stor.Stor, ixcols []int) func(x, y uint64) int {
	return func(x, y uint64) int {
		xr := offToRec(store, x)
		yr := offToRec(store, y)
		return comp.Compare(xr, yr, ixcols)
	}
}

func offToRec(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	return rt.Record(hacks.BStoS(buf))
}
