package parse

import (
	"github.com/jonfk/calc/ast"
	// "fmt"
)

type ParenDepth struct {
	Depth int
	Stack []*ast.ParenExpr
}

func (p *ParenDepth) push(paren *ast.ParenExpr) {
	p.Stack = append(p.Stack, paren)
	p.Depth += 1
	// fmt.Printf("PUSHING : %d slice : %v\n",p.Depth, p.Stack)
}

func (p *ParenDepth) pop() *ast.ParenExpr {
	// fmt.Printf("POPPING : %d slice : %v\n",p.Depth, p.Stack)
	var popped *ast.ParenExpr
	popped, p.Stack = p.Stack[len(p.Stack)-1], p.Stack[:len(p.Stack)-1]
	p.Depth -= 1
	return popped
}
