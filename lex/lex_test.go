package lex

import (
	"testing"
	"fmt"
)

func TestAdd(t *testing.T) {
	input := "4+4"
	lexer := Lex("TestAdd", input)
	var output []Token
	expected := []Token{
		Token{Typ:NUMBER,Val:"4"},
		Token{Typ:ADD,Val:"+"},
		Token{Typ:NUMBER,Val:"4"},
		Token{Typ:EOF,Val:""},
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
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
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
		Token{Typ:NUMBER,Val:"3.1"},
		Token{Typ:SUB,Val:"-"},
		Token{Typ:NUMBER,Val:"2.0"},
		Token{Typ:NUMBER,Val:"64."},
		Token{Typ:MUL,Val:"*"},
		Token{Typ:NUMBER,Val:"9.0"},
		Token{Typ:NUMBER,Val:"10."},
		Token{Typ:REM,Val:"%"},
		Token{Typ:NUMBER,Val:"2."},
		Token{Typ:NUMBER,Val:"9.9"},
		Token{Typ:QUO,Val:"/"},
		Token{Typ:NUMBER,Val:"3.1"},
		Token{Typ:EOF,Val:""},
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
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
		}
	}
}

func TestIfThenElse(t *testing.T) {
	input := `if true then 8 else 10 end`
	lexer := Lex("TestIfThenElse", input)
	var output []Token
	expected := []Token{
		Token{Typ:IF,Val:"if"},
		Token{Typ:BOOL,Val:"true"},
		Token{Typ:THEN,Val:"then"},
		Token{Typ:NUMBER,Val:"8"},
		Token{Typ:ELSE,Val:"else"},
		Token{Typ:NUMBER,Val:"10"},
		Token{Typ:END,Val:"end"},
		Token{Typ:EOF,Val:""},
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
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
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
		Token{Typ:NUMBER,Val:"9."},
		Token{Typ:GEQ,Val:">="},
		Token{Typ:NUMBER,Val:"9."},
		Token{Typ:NUMBER,Val:"8"},
		Token{Typ:LEQ,Val:"<="},
		Token{Typ:NUMBER,Val:"2"},
		Token{Typ:NUMBER,Val:"8"},
		Token{Typ:LSS,Val:"<"},
		Token{Typ:NUMBER,Val:"10"},
		Token{Typ:NUMBER,Val:"1."},
		Token{Typ:GTR,Val:">"},
		Token{Typ:NUMBER,Val:"2"},
		Token{Typ:BOOL,Val:"true"},
		Token{Typ:EQL,Val:"=="},
		Token{Typ:BOOL,Val:"true"},
		Token{Typ:NUMBER,Val:"2."},
		Token{Typ:NEQ,Val:"!="},
		Token{Typ:NUMBER,Val:"2"},
		Token{Typ:EOF,Val:""},
	}
	for {
		item := lexer.NextItem()
		output = append(output, item)
		fmt.Printf("%s ", item)
		if item.Typ == EOF {
			fmt.Println()
			break
		}
	}
	if len(output) != len(expected) {
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
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
		Token{Typ:NOT,Val:"!"},
		Token{Typ:LEFTPAREN,Val:"("},
		Token{Typ:BOOL,Val:"false"},
		Token{Typ:RIGHTPAREN,Val:")"},
		Token{Typ:BOOL,Val:"true"},
		Token{Typ:LOR,Val:"||"},
		Token{Typ:LEFTPAREN,Val:"("},
		Token{Typ:BOOL,Val:"false"},
		Token{Typ:RIGHTPAREN,Val:")"},
		Token{Typ:BOOL,Val:"false"},
		Token{Typ:LAND,Val:"&&"},
		Token{Typ:BOOL,Val:"true"},
		Token{Typ:EOF,Val:""},
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
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
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
		Token{Typ:LINECOMMENT,Val:"//aoeu"},
		Token{Typ:LINECOMMENT,Val:"///*test*/"},
		Token{Typ:BLOCKCOMMENT,Val:"/*test*/"},
		Token{Typ:EOF,Val:""},
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
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
	}
	for i, item := range output {
		if item.Typ != expected[i].Typ || item.Val != expected[i].Val {
			// For more information add %#v e.g:
			// t.Logf("\nExpected: %#v\n Got:     %#v\n", output, expected)
			t.Errorf("\nExpected: %+v\n Got:     %+v\n", output, expected)
		}
	}
}
