package compile

import (
	"bytes"
	"strings"

	"github.com/apmckinlay/gsuneido/value"
)

// Ast is the node type for an AST returned by parse
type Ast struct {
	Item
	val      value.Value
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
		if a.Token.String() == "" && a.Value == "" && a.val == nil {
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
	} else if a.val != nil {
		buf.WriteString(a.val.String())
	} else if a.Value != "" {
		buf.WriteString(a.Value)
	}
}

func ast(item Item, children ...Ast) Ast {
	return Ast{Item: item, Children: children}
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
	return Ast{Item: item, val: val, Children: children}
}

func (a *Ast) first() Ast {
	return a.Children[0]
}

func (a *Ast) second() Ast {
	return a.Children[1]
}
