// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckDatabase(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	if err := CheckDatabase("suneido.db"); err != "" {
		fmt.Println(err)
	} else {
		fmt.Println("checked database in", time.Since(t).Round(time.Millisecond))
	}
}
