package main

import (
	"jon/calc/lex"
	"io/ioutil"
	"fmt"
)

func main() {
	input, err := ioutil.ReadFile("test_input/test2.calc")
	if err != nil {
		fmt.Println("Error reading file: \n %s", err)
	}
	lexer := lex.Lex("mytest", string(input))

	for {
		item := lexer.NextItem()
		fmt.Printf("%s ", item)
		if item.Typ == lex.EOF {
			fmt.Println()
			break
		}
	}
}
