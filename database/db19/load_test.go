// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"testing"
	"time"
)

func TestLoadTable(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	n := LoadTable("stdlib.su")
	fmt.Println("loaded", n, "records in", time.Since(t))
}

func TestLoadDatabase(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	n := LoadDatabase()
	fmt.Println("loaded", n, "tables in", time.Since(t))
}
