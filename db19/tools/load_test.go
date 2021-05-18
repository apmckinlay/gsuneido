// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/db19"
)

func TestLoadTable(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	defer os.Remove("tmp.db")
	os.Remove("tmp.db")
	n := LoadTable("stdlib", "tmp.db")
	fmt.Println("loaded", n, "records in", time.Since(t).Round(time.Millisecond))
	ck(CheckDatabase("tmp.db"))
}

func TestLoadDatabase(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	defer os.Remove("tmp.db")
	n := LoadDatabase("../../database.su", "tmp.db")
	fmt.Println("loaded", n, "tables in", time.Since(t).Round(time.Millisecond))
	ck(CheckDatabase("tmp.db"))
}
