// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package lexer

import (
	"fmt"
	"strings"
	"testing"

	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hmap"
	"github.com/apmckinlay/gsuneido/util/str"
)

var T tok.Token
var S string

func BenchmarkKeyword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T, S = keyword(keywords[i%len(keywords)].kw)
	}
}

func BenchmarkKeywordLinear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T, S = keywordLinear(keywords[i%len(keywords)].kw)
	}
}

func BenchmarkQueryKeyword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T, S = queryKeyword(qrykwds[i%len(qrykwds)])
	}
}

func BenchmarkQueryKeywordMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T, S = queryKeywordMap(qrykwds[i%len(qrykwds)])
	}
}

func BenchmarkKeySwitch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T, S = keyword(keywords[i%len(keywords)].kw)
	}
}

func BenchmarkPerfect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T, S = perfect(keywords[i%len(keywords)].kw)
	}
}

func TestPerfectHash(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	fmt.Println(len(keywords), "keywords")
	hash2str := make(map[int]string)
	for _, pair := range keywords {
		kw := pair.kw
		hash := int(kw[0]) + int(kw[1])<<9 + len(kw)<<17
		assert.This(hash2str[hash]).Is("")
		hash2str[hash] = kw
	}
	fmt.Println(hash2str)
outer:
	for tblsize := len(hash2str); tblsize < 200; tblsize++ {
		tbl := make(map[int]string)
		for h, k := range hash2str {
			i := h % tblsize
			if tbl[i] != "" {
				continue outer
			}
			tbl[i] = k
		}
		fmt.Println("table size:", tblsize)
		for h, k := range tbl {
			// fmt.Print(h, `: {"`, k, `", tok.`, str.Capitalize(k), "},\n")
			fmt.Print(h, `: "`, k, "\",\n")
		}
		for h, k := range tbl {
			fmt.Print(h, `: tok.`, str.Capitalize(k), ",\n")
		}
		return
	}
	fmt.Println("failed")
}

var hashk = [128]string{
	24: "throw",
	17: "and",
	67: "super",
	49: "isnt",
	11: "continue",
	58: "new",
	10: "not",
	3:  "class",
	2:  "for",
	4:  "true",
	54: "return",
	33: "else",
	75: "switch",
	29: "is",
	12: "default",
	55: "or",
	27: "while",
	28: "do",
	23: "case",
	61: "if",
	32: "try",
	52: "this",
	74: "false",
	42: "forever",
	34: "break",
	71: "catch",
	53: "in",
	46: "function",
}

var hasht = []tok.Token{
	2:  tok.For,
	4:  tok.True,
	54: tok.Return,
	33: tok.Else,
	75: tok.Switch,
	29: tok.Is,
	3:  tok.Class,
	12: tok.Default,
	55: tok.Or,
	27: tok.While,
	28: tok.Do,
	23: tok.Case,
	61: tok.If,
	32: tok.Try,
	74: tok.False,
	42: tok.Forever,
	34: tok.Break,
	71: tok.Catch,
	53: tok.In,
	46: tok.Function,
	52: tok.This,
	24: tok.Throw,
	17: tok.And,
	67: tok.Super,
	49: tok.Isnt,
	11: tok.Continue,
	58: tok.New,
	10: tok.Not,
}

func perfect(s string) (tok.Token, string) {
	if 2 <= len(s) {
		h := (int(s[0]) + int(s[1])<<9 + len(s)<<17) % 76
		if hashk[h] == s {
			return hasht[h], hashk[h]
		}
	}
	return tok.Identifier, s
}

var qrykwds = []string{
	"by", // 1
	"in",
	"is",
	"or",
	"to",

	"and", // 1 *
	"key",
	"max",
	"min",
	"not",
	"set",

	"drop", // 2
	"into",
	"isnt",
	"join",
	"list",
	"sort",
	"true",
	"view",

	"alter", // 2 *
	"class",
	"count",
	"false",
	"index",
	"minus",
	"sview",
	"times",
	"total",
	"union",
	"where",

	"create", // 2 *
	"delete",
	"ensure",
	"extend",
	"insert",
	"remove",
	"rename",
	"unique",
	"update",

	"average", // 0
	"cascade",
	"destroy",
	"history",
	"project",
	"reverse",

	"function", // 0
	"leftjoin",

	"intersect", // 0
	"summarize",
	"tempindex",
}

