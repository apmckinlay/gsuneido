package hmap

import (
	"fmt"
	"io"
)

func (hm *Hmap) dump(w io.Writer) {
	fmt.Fprintf(w, "Hmap capacity %v size %v\n", len(hm.buckets)*load, hm.size)
	for b := 0; b < len(hm.buckets); b++ {
		before := fmt.Sprint("[", b, "] ")
		after := ""
		for buck := &hm.buckets[b]; buck != nil; buck = buck.overflow {
			for i := 0; i < bucketsize; i++ {
				if buck.tophash[i] != 0 {
					fmt.Fprint(w, before, buck.keys[i], ":", buck.vals[i])
					before = " "
					after = "\n"
				}
			}
			before = ", "
		}
		fmt.Fprint(w, after)
	}
}
