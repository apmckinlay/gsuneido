// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lexer

import (
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/str"
)

func NewQueryLexer(src string) *Lexer {
	return &Lexer{src: src, keyword: queryKeyword, nlwhite: true}
}

func queryKeyword(s string) (tok.Token, string) {
	ls := str.ToLower(s)
	if tok, ok := queryKeywords[ls]; ok {
		return tok, ls
	}
	return tok.Identifier, s
}

var queryKeywords = map[string]tok.Token{
	"alter":     tok.Alter,
	"and":       tok.And,
	"average":   tok.Average,
	"by":        tok.By,
	"cascade":   tok.Cascade,
	"count":     tok.Count,
	"create":    tok.Create,
	"delete":    tok.Delete,
	"destroy":   tok.Drop,
	"drop":      tok.Drop,
	"ensure":    tok.Ensure,
	"extend":    tok.Extend,
	"false":     tok.False,
	"history":   tok.History,
	"in":        tok.In,
	"index":     tok.Index,
	"insert":    tok.Insert,
	"intersect": tok.Intersect,
	"into":      tok.Into,
	"is":        tok.Is,
	"isnt":      tok.Isnt,
	"join":      tok.Join,
	"key":       tok.Key,
	"leftjoin":  tok.Leftjoin,
	"list":      tok.List,
	"max":       tok.Max,
	"min":       tok.Min,
	"minus":     tok.Minus,
	"not":       tok.Not,
	"or":        tok.Or,
	"project":   tok.Project,
	"remove":    tok.Remove,
	"rename":    tok.Rename,
	"reverse":   tok.Reverse,
	"set":       tok.Set,
	"sort":      tok.Sort,
	"summarize": tok.Summarize,
	"sview":     tok.Sview,
	"times":     tok.Times,
	"to":        tok.To,
	"total":     tok.Total,
	"true":      tok.True,
	"union":     tok.Union,
	"unique":    tok.Unique,
	"update":    tok.Update,
	"view":      tok.View,
	"where":     tok.Where,
}
