// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"testing"
)

func TestRepair(*testing.T) {
	if testing.Short() {
		return
	}
	_, err := Repair("../suneido.db", nil)
	fmt.Println(err)
}
