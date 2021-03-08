// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/bytes"
)

func TestCostModel(*testing.T) {
	if testing.Short() {
		return
	}
	const dbsize = 2 * 1024 * 1024 * 1024 // 2 gb
	// create file
	file := filepath.Join(os.TempDir(), "testcostmodel")
	db, err := stor.MmapStor(file, stor.CREATE)
	if err != nil {
		panic(err)
	}
	defer os.Remove(file)
	defer db.Close()
	for db.Size() < dbsize {
		_, buf := db.Alloc(4 * 1024)
		bytes.Fill(buf, 0x55)
	}
	// random reads
	sum := byte(0)
	const nreads = 10_000
	for readSize := 16; readSize <= 64*1024; readSize *= 2 {
		t := time.Now()
		for i := 0; i < nreads; i++ {
			off := rand.Intn(dbsize)
			buf := db.Data(uint64(off))
			if len(buf) > readSize {
				buf = buf[:readSize]
			}
			for _, b := range buf {
				sum += b
			}
		}
		fmt.Println("read", readSize, time.Since(t)/nreads)
	}
	fmt.Println(sum)
}
