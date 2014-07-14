package comp

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/env"
	"github.com/jxwr/doubi/parser"
	"github.com/jxwr/doubi/rt"
	"github.com/jxwr/doubi/token"
)

type Stack struct {
	cur  int
	mark int
	vals []rt.Object
}

func NewStack() *Stack {
	stack := &Stack{0, 0, []rt.Object{}}
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
		panic("Pop from empty stack, maybe missing return in some func")
	}
	self.cur--
	return self.vals[self.cur]
}

func (self *Stack) Mark() {
	self.mark = self.cur
}

func (self *Stack) Rewind() {
	self.cur = self.mark
}

type Eval struct {
	Debug bool
	E     *env.Env
	Stack *Stack
	Fun   *ast.FuncDeclExpr
	RT    *rt.Runtime

	lexer *parser.Lexer

	needReturn   bool
	loopDepth    int
	needBreak    bool
	needContinue bool
}

func NewEvaluater() *Eval {
	eval := &Eval{E: env.NewEnv(nil), Stack: NewStack()}
	return eval
}

func (self *Eval) SetRuntime(rt *rt.Runtime) {
	self.RT = rt
}

func (self *Eval) SetLexer(lexer *parser.Lexer) {
	self.lexer = lexer
}

func (self *Eval) Fatalf(pos token.Pos, format string, a ...interface{}) {
	if pos > 0 {
		self.lexer.PrintPosInfo(int(pos))
	}
	fmt.Printf("Error: "+format, a...)
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
		obj := self.RT.NewBoolObject(true)
		self.Stack.Push(obj)
	} else if node.Name == "false" {
		obj := self.RT.NewBoolObject(false)
		self.Stack.Push(obj)
	} else {
		obj, _ := self.E.LookUp(node.Name)
		if obj != nil {
			self.Stack.Push(obj.(rt.Object))
		} else {
			self.Fatalf(node.NamePos, "'%s' not Found", node.Name)
		}
	}
}

func (self *Eval) VisitBasicLit(node *ast.BasicLit) {
	self.debug(node)

	switch node.Kind {
	case token.INT:
		val, err := strconv.Atoi(node.Value)
		if err != nil {
			self.Fatalf(node.ValuePos, "%s convert to int failed: %v", node.Value, err)
		}
		obj := self.RT.NewIntegerObject(val)
		self.Stack.Push(obj)
	case token.FLOAT:
		val, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			self.Fatalf(node.ValuePos, "%s convert to float failed: %v", node.Value, err)
		}
		obj := self.RT.NewFloatObject(val)
		self.Stack.Push(obj)
	case token.STRING:
		val := strings.Trim(node.Value, "\"")
		obj := self.RT.NewStringObject(val)
		self.Stack.Push(obj)
	case token.CHAR:
		val := strings.Trim(node.Value, "'")
		obj := self.RT.NewStringObject(val)
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
	prop := self.RT.NewStringObject(node.Sel.Name)
	rets := rt.Invoke(self.RT, obj, "__get_property__", prop)
	self.Stack.Push(rets[0])
}

func (self *Eval) VisitIndexExpr(node *ast.IndexExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop()
	self.evalExpr(node.Index)
	index := self.Stack.Pop()
	rets := rt.Invoke(self.RT, obj, "__get_index__", index)
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

	rets := rt.Invoke(self.RT, obj, "__slice__", lowObj, highObj)
	self.Stack.Push(rets[0])
}

func (self *Eval) VisitCallExpr(node *ast.CallExpr) {
	self.evalExpr(node.Fun)

	switch fnobj := self.Stack.Pop().(type) {
	case *rt.FuncObject:
		// class methods
		if fnobj.IsBuiltin {
			args := []rt.Object{}
			for _, arg := range node.Args {
				self.evalExpr(arg)
				args = append(args, self.Stack.Pop())
			}
			self.E = env.NewEnv(self.E)
			fnobj.E = self.E
			rets := rt.Invoke(self.RT, fnobj, "__call__", args...)
			self.E = self.E.Outer
			for _, ret := range rets {
				self.Stack.Push(ret)
			}
		} else {
			fnDecl := fnobj.Decl
			fnBak := self.Fun
			self.Fun = fnDecl

			newEnv := env.NewEnv(fnobj.E)
			for i, arg := range node.Args {
				self.evalExpr(arg)
				newEnv.Put(fnDecl.Args[i].Name, self.Stack.Pop())
			}

			bakEnv := self.E
			self.E = newEnv
			fnobj.E = self.E
			self.needReturn = false
			fnDecl.Body.Accept(self)
			self.needReturn = false
			self.Fun = fnBak
			self.E = bakEnv
		}
	case *rt.GoFuncObject:
		args := []rt.Object{}
		for _, arg := range node.Args {
			self.evalExpr(arg)
			args = append(args, self.Stack.Pop())
		}
		self.E = env.NewEnv(self.E)
		rets := rt.Invoke(self.RT, fnobj, "__call__", args...)
		self.E = self.E.Outer
		for _, ret := range rets {
			self.Stack.Push(ret)
		}
	}
}

