// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/db19/meta"

type testTran struct{}

var testSchemas = map[string]*Schema{
	"tables": {Columns: []string{"table", "tablename"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"table"}},
			{Mode: 'k', Columns: []string{"tablename"}}}},
	"columns": {Columns: []string{"table", "column"},
		Indexes: []Index{{Mode: 'k', Columns: []string{"table", "column"}}}},
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
}

func (testTran) GetSchema(table string) *Schema {
	return testSchemas[table]
}

var testInfo = map[string]*meta.Info{
	"alias": {Nrows: 10, Size: 1000},
	"task":  {Nrows: 200, Size: 20000},
	"trans": {Nrows: 1000, Size: 100000},
	"hist2": {Nrows: 1000, Size: 100000},
}

func (testTran) GetInfo(table string) *meta.Info {
	if ti, ok := testInfo[table]; ok {
		return ti
	}
	return &meta.Info{Nrows: 100, Size: 10000}
}
