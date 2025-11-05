// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools_test

import (
	"os"
	"testing"
	"time"

	"slices"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/db19/tools"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

const dbName = "testtools.tmp"

const schema = "(one, two, three, four) key(one) index(four, two)"

var columns = []string{"one", "two", "three", "four"}

var data = [][]string{
	{"a", "b", "c", "d"},
	{"e", "f", "g", "h"},
}

func TestTools(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow TestTools")
	}
	createDb()
	defer os.Remove(dbName)
	_, _, err := tools.DumpDatabase(dbName, "dump_"+dbName)
	ck(err)
	defer os.Remove("dump_" + dbName)
	_, _, err = tools.LoadDatabase("dump_"+dbName, "loaded_"+dbName, "", "")
	ck(err)
	defer os.Remove("loaded" + dbName)
	ck(db19.CheckDatabase("loaded_" + dbName, true))
	defer os.Remove("loaded_" + dbName)
	tools.DumpDatabase("loaded_"+dbName, "dump2_"+dbName)
	defer os.Remove("dump2_" + dbName)
	compare("dump_"+dbName, "dump2_"+dbName)
	tools.Compact(dbName)
	defer os.Remove(dbName + ".bak")
	ck(db19.CheckDatabase(dbName, true))
	compareDb(dbName, "loaded_"+dbName)
	tools.DumpDatabase(dbName, "dump3_"+dbName)
	defer os.Remove("dump3_" + dbName)
	compare("dump_"+dbName, "dump3_"+dbName)
}

func createDb() {
	store, err := stor.MmapStor(dbName, stor.Create)
	ck(err)
	defer store.Close(true)
	db := db19.CreateDb(store)
	db19.StartConcur(db, 50*time.Millisecond)
	db19.MakeSuTran = func(ut *db19.UpdateTran) *core.SuTran {
		return core.NewSuTran(nil, true)
	}
	adm := func(admin string) {
		query.DoAdmin(db, admin, nil)
	}
	act := func(act string) {
		ut := db.NewUpdateTran()
		defer ut.Commit()
		query.DoAction(nil, ut, act)
	}
	for _, table := range []string{"foo", "bar"} {
		adm("create " + table + " " + schema)
		for _, d := range data {
			s := ""
			for i := range d {
				s += columns[i] + ": '" + d[i] + `' `
			}
			act("insert {" + s + "} into " + table)
		}
	}
	adm("alter bar drop (three)")
	db.Close()
}

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func compare(file1, file2 string) {
	b1, err := os.ReadFile(file1)
	ck(err)
	b2, err := os.ReadFile(file2)
	ck(err)
	if !slices.Equal(b1, b2) {
		panic(file1 + " is NOT the same as " + file2)
	}
}

func compareDb(db1, db2 string) {
	assert.This(getSchema(db1)).Is(getSchema(db2))
}

func getSchema(dbName string) string {
	db, err := db19.OpenDatabase(dbName)
	ck(err)
	defer db.Close()
	return db.Schema("foo") + "\n" + db.Schema("bar")
}
