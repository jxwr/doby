package comp

import (
	"fmt"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/token"
)

type PrettyPrinter struct {
	Debug       bool
	Indent      int
	ShowNewLine bool
}

func puts(str string) {
	fmt.Print(str)
}

func putTok(tok token.Token) {
	puts(" ")
	puts(token.Tokens[tok])
	puts(" ")
}

func (self *PrettyPrinter) putln() {
	if self.ShowNewLine {
		fmt.Println()
	}
}

func (self *PrettyPrinter) putIndent() {
	for i := 0; i < self.Indent; i++ {
		puts("  ")
	}
}

func (self *PrettyPrinter) debug(node interface{}) {
	if self.Debug {
		fmt.Printf("%#v\n", node)
	}
}

func (self *PrettyPrinter) VisitIdent(node *ast.Ident) {
	self.debug(node)

	puts(node.Name)
}

func (self *PrettyPrinter) VisitBasicLit(node *ast.BasicLit) {
	self.debug(node)

	puts(node.Value)
}

func (self *PrettyPrinter) VisitParenExpr(node *ast.ParenExpr) {
	self.debug(node)

	puts("(")
	node.X.Accept(self)
	puts(")")
}

func (self *PrettyPrinter) VisitSelectorExpr(node *ast.SelectorExpr) {
	self.debug(node)

	node.X.Accept(self)
	puts(".")
	node.Sel.Accept(self)
}

func (self *PrettyPrinter) VisitIndexExpr(node *ast.IndexExpr) {
	self.debug(node)

	node.X.Accept(self)
	puts("[")
	node.Index.Accept(self)
	puts("]")
}

func (self *PrettyPrinter) VisitSliceExpr(node *ast.SliceExpr) {
	self.debug(node)

	node.X.Accept(self)
	puts("[")
	node.Low.Accept(self)
	puts(":")
	node.High.Accept(self)
	puts("]")
}

func (self *PrettyPrinter) VisitCallExpr(node *ast.CallExpr) {
	self.debug(node)

	node.Fun.Accept(self)
	puts("(")
	for i, arg := range node.Args {
		self.ShowNewLine = false
		arg.Accept(self)
		self.ShowNewLine = true
		if i < len(node.Args)-1 {
			puts(", ")
		}
	}
	puts(")")

}

func (self *PrettyPrinter) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.debug(node)

	putTok(node.Op)
	node.X.Accept(self)
}

func (self *PrettyPrinter) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.debug(node)

	node.X.Accept(self)
	putTok(node.Op)
	node.Y.Accept(self)
}

func (self *PrettyPrinter) VisitArrayExpr(node *ast.ArrayExpr) {
	self.debug(node)

	puts("[")
	for i, elem := range node.Elems {
		elem.Accept(self)
		if i < len(node.Elems)-1 {
			puts(", ")
		}
	}
	puts("]")
}

func (self *PrettyPrinter) VisitSetExpr(node *ast.SetExpr) {
	self.debug(node)

	puts("#[")
	for i, elem := range node.Elems {
		elem.Accept(self)
		if i < len(node.Elems)-1 {
			puts(", ")
		}
	}
	puts("]")
}

func (self *PrettyPrinter) VisitDictExpr(node *ast.DictExpr) {
	self.debug(node)
	puts("#{")
	for i, field := range node.Fields {
		field.Name.Accept(self)
		puts(":")
		field.Value.Accept(self)

		if i < len(node.Fields)-1 {
			puts(", ")
		}
	}
	puts("}")
}

func (self *PrettyPrinter) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	self.debug(node)

	puts("func ")
	if node.Recv != nil {
		puts("(")
		node.Recv.Accept(self)
		puts(" ")
		node.RecvType.Accept(self)
		puts(") ")
	}

	if node.Name != nil {
		node.Name.Accept(self)
	}

	puts("(")
	for i, arg := range node.Args {
		arg.Accept(self)
		if i < len(node.Args)-1 {
			puts(", ")
		}
	}
	puts(") ")
	node.Body.Accept(self)
}

func (self *PrettyPrinter) VisitExprStmt(node *ast.ExprStmt) {
	self.debug(node)

	node.X.Accept(self)
	self.putln()
}

