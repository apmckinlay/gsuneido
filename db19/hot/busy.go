// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hot

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apmckinlay/gsuneido/util/exit"
	"github.com/apmckinlay/gsuneido/util/ss"
)

// BusyTally tracks columns that are frequently used in where clauses.
type BusyTally struct {
	lock sync.Mutex
	busy *ss.Sketch[string]
}

const busySize = 100               // AI recommendation
const busyInterval = 2 * time.Hour // ???
const busyFile = "hotcols.txt"

// Add adds a column with a given weight to the hot columns tracker
func (bt *BusyTally) Add(col string, weight int) {
	bt.lock.Lock()
	if bt.busy == nil {
		bt.busy = ss.New[string](busySize)
		go func() {
			ticker := time.NewTicker(busyInterval)
			for range ticker.C {
				bt.Save()
			}
		}()
		var startTime = time.Now()
		exit.Add("busy", func() { // for testing/development
			if time.Since(startTime) < busyInterval {
				bt.Save()
			}
		})
	}
	defer bt.lock.Unlock()
	bt.busy.AddWeight(col, weight)
}

// Save writes the busy columns to a file
func (bt *BusyTally) Save() {
	bt.lock.Lock()
	hot := bt.busy
	bt.busy = ss.New[string](busySize)
	bt.lock.Unlock()
	if hot == nil || hot.Len() == 0 {
		return
	}
	f, err := os.OpenFile(busyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, e := range hot.Top() {
		fmt.Fprintf(w, "%d\t%s\n", e.Count, e.Value)
	}
	if err := w.Flush(); err != nil {
		log.Println("ERROR:", err)
	}
}

// Busy is the most frequently used (in where clauses) columns for each table
type Busy map[string][]string

// LoadBusy reads the busy columns from a file
// and returns a map of table to frequently used column names
func LoadBusy() Busy {
	f, err := os.Open(busyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println("ERROR: reading hotcols:", err)
		}
		return nil
	}
	defer func() {
		f.Close()
		if err := os.Remove(busyFile); err != nil {
			log.Println("ERROR can't remove " + busyFile + ":", err)
			if err := os.Truncate(busyFile, 0); err != nil {
				log.Println("ERROR can't truncate " + busyFile + ":", err)
			}
		}
	}()
	counts := make(map[string]int)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		before, col, ok := strings.Cut(line, "\t")
		if !ok {
			continue
		}
		count, err := strconv.Atoi(before)
		if err != nil {
			continue
		}
		counts[col] += count
	}
	if err := scanner.Err(); err != nil {
		log.Println("ERROR: reading hotcols:", err)
		return nil
	}
	type kv struct {
		key   string
		count int
	}
	items := make([]kv, 0, len(counts))
	for k, v := range counts {
		items = append(items, kv{k, v})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].count > items[j].count
	})
	result := make(Busy)
	for _, item := range items {
		table, col, ok := strings.Cut(item.key, ".")
		if !ok {
			continue
		}
		if len(result[table]) < busySize {
			result[table] = append(result[table], col)
		}
	}
	return result
}
