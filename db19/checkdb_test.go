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
	if err := CheckDatabase("../suneido.db", true); err != nil {
		fmt.Println(err, err.(*errCorrupt).table)
	} else {
		fmt.Println("database checked in", time.Since(t).Round(time.Millisecond))
	}
}

func TestQuickCheck(*testing.T) {
	if testing.Short() {
		return
	}
	t := time.Now()
	db, e := OpenDb("../suneido.db", stor.Read, false)
	if e != nil {
		panic(e)
	}
	if err := db.QuickCheck(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("database checked in", time.Since(t).Round(time.Millisecond))
	}
}
