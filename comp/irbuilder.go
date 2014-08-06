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
	cc                 *instr.ClosureProto
	cs                 map[int]*instr.ClosureProto
	lexer              *parser.Lexer
	continueInstrStack []*instr.JumpInstr
	moduleNames        []string
}

func NewIRBuilder() *IRBuilder {
	c := instr.NewClosureProto(nil)
	irb := &IRBuilder{
		cc:                 c,
		cs:                 map[int]*instr.ClosureProto{0: c},
		continueInstrStack: []*instr.JumpInstr{},
	}
	return irb
}

func (self *IRBuilder) SetLexer(lexer *parser.Lexer) {
	self.lexer = lexer
}

func (self *IRBuilder) RootClosure() *instr.ClosureProto {
	return self.cc
}

func (self *IRBuilder) ClosureTable() map[int]*instr.ClosureProto {
	return self.cs
}

func (self *IRBuilder) buildExpr(expr ast.Expr) {
	expr.Accept(self)
}

func (self *IRBuilder) emit(instr instr.Instr) int {
	return self.cc.Emit(instr)
}

func (self *IRBuilder) PushClosureProto() int {
	c := instr.NewClosureProto(self.cc)
	self.cc.AddClosureProto(c)
	self.cc = c
	self.cs[c.Seq()] = c
	return c.Seq()
}

func (self *IRBuilder) PopClosureProto() *instr.ClosureProto {
	c := self.cc
	self.cc = self.cc.OuterClosureProto()
	return c
}

func (self *IRBuilder) DumpClosureProto() {
	self.cc.DumpClosureProto()
}

