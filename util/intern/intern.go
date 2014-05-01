// Package intern handles sharing common strings
package intern

import "sync"

// redundant to store string as key and value
// but no way to get key that matched lookup

var table = make(map[string]string)
var lock sync.Mutex

func Intern(s string) string {
	lock.Lock()
	defer lock.Unlock()
	if x, ok := table[s]; ok {
		return x
	}
	table[s] = s
	return s
}
