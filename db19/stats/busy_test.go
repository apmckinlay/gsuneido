// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stats

import (
	"os"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestBusyAdd(t *testing.T) {
	var bt BusyTally
	bt.Add("a", 1000)
	bt.Add("b", 2000)
	bt.Add("a", 3000)
	assert.T(t).This(bt.busy.Count()).Is(6000)
	assert.T(t).This(bt.busy.Len()).Is(2)
}

func TestBusyScans(t *testing.T) {
	t.Chdir(t.TempDir())
	var bt BusyTally
	getIndexes := func(_ string) [][]string {
		return [][]string{{"x"}, {"y"}}
	}
	bt.Add("t.x", 10000)
	bt.AddScan("t", 0, 20)
	bt.AddScan("t", 1, 15)
	bt.Save(getIndexes)
	bt.AddScan("t", 0, 30)
	bt.Save(getIndexes)
	busy, scans := LoadBusy()
	assert.T(t).This(busy["t"]).Is([]string{"x"})
	assert.T(t).This(scans["t"]["x"]).Is(50)
	assert.T(t).This(scans["t"]["y"]).Is(15)
	_, err := os.Stat(busyFile)
	assert.T(t).That(os.IsNotExist(err))
}

func TestBusyResetScans(t *testing.T) {
	t.Chdir(t.TempDir())
	var bt BusyTally
	getIndexes := func(table string) [][]string {
		if table == "t" {
			return [][]string{{"x"}}
		}
		return [][]string{{"y"}}
	}
	bt.Add("t.x", 10000)
	bt.AddScan("t", 0, 20)
	bt.AddScan("u", 0, 10)
	bt.ResetScans("t")
	bt.Save(getIndexes)
	busy, scans := LoadBusy()
	assert.T(t).This(scans["t"]).Is(nil)
	assert.T(t).This(scans["u"]["y"]).Is(10)
	assert.T(t).This(busy["t"]).Is([]string{"x"})
}

func TestHotcolsSaveLoad(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	var hc BusyTally
	hc.Add("t.x", 10000)
	hc.Add("t.y", 5000)
	hc.Add("t.z", 1000)
	hc.Save(nil)

	// Save should reset the sketch
	assert.That(hc.busy.Len() == 0)

	// Load should return top entries grouped by table
	result, _ := LoadBusy()
	assert.That(len(result) > 0)
	assert.That(len(result["t"]) > 0)
	assert.That(result["t"][0] == "x")
}

func TestHotcolsLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	result, _ := LoadBusy()
	assert.That(result == nil)
}

func TestHotcolsLoadAggregates(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	// First batch
	var hc BusyTally
	hc.Add("t.a", 5000)
	hc.Add("t.b", 3000)
	hc.Save(nil)

	// Second batch (same file, appended)
	hc.Add("t.a", 5000)
	hc.Add("t.c", 10000)
	hc.Save(nil)

	result, _ := LoadBusy()
	// "a" and "c" both have total count 10, so they should be in the top 2
	assert.That(len(result["t"]) >= 2)
	assert.That(result["t"][0] == "a" || result["t"][0] == "c")
	assert.That(result["t"][1] == "a" || result["t"][1] == "c")
}
