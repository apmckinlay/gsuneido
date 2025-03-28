// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDatabaseDropTable(t *testing.T) {
	db := CreateDb(stor.HeapStor(16 * 1024))
	createTbl(db)
	db.CheckerSync()
	assert.T(t).That(db.Drop("nonexistent") != nil)
	assert.T(t).That(db.Drop("mytable") == nil)
	assert.T(t).That(db.Drop("mytable") != nil)
}
