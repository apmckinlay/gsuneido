package ast

type Visitor interface {
	Before(node Node) bool // return false to skip children
	After(node Node)
}

// Traverse calls visitor.Before for node.
// If Before returns true,
// Traverse is called recursively for each child node,
// and then visitor.After is called for node.
// NOTE: it will not traverse nested functions and classes
// because they will be constants.
func Traverse(node Node, visitor Visitor) {
	if node == nil {
		return
	}
	if !visitor.Before(node) {
		return
	}
	node.Children(func(child Node) { Traverse(child, visitor) })
	visitor.After(node)
}
