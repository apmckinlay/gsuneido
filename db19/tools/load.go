// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/apmckinlay/gsuneido/core"
	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	btree3 "github.com/apmckinlay/gsuneido/db19/index/btree3"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/dbms/query"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/errs"
	"github.com/apmckinlay/gsuneido/util/sortlist"
	"github.com/apmckinlay/gsuneido/util/system"
)

type slBuilder = sortlist.Builder[uint64]

type loadJob struct {
	db    *Database
	list  *slBuilder
	ts    *meta.Schema
	nrecs int
	size  int64
}

// LoadDatabase imports a dumped database from a file using a worker pool.
// It returns the number of tables loaded. Errors are fatal.
// It does NOT check foreign key data
// because it assumes the dump was from a valid database.
func LoadDatabase(from, dbfile, privateKey, passphrase string) (
	nTables, nViews int, err error) {
	var errVal atomic.Value // error
	defer func() {
		if e := recover(); e != nil {
			err = errs.From(e)
		}
	}()
	f, r := loadOpen(from, privateKey, passphrase)
	defer f.Close()
	db, tmpfile := tmpdb()
	defer func() { db.Close(); os.Remove(tmpfile) }()

	// start the workers that build the indexes
	var wg sync.WaitGroup
	channel := make(chan *loadJob)
	for range options.Nworkers {
		wg.Go(func() {
			var job *loadJob
			defer func() {
				if e := recover(); e != nil {
					errVal.Store(fmt.Errorf("error loading %s: %v", job.ts.Table, e))
				}
			}()
			for job = range channel {
				loadTable2(job.db, job.ts, job.nrecs, job.size, job.list, false)
			}
		})
	}

	// load the tables
	nTables = 0
	for ; errVal.Load() == nil; nTables++ {
		schema := readLinePrefixed(r, "====== ")
		if schema == "" {
			break
		}
		sch := query.NewAdminParser(schema).Schema()
		ts := &meta.Schema{Schema: sch}
		nrecs, size, list := loadTable1(db, r, schema)
		if strings.HasPrefix(schema, "views ") {
			nViews = nrecs
			nTables--
		} else {
			channel <- &loadJob{db: db, ts: ts,
				nrecs: nrecs, size: size, list: list}
		}
	}
	close(channel)
	wg.Wait()
	if errVal.Load() != nil {
		return 0, 0, errVal.Load().(error)
	}
	trace("SIZE", db.Store.Size())
	db.CheckAllFkeys()
	db.GetState().Write()
	db.Close()
	ck(system.RenameBak(tmpfile, dbfile))
	return nTables, nViews, nil
}

// LoadTable is used by -load <table>
func LoadTable(table, dbfile string) (int, error) {
	var db *Database
	var err error
	if _, err = os.Stat(dbfile); os.IsNotExist(err) {
		db, err = CreateDatabase(dbfile)
	} else {
		db, err = OpenDatabase(dbfile)
	}
	if err != nil {
		return 0, fmt.Errorf("error loading %s: %w", table, err)
	}
	defer db.Close()
	nrecs, err := loadDbTable(table, table+".su", "", "", db)
	if err != nil {
		return 0, err
	}
	db.GetState().Write()
	return nrecs, nil
}

// LoadDbTable is used by dbms.Load / Database.Load
// i.e. on an already open database
func LoadDbTable(table, from, privateKey, passphrase string,
	db *Database) (n int, err error) {
	nrecs, err := loadDbTable(table, from, privateKey, passphrase, db)
	if err != nil {
		return 0, err
	}
	db.Persist() // for safety, not strictly required
	return nrecs, nil
}

func loadDbTable(table, from, privateKey, passphrase string,
	db *Database) (n int, err error) {
	if db.IsCorrupted() {
		return 0, fmt.Errorf("load not allowed when database is locked")
	}
	db.AddExclusive(table)
	defer func() {
		db.EndExclusive(table)
		if e := recover(); e != nil {
			err = fmt.Errorf("error loading %s: %v", table, e)
		}
	}()
	f, r := loadOpen(from, privateKey, passphrase)
	defer f.Close()
	schema := table + " " + readLinePrefixed(r, "====== ")
	ts := tableSchema(db, schema)
	nrecs, size, list := loadTable1(db, r, schema)
	loadTable2(db, ts, nrecs, size, list, true)
	return nrecs, nil
}

func loadOpen(filename, privateKey, passphrase string) (*os.File, *bufio.Reader) {
	f, err := os.Open(filename)
	ck(err)
	var r *bufio.Reader
	if privateKey == "" {
		r = bufio.NewReader(f)
	} else {
		r = bufio.NewReader(decryptor(privateKey, passphrase, f))
	}
	s, err := r.ReadString('\n')
	ck(err)
	if !strings.HasPrefix(s, dumpVersionBase) {
		panic("not a valid dump file")
	}
	if s != dumpVersion && s != dumpVersionPrev {
		panic("invalid dump version")
	}
	return f, r
}

