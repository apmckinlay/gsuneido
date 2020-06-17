// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestMetadata(t *testing.T) {
	base := NewTableInfoHtbl(0)
	base.Put(&TableInfo{
		table: 1,
		nrows: 100,
		size: 1000,
		indexes: []IndexInfo{
			{ root: 11111111111, treeLevels: 0 },
			{ root: 111111111111, treeLevels: 1 },
		},
	})
	base.Put(&TableInfo{
		table: 2,
		nrows: 200,
		size: 2000,
		indexes: []IndexInfo{
			{ root: 22222222222, treeLevels: 0 },
			{ root: 222222222222, treeLevels: 2 },
		},
	})
	over := NewTableInfoHtbl(0)
	over.Put(&TableInfo{
		table: 2,
		nrows: 9,
		size: 99,
		indexes: []IndexInfo{
			{ root: 22222222220, treeLevels: 1 },
			{ root: 222222222220, treeLevels: 2 },
		},
	})
	merged := base.Merge(over)
	buf := merged.Write()
	reread := ReadTablesInfo(buf)
	Assert(t).That(*reread.Get(1), Equals(*base.Get(1)))
	Assert(t).That(*reread.Get(2), Equals(TableInfo{
		table: 2,
		nrows: 209,
		size: 2099,
		indexes: []IndexInfo{
			{ root: 22222222220, treeLevels: 1 },
			{ root: 222222222220, treeLevels: 2 },
		},
	}))
}
