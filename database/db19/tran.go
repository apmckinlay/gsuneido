// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
)

type tran struct {
	num   int
	meta  *meta.Overlay
	store *stor.Stor
}

type ReadTran struct {
	tran
}

func NewReadTran() *ReadTran {
	state := GetState()
	return &ReadTran{tran: tran{num: int(atomic.AddInt64(&tranNum, 1)),
		meta: state.meta, store: state.store}}
}

type UpdateTran struct {
	tran
}

// tranNum should be accessed atomically
var tranNum int64

func NewUpdateTran() *UpdateTran {
	state := GetState()
	meta := state.meta.NewOverlay()
	return &UpdateTran{tran: tran{num: int(atomic.AddInt64(&tranNum, 1)),
		meta: meta, store: state.store}}
}

func (t *UpdateTran) Commit() int {
	UpdateState(func(state *DbState) {
		state.meta = t.meta.LayeredOnto(state.meta)
	})
	return t.num
}
