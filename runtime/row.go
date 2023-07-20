// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"math"
	"strings"

	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/set"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
	"golang.org/x/exp/slices"
)

// Row is the result of database queries.
// It is a list of DbRec so that operations e.g. join
// can avoid building new records.
type Row []DbRec

// DbRec is a Record along with its offset
type DbRec struct {
	Record
	Off uint64
}

func JoinRows(row1, row2 Row) Row {
	result := make(Row, 0, len(row1)+len(row2))
	return append(append(result, row1...), row2...)
}

// GetVal is used by query summarize and expr.
//   - returns "" for fld not in hdr.Columns
//   - returns stored value for Fields (rule ignored)
//   - calls rule for Columns not in Fields
func (row Row) GetVal(hdr *Header, fld string, th *Thread, tran *SuTran) Value {
	if !slices.Contains(hdr.Columns, fld) {
		return EmptyStr
	}
	if raw, ok := row.getRaw2(hdr, fld); ok {
		return Unpack(raw)
	}
	if strings.HasSuffix(fld, "_lower!") {
		base := fld[:len(fld)-7]
		x, _ := row.getRaw2(hdr, base)
		val := Unpack(x)
		return SuStr(str.ToLower(ToStr(val)))
	}
	if nil == getRule(th, fld) {
		return EmptyStr
	}
	// else construct SuRecord to handle rules
	return SuRecordFromRow(row, hdr, "", tran).Get(th, SuStr(fld))
}

// GetRaw handles _lower! but does NOT handle rules.
// It is used by SuRecord Get.
func (row Row) GetRaw(hdr *Header, fld string) string {
	if strings.HasSuffix(fld, "_lower!") {
		base := fld[:len(fld)-7]
		x, _ := row.getRaw2(hdr, base)
		return lowerRaw(x)
	}
	x, _ := row.getRaw2(hdr, fld)
	return x
}

// GetRawVal is like GetVal (i.e. handles rules) but returns a raw/packed value.
func (row Row) GetRawVal(hdr *Header, fld string, th *Thread, tran *SuTran) string {
	if strings.HasSuffix(fld, "_lower!") {
		base := fld[:len(fld)-7]
		x, _ := row.getRaw2(hdr, base)
		return lowerRaw(x)
	}
	if raw, ok := row.getRaw2(hdr, fld); ok {
		return raw
	}
	if nil == getRule(th, fld) {
		return ""
	}
	// else construct SuRecord to handle rules
	v := SuRecordFromRow(row, hdr, "", tran).Get(th, SuStr(fld))
	return Pack(v.(Packable))
}

func lowerRaw(x string) string {
	if x == "" || x[0] != PackString {
		return x
	}
	hasUpper := false
	for i := 0; i < len(x); i++ {
		hasUpper = hasUpper || ascii.IsUpper(x[i])
	}
	if !hasUpper {
		return x
	}
	buf := make([]byte, len(x))
	buf[0] = x[0]
	for i := 1; i < len(buf); i++ {
		buf[i] = ascii.ToLower(x[i])
	}
	return hacks.BStoS(buf)
}

func getRule(th *Thread, key string) Value {
	if th == nil || !validRule(key) {
		return nil
	}
	return Global.FindName(th, "Rule_"+key)
}

