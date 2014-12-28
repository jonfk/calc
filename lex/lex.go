package lex

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Pos represents a byte position in the original input text
type Pos int

func (p Pos) Position() Pos {
	return p
}

// token represents a token or text string returned from the scanner.
type Token struct {
	Typ TokenType // The type of this token.
	Pos Pos      // The starting position, in bytes, of this item in the input string.
	Val string   // The value of this item.
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
	ERROR        TokenType = iota // error occurred; value is text of error
	BOOL                         // boolean constant
	EOF
	NEWLINE                      // '\n'
	LINECOMMENT                  // // ..... includes symbol
	BLOCKCOMMENT                 // /* block comment includes surrounding symbols*/
	LEFTPAREN  // '('
	NUMBER     // simple number, including imaginary
	COMPLEX                      // complex constant (1+2i); imaginary is just a number
	RIGHTPAREN // ')'
	IDENTIFIER // alphanumeric identifier not starting with '.'
	// Keywords appear after all the rest.
	KEYWORD  // used only to delimit the keywords
	ELSE     // else keyword
	END      // end keyword
	IF       // if keyword
	THEN     // then keyword

	OPERATOR
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	LAND  // &&
	LOR   // ||

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
)

const eof = -1

var key = map[string]TokenType{
	"else":     ELSE,
	"end":      END,
	"if":       IF,
	"then":     THEN,
	"+":        ADD,
	"-":        SUB,
	"*":        MUL,
	"/":        QUO,
	"%":        REM,
	"&&":     LAND,
	"||":      LOR,
	"==":       EQL,
	"<":        LSS,
	">":        GTR,
	"=":        ASSIGN,
	"!":        NOT,
	"!=":       NEQ,
	"<=":       LEQ,
	">=":       GEQ,
}

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// lexer holds the state of the scanner.
type Lexer struct {
	name       string    // the name of the input; used only for error reports
	input      string    // the string being scanned
	state      stateFn   // the next lexing function to enter
	pos        Pos       // current position in the input
	start      Pos       // start position of this item
	width      Pos       // width of last rune read from input
	lastPos    Pos       // position of most recent item returned by nextItem
	items      chan Token // channel of scanned items
	parenDepth int       // nesting depth of ( ) exprs
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *Lexer) emit(t TokenType) {
	l.items <- Token{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *Lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Token{ERROR, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *Lexer) NextItem() Token {
	token := <-l.items
	l.lastPos = token.Pos
	return token
}

// lex creates a new scanner for the input string.
func Lex(name, input string) *Lexer {
	l := &Lexer{
		name:       name,
		input:      input,
		items:      make(chan Token),
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *Lexer) run() {
	for l.state = lexStart; l.state != nil; {
		l.state = l.state(l)
	}
}

// state functions

func lexStart(l *Lexer) stateFn {
	switch r := l.next(); {
	case r == eof:
		l.emit(EOF)
		return nil
	case isSpace(r):
		return lexSpace
	case isEndOfLine(r):
		return lexEndOfLine
	case ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	case isAlphaNumeric(r):
		l.backup()
		return lexIdentifier
	case isSpecialSym(r):
		l.backup()
		return lexOperator
	case r == '(':
		l.emit(LEFTPAREN)
		l.parenDepth++
		return lexStart
	case r == ')':
		l.emit(RIGHTPAREN)
		l.parenDepth--
		if l.parenDepth < 0 {
			return l.errorf("unexpected right paren %#U", r)
		}
		return lexStart
	default:
		return l.errorf("unknown syntax: %q", l.input[l.start:l.pos])
	}
}


// lexSpace scans a run of space characters.
// One space has already been seen.
func lexSpace(l *Lexer) stateFn {
	for isSpace(l.peek()) {
		l.next()
	}
	l.ignore()
	return lexStart
}

// lexEndOfLine scans a end of line character.
func lexEndOfLine(l *Lexer) stateFn {
	for isEndOfLine(l.peek()) {
		l.next()
	}
	l.emit(NEWLINE)
	return lexStart
}

// lexNumber scans a number: decimal, octal, hex, float, or imaginary. This
// isn't a perfect number scanner - for instance it accepts "." and "0x0.2"
// and "089" - but when it's wrong the input is invalid and the parser (via
// strconv) will notice.
func lexNumber(l *Lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
	}
	mark := l.pos
	if sign := l.peek(); sign == '+' || sign == '-' {
		// Complex: 1+2i. No spaces, must end in 'i'.
		if !l.scanNumber() || l.input[l.pos-1] != 'i' {
			l.pos = mark
			l.emit(NUMBER)
			return lexStart
			//return l.errorf("bad number syntax: %q", l.input[l.start:l.pos])
		}
		l.emit(COMPLEX)
	} else {
		l.emit(NUMBER)
	}
	return lexStart
}

func (l *Lexer) scanNumber() bool {
	// Is it hex?
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// Is it imaginary?
	l.accept("i")
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(l.peek()) {
		l.next()
		return false
	}
	return true
}

// lexIdentifier scans an alphanumeric.
func lexIdentifier(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case isAlphaNumeric(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			if !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}
			switch {
			case key[word] > KEYWORD:
				l.emit(key[word])
			case word == "true", word == "false":
				l.emit(BOOL)
			default:
				l.emit(IDENTIFIER)
			}
			break Loop
		}
	}
	return lexStart
}

