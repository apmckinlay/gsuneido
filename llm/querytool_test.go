// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestQueryTool(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	dbms := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbms }

	result, err := queryTool("tables")
	assert.That(err == nil)
	assert.This(result.Results).Is("[\n" +
		"[\"table\", \"nrows\", \"totalsize\"]\n" +
		"[\"columns\", 14, 0]\n" +
		"[\"indexes\", 0, 0]\n" +
		"[\"tables\", 4, 0]\n" +
		"[\"views\", 0, 0]\n" +
		"]\n")
}

func TestQueryToolSizeLimit(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	dbms := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbms }
	dbms.Admin("create big (k, a) key(k)", nil)

	th := core.NewThread(core.MainThread)
	defer th.Close()
	tran := dbms.Transaction(true)
	xs := strings.Repeat("x", 1500)
	for i := range 10 {
		n := tran.Action(th, fmt.Sprintf("insert { k: %d, a: %q } into big", i, xs))
		assert.This(n).Is(1)
	}
	tran.Complete()

	result, err := queryTool("big")
	assert.That(err == nil)
	assert.That(result.HasMore)
	head := "[\n[\"k\", \"a\"]\n" +
		"[0, \"" + xs + "\"]\n" +
		"[1, \"" + xs + "\"]\n"
	tail := "[4, \"" + xs + "\"]\n" +
		"[5, \"" + xs + "\"]\n" +
		"[6]\n" +
		"// too large, truncated\n" +
		"]\n"
	assert.That(strings.HasPrefix(result.Results, head))
	assert.That(strings.HasSuffix(result.Results, tail))
}

func TestQueryToolRowLimit(t *testing.T) {
	assert := assert.T(t)
	db := db19.CreateDb(stor.HeapStor(8192))
	db19.StartConcur(db, 50*time.Millisecond)
	dbms := dbms.NewDbmsLocal(db)
	core.GetDbms = func() core.IDbms { return dbms }
	dbms.Admin("create many (k) key(k)", nil)

	th := core.NewThread(core.MainThread)
	defer th.Close()
	tran := dbms.Transaction(true)
	for i := range 150 {
		n := tran.Action(th, fmt.Sprintf("insert { k: %d } into many", i))
		assert.This(n).Is(1)
	}
	tran.Complete()

	result, err := queryTool("many sort k")
	assert.That(err == nil)
	assert.That(result.HasMore)
	head := "[\n[\"k\"]\n[0]\n[1]\n[2]\n[3]\n[4]\n[5]\n[6]\n[7]\n[8]\n[9]\n"
	tail := "[90]\n[91]\n[92]\n[93]\n[94]\n[95]\n[96]\n[97]\n[98]\n[99]\n" +
		"// too large, truncated\n]\n"
	assert.That(strings.HasPrefix(result.Results, head))
	assert.That(strings.HasSuffix(result.Results, tail))
}
