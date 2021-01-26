// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Table struct {
	name string
}

func (tbl *Table) String() string {
	return tbl.name
}

func (tbl *Table) Init() {
}

func (tbl *Table) Columns() []string {
	return []string{} //TODO
}

func (tbl *Table) Transform() Query {
	return tbl
}
