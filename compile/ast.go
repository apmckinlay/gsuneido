package compile

import (
	"bytes"
	"strings"

	. "github.com/apmckinlay/gsuneido/base"
	. "github.com/apmckinlay/gsuneido/lexer"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// Ast is the node type for an AST returned by parse
type Ast struct {
	Item
	value    Value
	Children []Ast
}

// String formats a tree of Ast's in a relatively compact form
func (a Ast) String() string {
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
		if a.Text != "" {
			buf.WriteString("(" + a.Text + " " + a.value.String() + ")")
		} else {
			buf.WriteString(a.value.String())
		}
	} else if a.Text != "" {
		buf.WriteString(a.Text)
	}
}

func ast(item Item, children ...Ast) Ast {
	return fold(item, nil, children)
}

func ast2(name string, children ...Ast) Ast {
	return fold(Item{Text: name}, nil, children)
}

func astVal(name string, val Value) Ast {
	return fold(Item{Text: name}, val, []Ast{})
}

func astBuilder(item Item, nodes ...T) T {
	var val Value
	if len(nodes) >= 1 {
		if v, ok := nodes[0].(Value); ok {
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

func (a *Ast) fourth() Ast {
	return a.Children[3]
}

func fold(item Item, val Value, children []Ast) (x Ast) {
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
		return ast.commutative(Add, SuInt(0))
	case SUB:
		val = ast.unop(Uminus)
	case IS:
		val = ast.binop(Is)
	case ISNT:
		val = ast.binop(Isnt)
	case LT:
		val = ast.binop(Lt)
	case LTE:
		val = ast.binop(Lte)
	case GT:
		val = ast.binop(Gt)
	case GTE:
		val = ast.binop(Gte)
	case CAT:
		return ast.foldCat()
	case MUL:
		return ast.commutative(Mul, SuInt(1))
	case MOD:
		val = ast.binop(Mod)
	case LSHIFT:
		val = ast.binop(Lshift)
	case RSHIFT:
		val = ast.binop(Rshift)
	case BITOR:
		val = ast.binop(Bitor)
	case BITAND:
		val = ast.binop(Bitand)
	case BITXOR:
		val = ast.binop(Bitxor)
	case BITNOT:
		val = ast.unop(Bitnot)
	case NOT:
		val = ast.unop(Not)
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
	}
	cc := countConstant(a.Children)
	return cc == len(a.Children) || cc >= 2
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

type uopfn func(Value) Value
type bopfn func(Value, Value) Value

func (a *Ast) unop(uop uopfn) Value {
	verify.That(len(a.Children) == 1)
	return uop(a.Children[0].toVal())
}

func (a *Ast) binop(bop bopfn) Value {
	verify.That(len(a.Children) == 2)
	return bop(a.Children[0].toVal(), a.Children[1].toVal())
}

// for add and mul
func (a *Ast) commutative(bop bopfn, identity Value) Ast {
	k := identity
	i := 0
	for _, c := range a.Children {
		if c.Token == DIV && c.Children[0].isConstant() {
			k = Div(k, c.Children[0].toVal())
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
	empty := SuStr("")
	var k Value = empty
	i := 0
	for _, c := range a.Children {
		if c.isConstant() {
			k = Cat(k, c.toVal())
		} else {
			if k != empty {
				k = SuStr(k.ToStr()) // ensure not Concat
				a.Children[i] = valAst(k)
				k = empty
				i++
			}
			a.Children[i] = c
			i++
		}
	}
	k = SuStr(k.ToStr()) // ensure not Concat
	if i == 0 {          // all constant
		return valAst(k)
	} else if k != empty {
		a.Children[i] = valAst(k)
		i++
	}
	a.Children = a.Children[:i]
	return *a
}

func valAst(val Value) Ast {
	return Ast{value: val}
}

func (a *Ast) isConstant() bool {
	switch a.KeyTok() {
	case NUMBER, STRING, TRUE, FALSE:
		return true
	default:
		return a.value != nil && a.Text == ""
	}
}

func (a *Ast) toVal() Value {
	if a.value != nil {
		return a.value
	}
	switch a.KeyTok() {
	case NUMBER:
		return NumFromString(a.Text)
	case STRING:
		return SuStr(a.Text)
	case TRUE:
		return True
	case FALSE:
		return False
	default:
		panic("bad toVal")
	}
}
