// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query_test

import (
	"github.com/apmckinlay/gsuneido/db19/tools"
	"github.com/apmckinlay/gsuneido/dbms/query"
)

// Wire up database compaction for the bitmap benchmarks.
// It lives in the external test package because db19/tools imports query.
func init() {
	query.CompactDb = func(dbfile string) error {
		_, _, _, _, err := tools.Compact(dbfile)
		return err
	}
}
