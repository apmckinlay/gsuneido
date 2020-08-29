// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/comp"
	"github.com/apmckinlay/gsuneido/database/db19/ixspec"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

type Database struct {
	store *stor.Stor

	// state is the central immutable state of the database.
	// It must be accessed atomically and only updated via UpdateState.
	state stateHolder

	ck Checker
}

//-------------------------------------------------------------------

func init() {
	btree.GetLeafKey = getLeafKey
}

func getLeafKey(store *stor.Stor, is *ixspec.T, off uint64) string {
	rec := offToRec(store, off)
	return comp.Key(rt.Record(rec), is.Cols, is.Cols2)
}

func mkcmp(store *stor.Stor, is *ixspec.T) func(x, y uint64) int {
	return func(x, y uint64) int {
		xr := offToRec(store, x)
		yr := offToRec(store, y)
		return comp.Compare(xr, yr, is.Cols, is.Cols2)
	}
}

func offToRec(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	return rt.Record(hacks.BStoS(buf))
}