func keywordLinear(s string) (tok.Token, string) {
	if 2 <= len(s) && len(s) <= 8 && s[0] >= 'a' {
		for _, pair := range keywords {
			if pair.kw == s {
				return pair.tok, pair.kw
			}
		}
	}
	return tok.Identifier, strings.Clone(s)
}

// keywords doesn't use a map because we want to reuse the keyword string literals
// ordered by frequency of use to optimize successful searches
var keywords = []struct {
	kw  string
	tok tok.Token
}{
	{"return", tok.Return},
	{"if", tok.If},
	{"false", tok.False},
	{"is", tok.Is},
	{"true", tok.True},
	{"isnt", tok.Isnt},
	{"and", tok.And},
	{"function", tok.Function},
	{"for", tok.For},
	{"in", tok.In},
	{"not", tok.Not},
	{"super", tok.Super},
	{"or", tok.Or},
	{"else", tok.Else},
	{"class", tok.Class},
	{"this", tok.This},
	{"case", tok.Case},
	{"new", tok.New},
	{"continue", tok.Continue},
	{"throw", tok.Throw},
	{"try", tok.Try},
	{"catch", tok.Catch},
	{"while", tok.While},
	{"break", tok.Break},
	{"switch", tok.Switch},
	{"default", tok.Default},
	{"do", tok.Do},
	{"forever", tok.Forever},
}

var kwds = hmap.Hmap[string, tok.Token, helper]{}

func init() {
	for _, p := range keywords {
		kwds.Put(p.kw, p.tok)
	}
}

func keywordHmap(s string) (tok.Token, string) {
	if kw, tok, ok := kwds.Get2(s); ok {
		return tok, kw
	}
	return tok.Identifier, strings.Clone(s)
}

func BenchmarkKeywordHmap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T, S = keywordHmap(keywords[i%len(keywords)].kw)
		assert.That(S != "")
	}
}

type helper struct{}

func (h helper) Hash(s string) uint32 {
	return uint32(s[0]) + uint32(s[1])<<8 + uint32(len(s))<<16
}

func (h helper) Equal(s1, s2 string) bool {
	return s1 == s2
}

var kwdStr = [][]string{
	2: {"if",
		"is",
		"in",
		"or",
		"do"},
	3: {"and",
		"for",
		"not",
		"new",
		"try"},
	4: {"true",
		"isnt",
		"else",
		"this",
		"case"},
	5: {"false",
		"super",
		"class",
		"throw",
		"catch",
		"while",
		"break"},
	6: {"return",
		"switch"},
	7: {"default",
		"forever"},
	8: {"function",
		"continue"},
}

var kwdTok = [][]tok.Token{
	2: {tok.If,
		tok.Is,
		tok.In,
		tok.Or,
		tok.Do},
	3: {tok.And,
		tok.For,
		tok.Not,
		tok.New,
		tok.Try},
	4: {tok.True,
		tok.Isnt,
		tok.Else,
		tok.This,
		tok.Case},
	5: {tok.False,
		tok.Super,
		tok.Class,
		tok.Throw,
		tok.Catch,
		tok.While,
		tok.Break},
	6: {tok.Return,
		tok.Switch},
	7: {tok.Default,
		tok.Forever},
	8: {tok.Function,
		tok.Continue},
}

func keywordLinear2(s string) (tok.Token, string) {
	n := len(s)
	switch n {
	case 2, 3, 4, 5, 6, 7, 8:
		for i, k := range kwdStr[n] {
			if k == s {
				return kwdTok[n][i], k
			}
		}
	}
	return tok.Identifier, strings.Clone(s)
}

func BenchmarkKeywordLinear2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		T, S = keywordLinear2(keywords[i%len(keywords)].kw)
	}
}

func queryKeywordMap(s string) (tok.Token, string) {
	if tok, ok := queryKeywords[str.ToLower(s)]; ok {
		return tok, s
	}
	return tok.Identifier, s
}
