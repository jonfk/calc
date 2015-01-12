package main

import (
	// "jon/calc/lex"
	//"io/ioutil"
	"fmt"
	// "jon/calc/ast"
	"jon/calc/parse"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	// input, err := ioutil.ReadFile("test_input/test2.calc")
	// if err != nil {
	// 	fmt.Println("Error reading file: \n %s", err)
	// }
	// lexer := lex.Lex("mytest", string(input))

	// for {
	// 	item := lexer.NextItem()
	// 	fmt.Printf("%s ", item)
	// 	if item.Typ == lex.EOF {
	// 		fmt.Println()
	// 		break
	// 	}
	// }

	// var tree *ast.File = new(ast.File)
	// expr := &ast.BasicLit{lex.Token{lex.INT, 0, "1"}}
	// tree.List = append(tree.List, expr)

	// var test ast.Expr
	// test = expr

	// switch test.(type) {
	// case *ast.BasicLit:
	// 	fmt.Println("basic lit")
	// case nil:
	// 	fmt.Println("nil")
	// default:
	// 	fmt.Println("unknown")
	// }

	// //fmt.Printf("%+v", tree)
	// fmt.Println("tree")
	// spew.Dump(test)


	// Test parser
	input :=
`4+4`
	parser := parse.Parse("TestAdd", input)
	fmt.Println("parser")
	spew.Dump(parser)
	fmt.Println("\ntree")
	spew.Dump(parser.File)
}
