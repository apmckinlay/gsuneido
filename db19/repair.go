// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/system"
)

const dtfmt = "20060102.150405"

func Repair(dbfile string, err error) (string, error) {
	ec, _ := err.(*errCorrupt)
	store, err := stor.MmapStor(dbfile, stor.Read)
	if err != nil {
		return "", err
	}
	defer store.Close(true)
	r := repair{dbfile: dbfile, store: store, ec: ec}
	r.oldver = bufHasPrefix(store.Data(0), magicPrev)
	_, off, state := r.search()
	if off == 0 {
		return "", errors.New("repair failed - no valid states found")
	}
	truncated, err := r.fix(off)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprint("good state ", trace.Number(off+uint64(stateLen)), " ",
		time.UnixMilli(state.Asof).Format(dtfmt),
		" truncated ", trace.Number(truncated))
	return msg, nil
}

type repair struct {
	dbfile string
	oldver bool
	store  *stor.Stor
	ec     *errCorrupt
}

func (r *repair) search() (int, uint64, *DbState) {
	good := -1
	i := 0
	var prev int
	last := false
	var state *DbState
	scnr := newScanner(r.store)
	defer scnr.close()
	var offsets []uint64
	// search backwards for a good state
	// in increasing jumps to reduce the number of states checked
	for skip := 1; ; skip *= 2 {
		var done bool
		offsets, done = scnr.getUpTo(i)
		if done {
			i = len(offsets) - 1
			if i == prev {
				return 0, 0, nil // no more states
			}
			last = true
		}
		off := offsets[i]
		if state = r.check(i, off); state != nil {
			fmt.Println("+", i, "good")
			good = i
			scnr.stop()
			break
		}
		fmt.Println("+", i, r.ec)
		if last {
			return 0, 0, nil // no more states
		}
		prev = i
		i += skip
	}

	// binary search for where good changes to bad
	// This assumes that good and bad states are NOT mixed.
	// i.e. that it's good,good,good,bad,bad,bad
	// If they are mixed e.g. good,bad,good,bad
	// then we won't necessarily find the most recent good state.
	lo := prev
	hi := good
	for lo < hi-1 {
		mid := lo + (hi-lo)/2
		off := offsets[mid]
		if s := r.check(mid, off); s != nil {
			fmt.Println("-", mid, "good")
			hi = mid
			state = s
		} else {
			fmt.Println("-", mid, r.ec)
			lo = mid
		}
	}
	return hi, offsets[hi], state
}

func (r *repair) check(i int, off uint64) (state *DbState) {
	if i < 3 {
		// 	r.ec = &errCorrupt{err: "injected error"}
		// 	return nil
	}
	state = getState(r.store, off)
	if state == nil {
		r.ec = &errCorrupt{err: "ReadState failed"}
		return nil
	}
	ec := checkState(state, checkTable, r.ec.Table(), r.ec.Ixcols())
	if ec != nil {
		r.ec = ec
		return nil
	}
	return state
}

func getState(store *stor.Stor, off uint64) (state *DbState) {
	defer func() {
		if e := recover(); e != nil {
			state = nil
		}
	}()
	return ReadState(store, off)
}

func (r *repair) fix(off uint64) (uint64, error) {
	storeSize := r.store.Size()
	r.store.Close(true)
	size := off + uint64(stateLen)
	if size == storeSize || size == storeSize-tailSize {
		return 0, r.fixHead(size) // truncate 0
	}
	tmpfile, err := r.copySize(size)
	if err != nil {
		return 0, err
	}
	if !r.oldver {
		size += tailSize
	}
	return storeSize - size, system.RenameBak(tmpfile, r.dbfile)
}

func (r *repair) fixHead(size uint64) error {
	f, err := os.OpenFile(r.dbfile, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	if r.oldver {
		buf := make([]byte, stor.SmallOffsetLen)
		stor.WriteSmallOffset(buf, size)
		_, err = f.WriteAt(buf, int64(len(magic)))
	} else { // new version
		// add or overwrite shutdown marker
		_, err = f.WriteAt([]byte(shutdown), int64(size))
	}
	return err
}

func (r *repair) copySize(size uint64) (string, error) {
	src, err := os.Open(r.dbfile)
	if err != nil {
		return "", err
	}
	defer src.Close()
	dst, err := os.CreateTemp(".", "gs*.tmp")
	if err != nil {
		return "", err
	}
	defer dst.Close()
	tmpfile := dst.Name()
	_, err = io.CopyN(dst, src, int64(size))
	if err != nil {
		return "", err
	}
	if r.oldver {
		buf := make([]byte, stor.SmallOffsetLen)
		stor.WriteSmallOffset(buf, size)
		_, err = dst.WriteAt(buf, int64(len(magic)))
	} else { // new version
		_, err = dst.Write([]byte(shutdown))
	}
	return tmpfile, err
}

//-------------------------------------------------------------------

// scanner is used by repair to search for states (magic1 ... magic2)
// It uses a goroutine to scan ahead.
type scanner struct {
	offsets []uint64 // guarded by lock
	cond    sync.Cond
	lock    sync.Mutex
	wg      sync.WaitGroup
	done    uint32 // should be accessed atomically
}

func newScanner(store *stor.Stor) *scanner {
	var s scanner
	s.cond.L = &s.lock
	s.wg.Add(1)
	go s.scanner(store)
	return &s
}

func (s *scanner) getUpTo(i int) ([]uint64, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for len(s.offsets) <= i {
		if atomic.LoadUint32(&s.done) == 1 {
			return s.offsets, true
		}
		s.cond.Wait()
	}
	return s.offsets, false
}

func (s *scanner) scanner(store *stor.Stor) {
	off := store.Size()
	for {
		off = store.LastOffset(off, magic1, &s.done)
		if off == 0 { // finished (either at end or stopped)
			break
		}
		buf := store.Data(off)
		if string(buf[magic2at:magic2at+len(magic2)]) != magic2 {
			continue
		}
		s.lock.Lock()
		s.offsets = append(s.offsets, off)
		s.lock.Unlock()
		s.cond.Signal()
	}
	s.stop()
	s.cond.Signal()
	s.wg.Done()
}

func (s *scanner) stop() {
	atomic.StoreUint32(&s.done, 1)
}

func (s *scanner) close() {
	s.stop()
	s.wg.Wait()
}

//-------------------------------------------------------------------

func PrintStates(dbfile string, check bool) {
	store, err := stor.MmapStor(dbfile, stor.Read)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer store.Close(false)

	ec := &errCorrupt{}
	scnr := newScanner(store)
	for i := 0; ; i++ {
		offsets, done := scnr.getUpTo(i)
		if done {
			break
		}
		state := getState(store, offsets[i])
		if state == nil {
			fmt.Println(i, "read state failed")
			continue
		}
		msg := ""
		if check {
			msg = "good"
			ec = checkState(state, checkTable, ec.Table(), ec.Ixcols())
			if ec != nil {
				msg = ec.Error()
			}
		}
		fmt.Println(i, trace.Number(offsets[i]), time.UnixMilli(state.Asof).Format(dtfmt), msg)
	}
	scnr.close()
}
