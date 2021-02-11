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

func (tbl *Table) Columns() []string { //TODO
	switch tbl.name {
	case "customer":
		return []string{"id", "name", "city"}
	case "hist":
		return []string{"date", "item", "id", "cost"}
	case "hist2":
		return []string{"date", "item", "id", "cost"}
	case "trans":
		return []string{"item", "id", "cost", "date"}
	case "inven":
		return []string{"item", "qty"}
	case "table":
		return []string{"a", "b", "c"}
	case "table2":
		return []string{"c", "d", "e"}
	case "tables":
		return []string{"table", "tablename"}
	case "columns":
		return []string{"table", "column"}
	}
	panic("unknown table: " + tbl.name)
}

func (tbl *Table) Keys() [][]string { //TODO
	switch tbl.name {
	case "customer":
		return [][]string{{"id"}}
	case "hist":
		return [][]string{{"date", "item", "id"}}
	case "hist2":
		return [][]string{{"date"}}
	case "trans":
		return [][]string{{"date", "item", "id"}}
	case "inven":
		return [][]string{{"item"}}
	case "table":
		return [][]string{{"a"}}
	case "table2":
		return [][]string{{"e"}}
	case "tables":
		return [][]string{{"table"}, {"tablename"}}
	case "columns":
		return [][]string{{"table", "column"}}
	}
	panic("unknown table: " + tbl.name)
}

func (tbl *Table) Transform() Query {
	return tbl
}

func (tbl *Table) Fixed() []Fixed {
	return nil
}
