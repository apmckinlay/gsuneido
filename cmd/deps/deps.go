// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
)

const dir = "builtin/goc"
const file = "builtin/goc/cside.h"

func main() {
	mod, err := getLatestModTime(dir)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	const prefix = "// deps last modified "
	s := string(b)
	if i := strings.Index(s, prefix); i != -1 {
		s = s[:i]
	}
	s = strings.TrimRight(s, "\r\n")
	s += "\n\n" + prefix + mod.UTC().Format("2006-01-02 15:04:05 UTC") + "\n"
	err = os.WriteFile(file, []byte(s), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func getLatestModTime(dir string) (time.Time, error) {
	var latest time.Time
	err := fs.WalkDir(os.DirFS(dir), ".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() &&
				d.Name() != "cside.h" &&
				(strings.HasSuffix(d.Name(), ".h") ||
					strings.HasSuffix(d.Name(), ".c") ||
					strings.HasSuffix(d.Name(), ".cpp")) {
				info, err := d.Info()
				if err != nil {
					return err
				}
				t := info.ModTime()
				if t.Compare(latest) > 0 {
					latest = t
				}
			}
			return nil
		})
	return latest, err
}
