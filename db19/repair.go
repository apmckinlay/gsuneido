// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
)

const dtfmt = "20060102.150405"

func Repair(dbfile string, err error) error {
	ec, _ := err.(*ErrCorrupt)
	fmt.Println("repair", err, ec.Table())
	store, err := stor.MmapStor(dbfile, stor.READ)
	if err != nil {
		return err
	}
	off := store.Size()
	var state *DbState
	var t0, t time.Time
	for {
		off, state, t = prevState(store, off)
		if t0.IsZero() {
			t0 = t
		}
		if off == 0 {
			return errors.New("repair failed - no valid states found")
		}
		if state == nil {
			continue
		}
		if ec = checkState(state, ec.Table()); ec == nil {
			fmt.Println("good state", off, t.Format(dtfmt))
			fmt.Println("truncating", store.Size()-(off+uint64(stateLen)))
			return truncate(dbfile, store, off)
		}
		fmt.Println("bad state", off, t.Format(dtfmt))
	}
}

func prevState(store *stor.Stor, off uint64) (off2 uint64, state *DbState, t time.Time) {
	off2 = store.LastOffset(off, magic1)
	if off2 == 0 {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			state = nil
		}
	}()
	state, t = ReadState(store, off2)
	return off2, state, t
}

func checkState(state *DbState, table string) (ec *ErrCorrupt) {
	defer func() {
		if e := recover(); e != nil {
			ec = newErrCorrupt(e)
		}
	}()
	tcs := newTableCheckers(state, checkTable)
	defer tcs.finish()
	// If the previous check failed on a certain table,
	// then start by checking that table first (fail faster).
	if table != "" {
		ts := state.meta.GetRoSchema(table)
		tcs.work <- ts
	}
	tcs.state.meta.ForEachSchema(func(ts *meta.Schema) {
		if ts.Table == table {
			return // continue
		}
		select {
		case tcs.work <- ts:
		case <-tcs.stop:
			panic("") // overridden by finish
		}
	})
	return nil
}

func truncate(dbfile string, store *stor.Stor, off uint64) error {
	if off+uint64(stateLen) == store.Size() {
		return nil // nothing to do
	}
	store.Close()
	src, err := os.Open(dbfile)
	if err != nil {
		return err
	}
	dst, err := ioutil.TempFile(".", "gs*.tmp")
	if err != nil {
		return err
	}
	tmpfile := dst.Name()
	_, err = io.CopyN(dst, src, int64(off)+int64(stateLen))
	if err != nil {
		return err
	}
	buf := make([]byte, stor.SmallOffsetLen)
	stor.WriteSmallOffset(buf, off+uint64(stateLen))
	_, err = dst.WriteAt(buf, int64(len(magic)))
	if err != nil {
		return err
	}
	src.Close()
	dst.Close()
	if err = renameBak(tmpfile, dbfile); err != nil {
		return err
	}
	return ensureFlat(dbfile)
}

func renameBak(from string, to string) error {
	err := os.Remove(to + ".bak")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.Rename(to, to+".bak")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	err = os.Rename(from, to)
	if err != nil {
		return err
	}
	return nil
}

func ensureFlat(dbfile string) error {
	// ensure flattened (required by quick check)
	db, err := openDatabase(dbfile, stor.UPDATE, false)
	if err != nil {
		return fmt.Errorf("after rebuild: %w", err)
	}
	defer db.Close()
	db.Persist(true)
	return nil
}
