// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/pack"
)

var xlat = []byte{
	PackFalseOther:  PackFalse,
	PackTrueOther:   PackTrue,
	PackMinusOther:  PackMinus,
	PackPlusOther:   PackPlus,
	PackStringOther: PackString,
	PackDateOther:   PackDate,
	PackObjectOther: PackObject,
	PackRecordOther: PackRecord,
}

var rxlat = []byte{
	PackFalse:  PackFalseOther,
	PackTrue:   PackTrueOther,
	PackMinus:  PackMinusOther,
	PackPlus:   PackPlusOther,
	PackString: PackStringOther,
	PackDate:   PackDateOther,
	PackObject: PackObjectOther,
	PackRecord: PackRecordOther,
}

func ConvertRecord(buf []byte) {
	s := hacks.BStoS(buf)
	rec := Record(s)
	for i := 0; i < rec.Count(); i++ {
		pos, end := rec.GetRange(i)
		ConvertValue(xlat, buf[pos:end], s[pos:end])
	}
}

func RevertValue(buf []byte, s string) {
	ConvertValue(rxlat, buf, s)
}

func ConvertValue(xlat []byte, buf []byte, s string) {
	if len(buf) == 0 {
		return // ""
	}
	buf[0] = xlat[buf[0]] // handles simple types
	if len(buf) > 1 && (buf[0] == PackObjectOther || buf[0] == PackRecordOther) {
		convertObject(xlat, buf[1:], s[1:])
	}
}

func convertObject(xlat []byte, buf []byte, s string) {
	dcdr := pack.NewDecoder(s)
	n := int(dcdr.VarUint())
	for i := 0; i < n; i++ {
		convertSizedValue(xlat, buf, dcdr)
	}
	n = int(dcdr.VarUint())
	for i := 0; i < n; i++ {
		convertSizedValue(xlat, buf, dcdr)
		convertSizedValue(xlat, buf, dcdr)
	}

}

func convertSizedValue(xlat []byte, buf []byte, dcdr *pack.Decoder) {
	size := int(dcdr.VarUint())
	i := len(buf) - dcdr.Remaining()
	ConvertValue(xlat, buf[i:i+size], dcdr.Get(size))
}
