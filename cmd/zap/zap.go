// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	file := "suneido.db"
	bytesPerSpot := 1
	nspots := 1
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalln("usage: zap [filename [bytes [count]]]")
		}
		bytesPerSpot = n
	}
	if len(os.Args) > 3 {
		n, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Fatalln("usage: zap [filename [bytes [count]]]")
		}
		nspots = n
	}
	fmt.Println(file, bytesPerSpot, nspots)

	f, err := os.OpenFile(file, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	stat, _ := f.Stat()
	size := stat.Size()
	fmt.Println("size", size)
	for i := 0; i < nspots; i++ {
		// at := int64(rand.Int63n(size - int64(bytesPerSpot)))
		at := size - rand.Int63n(1_000_000) - int64(bytesPerSpot) // near end
		fmt.Println("zapped", file, "-", bytesPerSpot, "bytes at", at)
		f.Seek(at, 0)
		for j := 0; j < bytesPerSpot; j++ {
			f.WriteString(string(rune(rand.Intn(256))))
		}
	}
}
