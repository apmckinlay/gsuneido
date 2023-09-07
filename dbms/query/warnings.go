// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"log"
	"strings"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/str"
)

var seen = map[string]bool{}
var nWarn = 0

func Warnings(query string, q Query) {
	w := warnings(q)
	s := warnStr(query, w)
	if s != "" {
		if Suneido.Get(nil, SuStr("User")) == SuStr("default") {
			panic(s)
		} else {
			nWarn++
			if !seen[query] && nWarn < 10 {
				seen[query] = true
				log.Println("WARNING", s, query)
			}
		}
	}
}

const (
	projectNotUnique = 1
	unionNotDisjoint = 2
	joinManyToMany   = 4
)

func warnStr(query string, w int) string {
	var cb str.CommaBuilder
	if w&projectNotUnique != 0 &&
		!strings.Contains(query, "CHECKQUERY SUPPRESS: PROJECT NOT UNIQUE") {
		cb.Add("PROJECT NOT UNIQUE")
	}
	if w&unionNotDisjoint != 0 &&
		!strings.Contains(query, "CHECKQUERY SUPPRESS: UNION NOT DISJOINT") {
		cb.Add("UNION NOT DISJOINT")
	}
	if w&joinManyToMany != 0 &&
		!strings.Contains(query, "CHECKQUERY SUPPRESS: JOIN MANY TO MANY") {
		cb.Add("JOIN MANY TO MANY")
	}
	return cb.String()
}

func warnings(q Query) int { // recursive
	w := 0
	switch q := q.(type) {
	case *Project:
		if _, ok := q.source.(*Project); ok {
			log.Println("ERROR: transform did not merge project")
		}
		if !q.unique {
			w |= projectNotUnique
		}
	case *Union:
		if q.disjoint == "" {
			w |= unionNotDisjoint
		}
	case *Join:
		if q.joinType == n_n {
			w |= joinManyToMany
		}
	case *LeftJoin:
		if q.joinType == n_n {
			w |= joinManyToMany
		}
	case *Where:
		if _, ok := q.source.(*Where); ok {
			log.Println("ERROR: transform did not merge where")
		}
	case *Extend:
		if _, ok := q.source.(*Extend); ok {
			log.Println("ERROR: transform did not merge extend")
		}
	case *Rename:
		if _, ok := q.source.(*Rename); ok {
			log.Println("ERROR: transform did not merge rename")
		}
	}
	switch q := q.(type) {
	case q2i:
		w |= warnings(q.Source())
		w |= warnings(q.Source2())
	case q1i:
		w |= warnings(q.Source())
	}
	return w
}
