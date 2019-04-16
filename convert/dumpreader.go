package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	. "github.com/apmckinlay/gsuneido/database/record"
)

// Read a dump file in new format
func Read() {
	fin, err := os.Open("database.su2")
	ckerr(err)
	in := bufio.NewReader(fin)
	s, err := in.ReadString('\n')
	ckerr(err)
	hdr := "Suneido dump 2\n"
	if s != hdr {
		panic("\n\tgot: " + s + "\n\texpected: " + hdr)
	}
	n := 0
	for { // each table
		schema, err := in.ReadString('\n')
		if err == io.EOF {
			break
		}
		ckerr(err)
		if !strings.HasPrefix(schema, "====== ") {
			panic("bad schema: " + schema)
		}
		// fmt.Print(schema)
		readTable(in)
		n++
	}
	fin.Close()
	fmt.Println(n, "tables,", nrecs, "records")
}

func readTable(in *bufio.Reader) {
	for { // each record
		_, err := io.ReadFull(in, intbuf)
		if err == io.EOF {
			break
		}
		ckerr(err)
		size := int(binary.BigEndian.Uint32(intbuf))
		if size == 0 {
			break
		}
		_, err = io.ReadFull(in, inbuf[:size])
		ckerr(err)
		checkRecord(Record(inbuf[:size]))
		nrecs++
	}
}

func checkRecord(rec Record) {
	n := rec.Count()
	for i := 0; i < n; i++ {
		rec.GetVal(i) // unpack
	}
}
