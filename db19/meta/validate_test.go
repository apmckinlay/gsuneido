// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestValidateForeignKeys(t *testing.T) {
	assert := assert.T(t)

	// Create a base meta with a target table
	m := &Meta{}
	m.schema.Hamt = SchemaHamt{}.Mutable().Freeze()
	m.info.Hamt = InfoHamt{}.Mutable().Freeze()

	// Add target table with a key
	target := &Schema{Schema: schema.Schema{
		Table:   "target",
		Columns: []string{"id", "name"},
		Indexes: []schema.Index{
			{Mode: 'k', Columns: []string{"id"}},
		},
	}}
	target.Ixspecs(len(target.Indexes))
	targetInfo := NewInfo("target", nil, 0, 0)
	m = m.Put(target, targetInfo)

	// Test 1: Foreign key to nonexistent table
	assert.This(func() {
		mu := newMetaUpdate(m)
		refTable := &Schema{Schema: schema.Schema{
			Table:   "referring",
			Columns: []string{"id", "target_id"},
			Indexes: []schema.Index{
				{Mode: 'k', Columns: []string{"id"}},
				{Mode: 'i', Columns: []string{"target_id"},
					Fk: schema.Fkey{Table: "nonexistent", IIndex: 0}},
			},
		}}
		refTable.Ixspecs(len(refTable.Indexes))
		mu.putSchema(refTable)
		mu.freeze()
	}).Panics("foreign key references nonexistent table")

	// Test 2: Foreign key to nonexistent index
	assert.This(func() {
		mu := newMetaUpdate(m)
		refTable := &Schema{Schema: schema.Schema{
			Table:   "referring",
			Columns: []string{"id", "target_id"},
			Indexes: []schema.Index{
				{Mode: 'k', Columns: []string{"id"}},
				{Mode: 'i', Columns: []string{"target_id"},
					Fk: schema.Fkey{Table: "target", Columns: []string{"nonexistent"}, IIndex: 0}},
			},
		}}
		refTable.Ixspecs(len(refTable.Indexes))
		mu.putSchema(refTable)
		mu.freeze()
	}).Panics("foreign key references nonexistent index")

	// Test 3: Foreign key to non-key index
	assert.This(func() {
		// First add an index to target that is not a key
		mu := newMetaUpdate(m)
		targetWithIndex := *m.GetRoSchema("target") // copy
		targetWithIndex.Indexes = append(targetWithIndex.Indexes,
			schema.Index{Mode: 'i', Columns: []string{"name"}})
		targetWithIndex.Ixspecs(len(targetWithIndex.Indexes))
		mu.putSchema(&targetWithIndex)
		m2 := mu.freeze()

		// Now try to create FK to the non-key index
		mu2 := newMetaUpdate(m2)
		refTable := &Schema{Schema: schema.Schema{
			Table:   "referring",
			Columns: []string{"id", "target_name"},
			Indexes: []schema.Index{
				{Mode: 'k', Columns: []string{"id"}},
				{Mode: 'i', Columns: []string{"target_name"},
					Fk: schema.Fkey{Table: "target", Columns: []string{"name"}, IIndex: 1}},
			},
		}}
		refTable.Ixspecs(len(refTable.Indexes))
		mu2.putSchema(refTable)
		mu2.freeze()
	}).Panics("foreign key must point to key")

	// Test 4: Foreign key with wrong IIndex
	// First add another key to target so we have multiple indexes to reference
	mu4 := newMetaUpdate(m)
	targetWith2Keys := *m.GetRoSchema("target") // copy
	targetWith2Keys.Indexes = append(targetWith2Keys.Indexes,
		schema.Index{Mode: 'k', Columns: []string{"name"}})
	targetWith2Keys.Ixspecs(len(targetWith2Keys.Indexes))
	mu4.putSchema(&targetWith2Keys)
	m4 := mu4.freeze()

	assert.This(func() {
		mu := newMetaUpdate(m4)
		refTable := &Schema{Schema: schema.Schema{
			Table:   "referring",
			Columns: []string{"id", "target_id"},
			Indexes: []schema.Index{
				{Mode: 'k', Columns: []string{"id"}},
				{Mode: 'i', Columns: []string{"target_id"},
					Fk: schema.Fkey{Table: "target", Columns: []string{"id"}, IIndex: 1}}, // wrong IIndex - should be 0
			},
		}}
		refTable.Ixspecs(len(refTable.Indexes))
		mu.putSchema(refTable)
		mu.freeze()
	}).Panics("foreign key IIndex mismatch")

	// Test 5: Valid foreign key
	mu := newMetaUpdate(m)
	refTable := &Schema{Schema: schema.Schema{
		Table:   "referring",
		Columns: []string{"id", "target_id"},
		Indexes: []schema.Index{
			{Mode: 'k', Columns: []string{"id"}},
			{Mode: 'i', Columns: []string{"target_id"},
				Fk: schema.Fkey{Table: "target", Columns: []string{"id"}, IIndex: 0}},
		},
	}}
	refTable.Ixspecs(len(refTable.Indexes))
	mu.putSchema(refTable)
	mu.freeze() // Should not panic
}

