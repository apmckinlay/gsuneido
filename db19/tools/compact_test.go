// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"testing"

	"github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/stats"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestCreateStatsTable(t *testing.T) {
	test := func(tally stats.StatsTally, exists bool) {
		t.Helper()
		db := db19.CreateDb(stor.HeapStor(8192))
		defer db.Close()
		createStatsTable(db, tally)
		assert.T(t).This(db.GetState().Meta.GetRoSchema(db19.StatsTableName) != nil).Is(exists)
	}
	test(nil, false)
	test(stats.StatsTally{}, false)
	test(stats.StatsTally{"foo": {"bar": &stats.StatsColumn{}}}, false)
	col := &stats.StatsColumn{}
	col.Add("value")
	test(stats.StatsTally{"foo": {"bar": col}}, true)
}
