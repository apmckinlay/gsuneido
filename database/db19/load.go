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
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sortlist"
)

// LoadDatabase imports a dumped database from a file.
// It returns the number of tables loaded or panics on error.
func LoadDatabase(from,to string) int {
	defer func() {
		if e := recover(); e != nil {
			panic("load failed: " + fmt.Sprint(e))
		}
	}()
	f, r := open(from)
	defer f.Close()
	db := CreateDatabase(to) //TODO tmp & .bak
	defer db.Close()
	nTables := 0
	for ; ; nTables++ {
		schema := readLinePrefixed(r, "====== ")
		if schema == "" {
			break
		}
		loadTable(db.store, r, schema)
		trace()
		assert.That(nTables < 1010)
	}
	trace("SIZE", db.store.Size())
	return nTables
}

// LoadTable imports a dumped table from a file.
// It returns the number of records loaded or panics on error.
func (db *Database) LoadTable(table string) int {
	defer func() {
		if e := recover(); e != nil {
			panic("load failed: " + table + " " + fmt.Sprint(e))
		}
	}()
	f, r := open(table + ".su")
	defer f.Close()
	schema := table + " " + readLinePrefixed(r, "====== ")
	return loadTable(db.store, r, schema)
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
