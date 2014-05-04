package codegen

import (
	"fmt"
	"strconv"

	. "github.com/apmckinlay/gsuneido/compile/lex"
	"github.com/apmckinlay/gsuneido/compile/parse"
	i "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/util/varint"
	"github.com/apmckinlay/gsuneido/value"
)

func Codegen(ast parse.AstNode) *i.Function {
	fmt.Println("codegen", ast.String())
	cg := cgen{}
	cg.gen(ast)
	cg.emit(i.RETURN)
	return &i.Function{Code: cg.code, Values: cg.values}
}

type cgen struct {
	code   []byte
	values []value.Value
}

func (cg *cgen) gen(ast parse.AstNode) {
	fmt.Println("gen", ast.String())
	switch ast.Token {
	case NUMBER:
		n, _ := strconv.ParseInt(ast.Value, 10, 32)
		cg.emit(i.PUSHINT)
		cg.code = varint.EncodeInt32(int32(n), cg.code)
	case STRING:
		cg.emit(i.PUSHVAL)
		i := len(cg.values)
		cg.values = append(cg.values, value.StrVal(ast.Value))
		cg.code = varint.EncodeUint32(uint32(i), cg.code)
	case ADD:
		cg.binop(ast, i.ADD)
	case SUB:
		cg.binop(ast, i.SUB)
	case CAT:
		cg.binop(ast, i.CAT)
	case MUL:
		cg.binop(ast, i.MUL)
	case DIV:
		cg.binop(ast, i.DIV)
	case MOD:
		cg.binop(ast, i.MOD)
	default:
		panic("not implemented")
	}
}

func (cg *cgen) binop(ast parse.AstNode, op byte) {
	cg.gen(ast.Children[0])
	cg.gen(ast.Children[1])
	cg.emit(op)
}

func (cg *cgen) emit(b ...byte) {
	cg.code = append(cg.code, b...)
}
