package parse

import (
	"fmt"
	"github.com/jonfk/calc/ast"
	"github.com/jonfk/calc/lex"
	"os"
	"strings"
)

// parser holds the state of the scanner.
type Parser struct {
	name      string      // the name of the input; used only for error reports
	input     string      // the string being scanned
	pos       int         // the position of token in Items; pos == -1 when Items is nil
	Items     []lex.Token // the unreduced items received from the lexer
	lastToken lex.Token   // Used for error and debug messages

	Lexer *lex.Lexer // the lexer

	File     *ast.File  // the file being parsed
	topScope *ast.Scope // may be nil if topmost scope
	lastNode ast.Node   // last node parsed ??? currently only used by let. Is it necessary?

	pDepth *ParenDepth // paren depth for parsing expressions
}

// -----------------------------------------------------------------------------
// scoping support

func (p *Parser) openScope() {
	p.topScope = ast.NewScope(p.topScope)
}

func (p *Parser) closeScope() {
	p.topScope = p.topScope.Outer
}

// ------------------------------------------------------------------------------
// parsing support

// next return the next token in the input.
// calls getItem if at end of Items
func (p *Parser) next() lex.Token {
	if p.pos != -1 && p.Items[p.pos].Typ == lex.EOF {
		return p.Items[p.pos]
	}
	p.pos += 1
	if p.pos >= len(p.Items) {
		p.errorf("Internal error in next(): parser.pos moving out of bounds of lexed tokens\n")
	}
	// Ignore comments for now
	for p.Items[p.pos].Typ == lex.LINECOMMENT || p.Items[p.pos].Typ == lex.BLOCKCOMMENT {
		p.pos += 1
	}
	// call p.errorf if lexing error
	if p.Items[p.pos].Typ == lex.ERROR {
		p.errorf(p.Items[p.pos].String())
	}
	p.lastToken = p.Items[p.pos]
	return p.Items[p.pos]
}

func (p *Parser) nextNonNewline() lex.Token {
	t := p.next()
	for t.Typ == lex.NEWLINE {
		t = p.next()
	}
	return t
}

// peek returns the k forward token in items but does not move the pos.
func (p *Parser) peek(k int) lex.Token {
	for (p.pos + k) >= len(p.Items) {
		p.errorf("Internal error in peek(): parser.pos moving out of bounds of lexed tokens\n")
	}
	return p.Items[p.pos+k]
}

// backup steps back one token.
// Can only be called as many times as there are unreduced tokens in Items
// return error if there aren't enough tokens in Items
func (p *Parser) backup() error {
	if p.pos <= -1 {
		p.errorf("Internal error in backup: Cannot backup anymore pos is at start of Items\n")
	}
	p.pos -= 1
	// if p.pos != -1 {
	// 	p.lastToken = p.Items[p.pos]
	// }
	return nil
}

// accept consumes the next token if it's from the valid set.
func (p *Parser) accept(valid []lex.TokenType) bool {
	item := p.next()
	for _, tokTyp := range valid {
		if item.Typ == tokTyp {
			return true
		}
	}
	p.backup()
	return false
}

// acceptRun consumes a run of tokens from the valid set.
func (p *Parser) acceptRun(valid []lex.TokenType) {
	for p.accept(valid) {
	}
}

func (p *Parser) expect(valid lex.TokenType) bool {
	if p.Items[p.pos].Typ == valid {
		return true
	}
	return false
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by next.
func (p *Parser) lineNumber() int {
	item := p.Items[p.pos]
	return 1 + strings.Count(p.input[:item.Pos], "\n")
}

// lineNumber reports which line we're on, based a lex.Pos
func (p *Parser) lineNumberAt(pos lex.Pos) int {
	return 1 + strings.Count(p.input[:pos], "\n")
}

// errorf prints an error and terminates the scan
func (p *Parser) errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	os.Exit(1)
}

// Parse creates a new parser for the input string.
// It uses lex to tokenize the input
func Parse(name, input string) *Parser {
	l := lex.Lex(name, input)
	p := &Parser{
		name:     name,
		input:    input,
		pos:      -1,
		Lexer:    l,
		File:     ast.NewFile(),
		topScope: ast.NewScope(nil),
		pDepth:   new(ParenDepth),
	}
	p.run()
	return p
}

// runs the parser
func (p *Parser) run() {
	// for p.state = parseProg; p.state != nil; {
	// 	p.state = p.state(p)
	// }

	// lex everything
	t := p.Lexer.NextItem()
	for ; t.Typ != lex.EOF; t = p.Lexer.NextItem() {
		p.Items = append(p.Items, t)
	}
	p.Items = append(p.Items, t)

	parseFile(p)
	return
}

