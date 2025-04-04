// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"
	"fmt"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/pack"
)

// SuTimestamp is an extension of SuDate
// that adds an extra byte of precision
// for timestamps.
type SuTimestamp struct {
	SuDate
	extra uint8
}

func (d SuTimestamp) String() string {
	return fmt.Sprintf("#%04d%02d%02d.%02d%02d%02d%03d%03d",
		d.Year(), d.Month(), d.Day(),
		d.Hour(), d.Minute(), d.Second(), d.Millisecond(),
		d.extra)
}

func (d SuTimestamp) Hash() uint64 {
	h := d.SuDate.Hash()
	h = 31*h + uint64(d.extra)
	return h
}

// packing

func (SuTimestamp) PackSize(*packing) int {
	return 10
}

func (d SuTimestamp) Pack(pk *packing) {
	assert.That(d.extra != 0)
	d.SuDate.Pack(pk)
	pk.Put1(d.extra)
}

func UnpackTimestamp(sd SuDate, d *pack.Decoder) SuTimestamp {
	extra := d.Get1()
	assert.That(extra != 0)
	return SuTimestamp{SuDate: sd, extra: extra}
}

// compare

func (d SuTimestamp) Equal(other any) bool {
	return d == other
}

func (d SuTimestamp) Compare(other Value) int {
	if cmp := cmp.Compare(ordDate, Order(other)); cmp != 0 {
		return cmp * 2
	}
	if st, ok := other.(SuTimestamp); ok {
		return CompareSuTimestamp(d, st)
	}
	return CompareSuTimestamp(d, SuTimestamp{SuDate: other.(SuDate)})
}

func CompareSuTimestamp(d1, d2 SuTimestamp) int {
	if cmp := d1.SuDate.Compare(d2.SuDate); cmp != 0 {
		return cmp
	}
	return cmp.Compare(d1.extra, d2.extra)
}
