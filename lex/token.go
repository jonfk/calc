package lex

import (
	"fmt"
)

// Pos represents a byte position in the original input text
type Pos int

const NoPos Pos = 0

func (p Pos) Position() Pos {
	return p
}

// token represents a token or text string returned from the scanner.
type Token struct {
	Typ TokenType // The type of this token.
	Pos Pos       // The starting position, in bytes, of this item in the input string.
	Val string    // The value of this item.
}

func (t Token) String() string {
	switch {
	case t.Typ == EOF:
		return "EOF"
	case t.Typ == ERROR:
		return t.Val
	case t.Typ > KEYWORD && t.Typ < OPERATOR:
		return fmt.Sprintf("<%s>", t.Val)
	case t.Typ > OPERATOR:
		return fmt.Sprintf("[%s]", t.Val)
	case len(t.Val) > 10:
		return fmt.Sprintf("%.10q...", t.Val)
	}
	return fmt.Sprintf("%q", t.Val)
}

type TokenType int

const (
	ERROR TokenType = iota // error occurred; value is text of error
	BOOL                   // boolean constant
	EOF
	NEWLINE      // '\n'
	LINECOMMENT  // // ..... includes symbol
	BLOCKCOMMENT // /* block comment includes surrounding symbols*/
	LEFTPAREN    // '('
	//NUMBER       // simple number, including imaginary
	INT        // an int
	FLOAT      // a float
	STRING     // a string literal
	RIGHTPAREN // ')'
	SEMICOLON  // ';'
	COMMA      // ','
	IDENTIFIER // alphanumeric identifier not starting with '.'
	// Keywords appear after all the rest.
	KEYWORD // used only to delimit the keywords
	ELSE    // else keyword
	END     // end keyword
	IF      // if keyword
	THEN    // then keyword
	LET     // let keyword
	VAR     // var keyword
	VAL     // val keyword

	OPERATOR
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	LAND // &&
	LOR  // ||

	EQL // ==
	LSS // <
	GTR // >
	NOT // !

	NEQ // !=
	LEQ // <=
	GEQ // >=

	ASSIGN // = not a keyword or operator since it does not yield an expression
)

const eof = -1

var key = map[string]TokenType{
	"else": ELSE,
	"end":  END,
	"if":   IF,
	"then": THEN,
	"let":  LET,
	"val":  VAL,
	"var":  VAR,
	"+":    ADD,
	"-":    SUB,
	"*":    MUL,
	"/":    QUO,
	"%":    REM,
	"&&":   LAND,
	"||":   LOR,
	"==":   EQL,
	"<":    LSS,
	">":    GTR,
	"=":    ASSIGN,
	"!":    NOT,
	"!=":   NEQ,
	"<=":   LEQ,
	">=":   GEQ,
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func (op Token) Precedence() int {
	switch op.Typ {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB:
		return 4
	case MUL, QUO, REM:
		return 5
	}
	return LowestPrec
}

// Predicates

func (tok Token) IsLiteral() bool {
	switch tok.Typ {
	case INT, FLOAT, IDENTIFIER:
		return true
	default:
		return false
	}
}

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Token) IsOperator() bool { return tok.Typ > OPERATOR && tok.Typ != ASSIGN }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Token) IsKeyword() bool { return tok.Typ > KEYWORD && tok.Typ < OPERATOR }

// Compares Typ and Val but not position
// Used for debugging and testing
func (t Token) Equals(ot Token) bool { return t.Val == ot.Val && t.Typ == ot.Typ }
