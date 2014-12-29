calc
====
My toy calculator language

TODO:
- Split number into NUM, FLOAT, HEX, EXPONENT in Lexer (or parser?)

Identifiers can be alphanumeric with an underscore '_'

Grammar in BNF:

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