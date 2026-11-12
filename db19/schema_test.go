// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTablesNames(t *testing.T) {
	assert := assert.T(t)
	m := &meta.Meta{}
	add := func(name string) {
		m = m.Put(&meta.Schema{Schema: schema.Schema{
			Table: name, Columns: []string{"a"},
			Indexes: []schema.Index{{Mode: 'k', Columns: []string{"a"}}},
		}}, nil)
	}
	add("z")
	add("a")
	// Table names require only Schema, not Info.
	rt := &ReadTran{meta: m}
	assert.This(tablesSchema.Columns).Is([]string{"table"})
	check := func(rec core.Record, name string) {
		assert.This(rec.GetStr(0)).Is(name)
		assert.This(rec.Count()).Is(1)
	}
	st := tablesTable{}
	assert.This(st.nrows(&rt.tran)).Is(2)
	check(st.lookup(rt, 0, core.Pack(core.SuStr("a"))), "a")
	assert.This(st.lookup(rt, 0, core.Pack(core.SuStr("missing")))).Is(core.Record(""))
	it := newTablesIter(rt)
	it.Next(rt)
	check(it.GetRecord(rt), "a")
	it.Next(rt)
	check(it.GetRecord(rt), "z")
	it.Next(rt)
	assert.That(it.Eof())
}