func (self *Eval) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.debug(node)

	if node.Op == token.NOT {
		self.evalExpr(node.X)
		obj := self.Stack.Pop()
		rets := rt.Invoke(self.RT, obj, "__not__")
		for _, ret := range rets {
			self.Stack.Push(ret)
		}
	} else if node.Op == token.SUB {
		self.evalExpr(node.X)
		obj := self.Stack.Pop()
		var val rt.Object
		switch obj := obj.(type) {
		case *rt.IntegerObject:
			val = self.RT.NewIntegerObject(-obj.Val)
		case *rt.FloatObject:
			val = self.RT.NewFloatObject(-obj.Val)
		}
		self.Stack.Push(val)
	}
}

var OpFuncs = map[token.Token]string{
	token.ADD:            "__add__",
	token.SUB:            "__sub__",
	token.MUL:            "__mul__",
	token.QUO:            "__quo__",
	token.REM:            "__rem__",
	token.AND:            "__and__",
	token.OR:             "__or__",
	token.NOT:            "__not__",
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
	token.ADD_ASSIGN:     "__add_assign__",
	token.SUB_ASSIGN:     "__sub_assign__",
	token.MUL_ASSIGN:     "__mul_assign__",
	token.QUO_ASSIGN:     "__quo_assign__",
	token.REM_ASSIGN:     "__rem_assign__",
	token.AND_ASSIGN:     "__and_assign__",
	token.OR_ASSIGN:      "__or_assign__",
	token.XOR_ASSIGN:     "__xor_assign__",
	token.SHL_ASSIGN:     "__shl_assign__",
	token.SHR_ASSIGN:     "__shr_assign__",
	token.AND_NOT_ASSIGN: "__and_not_assign__",
}

func (self *Eval) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.debug(node)

	self.evalExpr(node.X)
	self.evalExpr(node.Y)

	robj := self.Stack.Pop()
	lobj := self.Stack.Pop()

	objs := rt.Invoke(self.RT, lobj, OpFuncs[node.Op], robj)
	self.Stack.Push(objs[0])
}

func (self *Eval) VisitArrayExpr(node *ast.ArrayExpr) {
	self.debug(node)

	elems := []rt.Object{}
	for _, elem := range node.Elems {
		self.evalExpr(elem)
		elems = append(elems, self.Stack.Pop())
	}
	obj := self.RT.NewArrayObject(elems)
	self.Stack.Push(obj)
}

func (self *Eval) VisitSetExpr(node *ast.SetExpr) {
	self.debug(node)

	elems := []rt.Object{}
	for _, elem := range node.Elems {
		self.evalExpr(elem)
		elems = append(elems, self.Stack.Pop())
	}
	obj := self.RT.NewSetObject(elems)
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
	obj := self.RT.NewDictObject(fieldMap)
	self.Stack.Push(obj)
}

func (self *Eval) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	self.debug(node)

	if node.Name != nil {
		fname := node.Name.Name
		self.E.Put(fname, self.RT.NewFuncObject(fname, node, self.E))
	} else {
		self.Stack.Push(self.RT.NewFuncObject("#<closure>", node, self.E))
	}
}

// stmts

func (self *Eval) VisitExprStmt(node *ast.ExprStmt) {
	self.debug(node)

	self.Stack.Mark()
	node.X.Accept(self)
	self.Stack.Rewind()
}

func (self *Eval) VisitSendStmt(node *ast.SendStmt) {
	self.debug(node)
}

