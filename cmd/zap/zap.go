// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	rand "math/rand/v2"
	"os"
	"strconv"
)

func main() {
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
	for range nspots {
		// at := int64(rand.Int63n(size - int64(bytesPerSpot)))
		at := size - rand.Int64N(1_000_000) - int64(bytesPerSpot) // near end
		fmt.Println("zapped", file, "-", bytesPerSpot, "bytes at", at)
		f.Seek(at, 0)
		for range bytesPerSpot {
			f.WriteString(string(rune(rand.IntN(256))))
		}
	}
}
