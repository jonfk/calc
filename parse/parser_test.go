package parse

import (
	// "fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jonfk/calc/ast"
	"github.com/jonfk/calc/lex"
	"testing"
)

func TestSimpleBinaryAdd(t *testing.T) {
	input := `4+4`
	parser := Parse("TestSimpleBinaryAdd", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			Y: &ast.BasicLit{
				Tok: lex.Token{
					Typ: lex.INT,
					Val: "4",
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
			Op: lex.Token{Typ: lex.SUB, Val: "-"},
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}},
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
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			Y: &ast.ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				X: &ast.BasicLit{
					Tok: lex.Token{
						Typ: lex.INT,
						Val: "4",
					},
				},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
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
			Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
			Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			X: &ast.BinaryExpr{
				X: &ast.ParenExpr{
					Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
					Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
					X:      &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
				},
				Op: lex.Token{Typ: lex.ADD, Val: "+"},
				Y: &ast.ParenExpr{
					Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
					Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
					X:      &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
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

func TestAssociativityADDSUB(t *testing.T) {
	input := `7-4+2`
	parser := Parse("TestAssociativityADDSUB", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X: &ast.BinaryExpr{
				X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "7"}},
				Op: lex.Token{Typ: lex.SUB, Val: "-"},
				Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			},
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
	}
}

func TestSimpleArithmeticPrecedence(t *testing.T) {
	input := `a+b*2-3/4%a`
	parser := Parse("TestSimpleArithmeticPrecedence", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X: &ast.BinaryExpr{
				X:  &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "a"}},
				Op: lex.Token{Typ: lex.ADD, Val: "+"},
				Y: &ast.BinaryExpr{
					X:  &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "b"}},
					Op: lex.Token{Typ: lex.MUL, Val: "*"},
					Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}},
				},
			},
			Op: lex.Token{Typ: lex.SUB, Val: "-"},
			Y: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "3"}},
					Op: lex.Token{Typ: lex.QUO, Val: "/"},
					Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
				},
				Op: lex.Token{Typ: lex.REM, Val: "%"},
				Y:  &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "a"}},
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
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			Y: &ast.ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				X: &ast.BasicLit{
					Tok: lex.Token{
						Typ: lex.INT,
						Val: "4",
					},
				},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			},
		},
		&ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.MUL, Val: "*"},
			Y: &ast.BasicLit{
				Tok: lex.Token{Typ: lex.INT, Val: "4"},
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

func TestMultiExprSemiColonSeperated(t *testing.T) {
	input :=
		`
4+(4);4*4;((100)/90);-aTest;
`
	parser := Parse("TestMultiExprSemiColonSeperated", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			Y: &ast.ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				X: &ast.BasicLit{
					Tok: lex.Token{
						Typ: lex.INT,
						Val: "4",
					},
				},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			},
		},
		&ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.MUL, Val: "*"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
		},
		&ast.ParenExpr{
			Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
			Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			X: &ast.BinaryExpr{
				X: &ast.ParenExpr{
					Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
					Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
					X:      &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "100"}},
				},
				Op: lex.Token{Typ: lex.QUO, Val: "/"},
				Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "90"}},
			},
		},
		&ast.UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"},
			X: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aTest"}},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestMultiExprMixedSemiColonSeperated(t *testing.T) {
	input :=
		`
4+(4)
         4*4
   ((100)/90);-aTest;
(aoeu)/2222
`
	parser := Parse("TestMultiExprMixedSemiColonSeperated", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			Y: &ast.ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				X: &ast.BasicLit{
					Tok: lex.Token{
						Typ: lex.INT,
						Val: "4",
					},
				},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			},
		},
		&ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.MUL, Val: "*"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
		},
		&ast.ParenExpr{
			Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
			Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			X: &ast.BinaryExpr{
				X: &ast.ParenExpr{
					Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
					Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
					X:      &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "100"}},
				},
				Op: lex.Token{Typ: lex.QUO, Val: "/"},
				Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "90"}},
			},
		},
		&ast.UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"},
			X: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aTest"}},
		},
		&ast.BinaryExpr{
			X: &ast.ParenExpr{
				Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
				X:      &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aoeu"}},
			},
			Op: lex.Token{Typ: lex.QUO, Val: "/"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2222"}},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestMultiLineExprs(t *testing.T) {
	input :=
		`
4+
(4)
         4+p*
4
   ((100)
/90);-aTest;
(aoeu)/2222
`
	parser := Parse("TestMultiLineExpressions", input)

	output := parser.File
	nodeList := []ast.Node{
		&ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			Y: &ast.ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				X: &ast.BasicLit{
					Tok: lex.Token{
						Typ: lex.INT,
						Val: "4",
					},
				},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			},
		},
		&ast.BinaryExpr{
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Y: &ast.BinaryExpr{
				Op: lex.Token{Typ: lex.MUL, Val: "*"},
				X:  &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "p"}},
				Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			},
		},
		&ast.ParenExpr{
			Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
			Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			X: &ast.BinaryExpr{
				X: &ast.ParenExpr{
					Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
					Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
					X:      &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "100"}},
				},
				Op: lex.Token{Typ: lex.QUO, Val: "/"},
				Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "90"}},
			},
		},
		&ast.UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"},
			X: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aTest"}},
		},
		&ast.BinaryExpr{
			X: &ast.ParenExpr{
				Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
				X:      &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aoeu"}},
			},
			Op: lex.Token{Typ: lex.QUO, Val: "/"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2222"}},
		},
	}
	expected := &ast.File{
		List: nodeList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}
