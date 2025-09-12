// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDatabaseDropTable(t *testing.T) {
	db := CreateDb(stor.HeapStor(16 * 1024))
	createTbl(db)
	db.CheckerSync()
	assert.T(t).That(db.Drop("nonexistent") != nil)
	assert.T(t).That(db.Drop("mytable") == nil)
	assert.T(t).That(db.Drop("mytable") != nil)
}

func TestForeignKeyToNonexistentIndex(t *testing.T) {
	db := CreateDb(stor.HeapStor(16 * 1024))
	defer db.Close()

	// Try to create a table with foreign key to a nonexistent table
	assert.T(t).This(func() {
		db.Create(&schema.Schema{
			Table:   "lin",
			Columns: []string{"c", "d"},
			Indexes: []schema.Index{
				{Mode: 'k', Columns: []string{"c"}},
				{Mode: 'i', Columns: []string{"d"}, Fk: schema.Fkey{Table: "hdr", Columns: []string{"b"}}},
			},
		})
	}).Panics("can't create foreign key to nonexistent table")

	db.Create(&schema.Schema{
		Table:   "hdr",
		Columns: []string{"a", "b"},
		Indexes: []schema.Index{{Mode: 'k', Columns: []string{"a"}}},
	})

	// Try to create a table with foreign key to a nonexistent index
	assert.T(t).This(func() {
		db.Create(&schema.Schema{
			Table:   "lin",
			Columns: []string{"c", "d"},
			Indexes: []schema.Index{
				{Mode: 'k', Columns: []string{"c"}},
				{Mode: 'i', Columns: []string{"d"}, Fk: schema.Fkey{Table: "hdr", Columns: []string{"b"}}},
			},
		})
	}).Panics("can't create foreign key to nonexistent index")
}
