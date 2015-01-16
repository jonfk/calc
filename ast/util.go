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
	switch a.(type) {
	case *Ident:
		switch b.(type) {
		case *Ident:
			if b.(*Ident).Tok.Equals(a.(*Ident).Tok) {
				return true
			}
			return false
		default:
			return false
		}
	case *BasicLit:
		switch b.(type) {
		case *BasicLit:
			if b.(*BasicLit).Tok.Equals(a.(*BasicLit).Tok) {
				return true
			}
			return false
		default:
			return false
		}
	case *ParenExpr:
		switch b.(type) {
		case *ParenExpr:
			aP, bP := a.(*ParenExpr), b.(*ParenExpr)
			return aP.Lparen.Equals(bP.Lparen) &&
				aP.Rparen.Equals(bP.Rparen) &&
				Equals(aP.X, bP.X)
		default:
			return false
		}
	case *UnaryExpr:
		switch b.(type) {
		case *UnaryExpr:
			return a.(*UnaryExpr).Op.Equals(b.(*UnaryExpr).Op) &&
				Equals(a.(*UnaryExpr).X, b.(*UnaryExpr).X)
		default:
			return false
		}
	case *BinaryExpr:
		switch b.(type) {
		case *BinaryExpr:
			return a.(*BinaryExpr).Op.Equals(b.(*BinaryExpr).Op) &&
				Equals(a.(*BinaryExpr).Y, b.(*BinaryExpr).Y) &&
				Equals(a.(*BinaryExpr).X, b.(*BinaryExpr).X)
		default:
			return false
		}
	case *File:
		switch b.(type) {
		case *File:
			afile, bfile := a.(*File), b.(*File)
			if len(afile.List) == len(bfile.List) {
				for i := range afile.List {
					if !Equals(afile.List[i], bfile.List[i]) {
						return false
					}
				}
				return true
			}
			return false
		default:
			return false
		}
	default:
		return false
	}
}

func Sprint(n Node) string {
	return sprintd(n, 0)
}

// sprintd print prints the node to standard output with the proper depth
// helper function to improve formatting
func sprintd(n Node, d int) string {
	switch n.(type) {
	case *Ident:
		return n.(*Ident).String()
	case *BasicLit:
		return n.(*BasicLit).String()
	case *ParenExpr:
		return n.(*ParenExpr).StringDepth(d)
	case *UnaryExpr:
		return n.(*UnaryExpr).StringDepth(d)
	case *BinaryExpr:
		return n.(*BinaryExpr).StringDepth(d)
	case *File:
		return n.(*File).String()
	default:
		return ""
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

func (n *File) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("(File ")
	for _, node := range n.List {
		buffer.WriteString("\n\n")
		buffer.WriteString(sprintd(node, 0))
	}
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
