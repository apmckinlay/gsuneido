// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hot

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestStatsColumnAdd(t *testing.T) {
	var sc StatsColumn
	assert.T(t).That(sc.hll == nil && sc.ss == nil && sc.kll == nil)
	sc.Add("a")
	sc.Add("b")
	sc.Add("a")
	assert.T(t).That(sc.hll != nil && sc.ss != nil && sc.kll != nil)
	assert.T(t).This(sc.ss.Count()).Is(3)
	assert.T(t).This(sc.kll.Count()).Is(3)
	// hll is approximate, small counts can be off slightly
	cardinality := int(sc.hll.Count())
	assert.T(t).That(1 <= cardinality && cardinality <= 3)
}

func TestStatsPack(t *testing.T) {
	assert := assert.T(t)
	stats := StatsTally{
		"customer": StatsTable{
			"name": &StatsColumn{},
			"city": &StatsColumn{},
		},
		"product": StatsTable{
			"code": &StatsColumn{},
		},
	}
	name := stats["customer"]["name"]
	for range 100 {
		name.Add("joe")
	}
	for range 50 {
		name.Add("sue")
	}
	for range 10 {
		name.Add("bob")
	}

	city := stats["customer"]["city"]
	for _, c := range []string{"a", "b", "c", "d", "e"} {
		for range 32 {
			city.Add(c)
		}
	}
	assert.This(city.kll.Count()).Is(name.kll.Count())

	stats["product"]["code"].Add("x1")

	stats.Complete()
	data := core.Pack(stats)
	assert.This(data[0]).Is(core.PackString)
	// DecodeStats expects the data without the leading PackString tag
	// (in production Record.GetStr strips it)
	decoded := UnpackStats(data[1:])

	assert.This(len(decoded)).Is(2)
	cus := decoded["customer"]
	assert.This(cus.Count).Is(160)
	assert.This(len(cus.Columns)).Is(2)
	prod := decoded["product"]
	assert.This(prod.Count).Is(1)
	assert.This(len(prod.Columns)).Is(1)

	nameStats, _ := cus.Columns["name"]
	assert.That(2 <= nameStats.Cardinality && nameStats.Cardinality <= 4)
	assert.This(len(nameStats.Tops)).Is(3)
	assert.This(nameStats.Tops[0].Value).Is("joe")
	assert.This(nameStats.Tops[0].Frac).Is(.625)
	assert.This(nameStats.Tops[1].Value).Is("sue")
	assert.This(nameStats.Tops[1].Frac).Is(.3125)
	assert.This(nameStats.Tops[2].Value).Is("bob")
	assert.This(nameStats.Tops[2].Frac).Is(.0625)

	cityStats, _ := cus.Columns["city"]
	assert.That(4 <= cityStats.Cardinality && cityStats.Cardinality <= 6)
	// kll is exact for small counts
	assert.This(len(cityStats.Quantiles)).Is(51)
	assert.This(cityStats.Quantiles[0]).Is("a")
	assert.This(cityStats.Quantiles[11]).Is("b")
	assert.This(cityStats.Quantiles[21]).Is("c")
	assert.This(cityStats.Quantiles[31]).Is("d")
	assert.This(cityStats.Quantiles[41]).Is("e")

	codeStats, _ := prod.Columns["code"]
	assert.This(codeStats.Cardinality).Is(1)
	assert.This(len(codeStats.Quantiles)).Is(51)
	assert.This(codeStats.Quantiles[0]).Is("x1")
	assert.This(len(codeStats.Tops)).Is(1)
	assert.This(codeStats.Tops[0].Value).Is("x1")
	assert.This(codeStats.Tops[0].Frac).Is(1)
}

func TestStatsSize(t *testing.T) {
	stats := StatsTally{
		"customer": StatsTable{
			"name": &StatsColumn{},
			"city": &StatsColumn{},
		},
		"product": StatsTable{
			"code": &StatsColumn{},
		},
		// table with nothing added should be skipped
		"empty": StatsTable{"x": &StatsColumn{}},
	}
	name := stats["customer"]["name"]
	for range 100 {
		name.Add("joe")
	}
	for range 50 {
		name.Add("sue")
	}
	city := stats["customer"]["city"]
	for _, c := range []string{"a", "b", "c", "d", "e"} {
		city.Add(c)
	}
	stats["product"]["code"].Add("x1")

	stats.Complete()
	assert.T(t).This(stats.PackSize(nil)).Is(len(core.Pack(stats)))
	assert.T(t).This(StatsTally{}.PackSize(nil)).Is(1)
}

func TestStatsPackUnaddedColumn(t *testing.T) {
	// a table with nothing added should be skipped entirely
	stats := StatsTally{"t": StatsTable{"x": &StatsColumn{}}}
	stats.Complete()
	data := core.Pack(stats)
	decoded := UnpackStats(data[1:])
	assert.T(t).This(len(decoded)).Is(0)
}

