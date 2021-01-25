// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Union struct {
	Query2
}

func (u *Union) String() string {
	return u.Query2.String("union")
}

func (u *Union) Transform() Query {
	u.source = u.source.Transform()
	u.source2 = u.source2.Transform()
	return u
}
