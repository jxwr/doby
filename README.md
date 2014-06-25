

name 

123
0x1213

(expr)

expr . name

expr [ index ]

slice : expr [ low : high ]

expr ( []expr )

op expr

expr op expr

-- 
ident : name
      
basiclit : INT | FLOAT | STRING | CHAR 

paren_expr : '(' expr ')'

selector_expr : expr '.' ident

index_expr : expr '[' expr ']'

slice_expr : expr '[' expr ':' expr ':' expr ']'

call_expr : expr '(' expr_list ')'

unary_expr : OP expr 

binary_expr : expr OP expr

expr : ident
     | basiclit
     | paren_expr
     | selector_expr
     | inden_expr
     | slice_expr
     | call_expr
     | unary_expr
     | binary_expr
     ;

