package comp

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/env"
	"github.com/jxwr/doubi/rt"
	"github.com/jxwr/doubi/token"
)

type Stack struct {
	cur  int
	vals []rt.Object
}

func NewStack() *Stack {
	stack := &Stack{0, []rt.Object{}}
	return stack
}

func (self *Stack) Push(obj rt.Object) {
	if len(self.vals) <= self.cur {
		self.vals = append(self.vals, obj)
	} else {
		self.vals[self.cur] = obj
	}
	self.cur++
}

func (self *Stack) Pop() rt.Object {
	if self.cur == 0 {
		panic("pop from empty stack")
	}
	self.cur--
	return self.vals[self.cur]
}

type Eval struct {
	Debug bool
	E     *env.Env
	Stack *Stack
	Fun   *ast.FuncDeclExpr
	RT    *rt.Runtime

	NeedReturn   bool
	LoopDepth    int
	NeedBreak    bool
	NeedContinue bool
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
		obj := rt.NewBoolObject(true)
		self.Stack.Push(obj)
	} else if node.Name == "false" {
		obj := rt.NewBoolObject(false)
		self.Stack.Push(obj)
	} else {
		obj, _ := self.E.LookUp(node.Name)
		if obj != nil {
			self.Stack.Push(obj.(rt.Object))
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
		obj := rt.NewIntegerObject(val)
		self.Stack.Push(obj)
	case token.FLOAT:
		val, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			self.fatal("%s convert to float failed: %v", node.Value, err)
		}
		obj := rt.NewFloatObject(val)
		self.Stack.Push(obj)
	case token.STRING:
		val := strings.Trim(node.Value, "\"")
		obj := rt.NewStringObject(val)
		self.Stack.Push(obj)
	case token.CHAR:
		val := strings.Trim(node.Value, "'")
		obj := rt.NewStringObject(val)
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
	prop := rt.NewStringObject(node.Sel.Name)
	rets := obj.Dispatch(self.RT, "__get_property__", prop)
	self.Stack.Push(rets[0])
}

func (self *Eval) VisitIndexExpr(node *ast.IndexExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop()
	self.evalExpr(node.Index)
	index := self.Stack.Pop()
	rets := obj.Dispatch(self.RT, "__get_index__", index)
	self.Stack.Push(rets[0])
}

func (self *Eval) VisitSliceExpr(node *ast.SliceExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop()

	var lowObj rt.Object
	var highObj rt.Object

	if node.Low != nil {
		self.evalExpr(node.Low)
		lowObj = self.Stack.Pop()
	}
	if node.High != nil {
		self.evalExpr(node.High)
		highObj = self.Stack.Pop()
	}

	rets := obj.Dispatch(self.RT, "__slice__", lowObj, highObj)
	self.Stack.Push(rets[0])
}

func (self *Eval) VisitCallExpr(node *ast.CallExpr) {
	var fnobj *rt.FuncObject
	ident, ok := node.Fun.(*ast.Ident)

	var val interface{}
	if ok {
		val, _ = self.E.LookUp(ident.Name)
	}
	if ok && val == nil {
		_, exist := rt.Builtins[ident.Name]
		if exist {
			fnobj = rt.NewFuncObject(ident.Name, nil, self.E).(*rt.FuncObject)
			args := []rt.Object{}
			for _, arg := range node.Args {
				self.evalExpr(arg)
				args = append(args, self.Stack.Pop())
			}
			self.E = env.NewEnv(self.E)
			fnobj.E = self.E
			rets := fnobj.Dispatch(self.RT, "__call__", args...)
			self.E = self.E.Outer
			for _, ret := range rets {
				self.Stack.Push(ret)
			}
		}
	} else {
		self.evalExpr(node.Fun)
		fnobj = self.Stack.Pop().(*rt.FuncObject)

		if fnobj.IsBuiltin {
			args := []rt.Object{}
			for _, arg := range node.Args {
				self.evalExpr(arg)
				args = append(args, self.Stack.Pop())
			}
			self.E = env.NewEnv(self.E)
			fnobj.E = self.E
			rets := fnobj.Dispatch(self.RT, "__call__", args...)
			self.E = self.E.Outer
			for _, ret := range rets {
				self.Stack.Push(ret)
			}
		} else {
			fnDecl := fnobj.Decl
			fnBak := self.Fun
			self.Fun = fnDecl

			self.E = env.NewEnv(self.E)
			for i, arg := range node.Args {
				self.evalExpr(arg)
				self.E.Put(fnDecl.Args[i].Name, self.Stack.Pop())
			}
			self.NeedReturn = false
			fnobj.E = self.E
			fnDecl.Body.Accept(self)
			self.NeedReturn = false

			self.Fun = fnBak
			self.E = self.E.Outer
		}
	}
}

