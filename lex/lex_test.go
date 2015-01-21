package lex

import (
	// "fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	input := "4+4"
	lexer := Lex("TestAdd", input)
	var output []Token
	expected := []Token{
		Token{Typ: INT, Val: "4"},
		Token{Typ: ADD, Val: "+"},
		Token{Typ: INT, Val: "4"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestInts(t *testing.T) {
	input :=
		`10
0b101011
0B11111
0c77627
0C77272
0x009abc
0X0293ABC
`
	lexer := Lex("TestInts", input)
	var output []Token
	expected := []Token{
		Token{Typ: INT, Val: "10"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "0b101011"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "0B11111"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "0c77627"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "0C77272"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "0x009abc"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "0X0293ABC"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestFloats(t *testing.T) {
	input :=
		`3.1 2.0e10
99.`
	lexer := Lex("TestFloats", input)
	var output []Token
	expected := []Token{
		Token{Typ: FLOAT, Val: "3.1"},
		Token{Typ: FLOAT, Val: "2.0e10"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: FLOAT, Val: "99."},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestFloatArith(t *testing.T) {
	input :=
		`3.1-2.0
64.*9.0
10.%2.
9.9/3.1`
	lexer := Lex("TestFloatArith", input)
	var output []Token
	expected := []Token{
		Token{Typ: FLOAT, Val: "3.1"},
		Token{Typ: SUB, Val: "-"},
		Token{Typ: FLOAT, Val: "2.0"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: FLOAT, Val: "64."},
		Token{Typ: MUL, Val: "*"},
		Token{Typ: FLOAT, Val: "9.0"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: FLOAT, Val: "10."},
		Token{Typ: REM, Val: "%"},
		Token{Typ: FLOAT, Val: "2."},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: FLOAT, Val: "9.9"},
		Token{Typ: QUO, Val: "/"},
		Token{Typ: FLOAT, Val: "3.1"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestIfThenElse(t *testing.T) {
	input := `if true then 8 else 10 end`
	lexer := Lex("TestIfThenElse", input)
	var output []Token
	expected := []Token{
		Token{Typ: IF, Val: "if"},
		Token{Typ: BOOL, Val: "true"},
		Token{Typ: THEN, Val: "then"},
		Token{Typ: INT, Val: "8"},
		Token{Typ: ELSE, Val: "else"},
		Token{Typ: INT, Val: "10"},
		Token{Typ: END, Val: "end"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestNumComp(t *testing.T) {
	input := `9.>=9.
8 <= 2
8<10
1.>2
true==true
2.!=2
`
	lexer := Lex("TestNumComp", input)
	var output []Token
	expected := []Token{
		Token{Typ: FLOAT, Val: "9."},
		Token{Typ: GEQ, Val: ">="},
		Token{Typ: FLOAT, Val: "9."},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "8"},
		Token{Typ: LEQ, Val: "<="},
		Token{Typ: INT, Val: "2"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "8"},
		Token{Typ: LSS, Val: "<"},
		Token{Typ: INT, Val: "10"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: FLOAT, Val: "1."},
		Token{Typ: GTR, Val: ">"},
		Token{Typ: INT, Val: "2"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: BOOL, Val: "true"},
		Token{Typ: EQL, Val: "=="},
		Token{Typ: BOOL, Val: "true"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: FLOAT, Val: "2."},
		Token{Typ: NEQ, Val: "!="},
		Token{Typ: INT, Val: "2"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		// fmt.Printf("%s ", item)
		if item.Typ == EOF {
			// fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestLComp(t *testing.T) {
	input := `

!(false)
true||(false)
false&&true`
	lexer := Lex("TestLComp", input)
	var output []Token
	expected := []Token{
		Token{Typ: NEWLINE, Val: "\n\n"},
		Token{Typ: NOT, Val: "!"},
		Token{Typ: LEFTPAREN, Val: "("},
		Token{Typ: BOOL, Val: "false"},
		Token{Typ: RIGHTPAREN, Val: ")"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: BOOL, Val: "true"},
		Token{Typ: LOR, Val: "||"},
		Token{Typ: LEFTPAREN, Val: "("},
		Token{Typ: BOOL, Val: "false"},
		Token{Typ: RIGHTPAREN, Val: ")"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: BOOL, Val: "false"},
		Token{Typ: LAND, Val: "&&"},
		Token{Typ: BOOL, Val: "true"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestComments(t *testing.T) {
	input := `
//aoeu
///*test*/
/*test*/`
	lexer := Lex("TestComments", input)
	var output []Token
	expected := []Token{
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: LINECOMMENT, Val: "//aoeu"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: LINECOMMENT, Val: "///*test*/"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: BLOCKCOMMENT, Val: "/*test*/"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF || item.Typ == ERROR {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestSemiColons(t *testing.T) {
	input :=
		`10;
1;
`
	lexer := Lex("TestInts", input)
	var output []Token
	expected := []Token{
		Token{Typ: INT, Val: "10"},
		Token{Typ: SEMICOLON, Val: ";"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: INT, Val: "1"},
		Token{Typ: SEMICOLON, Val: ";"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	input := `"aoeu" 9 + 8`

	lexer := Lex("TestStringLiterals", input)
	var output []Token
	expected := []Token{
		Token{Typ: STRING, Val: "\"aoeu\""},
		Token{Typ: INT, Val: "9"},
		Token{Typ: ADD, Val: "+"},
		Token{Typ: INT, Val: "8"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF || item.Typ == ERROR {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}

func TestValVarAssignments(t *testing.T) {
	input := `val a, b = 10, "big";
var a
val b = true`

	lexer := Lex("TestValVarAssign", input)
	var output []Token
	expected := []Token{
		Token{Typ: VAL, Val: "val"},
		Token{Typ: IDENTIFIER, Val: "a"},
		Token{Typ: COMMA, Val: ","},
		Token{Typ: IDENTIFIER, Val: "b"},
		Token{Typ: ASSIGN, Val: "="},
		Token{Typ: INT, Val: "10"},
		Token{Typ: COMMA, Val: ","},
		Token{Typ: STRING, Val: "\"big\""},
		Token{Typ: SEMICOLON, Val: ";"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: VAR, Val: "var"},
		Token{Typ: IDENTIFIER, Val: "a"},
		Token{Typ: NEWLINE, Val: "\n"},
		Token{Typ: VAL, Val: "val"},
		Token{Typ: IDENTIFIER, Val: "b"},
		Token{Typ: ASSIGN, Val: "="},
		Token{Typ: BOOL, Val: "true"},
		Token{Typ: EOF, Val: ""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		//fmt.Printf("%s ", item)
		if item.Typ == EOF || item.Typ == ERROR {
			//fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
		t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", expected, output)
		}
	}
}
