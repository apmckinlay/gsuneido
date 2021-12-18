// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19"
)

func DoAdmin(db *db19.Database, cmd string) {
	admin := ParseAdmin(cmd)
	admin.execute(db)
}

func checkForSystemTable(table string) {
	if isSystemTable(table) {
		panic("can't modify system table: " + table)
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
	checkForSystemTable(a.Table)
	db.Create(&a.Schema)
}

//-------------------------------------------------------------------

type ensureAdmin struct {
	Schema
}

func (a *ensureAdmin) String() string {
	return "ensure " + a.Schema.String()
}

func (a *ensureAdmin) execute(db *db19.Database) {
	checkForSystemTable(a.Table)
	db.Ensure(&a.Schema)
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
	checkForSystemTable(a.from)
	checkForSystemTable(a.to)
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
	checkForSystemTable(a.Table)
	db.AlterCreate(&a.Schema)
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
	checkForSystemTable(a.table)
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
	checkForSystemTable(a.Table)
	if !db.AlterDrop(&a.Schema) {
		panic("can't " + a.String())
	}
}

//-------------------------------------------------------------------

type viewAdmin struct {
	name string
	def  string
}

func (a *viewAdmin) String() string {
	return "view " + a.name + " = " + a.def
}

func (a *viewAdmin) execute(db *db19.Database) {
	checkForSystemTable(a.name)
	db.AddView(a.name, a.def)
}

//-------------------------------------------------------------------

type dropAdmin struct {
	table string
}

func (a *dropAdmin) String() string {
	return "drop " + a.table
}

func (a *dropAdmin) execute(db *db19.Database) {
	checkForSystemTable(a.table)
	if err := db.Drop(a.table); err != nil {
		panic(err)
	}
}
