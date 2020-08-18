// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"testing"
)

func TestLoad(*testing.T) {
	if testing.Short() {
		return
	}
	n := LoadTable("gl_transactions.su")
	fmt.Println("loaded", n, "records")
}
