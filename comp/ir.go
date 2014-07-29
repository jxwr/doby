package comp

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/parser"
	"github.com/jxwr/doubi/token"
	"github.com/jxwr/doubi/vm/instr"
)

/// IRBuilder

type IRBuilder struct {
	ModuleNames []string
	Fun         *ast.FuncDeclExpr
	C           *instr.ClosureProto
	CS          map[int]*instr.ClosureProto

	lexer *parser.Lexer
}

func NewIRBuilder() *IRBuilder {
	c := instr.NewClosureProto(nil)
	irb := &IRBuilder{
		C:  c,
		CS: map[int]*instr.ClosureProto{0: c},
	}
	return irb
}

func (self *IRBuilder) SetLexer(lexer *parser.Lexer) {
	self.lexer = lexer
}

func (self *IRBuilder) buildExpr(expr ast.Expr) {
	expr.Accept(self)
}

func (self *IRBuilder) emit(instr instr.Instr) {
	self.C.Emit(instr)
}

func (self *IRBuilder) PushClosureProto() int {
	c := instr.NewClosureProto(self.C)
	self.C.AddClosureProto(c)
	self.C = c
	self.CS[c.Seq] = c
	return c.Seq
}

func (self *IRBuilder) PopClosureProto() *instr.ClosureProto {
	c := self.C
	self.C = self.C.OuterClosureProto
	return c
}

func (self *IRBuilder) DumpClosureProto() {
	self.C.DumpClosureProto()
}

func (self *IRBuilder) Fatalf(pos token.Pos, format string, a ...interface{}) {
	if pos > 0 {
		self.lexer.PrintPosInfo(int(pos))
	}
	fmt.Printf("Error: "+format, a...)
	fmt.Println()
	os.Exit(1)
}

// exprs

func (self *IRBuilder) VisitIdent(node *ast.Ident) {
	if node.Name == "true" {
		self.emit(instr.PushTrue())
	} else if node.Name == "false" {
		self.emit(instr.PushFalse())
	} else {
		exist, offset := self.C.LookUpLocal(node.Name)
		if exist {
			self.emit(instr.LoadLocal(offset))
			return
		}
		exist, offset = self.C.LookUpUpval(node.Name)
		if exist {
			self.emit(instr.LoadUpval(offset))
			return
		}

		exist, depth, offset := self.C.LookUpOuter(node.Name)
		if exist {
			offset := self.C.AddUpvalVariable(node.Name, depth, offset)
			self.emit(instr.LoadUpval(offset))
		} else if ContainsString(self.ModuleNames, node.Name) {
			self.emit(instr.PushModule(node.Name))
		} else {
			self.Fatalf(node.NamePos, "'%s' not Found", node.Name)
		}
	}
}

func (self *IRBuilder) VisitBasicLit(node *ast.BasicLit) {
	switch node.Kind {
	case token.INT:
		val, err := strconv.Atoi(node.Value)
		if err != nil {
			self.Fatalf(node.ValuePos, "%s convert to int failed: %v", node.Value, err)
		}
		self.emit(instr.PushInt(val))
	case token.FLOAT:
		val, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			self.Fatalf(node.ValuePos, "%s convert to float failed: %v", node.Value, err)
		}
		self.emit(instr.PushFloat(val))
	case token.STRING:
		val := strings.Trim(node.Value, "\"")
		self.emit(instr.PushString(val))
	case token.CHAR:
		val := strings.Trim(node.Value, "'")
		self.emit(instr.PushString(val))
	}
}

func (self *IRBuilder) VisitParenExpr(node *ast.ParenExpr) {
	node.X.Accept(self)
}

func (self *IRBuilder) VisitSelectorExpr(node *ast.SelectorExpr) {
	self.emit(instr.PushString(node.Sel.Name))
	self.buildExpr(node.X)
	self.emit(instr.SendMethod("__get_property__", 1))
}

func (self *IRBuilder) VisitIndexExpr(node *ast.IndexExpr) {
	self.buildExpr(node.Index)
	self.buildExpr(node.X)
	self.emit(instr.SendMethod("__get_index__", 1))
}

