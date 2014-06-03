package compile

import (
	"bytes"
	"strings"

	"github.com/apmckinlay/gsuneido/util/verify"
	"github.com/apmckinlay/gsuneido/value"
)

// Ast is the node type for an AST returned by parse
type Ast struct {
	Item
	value    value.Value
	Children []Ast
}

// String formats a tree of Ast's in a relatively compact form
func (a *Ast) String() string {
	return string(a.bytes(0))
}

const maxline = 60 // allow for indenting

func (a *Ast) bytes(indent int) []byte {
	buf := bytes.Buffer{}
	if len(a.Children) == 0 {
		if a.Token.String() == "" && a.Text == "" && a.value == nil {
			buf.WriteString("()")
		} else {
			a.tokval(&buf)
		}
	} else {
		n := 0
		children := [][]byte{}
		for _, child := range a.Children {
			c := child.bytes(indent + 4)
			if bytes.IndexByte(c, '\n') != -1 {
				n = maxline
			} else {
				n += len(c) + 1
			}
			children = append(children, c)
		}
		if n < maxline {
			buf.WriteString("(")
			a.tokval(&buf)
			buf.WriteString(" ")
			buf.Write(bytes.Join(children, []byte(" ")))
		} else {
			buf.WriteString(strings.Repeat(" ", indent))
			buf.WriteString("(")
			a.tokval(&buf)
			sin := strings.Repeat(" ", indent+4)
			for _, c := range children {
				buf.WriteByte('\n')
				if bytes.IndexByte(c, '\n') == -1 {
					buf.WriteString(sin)
				}
				buf.Write(c)
			}
		}
		buf.WriteString(")")
	}
	return buf.Bytes()
}

func (a *Ast) tokval(buf *bytes.Buffer) {
	if ts := a.Token.String(); ts != "" {
		buf.WriteString(a.Token.String())
	} else if a.value != nil {
		buf.WriteString(a.value.String())
	} else if a.Text != "" {
		buf.WriteString(a.Text)
	}
}

func ast(item Item, children ...Ast) Ast {
	return fold(item, nil, children)
}

func ast2(name string, children ...Ast) Ast {
	return fold(Item{Token: INTERNAL, Text: name}, nil, children)
}

func astBuilder(item Item, nodes ...T) T {
	var val value.Value
	if len(nodes) >= 1 {
		if v, ok := nodes[0].(value.Value); ok {
			val = v
			nodes = nodes[1:]
		}
	}
	children := []Ast{}
	for _, node := range nodes {
		children = append(children, node.(Ast))
	}
	return fold(item, val, children)
}

func (a *Ast) first() Ast {
	return a.Children[0]
}

func (a *Ast) second() Ast {
	return a.Children[1]
}

func (a *Ast) third() Ast {
	return a.Children[2]
}

func fold(item Item, val value.Value, children []Ast) (x Ast) {
	//defer func() { fmt.Println("fold:", x) }()
	ast := Ast{item, val, children}
	if ast.isConstant() {
		return valAst(ast.toVal())
	}
	if !ast.foldable() {
		return ast
	}
	switch item.KeyTok() {
	case ADD:
		return ast.commutative(value.Add, value.SuInt(0))
	case SUB:
		val = ast.unop(value.Uminus)
	case IS:
		val = ast.binop(value.Is)
	case ISNT:
		val = ast.binop(value.Isnt)
	case LT:
		val = ast.binop(value.Lt)
	case LTE:
		val = ast.binop(value.Lte)
	case GT:
		val = ast.binop(value.Gt)
	case GTE:
		val = ast.binop(value.Gte)
	case CAT:
		return ast.foldCat()
	case MUL:
		return ast.commutative(value.Mul, value.SuInt(1))
	case MOD:
		val = ast.binop(value.Mod)
	case LSHIFT:
		val = ast.binop(value.Lshift)
	case RSHIFT:
		val = ast.binop(value.Rshift)
	case BITOR:
		val = ast.binop(value.Bitor)
	case BITAND:
		val = ast.binop(value.Bitand)
	case BITXOR:
		val = ast.binop(value.Bitxor)
	case BITNOT:
		val = ast.unop(value.Bitnot)
	case NOT:
		val = ast.unop(value.Not)
	default:
		return ast
	}
	return valAst(val)
}

func (a *Ast) foldable() bool {
	if len(a.Children) == 0 {
		return false
	}
	if a.Token == CAT {
		prev := false
		for _, c := range a.Children {
			cur := c.isConstant()
			if cur && prev {
				return true
			}
			prev = cur
		}
		return false
	} else {
		cc := countConstant(a.Children)
		return cc == len(a.Children) || cc >= 2
	}
}

func countConstant(children []Ast) int {
	n := 0
	for _, c := range children {
		if c.isConstant() ||
			(c.Token == DIV && c.Children[0].isConstant()) {
			n++
		}
	}
	return n
}

type uopfn func(value.Value) value.Value
type bopfn func(value.Value, value.Value) value.Value

func (a *Ast) unop(uop uopfn) value.Value {
	verify.That(len(a.Children) == 1)
	return uop(a.Children[0].toVal())
}

func (a *Ast) binop(bop bopfn) value.Value {
	verify.That(len(a.Children) == 2)
	return bop(a.Children[0].toVal(), a.Children[1].toVal())
}

func (a *Ast) ubop(uop uopfn, bop bopfn) value.Value {
	if len(a.Children) == 1 {
		return uop(a.Children[0].toVal())
	} else {
		result := a.Children[0].toVal()
		for _, c := range a.Children[1:] {
			result = bop(result, c.toVal())
		}
		return result
	}
}

// for add and mul
func (a *Ast) commutative(bop bopfn, identity value.Value) Ast {
	k := identity
	i := 0
	for _, c := range a.Children {
		if c.Token == DIV && c.Children[0].isConstant() {
			k = value.Div(k, c.Children[0].toVal())
		} else if c.isConstant() {
			k = bop(k, c.toVal())
		} else {
			a.Children[i] = c
			i++
		}
	}
	if i == 0 { // all constant
		return valAst(k)
	}
	a.Children[i] = valAst(k)
	a.Children = a.Children[:i+1]
	return *a
}

// cat is not commutative
func (a *Ast) foldCat() Ast {
	empty := value.SuStr("")
	var k value.Value = empty
	i := 0
	for _, c := range a.Children {
		if c.isConstant() {
			k = value.Cat(k, c.toVal())
		} else {
			if k != empty {
				k = value.SuStr(k.ToStr()) // ensure not Concat
				a.Children[i] = valAst(k)
				k = empty
				i++
			}
			a.Children[i] = c
			i++
		}
	}
	k = value.SuStr(k.ToStr()) // ensure not Concat
	if i == 0 {                // all constant
		return valAst(k)
	} else if k != empty {
		a.Children[i] = valAst(k)
		i++
	}
	a.Children = a.Children[:i]
	return *a
}

func valAst(val value.Value) Ast {
	return Ast{Item: Item{Token: VALUE}, value: val}
}

func (a *Ast) isConstant() bool {
	switch a.KeyTok() {
	case NUMBER, STRING, TRUE, FALSE, VALUE:
		return true
	default:
		return false
	}
}

func (a *Ast) toVal() value.Value {
	switch a.KeyTok() {
	case NUMBER:
		val, err := value.NumFromString(a.Text)
		if err != nil {
			panic("invalid number: " + a.Text)
		}
		return val
	case STRING:
		return value.SuStr(a.Text)
	case TRUE:
		return value.True
	case FALSE:
		return value.False
	case VALUE:
		return a.value
	default:
		panic("bad toVal")
	}
}
