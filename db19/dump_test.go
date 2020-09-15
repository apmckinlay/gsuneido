// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestDumpTable(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	defer os.Remove("tmp.su")
	n := DumpTable("../../suneido.db", "stdlib", "tmp.su")
	fmt.Println("dumped", n, "records in", time.Since(t).Round(time.Millisecond))
}

func TestDumpDatabase(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	defer os.Remove("tmp.su")
	n := DumpDatabase("../../suneido.db", "tmp.su")
	fmt.Println("dumped", n, "tables in", time.Since(t).Round(time.Millisecond))
}