// --------------------------------------------------------------------------------------------
// Recursive descent parser
// Mutually recursive functions

func parseFile(p *Parser) {
	switch t := p.nextNonNewline(); {
	case t.Typ == lex.IDENTIFIER || isLiteral(t) || t.Typ == lex.LEFTPAREN || isUnaryOp(t):
		p.backup()
		expr := parseStartExpr(p)
		p.File.List = append(p.File.List, expr)
		parseFile(p)
	case t.Typ == lex.EOF:
		return
	// case t.Typ == lex.LET:
	// 	assign := &ast.Assign{
	// 		Let: t,
	// 	}
	// 	p.File.List = append(p.File.List, assign)
	// 	p.lastNode = assign
	// 	parseLet(p)
	default:
		p.errorf("Invalid statement at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
}

// parses a let expression according to the grammar rule:
// let_expr ::= LET IDENTIFIER ASSIGN expr END

func parseStartExpr(p *Parser) ast.Expr {
	switch t := p.next(); {
	case t.Typ == lex.IDENTIFIER:
		ident := newIdentExpr(p, t)
		return parseLiteralOrIdent(p, ident, ident)
	case isLiteral(t):
		bLit := &ast.BasicLit{Tok: t}
		return parseLiteralOrIdent(p, bLit, bLit)
	case isUnaryOp(t):
		unary := &ast.UnaryExpr{Op: t}
		return parseUnaryExpr(p, unary, unary)
	case t.Typ == lex.LEFTPAREN:
		paren := newParenExpr(p, t)
		return parseParenExpr(p, paren, paren)
	default:
		p.errorf("Invalid start of expression at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
	return nil
}

func parseLiteralOrIdent(p *Parser, tree, last ast.Expr) ast.Expr {
	switch t := p.next(); {
	case t.IsOperator():
		oper := &ast.BinaryExpr{Op: t}
		tree, _ = ast.InsertExpr(tree, oper)
		return parseBinaryExpr(p, tree, oper)
	case t.Typ == lex.NEWLINE && len(p.pDepth.Stack) > 0:
		return parseLiteralOrIdent(p, tree, last)
	case atTerminator(t):
		return tree
	case t.Typ == lex.RIGHTPAREN:
		paren := p.pDepth.pop()
		paren.Rparen = t
		return parseParenExpr(p, tree, paren)
	default:
		p.errorf("Invalid expression at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
	return nil
}

func parseParenExpr(p *Parser, tree ast.Expr, last *ast.ParenExpr) ast.Expr {
	// if expr.X == nil && expr.Rparen.Val == ""
	// else if expr.X != nil && expr.Rparen.Val == ""
	// else if expr.Rparen.Val == ")"

	// Unclosed Empty paren expr
	if last.X == nil && last.Rparen.Val == "" {
		switch t := p.nextNonNewline(); {
		case t.Typ == lex.IDENTIFIER:
			ident := newIdentExpr(p, t)
			tree, _ = ast.InsertExpr(tree, ident)
			return parseLiteralOrIdent(p, tree, ident)
		case isLiteral(t):
			num := &ast.BasicLit{Tok: t}
			tree, _ = ast.InsertExpr(tree, num)
			return parseLiteralOrIdent(p, tree, num)
		case isUnaryOp(t):
			unary := &ast.UnaryExpr{Op: t}
			tree, _ = ast.InsertExpr(tree, unary)
			return parseUnaryExpr(p, tree, unary)
		case t.Typ == lex.LEFTPAREN:
			paren := newParenExpr(p, t)
			tree, _ = ast.InsertExpr(tree, paren)
			return parseParenExpr(p, tree, paren)
		case t.Typ == lex.RIGHTPAREN:
			paren := p.pDepth.pop()
			paren.Rparen = t
			if paren != last {
				p.errorf("Internal error in parseLiteralOrIdent closing paren not matching current paren expr at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			}
			return parseParenExpr(p, tree, last)
		default:
			p.errorf("Invalid expression at line %d:%d with token '%s' in file %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	} else if last.Rparen.Val == ")" {
		if last.X == nil {
			// empty closed paren expr ()
			// give value nil to ()
			fmt.Printf("Warning: empty paren expression has value nil")
		}
		// Closed non-empty paren expr
		switch t := p.next(); {
		case t.Typ == lex.IDENTIFIER:
			p.errorf("Invalid paren expression closed expression followed by identifier at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		case t.Typ == lex.INT || t.Typ == lex.FLOAT:
			p.errorf("Invalid paren expression closed expression followed by literal at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		case t.IsOperator():
			binary := &ast.BinaryExpr{Op: t}
			tree, _ = ast.InsertExpr(tree, binary)
			return parseBinaryExpr(p, tree, binary)
		case t.Typ == lex.LEFTPAREN:
			p.errorf("Invalid paren expression closed expression followed by opening parenthesis at line %d:%d with token '%s' in file %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		case t.Typ == lex.RIGHTPAREN:
			// close enclosing paren in case parenExpr{X:parenExpr{}}
			paren := p.pDepth.pop()
			paren.Rparen = t
			return parseParenExpr(p, tree, paren)
		case t.Typ == lex.NEWLINE && len(p.pDepth.Stack) > 0:
			return parseParenExpr(p, tree, last)
		case atTerminator(t):
			return tree
		default:
			p.errorf("Invalid expression at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	} else {
		p.errorf("Internal error in parseParenExpr: at line %d:%d with token '%s' in file: %s\n", p.lineNumber(), p.lastToken.Pos, p.lastToken.Val, p.name)
		return nil
	}
	return nil
}

func parseUnaryExpr(p *Parser, tree ast.Expr, last *ast.UnaryExpr) ast.Expr {
	if last.X == nil {
		switch t := p.next(); {
		case t.Typ == lex.IDENTIFIER:
			ident := newIdentExpr(p, t)
			tree, _ = ast.InsertExpr(tree, ident)
			return parseLiteralOrIdent(p, tree, ident)
		case isLiteral(t):
			bLit := &ast.BasicLit{Tok: t}
			tree, _ = ast.InsertExpr(tree, bLit)
			return parseLiteralOrIdent(p, tree, bLit)
		case t.Typ == lex.ADD || t.Typ == lex.SUB:
			unary := &ast.UnaryExpr{Op: t}
			tree, _ = ast.InsertExpr(tree, unary)
			return parseUnaryExpr(p, tree, unary)
		case t.Typ == lex.LEFTPAREN:
			paren := newParenExpr(p, t)
			tree, _ = ast.InsertExpr(tree, paren)
			return parseParenExpr(p, tree, paren)
		default:
			p.errorf("Invalid unary expression at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t, p.name)
		}
	} else {
		p.errorf("Internal error in parseUnaryExpr at line %d in file : %s\n", p.lineNumber(), p.name)
	}
	return nil
}

func parseBinaryExpr(p *Parser, tree ast.Expr, last *ast.BinaryExpr) ast.Expr {
	if last.Y == nil {
		switch t := p.nextNonNewline(); {
		case isUnaryOp(t):
			unary := &ast.UnaryExpr{Op: t}
			tree, _ = ast.InsertExpr(tree, unary)
			return parseUnaryExpr(p, tree, unary)
		case t.Typ == lex.IDENTIFIER:
			ident := newIdentExpr(p, t)
			tree, _ = ast.InsertExpr(tree, ident)
			return parseLiteralOrIdent(p, tree, ident)
		case isLiteral(t):
			bLit := &ast.BasicLit{Tok: t}
			tree, _ = ast.InsertExpr(tree, bLit)
			return parseLiteralOrIdent(p, tree, bLit)
		case t.Typ == lex.LEFTPAREN:
			paren := newParenExpr(p, t)
			tree, _ = ast.InsertExpr(tree, paren)
			return parseParenExpr(p, tree, paren)
		default:
			p.errorf("Invalid expression at line %d:%d with token '%s' in file : %S\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	} else {
		p.errorf("Internal Error: Invalid parser state in parseBinaryExpr")
		return nil
	}
	return nil
}

// --------------------------------------------------------------------------------------------
// Utility functions for parsing

func isUnaryOp(t lex.Token) bool {
	switch t.Typ {
	case lex.NOT, lex.ADD, lex.SUB:
		return true
	default:
		return false
	}
}

func isLiteral(t lex.Token) bool {
	switch t.Typ {
	case lex.BOOL, lex.INT, lex.FLOAT, lex.STRING:
		return true
	default:
		return false
	}
}

func atTerminator(t lex.Token) bool {
	if t.Typ == lex.NEWLINE || t.Typ == lex.SEMICOLON || t.Typ == lex.EOF || t.Typ == lex.THEN {
		return true
	}
	return false
}

func newParenExpr(p *Parser, t lex.Token) *ast.ParenExpr {
	paren := &ast.ParenExpr{Lparen: t}
	p.pDepth.push(paren)
	return paren
}

func newIdentExpr(p *Parser, t lex.Token) *ast.Ident {
	switch t.Typ {
	case lex.IDENTIFIER:
		ident := &ast.Ident{Tok: t}
		obj := p.topScope.Lookup(t.Val)
		if obj == nil {
			p.File.Unresolved = append(p.File.Unresolved, ident)
		} else {
			ident.Obj = obj
		}
		return ident
	default:
		p.errorf("Invalid expression at %d:%d expected an identifier but found '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
	return nil
}
