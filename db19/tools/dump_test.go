// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDumpTable(t *testing.T) {
	if testing.Short() {
		return
	}
	start := time.Now()
	defer os.Remove("tmp.su")
	n, err := DumpTable("../../suneido.db", "stdlib", "tmp.su")
	assert.T(t).This(err).Is(nil)
	fmt.Println("dumped", n, "records in", time.Since(start).Round(time.Millisecond))
}

func TestDumpDatabase(t *testing.T) {
	if testing.Short() {
		return
	}
	start := time.Now()
	defer os.Remove("tmp.su")
	n, err := DumpDatabase("../../suneido.db", "tmp.su")
	assert.T(t).This(err).Is(nil)
	fmt.Println("dumped", n, "tables in", time.Since(start).Round(time.Millisecond))
}