func (self *IRBuilder) VisitSliceExpr(node *ast.SliceExpr) {
	if node.Low == nil {
		self.emit(instr.PushNil())
	} else {
		self.buildExpr(node.Low)
	}

	if node.High == nil {
		self.emit(instr.PushNil())
	} else {
		self.buildExpr(node.High)
	}

	self.buildExpr(node.X)
	self.emit(instr.SendMethod("__slice__", 2))
}

func (self *IRBuilder) VisitCallExpr(node *ast.CallExpr) {
	for _, arg := range node.Args {
		self.buildExpr(arg)
	}

	self.buildExpr(node.Fun)
	self.emit(instr.SendMethod("__call__", len(node.Args)))
}

func (self *IRBuilder) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.buildExpr(node.X)
	if node.Op == token.NOT {
		self.emit(instr.SendMethod("__not__", 0))
	} else if node.Op == token.SUB {
		self.emit(instr.SendMethod("__minus__", 0))
	}
}

func (self *IRBuilder) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.buildExpr(node.Y)
	self.buildExpr(node.X)

	self.emit(instr.SendMethod(OpFuncs[node.Op], 1))
}

func (self *IRBuilder) VisitArrayExpr(node *ast.ArrayExpr) {
	for _, elem := range node.Elems {
		self.buildExpr(elem)
	}
	self.emit(instr.NewArray(len(node.Elems)))
}

func (self *IRBuilder) VisitSetExpr(node *ast.SetExpr) {
	for _, elem := range node.Elems {
		self.buildExpr(elem)
	}
	self.emit(instr.NewSet(len(node.Elems)))
}

func (self *IRBuilder) VisitDictExpr(node *ast.DictExpr) {
	for _, field := range node.Fields {
		self.buildExpr(field.Name)
		self.buildExpr(field.Value)
	}
	self.emit(instr.NewDict(len(node.Fields)))
}

func (self *IRBuilder) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	n := self.PushClosureProto()
	for _, arg := range node.Args {
		self.C.AddLocalVariable(arg.Name)
	}
	for i := len(node.Args) - 1; i >= 0; i-- {
		self.emit(instr.SetLocal(i))
	}
	node.Body.Accept(self)
	self.PopClosureProto()

	self.emit(instr.PushClosure(n))

	if node.Name != nil {
		exist, offset := self.C.LookUpLocal(node.Name.Name)
		if !exist {
			offset = self.C.AddLocalVariable(node.Name.Name)
		}
		self.emit(instr.SetLocal(offset))
	}
}

// stmts

func (self *IRBuilder) VisitExprStmt(node *ast.ExprStmt) {
	self.buildExpr(node.X)
}

func (self *IRBuilder) VisitSendStmt(node *ast.SendStmt) {
}

func (self *IRBuilder) VisitIncDecStmt(node *ast.IncDecStmt) {
	self.buildExpr(node.X)
	if node.Tok == token.INC {
		self.emit(instr.SendMethod("__inc__", 0))
	} else if node.Tok == token.DEC {
		self.emit(instr.SendMethod("__dec__", 0))
	}
}

func (self *IRBuilder) VisitAssignStmt(node *ast.AssignStmt) {
	if node.Tok == token.ASSIGN {
		for i := 0; i < len(node.Rhs); i++ {
			self.buildExpr(node.Rhs[i])
		}

		for i := len(node.Lhs) - 1; i >= 0; i-- {
			switch v := node.Lhs[i].(type) {
			case *ast.Ident:
				exist, offset := self.C.LookUpLocal(v.Name)
				if exist {
					self.emit(instr.SetLocal(offset))
					continue
				}
				exist, offset = self.C.LookUpUpval(v.Name)
				if exist {
					self.emit(instr.SetUpval(offset))
					continue
				}
				_, depth, offset := self.C.LookUpOuter(v.Name)
				if depth < 1 {
					offset := self.C.AddLocalVariable(v.Name)
					self.emit(instr.SetLocal(offset))
					continue
				} else {
					offset := self.C.AddUpvalVariable(v.Name, depth, offset)
					self.emit(instr.SetUpval(offset))
					continue
				}
			case *ast.IndexExpr:
				self.buildExpr(v.Index)
				self.buildExpr(v.X)
				self.emit(instr.SendMethod("__set_index__", 1))
			case *ast.SelectorExpr:
				self.emit(instr.PushString(v.Sel.Name))
				self.buildExpr(v.X)
				self.emit(instr.SendMethod("__set_property__", 1))
			}
		}
	} else {
		for i := 0; i < len(node.Lhs); i++ {
			self.buildExpr(node.Rhs[i])

			switch v := node.Lhs[i].(type) {
			case *ast.Ident:
				exist, offset := self.C.LookUpLocal(v.Name)
				if exist {
					self.emit(instr.LoadLocal(offset))
					continue
				}
				exist, offset = self.C.LookUpUpval(v.Name)
				if exist {
					self.emit(instr.LoadUpval(offset))
					continue
				}
			case *ast.IndexExpr:
				self.buildExpr(v.Index)
				self.buildExpr(v.X)
				self.emit(instr.SendMethod("__get_index__", 1))
			case *ast.SelectorExpr:
				self.emit(instr.PushString(v.Sel.Name))
				self.buildExpr(v.X)
				self.emit(instr.SendMethod("__get_property__", 1))
			}
			self.emit(instr.SendMethod(OpFuncs[node.Tok], 1))
		}
	}
}

