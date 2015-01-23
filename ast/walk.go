package ast

import "fmt"

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(v Visitor, list []*Ident) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkExprList(v Visitor, list []Expr) {
	for _, x := range list {
		Walk(v, x)
	}
}

// TODO(gri): Investigate if providing a closure to Walk leads to
//            simpler use (and may help eliminate Inspect in turn).

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Comments and fields
	case *Comment:
		// nothing to do

	case *CommentGroup:
		for _, c := range n.List {
			Walk(v, c)
		}

	// Expressions
	case *BadExpr, *Ident, *BasicLit:
		// nothing to do

	case *ParenExpr:
		Walk(v, n.X)

	case *UnaryExpr:
		Walk(v, n.X)

	case *BinaryExpr:
		Walk(v, n.X)
		Walk(v, n.Y)

	case *BlockExpr:
		walkExprList(v, n.List)

	case *IfExpr:
		Walk(v, n.Cond)
		Walk(v, n.Body)
		Walk(v, n.Else)

	// Files and packages
	case *File:
		if n.Doc != nil {
			Walk(v, n.Doc)
		}
		//Walk(v, n.Name)
		//Walk(v, n.List)
		for _, node := range n.List {
			Walk(v, node)
		}
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	default:
		fmt.Printf("ast.Walk: unexpected node type %T", n)
		panic("ast.Walk")
	}

	v.Visit(nil)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// for all the non-nil children of node, recursively.
//
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}
