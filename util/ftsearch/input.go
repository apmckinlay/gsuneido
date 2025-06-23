// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"github.com/apmckinlay/gsuneido/util/ascii"
	"github.com/apmckinlay/gsuneido/util/str"
	snowballeng "github.com/kljensen/snowball/english"
)

type input struct {
	text   string
	pos    int
	litpos int
}

func NewInput(text string) *input {
	return &input{text: text}
}

const minWordLength = 2
const maxWordLength = 32
const maxNumberLength = 16

var alphanum = str.MakeSet("a-zA-Z0-9")
var special = str.MakeSet("-_&")

func (src *input) Next() string {
	for {
		// skip non-alphanumeric
		for src.pos < len(src.text) && !alphanum.Contains(src.text[src.pos]) {
			src.pos++
		}
		if src.pos >= len(src.text) {
			return ""
		}
		pos := src.pos
		if src.pos >= src.litpos {
			i := pos
			var let, num, spec int
			for i < len(src.text) {
				if ascii.IsLetter(src.text[i]) {
					let = 1
				} else if ascii.IsDigit(src.text[i]) {
					num = 1
				} else if special.Contains(src.text[i]) {
					spec = 1
				} else {
					break
				}
				i++
			}
			if let+num+spec > 1 && i-pos > 1 {
				src.litpos = i + 1
				if i-pos <= maxWordLength {
					return str.ToLower(src.text[pos:i])
				}
			}
		}
		if ascii.IsLetter(src.text[src.pos]) {
			for src.pos < len(src.text) && ascii.IsLetter(src.text[src.pos]) {
				src.pos++
			}
			tok := src.text[pos:src.pos]
			if len(tok) >= minWordLength && len(tok) < maxWordLength {
				if _, ok := stopWords[str.ToLower(tok)]; !ok {
					return snowballeng.Stem(tok, true)
				}
			}
		} else if ascii.IsDigit(src.text[src.pos]) {
			for src.pos < len(src.text) && ascii.IsDigit(src.text[src.pos]) {
				src.pos++
			}
			tok := src.text[pos:src.pos]
			if len(tok) >= minWordLength && len(tok) < maxNumberLength {
				return tok
			}
		}
	}
}

// stopWords is copied from snowball
var stopWords = map[string]struct{}{
	"about":      {},
	"again":      {},
	"against":    {},
	"all":        {},
	"am":         {},
	"an":         {},
	"and":        {},
	"any":        {},
	"are":        {},
	"as":         {},
	"at":         {},
	"be":         {},
	"because":    {},
	"been":       {},
	"being":      {},
	"between":    {},
	"both":       {},
	"but":        {},
	"by":         {},
	"can":        {},
	"did":        {},
	"do":         {},
	"does":       {},
	"doing":      {},
	"don":        {},
	"during":     {},
	"each":       {},
	"few":        {},
	"for":        {},
	"from":       {},
	"further":    {},
	"had":        {},
	"has":        {},
	"have":       {},
	"having":     {},
	"he":         {},
	"her":        {},
	"here":       {},
	"hers":       {},
	"herself":    {},
	"him":        {},
	"himself":    {},
	"his":        {},
	"how":        {},
	"if":         {},
	"in":         {},
	"into":       {},
	"is":         {},
	"it":         {},
	"its":        {},
	"itself":     {},
	"just":       {},
	"me":         {},
	"more":       {},
	"most":       {},
	"my":         {},
	"myself":     {},
	"no":         {},
	"nor":        {},
	"not":        {},
	"now":        {},
	"of":         {},
	"once":       {},
	"only":       {},
	"or":         {},
	"other":      {},
	"our":        {},
	"ours":       {},
	"ourselves":  {},
	"out":        {},
	"own":        {},
	"same":       {},
	"she":        {},
	"should":     {},
	"so":         {},
	"some":       {},
	"such":       {},
	"than":       {},
	"that":       {},
	"the":        {},
	"their":      {},
	"theirs":     {},
	"them":       {},
	"themselves": {},
	"then":       {},
	"there":      {},
	"these":      {},
	"they":       {},
	"this":       {},
	"those":      {},
	"through":    {},
	"to":         {},
	"too":        {},
	"until":      {},
	"very":       {},
	"was":        {},
	"we":         {},
	"were":       {},
	"what":       {},
	"when":       {},
	"where":      {},
	"which":      {},
	"while":      {},
	"who":        {},
	"whom":       {},
	"why":        {},
	"will":       {},
	"with":       {},
	"you":        {},
	"your":       {},
	"yours":      {},
	"yourself":   {},
	"yourselves": {},
}
