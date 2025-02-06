// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/apmckinlay/gsuneido/core"
	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/system"
)

const dumpVersion = "Suneido dump 3\n"
const dumpVersionPrev = "Suneido dump 2\n"
const dumpVersionBase = "Suneido dump"

// DumpDatabase exports the entire database to a file
func DumpDatabase(dbfile, to string) (nTables, nViews int, err error) {
	db, err := OpenDb(dbfile, stor.Read, false)
	if err != nil {
		return 0, 0, err
	}
	defer db.Close()
	return Dump(db, to, "")
}

// Dump exports an entire open database to a file
// If checks as it dumps and could mark the database as corrupted.
func Dump(db *Database, to, publicKey string) (nTables, nViews int, err error) {
	if db.IsCorrupted() {
		return 0, 0, fmt.Errorf("dump not allowed when database is locked")
	}
	defer func() {
		if e := recover(); e != nil {
			if strings.HasPrefix(fmt.Sprint(e), "gopenpgp: ") {
				panic(e)
			}
			db.Corrupt()
			err = fmt.Errorf("dump failed: %v", e)
		}
	}()
	f, w, err := dumpOpen(to, publicKey)
	if err != nil {
		return 0, 0, err
	}
	tmpfile := f.Name()
	defer func() { f.Close(); os.Remove(tmpfile) }()
	nTables, nViews = dump(db, w)
	if err := w.Flush(); err != nil {
		return 0, 0, fmt.Errorf("dump failed: %v", err)
	}
	f.Close()
	if err := system.RenameBak(tmpfile, to); err != nil {
		return 0, 0, fmt.Errorf("dump failed: %v", err)
	}
	return nTables, nViews, nil
}

func dump(db *Database, w WriterPlus) (nTables, nViews int) {
	ics := newIndexCheckers()
	defer ics.finish()
	state := db.Persist()
	nViews = dumpViews(state, w)
	tables := make([]string, 0, 512)
	for sc := range state.Meta.Tables() {
		tables = append(tables, sc.Table)
	}
	sort.Strings(tables)
	for _, table := range tables {
		dumpTable2(db, state, table, true, w, ics)
	}
	return len(tables), nViews
}

// DumpTable exports a table to a binary file
func DumpTable(dbfile, table, to string) (nrecs int, err error) {
	db, err := OpenDb(dbfile, stor.Read, false)
	if err != nil {
		return 0, err
	}
	defer db.Close()
	return DumpDbTable(db, table, to, "")
}

// DumpDbTable exports a table to a binary file from an open database.
// It returns the number of records dumped or panics on error.
// If checks as it dumps and could mark the database as corrupted.
func DumpDbTable(db *Database, table, to, publicKey string) (nrecs int, err error) {
	if db.IsCorrupted() {
		return 0, fmt.Errorf("dump not allowed when database is locked")
	}
	defer func() {
		if e := recover(); e != nil {
			if strings.HasPrefix(fmt.Sprint(e), "gopenpgp: ") {
				panic(e)
			}
			db.Corrupt()
			err = fmt.Errorf("dump failed: %v", e)
		}
	}()
	f, w, err := dumpOpen(to, publicKey)
	if err != nil {
		return 0, err
	}
	tmpfile := f.Name()
	defer func() { f.Close(); os.Remove(tmpfile) }()
	nrecs = dumpDbTable(db, table, w)
	if nrecs == -1 {
		return 0, errors.New("can't find " + table)
	}
	if err := w.Flush(); err != nil {
		return 0, err
	}
	f.Close()
	if err := system.RenameBak(tmpfile, to); err != nil {
		return 0, err
	}
	return nrecs, nil
}

func dumpDbTable(db *Database, table string, w WriterPlus) int {
	ics := newIndexCheckers()
	defer ics.finish()
	state := db.Persist()
	return dumpTable2(db, state, table, false, w, ics)
}

