// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Table struct {
	name string
	t    QueryTran
}

func (tbl *Table) String() string {
	return tbl.name
}

func (tbl *Table) Init() {
}

func (tbl *Table) SetTran(t QueryTran) {
	tbl.t = t
}

func (tbl *Table) Columns() []string {
	schema := tbl.t.GetSchema(tbl.name)
	if schema == nil {
		panic("nonexistent table: " + tbl.name)
	}
	allcols := make([]string, 0, len(schema.Columns)+len(schema.Derived))
	allcols = append(allcols, schema.Columns...)
	allcols = append(allcols, schema.Derived...)
	return allcols
}

func (tbl *Table) Keys() [][]string {
	schema := tbl.t.GetSchema(tbl.name)
	if schema == nil {
		panic("nonexistent table: " + tbl.name)
	}
	keys := make([][]string, 0, 1)
	for _, ix := range schema.Indexes {
		if ix.Mode == 'k' {
			keys = append(keys, ix.Columns)
		}
	}
	return keys
}

func (tbl *Table) Transform() Query {
	return tbl
}

func (tbl *Table) Fixed() []Fixed {
	return nil
}
