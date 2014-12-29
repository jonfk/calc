calc
====
My toy calculator language. The lexer was inspired by the
[text/template/parse package](http://golang.org/pkg/text/template/parse/)
from the standard library and this talk by Rob Pike:
[Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE)

###TODO:
- Split number into NUM, FLOAT, HEX, EXPONENT in Lexer (or parser?)

###Notes:
- Identifiers can be alphanumeric with an underscore '_'

##Grammar in BNF:

    expr ::= num_expr
           | bool_expr
           | if_expr

    if_expr ::= IF bool_expr THEN expr ELSE expr END

    num_expr ::= NUMBER
               | num_expr ADD num_expr
               | num_expr SUB num_expr
               | num_expr MUL num_expr
               | num_expr QUO num_expr
               | num_expr REM num_expr

    bool_expr ::= BOOL
                | NOT bool_expr
                | bool_expr LAND bool_expr
                | bool_expr LOR bool_expr
                | num_expr EQL num_expr
                | num_expr LSS num_expr
                | num_expr GTR num_expr
                | num_expr NEQ num_expr
                | num_expr LEQ num_expr
                | num_expr GEQ num_expr

    let_expr ::= LET IDENTIFIER ASSIGN expr END
