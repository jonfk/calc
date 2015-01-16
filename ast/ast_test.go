package ast

import (
	// "fmt"
	"jon/calc/lex"
	"testing"
)

func TestSimpleInsert(t *testing.T) {
	// 4+(4)
	var err error
	var result Expr
	insert := func(tree, expr Expr) Expr {
		if err != nil {
			return nil
		}
		result, err = InsertExpr(tree, expr)
		return result
	}

	var testTree Expr = &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}}
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.ADD, Val: "+"}})
	paren := &ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("}}
	testTree = insert(testTree, paren)
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}})
	paren.Rparen = lex.Token{Typ: lex.RIGHTPAREN, Val: ")"}
	// fmt.Println(Sprint(testTree))

	expected := &BinaryExpr{
		X:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
		Op: lex.Token{Typ: lex.ADD, Val: "+"},
		Y: &ParenExpr{
			Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
			X:      &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
			Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
		},
	}
	if !Equals(testTree, expected) || err != nil {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\nWith error: %s\n", Sprint(expected), Sprint(testTree), err)
	}
}

func TestSimpleInsertIntoParens(t *testing.T) {
	// 4+((4)+100-(6.0))
	var err error
	var result Expr
	insert := func(tree, expr Expr) Expr {
		if err != nil {
			return nil
		}
		result, err = InsertExpr(tree, expr)
		return result
	}

	var testTree Expr = &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}}

	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.ADD, Val: "+"}})
	paren := &ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("}}
	testTree = insert(testTree, paren)
	paren2 := &ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("}}
	testTree = insert(testTree, paren2)
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}})
	paren2.Rparen = lex.Token{Typ: lex.RIGHTPAREN, Val: ")"}
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.ADD, Val: "+"}})
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "100"}})
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"}})
	paren3 := &ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("}}
	testTree = insert(testTree, paren3)
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.FLOAT, Val: "6.0"}})

	paren3.Rparen = lex.Token{Typ: lex.RIGHTPAREN, Val: ")"}
	paren.Rparen = lex.Token{Typ: lex.RIGHTPAREN, Val: ")"}
	//fmt.Println(Sprint(testTree))

	expected := &BinaryExpr{
		X:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
		Op: lex.Token{Typ: lex.ADD, Val: "+"},
		Y: &ParenExpr{
			Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
			X: &BinaryExpr{
				X: &BinaryExpr{
					X: &ParenExpr{
						Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
						X:      &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
						Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
					},
					Op: lex.Token{Typ: lex.ADD, Val: "+"},
					Y:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "100"}},
				},
				Op: lex.Token{Typ: lex.SUB, Val: "-"},
				Y: &ParenExpr{
					Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
					X:      &BasicLit{Tok: lex.Token{Typ: lex.FLOAT, Val: "6.0"}},
					Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
				},
			},
			Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
		},
	}
	if !Equals(testTree, expected) || err != nil {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n\nWith error: %s\n", Sprint(expected), Sprint(testTree), err)
	}
}

func TestInsertCompositePrecedence(t *testing.T) {
	// 1+2*3-4/5%6
	var err error
	var result Expr
	insert := func(tree, expr Expr) Expr {
		if err != nil {
			return nil
		}
		result, err = InsertExpr(tree, expr)
		return result
	}

	var testTree Expr = &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "1"}}
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.ADD, Val: "+"}})
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}})
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.MUL, Val: "*"}})
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "3"}})
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"}})
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}})
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.QUO, Val: "/"}})
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "5"}})
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.REM, Val: "%"}})
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "6"}})
	//fmt.Println(Sprint(testTree))

	expected := &BinaryExpr{
		X: &BinaryExpr{
			X:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "1"}},
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			Y: &BinaryExpr{
				X:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}},
				Op: lex.Token{Typ: lex.MUL, Val: "*"},
				Y:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "3"}},
			},
		},
		Op: lex.Token{Typ: lex.SUB, Val: "-"},
		Y: &BinaryExpr{
			X: &BinaryExpr{
				X:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "4"}},
				Op: lex.Token{Typ: lex.QUO, Val: "/"},
				Y:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "5"}},
			},
			Op: lex.Token{Typ: lex.REM, Val: "%"},
			Y:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "6"}},
		},
	}
	if !Equals(testTree, expected) || err != nil {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n\nWith error: %s\n", Sprint(expected), Sprint(testTree), err)
	}
}

func TestInsertUnary(t *testing.T) {
	// -+2*-(3)/8
	var err error
	var result Expr
	insert := func(tree, expr Expr) Expr {
		if err != nil {
			return nil
		}
		result, err = InsertExpr(tree, expr)
		return result
	}

	var testTree Expr = &UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"}}
	testTree = insert(testTree, &UnaryExpr{Op: lex.Token{Typ: lex.ADD, Val: "+"}})
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}})
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.MUL, Val: "*"}})
	testTree = insert(testTree, &UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"}})
	paren := &ParenExpr{Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("}}
	testTree = insert(testTree, paren)
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "3"}})
	paren.Rparen = lex.Token{Typ: lex.RIGHTPAREN, Val: ")"}
	testTree = insert(testTree, &BinaryExpr{Op: lex.Token{Typ: lex.QUO, Val: "/"}})
	testTree = insert(testTree, &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "8"}})
	//fmt.Println(Sprint(testTree))

	expected := &UnaryExpr{
		Op: lex.Token{Typ: lex.SUB, Val: "-"},
		X: &UnaryExpr{
			Op: lex.Token{Typ: lex.ADD, Val: "+"},
			X: &BinaryExpr{
				X: &BinaryExpr{
					X:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "2"}},
					Op: lex.Token{Typ: lex.MUL, Val: "*"},
					Y: &UnaryExpr{Op: lex.Token{Typ: lex.SUB, Val: "-"},
						X: &ParenExpr{
							Lparen: lex.Token{Typ: lex.LEFTPAREN, Val: "("},
							X:      &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "3"}},
							Rparen: lex.Token{Typ: lex.RIGHTPAREN, Val: ")"},
						},
					},
				},
				Op: lex.Token{Typ: lex.QUO, Val: "/"},
				Y:  &BasicLit{Tok: lex.Token{Typ: lex.INT, Val: "8"}},
			},
		},
	}
	if !Equals(testTree, expected) || err != nil {
		t.Errorf("\nExpected:\n%s\n\nGot:\n%s\n\nWith error: %s\n", Sprint(expected), Sprint(testTree), err)
	}
}
