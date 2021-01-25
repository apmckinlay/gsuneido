// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Intersect struct {
	Query2
}

func (it *Intersect) String() string {
	return it.Query2.String("intersect")
}

func (it *Intersect) Transform() Query {
	it.source = it.source.Transform()
	it.source2 = it.source2.Transform()
	return it
}
