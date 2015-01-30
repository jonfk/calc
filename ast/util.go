package ast

import (
	"bytes"
	"fmt"
)

// --------------------------------------------------------------------------------
// Utility Functions

// InsertExpr inserts an expression into an expression tree
// if the expression insertion is invalid it returns nil with an error
// An e.g. of an invalid expression insertion is Ident into a BasicLit
// Expressions have to be inserted in the order they would be encountered
func InsertExpr(tree, expr Expr) (Expr, error) {
	switch tree.(type) {
	case nil:
		return expr, nil
	case *Ident:
		switch t := expr.(type) {
		case *BinaryExpr:
			return insertBinaryExpr(tree, expr.(*BinaryExpr))
		default:
			return nil, fmt.Errorf("InsertExpr: cannot insert expr with type: %T into an Ident", t)
		}
	case *BasicLit:
		switch t := expr.(type) {
		case *BinaryExpr:
			return insertBinaryExpr(tree, expr.(*BinaryExpr))
		default:
			return nil, fmt.Errorf("InsertExpr: cannot insert expr with type: %T into a BasicLit", t)
		}
	case *ParenExpr:
		treeP := tree.(*ParenExpr)
		if treeP.X == nil && treeP.Rparen.Val == "" {
			tree.(*ParenExpr).X = expr
			return tree, nil
		} else if treeP.X != nil && treeP.Rparen.Val == "" {
			var err error
			tree.(*ParenExpr).X, err = InsertExpr(tree.(*ParenExpr).X, expr)
			return tree, err
		} else if treeP.Rparen.Val == ")" {
			switch t := expr.(type) {
			case *UnaryExpr:
				exprU := expr.(*UnaryExpr)
				exprU.X = tree
				return exprU, nil
			case *BinaryExpr:
				return insertBinaryExpr(tree, expr.(*BinaryExpr))
			default:
				return nil, fmt.Errorf("InsertExpr: cannot insert expr with type: %T into a closed ParenExpr", t)
			}
		}
	case *UnaryExpr:
		treeU := tree.(*UnaryExpr)
		if treeU.X == nil {
			treeU.X = expr
			return treeU, nil
		} else {
			switch expr.(type) {
			case *BinaryExpr:
				exprB := expr.(*BinaryExpr)
				if exprB.Op.Precedence() >= treeU.Op.Precedence() {
					var err error
					treeU.X, err = InsertExpr(treeU.X, expr)
					return treeU, err
				} else {
					exprB.X = tree
					return exprB, nil
				}
			default:
				var err error
				treeU.X, err = InsertExpr(treeU.X, expr)
				return treeU, err
			}
		}
	case *BinaryExpr:
		return insertIntoBinaryExpr(tree.(*BinaryExpr), expr)
	default:
		return nil, nil
	}
	return nil, nil
}

func insertBinaryExpr(tree Expr, expr *BinaryExpr) (Expr, error) {
	switch tree.(type) {
	case *BinaryExpr:
		return insertIntoBinaryExpr(tree.(*BinaryExpr), expr)
	default:
		expr.X = tree
		return expr, nil
	}
	return nil, nil

}

func insertIntoBinaryExpr(tree *BinaryExpr, expr Expr) (Expr, error) {
	if tree.X == nil {
		return nil, fmt.Errorf("InsertIntoBinaryExpr: Wrong insertion order. X cannot be nil")
	} else if tree.X != nil && tree.Y == nil {
		tree.Y = expr
		return tree, nil
	} else if tree.X != nil && tree.Y != nil {
		switch expr.(type) {
		case *BinaryExpr:
			exprB := expr.(*BinaryExpr)
			if tree.Op.Precedence() >= exprB.Op.Precedence() && !unclosedParen(tree.Y) {
				exprB.X = tree
				return exprB, nil
			} else {
				var err error
				tree.Y, err = InsertExpr(tree.Y, expr)
				return tree, err
			}
		default:
			var err error
			tree.Y, err = InsertExpr(tree.Y, expr)
			return tree, err
		}
	} else {
		return nil, fmt.Errorf("InsertExpr: Internal Error when inserting into BinaryExpr")
	}
}

