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

func (r *createAdmin) String() string {
	return "create " + r.Schema.String()
}

func (r *createAdmin) execute(db *db19.Database) {
	checkForSystemTable("create", r.Table)
	ts := &meta.Schema{Schema: r.Schema}
	ts.Ixspecs(ts.Indexes)
	ov := make([]*index.Overlay, len(ts.Indexes))
	for i := range ov {
		bt := btree.CreateBtree(db.Store, &ts.Indexes[i].Ixspec)
		ov[i] = index.OverlayFor(bt)
	}
	ti := &meta.Info{Table: r.Schema.Table, Indexes: ov}
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
	schema Schema
}

func (r *ensureAdmin) String() string {
	return "ensure " + r.schema.String()
}

func (r *ensureAdmin) execute(db *db19.Database) {
	//TODO
}

//-------------------------------------------------------------------

type renameAdmin struct {
	from string
	to   string
}

func (r *renameAdmin) String() string {
	return "rename " + r.from + " to " + r.to
}

func (r *renameAdmin) execute(db *db19.Database) {
	checkForSystemTable("rename", r.from)
	checkForSystemTable("rename to", r.to)
	if !db.RenameTable(r.from, r.to) {
		panic("can't " + r.String())
	}
}

//-------------------------------------------------------------------

type alterCreateAdmin struct {
	Schema
}

func (r *alterCreateAdmin) String() string {
	return "alter " + strings.Replace(r.Schema.String(), " ", " create ", 1)
}

func (r *alterCreateAdmin) execute(db *db19.Database) {
	checkForSystemTable("alter", r.Table)
	if !db.AlterCreate(&r.Schema) {
		panic("can't " + r.String())
	}
}

//-------------------------------------------------------------------

type alterRenameAdmin struct {
	table string
	from  []string
	to    []string
}

func (r *alterRenameAdmin) String() string {
	s := "alter " + r.table + " rename "
	sep := ""
	for i, from := range r.from {
		s += sep + from + " to " + r.to[i]
		sep = ", "
	}
	return s
}

func (r *alterRenameAdmin) execute(db *db19.Database) {
	checkForSystemTable("alter", r.table)
	if !db.AlterRename(r.table, r.from, r.to) {
		panic("can't " + r.String())
	}
}

//-------------------------------------------------------------------

type alterDropAdmin struct {
	Schema
}

func (r *alterDropAdmin) String() string {
	return "alter " + strings.Replace(r.Schema.String(), " ", " drop ", 1)
}

func (r *alterDropAdmin) execute(db *db19.Database) {
	checkForSystemTable("alter", r.Table)
	if !db.AlterDrop(&r.Schema) {
		panic("can't " + r.String())
	}
}

//-------------------------------------------------------------------

type dropAdmin struct {
	table string
}

func (r *dropAdmin) String() string {
	return "drop " + r.table
}

func (r *dropAdmin) execute(db *db19.Database) {
	checkForSystemTable("drop", r.table)
	if !db.DropTable(r.table) {
		panic("can't drop nonexistent table: " + r.table)
	}
}
