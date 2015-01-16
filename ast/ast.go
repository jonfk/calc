package ast

import (
	"bytes"
	"fmt"
	"jon/calc/lex"
	"strings"
	"unicode"
	"unicode/utf8"
)

// -------------------------------------------------------------------
// Interfaces

// All node types implement the Node interface.
type Node interface {
	Pos() lex.Pos
	End() lex.Pos
}

// All expression nodes implement the Expr interface.
type Expr interface {
	Node
	exprNode()
	Type() NodeType
}

// All declaration nodes implement the Decl interface.
type Decl interface {
	Node
	declNode()
}

/*
// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}

*/

// ---------------------------------------------------------------------
// Node Types

type NodeType int

const (
	CommentNode NodeType = iota
	CommentGroupNode
	BadExprNode
	IdentNode
	BasicLitNode
	ParenExprNode
	UnaryExprNode
	BinaryExprNode
	BlockExprNode
	IfExprNode
	AssignNode
)

// ---------------------------------------------------------------------
// Comments

// A Comment node represents a single //-style or /*-style comment.
type Comment struct {
	Slash lex.Pos // position of "/" starting the comment
	Text  string  // comment text (excluding '\n' for //-style comments)
}

func (c *Comment) Pos() lex.Pos   { return c.Slash }
func (c *Comment) End() lex.Pos   { return lex.Pos(int(c.Slash) + len(c.Text)) }
func (c *Comment) Type() NodeType { return CommentNode }

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
//
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() lex.Pos   { return g.List[0].Pos() }
func (g *CommentGroup) End() lex.Pos   { return g.List[len(g.List)-1].End() }
func (g *CommentGroup) Type() NodeType { return CommentGroupNode }

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func stripTrailingWhitespace(s string) string {
	i := len(s)
	for i > 0 && isWhitespace(s[i-1]) {
		i--
	}
	return s[0:i]
}