func TestStatsPackEmpty(t *testing.T) {
	empty := StatsTally{}
	empty.Complete()
	data := core.Pack(empty)
	assert.T(t).This(data).Is(string(rune(core.PackString)))
	decoded := UnpackStats(data[1:])
	assert.T(t).This(len(decoded)).Is(0)
}

func TestTailFrac(t *testing.T) {
	assert := assert.T(t)
	test := func(expected float64, nrows, topCount int) {
		t.Helper()
		frac := tailFrac(nrows, topCount)
		assert.This(frac).Is(expected)
		assert.That(frac >= 0)
	}
	test(1, 100, 0)
	test(.2, 100, 80)
	test(0, 100, 100)
	test(0, 100, 110) // approximate counts can overcount
	test(0, 0, 0)
}

func TestPointFrac(t *testing.T) {
	stats := Stats{
		"mytable": {
			Count: 100,
			Columns: map[string]ColStats{
				"status": {
					Cardinality: 5,
					Tops: []Top{
						{Value: "a", Frac: 0.5},
						{Value: "b", Frac: 0.3},
					},
					TailFrac: 0.2,
				},
				"type": {
					Cardinality: 3,
					Tops: []Top{
						{Value: "x", Frac: 0.6},
					},
					TailFrac: 0.4,
				},
			},
		},
	}
	assert := assert.T(t)

	frac, ok := stats.PointFrac("mytable", "status", "a")
	assert.That(ok)
	assert.This(frac).Is(float64(float32(0.5)))

	frac, ok = stats.PointFrac("mytable", "status", "b")
	assert.That(ok)
	assert.This(frac).Is(float64(float32(0.3)))

	frac, ok = stats.PointFrac("mytable", "status", "c")
	assert.That(ok)
	assert.This(frac).Is(float64(float32(0.2)) / float64(3))

	frac, ok = stats.PointFrac("mytable", "type", "x")
	assert.That(ok)
	assert.This(frac).Is(float64(float32(0.6)))

	frac, ok = stats.PointFrac("mytable", "type", "y")
	assert.That(ok)
	assert.This(frac).Is(float64(float32(0.4)) / float64(2))

	_, ok = stats.PointFrac("notable", "status", "a")
	assert.That(!ok)

	_, ok = stats.PointFrac("mytable", "nocol", "a")
	assert.That(!ok)
}

func TestPointFracTailCardGuard(t *testing.T) {
	// Cardinality is an HLL estimate and can be <= len(Tops),
	// so tailCard must be guarded against <= 0 (division by zero / negative).
	stats := Stats{
		"mytable": {
			Count: 100,
			Columns: map[string]ColStats{
				"status": {
					Cardinality: 2,
					Tops: []Top{
						{Value: "a", Frac: 0.5},
						{Value: "b", Frac: 0.5},
					},
					TailFrac: 0.0,
				},
				"under": {
					Cardinality: 1,
					Tops: []Top{
						{Value: "a", Frac: 0.6},
						{Value: "b", Frac: 0.4},
					},
					TailFrac: 0.0,
				},
			},
		},
	}
	assert := assert.T(t)

	frac, ok := stats.PointFrac("mytable", "status", "c")
	assert.That(ok)
	assert.This(frac).Is(0.0)

	frac, ok = stats.PointFrac("mytable", "under", "c")
	assert.That(ok)
	assert.This(frac).Is(0.0)
}

func TestStatsRangeFrac(t *testing.T) {
	assert := assert.T(t)
	stats := StatsTally{
		"customer": StatsTable{
			"city": &StatsColumn{},
		},
	}
	city := stats["customer"]["city"]
	for i := range 51 {
		city.Add(fmt.Sprintf("%02d", i))
	}
	// kll is exact for small counts
	stats.Complete()
	data := core.Pack(stats)
	decoded := UnpackStats(data[1:])

	c, _ := decoded["customer"].Columns["city"]
	fmt.Println(c.Quantiles[:10])
	fmt.Println(c.Quantiles[10:20])
	fmt.Println(c.Quantiles[20:30])
	fmt.Println(c.Quantiles[30:40])
	fmt.Println(c.Quantiles[40:50])
	fmt.Println(c.Quantiles[50:])

	frac, ok := decoded.RangeFrac("customer", "city", "", "25")
	assert.That(ok)
	assert.This(frac).Is(.5)

	frac, ok = decoded.RangeFrac("customer", "city", "25", "99")
	assert.That(ok)
	assert.This(frac).Is(.5)

	frac, ok = decoded.RangeFrac("customer", "city", "25", "26")
	assert.That(ok)
	assert.This(frac).Is(.02)

	frac, ok = decoded.RangeFrac("customer", "city", "", "99")
	assert.That(ok)
	assert.This(frac).Is(1)
}
