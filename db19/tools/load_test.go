// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLoadTable(*testing.T) {
	if testing.Short() {
		return
	}
	defer os.Remove("tmp.su")
	_, err := DumpTable("../../suneido.db", "configlib", "tmp.su")
	ck(err)
	t := time.Now()
	defer os.Remove("tmp.db")
	os.Remove("tmp.db")
	n, err := LoadTable("tmp", "tmp.db")
	ck(err)
	fmt.Println("loaded", n, "records in", time.Since(t).Round(time.Millisecond))
	ck(CheckDatabase("tmp.db"))
}

func TestLoadDbTable(*testing.T) {
	if testing.Short() {
		return
	}
	defer os.Remove("tmp.su")
	_, err := DumpTable("../../suneido.db", "configlib", "tmp.su")
	ck(err)

	defer os.Remove("tmp.db")
	os.Remove("tmp.db")
	db, err := CreateDatabase("tmp.db")
	ck(err)
	StartConcur(db, 100*time.Millisecond)

	t := time.Now()
	n, err := LoadDbTable("tmp", "tmp.su", "", "", db)
	ck(err)
	fmt.Println("loaded", n, "records in", time.Since(t).Round(time.Millisecond))

	db.Close()
	ck(CheckDatabase("tmp.db"))
}

func TestLoadDatabase(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	defer os.Remove("tmp.db")
	nTables, nViews, err := LoadDatabase("../../database.su", "tmp.db", "", "")
	ck(err)
	fmt.Println("loaded", nTables, "tables", nViews, "views in",
		time.Since(t).Round(time.Millisecond))
	ck(CheckDatabase("tmp.db"))
}

func TestLoadFkey(*testing.T) {
	if testing.Short() {
		return
	}
	db := CreateDb(stor.HeapStor(8192))
	doAdmin := func(cmd string) {
		query.DoAdmin(db, cmd, nil)
	}
	doAdmin("create tmp (a) key(a)")
	doAdmin("create tmp2 (k, a) key(k) index(a) in tmp")
	doAdmin("create tmp3 (k, a) key(k) index(a)")
	_, err := DumpDbTable(db, "tmp", "tmp.su", "")
	ck(err)
	defer os.Remove("tmp.su")
	_, err = DumpDbTable(db, "tmp2", "tmp2.su", "")
	ck(err)
	defer os.Remove("tmp2.su")
	_, err = DumpDbTable(db, "tmp3", "tmp3.su", "")
	ck(err)
	defer os.Remove("tmp3.su")
	_, err = LoadDbTable("tmp", "tmp.su", "", "", db)
	assert.That(strings.Contains(err.Error(),
		"can't overwrite table that foreign keys point to"))
	_, err = LoadDbTable("tmp2", "tmp2.su", "", "", db)
	assert.That(strings.Contains(err.Error(),
		"can't load single table with foreign keys"))
	_, err = LoadDbTable("tmp3", "tmp3.su", "", "", db)
	ck(err)
	doAdmin("drop tmp3")
	_, err = LoadDbTable("tmp3", "tmp3.su", "", "", db)
	ck(err)
}
