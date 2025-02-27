// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestJoin_nrows(t *testing.T) {
	test := func(n1, p1, n2, p2, expected int) {
		t.Helper()
		j1n := Join{}
		j1n.joinType = one_n
		jn1 := Join{}
		jn1.joinType = n_one
		assert.T(t).Msg(n1, "/", p1, one_n, n2, "/", p2, "=>", expected).
			This(j1n.nrows(n1, p1, n2, p2)).Is(expected)
		assert.T(t).Msg(n1, "/", p1, n_one, n2, "/", p2, "=>", expected).
			This(jn1.nrows(n2, p2, n1, p1)).Is(expected)
	}
	test(0, 100, 2000, 2000, 0)
	test(1, 100, 2000, 2000, 20)
	test(2, 100, 2000, 2000, 40)
	test(90, 100, 2000, 2000, 1800)
	test(100, 100, 2000, 2000, 2000)

	test(100, 100, 0, 2000, 0)
	test(100, 100, 1, 2000, 1)
	test(100, 100, 10, 2000, 10)
	test(100, 100, 200, 2000, 200)
	test(100, 100, 1800, 2000, 1800)

	test(2, 100, 200, 2000, 40)
	test(2, 100, 10, 2000, 10)
}

func TestJoin_SelectFixedBug(t *testing.T) {
	// Without handleFixed this test should give:
	// 		ASSERT FAILED: msg:  selEnd no data
	db := heapDb()
	db.adm("create cus (c3, ck) key(c3, ck)")
	db.act("insert { c3: 4 } into cus")
	db.adm("create ivc (ck, ik) key(ik)")
	db.adm("create bln (bk, ik) key (ik,bk)")
	query := `
			((cus extend bk = c3)
		join by(ck,bk)
				(bln
			join by(ik)
				(ivc where ik is 4)))
		where ck is "" `
	joinRev = impossible
	defer func() { joinRev = 0 }()
	tran := sizeTran{db.NewReadTran()}
	q := ParseQuery(query, tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	assert.This(Strategy(q)).Like(`
			{1_000 0+250_000} cus^(c3,ck)
			{500/1_000 0+250_000} where ck is ""
			{500/1_000 0+250_000} extend bk = c3
		{1/1_000 0+1_126_000} join n:1 by(ck,bk)
				{0.500x 1_000 0+125_500} bln^(ik,bk)
				{500/1_000 0+125_500} where ik is 4
			{1/1_000 0+376_000} join n:1 by(ik)
				{0.001x 1_000 0+500} ivc^(ik)
				{1/1_000 0+500} where*1 ik is 4 and ck is ""`)
	assert.This(queryAll2(q)).Is("")
}

func TestJoin_EmptyTempIndexBug(t *testing.T) {
	db := heapDb()
	db.adm("create ivc (ck, ik) key(ik)")
	// db.act("insert { c3: 4 } into cus")
	db.adm("create bln (bk, ik) key (ik,bk)")
	query := `
			(bln
		join by(ik)
			(ivc where ik is 4 and ck is ""))`
	joinRev = impossible
	defer func() { joinRev = 0 }()
	tran := sizeTran{db.NewReadTran()}
	q := ParseQuery(query, tran, nil)
	idx := []string{"ck", "ik"}
	q = setupIndex(q, ReadMode, idx, tran)
	// fmt.Println(Strategy(q))
	assert.T(t).Msg("empty TempIndex").
		That(!strings.Contains(Strategy(q), "TEMPINDEX()"))
	q.Select(idx, []string{"", ""})
}

func setupIndex(q Query, mode Mode, index []string, tran QueryTran) Query {
	q = q.Transform()
	fixcost, varcost := Optimize(q, mode, index, float64(1))
	if fixcost+varcost >= impossible {
		panic("impossible")
	}
	q = SetApproach(q, index, float64(1), tran)
	return q
}

// func rowstr(hdr *Header, row Row) string {
// 	if row == nil {
// 		return "nil"
// 	}
// 	var sb strings.Builder
// 	sep := ""
// 	for _, col := range hdr.Columns {
// 		val := row.GetVal(hdr, col, nil, nil)
// 		if val != EmptyStr {
// 			fmt.Fprint(&sb, sep, col, "=", AsStr(val))
// 			sep = " "
// 		}
// 	}
// 	return sb.String()
// }
