// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
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
	db, err := openDatabase(dbfile, stor.READ, false)
	ck(err)
	os.Remove(to) //TODO bak
	f, err := os.Create(to)
	ck(err)
	w := bufio.NewWriter(f)
	w.WriteString("Suneido dump 2\n")
	return db, f, w
}

func dumpTable(db *Database, schema *meta.Schema, multi bool, w *bufio.Writer) int {
	state := db.GetState()
	w.WriteString("====== ")
	if multi {
		w.WriteString(schema.Table + " ")
	}
	w.WriteString(schema.String() + "\n")
	info := state.meta.GetRoInfo(schema.Table)
	sum := uint64(0)
	count := info.Indexes[0].Check(func(off uint64) {
		sum += off                       // addition so order doesn't matter
		rec := offToRecCk(db.store, off) // verify data checksums
		writeInt(w, len(rec))
		w.WriteString(string(rec))
	})
	writeInt(w, 0)                      // end of table records
	checkOtherIndexes(info, count, sum) //TODO concurrent
	assert.This(count).Is(info.Nrows)
	return count
}

func writeInt(w *bufio.Writer, n int) {
	assert.That(0 <= n && n <= math.MaxUint32)
	w.WriteByte(byte(n >> 24))
	w.WriteByte(byte(n >> 16))
	w.WriteByte(byte(n >> 8))
	w.WriteByte(byte(n))
}

//TODO squeeze records when table has deleted fields
