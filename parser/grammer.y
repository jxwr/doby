%{

package parser

import (
       "github.com/jxwr/gocalc/ast"
       "github.com/jxwr/gocalc/token"
)

%}

// fields inside this union end up as the fields in a structure known
// as ${PREFIX}SymType, of which a reference is passed to the lexer.
%union {
    node ast.Node
    expr ast.Expr
    arg_list []ast.Expr
    stmt ast.Stmt
    lit string
}

%type <expr> expr ident basiclit
%type <expr> paren_expr selector_expr index_expr slice_expr 
%type <expr> call_expr unary_expr binary_expr prog
%type <arg_list> arg_list

// same for terminals
%token <lit> IDENT INT FLOAT STRING CHAR

%left '['
%left '.'
%left '|'
%left '&'
%left '+'  '-'
%left '*'  '/'  '%'
%left '('
%left UMINUS

%start prog

%%

ident : IDENT
      	{ $$ = ast.Ident{0, $1} }
      ;

basiclit : INT
	   { $$ = ast.BasicLit{0, token.INT, $1} }
	 | FLOAT 
	   { $$ = ast.BasicLit{0, token.FLOAT, $1} }
	 | STRING 
	   { $$ = ast.BasicLit{0, token.STRING, $1} }
	 | CHAR 
	   { $$ = ast.BasicLit{0, token.CHAR, $1} }
	 ;

paren_expr : '(' expr ')'
	     { $$ = ast.ParenExpr{0, $2, 0} }
	   ;

selector_expr : expr '.' ident
	      	{ $$ = ast.SelectorExpr{$1, $3.(ast.Ident)} }
	      ;

slice_expr : expr '[' expr ':' expr ']'
	     { $$ = ast.SliceExpr{$1, 0, $3, $5, 0} }	
	   ;

index_expr : expr '[' expr ']'
	     { $$ = ast.IndexExpr{$1, 0, $3, 0} }
	   ;

arg_list : expr
	   { $$ = []ast.Expr{$1} }
	 | arg_list ',' expr
	   { $$ = append($1, $3) }
	 ;

call_expr : expr '(' arg_list ')'
	    { $$ = ast.CallExpr{$1, 0, $3, 0} }
	  ;  

unary_expr : '-' expr %prec UMINUS
	     { $$ = ast.UnaryExpr{0, token.SUB, $2 } }
	   ;

binary_expr : expr '+' expr
	      { $$ = ast.BinaryExpr{$1, 0, token.ADD, $3 } }
            | expr '-' expr
	      { $$ = ast.BinaryExpr{$1, 0, token.SUB, $3 } }
            | expr '*' expr
	      { $$ = ast.BinaryExpr{$1, 0, token.MUL, $3 } }
            | expr '/' expr
	      { $$ = ast.BinaryExpr{$1, 0, token.QUO, $3 } }
            | expr '%' expr
	      { $$ = ast.BinaryExpr{$1, 0, token.REM, $3 } }
            | expr '&' expr
	      { $$ = ast.BinaryExpr{$1, 0, token.AND, $3 } }
            | expr '|' expr
	      { $$ = ast.BinaryExpr{$1, 0, token.OR, $3 } }
	    ;

expr : ident
     | basiclit
     | paren_expr
     | selector_expr
     | index_expr
     | slice_expr
     | call_expr
     | unary_expr
     | binary_expr
     ;

prog : expr
       { __yyfmt__.Printf("%#v\n", $1) }
     ;
