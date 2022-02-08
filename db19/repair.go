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

func Repair(dbfile string, err error) (string, error) {
	ec, _ := err.(*ErrCorrupt)
	fmt.Println("repair:", err, ec.Table())
	store, err := stor.MmapStor(dbfile, stor.READ)
	if err != nil {
		return "", err
	}
	defer store.Close()
	off := store.Size()
	var state *DbState
	var t0, t time.Time
	for {
		off, state, t = prevState(store, off)
		if t0.IsZero() {
			t0 = t
		}
		if off == 0 {
			return "", errors.New("repair failed - no valid states found")
		}
		if state == nil {
			continue
		}
		if ec = checkState(state, ec.Table()); ec == nil {
			msg := fmt.Sprintln("good state", off, t.Format(dtfmt)) +
				fmt.Sprintln("truncating", store.Size()-(off+uint64(stateLen)))
			return msg, truncate(dbfile, store, off)
		}
		fmt.Println("bad state", off, t.Format(dtfmt), ec)
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
		tcs.work <- table
	}
	tcs.state.Meta.ForEachSchema(func(ts *meta.Schema) {
		if ts.Table == table {
			return // continue
		}
		select {
		case tcs.work <- ts.Table:
		case <-tcs.stop:
			panic("") // overridden by finish
		}
	})
	return nil
}

func truncate(dbfile string, store *stor.Stor, off uint64) error {
	store.Close()
	size := off + uint64(stateLen)
	if size == store.Size() {
		return fixHeader(dbfile, size)
	}
	tmpfile, err := truncate2(dbfile, size)
	if err != nil {
		return err
	}
	return RenameBak(tmpfile, dbfile)
}

func fixHeader(dbfile string, size uint64) error {
	f, err := os.OpenFile(dbfile, os.O_WRONLY, 0)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.WriteAt([]byte(magic), 0)
	if err != nil {
		return err
	}
	buf := make([]byte, stor.SmallOffsetLen)
	stor.WriteSmallOffset(buf, size)
	_, err = f.WriteAt(buf, int64(len(magic)))
	return err
}

func truncate2(dbfile string, size uint64) (string, error) {
		src, err := os.Open(dbfile)
		if err != nil {
			return "", err
		}
		defer src.Close()
		dst, err := ioutil.TempFile(".", "gs*.tmp")
		if err != nil {
			return "", err
		}
		defer dst.Close()
		tmpfile := dst.Name()
		_, err = io.CopyN(dst, src, int64(size))
		if err != nil {
			return "", err
		}
		buf := make([]byte, stor.SmallOffsetLen)
		stor.WriteSmallOffset(buf, size)
		_, err = dst.WriteAt(buf, int64(len(magic)))
		if err != nil {
			return "", err
		}
		return tmpfile, nil
	}

func RenameBak(from string, to string) error { //TODO move to util
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
