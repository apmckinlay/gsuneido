package runtime

import "github.com/apmckinlay/gsuneido/util/str"

type Row []DbRec

func (row Row) Get(hdr *Header, fld string) Value {
	at, ok := hdr.Map[fld]
	if !ok || int(at.Reci) >= len(row) {
		return nil
	}
	return row[at.Reci].GetVal(int(at.Fldi))
}

func (row Row) GetRaw(hdr *Header, fld string) string {
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

// DbRec is a Record along with its address
type DbRec struct {
	Record
	Adr int
}

// Header specifies the fields (physical) and columns (logical) for a query
type Header struct {
	Fields  [][]string
	Columns []string
	Map     map[string]RowAt
}

// Rules is a list of the rule columns i.e. columns that are not fields
func (hdr *Header) Rules() []string {
	rules := []string{}
	for _, col := range hdr.Columns {
		if !str.ListHas(hdr.Fields[0], col) { //TODO handle multiple fields
			rules = append(rules, col)
		}
	}
	return rules
}
