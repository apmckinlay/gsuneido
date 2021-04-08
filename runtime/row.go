// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/sset"
)

// Row is the result of database queries.
// It is a list of DbRec so that operations e.g. join
// can avoid building new records.
type Row []DbRec

func (row Row) Get(hdr *Header, fld string) Value {
	assert.That(hdr.Map != nil)
	at, ok := hdr.Map[fld]
	if !ok || int(at.Reci) >= len(row) {
		return nil
	}
	return row[at.Reci].GetVal(int(at.Fldi))
}

func (row Row) GetRaw(hdr *Header, fld string) string {
	assert.That(hdr.Map != nil)
	at, ok := hdr.Map[fld]
	if !ok || int(at.Reci) >= len(row) {
		return ""
	}
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

//-------------------------------------------------------------------

// Header specifies the fields (physical) and columns (logical) for a query
type Header struct {
	Fields  [][]string
	Columns []string
	Map     map[string]RowAt
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
