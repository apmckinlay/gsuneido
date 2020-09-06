// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"bufio"
	"encoding/binary"
	"os"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// DumpDatabase exports a dumped database to a file.
// It returns the number of tables dumped or panics on error.
func DumpDatabase(from, to string) int {
	return 0
}

// DumpTable exports a dumped table to a file.
// It returns the number of records dumped or panics on error.
func DumpTable(dbfile, table, to string) int {
	ck := func(err error) {
		if err != nil {
			panic("dump failed: " + table + " " + err.Error())
		}
	}
	db := OpenDatabaseRead(dbfile)
	defer db.Close()
	os.Remove(to) //TODO bak
	f, err := os.Create(to)
	ck(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	w.WriteString("Suneido dump 2\n")

	state := db.GetState()
	schema := state.meta.GetRoSchema(table)
	w.WriteString("====== " + schema.String() + "\n")
	info := state.meta.GetRoInfo(table)
	n := 0
	intbuf := make([]byte, 4)
	iter := info.Indexes[0].Iter()
	for {
		_, off, ok := iter()
		if !ok {
			break
		}
		n++
		rec := offToRec(db.store, off)
		binary.BigEndian.PutUint32(intbuf, uint32(rec.Len()))
		w.Write(intbuf)
		w.WriteString(string(rec))
	}
	ck(w.Flush())
	assert.This(n).Is(info.Nrows)
	return n
}

//TODO squeeze records when table has deleted fields
