// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/util/sset"

type Intersect struct {
	Compatible
}

func (it *Intersect) String() string {
	return it.Query2.String("INTERSECT")
}

func (it *Intersect) Columns() []string {
	return sset.Intersect(it.source.Columns(), it.source2.Columns())
}

func (it *Intersect) Transform() Query {
	it.source = it.source.Transform()
	it.source2 = it.source2.Transform()
	return it
}
