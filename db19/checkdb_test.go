// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/stor"
)

func TestCheckDatabase(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	if err := CheckDatabase("../../suneido.db"); err != nil {
		fmt.Println(err, err.table)
	} else {
		fmt.Println("database checked in", time.Since(t).Round(time.Millisecond))
	}
}

func TestQuickCheck(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	db, e := openDatabase("../../suneido.db", stor.READ, false)
	if e != nil {
		panic(e)
	}
	db.QuickCheck()
	fmt.Println("database checked in", time.Since(t).Round(time.Millisecond))
}