func dumpOpen(to, publicKey string) (*os.File, WriterPlus, error) {
	to = strings.Replace(to, `\`, `/`, -1)
	dir := path.Dir(to)
	f, err := os.CreateTemp(dir, "gs*.tmp")
	if err != nil {
		return nil, nil, err
	}
	var w WriterPlus
	if publicKey == "" {
		w = bufio.NewWriter(f)
	} else {
		w = writerPlus{encryptor(publicKey, f)}
	}
	w.WriteString(dumpVersion)
	return f, w, nil
}

func encryptor(publicKey string, dst io.Writer) io.WriteCloser {
	publicKeyObj, err := crypto.NewKeyFromArmored(publicKey)
	ck(err)
	publicKeyRing, err := crypto.NewKeyRing(publicKeyObj)
	ck(err)
	encryptor, err := publicKeyRing.EncryptStreamWithCompression(dst, nil, nil)
	ck(err)
	return encryptor
}

type WriterPlus interface {
	io.Writer
	WriteString(s string) (n int, err error)
	WriteByte(b byte) error
	Flush() error
}

type writerPlus struct {
	io.WriteCloser
}

func (w writerPlus) WriteString(s string) (n int, err error) {
	return w.Write(hacks.Stobs(s))
}

func (w writerPlus) WriteByte(b byte) error {
	_, err := w.Write(hacks.Btobs(b))
	return err
}

func (w writerPlus) Flush() error {
	return w.WriteCloser.Close()
}

func dumpTable2(db *Database, state *DbState, table string, multi bool,
	w WriterPlus, ics *indexCheckers) int {
	w.WriteString("====== ")
	sc := state.Meta.GetRoSchema(table)
	if sc == nil {
		return -1
	}
	hasdel := sc.HasDeleted()
	schema := sc.DumpString()
	if !multi {
		schema = str.AfterFirst(schema, " ")
	}
	w.WriteString(schema + "\n")
	info := state.Meta.GetRoInfo(table)
	sum := uint64(0)
	nrows := info.Indexes[0].Check(func(off uint64) {
		sum += off                       // addition so order doesn't matter
		rec := OffToRecCk(db.Store, off) // verify data checksums
		if hasdel {
			rec = squeeze(rec, sc.Columns)
		}
		writeInt(w, len(rec))
		w.WriteString(string(rec))
	})
	writeInt(w, 0) // end of table records
	if nrows != info.Nrows {
		panic(fmt.Sprintln("dump", table, sc.Indexes[0].Columns,
			"count", nrows, "should equal info", info.Nrows))
	}
	ics.checkOtherIndexes(sc, info, nrows, sum) // concurrent
	return nrows
}

// squeeze removes deleted fields. It is used by dump and compact.
func squeeze(rec core.Record, cols []string) core.Record {
	var rb core.RecordBuilder
	for i, col := range cols {
		if col != "-" {
			rb.AddRaw(rec.GetRaw(i))
		}
	}
	return rb.Trim().Build()
}

func writeInt(w WriterPlus, n int) {
	assert.That(0 <= n && n <= math.MaxUint32)
	w.WriteByte(byte(n >> 24))
	w.WriteByte(byte(n >> 16))
	w.WriteByte(byte(n >> 8))
	w.WriteByte(byte(n))
}

func dumpViews(state *DbState, w WriterPlus) int {
	w.WriteString("====== views (view_name,view_definition) key(view_name)\n")
	nrecs := 0
	views := make([]string, 0, 32)
	for name := range state.Meta.Views() {
		views = append(views, name)
	}
	sort.Strings(views)
	for _, name := range views {
		def := state.Meta.GetView(name)
		var b core.RecordBuilder
		b.Add(core.SuStr(name))
		b.Add(core.SuStr(def))
		rec := b.Trim().Build()
		writeInt(w, len(rec))
		w.WriteString(string(rec))
		nrecs++
	}
	writeInt(w, 0) // end of table records
	return nrecs
}

// ------------------------------------------------------------------
// Concurrent checking of additional indexes. Also used by compact.

func newIndexCheckers() *indexCheckers {
	ics := indexCheckers{work: make(chan indexCheck, 32), // ???
		stop: make(chan void)}
	nw := options.Nworkers
	ics.wg.Add(nw)
	for range nw {
		go ics.worker()
	}
	return &ics
}

type void struct{}

type indexCheckers struct {
	err    atomic.Value // any
	work   chan indexCheck
	stop   chan void
	wg     sync.WaitGroup
	once   sync.Once
	closed bool
}

type indexCheck struct {
	table  string
	ixcols []string
	index  *index.Overlay
	count  int
	sum    uint64
}

func (ics *indexCheckers) checkOtherIndexes(sc *meta.Schema, info *meta.Info,
	count int, sum uint64) {
	for i := 1; i < len(info.Indexes); i++ {
		select {
		case ics.work <- indexCheck{table: info.Table, ixcols: sc.Indexes[0].Columns,
			index: info.Indexes[i], count: count, sum: sum}:
		case <-ics.stop:
			panic("") // overridden by finish
		}
	}
}

func (ics *indexCheckers) worker() {
	var table string
	defer func() {
		if e := recover(); e != nil {
			ics.err.Store(fmt.Errorf("%s: %v", table, e))
			ics.once.Do(func() { close(ics.stop) }) // notify main thread
		}
		ics.wg.Done()
	}()
	for ic := range ics.work {
		table = ic.table
		CheckOtherIndex(ic.ixcols, ic.index, ic.count, ic.sum)
	}
}

func (ics *indexCheckers) finish() {
	if !ics.closed {
		close(ics.work)
		ics.closed = true
	}
	ics.wg.Wait()
	if err := ics.err.Load(); err != nil {
		panic(err)
	}
}
