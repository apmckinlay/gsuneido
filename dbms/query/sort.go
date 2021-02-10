// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/util/str"

type Sort struct {
	Query1
	reverse bool
	columns []string
}

func (sort *Sort) String() string {
	s := sort.Query1.String() + " SORT "
	if sort.reverse {
		s += "reverse "
	}
	return s + str.Join(", ", sort.columns)
}

func (sort *Sort) Transform() Query {
	sort.source = sort.source.Transform()
	return sort
}
