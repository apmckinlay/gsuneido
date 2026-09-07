// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package stats tracks stats for frequently used (busy) columns
//
//  1. "busy" tracks columns that are frequently use in WHERE clauses.
//     This is appended to a text file either every two hours or on exit.
//  2. "stats" tracks data cardinality (HLL), common values (SS), and quantiles (KLL)
//     This data is collected by compact for the "busy" columns.
//
// "tally" as in BusyTally and StatsTally refers to the data collection.
package stats

import (
	"sort"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hll"
	"github.com/apmckinlay/gsuneido/util/kll"
	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/ss"
)

const topSize = 31  // AI recommendation
const topMin = .005 // AI recommendation

const kllMaxLen = 64 // truncate kll values to this many bytes
const ssMaxLen = 128 // skip ss values larger than this many bytes

const nQuantiles = 51 // i.e. 2% AI recommendation

// StatsTally maps table names to StatsTables
type StatsTally map[string]StatsTable

// StatsTable maps column names to StatsColumns
type StatsTable map[string]*StatsColumn

// StatsColumn collects stats for a column in a table
type StatsColumn struct {
	hll       *hll.HLL            // cardinality
	ss        *ss.Sketch[string]  // top values
	kll       *kll.Sketch[string] // quantiles
	quantiles []string
	top       []*ss.Entry[string]
}

// Add adds a value to the stats
func (sc *StatsColumn) Add(val string) {
	if sc.hll == nil {
		sc.hll = hll.New()
		sc.ss = ss.New[string](topSize)
		sc.kll = kll.New[string]()
	}
	sc.hll.Add(val)
	if len(val) <= ssMaxLen {
		sc.ss.Add(val)
	}
	if len(val) > kllMaxLen {
		val = val[:kllMaxLen]
	}
	sc.kll.Insert(val)
}

// use the Packable interface so we can use RecordBuilder (in compact)
var _ core.Packable = StatsTally{}

// Complete computes quantiles and top values so PackSize can stay read-only.
// Must be called once before PackSize/Pack.
func (stats StatsTally) Complete() {
	for _, st := range stats {
		for _, sc := range st {
			if sc.kll == nil {
				continue
			}
			assert.That(sc.quantiles == nil)
			sc.quantiles = sc.kll.Quantiles(nQuantiles)
			if sc.ss != nil {
				sc.top = sc.ss.TopMin(topMin)
				assert.That(len(sc.top) < 256)
			}
		}
	}
}

// PackSize returns the size of the byte slice that Encode will produce
func (stats StatsTally) PackSize(*uint64) int {
	strlen := func(s string) int {
		assert.That(len(s) < 256)
		return 1 + len(s)
	}
	size := 1 // PackString
	for table, st := range stats {
		cols := make([]string, 0, len(st))
		for col := range st {
			cols = append(cols, col)
		}

		count := 0
		// all cols should have the same count, so we just use the first one
		if kll := st[cols[0]].kll; kll != nil {
			count = kll.Count()
		}
		if count == 0 {
			continue
		}

		size += strlen(table)
		size += 5 // count, len(cols)
		for _, col := range cols {
			sc := st[col]
			assert.That(sc.quantiles != nil)
			size += strlen(col)
			size += 4 // cardinality

			for _, s := range sc.quantiles {
				size += strlen(s)
			}

			if sc.ss != nil {
				size++ // len(top)
				for _, t := range sc.top {
					size += strlen(t.Value)
					size += 4 // count
				}
			} else {
				size++ // len(top)
			}
		}
	}
	return size
}

func (stats StatsTally) PackSize2(*uint64, core.PackStack) int {
	return stats.PackSize(nil)
}

// Pack encodes the stats into a pack.Encoder starting with PackString
func (stats StatsTally) Pack(hash *uint64, e *pack.Encoder) {
	putstr := func(e *pack.Encoder, s string) *pack.Encoder {
		e.Put1(byte(len(s))).PutStr(s)
		return e
	}
	e.Put1(core.PackString) // so we can put it in a record in a table
	for table, st := range stats {
		cols := make([]string, 0, len(st))
		for col := range st {
			cols = append(cols, col)
		}

		count := 0
		// all cols should have the same count, so we just use the first one
		if kll := st[cols[0]].kll; kll != nil {
			count = kll.Count()
		}
		if count == 0 {
			continue
		}

		putstr(e, table).Uint32(uint32(count)).Put1(byte(len(cols)))
		for _, col := range cols {
			sc := st[col]
			putstr(e, col)
			cardinality := sc.hll.Count()
			e.Uint32(uint32(cardinality))

			for _, s := range sc.quantiles {
				putstr(e, s)
			}

			if sc.ss != nil {
				e.Put1(byte(len(sc.top)))
				for _, t := range sc.top {
					putstr(e, t.Value)
					e.Uint32(uint32(t.Count))
				}
			} else {
				e.Put1(0)
			}
		}
	}
}

