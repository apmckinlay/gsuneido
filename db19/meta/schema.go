// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math"
	"sort"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hamt"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hash"
	"golang.org/x/exp/slices"
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

func (*Schema) Hash(key string) uint32 {
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
		if ix.BestKey == nil {
			w.Put1(int(ix.Mode)).PutStrs(ix.Columns)
		} else {
			w.Put1(int(ix.Mode) + 1).PutStrs(ix.Columns).PutStrs(ix.BestKey)
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
		for i := 0; i < n; i++ {
			mode := byte(r.Get1())
			columns := r.GetStrs()
			var bestKey []string
			if mode == 'i'+1 || mode == 'u'+1 {
				mode--
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
		}
		ts.Ixspecs(ts.Indexes)
	}
	return &ts
}

// Ixspecs sets up the ixspecs for a table's indexes.
// In most cases idxs will be ts.Indexes.
func (ts *Schema) Ixspecs(idxs []schema.Index) {
	ts.setPrimary()
	ts.setContainsKey()
	short := ts.firstShortestKey() //TODO remove when everything has BestKey
	for i := range idxs {
		ix := &idxs[i]
		key := ix.BestKey
		if key == nil {
			key = short
		}
		switch ix.Mode {
		case 'u':
			cols := difference(key, ix.Columns)
			ix.Ixspec.Fields2 = ts.colsToFlds(cols)
			fallthrough
		case 'k':
			ix.Ixspec.Fields = ts.colsToFlds(ix.Columns)
		case 'i':
			cols := slc.With(ix.Columns, difference(key, ix.Columns)...)
			ix.Ixspec.Fields = ts.colsToFlds(cols)
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
		for j := 0; j < i; j++ {
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

// firstShortestKey is the old way, needed during the transition.
// Going forward it is replaced by BestKey
func (ts *Schema) firstShortestKey() []string {
	hasSpecial := func(cols []string) bool {
		for _, col := range cols {
			if strings.HasSuffix(col, "_lower!") {
				return true
			}
		}
		return false
	}
	usableKey := func(ix *schema.Index) bool {
		return ix.Mode == 'k' && len(ix.Columns) > 0 && !hasSpecial(ix.Columns)
	}
	var key []string
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if usableKey(ix) &&
			(key == nil || len(ix.Columns) < len(key)) {
			key = ix.Columns
		}
	}
	return key
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

// SetBestKey determines the BestKey for each index
// that requires the fewest additional columns.
// This should be done before IxSpecs.
// WARNING: this affects the index entries
// so it should only be called when there is no data.
func (ts *Schema) SetBestKey() {
	for i := range ts.Indexes {
		ix := &ts.Indexes[i]
		if (ix.Mode == 'i' || ix.Mode == 'u') && ix.BestKey == nil {
			best := -1
			bestLen := math.MaxInt
			for j := range ts.Indexes {
				key := &ts.Indexes[j]
				if key.Mode == 'k' {
					n := len(difference(key.Columns, ix.Columns))
					if n < bestLen {
						best = j
						bestLen = n
					}
				}
			}
			ix.BestKey = ts.Indexes[best].Columns
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
	return ts.Columns == nil && ts.Indexes == nil
}

func (ts *Schema) isView() bool {
	return !ts.IsTomb() && ts.Table[0] == '='
}

// isTable returns true if not a view and not a tombstone
func (ts *Schema) isTable() bool {
	return !ts.IsTomb() && !ts.isView()
}
