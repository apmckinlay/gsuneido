// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/comp"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/sortlist"
)

func LoadTable(file string) int {
	table := file
	if strings.HasSuffix(table, ".su") {
		table = file[:len(file)-3]
	}
	defer func() {
		if e := recover(); e != nil {
			panic("load failed: " + table + " " + fmt.Sprint(e))
		}
	}()
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	readLinePrefixed(r, "Suneido dump 2")
	schema := table + " " + readLinePrefixed(r, "====== ")[7:]
	fmt.Println(schema)
	req := compile.ParseRequest("create " + schema).(*compile.Schema)

	store, err := stor.MmapStor("tmp.db", stor.CREATE)
	store.Alloc(1) // don't use offset 0
	ckerr(err)
	list := sortlist.NewUnsorted()
	nrecs := readRecords(r, store, list)
	list.Finish()
	dn := store.Size()
	fmt.Println("data size", dn)
	key := firstShortestKey(req.Indexes)
	for i,ix := range req.Indexes {
		ixcols := ix.Fields
		if req.Indexes[i].Mode != 'k' {
			ixcols = append(ixcols, key...)
		}
		fmt.Println("index", ix, "cols", ixcols)
		if i > 0 {
			list.Sort(mkcmp(store, ixcols))
		}
		before := store.Size()
		bldr := btree.NewFbtreeBuilder(store)
		iter := list.Iter()
		n := 0
		for off := iter(); off != 0; off = iter() {
			bldr.Add(getLeafKey(store, ixcols, off), off)
			n++
		}
		assert.This(n).Is(nrecs)
		fmt.Println("index size", store.Size()-before)
	}
	fmt.Println("total indexes size", store.Size() - dn)
	return nrecs
}

func readLinePrefixed(r *bufio.Reader, pre string) string {
	s, err := r.ReadString('\n') // file header
	ckerr(err)
	if !strings.HasPrefix(s, pre) {
		panic("not a valid dump file")
	}
	return s
}

func readRecords(in *bufio.Reader, store *stor.Stor, list *sortlist.Builder) int {
	nrecs := 0
	intbuf := make([]byte, 4)
	for { // each record
		_, err := io.ReadFull(in, intbuf)
		if err == io.EOF {
			break
		}
		ckerr(err)
		size := int(binary.BigEndian.Uint32(intbuf))
		if size == 0 {
			break
		}
		off, buf := store.Alloc(size)
		_, err = io.ReadFull(in, buf)
		ckerr(err)
		list.Add(off)
		nrecs++
	}
	return nrecs
}

func ckerr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func firstShortestKey(indexes []*compile.Index) []int {
	var key []int
	for _, ix := range indexes {
		if key == nil || len(ix.Fields) < len(key) {
			key = ix.Fields
		}
	}
	return key
}

func getLeafKey(store *stor.Stor, ixspec interface{}, off uint64) string {
	rec := offToRec(store, off)
	ixcols := ixspec.([]int)
	return comp.Key(rt.Record(rec), ixcols)
}

func mkcmp(store *stor.Stor, ixcols []int) func(x, y uint64) int {
	return func(x, y uint64) int {
		xr := offToRec(store, x)
		yr := offToRec(store, y)
		return comp.Compare(xr, yr, ixcols)
	}
}

func offToRec(store *stor.Stor, off uint64) rt.Record {
	buf := store.Data(off)
	return rt.Record(hacks.BStoS(buf))
}