// Stats holds the statistical information.
type Stats map[string]TableStats

// TableStats holds the column statistics for a single table.
// @immutable
type TableStats struct {
	Count   int // kll count - total values inserted (same for all columns)
	Columns map[string]ColStats
}

// ColStats holds the statistics for a single column.
// @immutable
type ColStats struct {
	Cardinality int
	Quantiles   []string
	Tops        []Top // sorted by frequency, most frequent first
	TailFrac    float64
}

// Top holds a single value frequency entry for the "heavy hitters"
// @immutable
type Top struct {
	Value string
	Frac  float64
}

// UnpackStats decodes the stats from a string into a Stats map
func UnpackStats(s string) Stats {
	getstr := func(d *pack.Decoder) string {
		return d.Get(int(d.Get1()))
	}
	stats := make(Stats)
	d := pack.NewDecoder(s)
	for d.Remaining() > 0 {
		table := getstr(d)
		nrows := int(d.Uint32())
		assert.That(nrows > 0)
		nCols := int(d.Get1())
		columns := make(map[string]ColStats, nCols)
		for range nCols {
			col := getstr(d)
			cardinality := int(d.Uint32())
			quantiles := make([]string, nQuantiles)
			for i := range nQuantiles {
				quantiles[i] = getstr(d)
			}
			nTop := d.Get1()
			top := make([]Top, 0, nTop)
			topCount := 0
			for range nTop {
				val := getstr(d)
				count := d.Uint32()
				top = append(top, Top{
					Value: val,
					Frac:  float64(count) / float64(nrows),
				})
				topCount += int(count)
			}
			columns[col] = ColStats{
				Cardinality: cardinality,
				Quantiles:   quantiles,
				Tops:        top,
				TailFrac:    tailFrac(nrows, topCount),
			}
		}
		stats[table] = TableStats{
			Count:   nrows,
			Columns: columns,
		}
	}
	return stats
}

// tailFrac returns the fraction of rows not accounted for by the top values.
// Top counts are sketch estimates and may sum to more than nrows,
// so the result is clamped to zero.
func tailFrac(nrows, topCount int) float64 {
	return float64(max(0, nrows-topCount)) / float64(max(1, nrows))
}

// PointFrac returns the estimated fraction of rows where the column is
// equal to value. The result is clamped to [0,1].
func (stats Stats) PointFrac(table, col, value string) (float64, bool) {
	ts, ok := stats[table]
	if !ok {
		return 0, false
	}
	cs, ok := ts.Columns[col]
	if !ok {
		return 0, false
	}
	for _, t := range cs.Tops {
		if t.Value == value {
			return clamp01(t.Frac), true
		}
	}
	tailCard := max(cs.Cardinality-len(cs.Tops), 1)
	return clamp01(cs.TailFrac / float64(tailCard)), true
}

// RangeFrac returns the estimated fraction of rows in the specified table
// and column that have a value in the specified range.
// `from` is inclusive, `to` is exclusive.
// The result is clamped to [0,1].
func (stats Stats) RangeFrac(table, col, from, to string) (float64, bool) {
	ts, ok := stats[table]
	if !ok {
		return 0, false
	}
	cs, ok := ts.Columns[col]
	if !ok {
		return 0, false
	}
	if from >= to {
		return 0, true
	}
	frac := fracLess(cs.Quantiles, to) - fracLess(cs.Quantiles, from)
	return clamp01(frac), true
}

// clamp01 clamps a fraction to the range [0,1]
func clamp01(f float64) float64 {
	return max(0, min(f, 1))
}

// fracLess returns the estimated fraction of values less than x,
// estimated from the quantiles.
// Each quantile represents a boundary at rank i/(n-1),
// so the number of quantiles less than x gives the fraction.
func fracLess(quantiles []string, x string) float64 {
	n := len(quantiles)
	if n <= 1 {
		return 0
	}
	// pos = number of quantiles strictly less than x
	pos := sort.Search(n, func(i int) bool { return quantiles[i] >= x })
	frac := float64(pos) / float64(n-1)
	if frac > 1 {
		frac = 1
	}
	return frac
}
