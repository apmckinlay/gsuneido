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
	t := time.Now()
	defer os.Remove("tmp.db")
	os.Remove("tmp.db")
	n, err := LoadTable("stdlib", "tmp.db")
	assert.This(err).Is(nil)
	fmt.Println("loaded", n, "records in", time.Since(t).Round(time.Millisecond))
	ck(CheckDatabase("tmp.db"))
}

func TestLoadDatabase(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	defer os.Remove("tmp.db")
	nTables, nViews, e := LoadDatabase("../../database.su", "tmp.db")
	assert.This(e).Is(nil)
	fmt.Println("loaded", nTables, "tables", nViews, "views in",
		time.Since(t).Round(time.Millisecond))
	ck(CheckDatabase("tmp.db"))
}

func TestLoadFkey(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := CreateDb(store)
	ck(err)
	doAdmin := func(cmd string) {
		query.DoAdmin(db, cmd, nil)
	}
	doAdmin("create tmp (a) key(a)")
	doAdmin("create tmp2 (k, a) key(k) index(a) in tmp")
	_, err = DumpDbTable(db, "tmp2", "tmp2.su")
	ck(err)
	doAdmin("drop tmp2")
	doAdmin("drop tmp")
	_, err = LoadDbTable("tmp2", db)
	os.Remove("tmp2.su")
	assert.That(strings.Contains(err.Error(),
		"can't create foreign key to nonexistent index"))
}
