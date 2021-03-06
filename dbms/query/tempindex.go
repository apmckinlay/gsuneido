// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/util/str"

type TempIndex struct {
	Query1
	order []string
	tran  QueryTran
}

func (ti *TempIndex) Transform() Query {
	return ti
}

func (ti *TempIndex) String() string {
	return ti.source.String() + " TEMPINDEX" + str.Join("(,)", ti.order)
}
