// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// TestColumns covers the edge cases that fuzzColumns' random tables cannot
// produce: a hidden "-" column and a derived column that sorts between stored
// columns (NewFuzzTable's derived columns are named Cr0.. and so always sort
// last).
func TestColumns(t *testing.T) {
	db := heapDb()
	defer db.Close()
	db.adm("create withhidden (a, -, b) key(a)")
	db.adm("create withderived (a, B, c) key(a)")

	// hidden "-" excluded; b's field number still reflects its physical
	// position after the "-" placeholder
	assert.T(t).This(queryAll(db.Database, "columns where table = 'withhidden'")).
		Is("column=a field=0 table=withhidden | column=b field=2 table=withhidden")
	// derived reported uncapitalized (b) with field -1, sorted between a and c
	assert.T(t).This(queryAll(db.Database, "columns where table = 'withderived'")).
		Is("column=a field=0 table=withderived | column=b field=-1 table=withderived |" +
			" column=c field=1 table=withderived")

	// lookups by (table, column)
	assert.T(t).This(queryAll(db.Database, "columns where table = 'withderived' and column = 'b'")).
		Is("column=b field=-1 table=withderived")
	assert.T(t).This(queryAll(db.Database, "columns where table = 'withhidden' and column = '-'")).
		Is("")
}

// TestSchemaJoinNonStringKey verifies that joining a real table to a schema
// table on a non-string key yields an unmatched row rather than panicking in
// the schema table's Lookup (which must not assume string keys).
func TestSchemaJoinNonStringKey(t *testing.T) {
	db := heapDb()
	defer db.Close()

	db.adm("create badtable (table, x) key(table)")
	db.act("insert { table: 123, x: 1 } into badtable")
	db.adm("create badcolumn (table, column, y) key(table, column)")
	db.act("insert { table: 123, column: 456, y: 1 } into badcolumn")
	db.adm("create badview (view_name, z) key(view_name)")
	db.act("insert { view_name: 123, z: 1 } into badview")

	assert.T(t).This(queryAll(db.Database, "badtable leftjoin tables")).
		Is("table=123 x=1")
	assert.T(t).This(queryAll(db.Database, "badcolumn leftjoin columns")).
		Is("column=456 table=123 y=1")
	assert.T(t).This(queryAll(db.Database, "badview leftjoin views")).
		Is("view_name=123 z=1")
}