func unclosedParen(tree Expr) bool {
	switch tree.(type) {
	case *ParenExpr:
		treeP := tree.(*ParenExpr)
		if treeP.Rparen.Val == "" {
			return true
		}
	case *UnaryExpr:
		return unclosedParen(tree.(*UnaryExpr).X)
	case *BinaryExpr:
		return unclosedParen(tree.(*BinaryExpr).Y)
	default:
		return false
	}
	return false
}

// Does deep comparison of Nodes
// Compares values of nodes and not position
// Used for testing
func Equals(a, b Node) bool {
	switch av := a.(type) {
	case *Ident:
		switch bv := b.(type) {
		case *Ident:
			if bv.Tok.Equals(av.Tok) {
				return true
			}
			return false
		}
	case *BasicLit:
		switch bv := b.(type) {
		case *BasicLit:
			if bv.Tok.Equals(av.Tok) {
				return true
			}
			return false
		}
	case *ParenExpr:
		switch bv := b.(type) {
		case *ParenExpr:
			return av.Lparen.Equals(bv.Lparen) &&
				av.Rparen.Equals(bv.Rparen) &&
				Equals(av.X, bv.X)
		}
	case *UnaryExpr:
		switch bv := b.(type) {
		case *UnaryExpr:
			return av.Op.Equals(bv.Op) &&
				Equals(av.X, bv.X)
		}
	case *BinaryExpr:
		switch bv := b.(type) {
		case *BinaryExpr:
			return av.Op.Equals(bv.Op) &&
				Equals(av.Y, bv.Y) &&
				Equals(av.X, bv.X)
		}
	case *ExprStmt:
		switch bv := b.(type) {
		case *ExprStmt:
			return Equals(av.X, bv.X)
		}
	case *AssignStmt:
		switch bv := b.(type) {
		case *AssignStmt:
			return av.Tok.Equals(bv.Tok) &&
				Equals(av.Lhs, bv.Lhs) &&
				Equals(av.Rhs, bv.Rhs)
		}
	case *DeclStmt:
		switch bv := b.(type) {
		case *DeclStmt:
			return Equals(av.Decl, bv.Decl)
		}
	case *ValueSpec:
		switch bv := b.(type) {
		case *ValueSpec:
			return Equals(av.Name, bv.Name) &&
				Equals(av.Value, bv.Value)
		}
	case *GenDecl:
		switch bv := b.(type) {
		case *GenDecl:
			return av.Tok.Equals(bv.Tok) &&
				Equals(av.Spec, bv.Spec)
		}
	case *File:
		switch bv := b.(type) {
		case *File:
			if len(av.List) == len(bv.List) {
				for i := range av.List {
					if !Equals(av.List[i], bv.List[i]) {
						return false
					}
				}
				return true
			}
			return false
		}
	default:
		return false
	}
	return false
}

func Sprint(n Node) string {
	return sprintd(n, 0)
}

// sprintd print prints the node to standard output with the proper depth
// helper function to improve formatting
func sprintd(n Node, d int) string {
	switch nt := n.(type) {
	case *Ident:
		return nt.String()
	case *BasicLit:
		return nt.String()
	case *ParenExpr:
		return nt.StringDepth(d)
	case *UnaryExpr:
		return nt.StringDepth(d)
	case *BinaryExpr:
		return nt.StringDepth(d)
	case *ExprStmt:
		return nt.String()
	case *AssignStmt:
		return nt.StringDepth(d)
	case *DeclStmt:
		return nt.StringDepth(d)
	case *ValueSpec:
		return nt.StringDepth(d)
	case *GenDecl:
		return nt.StringDepth(d)
	case *File:
		return nt.String()
	case nil:
		return "<nil>"
	default:
		return "<<UNKNOWN>>"
	}
	return ""
}

