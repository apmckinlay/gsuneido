// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/apmckinlay/gsuneido/database/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// DumpDatabase exports a dumped database to a file.
// It returns the number of tables dumped or panics on error.
func DumpDatabase(dbfile, to string) int {
	defer func() {
		if e := recover(); e != nil {
			panic("dump failed: " + fmt.Sprint(e))
		}
	}()
	db, f, w := dumpOpen(dbfile, to)
	defer db.Close()
	defer f.Close()
	state := db.GetState()
	ntables := 0
	state.meta.ForEachSchema(func(sc *schema.Schema) {
		dumpTable(db, sc.Table, true, w)
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
			panic("dump failed: " + table + " " + fmt.Sprint(e))
		}
	}()
	db, f, w := dumpOpen(dbfile, to)
	defer db.Close()
	defer f.Close()
	n := dumpTable(db, table, false, w)
	ck(w.Flush())
	return n
}

func dumpOpen(dbfile, to string) (*Database, *os.File, *bufio.Writer) {
	db := OpenDatabaseRead(dbfile)
	os.Remove(to) //TODO bak
	f, err := os.Create(to)
	ck(err)
	w := bufio.NewWriter(f)
	w.WriteString("Suneido dump 2\n")
	return db, f, w
}

var intbuf [4]byte

func dumpTable(db *Database, table string, multi bool, w *bufio.Writer) int {
	state := db.GetState()
	schema := state.meta.GetRoSchema(table)
	if schema == nil {
		panic("can't find " + table)
	}
	w.WriteString("====== ")
	if multi {
		w.WriteString(table + " ")
	}
	w.WriteString(schema.String() + "\n")
	info := state.meta.GetRoInfo(table)
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
		rec := offToRec(db.store, off)
		writeInt(rec.Len())
		w.WriteString(string(rec))
	}
	writeInt(0)
	assert.This(n).Is(info.Nrows)
	return n
}

//TODO squeeze records when table has deleted fields
