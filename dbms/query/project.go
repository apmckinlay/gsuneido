// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import "github.com/apmckinlay/gsuneido/util/str"

type Project struct {
	Query1
	columns []string
}

func (p *Project) String() string {
	return p.Query1.String() + " project " + str.Join(", ", p.columns...)
}
