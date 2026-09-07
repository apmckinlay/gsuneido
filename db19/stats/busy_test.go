// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stats

import (
	"os"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestHotcolsAdd(t *testing.T) {
	var hc BusyTally
	hc.Add("a", 1000)
	hc.Add("b", 2000)
	hc.Add("a", 3000)
	assert.T(t).This(hc.busy.Count()).Is(6000)
	assert.T(t).This(hc.busy.Len()).Is(2)
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
	hc.Save()

	// Save should reset the sketch
	assert.That(hc.busy.Len() == 0)

	// Load should return top entries grouped by table
	result := LoadBusy()
	assert.That(len(result) > 0)
	assert.That(len(result["t"]) > 0)
	assert.That(result["t"][0] == "x")
}

func TestHotcolsLoadEmpty(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(orig)

	result := LoadBusy()
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
	hc.Save()

	// Second batch (same file, appended)
	hc.Add("t.a", 5000)
	hc.Add("t.c", 10000)
	hc.Save()

	result := LoadBusy()
	// "a" and "c" both have total count 10, so they should be in the top 2
	assert.That(len(result["t"]) >= 2)
	assert.That(result["t"][0] == "a" || result["t"][0] == "c")
	assert.That(result["t"][1] == "a" || result["t"][1] == "c")
}