func decryptor(privateKey, passphrase string, src io.Reader) io.Reader {
	privateKeyObj, err := crypto.NewKeyFromArmored(privateKey)
	ck(err)
	privateKeyUnlocked, err := privateKeyObj.Unlock([]byte(passphrase))
	ck(err)
	defer privateKeyUnlocked.ClearPrivateParams()
	privateKeyRing, err := crypto.NewKeyRing(privateKeyUnlocked)
	ck(err)
	decryptor, err := privateKeyRing.DecryptStream(src, nil, 0)
	ck(err)
	return decryptor
}

func tableSchema(db *Database, schem string) *meta.Schema {
	sch := query.NewAdminParser(schem).Schema()
	ts := &meta.Schema{Schema: sch}
	tsCur := db.GetState().Meta.GetRoSchema(ts.Table)
	if tsCur != nil && tsCur.HasFkeyToHere() {
		panic("can't overwrite table that foreign keys point to")
	}
	if ts.HasFkey() {
		panic("can't load single table with foreign keys")
	}
	return ts
}

// loadTable1 reads the data
func loadTable1(db *Database, r *bufio.Reader, schema string) (
	nrows int, size int64, list *sortlist.Builder[uint64]) {
	trace(schema)
	if strings.HasPrefix(schema, "views ") {
		return loadViews(db, r, schema), 0, nil
	}
	store := db.Store
	list = sortlist.NewUnsorted(func(x uint64) bool { return x == 0 })
	nrows, size = readRecords(r, store, list)
	trace("nrecs", nrows, "data size", size)
	list.Finish()
	return nrows, size, list
}

// loadTable2 builds the indexes.
// It is multi-threaded when loading an entire database
func loadTable2(db *Database, ts *meta.Schema,
	nrows int, size int64, list *slBuilder, overwrite bool) {
	indexes := buildIndexes(ts, list, db.Store, nrows)
	ti := meta.NewInfo(ts.Table, indexes, nrows, size)
	if overwrite {
		db.OverwriteTable(ts, ti)
	} else {
		db.AddNewTable(ts, ti)
	}
}

func readLinePrefixed(r *bufio.Reader, pre string) string {
	s, err := r.ReadString('\n')
	if err == io.EOF {
		return ""
	}
	ck(err)
	if !strings.HasPrefix(s, pre) {
		panic("not a valid dump file")
	}
	return s[len(pre):]
}

func readRecords(in *bufio.Reader, store *stor.Stor, list *slBuilder) (
	nrecs int, size int64) {
	intbuf := make([]byte, 4)
	for { // each record
		_, err := io.ReadFull(in, intbuf)
		if err == io.EOF {
			break
		}
		ck(err)
		n := int(binary.BigEndian.Uint32(intbuf))
		if n == 0 {
			break
		}
		off, buf := store.Alloc(n + cksum.Len)
		_, err = io.ReadFull(in, buf[:n])
		ck(err)
		cksum.Update(buf)
		list.Add(off)
		nrecs++
		size += int64(n)
	}
	return nrecs, size
}

func buildIndexes(ts *meta.Schema, list *slBuilder, store *stor.Stor,
	nrecs int) []*index.Overlay {
	i := -1
	defer func() {
		if e := recover(); e != nil {
			index := ""
			if i != -1 {
				index = ts.Indexes[i].String() + " "
			}
			panic(fmt.Sprintf("%s%v", index, e))
		}
	}()
	ts.SetupIndexes()
	ov := make([]*index.Overlay, len(ts.Indexes))
	for i = range ts.Indexes {
		ix := ts.Indexes[i]
		trace(ix)
		if i > 0 || ix.Mode != 'k' {
			list.Sort(MakeLess(store, &ix.Ixspec))
		}
		before := store.Size()
		// only used by load and compact, so hard wired to new btree version
		bldr := btree3.Builder(store)
		iter := list.Iter()
		n := 0
		for off := iter(); off != 0; off = iter() {
			if !bldr.Add(IndexKey(store, &ix.Ixspec, off), off) {
				panic("cannot build index: duplicate value: " +
					ts.Table + " " + ix.String())
			}
			n++
		}
		bt := bldr.Finish()
		if bt.TreeLevels() > 6 {
			log.Println("ERROR: btree levels > 6 in", ts.Table, "nrecs", nrecs, "treeLevels", bt.TreeLevels(), "index", ts.Indexes[i].Columns)
		}
		bt.SetIxspec(&ix.Ixspec)
		ov[i] = index.OverlayFor(bt)
		assert.That(n == nrecs)
		trace("size", store.Size()-before)
	}
	return ov
}

func loadViews(db *Database, in *bufio.Reader, schema string) int {
	assert.That(strings.HasPrefix(schema, "views (view_name,view_definition)"))
	intbuf := make([]byte, 4)
	buf := make([]byte, 32768)
	nrecs := 0
	for { // each record
		_, err := io.ReadFull(in, intbuf)
		if err == io.EOF {
			break
		}
		ck(err)
		n := int(binary.BigEndian.Uint32(intbuf))
		if n == 0 {
			break
		}
		_, err = io.ReadFull(in, buf[:n])
		ck(err)

		rec := core.Record(string(buf[:n]))
		name := rec.GetStr(0)
		def := rec.GetStr(1)
		db.AddView(name, def)

		nrecs++
	}
	return nrecs
}

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func trace(...any) {
	// fmt.Println(args...) // comment out to disable tracing
}