// Text returns the text of the comment.
// Comment markers (//, /*, and */), the first space of a line comment, and
// leading and trailing empty lines are removed. Multiple empty lines are
// reduced to one, and trailing space on lines is trimmed. Unless the result
// is empty, it is newline-terminated.
//
func (g *CommentGroup) Text() string {
	if g == nil {
		return ""
	}
	comments := make([]string, len(g.List))
	for i, c := range g.List {
		comments[i] = string(c.Text)
	}

	lines := make([]string, 0, 10) // most comments are less than 10 lines
	for _, c := range comments {
		// Remove comment markers.
		// The parser has given us exactly the comment text.
		switch c[1] {
		case '/':
			//-style comment (no newline at the end)
			c = c[2:]
			// strip first space - required for Example tests
			if len(c) > 0 && c[0] == ' ' {
				c = c[1:]
			}
		case '*':
			/*-style comment */
			c = c[2 : len(c)-2]
		}

		// Split on newlines.
		cl := strings.Split(c, "\n")

		// Walk lines, stripping trailing white space and adding to list.
		for _, l := range cl {
			lines = append(lines, stripTrailingWhitespace(l))
		}
	}

	// Remove leading blank lines; convert runs of
	// interior blank lines to a single blank line.
	n := 0
	for _, line := range lines {
		if line != "" || n > 0 && lines[n-1] != "" {
			lines[n] = line
			n++
		}
	}
	lines = lines[0:n]

	// Add final "" entry to get trailing newline from Join.
	if n > 0 && lines[n-1] != "" {
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

// ---------------------------------------------------------------
// Expressions

// An expression is represented by a tree consisting of one
// or more of the following concrete expression nodes.
//
type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	//
	BadExpr struct {
		From, To lex.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		Tok lex.Token // identifier token
		Obj *Object   // denoted object; or nil
	}

	// A BasicLit node represents a literal of basic type.
	BasicLit struct {
		Tok lex.Token // token.INT, token.FLOAT, not supported yet: token.CHAR, or token.STRING
	}

	// A ParenExpr node represents a parenthesized expression.
	ParenExpr struct {
		Lparen lex.Token // position of "("
		X      Expr      // parenthesized expression
		Rparen lex.Token // position of ")"
	}

	// A UnaryExpr node represents a unary expression.
	UnaryExpr struct {
		Op lex.Token // operator
		X  Expr      // operand
	}

	// A BinaryExpr node represents a binary expression.
	BinaryExpr struct {
		X  Expr      // left operand
		Op lex.Token // operator
		Y  Expr      // right operand
	}

	// A BlockExpr node represents a list of expressions.
	// Evaluates to the last expression in the list.
	//
	BlockExpr struct {
		StartPos lex.Pos
		List     []Expr
		EndPos   lex.Pos
	}

	// An IfExpr node represents an if expression.
	IfExpr struct {
		If     lex.Token
		Cond   Expr
		Body   *BlockExpr
		Else   *BlockExpr
		EndTok lex.Token
	}
)

// Pos and End implementations for expression/type nodes.
//
func (x *BadExpr) Pos() lex.Pos    { return x.From }
func (x *Ident) Pos() lex.Pos      { return x.Tok.Pos }
func (x *BasicLit) Pos() lex.Pos   { return x.Tok.Pos }
func (x *ParenExpr) Pos() lex.Pos  { return x.Lparen.Pos }
func (x *UnaryExpr) Pos() lex.Pos  { return x.Op.Pos }
func (x *BinaryExpr) Pos() lex.Pos { return x.X.Pos() }
func (x *BlockExpr) Pos() lex.Pos  { return x.StartPos }
func (x *IfExpr) Pos() lex.Pos     { return x.If.Pos }

func (x *BadExpr) End() lex.Pos    { return x.To }
func (x *Ident) End() lex.Pos      { return lex.Pos(int(x.Tok.Pos) + len(x.Tok.Val)) }
func (x *BasicLit) End() lex.Pos   { return lex.Pos(int(x.Tok.Pos) + len(x.Tok.Val)) }
func (x *ParenExpr) End() lex.Pos  { return x.Rparen.Pos + 1 }
func (x *UnaryExpr) End() lex.Pos  { return x.X.End() }
func (x *BinaryExpr) End() lex.Pos { return x.Y.End() }
func (x *BlockExpr) End() lex.Pos  { return x.EndPos }
func (x *IfExpr) End() lex.Pos     { return x.EndTok.Pos }

func (x *BadExpr) Type() NodeType    { return BadExprNode }
func (x *Ident) Type() NodeType      { return IdentNode }
func (x *BasicLit) Type() NodeType   { return BasicLitNode }
func (x *ParenExpr) Type() NodeType  { return ParenExprNode }
func (x *UnaryExpr) Type() NodeType  { return UnaryExprNode }
func (x *BinaryExpr) Type() NodeType { return BinaryExprNode }
func (x *BlockExpr) Type() NodeType  { return BlockExprNode }
func (x *IfExpr) Type() NodeType     { return IfExprNode }

// exprNode() ensures that only expression/type nodes can be
// assigned to an ExprNode.
//
func (*BadExpr) exprNode()    {}
func (*Ident) exprNode()      {}
func (*BasicLit) exprNode()   {}
func (*ParenExpr) exprNode()  {}
func (*UnaryExpr) exprNode()  {}
func (*BinaryExpr) exprNode() {}
func (*BlockExpr) exprNode()  {}

// ----------------------------------------------------------------------------
// Convenience functions for Idents

// NewIdent creates a new Ident without position.
// Useful for ASTs generated by code other than the Go parser.
//
func NewIdent(name string) *Ident { return &Ident{lex.Token{lex.IDENTIFIER, lex.NoPos, name}, nil} }

// IsExported reports whether name is an exported Go symbol
// (that is, whether it begins with an upper-case letter).
//
func IsExported(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}

// IsExported reports whether id is an exported Go symbol
// (that is, whether it begins with an uppercase letter).
//
func (id *Ident) IsExported() bool { return IsExported(id.Tok.Val) }

// ----------------------------------------------------------------------------
// Declarations

type (
	// An Assign node represents an assignment expression.
	Assign struct {
		Let    lex.Token
		Lhs    *Ident
		Assign lex.Token
		Rhs    Expr
		EndTok lex.Token
	}
)

func (x *Assign) Pos() lex.Pos { return x.Let.Pos }

func (x *Assign) End() lex.Pos { return x.EndTok.Pos }

func (x *Assign) Type() NodeType { return AssignNode }

func (*Assign) declNode() {}

// ----------------------------------------------------------------------------
// Files and packages

// A File node represents a Go source file.
//
// The Comments list contains all comments in the source file in order of
// appearance, including the comments that are pointed to from other nodes
// via Doc and Comment fields.
//
type File struct {
	Doc      *CommentGroup // associated documentation; or nil
	StartPos lex.Pos
	EndPos   lex.Pos
	//Package    lex.Pos       // position of "package" keyword
	//Name       *Ident          // package name
	Scope *Scope // package scope (this file only)
	List  []Node // list of nodes in file
	// Block *BlockExpr // Expressions in this file
	//Imports    []*ImportSpec   // imports in this file
	Unresolved []*Ident        // unresolved identifiers in this file
	Comments   []*CommentGroup // list of all comments in the source file
}

func (f *File) Pos() lex.Pos { return f.StartPos }
func (f *File) End() lex.Pos { return f.EndPos }

func NewFile() *File {
	return &File{
		Doc:   new(CommentGroup),
		Scope: new(Scope),
		// Block: new(BlockExpr),
	}
}

/*
// A Package node represents a set of source files
// collectively building a Go package.
//
type Package struct {
	Name    string             // package name
	Scope   *Scope             // package scope across all files
	Imports map[string]*Object // map of package id -> package object
	Files   map[string]*File   // Go source files by filename
}

func (p *Package) Pos() token.Pos { return token.NoPos }
func (p *Package) End() token.Pos { return token.NoPos }
*/

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
