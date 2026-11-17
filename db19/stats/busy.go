// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package stats

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/apmckinlay/gsuneido/util/ss"
)

// BusyTally tracks frequent where columns and potential index scan savings.
type BusyTally struct {
	lock  sync.Mutex
	busy  *ss.Sketch[string]
	scans map[string]map[int]int
}

const busySize = 100 // AI recommendation
const busyFile = "busycols.txt"

// Add adds a weighted column (table.column).
func (bt *BusyTally) Add(col string, weight int) {
	bt.lock.Lock()
	defer bt.lock.Unlock()
	if bt.busy == nil {
		bt.busy = ss.New[string](busySize)
	}
	bt.busy.AddWeight(col, weight)
}

// AddScan adds a weighted index scan by table and index ordinal.
func (bt *BusyTally) AddScan(table string, index, cost int) {
	bt.lock.Lock()
	defer bt.lock.Unlock()
	if bt.scans == nil {
		bt.scans = make(map[string]map[int]int)
	}
	m := bt.scans[table]
	if m == nil {
		m = make(map[int]int)
		bt.scans[table] = m
	}
	m[index] += cost
}

// ResetScans discards the scan weights for a table.
// It should be called when the table's indexes change so that weights
// recorded by index ordinal are not applied to different indexes.
func (bt *BusyTally) ResetScans(table string) {
	bt.lock.Lock()
	defer bt.lock.Unlock()
	delete(bt.scans, table)
}

// Save appends busy columns and scanned indexes to a file.
// getIndexes gets the indexes for a table or nil if not found
func (bt *BusyTally) Save(getIndexes func(table string) [][]string) {
	bt.lock.Lock()
	busy := bt.busy
	scans := bt.scans
	bt.busy = ss.New[string](busySize)
	bt.scans = make(map[string]map[int]int)
	bt.lock.Unlock()
	if (busy == nil || busy.Len() == 0) && len(scans) == 0 {
		return
	}
	f, err := os.OpenFile(busyFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, e := range busy.Top() {
		fmt.Fprintf(w, "%d\t%s\n", e.Count, e.Value)
	}
	for table, indexCounts := range scans {
		indexes := getIndexes(table)
		for ixi, count := range indexCounts {
			if ixi < len(indexes) {
				ix := strings.Join(indexes[ixi], ",")
				fmt.Fprintf(w, "%d\t%s\t%s\n", count, table, ix)
			}
		}
	}
	if err := w.Flush(); err != nil {
		log.Println("ERROR:", err)
	}
}

// Busy is the most frequently used (in where clauses) columns for each table
type Busy map[string][]string

// Scans maps tables to indexes to their estimated records scanned
type Scans map[string]map[string]int

// LoadBusy consumes the file, returning busy columns and index scan weights.
func LoadBusy() (Busy, Scans) {
	f, err := os.Open(busyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println("ERROR: reading", busyFile, ":", err)
		}
		return nil, nil
	}
	defer func() {
		f.Close()
		if err := os.Remove(busyFile); err != nil {
			log.Println("ERROR can't remove "+busyFile+":", err)
			if err := os.Truncate(busyFile, 0); err != nil {
				log.Println("ERROR can't truncate "+busyFile+":", err)
			}
		}
	}()
	counts := make(map[string]int)
	scans := make(Scans)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		before, col, ok := strings.Cut(line, "\t")
		if !ok {
			continue
		}
		count, err := strconv.Atoi(before)
		if err != nil || count <= 0 {
			continue
		}
		if table, index, ok := strings.Cut(col, "\t"); ok {
			if index == "" {
				continue
			}
			if scans[table] == nil {
				scans[table] = make(map[string]int)
			}
			scans[table][index] += count
			continue
		}
		counts[col] += count
	}
	if err := scanner.Err(); err != nil {
		log.Println("ERROR: reading", busyFile, ":", err)
		return nil, nil
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
	n := 0
	for _, item := range items {
		if n++; n > busySize {
			break
		}
		table, col, ok := strings.Cut(item.key, ".")
		if !ok {
			continue
		}
		result[table] = append(result[table], col)
	}
	return result, scans
}
