// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var idxUse = make(map[string]int)
var idxUseLock sync.Mutex // guards idxUse

func IdxUse(table string, idx []string) {
	if strings.HasPrefix(table, "tests_20") {
		return
	}
	var sb strings.Builder
	sb.WriteString(table)
	sep := "^"
	for _, fld := range idx {
		sb.WriteString(sep)
		sb.WriteString(fld)
		sep = ","
	}
	k := sb.String()
	idxUseLock.Lock()
	idxUse[k]++
	idxUseLock.Unlock()
}

func init() { go idxUseSave() }

func idxUseSave() {
	for {
		time.Sleep(8 * time.Hour)
		idxUseWrite(PullIdxUse())
	}
}

// PullIdxUse returns the current index usage and resets it
func PullIdxUse() map[string]int {
	idxUseLock.Lock()
	iu := idxUse
	idxUse = make(map[string]int, len(iu))
	idxUseLock.Unlock()
	return iu
}

func idxUseWrite(iu map[string]int) {
	if len(iu) == 0 {
		return
	}
	f, err := os.OpenFile("idxuse.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for k, v := range iu {
		fmt.Fprintf(w, "%d\t%s\n", v, k)
	}
	if err := w.Flush(); err != nil {
		log.Println("ERROR: writing idxuse:", err)
	}
}
