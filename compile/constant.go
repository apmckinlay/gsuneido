package compile

import (
	v "github.com/apmckinlay/gsuneido/value"
)

// Constant compiles a Suneido constant (e.g. a library record)
// to a Suneido Value
func Constant(src string) v.Value {
	p := newParser(src)
	return p.constant()
}

func (p *parser) constant() v.Value {
	switch p.KeyTok() {
	case FUNCTION:
		ast := p.function()
		return codegen(ast)
	default:
		panic("constant: not implemented: " + p.Value)
	}
}
