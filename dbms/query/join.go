// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/sset"
	"github.com/apmckinlay/gsuneido/util/ssset"
	"github.com/apmckinlay/gsuneido/util/str"
)

type Join struct {
	Query2
	by []string
	joinType
	nr int
}

type joinApproach struct {
	reverse bool
}

type joinType int

const (
	one_one joinType = iota + 1
	one_n
	n_one
)

func (jt joinType) String() string {
	switch jt {
	case 0:
		return ""
	case one_one:
		return "1:1"
	case one_n:
		return "1:n"
	case n_one:
		return "n:1"
	default:
		panic("bad joinType")
	}
}

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
		panic("join: does not support many to many")
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
	return parenQ2(jn.source) + " " + op + " " +
		str.Opt(jn.joinType.String(), " ") + by + paren(jn.source2)
}

func (jn *Join) Columns() []string {
	return sset.Union(jn.source.Columns(), jn.source2.Columns())
}

func (jn *Join) Indexes() [][]string {
	// can really only provide source.indexes() but optimize may swap.
	// optimize will return impossible for source2 indexes.
	return ssset.Union(jn.source.Indexes(), jn.source2.Indexes())
}

func (jn *Join) Keys() [][]string {
	switch jn.joinType {
	case one_one:
		return ssset.Union(jn.source.Keys(), jn.source2.Keys())
	case one_n:
		return jn.source2.Keys()
	case n_one:
		return jn.source.Keys()
	default:
		panic("unknown join type")
	}
}

func (jn *Join) Fixed() []Fixed {
	return combineFixed(jn.source.Fixed(), jn.source2.Fixed())
}

func (jn *Join) Transform() Query {
	jn.source = jn.source.Transform()
	jn.source2 = jn.source2.Transform()
	return jn
}

func (jn *Join) optimize(mode Mode, index []string) (Cost, interface{}) {
	defer be(gin("Join", jn, index))
	cost := jn.opt(jn.source, jn.source2, jn.joinType, mode, index)
	cost2 := outOfOrder +
		jn.opt(jn.source2, jn.source, jn.joinType.reverse(), mode, index)
	trace("cost1", cost, "cost2 (reverse)", cost2)
	approach := &joinApproach{}
	if cost2 < cost {
		cost = cost2
		approach.reverse = true
	}
	if cost >= impossible {
		return impossible, nil
	}
	return cost, approach
}

func (jt joinType) reverse() joinType {
	switch jt {
	case one_n:
		return n_one
	case n_one:
		return one_n
	}
	return jt
}

func (jn *Join) opt(src1, src2 Query, joinType joinType,
	mode Mode, index []string) Cost {
	trace("OPT", paren(src1), "JOIN", joinType, paren(src2))
	// always have to read all of source 1
	cost1 := Optimize(src1, mode, index)
	if cost1 >= impossible {
		return impossible
	}
	cost2 := Optimize(src2, mode, jn.by)
	if cost2 >= impossible {
		return impossible
	}
	nrows1 := src1.nrows()
	// should only be taking a portion of the variable cost2,
	// not the fixed temp index cost2 (so 2/3 instead of 1/2)
	cost := cost1 + (nrows1 * src2.lookupCost()) + (cost2 * 2 / 3)
	trace("join opt", cost1, "+ (", nrows1, "*", src2.lookupCost(), ") + (",
		cost2, "* 2/3 ) =", cost)
	return cost
}

func (jn *Join) setApproach(index []string, approach interface{}, tran QueryTran) {
	ap := approach.(*joinApproach)
	if ap.reverse {
		jn.source, jn.source2 = jn.source2, jn.source
		jn.joinType = jn.joinType.reverse()
	}
	jn.source = SetApproach(jn.source, index, tran)
	jn.source2 = SetApproach(jn.source2, jn.by, tran)
}

func (jn *Join) nrows() int {
	// n_one and one_n assume records will have matching counterparts
	nrows1 := jn.source.nrows()
	nrows2 := jn.source2.nrows()
	var nrows int
	switch jn.joinType {
	case one_one:
		nrows = ints.Min(nrows1, nrows2)
	case n_one:
		nrows = nrows1
	case one_n:
		nrows = nrows2
	default:
		panic("shouldn't reach here")
	}
	return nrows / 2 // actual will be between 0 and nrows so estimate halfway
}

func (jn *Join) rowSize() int {
	return jn.source.rowSize() + jn.source2.rowSize()
}

func (jn *Join) lookupCost() int {
	return jn.source.lookupCost() * 2 // ???
}

// LeftJoin ---------------------------------------------------------

type LeftJoin struct {
	Join
}

func (lj *LeftJoin) String() string {
	return lj.string("LEFTJOIN")
}

func (lj *LeftJoin) Keys() [][]string {
	// can't use source2.Keys() like Join.Keys()
	// because multiple right sides can be missing/blank
	switch lj.joinType {
	case one_one, n_one:
		return lj.source.Keys()
	case one_n:
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

func (lj *LeftJoin) optimize(mode Mode, index []string) (Cost, interface{}) {
	approach := &joinApproach{}
	cost := lj.opt(lj.source, lj.source2, lj.joinType, mode, index)
	return cost, approach
}

func (lj *LeftJoin) setApproach(index []string, _ interface{}, tran QueryTran) {
	lj.source = SetApproach(lj.source, index, tran)
	lj.source2 = SetApproach(lj.source2, lj.by, tran)
}

func (lj *LeftJoin) nrows() int {
	return lj.source.nrows()
}
