package cksum

import (
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"math/rand"
	"testing"
)

var Sum uint32

const ndata = 1023

var Data = make([]byte, ndata)

func init() {
	for i := range ndata {
		Data[i] = byte(rand.Int31n(256))
	}
}

func BenchmarkAdler32(b *testing.B) {
	for b.Loop() {
		Sum += adler32.Checksum(Data)
	}
}

func BenchmarkCrc32IEEE(b *testing.B) {
	for b.Loop() {
		Sum += crc32.Checksum(Data, crc32.IEEETable)
	}
}

func BenchmarkCrc32Cast(b *testing.B) {
	table := crc32.MakeTable(crc32.Castagnoli)
	for b.Loop() {
		Sum += crc32.Checksum(Data, table)
	}
}

func TestDetection(*testing.T) {
	if testing.Short() {
		return
	}
	table := crc32.MakeTable(crc32.Castagnoli)
	var diff, lodiff, hidiff int
	for range 100000 {
		a := crc32.Checksum(Data, table)
		for range 4 {
			Data[rand.Int31n(ndata)] = byte(rand.Int31n(256))
		}
		b := crc32.Checksum(Data, table)
		if a != b {
			diff++
		}
		if uint16(a) != uint16(b) {
			lodiff++
		}
		if a>>16 != b>>16 {
			hidiff++
		}
	}
	fmt.Println("diff", diff, "lo", lodiff, "hi", hidiff)
}
