// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/apmckinlay/gsuneido/core/trace"
	"github.com/apmckinlay/gsuneido/util/str"
)

func Format(t QueryTran, query string) string {
	q := parseQuery(query, t, nil, nil, true)
	return format(0, q, 0)
}

func format(indent int, q Query, parens int) string { // recursive
	in := strings.Repeat(" ", indent*4)
	var s string
	switch q := q.(type) {
	case q2i:
		indent++
		leftin := indent
		left := q.Source()
		if _, ok := left.(q2i); ok && which(left) == which(q) {
			leftin--
		}
		s = format(leftin, q.Source(), 1) + "\n" +
			in + q.String() + "\n" +
			format(indent, q.Source2(), 1)
		if parens >= 1 {
			s = addParens(s)
		}
	case *Sort:
		s = format(indent, q.Source(), 0) + "\n" +
			in + q.String()
	case *View:
		s = in + q.String()
	case q1i:
		s = format(indent, q.Source(), 2) + "\n" +
			in + q.String()
		if parens == 1 {
			s = addParens(s)
		}
	default:
		s = in + q.String()
	}
	return s
}

func which(x any) string {
	t := reflect.TypeOf(x)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Name()
}

func addParens(s string) string {
	i := 0
	for s[i] == ' ' {
		i++
	}
	return s[:i] + "(" + s[i:] + ")"
}

//-------------------------------------------------------------------

// String prints the full query, including child sources
// whereas query.String only shows that operation
func String(q Query) string {
	switch qi := q.(type) {
	case q2i:
		return paren2(qi.Source()) + " " + q.String() + " " + paren1(qi.Source2())
	case *Sort:
		return String(qi.Source()) + str.Opt(" ", q.String()) // no parens
	case *View:
		return q.String()
	case q1i:
		return paren2(qi.Source()) + str.Opt(" ", q.String())
	default:
		return q.String()
	}
}

func paren1(q Query) string {
	switch q.(type) {
	case *Table, *Tables, *TablesLookup, *Columns, *Indexes, *Views,
		*Nothing, *ProjectNone:
		return String(q)
	}
	return "(" + String(q) + ")"
}

func paren2(q Query) string {
	if _, ok := q.(q2i); ok {
		return "(" + String(q) + ")"
	}
	return String(q)
}

//-------------------------------------------------------------------

func Strategy(q Query) string {
	return strategy(q, 0)
}

const indent1 = "    "

func strategy(q Query, indent int) string { // recursive
	in := strings.Repeat(indent1, indent)
	nrows, pop := q.Nrows()
	m := q.Metrics()
	cost := "{"
	if m.frac != 1 {
		cost += fmt.Sprintf("%.3fx ", m.frac)
	}
	cost += trace.Number(nrows)
	if nrows != pop {
		cost += "/" + trace.Number(pop)
	}
	cost += " " + trace.Number(m.fixcost) + "+" + trace.Number(m.varcost)
	cost += "} "
	switch q := q.(type) {
	case *Sort:
		if q.String() == "" {
			return strategy(q.Source(), indent)
		} else {
			return strategy(q.Source(), indent) + "\n" +
				in + cost + q.String()
		}
	case q2i:
		return strategy(q.Source(), indent+1) + "\n" +
			in + cost + q.String() + "\n" +
			strategy(q.Source2(), indent+1)
	case q1i:
		return strategy(q.Source(), indent) + "\n" +
			in + cost + q.String()
	default:
		return in + cost + q.String()
	}
}

// Strategy2 is like Strategy but without the cost/size estimates
// so it is more stable for tests
func Strategy2(q Query) string {
	return strategy2(q, 0)
}

func strategy2(q Query, indent int) string { // recursive
	in := strings.Repeat(indent1, indent)
	switch q := q.(type) {
	case *Sort:
		if q.String() == "" {
			return strategy2(q.Source(), indent)
		} else {
			return strategy2(q.Source(), indent) + "\n" +
				in + q.String()
		}
	case q2i:
		return strategy2(q.Source(), indent+1) + "\n" +
			in + q.String() + "\n" +
			strategy2(q.Source2(), indent+1)
	case q1i:
		return strategy2(q.Source(), indent) + "\n" +
			in + q.String()
	default:
		return in + q.String()
	}
}
