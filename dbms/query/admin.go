// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
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
	schema Schema
}

func (r *createAdmin) String() string {
	return "create " + r.schema.String()
}

func (r *createAdmin) execute(db *db19.Database) {
	checkForSystemTable("create", r.schema.Table)
	ts := &meta.Schema{Schema: r.schema}
	ts.Ixspecs()
	ov := make([]*index.Overlay, len(ts.Indexes))
	for i := range ov {
		bt := btree.CreateBtree(db.Store, &ts.Indexes[i].Ixspec)
		ov[i] = index.OverlayFor(bt)
	}
	ti := &meta.Info{Table: r.schema.Table, Indexes: ov}
	db.LoadedTable(ts, ti)
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
		panic("can't rename: " + r.from + " to " + r.to)
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
	//TODO
}

//-------------------------------------------------------------------

type alterCreateAdmin struct {
	schema Schema
}

func (r *alterCreateAdmin) String() string {
	return "alter " + strings.Replace(r.schema.String(), " ", " create ", 1)
}

func (r *alterCreateAdmin) execute(db *db19.Database) {
	//TODO
}

//-------------------------------------------------------------------

type alterDropAdmin struct {
	schema Schema
}

func (r *alterDropAdmin) String() string {
	return "alter " + strings.Replace(r.schema.String(), " ", " drop ", 1)
}

func (r *alterDropAdmin) execute(db *db19.Database) {
	//TODO
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
