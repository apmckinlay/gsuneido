// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile/lexer"
	"github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
)

type suScanner struct {
	ValueBase[*suScanner]
	MayLock
	lxr  lexer.Lexer
	item lexer.Item
	// name is either "Scanner" or "QueryScanner"
	name string
}

var _ = builtin(Scanner, "(string)")

func Scanner(arg Value) Value {
	return &suScanner{lxr: *lexer.NewLexer(ToStr(arg)), name: "Scanner"}
}

var _ Value = (*suScanner)(nil)

func (sc *suScanner) Equal(other any) bool {
	return sc == other
}

func (*suScanner) Lookup(_ *Thread, method string) Callable {
	return scannerMethods[method]
}

var scannerMethods = methods()

var _ = method(scan_KeywordQ, "()")

func scan_KeywordQ(this Value) Value {
	return SuBool(this.(*suScanner).isKeyword())
}

var _ = method(scan_Length, "()")

func scan_Length(this Value) Value {
	sc := this.(*suScanner)
	if sc.Lock() {
		defer sc.Unlock()
	}
	from := sc.item.Pos
	to := sc.lxr.Position()
	return IntVal(to - int(from))
}

var _ = method(scan_Next, "()")

func scan_Next(this Value) Value {
	return this.(*suScanner).next()
}

var _ = method(scan_Next2, "()")

func scan_Next2(this Value) Value {
	sc := this.(*suScanner)
	if sc.Lock() {
		defer sc.Unlock()
	}
	sc.item = sc.lxr.Next()
	if sc.item.Token == tokens.Eof {
		return sc
	}
	return SuStr(sc.type2())
}

var _ = method(scan_Position, "()")

func scan_Position(this Value) Value {
	sc := this.(*suScanner)
	if sc.Lock() {
		defer sc.Unlock()
	}
	return IntVal(sc.lxr.Position())
}

var _ = method(scan_Text, "()")

func scan_Text(this Value) Value {
	return this.(*suScanner).text()
}

var _ = method(scan_Type, "()")

func scan_Type(this Value) Value {
	sc := this.(*suScanner)
	if sc.Lock() {
		defer sc.Unlock()
	}
	return SuStr(sc.type2())
}

var _ = method(scan_Value, "()")

func scan_Value(this Value) Value {
	sc := this.(*suScanner)
	if sc.Lock() {
		defer sc.Unlock()
	}
	return SuStr(sc.item.Text)
}

func (sc *suScanner) next() Value {
	if sc.Lock() {
		defer sc.Unlock()
	}
	sc.item = sc.lxr.Next()
	if sc.item.Token == tokens.Eof {
		return sc
	}
	return sc.text()
}

func (sc *suScanner) text() Value {
	if sc.Lock() {
		defer sc.Unlock()
	}
	src := sc.lxr.Source()
	from := sc.item.Pos
	to := sc.lxr.Position()
	return SuStr(src[from:to])
}

// type2 caller must lock
func (sc *suScanner) type2() string {
	if sc.item.Token.IsOperator() {
		return ""
	}
	if sc.item.Token.IsIdent() {
		return "IDENTIFIER"
	}
	switch sc.item.Token {
	case tokens.Error:
		return "ERROR"
	case tokens.Identifier:
		return "IDENTIFIER"
	case tokens.Number:
		return "NUMBER"
	case tokens.String, tokens.Symbol:
		return "STRING"
	case tokens.Whitespace:
		return "WHITESPACE"
	case tokens.Comment:
		return "COMMENT"
	case tokens.Newline:
		return "NEWLINE"
	default:
		return ""
	}
}

func (sc *suScanner) isKeyword() bool {
	if sc.Lock() {
		defer sc.Unlock()
	}
	return sc.item.Token != tokens.Identifier && sc.item.Token.IsIdent()
}

// iterator ---------------------------------------------------------

func (sc *suScanner) Iter() Iter {
	return sc
}

func (sc *suScanner) Next() Value {
	if tok := sc.next(); tok != sc {
		return tok
	}
	return nil
}

func (sc *suScanner) Dup() Iter {
	if sc.Lock() {
		defer sc.Unlock()
	}
	return &suScanner{lxr: *sc.lxr.Dup()}
}

func (sc *suScanner) Infinite() bool {
	return false
}

func (sc *suScanner) SetConcurrent() {
	sc.MayLock.SetConcurrent()
}

func (sc *suScanner) Instantiate() *SuObject {
	return InstantiateIter(sc)
}

var _ Iter = (*suScanner)(nil)
