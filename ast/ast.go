package ast

import (
	"github.com/jonfk/calc/lex"
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
}

// All declaration nodes implement the Decl interface.
type Decl interface {
	Node
	declNode()
}

// All statement nodes implement the Stmt interface.
type Stmt interface {
	Node
	stmtNode()
}

// ---------------------------------------------------------------------
// Comments

// A Comment node represents a single //-style or /*-style comment.
type Comment struct {
	Slash lex.Pos // position of "/" starting the comment
	Text  string  // comment text (excluding '\n' for //-style comments)
}

func (c *Comment) Pos() lex.Pos { return c.Slash }
func (c *Comment) End() lex.Pos { return lex.Pos(int(c.Slash) + len(c.Text)) }

// A CommentGroup represents a sequence of comments
// with no other tokens and no empty lines between.
//
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (g *CommentGroup) Pos() lex.Pos { return g.List[0].Pos() }
func (g *CommentGroup) End() lex.Pos { return g.List[len(g.List)-1].End() }

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
// Statements

// A statement is represented by a tree consisting of one
// or more of the following concrete statement nodes.
//
type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	//
	BadStmt struct {
		From, To token.Pos // position range of bad statement
	}

	// A DeclStmt node represents a declaration in a statement list.
	DeclStmt struct {
		Decl Decl // *GenDecl with CONST, TYPE, or VAR token
	}

	// An ExprStmt node represents a (stand-alone) expression
	// in a statement list.
	//
	ExprStmt struct {
		X Expr // expression
	}

	// An AssignStmt node represents an assignment or
	// a short variable declaration.
	//
	AssignStmt struct {
		Lhs    []Expr
		TokPos token.Pos   // position of Tok
		Tok    token.Token // assignment token, DEFINE
		Rhs    []Expr
	}
)

// ----------------------------------------------------------------------------
// Declarations

// A declaration is represented by one of the following declaration nodes.
//
type (
	// var and val declarations
	GenDecl struct {
		Doc    *CommentGroup // associated documentation; or nil
		TokPos token.Pos     // position of Tok
		Tok    token.Token   // IMPORT, CONST, TYPE, VAR
		Lparen token.Pos     // position of '(', if any
		Specs  []Spec
		Rparen token.Pos // position of ')', if any
	}

	// A FuncDecl node represents a function declaration.
	FuncDecl struct {
		Doc  *CommentGroup // associated documentation; or nil
		Recv *FieldList    // receiver (methods); or nil (functions)
		Name *Ident        // function/method name
		Type *FuncType     // function signature: parameters, results, and position of "func" keyword
		Body *BlockStmt    // function body; or nil (forward declaration)
	}
)

func (x *Assign) Pos() lex.Pos { return x.Let.Pos }

func (x *Assign) End() lex.Pos { return x.EndTok.Pos }

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
