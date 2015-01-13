package parse

import (
	// "fmt"
	"testing"
	"jon/calc/lex"
	"jon/calc/ast"
	"github.com/davecgh/go-spew/spew"
)

func TestSimpleBinaryAdd(t *testing.T) {
	input := `4+4`
	parser := Parse("TestSimpleBinaryAdd", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X: &ast.BasicLit{Tok:lex.Token{Typ:lex.INT, Val:"4"}},
			Op:lex.Token{Typ:lex.ADD, Val:"+"},
			Y: &ast.BasicLit{
				Tok: lex.Token{
					Typ:lex.INT,
					Val:"4",
				},
			},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
	}
}

func TestSimpleUnary(t *testing.T) {
	input := `-2`
	parser := Parse("TestSimpleUnary", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.UnaryExpr{
			Op:lex.Token{Typ:lex.SUB, Val:"-"},
			X: &ast.BasicLit{Tok:lex.Token{Typ:lex.INT, Val:"2"}},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
	}
}

func TestSimpleParen(t *testing.T) {
	input := `4+(4)`
	parser := Parse("TestSimpleParen", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X: &ast.BasicLit{Tok:lex.Token{Typ:lex.INT, Val:"4"}},
			Op:lex.Token{Typ:lex.ADD, Val:"+"},
			Y: &ast.ParenExpr{Lparen: lex.Token{Typ:lex.LEFTPAREN, Val:"("},
				X: &ast.BasicLit{
					Tok: lex.Token{
						Typ:lex.INT,
						Val:"4",
					},
				},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val:")"},
			},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
	}
}

func TestCompositeParen(t *testing.T) {
	input := `(( 4 ) + ( 4 ))`
	parser := Parse("TestCompositeParen", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.ParenExpr{
			Lparen: lex.Token{Typ: lex.LEFTPAREN, Val:"("},
			Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val:")"},
			X: &ast.BinaryExpr{
				X: &ast.ParenExpr{
					Lparen: lex.Token{Typ: lex.LEFTPAREN, Val:"("},
					Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val:")"},
					X: &ast.BasicLit{ Tok:lex.Token{Typ:lex.INT, Val:"4"}},
				},
				Op:lex.Token{Typ:lex.ADD, Val:"+"},
				Y: &ast.ParenExpr{
					Lparen: lex.Token{Typ:lex.LEFTPAREN, Val:"("},
					Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val:")"},
					X: &ast.BasicLit{ Tok: lex.Token{Typ:lex.INT, Val:"4"}},
				},
			},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
	}
}

func TestMultiExpr(t *testing.T) {
	input :=
		`
4+(4)
4*4
`
	parser := Parse("TestMultipleExpr", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X: &ast.BasicLit{Tok:lex.Token{Typ:lex.INT, Val:"4"}},
			Op:lex.Token{Typ:lex.ADD, Val:"+"},
			Y: &ast.ParenExpr{Lparen: lex.Token{Typ:lex.LEFTPAREN, Val:"("},
				X: &ast.BasicLit{
					Tok: lex.Token{
						Typ:lex.INT,
						Val:"4",
					},
				},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val:")"},
			},
		},
		&ast.BinaryExpr{
			X: &ast.BasicLit{Tok:lex.Token{Typ:lex.INT, Val:"4"}},
			Op:lex.Token{Typ:lex.MUL, Val:"*"},
			Y: &ast.BasicLit{
				Tok: lex.Token{ Typ:lex.INT, Val:"4"},
			},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
	}
}
