// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"slices"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hash"
)

// Note: views are stored with the name in Schema.Table prefixed by '='
// and the definition in Schema.Columns[0]

type SchemaHamt = hamt.Hamt[string, *Schema]

type Schema struct {
	schema.Schema
	// lastMod must be set to Meta.infoClock on new or modified items.
	// It is used for persist meta chaining/flattening.
	lastMod int
	// created is used to avoid tombstones (and persisting them)
	// for temporary tables (e.g. from tests)
	created int
}

func (ts *Schema) Key() string {
	return ts.Table
}

func (*Schema) Hash(key string) uint64 {
	return hash.String(key)
}

func (ts *Schema) LastMod() int {
	return ts.lastMod
}

func (ts *Schema) SetLastMod(mod int) {
	ts.lastMod = mod
}

func (ts *Schema) StorSize() int {
	size := stor.LenStr(ts.Table) +
		stor.LenStrs(ts.Columns) + stor.LenStrs(ts.Derived) + 1
	for i := range ts.Indexes {
		idx := ts.Indexes[i]
		size += 1 + stor.LenStrs(idx.Columns) +
			stor.LenStr(idx.Fk.Table) + 1 + stor.LenStrs(idx.Fk.Columns)
		if idx.BestKey != nil {
			size += stor.LenStrs(idx.BestKey)
		}
	}
	return size
}

func (ts *Schema) Write(w *stor.Writer) {
	w.PutStr(ts.Table)
	w.PutStrs(ts.Columns)
	w.PutStrs(ts.Derived)
	w.Put1(len(ts.Indexes))
	for _, ix := range ts.Indexes {
		if ix.Fk.Table == "" && len(ix.Fk.Columns) != 0 {
			// TEMPORARY - old bug filled in Columns when it shouldn't
			ix.Fk.Columns = nil
			fmt.Println("Write Schema: Fk.Table empty but Fk.Columns:", ix.Fk.Columns)
		}
		if ix.Mode == 'k' {
			w.Put1(int(ix.Mode)).PutStrs(ix.Columns)
		} else {
			assert.That(ix.BestKey != nil)
			w.Put1(int(ix.Mode)).PutStrs(ix.Columns).PutStrs(ix.BestKey)
		}
		w.PutStr(ix.Fk.Table).Put1(int(ix.Fk.Mode)).PutStrs(ix.Fk.Columns)
	}
}

func ReadSchema(_ *stor.Stor, r *stor.Reader) *Schema {
	ts := Schema{}
	ts.Table = r.GetStr()
	ts.Columns = r.GetStrs()
	ts.Derived = r.GetStrs()
	if n := r.Get1(); n > 0 {
		ts.Indexes = make([]schema.Index, n)
		for i := range n {
			mode := byte(r.Get1())
			columns := r.GetStrs()
			var bestKey []string
			if mode != 'k' {
				bestKey = r.GetStrs()
			}
			ts.Indexes[i] = schema.Index{
				Mode:    mode,
				Columns: columns,
				BestKey: bestKey,
				Fk: schema.Fkey{
					Table:   r.GetStr(),
					Mode:    byte(r.Get1()),
					Columns: r.GetStrs()},
			}
			if ts.Indexes[i].Fk.Table == "" && len(ts.Indexes[i].Fk.Columns) != 0 {
				// TEMPORARY - old bug filled in Columns when it shouldn't
				ts.Indexes[i].Fk.Columns = nil
			}
		}
		ts.Ixspecs(0)
	}
	return &ts
}

// Ixspecs sets up the ixspecs for a table's indexes.
func (ts *Schema) Ixspecs(nold int) {
	ts.setPrimary()
	ts.setContainsKey()
	for i := nold; i < len(ts.Indexes); i++ {
		ix := &ts.Indexes[i]
		assert.That(ix.Mode == 'k' || ix.BestKey != nil)
		switch ix.Mode {
		case 'u':
			cols := difference(ix.BestKey, ix.Columns)
			ix.Ixspec.Fields2 = ts.colsToFlds(cols)
			fallthrough
		case 'k':
			ix.Fields = ix.Columns
			ix.Ixspec.Fields = ts.colsToFlds(ix.Columns)
		case 'i':
			ix.Fields = slc.With(ix.Columns, difference(ix.BestKey, ix.Columns)...)
			ix.Ixspec.Fields = ts.colsToFlds(ix.Fields)
		default:
			panic("Ixspecs invalid mode")
		}
	}
}