func (self *Eval) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop().(*rt.IntegerObject)
	self.Stack.Push(rt.NewIntegerObject(-obj.Val))
}

var OpFuncs = map[token.Token]string{
	token.ADD:            "__add__",
	token.SUB:            "__sub__",
	token.MUL:            "__mul__",
	token.QUO:            "__quo__",
	token.REM:            "__rem__",
	token.AND:            "__and__",
	token.OR:             "__or__",
	token.XOR:            "__xor__",
	token.SHL:            "__shl__",
	token.SHR:            "__shr__",
	token.AND_NOT:        "__and_not__",
	token.LAND:           "__land__",
	token.LOR:            "__lor__",
	token.EQL:            "__eql__",
	token.LSS:            "__lss__",
	token.GTR:            "__gtr__",
	token.LEQ:            "__leq__",
	token.GEQ:            "__geq__",
	token.NEQ:            "__neq__",
	token.ADD_ASSIGN:     "__+=__",
	token.SUB_ASSIGN:     "__-=__",
	token.MUL_ASSIGN:     "__*=__",
	token.QUO_ASSIGN:     "__/=__",
	token.REM_ASSIGN:     "__%=__",
	token.AND_ASSIGN:     "__&=__",
	token.OR_ASSIGN:      "__|=__",
	token.XOR_ASSIGN:     "__^=__",
	token.SHL_ASSIGN:     "__<<=__",
	token.SHR_ASSIGN:     "__>>=__",
	token.AND_NOT_ASSIGN: "__&^=__",
}

func (self *Eval) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	self.evalExpr(node.Y)

	robj := self.Stack.Pop()
	lobj := self.Stack.Pop()

	objs := lobj.Dispatch(self.RT, OpFuncs[node.Op], robj)
	self.Stack.Push(objs[0])
}

func (self *Eval) VisitArrayExpr(node *ast.ArrayExpr) {
	self.debug(node)

	elems := []rt.Object{}
	for _, elem := range node.Elems {
		self.evalExpr(elem)
		elems = append(elems, self.Stack.Pop())
	}
	obj := rt.NewArrayObject(elems)
	self.Stack.Push(obj)
}

func (self *Eval) VisitSetExpr(node *ast.SetExpr) {
	self.debug(node)

	elems := []rt.Object{}
	for _, elem := range node.Elems {
		self.evalExpr(elem)
		elems = append(elems, self.Stack.Pop())
	}
	obj := rt.NewSetObject(elems)
	self.Stack.Push(obj)
}

func (self *Eval) VisitDictExpr(node *ast.DictExpr) {
	self.debug(node)

	fieldMap := map[string]rt.Object{}
	for _, field := range node.Fields {
		self.evalExpr(field.Name)
		key := self.Stack.Pop()
		self.evalExpr(field.Value)
		val := self.Stack.Pop()
		fieldMap[key.HashCode()] = val
	}
	obj := rt.NewDictObject(&fieldMap)
	self.Stack.Push(obj)
}

func (self *Eval) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	self.debug(node)

	if node.Name != nil {
		fname := node.Name.Name
		self.E.Put(fname, rt.NewFuncObject(fname, node, self.E))
	} else {
		self.Stack.Push(rt.NewFuncObject("#<closure>", node, self.E))
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
		obj.Dispatch(self.RT, "__inc__")
	} else if node.Tok == token.DEC {
		obj.Dispatch(self.RT, "__dec__")
	}
}

func ContainsString(ss []string, s string) bool {
	found := false
	for _, v := range ss {
		if v == s {
			found = true
			break
		}
	}
	return found
}