func (self *IRBuilder) Fatalf(pos token.Pos, format string, a ...interface{}) {
	if pos > 0 {
		self.lexer.PrintPosInfo(int(pos))
	}
	fmt.Printf("Error: "+format, a...)
	fmt.Println()
	os.Exit(1)
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

// exprs

func (self *IRBuilder) VisitIdent(node *ast.Ident) {
	if node.Name == "nil" {
		self.emit(instr.PushNil())
	} else if node.Name == "true" {
		self.emit(instr.PushTrue())
	} else if node.Name == "false" {
		self.emit(instr.PushFalse())
	} else {
		exist, offset := self.cc.LookUpLocal(node.Name)
		if exist {
			self.emit(instr.LoadLocal(offset))
			return
		}
		exist, offset = self.cc.LookUpUpval(node.Name)
		if exist {
			self.emit(instr.LoadUpval(offset))
			return
		}

		exist, depth, offset := self.cc.LookUpOuter(node.Name)
		if exist {
			offset := self.cc.AddUpvalVariable(node.Name, depth, offset)
			self.emit(instr.LoadUpval(offset))
		} else if ContainsString(self.moduleNames, node.Name) {
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
	funNameOffset := 0
	if node.Name != nil {
		exist, offset := self.cc.LookUpLocal(node.Name.Name)
		if !exist {
			offset = self.cc.AddLocalVariable(node.Name.Name)
		}
		funNameOffset = offset
	}

	n := self.PushClosureProto()
	for _, arg := range node.Args {
		self.cc.AddLocalVariable(arg.Name)
	}
	for i := len(node.Args) - 1; i >= 0; i-- {
		self.emit(instr.SetLocal(i))
	}
	node.Body.Accept(self)
	self.PopClosureProto()

	self.emit(instr.PushClosure(n))

	if node.Name != nil {
		self.emit(instr.SetLocal(funNameOffset))
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
				exist, offset := self.cc.LookUpLocal(v.Name)
				if exist {
					self.emit(instr.SetLocal(offset))
					continue
				}
				exist, offset = self.cc.LookUpUpval(v.Name)
				if exist {
					self.emit(instr.SetUpval(offset))
					continue
				}
				_, depth, offset := self.cc.LookUpOuter(v.Name)
				if depth < 1 {
					offset := self.cc.AddLocalVariable(v.Name)
					self.emit(instr.SetLocal(offset))
					continue
				} else {
					offset := self.cc.AddUpvalVariable(v.Name, depth, offset)
					self.emit(instr.SetUpval(offset))
					continue
				}
			case *ast.IndexExpr:
				self.buildExpr(v.Index)
				self.buildExpr(node.Rhs[0])
				self.buildExpr(v.X)
				self.emit(instr.SendMethod("__set_index__", 2))
			case *ast.SelectorExpr:
				self.emit(instr.PushString(v.Sel.Name))
				self.buildExpr(node.Rhs[0])
				self.buildExpr(v.X)
				self.emit(instr.SendMethod("__set_property__", 2))
			}
		}
	} else {
		for i := 0; i < len(node.Lhs); i++ {
			self.buildExpr(node.Rhs[i])

			switch v := node.Lhs[i].(type) {
			case *ast.Ident:
				exist, offset := self.cc.LookUpLocal(v.Name)
				if exist {
					self.emit(instr.LoadLocal(offset))
					goto out
				}
				exist, offset = self.cc.LookUpUpval(v.Name)
				if exist {
					self.emit(instr.LoadUpval(offset))
					goto out
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
		out:
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
		continueInstr := instr.Jump(-1)
		self.emit(continueInstr)
		self.continueInstrStack = append(self.continueInstrStack, continueInstr)
	}
}

func (self *IRBuilder) VisitBlockStmt(node *ast.BlockStmt) {
	for _, stmt := range node.List {
		stmt.Accept(self)
	}
}

func (self *IRBuilder) VisitIfStmt(node *ast.IfStmt) {
	self.buildExpr(node.Cond)
	jmpInstr := instr.JumpIfFalse(-1)
	self.emit(jmpInstr)
	node.Body.Accept(self)
	if node.Else != nil {
		// else {
		endInstr := instr.Jump(-1)
		self.emit(endInstr)
		elsePc := self.emit(instr.Label("if_else_label"))
		jmpInstr.Target = elsePc
		node.Else.Accept(self)
		// }
		endPc := self.emit(instr.Label("if_end_label"))
		endInstr.Target = endPc
	} else {
		endPc := self.emit(instr.Label("if_end_label"))
		jmpInstr.Target = endPc
	}
}

func (self *IRBuilder) VisitCaseClause(node *ast.CaseClause) {
	// dummy
}

var switchSeq int = 0

func (self *IRBuilder) VisitSwitchStmt(node *ast.SwitchStmt) {
	node.Init.(*ast.ExprStmt).X.Accept(self)
	switchExprOffset := self.cc.AddLocalVariable(fmt.Sprintf("#switch%d#", switchSeq))

	var lastCaseJumpInstr *instr.JumpIfFalseInstr
	endJmpList := []*instr.JumpInstr{}

	self.emit(instr.SetLocal(switchExprOffset))
	for _, c := range node.Body.List {
		caseClause := c.(*ast.CaseClause)

		// set last case clause jumpiffalse target
		caseStart := self.emit(instr.Label("case_start_label"))
		if lastCaseJumpInstr != nil {
			lastCaseJumpInstr.Target = caseStart
		}

		// emit the cmp
		if caseClause.List != nil {
			for _, e := range caseClause.List {
				_, ok := e.(*ast.BasicLit)
				self.buildExpr(e)
				self.emit(instr.LoadLocal(switchExprOffset))
				if ok {
					self.emit(instr.SendMethod("__eql__", 1))
				}
				jmp := instr.JumpIfFalse(-1)
				lastCaseJumpInstr = jmp
				self.emit(jmp)
			}
		}

		// emit the case body
		for _, s := range caseClause.Body {
			s.Accept(self)
		}

		// jump to end after every case body
		jmp := instr.Jump(-1)
		endJmpList = append(endJmpList, jmp)
		self.emit(jmp)
	}

	outpc := self.emit(instr.Label("switch_out_label"))
	for _, ins := range endJmpList {
		ins.Target = outpc
	}

	switchSeq++
}

func (self *IRBuilder) VisitSelectStmt(node *ast.SelectStmt) {
}

func (self *IRBuilder) VisitForStmt(node *ast.ForStmt) {
	pushBlockInstr := instr.PushBlock(-1)
	self.emit(pushBlockInstr)

	if node.Init != nil {
		node.Init.Accept(self)
	}

	condPc := self.emit(instr.Label("for_cond_label"))
	self.buildExpr(node.Cond)
	condJumpInstr := instr.JumpIfFalse(-1)
	self.emit(condJumpInstr)
	node.Body.Accept(self)
	postPc := self.emit(instr.Label("for_post_label"))
	for _, instr := range self.continueInstrStack {
		instr.Target = postPc
	}
	self.continueInstrStack = self.continueInstrStack[:0]

	if node.Post != nil {
		node.Post.Accept(self)
	}
	self.emit(instr.Jump(condPc))
	endPc := self.emit(instr.Label("for_end"))
	condJumpInstr.Target = endPc

	popBlockInstr := instr.PopBlock(-1)
	pc := self.emit(popBlockInstr)
	pushBlockInstr.Target = pc
}

var iterSeq int = 0

func (self *IRBuilder) VisitRangeStmt(node *ast.RangeStmt) {
	pushBlockInstr := instr.PushBlock(-1)
	self.emit(pushBlockInstr)

	keyName := node.KeyValue[0].(*ast.Ident).Name
	valName := node.KeyValue[1].(*ast.Ident).Name

	keyOffset := self.cc.AddLocalVariable(keyName)
	valOffset := self.cc.AddLocalVariable(valName)
	iterOffset := self.cc.AddLocalVariable(fmt.Sprintf("#iter%d#", iterSeq))
	xOffset := self.cc.AddLocalVariable(fmt.Sprintf("#iter#x%d#", iterSeq))
	iterSeq++

	self.emit(instr.PushInt(0))
	self.emit(instr.SetLocal(iterOffset))

	self.buildExpr(node.X)
	self.emit(instr.SetLocal(xOffset))

	beginLabel := self.emit(instr.Label("for_range_start"))
	self.emit(instr.LoadLocal(iterOffset))
	self.emit(instr.LoadLocal(xOffset))

	self.emit(instr.SendMethod("__iter__", 1))
	condJump := instr.JumpIfFalse(-1)
	self.emit(condJump)
	self.emit(instr.SetLocal(valOffset))
	self.emit(instr.SetLocal(keyOffset))
	node.Body.Accept(self)
	self.emit(instr.LoadLocal(iterOffset))
	self.emit(instr.SendMethod("__inc__", 0))
	self.emit(instr.Jump(beginLabel))
	endLabel := self.emit(instr.Label("for_range_end"))
	condJump.Target = endLabel

	popBlockInstr := instr.PopBlock(-1)
	pc := self.emit(popBlockInstr)
	pushBlockInstr.Target = pc
}

func (self *IRBuilder) VisitImportStmt(node *ast.ImportStmt) {
	if len(node.Modules) == 1 {
		modname := node.Modules[0]
		modname = strings.Trim(modname, "\" ")
		self.emit(instr.Import(modname))
		xs := strings.Split(modname, "/")
		self.moduleNames = append(self.moduleNames, xs[len(xs)-1])
	}
}
