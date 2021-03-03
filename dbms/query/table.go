// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Table struct {
	name   string
	t      QueryTran
	schema *Schema // cached
}

func (tbl *Table) String() string {
	return tbl.name
}

func (tbl *Table) Init() {
}

func (tbl *Table) SetTran(t QueryTran) {
	tbl.t = t
	tbl.schema = t.GetSchema(tbl.name)
	if tbl.schema == nil {
		panic("nonexistent table: " + tbl.name)
	}
}

func (tbl *Table) Columns() []string {
	allcols := make([]string, 0, len(tbl.schema.Columns)+len(tbl.schema.Derived))
	allcols = append(allcols, tbl.schema.Columns...)
	allcols = append(allcols, tbl.schema.Derived...)
	return allcols
}

func (tbl *Table) Indexes() [][]string {
	idxs := make([][]string, 0, 1)
	for _, ix := range tbl.schema.Indexes {
		idxs = append(idxs, ix.Columns)
	}
	return idxs
}

func (tbl *Table) Keys() [][]string {
	keys := make([][]string, 0, 1)
	for _, ix := range tbl.schema.Indexes {
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

func (tbl *Table) Updateable() bool {
	return true
}
