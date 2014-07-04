package comp

import (
	"fmt"
	"reflect"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/token"
)

type Attr struct {
	Debug bool
	Env   map[string]interface{}
}

func (self *Attr) log(fmtstr string, args ...interface{}) {
	fmt.Printf(fmtstr, args...)
	fmt.Println()
}

func (self *Attr) debug(node interface{}) {
	if self.Debug {
		fmt.Printf("%s(%#v)\n", reflect.TypeOf(node).Name(), node)
	}
}

func (self *Attr) checkIdentRef(node ast.Expr) {
	switch arg := node.(type) {
	case *ast.Ident:
		if _, ok := self.Env[arg.Name]; !ok {
			self.log("'%s' not found", arg.Name)
		}
	default:
		arg.Accept(self)
	}
}

func (self *Attr) checkIdentListRef(nodes []ast.Expr) {
	for _, node := range nodes {
		self.checkIdentRef(node)
	}
}

// exprs

func (self *Attr) VisitIdent(node *ast.Ident) {
	self.debug(node)
}

func (self *Attr) VisitBasicLit(node *ast.BasicLit) {
	self.debug(node)
}

func (self *Attr) VisitParenExpr(node *ast.ParenExpr) {
	self.debug(node)
}

func (self *Attr) VisitSelectorExpr(node *ast.SelectorExpr) {
	self.debug(node)

	self.checkIdentRef(node.X)
}

func (self *Attr) VisitIndexExpr(node *ast.IndexExpr) {
	self.debug(node)

	self.checkIdentRef(node.X)
	self.checkIdentRef(node.Index)
}

func (self *Attr) VisitSliceExpr(node *ast.SliceExpr) {
	self.debug(node)

	self.checkIdentRef(node.X)
	self.checkIdentRef(node.Low)
	self.checkIdentRef(node.High)
}

func (self *Attr) VisitCallExpr(node *ast.CallExpr) {
	self.debug(node)

	self.checkIdentRef(node.Fun)
	self.checkIdentListRef(node.Args)
}

func (self *Attr) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.debug(node)

	self.checkIdentRef(node.X)
}

func (self *Attr) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.debug(node)

	self.checkIdentRef(node.X)
	self.checkIdentRef(node.Y)
}

func (self *Attr) VisitArrayExpr(node *ast.ArrayExpr) {
	self.debug(node)

	self.checkIdentListRef(node.Elems)
}

func (self *Attr) VisitSetExpr(node *ast.SetExpr) {
	self.debug(node)

	self.checkIdentListRef(node.Elems)
}

func (self *Attr) VisitDictExpr(node *ast.DictExpr) {
	self.debug(node)

	for _, field := range node.Fields {
		self.checkIdentRef(field.Name)
		self.checkIdentRef(field.Value)
	}
}

func (self *Attr) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	self.debug(node)

	node.Body.Accept(self)
}

// stmts

func (self *Attr) VisitExprStmt(node *ast.ExprStmt) {
	self.debug(node)

	node.X.Accept(self)
}

func (self *Attr) VisitSendStmt(node *ast.SendStmt) {
	self.debug(node)

	self.checkIdentRef(node.Chan)
	self.checkIdentRef(node.Value)
}

func (self *Attr) VisitIncDecStmt(node *ast.IncDecStmt) {
	self.debug(node)

	self.checkIdentRef(node.X)
}

func (self *Attr) VisitAssignStmt(node *ast.AssignStmt) {
	self.debug(node)

	if node.Tok != token.ASSIGN {
		for _, arg := range node.Lhs {
			self.checkIdentRef(arg)
		}
	}

	for _, arg := range node.Rhs {
		self.checkIdentRef(arg)
	}
}

func (self *Attr) VisitGoStmt(node *ast.GoStmt) {
	self.debug(node)

	node.Call.Accept(self)
}

func (self *Attr) VisitReturnStmt(node *ast.ReturnStmt) {
	self.debug(node)

	self.checkIdentListRef(node.Results)
}

func (self *Attr) VisitBranchStmt(node *ast.BranchStmt) {
	self.debug(node)
}

func (self *Attr) VisitBlockStmt(node *ast.BlockStmt) {
	self.debug(node)

	for _, stmt := range node.List {
		stmt.Accept(self)
	}
}

func (self *Attr) VisitIfStmt(node *ast.IfStmt) {
	self.debug(node)

	self.checkIdentRef(node.Cond)
	node.Body.Accept(self)
	if node.Else != nil {
		node.Else.Accept(self)
	}
}

func (self *Attr) VisitCaseClause(node *ast.CaseClause) {
	self.debug(node)

	self.checkIdentListRef(node.List)
	for _, stmt := range node.Body {
		stmt.Accept(self)
	}
}

func (self *Attr) VisitSwitchStmt(node *ast.SwitchStmt) {
	self.debug(node)

	node.Init.Accept(self)
	node.Body.Accept(self)
}

func (self *Attr) VisitSelectStmt(node *ast.SelectStmt) {
	self.debug(node)

	node.Body.Accept(self)
}

func (self *Attr) VisitForStmt(node *ast.ForStmt) {
	self.debug(node)

	node.Init.Accept(self)
	self.checkIdentRef(node.Cond)
	node.Body.Accept(self)
}

func (self *Attr) VisitRangeStmt(node *ast.RangeStmt) {
	self.debug(node)

	self.checkIdentListRef(node.KeyValue)
	self.checkIdentRef(node.X)
	node.Body.Accept(self)
}
