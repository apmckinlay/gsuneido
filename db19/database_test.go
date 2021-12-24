// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"os"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDatabaseDropTable(t *testing.T) {
	db := createDb()
	db.CheckerSync()
	defer func() { db.Close(); os.Remove("tmp.db") }()
	assert.T(t).That(db.Drop("nonexistent") != nil)
	assert.T(t).That(db.Drop("mytable") == nil)
	assert.T(t).That(db.Drop("mytable") != nil)
}
