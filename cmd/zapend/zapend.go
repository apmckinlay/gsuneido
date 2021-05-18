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
	bytes := 1
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	if len(os.Args) > 2 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalln("usage: zapend [filename [nbytes]]")
		}
		bytes = n
	}
	fmt.Println(file, bytes)

	f, err := os.OpenFile(file, os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	stat, _ := f.Stat()
	size := stat.Size()
	fmt.Println("size", size)
	// at := int64(rand.Int63n(size - int64(bytesPerSpot)))
	at := size - int64(bytes)
	f.Seek(at, 0)
	buf := make([]byte, 1)
	for j := 0; j < bytes; j++ {
		buf[0] = byte(rand.Intn(256))
		f.Write(buf)
	}
}
