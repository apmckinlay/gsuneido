// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestLoadTable(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	defer os.Remove("tmp.db")
	db := CreateDatabase("tmp.db")
	defer db.Close()
	n := db.LoadTable("stdlib")
	fmt.Println("loaded", n, "records in", time.Since(t))
}

func TestLoadDatabase(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	defer os.Remove("tmp.db")
	n := LoadDatabase("database.su", "tmp.db")
	fmt.Println("loaded", n, "tables in", time.Since(t))
}
