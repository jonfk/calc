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
- Add support for val declarations
- Add support for var declarations
- Add if expressions
- Add eval package and implement an interpreter
- Add character literals
- Add support for let statements
- Add "=>" token to lexer
- keep parsing expression if in a paren
- Add support for function literals
- Add support for function declarations
- Add support for lists
- Add datatypes
- Add references and probably some for of gc
- Add records or structs?
- Add doc comment support to parsing
- Add typing system(?) or go with dynamic typing
- Add pattern matching(?)
- Vendor dependencies(or remove them)

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

##Grammar in EBNF

    literal = NUMBER
            | IDENTIFIER
            | BOOL


    if_expr = "if" , bool_expr , "then" , expr , "else" , expr "end"

    num_expr = "+" , expr
               | "-" , expr
               | expr , "+" , expr
               | expr , "-" , expr
               | expr , "*" , expr
               | expr , "/" , expr
               | expr , "%" , expr

    bool_expr = "!" , expr
                | expr , "&&" , expr
                | expr , "||" , expr
                | expr , "==" , expr
                | expr , ">" , expr
                | expr , "<" , expr
                | expr , "!=" , expr
                | expr , ">=" , expr
                | expr , "<=" , expr

    tuple_expr = "(" , expr , "," , expr , { "," , expr } , ")" # n > 1

    function = "fn" , "(" , ident_stmt , ")" , "=>" , expr , "end" # remove end keyword ?

    block = expr , { ("\n" | ";") , expr}

    let_expr = "let" , val_decl , "in" , expr , "end"

    expr = literal
           | num_expr
           | bool_expr
           | if_expr
           | "(" , expr , ")"
           | tuple_expr
           | let_expr
           | block
           | function
           | func_apcl

    func_apcl = (IDENTIFIER | function ) , "(" , expr , { "," , expr } , ")"

    ident_stmt = IDENTIFIER , { "," , IDENTIFIER }

    decl = val_decl
         | var_decl

    val_decl = "val" , ident_stmt , "=" , expr

    var_decl = "var" , ident_stmt , "=" , expr

    func_decl = "def" , IDENTIFIER , "(" , ident_stmt , ")" , "=" , expr , "end"


###Planned Extensions to grammar
- Add pattern matching to grammar
```
pattern = "(" , pat , "," , pat , { "," , pat } , ")" # n > 1
        | literal

val_decl = "val" , pat , "=" , expr
```

##Dependencies
Depedencies are kept to a minimum.
- https://github.com/davecgh/go-spew
```bash
# Used for testing and debugging
go get github.com/davecgh/go-spew/spew
```