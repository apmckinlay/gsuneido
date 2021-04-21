// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sset"
)

// Row is the result of database queries.
// It is a list of DbRec so that operations e.g. join
// can avoid building new records.
type Row []DbRec

func JoinRows(row1, row2 Row) Row {
	result := make(Row, 0, len(row1)+len(row2))
	return append(append(result, row1...), row2...)
}

func (row Row) Get(hdr *Header, fld string) Value {
	return Unpack(row.GetRaw(hdr, fld))
}

func (row Row) GetRaw(hdr *Header, fld string) string {
	for {
		assert.That(hdr.Map != nil)
		at, ok := hdr.Map[fld]
		if !ok || int(at.Reci) >= len(row) {
			return ""
		}
		if row[at.Reci].Record != "" {
			return row[at.Reci].GetRaw(int(at.Fldi))
		}
		if hdr.Next == nil {
			return ""
		}
		row = row[len(hdr.Fields):]
		hdr = hdr.Next
	}
}

func (row Row) GetRawAt(at RowAt) string {
	return row[at.Reci].GetRaw(int(at.Fldi))
}

// RowAt specifies the position of a field within a Row
type RowAt struct {
	Reci int16
	Fldi int16
}

// DbRec is a Record along with its offset
type DbRec struct {
	Record
	Off uint64
}

// SameAs returns true if the db records have the same Off's
// and derived records (with Off == 0) are equal
func (row Row) SameAs(row2 Row) bool {
	if len(row) != len(row2) {
		return false
	}
	for i := range row {
		if (row[i].Off == 0) != (row2[i].Off == 0) {
			return false
		}
		if row[i].Off == 0 {
			if row[i].Record != row2[i].Record {
				return false
			}
		} else {
			if row[i].Off != row2[i].Off {
				return false
			}
		}
	}
	return true
}

//-------------------------------------------------------------------

// Header specifies the fields (physical) and columns (logical) for a query
type Header struct {
	Fields  [][]string
	Columns []string
	Map     map[string]RowAt
	// Next is used by Union
	Next *Header
}

func NewHeader(fields [][]string, columns []string) *Header {
	hdr := Header{Fields: fields, Columns: columns}
	hdr.Map = make(map[string]RowAt, len(hdr.Fields))
	for ri, r := range hdr.Fields {
		for fi, f := range r {
			hdr.Map[f] = RowAt{int16(ri), int16(fi)}
		}
	}
	return &hdr
}

func JoinHeaders(x, y *Header) *Header {
	fields := make([][]string, 0, len(x.Fields)+len(y.Fields))
	fields = append(append(fields, x.Fields...), y.Fields...)
	columns := sset.Union(x.Columns, y.Columns)
	return NewHeader(fields, columns)
}

func (hdr *Header) String() string {
	s := "Header\n"
	for hdr != nil {
		s += fmt.Sprintln("Fields:", hdr.Fields)
		s += fmt.Sprintln("Columns:", hdr.Columns)
		hdr = hdr.Next
	}
	return s
}

// Rules is a list of the rule columns i.e. columns that are not fields
func (hdr *Header) Rules() []string {
	rules := []string{}
	for _, col := range hdr.Columns {
		if !hdr.hasField(col) {
			rules = append(rules, col)
		}
	}
	return rules
}

func (hdr *Header) hasField(col string) bool {
	for _, fields := range hdr.Fields {
		if sset.Contains(fields, col) {
			return true
		}
	}
	return false
}

func (hdr *Header) GetFields() []string {
	if len(hdr.Fields) == 1 {
		return hdr.Fields[0]
	}
	result := make([]string, len(hdr.Fields[0]))
	copy(result, hdr.Fields[0])
	for _, fields := range hdr.Fields[1:] {
		for _, fld := range fields {
			result = sset.AddUnique(result, fld)
		}
	}
	return result
}

func (hdr *Header) EqualRows(r1, r2 Row) bool {
	for _, col := range hdr.Columns {
		if r1.GetRaw(hdr, col) != r2.GetRaw(hdr, col) {
			return false
		}
	}
	return true
}
