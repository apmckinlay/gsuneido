// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lexer

import (
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

func NewQueryLexer(src string) *Lexer {
	return &Lexer{src: src, keyword: queryKeyword, nlwhite: true}
}

func queryKeyword(s string) (tok.Token, string) {
	switch len(s) {
	case 2:
		if s == "by" {
			return tok.By, s
		}
		if s == "in" {
			return tok.In, s
		}
		if s == "is" {
			return tok.Is, s
		}
		if s == "or" {
			return tok.Or, s
		}
		if s == "to" {
			return tok.To, s
		}
	case 3:
		if s == "and" {
			return tok.And, s
		}
		if s == "key" {
			return tok.Key, s
		}
		if s == "set" {
			return tok.Set, s
		}
		if s == "max" {
			return tok.Max, s
		}
		if s == "min" {
			return tok.Min, s
		}
		if s == "not" {
			return tok.Not, s
		}
	case 4:
		if s == "drop" {
			return tok.Drop, s
		}
		if s == "into" {
			return tok.Into, s
		}
		if s == "isnt" {
			return tok.Isnt, s
		}
		if s == "join" {
			return tok.Join, s
		}
		if s == "list" {
			return tok.List, s
		}
		if s == "sort" {
			return tok.Sort, s
		}
		if s == "true" {
			return tok.True, s
		}
		if s == "view" {
			return tok.View, s
		}
	case 5:
		if s == "alter" {
			return tok.Alter, s
		}
		if s == "total" {
			return tok.Total, s
		}
		if s == "class" {
			return tok.Class, s
		}
		if s == "count" {
			return tok.Count, s
		}
		if s == "false" {
			return tok.False, s
		}
		if s == "index" {
			return tok.Index, s
		}
		if s == "minus" {
			return tok.Minus, s
		}
		if s == "sview" {
			return tok.Sview, s
		}
		if s == "union" {
			return tok.Union, s
		}
		if s == "times" {
			return tok.Times, s
		}
		if s == "where" {
			return tok.Where, s
		}
	case 6:
		if s == "create" {
			return tok.Create, s
		}
		if s == "delete" {
			return tok.Delete, s
		}
		if s == "ensure" {
			return tok.Ensure, s
		}
		if s == "insert" {
			return tok.Insert, s
		}
		if s == "extend" {
			return tok.Extend, s
		}
		if s == "remove" {
			return tok.Remove, s
		}
		if s == "rename" {
			return tok.Rename, s
		}
		if s == "unique" {
			return tok.Unique, s
		}
		if s == "update" {
			return tok.Update, s
		}
	case 7:
		if s == "average" {
			return tok.Average, s
		}
		if s == "cascade" {
			return tok.Cascade, s
		}
		if s == "destroy" {
			return tok.Drop, s
		}
		if s == "history" {
			return tok.History, s
		}
		if s == "project" {
			return tok.Project, s
		}
		if s == "reverse" {
			return tok.Reverse, s
		}
	case 8:
		if s == "function" {
			return tok.Function, s
		}
		if s == "leftjoin" {
			return tok.Leftjoin, s
		}
	case 9:
		if s == "intersect" {
			return tok.Intersect, s
		}
		if s == "summarize" {
			return tok.Summarize, s
		}
		if s == "tempindex" {
			return tok.TempIndex, s
		}
	}
	return tok.Identifier, s
}
