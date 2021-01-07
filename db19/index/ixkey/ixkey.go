// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package ixkey handles specifying and encoding index key strings
// that are directly comparable.
// Single field keys are not encoded.
// Fields are separated by two zero bytes 0,0.
// Zero bytes are encoded as 0,1.
// Normally the values will be packed,
// but this is not required as long as they compare directly.
package ixkey

import (
	"fmt"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

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

// Encoder builds keys incrementally.
// Note: Do not use this for single field keys - they should not be encoded.
type Encoder struct {
	buf []byte
}

// Add appends a field value
func (b *Encoder) Add(fld string) {
	if b.buf == nil {
		b.buf = make([]byte, 0, 2*(len(fld)+2))
	} else {
		b.buf = append(b.buf, 0, 0) // separator
	}
	b.buf = encode(b.buf, fld)
}

// String returns the key and resets the Encoder to be empty.
func (b *Encoder) String() string {
	s := hacks.BStoS(b.buf)
	b.buf = nil // reset
	return s
}

// Key builds a key from a data Record using a Spec.
func (spec *Spec) Key(rec Record) string {
	fields := spec.Fields
	if len(fields) == 0 {
		return ""
	}
	if len(fields) == 1 && len(spec.Fields2) == 0 {
		return getRaw(rec, fields[0]) // don't need to encode single field keys
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

// Increment adds the smallest amount to a key,
// used to convert > to >= or <= to <
func (spec *Spec) Increment(key string) string {
	if spec.raw() {
		return key + "\x00"
	}
	// encoded
	return key + "\x00\x00" // add empty field trailing field
}

func (spec *Spec) raw() bool {
	return len(spec.Fields) == 0 ||
		(len(spec.Fields) == 1 && len(spec.Fields2) == 0)
}
