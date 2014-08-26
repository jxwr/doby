package comp

import (
	"fmt"
	"strings"

	"github.com/jxwr/doby/ast"
	"github.com/jxwr/doby/env"
	"github.com/jxwr/doby/token"
)

type Attr struct {
	env *env.Env
	fun *ast.FuncDeclExpr
}

func NewAttr() *Attr {
	attr := &Attr{env: env.NewEnv(nil)}
	return attr
}

func (self *Attr) log(fmtstr string, args ...interface{}) {
	fmt.Printf(fmtstr, args...)
	fmt.Println()
}

func (self *Attr) checkIdentRef(node ast.Expr) {
	switch arg := node.(type) {
	case *ast.Ident:
		keyword := false
		for i := token.BREAK; i <= token.VAR; i++ {
			if token.Tokens[i] == arg.Name {
				keyword = true
			}
		}
		_, env := self.env.LookUp(arg.Name)
		if arg.Name != "_" && env == nil && !keyword {
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
}

func (self *Attr) VisitBasicLit(node *ast.BasicLit) {
}

func (self *Attr) VisitParenExpr(node *ast.ParenExpr) {
	self.checkIdentRef(node.X)
}

func (self *Attr) VisitSelectorExpr(node *ast.SelectorExpr) {
	self.checkIdentRef(node.X)
}

func (self *Attr) VisitIndexExpr(node *ast.IndexExpr) {
	self.checkIdentRef(node.X)
	self.checkIdentRef(node.Index)
}

func (self *Attr) VisitSliceExpr(node *ast.SliceExpr) {
	self.checkIdentRef(node.X)
	if node.Low != nil {
		self.checkIdentRef(node.Low)
	}
	if node.High != nil {
		self.checkIdentRef(node.High)
	}
}

func (self *Attr) VisitCallExpr(node *ast.CallExpr) {
	self.checkIdentRef(node.Fun)
	self.checkIdentListRef(node.Args)
}

func (self *Attr) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.checkIdentRef(node.X)
}

func (self *Attr) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.checkIdentRef(node.X)
	self.checkIdentRef(node.Y)
}

func (self *Attr) VisitArrayExpr(node *ast.ArrayExpr) {
	self.checkIdentListRef(node.Elems)
}

func (self *Attr) VisitSetExpr(node *ast.SetExpr) {
	self.checkIdentListRef(node.Elems)
}

func (self *Attr) VisitDictExpr(node *ast.DictExpr) {
	for _, field := range node.Fields {
		self.checkIdentRef(field.Name)
		self.checkIdentRef(field.Value)
	}
}

func (self *Attr) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	if node.Name != nil {
		self.env.Put(node.Name.Name, node.Name)
	}

	self.Enter()
	for _, arg := range node.Args {
		self.env.Put(arg.Name, arg)
	}

	fnBak := self.fun
	self.fun = node
	node.Body.Accept(self)
	self.fun = fnBak
	self.Leave()
}

// stmts

func (self *Attr) VisitExprStmt(node *ast.ExprStmt) {
	node.X.Accept(self)
}

func (self *Attr) VisitSendStmt(node *ast.SendStmt) {
	self.checkIdentRef(node.Chan)
	self.checkIdentRef(node.Value)
}

func (self *Attr) VisitIncDecStmt(node *ast.IncDecStmt) {
	self.checkIdentRef(node.X)
}

func (self *Attr) VisitAssignStmt(node *ast.AssignStmt) {
	if node.Tok != token.ASSIGN {
		for _, arg := range node.Lhs {
			self.checkIdentRef(arg)
		}
	} else {
		for _, arg := range node.Lhs {
			if ident, ok := arg.(*ast.Ident); ok {
				val, _ := self.env.LookUp(ident.Name)
				// set lexical variable
				if val == nil && self.fun != nil {
					self.fun.LocalNames = append(self.fun.LocalNames, ident.Name)
				}
				self.env.Put(ident.Name, ident)
			}
		}
	}

	for _, arg := range node.Rhs {
		self.checkIdentRef(arg)
	}

	for _, lh := range node.Lhs {
		lh, ok := lh.(*ast.Ident)
		if ok {
			self.env.Put(lh.Name, lh)
		}
	}
}

func (self *Attr) VisitGoStmt(node *ast.GoStmt) {
	node.Call.Accept(self)
}

func (self *Attr) VisitReturnStmt(node *ast.ReturnStmt) {
	self.checkIdentListRef(node.Results)
}

func (self *Attr) VisitBranchStmt(node *ast.BranchStmt) {
}

func (self *Attr) VisitBlockStmt(node *ast.BlockStmt) {
	for _, stmt := range node.List {
		stmt.Accept(self)
	}
}

func (self *Attr) VisitIfStmt(node *ast.IfStmt) {
	self.checkIdentRef(node.Cond)
	node.Body.Accept(self)
	if node.Else != nil {
		node.Else.Accept(self)
	}
}

func (self *Attr) VisitCaseClause(node *ast.CaseClause) {
	self.checkIdentListRef(node.List)
	for _, stmt := range node.Body {
		stmt.Accept(self)
	}
}

func (self *Attr) VisitSwitchStmt(node *ast.SwitchStmt) {
	node.Init.Accept(self)
	node.Body.Accept(self)
}

func (self *Attr) VisitSelectStmt(node *ast.SelectStmt) {
	node.Body.Accept(self)
}

func (self *Attr) VisitForStmt(node *ast.ForStmt) {
	if node.Init != nil {
		node.Init.Accept(self)
	}
	self.checkIdentRef(node.Cond)
	node.Body.Accept(self)
}

func (self *Attr) VisitRangeStmt(node *ast.RangeStmt) {
	key := node.KeyValue[0].(*ast.Ident)
	self.env.Put(key.Name, key)

	val := node.KeyValue[1].(*ast.Ident)
	self.env.Put(val.Name, val)

	self.checkIdentRef(node.X)
	self.Enter()
	node.Body.Accept(self)
	self.Leave()
}

func (self *Attr) VisitImportStmt(node *ast.ImportStmt) {
	if len(node.Modules) == 1 {
		modname := node.Modules[0]
		modname = strings.Trim(modname, "\"")
		xs := strings.Split(modname, "/")
		self.env.Put(xs[len(xs)-1], node)
	}
}

func (self *Attr) Enter() {
	self.env = env.NewEnv(self.env)
}

func (self *Attr) Leave() {
	self.env = self.env.Outer
}
