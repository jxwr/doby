package comp

import (
	"fmt"
	"reflect"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/env"
	"github.com/jxwr/doubi/rt"
	"github.com/jxwr/doubi/token"
)

type Attr struct {
	Debug bool
	E     *env.Env
	Fun   *ast.FuncDeclExpr
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
		_, builtin := rt.Builtins[arg.Name]

		keyword := false
		for i := token.BREAK; i <= token.VAR; i++ {
			if token.Tokens[i] == arg.Name {
				keyword = true
			}
		}
		_, env := self.E.LookUp(arg.Name)
		if arg.Name != "_" && env == nil && !builtin && !keyword {
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

	self.checkIdentRef(node.X)
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
	if node.Low != nil {
		self.checkIdentRef(node.Low)
	}
	if node.High != nil {
		self.checkIdentRef(node.High)
	}
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

	if node.Name != nil {
		self.E.Put(node.Name.Name, node.Name)
	}

	self.Enter()
	for _, arg := range node.Args {
		self.E.Put(arg.Name, arg)
	}

	fnBak := self.Fun
	self.Fun = node
	node.Body.Accept(self)
	self.Fun = fnBak
	self.Leave()
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
	} else {
		for _, arg := range node.Lhs {
			if ident, ok := arg.(*ast.Ident); ok {
				val, _ := self.E.LookUp(ident.Name)
				// set lexical variable
				if val == nil && self.Fun != nil {
					self.Fun.LocalNames = append(self.Fun.LocalNames, ident.Name)
				}
				self.E.Put(ident.Name, ident)
			}
		}
	}

	for _, arg := range node.Rhs {
		self.checkIdentRef(arg)
	}

	for _, lh := range node.Lhs {
		lh, ok := lh.(*ast.Ident)
		if ok {
			self.E.Put(lh.Name, lh)
		}
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

	if node.Init != nil {
		node.Init.Accept(self)
	}
	self.checkIdentRef(node.Cond)
	node.Body.Accept(self)
}

func (self *Attr) VisitRangeStmt(node *ast.RangeStmt) {
	self.debug(node)

	key := node.KeyValue[0].(*ast.Ident)
	self.E.Put(key.Name, key)

	val := node.KeyValue[1].(*ast.Ident)
	self.E.Put(val.Name, val)

	self.checkIdentRef(node.X)
	self.Enter()
	node.Body.Accept(self)
	self.Leave()
}

func (self *Attr) VisitImportStmt(node *ast.ImportStmt) {
}

func (self *Attr) Enter() {
	self.E = env.NewEnv(self.E)
}

func (self *Attr) Leave() {
	self.E = self.E.Outer
}
