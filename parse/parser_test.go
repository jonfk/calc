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
	stmtList := []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.BinaryExpr{
				X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
				Op: lex.Token{Typ: lex.ADD, Val: "+"},
				Y: &ast.BasicLit{
					Tok: lex.Token{
						Typ: lex.INT,
						Val: "4",
					},
				},
			},
		},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestSimpleUnary(t *testing.T) {
	input := `-2`
	parser := Parse("TestSimpleUnary", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.UnaryExpr{
				Op: lex.Token{Typ: lex.SUB, Val: "-"},
				X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}},
			},
		},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestSimpleParen(t *testing.T) {
	input := `4+(4)`
	parser := Parse("TestSimpleParen", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.BinaryExpr{
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
		},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestCompositeParen(t *testing.T) {
	input := `(( 4 ) + ( 4 ))`
	parser := Parse("TestCompositeParen", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.ParenExpr{
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
		},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestAssociativityADDSUB(t *testing.T) {
	input := `7-4+2`
	parser := Parse("TestAssociativityADDSUB", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "7"}},
					Op: lex.Token{Typ: lex.SUB, Val: "-"},
					Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
				},
				Op: lex.Token{Typ: lex.ADD, Val: "+"},
				Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}},
			},
		},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestSimpleArithmeticPrecedence(t *testing.T) {
	// a+b*2-3/4%a = (a + (b*2)) - ((3/4)%a)
	input := `a+b*2-3/4%a`
	parser := Parse("TestSimpleArithmeticPrecedence", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{
			X: &ast.BinaryExpr{
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
		},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
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
	stmtList := []ast.Stmt{
		&ast.ExprStmt{X: &ast.BinaryExpr{
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
		},
		&ast.ExprStmt{X: &ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.MUL, Val: "*"},
			Y: &ast.BasicLit{
				Tok: lex.Token{Typ: lex.INT, Val: "4"},
			},
		},
		},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestMultiExprSemiColonSeperated(t *testing.T) {
	input :=
		`
4+(4);4*4;((100)/90);-aTest;
`
	parser := Parse("TestMultiExprSemiColonSeperated", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{X: &ast.BinaryExpr{
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
		}},
		&ast.ExprStmt{X: &ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.MUL, Val: "*"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
		}},
		&ast.ExprStmt{X: &ast.ParenExpr{
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
		}},
		&ast.ExprStmt{X: &ast.UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"},
			X: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aTest"}},
		}},
	}
	expected := &ast.File{
		List: stmtList,
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
	stmtList := []ast.Stmt{
		&ast.ExprStmt{X: &ast.BinaryExpr{
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
		}},
		&ast.ExprStmt{X: &ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Op: lex.Token{Typ: lex.MUL, Val: "*"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
		}},
		&ast.ExprStmt{X: &ast.ParenExpr{
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
		}},
		&ast.ExprStmt{X: &ast.UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"},
			X: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aTest"}},
		}},
		&ast.ExprStmt{X: &ast.BinaryExpr{
			X: &ast.ParenExpr{
				Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
				X:      &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aoeu"}},
			},
			Op: lex.Token{Typ: lex.QUO, Val: "/"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2222"}},
		}},
	}
	expected := &ast.File{
		List: stmtList,
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
         4+p*4
   ((100)/
90);-aTest;
(aoeu)/2222
`
	parser := Parse("TestMultiLineExpressions", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{X: &ast.BinaryExpr{
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
		}},
		&ast.ExprStmt{X: &ast.BinaryExpr{
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Y: &ast.BinaryExpr{
				Op: lex.Token{Typ: lex.MUL, Val: "*"},
				X:  &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "p"}},
				Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			},
		}},
		&ast.ExprStmt{X: &ast.ParenExpr{
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
		}},
		&ast.ExprStmt{X: &ast.UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"},
			X: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aTest"}},
		}},
		&ast.ExprStmt{X: &ast.BinaryExpr{
			X: &ast.ParenExpr{
				Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
				Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
				X:      &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "aoeu"}},
			},
			Op: lex.Token{Typ: lex.QUO, Val: "/"},
			Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2222"}},
		}},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestSimpleBoolExpr(t *testing.T) {
	input :=
		`true || false && test`
	parser := Parse("TestSimpleBoolExpr", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{X: &ast.BinaryExpr{
			X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.BOOL, Val: "true"}},
			Op: lex.Token{Typ: lex.LOR, Val: "||"},
			Y: &ast.BinaryExpr{
				X:  &ast.BasicLit{Tok: lex.Token{Typ: lex.BOOL, Val: "false"}},
				Op: lex.Token{Typ: lex.LAND, Val: "&&"},
				Y:  &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "test"}},
			},
		}},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestMultiLineParenExpr(t *testing.T) {
	input :=
		`((10)
*1
/9)`
	parser := Parse("TestMultiLineParenExpr", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.ExprStmt{X: &ast.ParenExpr{
			Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
			Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
			X: &ast.BinaryExpr{
				X: &ast.BinaryExpr{
					X: &ast.ParenExpr{
						Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
						Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
						X:      &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "10"}},
					},
					Op: lex.Token{Typ: lex.MUL, Val: "*"},
					Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "1"}},
				},
				Op: lex.Token{Typ: lex.QUO, Val: "/"},
				Y:  &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "9"}}},
		}},
	}
	expected := &ast.File{
		List: stmtList,
	}
	if !ast.Equals(parser.File, expected) {
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestSimpleAssignStmt(t *testing.T) {
	input := `c = 9`
	parser := Parse("TestSimpleAssignStmt", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "c"}},
			Tok: lex.Token{Typ: lex.ASSIGN, Val: "="},
			Rhs: &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "9"}},
		},
	}
	expected := &ast.File{
		List: stmtList,
	}

	// if !ast.Equals(parser.File.List[0].(*ast.AssignStmt).Rhs, expected.List[0].(*ast.AssignStmt).Rhs) {
	// 	t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	// }

	if !ast.Equals(parser.File, expected) {
		//t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestMultyAssignStmt(t *testing.T) {
	input := `a = true;b=10;c=b`
	parser := Parse("TestSimpleAssignStmt", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.AssignStmt{
			Lhs: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "a"}},
			Tok: lex.Token{Typ: lex.ASSIGN, Val: "="},
			Rhs: &ast.BasicLit{Tok: lex.Token{Typ: lex.BOOL, Val: "true"}},
		},
		&ast.AssignStmt{
			Lhs: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "b"}},
			Tok: lex.Token{Typ: lex.ASSIGN, Val: "="},
			Rhs: &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "10"}},
		},
		&ast.AssignStmt{
			Lhs: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "c"}},
			Tok: lex.Token{Typ: lex.ASSIGN, Val: "="},
			Rhs: &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "b"}},
		},
	}
	expected := &ast.File{
		List: stmtList,
	}

	if !ast.Equals(parser.File, expected) {
		//t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}

func TestSimpleDeclStmt(t *testing.T) {
	input := `var c = 9`
	parser := Parse("TestSimpleDeclStmt", input)

	output := parser.File
	stmtList := []ast.Stmt{
		&ast.DeclStmt{
			Decl: &ast.GenDecl{
				Tok: lex.Token{Typ: lex.VAR, Val: "var"},
				Spec: &ast.ValueSpec{
					Name:  &ast.Ident{Tok: lex.Token{Typ: lex.IDENTIFIER, Val: "c"}},
					Value: &ast.BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "9"}},
				},
			},
		},
	}
	expected := &ast.File{
		List: stmtList,
	}

	if !ast.Equals(parser.File, expected) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
		// t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", expected.String(), output.String())
	}
}
