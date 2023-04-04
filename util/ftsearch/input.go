// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ftsearch

import (
	"github.com/apmckinlay/gsuneido/util/ascii"
	snowballeng "github.com/kljensen/snowball/english"
)

type input struct {
	text string
	pos  int
}

func newInput(text string) *input {
	return &input{text, 0}
}

const maxTokenLength = 32

func (src *input) Next() string {
	for {
		// skip non-letter
		for src.pos < len(src.text) && !ascii.IsLetter(src.text[src.pos]) {
			src.pos++
		}
		pos := src.pos
		for src.pos < len(src.text) && ascii.IsLetter(src.text[src.pos]) {
			src.pos++
		}
		tok := src.text[pos:src.pos]
		if tok == "" {
			return ""
		}
		if len(tok) > 1 && len(tok) < maxTokenLength {
			if _, ok := stopWords[tok]; !ok {
				return snowballeng.Stem(tok, true)
			}
		}
	}
}

// stopWords is copied from snowball
var stopWords = map[string]struct{}{
	"about":      {},
	"above":      {},
	"after":      {},
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
	"before":     {},
	"being":      {},
	"below":      {},
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
	"down":       {},
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
	"off":        {},
	"on":         {},
	"once":       {},
	"only":       {},
	"or":         {},
	"other":      {},
	"our":        {},
	"ours":       {},
	"ourselves":  {},
	"out":        {},
	"over":       {},
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
	"under":      {},
	"until":      {},
	"up":         {},
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
