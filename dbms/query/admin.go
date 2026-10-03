// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
)

func DoAdmin(db *db19.Database, cmd string, sv *Sviews, perms *Perms) {
	admin := ParseAdmin(cmd)
	admin.execute(db, sv, perms)
}

func DoAdminTest(db *db19.Database, cmd string) {
	DoAdmin(db, cmd, nil, AllPerms)
}

// checkSchemaPerm panics unless the permissions allow the schema change.
// PermSchema only allows creating tables and adding columns
// (create, alter ... create, ensure), not other schema changes.
// Nil permissions deny all schema changes.
func checkSchemaPerm(perms *Perms, action SchemaPerm) {
	if !perms.SchemaActAllowed(action) {
		panic("not authorized")
	}
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

func (a *createAdmin) execute(db *db19.Database, _ *Sviews, perms *Perms) {
	checkForSystemTable(a.Table)
	checkSchemaPerm(perms, PermCreate)
	db.Create(&a.Schema)
	perms.AddTable(a.Table, "write")
}

//-------------------------------------------------------------------

type ensureAdmin struct {
	Schema
}

func (a *ensureAdmin) String() string {
	return "ensure " + a.Schema.String()
}

func (a *ensureAdmin) execute(db *db19.Database, _ *Sviews, perms *Perms) {
	checkForSystemTable(a.Table)
	checkSchemaPerm(perms, PermCreate)
	ts := db.GetState().Meta.GetRoSchema(a.Table)
	if ts != nil && !perms.SchemaActAllowed(PermUpdate) {
		// with only create permission, may restate existing indexes
		// and keys, but not add new ones
		for _, ix := range a.Indexes {
			if ts.FindIndex(ix.Columns) == nil {
				panic("not authorized: " + a.Table)
			}
		}
	}
	db.Ensure(&a.Schema)
	if ts == nil && perms != AllPerms {
		perms.AddTable(a.Table, "write")
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

func (a *renameAdmin) execute(db *db19.Database, _ *Sviews, perms *Perms) {
	checkSchemaPerm(perms, PermUpdate)
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

func (a *alterCreateAdmin) execute(db *db19.Database, _ *Sviews, perms *Perms) {
	if len(a.Indexes) == 0 {
		checkSchemaPerm(perms, PermCreate)
	} else {
		checkSchemaPerm(perms, PermUpdate)
	}
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
	var sb strings.Builder
	sb.WriteString("alter ")
	sb.WriteString(a.table)
	sb.WriteString(" rename ")
	sep := ""
	for i, from := range a.from {
		sb.WriteString(sep)
		sb.WriteString(from)
		sb.WriteString(" to ")
		sb.WriteString(a.to[i])
		sep = ", "
	}
	return sb.String()
}

func (a *alterRenameAdmin) execute(db *db19.Database, _ *Sviews, perms *Perms) {
	checkSchemaPerm(perms, PermUpdate)
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

func (a *alterDropAdmin) execute(db *db19.Database, _ *Sviews, perms *Perms) {
	checkSchemaPerm(perms, PermUpdate)
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

func (a *viewAdmin) execute(db *db19.Database, _ *Sviews, perms *Perms) {
	checkSchemaPerm(perms, PermUpdate)
	checkForSystemTable(a.name)
	db.AddView(a.name, a.def)
}

//-------------------------------------------------------------------

type sviewAdmin viewAdmin

func (a *sviewAdmin) String() string {
	return "sview " + a.name + " = " + a.def
}

func (a *sviewAdmin) execute(_ *db19.Database, sv *Sviews, perms *Perms) {
	checkSchemaPerm(perms, PermUpdate)
	checkForSystemTable(a.name)
	sv.AddSview(a.name, a.def)
}

//-------------------------------------------------------------------

type dropAdmin struct {
	table string
}

func (a *dropAdmin) String() string {
	return "drop " + a.table
}

func (a *dropAdmin) execute(db *db19.Database, sv *Sviews, perms *Perms) {
	checkSchemaPerm(perms, PermUpdate)
	checkForSystemTable(a.table)
	if sv != nil && sv.DropSview(a.table) {
		return
	}
	if err := db.Drop(a.table); err != nil {
		panic(err)
	}
}
