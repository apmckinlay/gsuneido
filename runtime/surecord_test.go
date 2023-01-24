// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"math/rand"
	"strings"
	"sync"
	"testing"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestSuRecord(t *testing.T) {
	r := new(SuRecord)
	assert.T(t).This(r.Type()).Is(types.Record)
	assert.T(t).This(r.String()).Is("[]")
	r.Set(SuStr("a"), SuInt(123))
	assert.T(t).This(r.String()).Is("[a: 123]")
}

func TestSuRecord_ReadonlyUnpack(t *testing.T) {
	b := RecordBuilder{}
	b.Add(SuInt(123))
	b.Add(SuStr("foobar"))
	rec := b.Build()
	dbrec := DbRec{Record: rec}
	row := Row{dbrec}

	hdr := NewHeader([][]string{{"num", "str"}}, []string{"num", "str"})
	surec := SuRecordFromRow(row, hdr, "", nil)

	assert.T(t).This(surec.Get(nil, SuStr("str"))).Is(SuStr("foobar"))
	surec.SetReadOnly()
	assert.T(t).This(surec.Get(nil, SuStr("num"))).Is(SuInt(123))
}

func TestSuRecord_Concurrency(t *testing.T) {
	var nThreads = 8
	var nActions = 400000
	if testing.Short() {
		nThreads = 4
		nActions = 4000
	}
	cols := []string{"a", "b", "c", "d"}
	randCol := func() SuStr {
		return SuStr(cols[rand.Intn(len(cols))])
	}
	rb := RecordBuilder{}
	rb.Add(randCol())
	rb.Add(randCol())
	row := Row{DbRec{Record: rb.Build()}}
	getrec, setrec := func() (func() *SuRecord, func(*SuRecord)) {
		rec := NewSuRecord()
		rec.SetConcurrent()
		var lock sync.Mutex
		return func() *SuRecord {
				lock.Lock()
				defer lock.Unlock()
				return rec
			}, func(r *SuRecord) {
				lock.Lock()
				defer lock.Unlock()
				rec = r
			}
	}()

	rule := &SuBuiltinMethod{
		Fn: func(th *Thread, this Value, _ []Value) Value {
			switch rand.Intn(2) {
			case 0:
				return this.Get(th, randCol())
			case 1:
				this.Put(th, randCol(), SuStr(randCol()))
			}
			return nil
		},
		BuiltinParams: BuiltinParams{ParamSpec: ParamSpec0}}

	var wg sync.WaitGroup
	run := func() {
		th := &Thread{}
		defer wg.Done()
		for i := 0; i < nActions; i++ {
			switch rand.Intn(12) {
			case 0:
				getrec().Get(th, randCol())
			case 1:
				getrec().Put(th, randCol(), SuStr(randCol()))
			case 2:
				r := NewSuRecord()
				r.SetConcurrent()
				r.AttachRule(SuStr("c"), rule)
				r.AttachRule(SuStr("d"), rule)
				setrec(r)
			case 3:
				r := SuRecordFromRow(row, SimpleHeader(cols), "", nil)
				r.SetConcurrent()
				r.AttachRule(SuStr("c"), rule)
				r.AttachRule(SuStr("d"), rule)
				setrec(r)
			case 4:
				r := getrec().Copy().(*SuRecord)
				r.SetConcurrent()
				setrec(r)
			case 5:
				getrec().ToObject()
			case 6:
				getrec().Delete(th, randCol())
			case 7:
				getrec().Invalidate(th, string(randCol()))
			case 8:
				getrec().ToRecord(th, SimpleHeader(cols))
			case 9:
				getrec().GetDeps(string(randCol()))
			case 10:
				getrec().SetDeps(string(randCol()), "a,c")
			case 11:
				func() {
					// ignore object-modified-during-iteration
					defer func() { recover() }()
					getrec().Display(th)
				}()
			}
		}
	}
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		go run()
	}
	wg.Wait()
}

func TestSuRecord_validRule(t *testing.T) {
	test := func(s string, expected bool) {
		assert.This(validRule(s)).Is(expected)
	}
	test("foo", true)
	test("foo_bar", true)
	test("foo?_bar", true)
	test("123", true)

	test("foo\n", false)
	test("\nfoo", false)
	test("foo bar", false)
	test(strings.Repeat("x", maxRule+1), false)
}

func BenchmarkSuRecord(b *testing.B) {
	th := &Thread{}
	var rec SuRecord
	for i := 0; i < b.N; i++ {
		rec = SuRecord{}
		rec.Get(th, SuStr("foo"))
	}
}
