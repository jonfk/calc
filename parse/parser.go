package parse

import (
	"fmt"
	"os"
	"jon/calc/ast"
	"jon/calc/lex"
	"strings"
	// "github.com/davecgh/go-spew/spew"
)


// parser holds the state of the scanner.
type Parser struct {
	name   string      // the name of the input; used only for error reports
	input  string      // the string being scanned
	pos    int         // the position of token in Items; pos == -1 when Items is nil
	Items  []lex.Token // the unreduced items received from the lexer
	lastToken lex.Token // Used for error and debug messages

	Lexer  *lex.Lexer  // the lexer

	File   *ast.File    // the file being parsed
	topScope *ast.Scope // may be nil if topmost scope
	lastNode ast.Node   // last node parsed ??? currently only used by let. Is it necessary?

	pDepth  *ParenDepth  // paren depth for parsing expressions
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

// ignore skips over the pending input before this point.
func (p *Parser) ignore() {
	p.Items = append(p.Items[:p.pos], p.Items[p.pos+1:]...)
	p.pos -= 1
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
		name:  name,
		input: input,
		pos: -1,
		Lexer: l,
		File: ast.NewFile(),
		topScope: ast.NewScope(nil),
		pDepth: new(ParenDepth),
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
	for ;t.Typ != lex.EOF; t = p.Lexer.NextItem() {
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
	case t.Typ == lex.IDENTIFIER || t.Typ == lex.INT || t.Typ == lex.FLOAT || t.Typ == lex.LEFTPAREN || t.Typ == lex.ADD || t.Typ == lex.SUB:
		p.backup()
		expr := parseExpr(p, nil)
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
		p.errorf("Invalid statement at line %v:%v with token '%s' in file : %s\n",p.lineNumber(), t.Pos, t.Val, p.name)
	}
}

func parseIdent(p *Parser) *ast.Ident {
	switch t := p.nextNonNewline(); {
	case t.Typ == lex.IDENTIFIER:
		ident := &ast.Ident{Tok:t}
		obj := p.topScope.Lookup(t.Val)
		if obj == nil {
			p.File.Unresolved = append(p.File.Unresolved, ident)
		} else {
			ident.Obj = obj
		}
		return ident
	default:
		p.errorf("Invalid expression at %d:%d expected an identifier but found %v in file : %S\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
	return nil
}

// parses a let expression according to the grammar rule:
// let_expr ::= LET IDENTIFIER ASSIGN expr END
func parseLet(p *Parser) {
	assign := p.lastNode.(*ast.Assign)
	switch t := p.next(); {
	case t.Typ == lex.IDENTIFIER:
		id := ast.NewIdent(t.Val)
		id.Obj = ast.NewObj(ast.Val, t.Val)
		id.Obj.Decl = assign
		assign.Lhs = id
		if eq := p.next(); eq.Typ == lex.ASSIGN {
			assign.Assign = eq
		}
		assign.Rhs = parseExpr(p, nil)
	default:
		p.errorf("Invalid let expression at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
}

func parseExpr(p *Parser, expr ast.Expr) ast.Expr {
	switch expr.(type) {
	case nil:
		return parseStartExpr(p)
	case *ast.Ident, *ast.BasicLit:
		mark := p.pos
		switch t := p.nextNonNewline(); {
		case t.IsOperator():
			oper := &ast.BinaryExpr{X:expr, Op:t}
			return parseExpr(p, oper)
		default:
			p.pos = mark
			t = p.next()
			if atTerminator(t) {
				return expr
			} else if t.Typ == lex.RIGHTPAREN {
				paren := p.pDepth.pop()
				paren.Rparen = t
				return expr
			} else {
				p.errorf("Invalid expression at line %d:%d with token %v in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			}
		}
	case *ast.ParenExpr:
		return parseParenExpr(p, expr.(*ast.ParenExpr))
	case *ast.UnaryExpr:
		return parseUnaryExpr(p, expr.(*ast.UnaryExpr))
	case *ast.BinaryExpr:
		return parseBinaryExpr(p, expr.(*ast.BinaryExpr))
	default:
		p.errorf("Internal error unknown expression at line %d:%d after token '%s' in file : %s\n", p.lineNumber(), p.lastToken.Pos, p.lastToken.Val, p.name)
		return nil
	}
	return nil
}

func  parseStartExpr(p *Parser) ast.Expr {
	switch t := p.nextNonNewline(); t.Typ {
	case lex.IDENTIFIER:
		p.backup()
		ident := parseIdent(p)
		return parseExpr(p, ident)
	case lex.INT, lex.FLOAT:
		bLit := &ast.BasicLit{Tok:t}
		return parseExpr(p, bLit)
	case lex.ADD, lex.SUB:
		unary := &ast.UnaryExpr{Op:t}
		return parseExpr(p, unary)
	case lex.LEFTPAREN:
		paren := newParenExpr(p, t)
		return parseExpr(p, paren)
	default:
		p.errorf("Invalid start of expression at line %d:%d with token %v in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
	}
	return nil
}

func atTerminator(t lex.Token) bool {
	if t.Typ == lex.NEWLINE || t.Typ == lex.SEMICOLON || t.Typ == lex.EOF || t.Typ == lex.THEN {
		return true
	}
	return false
}

func newParenExpr(p *Parser, t lex.Token) *ast.ParenExpr {
	paren := &ast.ParenExpr{Lparen:t}
	p.pDepth.push(paren)
	return paren
}

func parseParenExpr(p *Parser, expr *ast.ParenExpr) ast.Expr {
	// if expr.X == nil && expr.Rparen.Val == ""
	// else if expr.X != nil && expr.Rparen.Val == ""
	// else if expr.Rparen.Val == ")"

	// Unclosed Empty paren expr
	if expr.X == nil && expr.Rparen.Val == "" {
		switch t := p.nextNonNewline(); {
		case t.Typ == lex.IDENTIFIER:
			p.backup()
			ident := parseIdent(p)
			expr.X = parseExpr(p, ident)
			return parseExpr(p, expr)
		case t.Typ == lex.INT || t.Typ == lex.FLOAT:
			num := &ast.BasicLit{Tok:t}
			expr.X = parseExpr(p, num)
			return parseExpr(p, expr)
		case t.Typ == lex.SUB || t.Typ == lex.ADD:
			unary := &ast.UnaryExpr{Op:t}
			expr.X = parseExpr(p, unary)
			return parseExpr(p, expr)
		case t.Typ == lex.LEFTPAREN:
			paren := newParenExpr(p, t)
			expr.X = parseExpr(p, paren)
			return parseExpr(p, expr)
		case t.Typ == lex.RIGHTPAREN:
			expr.Rparen = t
			paren := p.pDepth.pop()
			if paren != expr {
				p.errorf("Internal error in parseExpr closing paren not matching current paren expr at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			}
			return parseExpr(p, expr)
		default:
			p.errorf("Invalid expression at line %d:%d with token '%s' in file %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	} else if expr.X != nil && expr.Rparen.Val == "" {
		// Unclosed non-empty paren expr
		switch t := p.nextNonNewline(); {
		case t.Typ == lex.IDENTIFIER:
			p.errorf("Invalid paren expression closed expression followed by identifier at line %d:%d with token '%s' in file : %s", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		case t.Typ == lex.INT || t.Typ == lex.FLOAT:
			p.errorf("Invalid paren expression closed expression followed by literal at line %d:%d with token '%s' in file : %s", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		case t.IsOperator():
			binary := &ast.BinaryExpr{
				X:expr.X,
				Op:t}
			expr.X = parseExpr(p, binary)
			return parseExpr(p, expr)
		case t.Typ == lex.LEFTPAREN:
			p.errorf("Invalid paren expression closed expression followed by opening parenthesis at line %d:%d with token '%s' in file : %s", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		case t.Typ == lex.RIGHTPAREN:
			expr.Rparen = t
			paren := p.pDepth.pop()
			if paren != expr {
				p.errorf("Internal error in parseExpr closing paren not matching current paren expr at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			}
			return parseExpr(p, expr)
		default:
			p.errorf("Invalid expression at line %d:%d with token '%s' in file %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	} else if expr.Rparen.Val == ")" {
		if expr.X == nil {
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
			binary := &ast.BinaryExpr{
				X:expr,
				Op:t}
			return parseExpr(p, binary)
		case t.Typ == lex.LEFTPAREN:
			p.errorf("Invalid paren expression closed expression followed by opening parenthesis at line %d:%d with token '%s' in file %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
			return nil
		case t.Typ == lex.RIGHTPAREN:
			// close enclosing paren in case parenExpr{X:parenExpr{}}
			//expr.Rparen = t
			paren := p.pDepth.pop()
			paren.Rparen = t
			return expr
		case atTerminator(t):
			return expr
		default:
			p.errorf("Invalid expression at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	} else {
		p.errorf("Internal error in parseParenExpr: at line %d:%d with token '%s' in file: %s\n", p.lineNumber(), p.lastToken.Pos, p.lastToken.Val, p.name)
		return nil
	}
	return nil
}

func parseUnaryExpr(p *Parser, expr *ast.UnaryExpr) ast.Expr {
	if expr.X == nil {
		switch t := p.next(); {
		case t.Typ == lex.IDENTIFIER:
			p.backup()
			ident := parseIdent(p)
			expr.X = ident
			return parseExpr(p, expr)
		case t.Typ == lex.INT || t.Typ == lex.FLOAT:
			bLit := &ast.BasicLit{Tok:t}
			expr.X = bLit
			return parseExpr(p, expr)
		case t.Typ == lex.ADD || t.Typ == lex.SUB:
			unary := &ast.UnaryExpr{Op:t}
			expr.X = parseExpr(p, unary)
			return expr
		case t.Typ == lex.LEFTPAREN:
			paren := newParenExpr(p, t)
			expr.X = parseExpr(p, paren)
			return expr
		default:
			p.errorf("Invalid unary expression at line %d:%d with token '%s' in file : 5S\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	} else {
		switch t:= p.next(); {
		case atTerminator(t):
			return expr
		case t.Typ == lex.RIGHTPAREN:
			paren := p.pDepth.pop()
			paren.Rparen = t
			return paren
		case t.IsOperator():
			if t.Precedence() > expr.Op.Precedence() {
				binary := &ast.BinaryExpr{X:expr.X, Op:t}
				expr.X = parseExpr(p, binary)
				return expr
			} else {
				binary := &ast.BinaryExpr{X:expr, Op:t}
				return parseExpr(p, binary)
			}
		default:
			p.errorf("Invalid expression at line %d:%d with token '%s' in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	}
	return nil
}

func parseBinaryExpr(p *Parser, expr *ast.BinaryExpr) ast.Expr {
	if expr.Y == nil {
		switch t := p.nextNonNewline(); {
		case t.Typ == lex.ADD || t.Typ == lex.SUB:
			unary := &ast.UnaryExpr{Op:t}
			expr.Y = parseExpr(p, unary)
			return expr
		case t.Typ == lex.IDENTIFIER:
			p.backup()
			ident := parseIdent(p)
			expr.Y = ident
			return parseExpr(p, expr)
		case t.Typ == lex.INT || t.Typ == lex.FLOAT:
			bLit := &ast.BasicLit{Tok:t}
			expr.Y = bLit
			return parseExpr(p, expr)
		case t.Typ == lex.LEFTPAREN:
			paren := newParenExpr(p, t)
			expr.Y = parseExpr(p, paren)
			return expr
		default:
			p.errorf("Invalid expression at line %d:%d with token '%s' in file : %S\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	} else {
		switch t := p.next(); {
		case atTerminator(t):
			return expr
		case t.IsOperator():
			if t.Precedence() > expr.Op.Precedence() {
				binary := &ast.BinaryExpr{X:expr.Y, Op:t}
				expr.Y = parseExpr(p, binary)
				return expr
			} else {
				binary := &ast.BinaryExpr{X:expr, Op:t}
				return parseExpr(p, binary)
			}
		default:
			p.errorf("Invalid binary expression at line %d:%d with token %v in file : %s\n", p.lineNumber(), t.Pos, t.Val, p.name)
		}
	}
	return nil
}
