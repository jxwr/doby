package comp

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/token"
)

type Stack struct {
	cur  int
	vals []Object
}

func NewStack() *Stack {
	stack := &Stack{0, []Object{}}
	return stack
}

func (self *Stack) Push(obj Object) {
	if len(self.vals) <= self.cur {
		self.vals = append(self.vals, obj)
	} else {
		self.vals[self.cur] = obj
	}
	self.cur++
}

func (self *Stack) Pop() Object {
	if self.cur == 0 {
		panic("pop from empty stack")
	}
	self.cur--
	return self.vals[self.cur]
}

type Eval struct {
	Debug bool
	E     *Env
	Stack *Stack
}

func (self *Eval) log(fmtstr string, args ...interface{}) {
	fmt.Printf(fmtstr, args...)
	fmt.Println()
}

func (self *Eval) fatal(fmtstr string, args ...interface{}) {
	fmt.Printf(fmtstr, args...)
	fmt.Println()
	os.Exit(1)
}

func (self *Eval) evalExpr(expr ast.Expr) {
	expr.Accept(self)
}

func (self *Eval) debug(node interface{}) {
	if self.Debug {
		fmt.Printf("%s(%#v)\n", reflect.TypeOf(node).Name(), node)
	}
}

// exprs

func (self *Eval) VisitIdent(node *ast.Ident) {
	self.debug(node)

	obj := self.E.LookUp(node.Name)
	if obj != nil {
		self.Stack.Push(obj.(Object))
	} else {
		panic(node.Name + " not found")
	}
}

func (self *Eval) VisitBasicLit(node *ast.BasicLit) {
	self.debug(node)

	switch node.Kind {
	case token.INT:
		val, err := strconv.Atoi(node.Value)
		if err != nil {
			self.fatal("%s convert to int failed: %v", node.Value, err)
		}
		obj := NewIntegerObject(val)
		self.Stack.Push(obj)
	case token.FLOAT:
		val, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			self.fatal("%s convert to float failed: %v", node.Value, err)
		}
		obj := NewFloatObject(val)
		self.Stack.Push(obj)
	case token.STRING:
		val := strings.Trim(node.Value, "\"")
		obj := NewStringObject(val)
		self.Stack.Push(obj)
	case token.CHAR:
		val := strings.Trim(node.Value, "'")
		obj := NewStringObject(val)
		self.Stack.Push(obj)
	}
}

func (self *Eval) VisitParenExpr(node *ast.ParenExpr) {
	self.debug(node)
}

func (self *Eval) VisitSelectorExpr(node *ast.SelectorExpr) {
	self.debug(node)
}

func (self *Eval) VisitIndexExpr(node *ast.IndexExpr) {
	self.debug(node)
}

func (self *Eval) VisitSliceExpr(node *ast.SliceExpr) {
	self.debug(node)
}

func (self *Eval) VisitCallExpr(node *ast.CallExpr) {
	self.debug(node)

	var fnobj Object
	ident, ok := node.Fun.(*ast.Ident)
	if ok && self.E.LookUp(ident.Name) == nil {
		_, exist := Builtins[ident.Name]
		if exist {
			fnobj = NewFuncObject(ident.Name, nil)
		}
	} else {
		self.evalExpr(node.Fun)
		fnobj = self.Stack.Pop()
	}

	args := []Object{}
	for _, arg := range node.Args {
		self.evalExpr(arg)
		args = append(args, self.Stack.Pop())
	}

	fnobj.Dispatch("__call__", args...)
}

func (self *Eval) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.debug(node)
}

func (self *Eval) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.debug(node)

	self.evalExpr(node.Y)
	self.evalExpr(node.X)

	lobj := self.Stack.Pop()
	robj := self.Stack.Pop()

	switch node.Op {
	case token.ADD:
		objs := lobj.Dispatch("__add__", robj)
		self.Stack.Push(objs[0])
	case token.SUB:
		objs := lobj.Dispatch("__sub__", robj)
		self.Stack.Push(objs[0])
	case token.MUL:
		objs := lobj.Dispatch("__mul__", robj)
		self.Stack.Push(objs[0])
	case token.QUO:
		objs := lobj.Dispatch("__QUO__", robj)
		self.Stack.Push(objs[0])
	case token.REM:
		objs := lobj.Dispatch("__REM__", robj)
		self.Stack.Push(objs[0])
	}
}

func (self *Eval) VisitArrayExpr(node *ast.ArrayExpr) {
	self.debug(node)

	elems := []Object{}
	for _, elem := range node.Elems {
		self.evalExpr(elem)
		elems = append(elems, self.Stack.Pop())
	}
	obj := NewArrayObject(elems)
	self.Stack.Push(obj)
}

func (self *Eval) VisitSetExpr(node *ast.SetExpr) {
	self.debug(node)

	elems := []Object{}
	for _, elem := range node.Elems {
		self.evalExpr(elem)
		elems = append(elems, self.Stack.Pop())
	}
	obj := NewSetObject(elems)
	self.Stack.Push(obj)
}

func (self *Eval) VisitDictExpr(node *ast.DictExpr) {
	self.debug(node)
}

func (self *Eval) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	self.debug(node)
}

// stmts

func (self *Eval) VisitExprStmt(node *ast.ExprStmt) {
	self.debug(node)

	node.X.Accept(self)
}

func (self *Eval) VisitSendStmt(node *ast.SendStmt) {
	self.debug(node)
}

func (self *Eval) VisitIncDecStmt(node *ast.IncDecStmt) {
	self.debug(node)
}

func (self *Eval) VisitAssignStmt(node *ast.AssignStmt) {
	self.debug(node)

	for i := 0; i < len(node.Lhs); i++ {
		self.evalExpr(node.Rhs[i])
		robj := self.Stack.Pop()

		switch v := node.Lhs[i].(type) {
		case *ast.Ident:
			self.E.Put(v.Name, robj)
		case *ast.IndexExpr:
			self.evalExpr(v.X)
			lobj := self.Stack.Pop()
			self.evalExpr(v.Index)
			idx := self.Stack.Pop()
			lobj.Dispatch("__set_index__", idx, robj)
		case *ast.SelectorExpr:
			self.evalExpr(v.X)
			lobj := self.Stack.Pop()
			sel := NewStringObject(v.Sel.Name)
			lobj.Dispatch("__set_property__", sel, robj)
		}
	}
}

func (self *Eval) VisitGoStmt(node *ast.GoStmt) {
	self.debug(node)
}

func (self *Eval) VisitReturnStmt(node *ast.ReturnStmt) {
	self.debug(node)
}

func (self *Eval) VisitBranchStmt(node *ast.BranchStmt) {
	self.debug(node)
}

func (self *Eval) VisitBlockStmt(node *ast.BlockStmt) {
	self.debug(node)
}

func (self *Eval) VisitIfStmt(node *ast.IfStmt) {
	self.debug(node)
}

func (self *Eval) VisitCaseClause(node *ast.CaseClause) {
	self.debug(node)
}

func (self *Eval) VisitSwitchStmt(node *ast.SwitchStmt) {
	self.debug(node)
}

func (self *Eval) VisitSelectStmt(node *ast.SelectStmt) {
	self.debug(node)
}

func (self *Eval) VisitForStmt(node *ast.ForStmt) {
	self.debug(node)
}

func (self *Eval) VisitRangeStmt(node *ast.RangeStmt) {
	self.debug(node)
}