// getRaw2 gets a stored field.
// It handles a field having multiple possible locations
// due to union supplying one of two records.
// It returns false if fld is not in hdr.Fields i.e. derived (non-stored) fields.
func (row Row) getRaw2(hdr *Header, fld string) (string, bool) {
	// find will only return the first possible location
	at, ok := hdr.find(fld)
	if ok && row[at.Reci].Record != "" { // not empty side of union
		return row[at.Reci].GetRaw(int(at.Fldi)), true
	}
	// handle nil records from Union
	for reci := int(at.Reci + 1); reci < len(hdr.Fields); reci++ {
		if fldi := slices.Index(hdr.Fields[reci], fld); fldi >= 0 {
			if row[reci].Record != "" {
				return row[reci].GetRaw(int(fldi)), true
			}
		}
	}
	return "", false
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

func (row Row) Derived() int {
	n := 0
	for _, dbrec := range row {
		if dbrec.Off == 0 { // derived record e.g. from extend or summarize
			n += len(dbrec.Record)
		}
	}
	return n
}

//-------------------------------------------------------------------

// Header specifies the fields (physical) and columns (logical) for a query
type Header struct {
	// cache the location of fields.
	// WARNING: assumed to not be concurrent (no locking)
	cache map[string]rowAt
	// Fields is the physical fields, including deleted ("-")
	Fields [][]string
	// Columns is the logical fields
	Columns []string
}

// rowAt specifies the position of a field within a Row
type rowAt struct {
	Reci int16
	Fldi int16
}

func NewHeader(fields [][]string, columns []string) *Header {
	return &Header{Fields: fields, Columns: columns,
		cache: make(map[string]rowAt)}
}

func SimpleHeader(fields []string) *Header {
	return NewHeader([][]string{fields}, fields)
}

func (hdr *Header) Dup() *Header {
	if hdr == nil {
		return nil
	}
	return NewHeader(hdr.Fields, hdr.Columns)
}

func JoinHeaders(x, y *Header) *Header {
	//TODO what if rule on one and field on the other ???
	fields := make([][]string, 0, len(x.Fields)+len(y.Fields))
	fields = append(append(fields, x.Fields...), y.Fields...)
	columns := set.Union(x.Columns, y.Columns)
	return NewHeader(fields, columns)
}

// Rules is a list of the rule columns i.e. columns that are not fields
func (hdr *Header) Rules() []string {
	rules := []string{}
	for _, col := range hdr.Columns {
		if !hdr.HasField(col) {
			rules = append(rules, col)
		}
	}
	return rules
}

func (hdr *Header) HasField(col string) bool {
	for _, fields := range hdr.Fields {
		if slices.Contains(fields, col) {
			return true
		}
	}
	return false
}

// find returns the location of the first occurence of fld in hdr.Fields
// and caches the result. (Multiple occurrences come from union.)
func (hdr *Header) find(fld string) (rowAt, bool) {
	if at, ok := hdr.cache[fld]; ok {
		return at, at != missing
	}
	for reci, fields := range hdr.Fields {
		if fldi := slices.Index(fields, fld); fldi >= 0 {
			at := rowAt{Reci: int16(reci), Fldi: int16(fldi)}
			hdr.cache[fld] = at
			return at, true
		}
	}
	hdr.cache[fld] = missing
	return missing, false
}

var missing = rowAt{Reci: 0, Fldi: math.MaxInt16}

func (hdr *Header) HasDeleted() bool {
	_, ok := hdr.find("-") // will cache if found
	return ok
}

// GetFields returns a list of the fields, including deleted ("-"),
// without duplicates (e.g. from Join or Union)
func (hdr *Header) GetFields() []string {
	if len(hdr.Fields) == 1 {
		return slices.Clip(hdr.Fields[0])
	}
	result := make([]string, 0, len(hdr.Columns))
	result = append(result, hdr.Fields[0]...)
	for _, fields := range hdr.Fields[1:] {
		for _, fld := range fields {
			if fld == "-" || !slices.Contains(result, fld) {
				result = append(result, fld)
			}
		}
	}
	return result
}

// EqualRows is used by Project
func (hdr *Header) EqualRows(r1, r2 Row, th *Thread, st *SuTran) bool {
	return EqualRows(hdr, r1, hdr, r2, hdr.Columns, th, st)
}

// EqualRows is used by Compatible
func EqualRows(hdr1 *Header, r1 Row, hdr2 *Header, r2 Row, cols []string,
	th *Thread, st *SuTran) bool {
	for _, col := range cols {
		if r1.GetRawVal(hdr1, col, th, st) != r2.GetRawVal(hdr2, col, th, st) {
			return false
		}
	}
	return true
}

// func (hdr *Header) Equal(hdr2 *Header) bool {
// 	if !slices.Equal(hdr.Columns, hdr2.Columns) {
// 		return false
// 	}
// 	if len(hdr.Fields) != len(hdr2.Fields) {
// 		return false
// 	}
// 	for i := range hdr.Fields {
// 		if !slices.Equal(hdr.Fields[i], hdr2.Fields[i]) {
// 			return false
// 		}
// 	}
// 	return true
// }

// Schema is GetFields() plus the derived columns with the rules capitalized
func (hdr *Header) Schema() []string {
	fields := hdr.GetFields()
	return hdr.AppendDerived(fields)
}

func (hdr *Header) AppendDerived(fields []string) []string {
	for _, col := range hdr.Columns {
		if !slices.Contains(fields, col) {
			if !strings.HasSuffix(col, "_lower!") {
				col = str.Capitalize(col)
			}
			fields = append(fields, col)
		}
	}
	assert.That(!slices.Contains(fields, ""))
	return fields
}

// Physical is the fields without deleted ("-")
func (hdr *Header) Physical() []string {
	if len(hdr.Fields) == 1 && !slices.Contains(hdr.Fields[0], "-") {
		return hdr.Fields[0]
	}
	result := make([]string, 0, len(hdr.Columns))
	for _, flds := range hdr.Fields {
		for _, fld := range flds {
			if fld != "-" {
				result = append(result, fld)
			}
		}
	}
	return result
}
