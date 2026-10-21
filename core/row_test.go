// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"strings"
	"testing"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestGetRaw(*testing.T) {
	hdr := SimpleHeader([]string{"one", "two", "two_lower!"})
	var rb RecordBuilder
	rb.Add(SuStr("hello"))
	rb.Add(SuStr("Hello World"))
	rec := rb.Build()
	row := NewRow(rec)
	assert.This(row.GetRaw(hdr, "one")).Is(Pack(SuStr("hello")))
	assert.This(row.GetRaw(hdr, "two")).Is(Pack(SuStr("Hello World")))
	assert.This(row.GetRaw(hdr, "two_lower!")).Is(Pack(SuStr("hello world")))
}

func TestRowSameAs(t *testing.T) {
	rec := new(RecordBuilder).Add(SuStr("hello")).Build()
	copy := Record(strings.Clone(string(rec)))
	test := func(row, other Row, expected bool) {
		t.Helper()
		assert.T(t).This(row.EqualRecs(other)).Is(expected)
	}
	test(NewRow(rec).With(""), NewRow(copy).With(""), true)
	test(NewRow(rec), NewRow(copy).With(""), false)
	test(NewRow(rec), NewRow("\x00"), false)
	test(NewRow(""), NewRow("\x00"), false)
	test(Row{}, Row{}, true)
}

func TestRowGet(t *testing.T) {
	test := func(rec Record) {
		t.Helper()
		row := NewRow(rec)
		assert.T(t).This(row.Len()).Is(1)
		assert.T(t).This(row.Get(0)).Is(rec)
		if rec != "" {
			assert.T(t).That(unsafe.StringData(string(row.Get(0))) == unsafe.StringData(string(rec)))
		} else {
			assert.T(t).This(row.Ptrs.Get(0) == nil).Is(true)
		}
		assert.T(t).This(testing.AllocsPerRun(100, func() {
			_ = row.Get(0)
		})).Is(0)
		assert.T(t).This(testing.AllocsPerRun(100, func() {
			row = NewRow(rec)
		})).Is(0)
		row = NewRow("").With(rec).With("\x00")
		assert.T(t).This(row.Get(1)).Is(rec)
		assert.T(t).This(testing.AllocsPerRun(100, func() {
			_ = row.Get(1)
		})).Is(0)
	}
	test("")
	s := strings.Clone("not a record")
	test(Record(s[:0]))
	test("\x00")
	test(new(RecordBuilder).Add(SuStr("hello")).Build())
	test(new(RecordBuilder).Add(SuStr(strings.Repeat("a", 256))).Build())
	test(new(RecordBuilder).Add(SuStr(strings.Repeat("a", 65536))).Build())
	row := NewRow("\x00")
	assert.T(t).This(func() { row.Get(-1) }).Panics("index out of range")
	assert.T(t).This(func() { row.Get(1) }).Panics("index out of range")
}

func TestJoinRows(t *testing.T) {
	rec1 := new(RecordBuilder).Add(SuStr("one")).Build()
	rec2 := new(RecordBuilder).Add(SuStr("two")).Build()
	one := NewRow(rec1)
	two := NewRow(rec2).With("")
	row := JoinRows(one, two)
	assert.T(t).This(row.Len()).Is(3)
	assert.T(t).This(row.Get(0)).Is(rec1)
	assert.T(t).This(row.Get(1)).Is(rec2)
	assert.T(t).This(row.Get(2)).Is(Record(""))
	assert.T(t).This(one.Len()).Is(1)
	assert.T(t).This(two.Len()).Is(2)
	assert.T(t).This(JoinRows(Row{}, one).EqualRecs(one)).Is(true)
	assert.T(t).This(JoinRows(one, Row{}).EqualRecs(one)).Is(true)
	assert.T(t).This(JoinRows(Row{}, Row{}).Len()).Is(0)
	assert.T(t).This(Row{}.With(rec1).EqualRecs(one)).Is(true)
	assert.T(t).This(row.String()).Is("[{\"one\"} {\"two\"} {}]")
}

func TestRowUnion(t *testing.T) {
	hdr := NewHeader([][]string{{"one"}, {"one"}}, []string{"one"})
	rec := new(RecordBuilder).Add(SuStr("value")).Build()
	test := func(row Row, expected string) {
		t.Helper()
		assert.T(t).This(row.GetRaw(hdr, "one")).Is(expected)
	}
	test(NewRow(rec).With(""), Pack(SuStr("value")))
	test(NewRow("").With(rec), Pack(SuStr("value")))
	test(NewRow("").With(""), "")
}

func BenchmarkRow_GetRawVal(b *testing.B) {
	var row Row
	hdr := SimpleHeader([]string{"foo", "bar"})
	for b.Loop() {
		row.GetRawVal(hdr, "baz", nil, nil)
	}
}
