// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ixkey handles specifying and encoding index key strings
// that are directly comparable.
// Single field index keys are not encoded.
// But a single value for a multi-field index still needs to be encoded.
// Fields are separated by two zero bytes 0,0.
// Zero bytes are encoded as 0,1.
// Normally the values will be packed,
// but this is not required as long as they compare directly.
package ixkey

import (
	"fmt"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

const Min = ""
const Max = "\xff\xff\xff\xff\xff\xff\xff\xff"
const Sep = "\x00\x00"

// Technically there is no maximum key string.
// However, in practice keys are packed values, encoded when composite.
// Packed values start with a type byte from 0 to 7 so 0xff will be larger.
// And 0xff will be larger than any ascii strings.

// Spec specifies the field(s) in an index key
type Spec struct {
	// Fields specifies the fields in the key.
	Fields []int
	// Fields2 is used for unique indexes (that allow multiple empty keys).
	// It will only be used if all of the Fields value are empty.
	Fields2 []int
}

func (spec *Spec) String() string {
	return fmt.Sprint("ixspec ", spec.Fields, ",", spec.Fields2)
}

// Encoder builds keys incrementally. Zero value is ready to use.
// Note: Do not use this for single field keys - they should not be encoded.
type Encoder struct {
	buf []byte
}

// Add appends a field value
func (e *Encoder) Add(fld string) {
	cklen(len(e.buf) + len(fld))
	if e.buf == nil {
		e.buf = make([]byte, 0, 2*(len(fld)+2))
	} else {
		e.buf = append(e.buf, 0, 0) // separator
	}
	e.buf = encode(e.buf, fld)
}

// String returns the key and resets the Encoder to be empty.
// Trailing field separators (empty fields) are trimmed.
func (e *Encoder) String() string {
	s := hacks.BStoS(e.buf)
	e.buf = nil // reset
	for strings.HasSuffix(s, "\x00\x00") {
		s = s[:len(s)-2]
	}
	return s
}

func (e *Encoder) Dup() *Encoder {
	var e2 Encoder
	// use append so if e.buf is nil, e2.buf will be nil
	e2.buf = append([]byte(nil), e.buf...)
	return &e2
}

// Key builds a key from a data Record using a Spec.
func (spec *Spec) Key(rec Record) string {
	assert.That(spec.Fields != nil)
	fields := spec.Fields
	if len(fields) == 0 {
		return ""
	}
	if len(fields) == 1 && len(spec.Fields2) == 0 {
		s := getRaw(rec, fields[0]) // don't need to encode single field keys
		cklen(len(s))
		return s
	}
	n := 0
	lastNonEmpty := -1
	for i, field := range fields {
		fldlen := len(rec.GetRaw(field))
		if fldlen > 0 {
			lastNonEmpty = i
		}
		n += fldlen
	}
	if lastNonEmpty == -1 { // fields all empty
		if len(spec.Fields2) == 0 {
			return ""
		}
		for _, field := range spec.Fields2 {
			n += fieldLen(rec, field)
		}
	} else {
		fields = fields[:lastNonEmpty+1]
	}
	n += 2 * len(fields) // for separators (2 bytes extra)
	n += n / 16          // allow for some escapes
	cklen(n)
	buf := make([]byte, 0, n)
	if lastNonEmpty == -1 {
		for range fields {
			buf = append(buf, 0, 0) // separator
		}
		fields = spec.Fields2
	}
	for i, f := range fields {
		if i > 0 {
			buf = append(buf, 0, 0) // separator
		}
		buf = encode(buf, getRaw(rec, f))
	}
	return hacks.BStoS(buf)
}

func cklen(n int) {
	if n > maxKey {
		panic(fmt.Sprintf("key too large (%d > %d)", n, maxKey))
	}
}

const maxKey = 4096 // ???

// Encodes returns whether the Spec requires encoding.
func (spec *Spec) Encodes() bool {
	return len(spec.Fields) > 1 || len(spec.Fields2) > 0
}

func Encode(s string) string {
	if i := strings.IndexByte(s, 0); i == -1 {
		return s
	}
	buf := make([]byte, 0, len(s)+4)
	buf = encode(buf, s)
	return hacks.BStoS(buf)
}

// encode appends encoded b to buf and returns resulting buf
func encode(buf []byte, b string) []byte {
	for len(b) > 0 {
		i := strings.IndexByte(b, 0)
		if i == -1 { // no zero bytes
			buf = append(buf, b...)
			break
		}
		// b[i] == 0
		i++
		buf = append(buf, b[:i]...) // copy up to and including zero
		buf = append(buf, 1)
		b = b[i:]
	}
	return buf
}

func fieldLen(rec Record, field int) int {
	if field < 0 {
		field = -field - 2 // _lower!
	}
	return len(rec.GetRaw(field))
}

// getRaw does Record.GetRaw and handles _lower! fields
func getRaw(rec Record, field int) string {
	if field >= 0 {
		return rec.GetRaw(field)
	}
	field = -field - 2 // _lower!
	return PackedToLower(rec.GetRaw(field))
}

// Compare compares the specified fields of the two records
// without building keys for them
func (spec *Spec) Compare(r1, r2 Record) int {
	empty := true
	for _, f := range spec.Fields {
		var x1, x2 string
		var cmp int
		if f < 0 { // _lower!
			f = -f - 2
			x1 = r1.GetRaw(f)
			x2 = r2.GetRaw(f)
			cmp = PackedCmpLower(x1, x2)
		} else {
			x1 = r1.GetRaw(f)
			x2 = r2.GetRaw(f)
			cmp = strings.Compare(x1, x2)
		}
		if cmp != 0 {
			return cmp
		}
		if x1 != "" || x2 != "" {
			empty = false
		}
	}
	if empty {
		for _, f := range spec.Fields2 {
			// NOTE: assumes fields2 will not be _lower!
			if cmp := strings.Compare(r1.GetRaw(f), r2.GetRaw(f)); cmp != 0 {
				return cmp
			}
		}
	}
	return 0
}

func (spec *Spec) Increment(key string) string {
	if spec.raw() {
		return key + "\x00"
	}
	// encoded
	return key + Sep // add empty field trailing field
}

func (spec *Spec) raw() bool {
	return len(spec.Fields) == 0 ||
		(len(spec.Fields) == 1 && len(spec.Fields2) == 0)
}

func (spec *Spec) Trunc(n int) *Spec {
	return &Spec{Fields: spec.Fields[:n]}
}

// Decode is for tests and debugging
func Decode(comp string) []string {
	if comp == "" {
		return nil
	}
	parts := strings.Split(comp, Sep)
	result := make([]string, len(parts))
	for i, p := range parts {
		result[i] = strings.ReplaceAll(p, "\x00\x01", "\x00")
	}
	return result
}

func DecodeValues(comp string) []Value {
	if comp == "" {
		return nil
	}
	parts := strings.Split(comp, Sep)
	result := make([]Value, len(parts))
	for i, p := range parts {
		s := strings.ReplaceAll(p, "\x00\x01", "\x00")
		if s == Max {
			result[i] = SuStr("<max>")
		} else {
			result[i] = Unpack(s)
		}
	}
	return result
}

// HasPrefix is prefix by field.
// i.e. the prefix must be followed by a field separator (or at the end)
func HasPrefix(s, prefix string) bool {
	sn := len(s)
	pn := len(prefix)
	return sn >= pn && s[0:pn] == prefix && // byte-wise prefix
		(sn == pn ||
			(sn >= pn+2 && s[pn:pn+2] == Sep))
}

func Make(row Row, hdr *Header, cols []string, th *Thread, st *SuTran) string {
	if len(cols) == 1 { // WARNING: only correct for keys
		s := row.GetRawVal(hdr, cols[0], th, st)
		cklen(len(s))
		return s
	}
	enc := Encoder{}
	for _, col := range cols {
		enc.Add(row.GetRawVal(hdr, col, th, st))
	}
	return enc.String()
}
