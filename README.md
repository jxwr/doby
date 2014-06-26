DOUBI GRAMMER

;;; expression 

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

;;; Stmts

expr_stmt : expr

send_stmt : expr "<-" expr

incdec_stmt : expr "++"
	    | expr "--"

assign_stmt : expr_list = expr_list

go_stmt : go call_expr

return_stmt : "return" expr_list

branch_stmt : "break" "continue"

block_stmt : '{' stmt_list '}'

if_stmt : "if" expr block_stmt 'else' stmt

case_stmt : "case" expr_list ':' stmt_list

switch_stmt : "switch" stmt block_stmt

select_stmt : "select" block_stmt

for_stmt : "for" stmt ';' expr ';' stmt block_stmt

range_stmt : "for" expr ',' expr '=' expr block_stmt
