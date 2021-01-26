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

func (jn *Join) String() string {
	return jn.string("join")
}

func (jn *Join) string(op string) string {
	by := ""
	if len(jn.by) > 0 {
		by = "by" + str.Join("(,)", jn.by...) + " "
	}
	return jn.Query1.String() + " " + op + " " + by + jn.source2.String()
}

func (jn *Join) Columns() []string {
	return sset.Union(jn.source.Columns(), jn.source2.Columns())
}

func (jn *Join) Transform() Query {
	jn.source = jn.source.Transform()
	jn.source2 = jn.source2.Transform()
	return jn
}

type LeftJoin struct {
	Join
}

func (lj *LeftJoin) String() string {
	return lj.string("leftjoin")
}

func (lj *LeftJoin) Transform() Query {
	lj.source = lj.source.Transform()
	lj.source2 = lj.source2.Transform()
	return lj
}
