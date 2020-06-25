// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/database/db19/metadata"
)

type UpdateTran struct {
	num int
	//state *DbState
	meta *metadata.Overlay
}

// tranNum should be accessed atomically
var tranNum int64

func NewUpdateTran() *UpdateTran {
	state := GetState()
	info := metadata.NewOverlay(
		state.baseSchema,
		state.baseInfo,
		state.memMeta,
		metadata.NewTableInfoHtbl(0))
	return &UpdateTran{num: int(atomic.AddInt64(&tranNum, 1)), meta: info}
}

func (t *UpdateTran) Commit() {
	UpdateState(func(state *DbState) {
		state.memMeta = t.meta.LayeredOnto(state.memMeta)
	})
	Merge(t.num) //TODO async
}
