// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Join struct {
	Query2
	by []string
}

func (jn *Join) Init() {
	jn.source.Init()
	jn.source2.Init()
	by := sset.Intersect(jn.source.Columns(), jn.source2.Columns())
	if len(by) == 0 {
		panic("join: common columns required")
	}
	if jn.by == nil {
		jn.by = by
	} else if !sset.Equal(jn.by, by) {
		panic("join: by does not match common columns")
	}
	//TODO whether by contains key
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

func (lj *LeftJoin) Transform() Query {
	lj.source = lj.source.Transform()
	lj.source2 = lj.source2.Transform()
	return lj
}
