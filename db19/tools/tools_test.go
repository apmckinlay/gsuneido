// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools_test

import (
	"fmt"
	"os"
	"strings"
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

func TestCompactScanOrder(t *testing.T) {
	assert.TestOnlyIndividually(t)
	t.Chdir(t.TempDir())
	makeSuTran := db19.MakeSuTran
	db19.MakeSuTran = func(*db19.UpdateTran) *core.SuTran {
		return core.NewSuTran(nil, true)
	}
	defer func() { db19.MakeSuTran = makeSuTran }()
	db, err := db19.CreateDatabase(dbName)
	assert.T(t).This(err).Is(nil)
	db19.StartConcur(db, 50*time.Millisecond)
	query.DoAdminTest(db, "create parent (id,code) key(id) key(code)")
	query.DoAdminTest(db, "create child (id,parent) key(id) index(parent) in parent(id)")
	for i := range 20 {
		ut := db.NewUpdateTran(core.AllPerms)
		query.DoAction(nil, ut, fmt.Sprintf("insert {id: %d, code: %d} into parent", i, 20-i))
		query.DoAction(nil, ut, fmt.Sprintf("insert {id: %d, parent: %d} into child", 19-i, i))
		ut.Commit()
	}
	db.Close()
	compact := func(busy string, first string) {
		t.Helper()
		if busy != "" {
			assert.T(t).This(os.WriteFile("busycols.txt", []byte(busy), 0644)).Is(nil)
		}
		_, _, _, _, err := tools.Compact(dbName)
		assert.T(t).This(err).Is(nil)
		assert.T(t).This(db19.CheckDatabase(dbName, true)).Is(nil)
		db, err := db19.OpenDb(dbName, stor.Read, false)
		assert.T(t).This(err).Is(nil)
		defer db.Close()
		assert.T(t).This(db.GetState().Meta.GetRoSchema("parent").Indexes[0].Columns).
			Is([]string{first})
		assert.T(t).This(db.GetState().Meta.GetRoSchema("child").Indexes[0].Columns).
			Is([]string{"parent"})
		rt := db.NewReadTran(core.AllPerms)
		defer rt.Abort()
		for _, table := range []string{"parent", "child"} {
			iter := rt.IndexIter(table, 0)
			prev, n := uint64(0), 0
			for iter.Next(rt); !iter.Eof(); iter.Next(rt) {
				off := iter.CurOff()
				assert.T(t).That(off > prev)
				prev = off
				n++
			}
			assert.T(t).This(n).Is(20)
		}
	}
	compact("100\tparent\tcode\n100\tchild\tparent\n", "code")
	_, _, err = tools.DumpDatabase(dbName, "dump.tmp")
	assert.T(t).This(err).Is(nil)
	_, _, err = tools.LoadDatabase("dump.tmp", "loaded.tmp", "", "")
	assert.T(t).This(err).Is(nil)
	assert.T(t).This(db19.CheckDatabase("loaded.tmp", true)).Is(nil)
	loaded, err := db19.OpenDb("loaded.tmp", stor.Read, false)
	assert.T(t).This(err).Is(nil)
	assert.T(t).This(loaded.GetState().Meta.GetRoSchema("parent").Indexes[0].Columns).
		Is([]string{"code"})
	assert.T(t).This(loaded.GetState().Meta.GetRoSchema("child").Indexes[0].Columns).
		Is([]string{"parent"})
	loaded.Close()
	compact("", "code")
	compact("2\tparent\tid\n", "code")
	compact("100\tparent\tcode\n100\tparent\tid\n", "code")
	compact("200\tparent\tid\n100\tparent\tcode\n", "id")
}

func TestTools(t *testing.T) {
	assert.TestOnlyIndividually(t)
	t.Chdir(t.TempDir()) // isolate database files and busycols.txt
	createDb()
	defer os.Remove(dbName)
	_, _, err := tools.DumpDatabase(dbName, "dump_"+dbName)
	ck(err)
	_, _, err = tools.LoadDatabase("dump_"+dbName, "loaded_"+dbName, "", "")
	ck(err)
	ck(db19.CheckDatabase("loaded_"+dbName, true))
	_, _, err = tools.DumpDatabase("loaded_"+dbName, "dump2_"+dbName)
	ck(err)
	compare("dump_"+dbName, "dump2_"+dbName)
	_, _, _, _, err = tools.Compact(dbName)
	ck(err)
	ck(db19.CheckDatabase(dbName, true))
	compareDb(dbName, "loaded_"+dbName)
	_, _, err = tools.DumpDatabase(dbName, "dump3_"+dbName)
	ck(err)
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
		query.DoAdminTest(db, admin)
	}
	act := func(act string) {
		ut := db.NewUpdateTran(core.AllPerms)
		defer ut.Commit()
		query.DoAction(nil, ut, act)
	}
	for _, table := range []string{"foo", "bar"} {
		adm("create " + table + " " + schema)
		for _, d := range data {
			var sb strings.Builder
			for i := range d {
				sb.WriteString(columns[i])
				sb.WriteString(": '")
				sb.WriteString(d[i])
				sb.WriteString("' ")
			}
			act("insert {" + sb.String() + "} into " + table)
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