func TestValidateCreateDuplicateColumn(t *testing.T) {
	assert := assert.T(t)
	store := stor.HeapStor(8192)

	m := &Meta{}
	m.schema.Hamt = SchemaHamt{}.Mutable().Freeze()
	m.info.Hamt = InfoHamt{}.Mutable().Freeze()

	// Add a base table
	base := &Schema{Schema: schema.Schema{
		Table:   "test",
		Columns: []string{"a", "b"},
		Indexes: []schema.Index{
			{Mode: 'k', Columns: []string{"a"}},
		},
	}}
	base.Ixspecs(len(base.Indexes))
	baseInfo := NewInfo("test", nil, 0, 0)
	m = m.Put(base, baseInfo)

	// Try to create a duplicate column via AlterCreate
	assert.This(func() {
		ac := &schema.Schema{
			Table:   "test",
			Columns: []string{"a"}, // duplicate of existing column
		}
		m.AlterCreate(ac, store)
	}).Panics("can't create existing column")
}

func TestValidateForeignKeyAfterRename(t *testing.T) {
	m := &Meta{}
	m.schema.Hamt = SchemaHamt{}.Mutable().Freeze()
	m.info.Hamt = InfoHamt{}.Mutable().Freeze()

	// Create target table
	target := &Schema{Schema: schema.Schema{
		Table:   "target",
		Columns: []string{"id", "name"},
		Indexes: []schema.Index{
			{Mode: 'k', Columns: []string{"id"}},
		},
	}}
	target.Ixspecs(len(target.Indexes))
	targetInfo := NewInfo("target", nil, 0, 0)
	m = m.Put(target, targetInfo)

	// Create referring table with explicit FK columns
	mu := newMetaUpdate(m)
	referring := &Schema{Schema: schema.Schema{
		Table:   "referring",
		Columns: []string{"rid", "target_id"},
		Indexes: []schema.Index{
			{Mode: 'k', Columns: []string{"rid"}},
			{Mode: 'i', Columns: []string{"target_id"},
				Fk: schema.Fkey{Table: "target", Columns: []string{"id"}, IIndex: 0}},
		},
	}}
	referring.Ixspecs(len(referring.Indexes))
	mu.putSchema(referring)
	mu.putInfo(NewInfo("referring", nil, 0, 0))
	m.createFkeys(mu, &referring.Schema, &referring.Schema)
	m = mu.freeze()

	// Rename column in target table - should update FK references
	m2 := m.AlterRename("target", []string{"id"}, []string{"new_id"})

	// Verify the FK was updated and still validates
	ref2 := m2.GetRoSchema("referring")
	checkForeignKeys(m2.GetRoSchema, &ref2.Schema)
	// Should not panic - FK.Columns should have been updated from "id" to "new_id"
}