func (id *Ident) String() string {
	if id != nil {
		return id.Tok.Val
	}
	return "<nil>"
}

func (n *BasicLit) String() string {
	return fmt.Sprintf("(basiclit %s)", n.Tok)
}

func (n *ParenExpr) String() string {
	return n.StringDepth(0)
}

func (n *UnaryExpr) String() string {
	return n.StringDepth(0)
	// return fmt.Sprintf("(UnaryExpr \n\tOp:%s \n\tX:%s)", n.Op, sprintd(n.X, 0))
}

func (n *BinaryExpr) String() string {
	return n.StringDepth(0)
	// return fmt.Sprintf("(BinaryExpr \n\tOp:%s \n\tX:%s \n\tY:%s)", n.Op, sprintd(n.X, 0), sprintd(n.Y, 0))
}

func (n *ExprStmt) String() string {
	return Sprint(n.X)
}

func (n *AssignStmt) String() string {
	return n.StringDepth(0)
}

func (n *DeclStmt) String() string {
	return n.StringDepth(0)
}

func (n *GenDecl) String() string {
	return n.StringDepth(0)
}

func (n *ValueSpec) String() string {
	return n.StringDepth(0)
}

func (n *File) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("(File ")
	for _, node := range n.List {
		buffer.WriteString("\n\n")
		buffer.WriteString(sprintd(node, 0))
	}
	buffer.WriteString(")")
	return buffer.String()
}

// Util print functions to print the correct depth

func (n *ParenExpr) StringDepth(d int) string {
	return fmt.Sprintf("(ParenExpr '%s' %s '%s')", n.Lparen.Val, sprintd(n.X, d+1), n.Rparen.Val)
}

func (n *UnaryExpr) StringDepth(d int) string {
	var buffer bytes.Buffer
	buffer.WriteString("(UnaryExpr ")
	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}

	buffer.WriteString("Op: ")
	buffer.WriteString(n.Op.String())
	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString(sprintd(n.X, d+1))
	buffer.WriteString(")")

	return buffer.String()
}

func (n *BinaryExpr) StringDepth(d int) string {
	var buffer bytes.Buffer
	buffer.WriteString("(BinaryExpr ")
	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}

	buffer.WriteString("Op: ")
	buffer.WriteString(n.Op.String())
	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("X: ")
	buffer.WriteString(sprintd(n.X, d+1))

	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("Y: ")
	buffer.WriteString(sprintd(n.Y, d+1))
	buffer.WriteString(")")

	return buffer.String()
}

func (n *AssignStmt) StringDepth(d int) string {
	var buffer bytes.Buffer
	buffer.WriteString("(AssignStmt ")
	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}

	buffer.WriteString("Lhs: ")
	buffer.WriteString(sprintd(n.Lhs, d+1))
	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("Tok: ")
	buffer.WriteString(n.Tok.String())

	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("Rhs: ")
	buffer.WriteString(sprintd(n.Rhs, d+1))
	buffer.WriteString(")")

	return buffer.String()
}

func (n *DeclStmt) StringDepth(d int) string {
	return fmt.Sprintf("(DeclStmt %s)", sprintd(n.Decl, d+1))
}

func (n *ValueSpec) StringDepth(d int) string {
	var buffer bytes.Buffer
	buffer.WriteString("(ValueSpec ")

	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("Name: ")
	buffer.WriteString(sprintd(n.Name, d+1))

	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString("Value: ")
	buffer.WriteString(sprintd(n.Value, d+1))
	buffer.WriteString(")")

	return buffer.String()
}

func (n *GenDecl) StringDepth(d int) string {
	var buffer bytes.Buffer
	buffer.WriteString("(GenDecl ")
	buffer.WriteString(n.Tok.String())

	buffer.WriteString("\n")
	for i := 0; i < d; i++ {
		buffer.WriteString("\t")
	}
	buffer.WriteString(sprintd(n.Spec, d+1))
	buffer.WriteString(")")

	return buffer.String()
}
