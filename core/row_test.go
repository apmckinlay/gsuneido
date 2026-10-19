// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestGetRaw(*testing.T) {
	hdr := SimpleHeader([]string{"one", "two", "two_lower!"})
	var rb RecordBuilder
	rb.Add(SuStr("hello"))
	rb.Add(SuStr("Hello World"))
	rec := rb.Build()
	row := Row{rec}
	assert.This(row.GetRaw(hdr, "one")).Is(Pack(SuStr("hello")))
	assert.This(row.GetRaw(hdr, "two")).Is(Pack(SuStr("Hello World")))
	assert.This(row.GetRaw(hdr, "two_lower!")).Is(Pack(SuStr("hello world")))
}

func TestRowSameAs(t *testing.T) {
	rec := new(RecordBuilder).Add(SuStr("hello")).Build()
	copy := Record(strings.Clone(string(rec)))
	test := func(row, other Row, expected bool) {
		t.Helper()
		assert.T(t).This(row.SameAs(other)).Is(expected)
	}
	test(Row{rec, ""}, Row{copy, ""}, true)
	test(Row{rec}, Row{copy, ""}, false)
	test(Row{rec}, Row{"\x00"}, false)
	test(Row{""}, Row{"\x00"}, false)
	test(nil, Row{}, true)
}

func BenchmarkRow_GetRawVal(b *testing.B) {
	var row Row
	hdr := SimpleHeader([]string{"foo", "bar"})
	for b.Loop() {
		row.GetRawVal(hdr, "baz", nil, nil)
	}
}
