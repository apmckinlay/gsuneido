// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package query

import (
	"strings"

	"github.com/apmckinlay/gsuneido/compile/lexer"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
)

// StripSort returns the query without a sort (if present).
// The result is undefined if the query is not valid.
func StripSort(query string) string {
	if !strings.Contains(query, "sort") {
		return query
	}
	lxr := lexer.NewQueryLexer(query)
	for {
		item := lxr.NextSkip()
		switch item.Token {
		case tok.Sort:
			if ok, pos := matchSort(lxr); ok {
				return strings.TrimRight(query[:item.Pos]+query[pos:], " ")
			}
		case tok.Eof:
			return query
		}
	}
}

func matchSort(lxr *lexer.Lexer) (ok bool, pos int) {
	start := lxr.Position()
	item := lxr.NextSkip()
	if item.Token == tok.Reverse {
		item = lxr.NextSkip()
	}
	for {
		if !item.Token.IsIdent() {
			lxr.SetPos(start)
			return
		}
		item = lxr.NextSkip()
		switch item.Token {
		case tok.Comma:
			item = lxr.NextSkip()
			//continue
		case tok.Eof, tok.Where:
			return true, int(item.Pos)
		default:
			lxr.SetPos(start)
			return
		}
	}
}

func GetSort(query string) string {
	if !strings.Contains(query, "sort") {
		return ""
	}
	lxr := lexer.NewQueryLexer(query)
	for {
		item := lxr.NextSkip()
		switch item.Token {
		case tok.Sort:
			if sort := getSort(lxr); sort != "" {
				return sort
			}
		case tok.Eof:
			return ""
		}
	}
}

func getSort(lxr *lexer.Lexer) string {
	start := lxr.Position()
	var sb strings.Builder
	item := lxr.NextSkip()
	if item.Token == tok.Reverse {
		sb.WriteString("reverse ")
		item = lxr.NextSkip()
	}
	for {
		if !item.Token.IsIdent() {
			lxr.SetPos(start)
			return ""
		}
		sb.WriteString(item.Text)
		item = lxr.NextSkip()
		switch item.Token {
		case tok.Comma:
			sb.WriteString(",")
			item = lxr.NextSkip()
			//continue
		case tok.Eof, tok.Where:
			return sb.String()
		default:
			lxr.SetPos(start)
			return ""
		}
	}
}

// JustTable returns the table name if the query is just a table name,
// ignoring whitespace and comments. Returns an empty string if the query
// contains anything other than a single table name.
func JustTable(query string) string {
	lxr := lexer.NewQueryLexer(query)
	id := ""
	for {
		item := lxr.Next()
		switch item.Token {
		case tok.Whitespace, tok.Newline, tok.Comment:
			continue
		case tok.Eof:
			return id
		default:
			if id == "" && item.Token.IsIdent() {
				id = item.Text
				continue
			}
			// If we get here, we found something other than whitespace/comments
			// after an identifier, or a non-identifier before EOF
			return ""
		}
	}
}