func (self *IRBuilder) VisitGoStmt(node *ast.GoStmt) {
	node.Call.Accept(self)
}

func (self *IRBuilder) VisitReturnStmt(node *ast.ReturnStmt) {
	for _, res := range node.Results {
		self.buildExpr(res)
	}
	self.emit(instr.RaiseReturn(len(node.Results)))
}

func (self *IRBuilder) VisitBranchStmt(node *ast.BranchStmt) {
	if node.Tok == token.BREAK {
		self.emit(instr.RaiseBreak())
	}

	if node.Tok == token.CONTINUE {
		self.emit(instr.RaiseContinue())
	}
}

func (self *IRBuilder) VisitBlockStmt(node *ast.BlockStmt) {
	for _, stmt := range node.List {
		stmt.Accept(self)
	}
}

func (self *IRBuilder) VisitIfStmt(node *ast.IfStmt) {
	self.buildExpr(node.Cond)
	if node.Else != nil {
		self.emit(instr.JumpIfFalse(0))
	} else {
		self.emit(instr.JumpIfFalse(0))
	}
	node.Body.Accept(self)
	if node.Else != nil {
		self.emit(instr.Jump(0))
		self.emit(instr.Label("false"))
		node.Else.Accept(self)
	}
	self.emit(instr.Label("end"))
}

func (self *IRBuilder) VisitCaseClause(node *ast.CaseClause) {
	// default
	if node.List != nil {
		for _, e := range node.List {
			_, ok := e.(*ast.BasicLit)
			self.buildExpr(e)
			if ok {
				self.emit(instr.SendMethod("__eql__", 0))
			}
			self.emit(instr.JumpIfFalse(0))
		}
	}

	for _, s := range node.Body {
		s.Accept(self)
	}

	self.emit(instr.Label("out"))
}

func (self *IRBuilder) VisitSwitchStmt(node *ast.SwitchStmt) {
	node.Init.(*ast.ExprStmt).X.Accept(self)

	for _, c := range node.Body.List {
		c.Accept(self)
	}
}

func (self *IRBuilder) VisitSelectStmt(node *ast.SelectStmt) {
}

func (self *IRBuilder) VisitForStmt(node *ast.ForStmt) {
	if node.Init != nil {
		node.Init.Accept(self)
	}

	self.emit(instr.Label("cond"))
	self.buildExpr(node.Cond)
	self.emit(instr.JumpIfFalse(0))
	node.Body.Accept(self)
	node.Post.Accept(self)
	self.emit(instr.Jump(0))
	self.emit(instr.Label("out"))
}

func (self *IRBuilder) VisitRangeStmt(node *ast.RangeStmt) {
	self.buildExpr(node.X)
}

func (self *IRBuilder) VisitImportStmt(node *ast.ImportStmt) {
	if len(node.Modules) == 1 {
		modname := node.Modules[0]
		modname = strings.Trim(modname, "\" ")
		self.emit(instr.Import(modname))
		xs := strings.Split(modname, "/")
		self.ModuleNames = append(self.ModuleNames, xs[len(xs)-1])
	}
}