func (self *Eval) VisitAssignStmt(node *ast.AssignStmt) {
	self.debug(node)

	if node.Tok == token.ASSIGN {
		for i := 0; i < len(node.Lhs); i++ {
			self.evalExpr(node.Rhs[i])
			robj := self.Stack.Pop()

			switch v := node.Lhs[i].(type) {
			case *ast.Ident:
				// closure
				val, env := self.E.LookUp(v.Name)
				if val == nil {
					self.E.Put(v.Name, robj)
				} else if self.Fun != nil && ContainsString(self.Fun.LocalNames, v.Name) && env != self.E {
					self.E.Put(v.Name, robj)
				} else {
					env.Put(v.Name, robj)
				}
			case *ast.IndexExpr:
				self.evalExpr(v.X)
				lobj := self.Stack.Pop()
				self.evalExpr(v.Index)
				idx := self.Stack.Pop()
				lobj.Dispatch(self.RT, "__set_index__", idx, robj)
			case *ast.SelectorExpr:
				self.evalExpr(v.X)
				lobj := self.Stack.Pop()
				sel := rt.NewStringObject(v.Sel.Name)
				lobj.Dispatch(self.RT, "__set_property__", sel, robj)
			}
		}
	} else {
		for i := 0; i < len(node.Lhs); i++ {
			self.evalExpr(node.Rhs[i])
			robj := self.Stack.Pop()

			switch v := node.Lhs[i].(type) {
			case *ast.Ident:
				val, _ := self.E.LookUp(v.Name)
				val.(rt.Object).Dispatch(self.RT, OpFuncs[node.Tok], robj)
			case *ast.IndexExpr:
				// a[b] += c
				self.evalExpr(v.X)
				lobj := self.Stack.Pop()
				self.evalExpr(v.Index)
				idx := self.Stack.Pop()
				rets := lobj.Dispatch(self.RT, "__get_index__", idx)
				rets[0].Dispatch(self.RT, OpFuncs[node.Tok], robj)
			case *ast.SelectorExpr:
				self.evalExpr(v.X)
				lobj := self.Stack.Pop()
				sel := rt.NewStringObject(v.Sel.Name)
				rets := lobj.Dispatch(self.RT, "__get_property__", sel)
				rets[0].Dispatch(self.RT, OpFuncs[node.Tok], robj)
			}
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

	if node.Tok == token.BREAK {
		self.NeedBreak = true
	}

	if node.Tok == token.CONTINUE {
		self.NeedContinue = true
	}

}

func (self *Eval) VisitBlockStmt(node *ast.BlockStmt) {
	self.E = env.NewEnv(self.E)
	for _, stmt := range node.List {
		// need break in all loop
		if self.NeedReturn {
			break
		}
		if self.LoopDepth > 0 && self.NeedBreak {
			break
		}
		if self.LoopDepth > 0 && self.NeedContinue {
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

	if cond.(*rt.BoolObject).Val {
		node.Body.Accept(self)
	} else if node.Else != nil {
		node.Else.Accept(self)
	}
}

func (self *Eval) VisitCaseClause(node *ast.CaseClause) {
	self.debug(node)

	initObj := self.Stack.Pop()

	// default
	if node.List != nil {
		for _, e := range node.List {
			_, ok := e.(*ast.BasicLit)
			self.evalExpr(e)
			if ok {
				v := self.Stack.Pop()
				rets := initObj.Dispatch(self.RT, "__eql__", v)
				if rets[0].(*rt.BoolObject).Val == false {
					self.Stack.Push(rt.NewBoolObject(false))
					return
				}
			} else {
				v := self.Stack.Pop().(*rt.BoolObject)
				if v.Val == false {
					self.Stack.Push(rt.NewBoolObject(false))
					return
				}
			}
		}
	}

	for _, s := range node.Body {
		// need break in all loop
		if self.NeedReturn {
			break
		}
		if self.LoopDepth > 0 && self.NeedBreak {
			break
		}
		if self.LoopDepth > 0 && self.NeedContinue {
			break
		}
		s.Accept(self)
	}

	self.Stack.Push(rt.NewBoolObject(true))
}

func (self *Eval) VisitSwitchStmt(node *ast.SwitchStmt) {
	self.debug(node)

	node.Init.Accept(self)
	initObj := self.Stack.Pop()
	for _, c := range node.Body.List {
		self.Stack.Push(initObj)
		c.Accept(self)
		hit := self.Stack.Pop().(*rt.BoolObject)
		if hit.Val {
			break
		}
	}
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
		if !cond.(*rt.BoolObject).Val {
			break
		}

		self.LoopDepth++
		node.Body.Accept(self)
		self.LoopDepth--

		if self.NeedReturn {
			break
		}
		if self.NeedBreak {
			self.NeedBreak = false
			break
		}
		if self.NeedContinue {
			self.NeedContinue = false
		}
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

	self.E = env.NewEnv(self.E)

	switch v := obj.(type) {
	case *rt.ArrayObject:
		for i, val := range v.Vals {
			self.E.Put(keyName, rt.NewIntegerObject(i))
			self.E.Put(valName, val)

			self.LoopDepth++
			node.Body.Accept(self)
			self.LoopDepth--

			if self.NeedReturn {
				break
			}
			if self.NeedBreak {
				self.NeedBreak = false
				break
			}
			if self.NeedContinue {
				self.NeedContinue = false
			}
		}
	case *rt.SetObject:
		for i, val := range v.Vals {
			self.E.Put(keyName, rt.NewIntegerObject(i))
			self.E.Put(valName, val)

			self.LoopDepth++
			node.Body.Accept(self)
			self.LoopDepth--

			if self.NeedReturn {
				break
			}
			if self.NeedBreak {
				self.NeedBreak = false
				break
			}
			if self.NeedContinue {
				self.NeedContinue = false
			}
		}
	}

	self.E = self.E.Outer
}
