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

	NeedReturn bool
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

	if node.Name == "true" {
		obj := NewBoolObject(true)
		self.Stack.Push(obj)
	} else if node.Name == "false" {
		obj := NewBoolObject(false)
		self.Stack.Push(obj)
	} else {
		obj := self.E.LookUp(node.Name)
		if obj != nil {
			self.Stack.Push(obj.(Object))
		} else {
			panic(node.Name + " not found")
		}
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

	node.X.Accept(self)
}

func (self *Eval) VisitSelectorExpr(node *ast.SelectorExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop()
	prop := NewStringObject(node.Sel.Name)
	rets := obj.Dispatch(self, "__get_property__", prop)
	self.Stack.Push(rets[0])
}

func (self *Eval) VisitIndexExpr(node *ast.IndexExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop()
	self.evalExpr(node.Index)
	index := self.Stack.Pop()
	rets := obj.Dispatch(self, "__get_index__", index)
	self.Stack.Push(rets[0])
}

func (self *Eval) VisitSliceExpr(node *ast.SliceExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop()

	var lowObj Object
	var highObj Object

	if node.Low != nil {
		self.evalExpr(node.Low)
		lowObj = self.Stack.Pop()
	}
	if node.High != nil {
		self.evalExpr(node.High)
		highObj = self.Stack.Pop()
	}

	rets := obj.Dispatch(self, "__slice__", lowObj, highObj)
	self.Stack.Push(rets[0])
}

func (self *Eval) VisitCallExpr(node *ast.CallExpr) {
	var fnobj Object
	ident, ok := node.Fun.(*ast.Ident)

	if ok && self.E.LookUp(ident.Name) == nil {
		_, exist := Builtins[ident.Name]
		if exist {
			fnobj = NewFuncObject(ident.Name, nil)
			args := []Object{}
			for _, arg := range node.Args {
				self.evalExpr(arg)
				args = append(args, self.Stack.Pop())
			}
			self.E = NewEnv(self.E)
			rets := fnobj.Dispatch(self, "__call__", args...)
			self.E = self.E.Outer
			for _, ret := range rets {
				self.Stack.Push(ret)
			}
		}
	} else {
		self.evalExpr(node.Fun)
		fnobj = self.Stack.Pop()

		fn := fnobj.(*FuncObject)
		if fn.IsBuiltin {
			args := []Object{}
			for _, arg := range node.Args {
				self.evalExpr(arg)
				args = append(args, self.Stack.Pop())
			}
			self.E = NewEnv(self.E)
			rets := fnobj.Dispatch(self, "__call__", args...)
			self.E = self.E.Outer
			for _, ret := range rets {
				self.Stack.Push(ret)
			}
		} else {
			fnDecl := fn.Decl

			self.E = NewEnv(self.E)
			for i, arg := range node.Args {
				self.evalExpr(arg)
				self.E.Put(fnDecl.Args[i].Name, self.Stack.Pop())
			}
			self.NeedReturn = false
			fnDecl.Body.Accept(self)
			self.NeedReturn = false
			self.E = self.E.Outer
		}
	}
}

func (self *Eval) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop().(*IntegerObject)
	self.Stack.Push(NewIntegerObject(-obj.val))
}

var BinaryFuncs = map[token.Token]string{
	token.ADD:     "__add__",
	token.SUB:     "__sub__",
	token.MUL:     "__mul__",
	token.QUO:     "__quo__",
	token.REM:     "__rem__",
	token.AND:     "__and__",
	token.OR:      "__or__",
	token.XOR:     "__xor__",
	token.SHL:     "__shl__",
	token.SHR:     "__shr__",
	token.AND_NOT: "__and_not__",
	token.LAND:    "__land__",
	token.LOR:     "__lor__",
	token.EQL:     "__eql__",
	token.LSS:     "__lss__",
	token.GTR:     "__gtr__",
	token.LEQ:     "__leq__",
	token.GEQ:     "__geq__",
	token.NEQ:     "__neq__",
}

func (self *Eval) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	self.evalExpr(node.Y)

	robj := self.Stack.Pop()
	lobj := self.Stack.Pop()

	objs := lobj.Dispatch(self, BinaryFuncs[node.Op], robj)
	self.Stack.Push(objs[0])
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

	fieldMap := map[string]Object{}
	for _, field := range node.Fields {
		self.evalExpr(field.Name)
		key := self.Stack.Pop()
		self.evalExpr(field.Value)
		val := self.Stack.Pop()
		fieldMap[key.HashCode()] = val
	}
	obj := NewDictObject(&fieldMap)
	self.Stack.Push(obj)
}

func (self *Eval) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	self.debug(node)

	if node.Name != nil {
		fname := node.Name.Name
		self.E.Put(fname, NewFuncObject(fname, node))
	} else {
		self.Stack.Push(NewFuncObject("#<closure>", node))
	}
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

	self.evalExpr(node.X)
	obj := self.Stack.Pop()

	if node.Tok == token.INC {
		obj.Dispatch(self, "__inc__")
	} else if node.Tok == token.DEC {
		obj.Dispatch(self, "__dec__")
	}
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
			lobj.Dispatch(self, "__set_index__", idx, robj)
		case *ast.SelectorExpr:
			self.evalExpr(v.X)
			lobj := self.Stack.Pop()
			sel := NewStringObject(v.Sel.Name)
			lobj.Dispatch(self, "__set_property__", sel, robj)
		}
	}
}

func (self *Eval) VisitGoStmt(node *ast.GoStmt) {
	self.debug(node)

	go node.Call.Accept(self)
}

func (self *Eval) VisitReturnStmt(node *ast.ReturnStmt) {
	self.debug(node)

	for _, res := range node.Results {
		self.evalExpr(res)
	}

	self.NeedReturn = true
}

func (self *Eval) VisitBranchStmt(node *ast.BranchStmt) {
	self.debug(node)
}

func (self *Eval) VisitBlockStmt(node *ast.BlockStmt) {
	self.E = NewEnv(self.E)

	for _, stmt := range node.List {
		if self.NeedReturn {
			break
		}
		stmt.Accept(self)
	}
	self.E = self.E.Outer
}

func (self *Eval) VisitIfStmt(node *ast.IfStmt) {
	self.debug(node)

	self.evalExpr(node.Cond)
	cond := self.Stack.Pop()

	if cond.(*BoolObject).val {
		node.Body.Accept(self)
	} else if node.Else != nil {
		node.Else.Accept(self)
	}
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

	if node.Init != nil {
		node.Init.Accept(self)
	}

	for {
		self.evalExpr(node.Cond)
		cond := self.Stack.Pop()
		if !cond.(*BoolObject).val {
			break
		}
		node.Body.Accept(self)
		if node.Post != nil {
			node.Post.Accept(self)
		}
	}
}

func (self *Eval) VisitRangeStmt(node *ast.RangeStmt) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop()

	keyName := node.KeyValue[0].(*ast.Ident).Name
	valName := node.KeyValue[1].(*ast.Ident).Name

	self.E = NewEnv(self.E)

	switch v := obj.(type) {
	case *ArrayObject:
		for i, val := range v.vals {
			self.E.Put(keyName, NewIntegerObject(i))
			self.E.Put(valName, val)
			node.Body.Accept(self)
		}
	case *SetObject:
		for i, val := range v.vals {
			self.E.Put(keyName, NewIntegerObject(i))
			self.E.Put(valName, val)
			node.Body.Accept(self)
		}

	}

	self.E = self.E.Outer
}
