package builtin

import (
	"github.com/apmckinlay/gsuneido/lexer"
	"github.com/apmckinlay/gsuneido/lexer/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
)

type SuScanner struct {
	CantConvert
	lxr  lexer.Lexer
	item lexer.Item
	name string
}

var _ = builtin1("Scanner(string)",
	func(arg Value) Value {
		return &SuScanner{lxr: *lexer.NewLexer(ToStr(arg)), name: "Scanner"}
	})

var _ Value = (*SuScanner)(nil)

func (sc *SuScanner) Get(*Thread, Value) Value {
	panic(sc.name + " does not support get")
}

func (sc *SuScanner) Put(*Thread, Value, Value) {
	panic(sc.name + " does not support put")
}

func (sc *SuScanner) RangeTo(int, int) Value {
	panic(sc.name + " does not support range")
}

func (sc *SuScanner) RangeLen(int, int) Value {
	panic(sc.name + " does not support range")
}

func (sc *SuScanner) Hash() uint32 {
	panic(sc.name + " hash not implemented")
}

func (sc *SuScanner) Hash2() uint32 {
	panic(sc.name + " hash not implemented")
}

func (sc *SuScanner) Compare(Value) int {
	panic(sc.name + " compare not implemented")
}

func (sc *SuScanner) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call " + sc.name)
}

func (sc *SuScanner) String() string {
	return "a" + sc.name
}

func (*SuScanner) Type() types.Type {
	return types.BuiltinInstance
}

func (sc *SuScanner) Equal(other interface{}) bool {
	if sc2, ok := other.(*SuScanner); ok {
		return sc == sc2
	}
	return false
}

func (*SuScanner) Lookup(_ *Thread, method string) Callable {
	return scannerMethods[method]
}

var scannerMethods = Methods{
	"Keyword?": method0(func(this Value) Value {
		return SuBool(this.(*SuScanner).isKeyword())
	}),
	"Length": method0(func(this Value) Value {
		sc := this.(*SuScanner)
		from := sc.item.Pos
		to := sc.lxr.Position()
		return IntVal(to - int(from))
	}),
	"Next": method0(func(this Value) Value {
		return this.(*SuScanner).next()
	}),
	"Next2": method0(func(this Value) Value {
		sc := this.(*SuScanner)
		sc.item = sc.lxr.Next()
		if sc.item.Token == tokens.Eof {
			return sc
		}
		return SuStr(sc.type2())
	}),
	"Position": method0(func(this Value) Value {
		return IntVal(this.(*SuScanner).lxr.Position())
	}),
	"Text": method0(func(this Value) Value {
		return this.(*SuScanner).text()
	}),
	"Type": method0(func(this Value) Value {
		return SuStr(this.(*SuScanner).type2())
	}),
	// TODO remove after everyone has switched to new Type
	"Type2": method0(func(this Value) Value {
		return SuStr(this.(*SuScanner).type2())
	}),
	"Value": method0(func(this Value) Value {
		return SuStr(this.(*SuScanner).item.Text)
	}),
}

func (sc *SuScanner) next() Value {
	sc.item = sc.lxr.Next()
	if sc.item.Token == tokens.Eof {
		return sc
	}
	return sc.text()
}

func (sc *SuScanner) text() Value {
	src := sc.lxr.Source()
	from := sc.item.Pos
	to := sc.lxr.Position()
	return SuStr(src[from:to])
}

func (sc *SuScanner) type2() string {
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
	case tokens.String:
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

func (sc *SuScanner) isKeyword() bool {
	return sc.item.Token != tokens.Identifier && sc.item.Token.IsIdent()
}

// iterator ---------------------------------------------------------

func (sc *SuScanner) Iter() Iter {
	return sc
}

func (sc *SuScanner) Next() Value {
	if tok := sc.next(); tok != sc {
		return tok
	}
	return nil
}

func (sc *SuScanner) Dup() Iter {
	return &SuScanner{lxr: *sc.lxr.Dup()}
}

func (sc *SuScanner) Infinite() bool {
	return false
}
