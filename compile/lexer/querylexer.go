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
			return tok.By, "by"
		}
		if s == "in" {
			return tok.In, "in"
		}
		if s == "is" {
			return tok.Is, "is"
		}
		if s == "or" {
			return tok.Or, "or"
		}
		if s == "to" {
			return tok.To, "to"
		}
	case 3:
		if s == "and" {
			return tok.And, "and"
		}
		if s == "key" {
			return tok.Key, "key"
		}
		if s == "set" {
			return tok.Set, "set"
		}
		if s == "max" {
			return tok.Max, "max"
		}
		if s == "min" {
			return tok.Min, "min"
		}
		if s == "not" {
			return tok.Not, "not"
		}
	case 4:
		if s == "drop" {
			return tok.Drop, "drop"
		}
		if s == "into" {
			return tok.Into, "into"
		}
		if s == "isnt" {
			return tok.Isnt, "isnt"
		}
		if s == "join" {
			return tok.Join, "join"
		}
		if s == "list" {
			return tok.List, "list"
		}
		if s == "sort" {
			return tok.Sort, "sort"
		}
		if s == "true" {
			return tok.True, "true"
		}
		if s == "view" {
			return tok.View, "view"
		}
	case 5:
		if s == "alter" {
			return tok.Alter, "alter"
		}
		if s == "total" {
			return tok.Total, "total"
		}
		if s == "class" {
			return tok.Class, "class"
		}
		if s == "count" {
			return tok.Count, "count"
		}
		if s == "false" {
			return tok.False, "false"
		}
		if s == "index" {
			return tok.Index, "index"
		}
		if s == "minus" {
			return tok.Minus, "minus"
		}
		if s == "sview" {
			return tok.Sview, "sview"
		}
		if s == "union" {
			return tok.Union, "union"
		}
		if s == "times" {
			return tok.Times, "times"
		}
		if s == "where" {
			return tok.Where, "where"
		}
	case 6:
		if s == "create" {
			return tok.Create, "create"
		}
		if s == "delete" {
			return tok.Delete, "delete"
		}
		if s == "ensure" {
			return tok.Ensure, "ensure"
		}
		if s == "insert" {
			return tok.Insert, "insert"
		}
		if s == "extend" {
			return tok.Extend, "extend"
		}
		if s == "remove" {
			return tok.Remove, "remove"
		}
		if s == "rename" {
			return tok.Rename, "rename"
		}
		if s == "unique" {
			return tok.Unique, "unique"
		}
		if s == "update" {
			return tok.Update, "update"
		}
	case 7:
		if s == "average" {
			return tok.Average, "average"
		}
		if s == "cascade" {
			return tok.Cascade, "cascade"
		}
		if s == "destroy" {
			return tok.Drop, "destroy"
		}
		if s == "history" {
			return tok.History, "history"
		}
		if s == "project" {
			return tok.Project, "project"
		}
		if s == "reverse" {
			return tok.Reverse, "reverse"
		}
	case 8:
		if s == "function" {
			return tok.Function, "function"
		}
		if s == "leftjoin" {
			return tok.Leftjoin, "leftjoin"
		}
	case 9:
		if s == "intersect" {
			return tok.Intersect, "intersect"
		}
		if s == "summarize" {
			return tok.Summarize, "summarize"
		}
		if s == "tempindex" {
			return tok.TempIndex, "tempindex"
		}
	}
	return tok.Nil, ""
}
