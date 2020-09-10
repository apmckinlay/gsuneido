// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package cksum handles adding and checking checksums on data.
// crc32 since this is faster than adler32.
// Castagnoli because it slows down less for random data lengths.
// It only uses the low 16 bits of the crc32 to save space.
// From testing this doesn't impair error detection much.
package cksum

import "hash/crc32"

const Len = 2 // 16 bits

var crc32table = crc32.MakeTable(crc32.Castagnoli)

// Update computes the checksum of data (excluding the checksum space)
// and stores it in the last Len bytes of data.
// The caller must allocate the space for the checksum.
func Update(data []byte) {
	n := len(data) - Len
	cs := crc32.Checksum(data[:n], crc32table)
	data[n] = byte(cs)
	data[n+1] = byte(cs >> 8)
}

// Check computes the checksum of data (excluding the checksum space)
// and compares it to the stored value at the end of the data.
// Panics on failure.
func Check(data []byte) bool {
	n := len(data) - Len
	cs := crc32.Checksum(data[:n], crc32table)
	return data[n] == byte(cs) && data[n+1] == byte(cs>>8)
}

func MustCheck(data []byte) {
	if !Check(data) {
		panic("checksum failure")
	}
}
