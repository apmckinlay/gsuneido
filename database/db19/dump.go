// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// DumpDatabase exports a dumped database to a file.
// It returns the number of tables dumped or panics on error.
func DumpDatabase(dbfile, to string) int {
	defer func() {
		if e := recover(); e != nil {
			os.Remove(to)
			panic("dump failed: " + fmt.Sprint(e))
		}
	}()
	db, f, w := dumpOpen(dbfile, to)
	defer db.Close()
	defer f.Close()
	state := db.GetState()
	ntables := 0
	state.meta.ForEachSchema(func(sc *meta.Schema) {
		dumpTable(db, sc, true, w)
		ntables++
	})
	ck(w.Flush())
	return ntables
}

// DumpTable exports a dumped table to a file.
// It returns the number of records dumped or panics on error.
func DumpTable(dbfile, table, to string) int {
	defer func() {
		if e := recover(); e != nil {
			os.Remove(to)
			panic("dump failed: " + table + " " + fmt.Sprint(e))
		}
	}()
	db, f, w := dumpOpen(dbfile, to)
	defer db.Close()
	defer f.Close()
	state := db.GetState()
	schema := state.meta.GetRoSchema(table)
	if schema == nil {
		panic("can't find " + table)
	}
	n := dumpTable(db, schema, false, w)
	ck(w.Flush())
	return n
}

func dumpOpen(dbfile, to string) (*Database, *os.File, *bufio.Writer) {
	db, err := OpenDatabaseRead(dbfile)
	ck(err)
	os.Remove(to) //TODO bak
	f, err := os.Create(to)
	ck(err)
	w := bufio.NewWriter(f)
	w.WriteString("Suneido dump 2\n")
	return db, f, w
}

var intbuf [4]byte

func dumpTable(db *Database, schema *meta.Schema, multi bool, w *bufio.Writer) int {
	state := db.GetState()
	w.WriteString("====== ")
	if multi {
		w.WriteString(schema.Table + " ")
	}
	w.WriteString(schema.String() + "\n")
	info := state.meta.GetRoInfo(schema.Table)
	n := 0
	writeInt := func(n int) {
		binary.BigEndian.PutUint32(intbuf[:], uint32(n))
		w.Write(intbuf[:])
	}
	iter := info.Indexes[0].Iter()
	for {
		_, off, ok := iter()
		if !ok {
			break
		}
		n++
		rec := offToRecCk(db.store, off) // verify data checksums
		writeInt(rec.Len())
		w.WriteString(string(rec))
	}
	writeInt(0)
	assert.This(n).Is(info.Nrows)
	return n
}

//TODO squeeze records when table has deleted fields
