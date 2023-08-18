// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

import (
	"fmt"
	"hash/crc64"
	"testing"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	. "github.com/apmckinlay/gsuneido/dbms/query"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/apmckinlay/gsuneido/util/generic/hmap"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/hash"
)

func TestQuery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	// Global.TestDef("Rule_c",
	// 	compile.Constant("function() { return .b }"))
	db, err := db19.OpenDatabaseRead("../suneido.db")
	if err != nil {
		panic(err.Error())
	}
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	tran := db.NewReadTran()
	s := `(((cus join ivc) join (((bln extend c1 = ik)) where bk is "82")) union ((cus join (ivc union ivc)) join bln)) sort b2,c4,i4`
	q := ParseQuery(s, tran, nil)
	// trace.QueryOpt.Set()
	// trace.JoinOpt.Set()
	q, _, _ = Setup(q, ReadMode, tran)

	fmt.Println("----------------")
	fmt.Println(Strategy(q))
	th := &Thread{}
	n := 0
	hdr := q.Header()
	fields := slc.Without(hdr.GetFields(), "-")
	hashes := make(map[uint64]struct{})
	for {
		row := q.Get(th, Next)
		if row == nil {
			break
		}
		hash := hashRow(hdr, fields, row)
		if _, ok := hashes[hash]; ok {
			panic("duplicate hash")
		}
		hashes[hash] = struct{}{}
		n++
	}
	fmt.Println(n, "rows")
	exit.RunFuncs()
}

func hashRow(hdr *Header, fields []string, row Row) uint64 {
	hash := uint64(0)
	for _, fld := range fields {
		hash = hash*31 + hashPacked(row.GetRaw(hdr, fld))
	}
	return hash
}

var ecma = crc64.MakeTable(crc64.ECMA)

func hashPacked(p string) uint64 {
	if len(p) > 0 && p[0] >= PackObject {
		return hashObject(p)
	}
	return crc64.Checksum(hacks.Stobs(p), ecma)
}

func hashObject(p string) uint64 {
	hash := uint64(0)
	for i := 0; i < len(p); i++ {
		// use simple addition to be insensitive to member order
		hash += uint64(p[i])
	}
	return hash
}

func TestQuery2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	db, err := db19.OpenDatabaseRead("../suneido.db")
	if err != nil {
		panic(err.Error())
	}
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	tran := db.NewReadTran()
	s := `aln where ik is "67" sort ik`
	q := ParseQuery(s, tran, nil)
	// trace.QueryOpt.Set()
	q = q.Transform()
	const frac = 100
	var index []string
	fixcost, varcost := Optimize(q, ReadMode, index, frac)
	if fixcost+varcost >= 9999999999 {
		panic("invalid query: " + q.String())
	}
	q = SetApproach(q, index, frac, tran)

	fmt.Println("----------------")
	fmt.Println(Strategy(q))
	th := &Thread{}
	fmt.Println(q.Get(th, Next))
}

func BenchmarkProject_Old(b *testing.B) {
	db, err := db19.OpenDatabaseRead("../suneido.db")
	if err != nil {
		panic(err.Error())
	}
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	tran := db.NewReadTran()
	q := ParseQuery("gl_transactions", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	data := make([]Row, 0, 1000)
	for row := q.Get(nil, Next); row != nil; row = q.Get(nil, Next) {
		data = append(data, row)
	}
	hdr := q.Header()
	cols := []string{"gltran_currency", "gltran_date", "glacct_num"}
	for i := 0; i < b.N; i++ {
		results := make(map[string]Row)
		for _, row := range data {
			key := ixkey.Make(row, hdr, cols, nil, nil)
			if _, ok := results[key]; !ok {
				results[key] = row
			}
		}
	}
}

func BenchmarkProject_Hmap(b *testing.B) {
	db, err := db19.OpenDatabaseRead("../suneido.db")
	if err != nil {
		panic(err.Error())
	}
	MakeSuTran = func(qt QueryTran) *SuTran {
		return nil
	}
	tran := db.NewReadTran()
	q := ParseQuery("gl_transactions", tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	data := make([]Row, 0, 1000)
	for row := q.Get(nil, Next); row != nil; row = q.Get(nil, Next) {
		data = append(data, row)
	}
	hdr := q.Header()
	cols := []string{"gltran_currency", "gltran_date", "glacct_num"}
	hn, en := 0, 0
	type T struct {
		row  Row
		hash uint32
	}
	hfn := func(t T) uint32 {
		return t.hash
	}
	eqfn := func(x, y T) bool {
		en++
		return x.hash == y.hash &&
			equalCols(x.row, y.row, hdr, cols, nil, nil)
	}
	for i := 0; i < b.N; i++ {
		results := hmap.NewHmapFuncs[T, struct{}](hfn, eqfn)
		for _, row := range data {
			hn++
			t := T{row: row, hash: hashCols(row, hdr, cols, nil, nil)}
			results.GetPut(t, struct{}{})
		}
	}
	fmt.Println("rows", len(data), "hn", hn, "en", en)
}

func hashCols(row Row, hdr *Header, cols []string, th *Thread, st *SuTran) uint32 {
	h := uint32(31)
	for _, col := range cols {
		x := row.GetRawVal(hdr, col, th, st)
		h = 31*h + hash.String(x)
	}
	return h
}
func equalCols(x, y Row, hdr *Header, cols []string, th *Thread, st *SuTran) bool {
	for _, col := range cols {
		if x.GetRawVal(hdr, col, th, st) != y.GetRawVal(hdr, col, th, st) {
			return false
		}
	}
	return true
}
