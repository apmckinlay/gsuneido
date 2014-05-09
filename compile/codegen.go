package compile

import (
	"strconv"

	i "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/util/varint"
	"github.com/apmckinlay/gsuneido/value"
)

// codegen compiles a Function from an Ast
func codegen(ast Ast) *value.SuFunc {
	//fmt.Println("codegen", ast.String())
	cg := cgen{}
	cg.gen(ast)
	cg.emit(i.RETURN)
	return &value.SuFunc{
		Code:    cg.code,
		Values:  cg.values,
		Strings: cg.strings,
		Nlocals: len(cg.strings), // ultimately WRONG!
	}
}

type cgen struct {
	code    []byte
	values  []value.Value
	strings []string
}

func (cg *cgen) gen(ast Ast) {
	//fmt.Println("gen", ast.String())
	switch ast.KeyTok() {
	case FUNCTION:
		// TODO params
		cg.gen(ast.Children[1]) // statements
	case STATEMENTS:
		for _, stmt := range ast.Children {
			cg.gen(stmt)
		}
	case EXPRESSION:
		cg.gen(ast.Children[0])
	case NUMBER:
		n, err := strconv.ParseInt(ast.Value, 0, 32)
		if err == nil {
			cg.emit(i.PUSHINT)
			cg.code = varint.EncodeInt32(int32(n), cg.code)
		} else {
			v, err := value.ParseNum(ast.Value)
			if err != nil {
				panic("invalid number: " + ast.Value)
			}
			cg.value(v)
		}
	case STRING:
		// TODO: copy so no ref to source
		cg.emit(i.PUSHVAL)
		i := cg.value(value.SuStr(ast.Value))
		cg.emitUint(i)
	case ADD:
		cg.ubinop(ast, i.UPLUS, i.ADD)
	case SUB:
		cg.ubinop(ast, i.UMINUS, i.SUB)
	case CAT:
		cg.binop(ast, i.CAT)
	case MUL:
		cg.binop(ast, i.MUL)
	case DIV:
		cg.binop(ast, i.DIV)
	case MOD:
		cg.binop(ast, i.MOD)
	case IDENTIFIER:
		if isLocal(ast.Value) {
			cg.emit(i.LOAD)
			cg.emitUint(cg.local(ast.Value))
		} else {
			panic("not implemented")
		}
	case EQ:
		id := ast.Children[0]
		if isLocal(id.Value) {
			cg.gen(ast.Children[1]) // expr
			cg.emit(i.STORE)
			cg.emitUint(cg.local(id.Value))
		} else {
			panic("not implemented")
		}
	default:
		panic("not implemented" + ast.String())
	}
}

// value returns an index for the value
// reusing if duplicate, adding otherwise
func (cg *cgen) value(v value.Value) int {
	for i, v2 := range cg.values {
		if v.Equals(v2) {
			return i
		}
	}
	i := len(cg.values)
	cg.values = append(cg.values, v)
	return i
}

func isLocal(s string) bool {
	return 'a' <= s[0] && s[0] <= 'z'
}

// local returns the index for a local variable
func (cg *cgen) local(s string) int {
	for i, s2 := range cg.strings {
		if s == s2 {
			return i
		}
	}
	i := len(cg.strings)
	// TODO intern to avoid ref to source
	cg.strings = append(cg.strings, s)
	return i
}

// ubinop is called for ops that can be unary or binary
func (cg *cgen) ubinop(ast Ast, uop, bop byte) {
	if len(ast.Children) == 1 {
		cg.unop(ast, uop)
	} else {
		cg.binop(ast, bop)
	}
}

func (cg *cgen) unop(ast Ast, op byte) {
	cg.gen(ast.Children[0])
	cg.emit(op)
}

func (cg *cgen) binop(ast Ast, op byte) {
	cg.gen(ast.Children[0])
	cg.gen(ast.Children[1])
	cg.emit(op)
}

func (cg *cgen) emit(b ...byte) {
	cg.code = append(cg.code, b...)
}

func (cg *cgen) emitUint(i int) {
	cg.code = varint.EncodeUint32(uint32(i), cg.code)
}