// atTerminator reports whether the input is at valid termination character to
// appear after an identifier. Breaks .X.Y into two pieces. Also catches cases
// like "$x+2" not being acceptable without a space, in case we decide one
// day to implement arithmetic.
func (l *Lexer) atTerminator() bool {
	r := l.peek()
	if isSpace(r) || isEndOfLine(r) {
		return true
	}
	switch r {
	case eof, '.', ',', ':', ')', '(':
		return true
	}
	return false
}

// lexOperator scans an special characters.
func lexOperator(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case r == '/' && l.peek() == '/':
			return lexLineComment
		case r == '/' && l.peek() == '*':
			return lexBlockComment
		case isSpecialSym(r):
			// absorb.
		default:
			l.backup()
			word := l.input[l.start:l.pos]
			// Is it necessary? But cannot be used or will error between arithmetic e.g. 2+2 as 2 is not a terminator
			// if !l.atTerminator() {
			// 	return l.errorf("bad character %#U", r)
			// }
			switch {
			case key[word] > OPERATOR:
				l.emit(key[word])
			default:
				return l.errorf("bad character %#U", r)
			}
			break Loop
		}
	}
	return lexStart
}

func lexLineComment(l *Lexer) stateFn {
Loop:
	for {
		switch r := l.next(); {
		case !isEndOfLine(r):
			// absorb.
		default:
			l.backup()
			if !l.atTerminator() {
				return l.errorf("bad character %#U", r)
			}
			l.emit(LINECOMMENT)
			break Loop
		}
	}
	return lexStart
}

func lexBlockComment(l *Lexer) stateFn {
Loop:
	for {
		// if we find '*' and the next is  '/'
		switch r := l.next(); {
		case !l.atEndBlockComment():
			// absorb.
		case r == eof:
			return l.errorf("Non-terminating block comment at %#U", r)
		default:
			// l.backup()
			// l.next()
			word := l.input[l.start:l.pos]
			switch {
			case strings.Index(word, "*/")== len(word)-len("*/"):
				l.emit(BLOCKCOMMENT)
			default:
				return l.errorf("error in  block comment at %#U", r)
			}
			break Loop
		}
	}
	return lexStart
}

func (l *Lexer) atEndBlockComment() bool {
	word := l.input[l.pos-2:l.pos]
	if strings.Index(word, "*/")== len(word)-len("*/") {
		return true
	}
	return false
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

// isSpecialSym reports whether r is a special symbol used as operators.
func isSpecialSym(r rune) bool {
	if strings.IndexRune("+-*/%&|=><!", r) >= 0 {
		return true
	}
	return false
}
