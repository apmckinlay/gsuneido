// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

type Where struct {
	Query1
	expr Expr
}

func (w *Where) String() string {
	return w.Query1.String() + " where " + w.expr.String()
}
