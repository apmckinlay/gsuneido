// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Join struct {
	Query2
	by []string
	joinType
}

type joinType int

const (
	one_one joinType = iota
	one_n
	n_one
	n_n
)

func (jn *Join) Init() {
	jn.Query2.Init()

	by := sset.Intersect(jn.source.Columns(), jn.source2.Columns())
	if len(by) == 0 {
		panic("join: common columns required")
	}
	if jn.by == nil {
		jn.by = by
	} else if !sset.Equal(jn.by, by) {
		panic("join: by does not match common columns")
	}

	k1 := containsKey(jn.by, jn.source.Keys())
	k2 := containsKey(jn.by, jn.source2.Keys())
	if k1 && k2 {
		jn.joinType = one_one
	} else if k1 {
		jn.joinType = one_n
	} else if k2 {
		jn.joinType = n_one
	} else {
		jn.joinType = n_n
	}
}

func (jn *Join) String() string {
	return jn.string("JOIN")
}

func (jn *Join) string(op string) string {
	by := ""
	if len(jn.by) > 0 {
		by = "by" + str.Join("(,)", jn.by) + " "
	}
	return paren(jn.source) + " " + op + " " + by + paren(jn.source2)
}

func (jn *Join) Columns() []string {
	return sset.Union(jn.source.Columns(), jn.source2.Columns())
}

func (jn *Join) Indexes() [][]string {
	switch jn.joinType {
	case one_one:
		return ssset.Union(jn.source.Indexes(), jn.source2.Indexes())
	case one_n:
		return jn.source2.Indexes()
	case n_one:
		return jn.source.Indexes()
	case n_n:
		// union of indexes that don't include joincols
		idxs := [][]string{}
		for _, ix := range jn.source.Indexes() {
			if sset.Disjoint(ix, jn.by) {
				idxs = append(idxs, ix)
			}
		}
		for _, ix := range jn.source2.Indexes() {
			if sset.Disjoint(ix, jn.by) {
				ssset.AddUnique(idxs, ix)
			}
		}
		return idxs
	default:
		panic("unknown join type")
	}
}

func (jn *Join) Keys() [][]string {
	switch jn.joinType {
	case one_one:
		return ssset.Union(jn.source.Keys(), jn.source2.Keys())
	case one_n:
		return jn.source2.Keys()
	case n_one:
		return jn.source.Keys()
	case n_n:
		return jn.keypairs()
	default:
		panic("unknown join type")
	}
}

func (jn *Join) Transform() Query {
	jn.source = jn.source.Transform()
	jn.source2 = jn.source2.Transform()
	return jn
}

func (jn *Join) Fixed() []Fixed {
	return combineFixed(jn.source.Fixed(), jn.source2.Fixed())
}

// LeftJoin ---------------------------------------------------------

type LeftJoin struct {
	Join
}

func (lj *LeftJoin) String() string {
	return lj.string("LEFTJOIN")
}

func (lj *LeftJoin) Keys() [][]string {
	switch lj.joinType {
	case one_one, n_one:
		return lj.source.Keys()
	case one_n, n_n:
		return lj.keypairs()
	default:
		panic("unknown join type")
	}
}

func (lj *LeftJoin) Transform() Query {
	lj.source = lj.source.Transform()
	lj.source2 = lj.source2.Transform()
	return lj
}