func (self *Eval) VisitIncDecStmt(node *ast.IncDecStmt) {
	self.debug(node)

	self.evalExpr(node.X)
	obj := self.Stack.Pop()

	if node.Tok == token.INC {
		rt.Invoke(self.RT, obj, "__inc__")
	} else if node.Tok == token.DEC {
		rt.Invoke(self.RT, obj, "__dec__")
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
		rhs := []rt.Object{}

		for i := 0; i < len(node.Rhs); i++ {
			self.evalExpr(node.Rhs[i])
		}

		llen := len(node.Lhs)
		for i := 0; i < llen; i++ {
			robj := self.Stack.Pop()
			rhs = append(rhs, robj)
		}

		for i := 0; i < llen; i++ {
			robj := rhs[llen-i-1]

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
				rt.Invoke(self.RT, lobj, "__set_index__", idx, robj)
			case *ast.SelectorExpr:
				self.evalExpr(v.X)
				lobj := self.Stack.Pop()
				sel := self.RT.NewStringObject(v.Sel.Name)
				rt.Invoke(self.RT, lobj, "__set_property__", sel, robj)
			}
		}
	} else {
		for i := 0; i < len(node.Lhs); i++ {
			self.evalExpr(node.Rhs[i])
			robj := self.Stack.Pop()

			switch v := node.Lhs[i].(type) {
			case *ast.Ident:
				val, _ := self.E.LookUp(v.Name)
				rt.Invoke(self.RT, val.(rt.Object), OpFuncs[node.Tok], robj)
			case *ast.IndexExpr:
				// a[b] += c
				self.evalExpr(v.X)
				lobj := self.Stack.Pop()
				self.evalExpr(v.Index)
				idx := self.Stack.Pop()
				rets := rt.Invoke(self.RT, lobj, "__get_index__", idx)
				rt.Invoke(self.RT, rets[0], OpFuncs[node.Tok], robj)
			case *ast.SelectorExpr:
				self.evalExpr(v.X)
				lobj := self.Stack.Pop()
				sel := self.RT.NewStringObject(v.Sel.Name)
				rets := rt.Invoke(self.RT, lobj, "__get_property__", sel)
				rt.Invoke(self.RT, rets[0], OpFuncs[node.Tok], robj)
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

	self.needReturn = true
}

func (self *Eval) VisitBranchStmt(node *ast.BranchStmt) {
	self.debug(node)

	if node.Tok == token.BREAK {
		self.needBreak = true
	}

	if node.Tok == token.CONTINUE {
		self.needContinue = true
	}

}

func (self *Eval) VisitBlockStmt(node *ast.BlockStmt) {
	self.E = env.NewEnv(self.E)
	for _, stmt := range node.List {
		// need break in all loop
		if self.needReturn {
			break
		}
		if self.loopDepth > 0 && self.needBreak {
			break
		}
		if self.loopDepth > 0 && self.needContinue {
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
				rets := rt.Invoke(self.RT, initObj, "__eql__", v)
				if rets[0].(*rt.BoolObject).Val == false {
					self.Stack.Push(self.RT.NewBoolObject(false))
					return
				}
			} else {
				v := self.Stack.Pop().(*rt.BoolObject)
				if v.Val == false {
					self.Stack.Push(self.RT.NewBoolObject(false))
					return
				}
			}
		}
	}

	for _, s := range node.Body {
		// need break in all loop
		if self.needReturn {
			break
		}
		if self.loopDepth > 0 && self.needBreak {
			break
		}
		if self.loopDepth > 0 && self.needContinue {
			break
		}
		s.Accept(self)
	}

	self.Stack.Push(self.RT.NewBoolObject(true))
}

func (self *Eval) VisitSwitchStmt(node *ast.SwitchStmt) {
	self.debug(node)

	// dirty hack, we need keep the stack clean
	node.Init.(*ast.ExprStmt).X.Accept(self)
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

		self.loopDepth++
		node.Body.Accept(self)
		self.loopDepth--

		if self.needReturn {
			break
		}
		if self.needBreak {
			self.needBreak = false
			break
		}
		if self.needContinue {
			self.needContinue = false
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
			self.E.Put(keyName, self.RT.NewIntegerObject(i))
			self.E.Put(valName, val)

			self.loopDepth++
			node.Body.Accept(self)
			self.loopDepth--

			if self.needReturn {
				break
			}
			if self.needBreak {
				self.needBreak = false
				break
			}
			if self.needContinue {
				self.needContinue = false
			}
		}
	case *rt.SetObject:
		for i, val := range v.Vals {
			self.E.Put(keyName, self.RT.NewIntegerObject(i))
			self.E.Put(valName, val)

			self.loopDepth++
			node.Body.Accept(self)
			self.loopDepth--

			if self.needReturn {
				break
			}
			if self.needBreak {
				self.needBreak = false
				break
			}
			if self.needContinue {
				self.needContinue = false
			}
		}
	case *rt.DictObject:
		for i, val := range v.Property.Slots {
			self.E.Put(keyName, self.RT.NewStringObject(i))
			self.E.Put(valName, val)

			self.loopDepth++
			node.Body.Accept(self)
			self.loopDepth--

			if self.needReturn {
				break
			}
			if self.needBreak {
				self.needBreak = false
				break
			}
			if self.needContinue {
				self.needContinue = false
			}
		}
	}

	self.E = self.E.Outer
}

func (self *Eval) VisitImportStmt(node *ast.ImportStmt) {
	if len(node.Modules) == 1 {
		modname := node.Modules[0]
		modname = strings.Trim(modname, "\" ")
		mod, _ := self.RT.Env.LookUp(modname)
		xs := strings.Split(modname, "/")
		self.E.Put(xs[len(xs)-1], mod)
	}
}
