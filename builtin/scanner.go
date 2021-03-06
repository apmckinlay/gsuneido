// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"github.com/apmckinlay/gsuneido/compile/lexer"
	"github.com/apmckinlay/gsuneido/compile/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type suScanner struct {
	CantConvert
	MayLock
	lxr  lexer.Lexer
	item lexer.Item
	// name is either "Scanner" or "QueryScanner"
	name string
}

var _ = builtin1("Scanner(string)",
	func(arg Value) Value {
		return &suScanner{lxr: *lexer.NewLexer(ToStr(arg)), name: "Scanner"}
	})

var _ Value = (*suScanner)(nil)

func (sc *suScanner) Get(*Thread, Value) Value {
	panic(sc.name + " does not support get")
}

func (sc *suScanner) Put(*Thread, Value, Value) {
	panic(sc.name + " does not support put")
}

func (sc *suScanner) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic(sc.name + " does not support update")
}

func (sc *suScanner) RangeTo(int, int) Value {
	panic(sc.name + " does not support range")
}

func (sc *suScanner) RangeLen(int, int) Value {
	panic(sc.name + " does not support range")
}

func (sc *suScanner) Hash() uint32 {
	panic(sc.name + " hash not implemented")
}

func (sc *suScanner) Hash2() uint32 {
	panic(sc.name + " hash not implemented")
}

func (sc *suScanner) Compare(Value) int {
	panic(sc.name + " compare not implemented")
}

func (sc *suScanner) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call " + sc.name)
}

func (sc *suScanner) String() string {
	return "a" + sc.name
}

func (*suScanner) Type() types.Type {
	return types.BuiltinInstance
}

func (sc *suScanner) Equal(other interface{}) bool {
	sc2, ok := other.(*suScanner)
	return ok && sc == sc2
}

func (*suScanner) Lookup(_ *Thread, method string) Callable {
	return scannerMethods[method]
}

var scannerMethods = Methods{
	"Keyword?": method0(func(this Value) Value {
		return SuBool(this.(*suScanner).isKeyword())
	}),
	"Length": method0(func(this Value) Value {
		sc := this.(*suScanner)
		if sc.Lock() {
			defer sc.Unlock()
		}
		from := sc.item.Pos
		to := sc.lxr.Position()
		return IntVal(to - int(from))
	}),
	"Next": method0(func(this Value) Value {
		return this.(*suScanner).next()
	}),
	"Next2": method0(func(this Value) Value {
		sc := this.(*suScanner)
		if sc.Lock() {
			defer sc.Unlock()
		}
		sc.item = sc.lxr.Next()
		if sc.item.Token == tokens.Eof {
			return sc
		}
		return SuStr(sc.type2())
	}),
	"Position": method0(func(this Value) Value {
		sc := this.(*suScanner)
		if sc.Lock() {
			defer sc.Unlock()
		}
		return IntVal(sc.lxr.Position())
	}),
	"Text": method0(func(this Value) Value {
		return this.(*suScanner).text()
	}),
	"Type": method0(func(this Value) Value {
		sc := this.(*suScanner)
		if sc.Lock() {
			defer sc.Unlock()
		}
		return SuStr(sc.type2())
	}),
	"Value": method0(func(this Value) Value {
		sc := this.(*suScanner)
		if sc.Lock() {
			defer sc.Unlock()
		}
		return SuStr(sc.item.Text)
	}),
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
