package eval

import (
	"fmt"

	"github.com/jxwr/doubi/ast"
)

func EvalExprStmt(stmt *ast.ExprStmt) {
	fmt.Printf("EXPR %#v\n", *stmt)
}

func EvalAssignStmt(stmt *ast.AssignStmt) {
	fmt.Printf("ASSIGN %#v\n", *stmt)
}

func EvalStmt(stmt *ast.Stmt) {
	//	fmt.Printf("%#v\n", *stmt)

	switch s := (*stmt).(type) {
	case ast.ExprStmt:
		EvalExprStmt(&s)
	case ast.AssignStmt:
		EvalAssignStmt(&s)
	}
}

func Eval(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		EvalStmt(&stmt)
	}
}
