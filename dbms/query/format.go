// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"reflect"
	"strings"

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
			in + format1(q) + "\n" +
			format(indent, q.Source2(), 1)
		if parens >= 1 {
			s = addParens(s)
		}
	case *Sort:
		s = format(indent, q.Source(), 0) + "\n" +
			in + format1(q)
	case *View:
		s = in + q.name + " /*view*/"
	case q1i:
		s = format(indent, q.Source(), 2) + "\n" +
			in + format1(q)
		if parens == 1 {
			s = addParens(s)
		}
	default:
		s = in + q.String()
	}
	return s
}

type formatable interface {
	format() string
}

func format1(q q1i) string {
	if f, ok := q.(formatable); ok {
		return f.format()
	}
	op, rest, _ := strings.Cut(q.stringOp(), " ")
	return str.ToLower(op) + str.Opt(" ", rest)
}

func which(x any) string {
	t := reflect.TypeOf(x)
	for t.Kind() == reflect.Ptr {
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
