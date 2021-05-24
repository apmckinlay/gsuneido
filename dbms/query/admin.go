// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
)

func DoAdmin(db *db19.Database, cmd string) {
	admin := ParseAdmin(cmd)
	admin.execute(db)
}

func checkForSystemTable(op, table string) {
	if isSystemTable(table) {
		panic("can't " + op + " system table: " + table)
	}
}

func isSystemTable(table string) bool {
	switch table {
	case "tables", "columns", "indexes", "views":
		return true
	}
	return false
}

//-------------------------------------------------------------------

type createAdmin struct {
	Schema
}

func (a *createAdmin) String() string {
	return "create " + a.Schema.String()
}

func (a *createAdmin) execute(db *db19.Database) {
	checkForSystemTable("create", a.Table)
	ts := &meta.Schema{Schema: a.Schema}
	ts.Ixspecs(ts.Indexes)
	ov := make([]*index.Overlay, len(ts.Indexes))
	for i := range ov {
		bt := btree.CreateBtree(db.Store, &ts.Indexes[i].Ixspec)
		ov[i] = index.OverlayFor(bt)
	}
	ti := &meta.Info{Table: a.Schema.Table, Indexes: ov}
	db.LoadedTable(ts, ti)
}

func createIndexes(db *db19.Database, idxs []schema.Index) []*index.Overlay {
	ov := make([]*index.Overlay, len(idxs))
	for i := range ov {
		bt := btree.CreateBtree(db.Store, &idxs[i].Ixspec)
		ov[i] = index.OverlayFor(bt)
	}
	return ov
}

//-------------------------------------------------------------------

type ensureAdmin struct {
	Schema
}

func (a *ensureAdmin) String() string {
	return "ensure " + a.Schema.String()
}

func (a *ensureAdmin) execute(db *db19.Database) {
	checkForSystemTable("ensure", a.Table)
	if !db.Ensure(&a.Schema) {
		panic("can't " + a.String())
	}
}

//-------------------------------------------------------------------

type renameAdmin struct {
	from string
	to   string
}

func (a *renameAdmin) String() string {
	return "rename " + a.from + " to " + a.to
}

func (a *renameAdmin) execute(db *db19.Database) {
	checkForSystemTable("rename", a.from)
	checkForSystemTable("rename to", a.to)
	if !db.RenameTable(a.from, a.to) {
		panic("can't " + a.String())
	}
}

//-------------------------------------------------------------------

type alterCreateAdmin struct {
	Schema
}

func (a *alterCreateAdmin) String() string {
	return "alter " + strings.Replace(a.Schema.String(), " ", " create ", 1)
}

func (a *alterCreateAdmin) execute(db *db19.Database) {
	checkForSystemTable("alter", a.Table)
	if !db.AlterCreate(&a.Schema) {
		panic("can't " + a.String())
	}
}

//-------------------------------------------------------------------

type alterRenameAdmin struct {
	table string
	from  []string
	to    []string
}

func (a *alterRenameAdmin) String() string {
	s := "alter " + a.table + " rename "
	sep := ""
	for i, from := range a.from {
		s += sep + from + " to " + a.to[i]
		sep = ", "
	}
	return s
}

func (a *alterRenameAdmin) execute(db *db19.Database) {
	checkForSystemTable("alter", a.table)
	if !db.AlterRename(a.table, a.from, a.to) {
		panic("can't " + a.String())
	}
}

//-------------------------------------------------------------------

type alterDropAdmin struct {
	Schema
}

func (a *alterDropAdmin) String() string {
	return "alter " + strings.Replace(a.Schema.String(), " ", " drop ", 1)
}

func (a *alterDropAdmin) execute(db *db19.Database) {
	checkForSystemTable("alter", a.Table)
	if !db.AlterDrop(&a.Schema) {
		panic("can't " + a.String())
	}
}

//-------------------------------------------------------------------

type dropAdmin struct {
	table string
}

func (a *dropAdmin) String() string {
	return "drop " + a.table
}

func (a *dropAdmin) execute(db *db19.Database) {
	checkForSystemTable("drop", a.table)
	if !db.DropTable(a.table) {
		panic("can't drop nonexistent table: " + a.table)
	}
}
