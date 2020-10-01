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

func Repair(dbfile string, ec *ErrCorrupt) error {
	fmt.Println("repair")
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
			fmt.Println("truncating", store.Size()-off,
				"=", store.Size(), "-", off)
			fmt.Println("repairing to", t.Format(dtfmt), "from", t0.Format(dtfmt))
			return truncate(dbfile, store, off)
		}
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
			ec = NewErrCorrupt(e)
		}
	}()
	dc := (*dbcheck)(state)
	// If the previous check failed on a certain table,
	// then start by checking that table.
	if table != "" {
		sc := state.meta.GetRoSchema(table)
		dc.checkTable(sc)
	}
	dc.forEachTable(func(sc *meta.Schema) {
		if sc.Table != table {
			dc.checkTable(sc)
		}
	})
	return nil
}

func truncate(dbfile string, store *stor.Stor, off uint64) error {
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