func (ts *Schema) setPrimary() {
	keys := make([]*schema.Index, 0, 4)
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if ix.Mode == 'k' {
			ix.Primary = false
			keys = append(keys, ix)
		}
	}
	lt := func(i, j int) bool {
		icols := keys[i].Columns
		jcols := keys[j].Columns
		ni := len(icols) * 2
		nj := len(jcols) * 2
		// choose key(foo_lower!) before key(foo)
		if ni == 2 && strings.HasSuffix(icols[0], "_lower!") {
			ni--
		}
		if nj == 2 && strings.HasSuffix(jcols[0], "_lower!") {
			nj--
		}
		return ni < nj
	}
	sort.SliceStable(keys, lt)
	keys[0].Primary = true
	if len(keys[0].Columns) == 0 {
		return
	}
outer:
	for i := 1; i < len(keys); i++ {
		for j := range i {
			if keys[j].Primary && subset(keys[i].Columns, keys[j].Columns) {
				continue outer
			}
		}
		keys[i].Primary = true
	}
}

func (ts *Schema) setContainsKey() {
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if ix.Mode == 'u' {
			for j := range ts.Indexes {
				key := &ts.Indexes[j]
				if key.Mode == 'k' && subset(ix.Columns, key.Columns) {
					ix.ContainsKey = true
					break
				}
			}
		}
	}
}

func (ts *Schema) colsToFlds(cols []string) []int {
	flds := make([]int, len(cols))
	for i, col := range cols {
		c := slices.Index(ts.Columns, col)
		if strings.HasSuffix(col, "_lower!") {
			if c = slices.Index(ts.Columns, col[:len(col)-7]); c != -1 {
				c = -c - 2
			}
		}
		assert.That(c != -1)
		flds[i] = c
	}
	return flds
}

func (ts *Schema) SetupIndexes() {
	ts.SetupNewIndexes(0)
}

// SetupNewIndexes sets BestKey and creates Ixspecs.
func (ts *Schema) SetupNewIndexes(nold int) []schema.Index {
	ts.SetBestKeys(nold)
	ts.Ixspecs(nold)
	return ts.Indexes[nold:]
}

// SetBestKeys determines the BestKey for each index
// that requires the fewest additional columns.
// This should be done before IxSpecs.
// WARNING: this affects the stored index entries
// so it should only be used for empty indexes.
func (ts *Schema) SetBestKeys(nold int) {
	for i := nold; i < len(ts.Indexes); i++ {
		ix := &ts.Indexes[i]
		if (ix.Mode == 'i' || ix.Mode == 'u') && ix.BestKey == nil {
			var best *schema.Index
			bestLen := math.MaxInt
			for j := range ts.Indexes {
				key := &ts.Indexes[j]
				if key.Mode == 'k' {
					n := len(difference(key.Columns, ix.Columns))
					if n < bestLen {
						best = key
						bestLen = n
					}
				}
			}
			assert.That(ix.BestKey == nil) // should never overwrite
			ix.BestKey = best.Columns
		}
	}
}

// difference is like set.Difference
// but it takes _lower! into account, which makes it asymmetric
func difference(key, idx []string) []string {
	z := make([]string, 0, len(key))
outer:
	for _, ke := range key {
		ket := strings.TrimSuffix(ke, "_lower!")
		for _, ie := range idx {
			if ie == ke || ie == ket {
				continue outer
			}
		}
		z = append(z, ket)
	}
	return z
}

// subset is like set.Subset
// but it takes _lower! into account, which makes it asymmetric
func subset(idx, key []string) bool {
outer:
	for _, ke := range key {
		ket := strings.TrimSuffix(ke, "_lower!")
		for _, ie := range idx {
			if ie == ke || ie == ket {
				continue outer
			}
		}
		return false
	}
	return true
}

func (m *Meta) newSchemaTomb(table string) *Schema {
	return &Schema{Schema: schema.Schema{Table: table}}
}

func (m *Meta) newSchemaView(name, def string) *Schema {
	return &Schema{Schema: schema.Schema{Table: "=" + name, Columns: []string{def}}}
}

func (ts *Schema) IsTomb() bool {
	return len(ts.Columns) == 0 && len(ts.Indexes) == 0
}

func (ts *Schema) isView() bool {
	return !ts.IsTomb() && ts.Table[0] == '='
}

// isTable returns true if not a view and not a tombstone
func (ts *Schema) isTable() bool {
	return !ts.IsTomb() && !ts.isView()
}
