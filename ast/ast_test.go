package ast

import (
	// "fmt"
	"jon/calc/lex"
	"testing"
)

func TestSimpleInsert(t *testing.T) {
	var testingTree Expr = &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}}
	testingTree = InsertExpr(testingTree, &BinaryExpr{Op: lex.Token{Typ: lex.ADD, Val: "+"}})
	testingTree = InsertExpr(&ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("}})
	testingTree = InsertExpr(&ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("}})

	expertedTree :=
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
	if !ast.Equals(testingTree, expectedTree) {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n", spew.Sdump(expected), spew.Sdump(output))
	}
}

func TestPrecedenceInsert(t *testing.T) {
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
