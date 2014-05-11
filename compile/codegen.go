package compile

import (
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
		stmts := ast.Children[1].Children
		for si, stmt := range stmts {
			cg.gen(stmt)
			if si < len(stmts)-1 {
				cg.emit(i.POP)
			}
		}
	case EXPRESSION:
		cg.gen(ast.Children[0])
	case NUMBER:
		val, err := value.NumFromString(ast.Value)
		if err != nil {
			panic("invalid number: " + ast.Value)
		}
		if si, ok := val.(value.SuInt); ok {
			cg.emit(i.INT)
			cg.code = varint.EncodeInt32(int32(si), cg.code)
		} else {
			cg.emit(i.VALUE)
			cg.emitUint(cg.value(val))
		}
	case STRING:
		// TODO: copy so no ref to source
		cg.emit(i.VALUE)
		i := cg.value(value.SuStr(ast.Value))
		cg.emitUint(i)
	case IS:
		cg.binop(ast, i.IS)
	case ISNT:
		cg.binop(ast, i.ISNT)
	case LT:
		cg.binop(ast, i.LT)
	case LTE:
		cg.binop(ast, i.LTE)
	case GT:
		cg.binop(ast, i.GT)
	case GTE:
		cg.binop(ast, i.GTE)
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
	case LSHIFT:
		cg.binop(ast, i.LSHIFT)
	case RSHIFT:
		cg.binop(ast, i.RSHIFT)
	case BITOR:
		cg.binop(ast, i.BITOR)
	case BITAND:
		cg.binop(ast, i.BITAND)
	case BITXOR:
		cg.binop(ast, i.BITXOR)
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
	// TODO merge pop with previous instruction
	// TODO detect pop after side effect free instruction
	cg.code = append(cg.code, b...)
}

func (cg *cgen) emitUint(i int) {
	cg.code = varint.EncodeUint32(uint32(i), cg.code)
}
