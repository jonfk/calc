[![Build Status](https://travis-ci.org/jonfk/calc.svg)](https://travis-ci.org/jonfk/calc)
[![GoDoc](https://godoc.org/github.com/jonfk/calc?status.svg)](http://godoc.org/github.com/jonfk/calc)

calc
====

My toy calculator language. The lexer was inspired by the
[text/template/parse package](http://golang.org/pkg/text/template/parse/)
from the standard library and this talk by Rob Pike:
[Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE)
The parser is a custom recursive descent parser. The parsing support was
also inspired by go/parser and text/template/parse.

###TODO
- Add more tests for parser
- Add comment support to parsing
- Add eval package and implement an interpreter
- Add character literals
- Add support for val declarations
- Add support for let statements
- Add support for function literals
- Add support for function declarations
- Add support for lists
- Add typing system
- Vendor dependencies
- Make project go gettable with cannonical import paths based on where it is hosted

###Notes
- Identifiers can be alphanumeric with an underscore '_'
- Expressions can end with a ';' but ';' are not strictly necessary. They can be used
to disambiguate certain expressions.
- Unary Expressions cannot span multiple lines
- Binary Expressions can span multiple lines only if line ends with operator
- Operator precedence are left binding and as follows:

```
Highest(5): *, /, %
       (4): +, -
       (3): ==, !=, <, >, <=, >=
       (2): &&
Lower  (1): ||
Lowest (0): anything else

e.g 4+2/3 == 4 + (2/3)
    4-5+4%a+5 == ((4 - 5) + (4%a)) + 5
```

##Grammar in BNF

    expr ::= num_expr
           | bool_expr
           | if_expr
           | LPAREN expr RPAREN

    if_expr ::= IF bool_expr THEN expr ELSE expr END

    num_expr ::= NUMBER
               | IDENTIFIER
               | num_expr ADD num_expr
               | num_expr SUB num_expr
               | num_expr MUL num_expr
               | num_expr QUO num_expr
               | num_expr REM num_expr

    bool_expr ::= BOOL
                | IDENTIFIER
                | NOT bool_expr
                | bool_expr LAND bool_expr
                | bool_expr LOR bool_expr
                | num_expr EQL num_expr
                | num_expr LSS num_expr
                | num_expr GTR num_expr
                | num_expr NEQ num_expr
                | num_expr LEQ num_expr
                | num_expr GEQ num_expr

    ident_stmt ::= IDENTIFIER

    val_decl ::= val ident ASSIGN expr

    let_decl ::= LET val_decl IN block END

    block ::=

###Dependencies
Depedencies are kept to a minimum.
- https://github.com/davecgh/go-spew
```bash
# Used for testing and debugging
go get github.com/davecgh/go-spew/spew
```