func (self *PrettyPrinter) VisitSendStmt(node *ast.SendStmt) {
	self.debug(node)

	node.Chan.Accept(self)
	puts(" <- ")
	node.Value.Accept(self)
	self.putln()
}

func (self *PrettyPrinter) VisitIncDecStmt(node *ast.IncDecStmt) {
	self.debug(node)

	node.X.Accept(self)
	putTok(node.Tok)
	self.putln()
}

func (self *PrettyPrinter) VisitAssignStmt(node *ast.AssignStmt) {
	self.debug(node)

	for i, arg := range node.Lhs {
		arg.Accept(self)
		if i < len(node.Lhs)-1 {
			puts(", ")
		}
	}

	putTok(node.Tok)

	for i, arg := range node.Rhs {
		arg.Accept(self)
		if i < len(node.Rhs)-1 {
			puts(", ")
		}
	}
	self.putln()
}

func (self *PrettyPrinter) VisitGoStmt(node *ast.GoStmt) {
	self.debug(node)

	puts("go ")
	node.Call.Accept(self)
	self.putln()
}

func (self *PrettyPrinter) VisitReturnStmt(node *ast.ReturnStmt) {
	self.debug(node)

	puts("return ")
	for i, ret := range node.Results {
		ret.Accept(self)
		if i < len(node.Results)-1 {
			puts(", ")
		}
	}
	self.putln()
}

func (self *PrettyPrinter) VisitBranchStmt(node *ast.BranchStmt) {
	self.debug(node)

	putTok(node.Tok)
	self.putln()
}

func (self *PrettyPrinter) VisitBlockStmt(node *ast.BlockStmt) {
	self.debug(node)

	nl := self.ShowNewLine
	puts("{")
	self.ShowNewLine = true
	self.putln()
	self.Indent++
	for _, stmt := range node.List {
		self.putIndent()
		stmt.Accept(self)
	}
	self.Indent--
	self.ShowNewLine = nl
	puts("}")
	self.putln()
}

func (self *PrettyPrinter) VisitIfStmt(node *ast.IfStmt) {
	self.debug(node)

	puts("if ")
	node.Cond.Accept(self)
	puts(" ")
	node.Body.Accept(self)

	if node.Else != nil {
		puts(" else ")
		node.Else.Accept(self)
	}
	self.putln()
}

func (self *PrettyPrinter) VisitCaseClause(node *ast.CaseClause) {
	self.debug(node)

	puts("case ")
	for i, stmt := range node.List {
		stmt.Accept(self)
		if i < len(node.List)-1 {
			puts(", ")
		}
	}
	puts(":")
	self.putln()
	self.Indent++
	for _, stmt := range node.Body {
		self.putIndent()
		stmt.Accept(self)
	}
	self.Indent--
}

func (self *PrettyPrinter) VisitSwitchStmt(node *ast.SwitchStmt) {
	self.debug(node)

	puts("switch ")
	self.ShowNewLine = false
	node.Init.Accept(self)
	self.ShowNewLine = true
	puts(" ")
	node.Body.Accept(self)
}

func (self *PrettyPrinter) VisitSelectStmt(node *ast.SelectStmt) {
	self.debug(node)

	puts("select ")
	node.Body.Accept(self)
}

func (self *PrettyPrinter) VisitForStmt(node *ast.ForStmt) {
	self.debug(node)

	puts("for ")
	self.ShowNewLine = false
	node.Init.Accept(self)
	puts("; ")
	node.Cond.Accept(self)
	puts("; ")
	node.Post.Accept(self)
	puts(" ")
	self.ShowNewLine = true
	node.Body.Accept(self)
}

func (self *PrettyPrinter) VisitRangeStmt(node *ast.RangeStmt) {
	self.debug(node)

	puts("for ")
	self.ShowNewLine = false
	node.KeyValue[0].Accept(self)
	puts(", ")
	node.KeyValue[1].Accept(self)
	puts(" = range ")
	node.X.Accept(self)
	puts(" ")
	self.ShowNewLine = true
	node.Body.Accept(self)
}
