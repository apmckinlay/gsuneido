// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/runtime"
)

// testTran has hard coded table schemas for tests
// See also: sizeTran
type testTran struct{}

var testSchemas = map[string]*Schema{
	"table": {Columns: []string{"a", "b", "c"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}}},
	"table2": {Columns: []string{"c", "d", "e"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"e"}}}},
	"customer": {Columns: []string{"id", "name", "city"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"id"}}}},
	"supplier": {Columns: []string{"supplier", "name", "city"}, Indexes: []Index{
		{Mode: 'k', Columns: []string{"supplier"}},
		{Mode: 'i', Columns: []string{"city"}}}},
	"hist": {Columns: []string{"date", "item", "id", "cost"}, Indexes: []Index{
		{Mode: 'i', Columns: []string{"date"}},
		{Mode: 'k', Columns: []string{"date", "item", "id"}}}},
	"hist2": {Columns: []string{"date", "item", "id", "cost"}, Indexes: []Index{
		{Mode: 'i', Columns: []string{"id"}},
		{Mode: 'k', Columns: []string{"date"}}}},
	"trans": {Columns: []string{"item", "id", "cost", "date"}, Indexes: []Index{
		{Mode: 'k', Columns: []string{"date", "item", "id"}},
		{Mode: 'i', Columns: []string{"item"}}}},
	"inven": {Columns: []string{"item", "qty"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"item"}}}},
	"abc": {Columns: []string{"a", "b", "c"}, Indexes: []Index{
		{Mode: 'i', Columns: []string{"a"}},
		{Mode: 'k', Columns: []string{"b"}},
		{Mode: 'k', Columns: []string{"c"}}}},
	"bcd": {Columns: []string{"b", "c", "d"}, Indexes: []Index{
		{Mode: 'k', Columns: []string{"b"}},
		{Mode: 'k', Columns: []string{"c"}},
		{Mode: 'i', Columns: []string{"d"}}}},
	"cus": {Columns: []string{"cnum", "abbrev", "name"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"cnum"}}}},
	"task": {Columns: []string{"tnum", "cnum"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"tnum"}}}},
	"co": {Columns: []string{"tnum", "signed"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"tnum"}}}},
	"alias": {Columns: []string{"id", "name2"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"id"}}}},
	"comp": {Columns: []string{"a", "b", "c"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"a", "b", "c"}}}},
	"withdeps": {Columns: []string{"a", "b", "b_deps", "c", "c_deps"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"a"}}}},
}

func (testTran) GetSchema(table string) *Schema {
	return testSchemas[table]
}

var testInfo = map[string]*meta.Info{
	"alias":   {Nrows: 10, Size: 1000},
	"task":    {Nrows: 200, Size: 20000},
	"columns": {Nrows: 1000, Size: 100000},
	"trans":   {Nrows: 1000, Size: 100000},
	"hist2":   {Nrows: 1000, Size: 100000},
	"comp":    {Nrows: 1000, Size: 100000},
}

func (testTran) GetInfo(table string) *meta.Info {
	if ti, ok := testInfo[table]; ok {
		return ti
	}
	return &meta.Info{Nrows: 100, Size: 10000}
}

func (testTran) GetAllInfo() []*meta.Info {
	infos := make([]*meta.Info, 0, len(testSchemas))
	for table := range testSchemas {
		info := meta.Info{Table: table, Nrows: 100, Size: 10000}
		if ti, ok := testInfo[table]; ok {
			info.Nrows, info.Size = ti.Nrows, ti.Size
		}
		infos = append(infos, &info)
	}
	return infos
}

func (testTran) GetAllSchema() []*meta.Schema {
	schemas := make([]*meta.Schema, 0, len(testSchemas))
	for table, schema := range testSchemas {
		schema.Table = table
		schemas = append(schemas,
			&meta.Schema{Schema: *schema})
	}
	return schemas
}

func (testTran) GetAllViews() []string {
	return nil
}

func (t testTran) GetView(table string) string {
	if table == "myview" {
		return "cus join task"
	}
	return ""
}

func (t testTran) GetStore() *stor.Stor {
	return nil
}

func (t testTran) RangeFrac(table string, iIndex int, org, end string) float64 {
	schema := t.GetSchema(table)
	ix := schema.Indexes[iIndex]
	decode := len(ix.Columns) > 1 || ix.Mode != 'k'
	orgPos := t.fracPos(org, decode)
	endPos := t.fracPos(end, decode)
	return endPos - orgPos
}

// fracPos treats keys as multi-digit decimal numbers.
// Each component of the key should be an integer from 0 to 9.
func (t testTran) fracPos(key string, decode bool) float64 {
	var vals []string
	if decode {
		vals = ixkey.Decode(key)
	} else {
		vals = []string{key}
	}
	var f float64
	mul := float64(.1)
	for i, s := range vals {
		var n int
		if s == ixkey.Max {
			n = 10
		} else {
			n = runtime.ToInt(runtime.Unpack(s))
			if i+1 < len(vals) && vals[i+1] == "" {
				n++
			}
		}
		f += mul * float64(n)
		mul /= 10
	}
	return f
}

func (t testTran) Lookup(_ string, _ int, key string) *runtime.DbRec {
	// WARNING: assumes key columns match table columns
	var vals []string
	if strings.Contains(key, "\x00\x00") {
		vals = ixkey.Decode(key)
	} else {
		vals = []string{key}
	}
	var rb runtime.RecordBuilder
	for _, v := range vals {
		rb.AddRaw(v)
	}
	return &runtime.DbRec{Record: rb.Build()}
}

func (t testTran) Output(*runtime.Thread, string, runtime.Record) {
	panic("should not be called")
}

func (t testTran) GetIndexI(string, int) *index.Overlay {
	panic("should not be called")
}

func (t testTran) GetRecord(uint64) runtime.Record {
	panic("should not be called")
}

func (t testTran) MakeLess(*ixkey.Spec) func(x, y uint64) bool {
	panic("should not be called")
}

func (t testTran) Read(string, int, string, string) {
}
