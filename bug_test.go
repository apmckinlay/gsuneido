// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"testing"

	. "github.com/apmckinlay/gsuneido/core"
	. "github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBug(t *testing.T) {
	assert.TestOnlyIndividually(t)

	openDbms()
	defer db.CloseKeepMapped()

	query := `(cus leftjoin ((ivc extend c3 = ik) union ivc)) where ik is "12"`
	th := &Thread{}

	tran := db.NewReadTran()
	q := ParseQuery(query, tran, nil)
	q, _, _ = Setup(q, ReadMode, tran)
	fmt.Println("optimized:", String(q))

	hdr := q.Header()
	fmt.Println("\nheader fields:", hdr.GetFields())
	fmt.Println("header columns:", hdr.Columns)

	fmt.Println("\n--- Simple result ---")
	simple := q.Simple(th)
	h2 := NewQueryHasher(hdr).CheckDups()
	for i, row := range simple {
		h2.Row(row)
		fmt.Printf("simple[%d]: ", i)
		for _, fld := range hdr.GetFields() {
			if fld != "-" {
				fmt.Printf("%s=%q ", fld, row.GetRawVal(hdr, fld, nil, nil))
			}
		}
		fmt.Println()
	}

	fmt.Println("\n--- Get result ---")
	q.Rewind()
	h1 := NewQueryHasher(hdr).CheckDups()
	for row := q.Get(th, Next); row != nil; row = q.Get(th, Next) {
		fmt.Print("row: ")
		for _, fld := range hdr.GetFields() {
			if fld != "-" {
				fmt.Printf("%s=%q ", fld, row.GetRawVal(hdr, fld, nil, nil))
			}
		}
		fmt.Println()
		h1.Row(row)
	}

	assert.This(h2.Result(true)).Is(h1.Result(true))
}
