// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"testing"
)

func TestIntersect(t *testing.T) {
	db := heapDb()
	db.adm("create t1 (a, b) key (a)")
	db.act(`insert {a: "", b: 2} into t1`)
	db.adm("create t2 (b, c) key (c)")
	db.act(`insert {b: 2, c: ""} into t2`)
	queryHashAll(db.Database, `t1 intersect t2`)
}
