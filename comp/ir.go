package comp

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/parser"
	"github.com/jxwr/doubi/rt"
	"github.com/jxwr/doubi/token"
)

type Closure struct {
	LocalVariables []string
	InnerClosures  []*Closure
	OuterClosure   *Closure
	Instrs         []string
	Args           []string
}

func NewClosure(outer *Closure) *Closure {
	c := &Closure{
		LocalVariables: []string{},
		InnerClosures:  []*Closure{},
		OuterClosure:   outer,
		Instrs:         []string{},
		Args:           []string{},
	}
	return c
}

func (self *Closure) emit(instr string) {
	self.Instrs = append(self.Instrs, instr)
}

func (self *Closure) DumpClosure() {
	fmt.Println("------------------------------")
	for i, instr := range self.Instrs {
		fmt.Printf("%3d: %s\n", i, instr)
	}
	for _, ic := range self.InnerClosures {
		ic.DumpClosure()
	}
}

func (self *Closure) AddClosure(c *Closure) int {
	self.InnerClosures = append(self.InnerClosures, c)
	return len(self.InnerClosures)
}

func (self *Closure) AddArg(name string) int {
	self.Args = append(self.Args, name)
	return len(self.Args)
}

func (self *Closure) Put(name string) {
	self.LocalVariables = append(self.LocalVariables, name)
}

func (self *Closure) LookUp(name string) (exist bool, depth int) {
	c := self
	for c != nil {
		if ContainsString(c.LocalVariables, name) {
			exist = true
			return
		}
		depth++
		c = c.OuterClosure
	}
	exist = false
	return
}

type IRBuilder struct {
	ModuleNames []string
	Fun         *ast.FuncDeclExpr
	RT          *rt.Runtime
	C           *Closure

	lexer *parser.Lexer
}

func NewIRBuilder() *IRBuilder {
	irb := &IRBuilder{C: NewClosure(nil)}
	return irb
}

func (self *IRBuilder) SetRuntime(rt *rt.Runtime) {
	self.RT = rt
}

func (self *IRBuilder) SetLexer(lexer *parser.Lexer) {
	self.lexer = lexer
}

func (self *IRBuilder) buildExpr(expr ast.Expr) {
	expr.Accept(self)
}

func (self *IRBuilder) emit(instr string) {
	self.C.emit(instr)
}

func (self *IRBuilder) PushClosure() int {
	c := NewClosure(self.C)
	n := self.C.AddClosure(c)
	self.C = c
	return n
}

func (self *IRBuilder) PopClosure() *Closure {
	c := self.C
	self.C = self.C.OuterClosure
	return c
}

func (self *IRBuilder) DumpClosure() {
	self.C.DumpClosure()
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
		self.emit("push_true")
	} else if node.Name == "false" {
		self.emit("push_false")
	} else {
		// self.Fatalf(node.NamePos, "'%s' not Found", node.Name)
	}
}

func (self *IRBuilder) VisitBasicLit(node *ast.BasicLit) {
	switch node.Kind {
	case token.INT:
		val, err := strconv.Atoi(node.Value)
		if err != nil {
			self.Fatalf(node.ValuePos, "%s convert to int failed: %v", node.Value, err)
		}
		self.emit(fmt.Sprintf("push_int %d", val))
	case token.FLOAT:
		val, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			self.Fatalf(node.ValuePos, "%s convert to float failed: %v", node.Value, err)
		}
		self.emit(fmt.Sprintf("push_float %f", val))
	case token.STRING:
		val := strings.Trim(node.Value, "\"")
		self.emit(fmt.Sprintf("push_string %s", val))
	case token.CHAR:
		val := strings.Trim(node.Value, "'")
		self.emit(fmt.Sprintf("push_string %s", val))
	}
}

func (self *IRBuilder) VisitParenExpr(node *ast.ParenExpr) {
	node.X.Accept(self)
}

func (self *IRBuilder) VisitSelectorExpr(node *ast.SelectorExpr) {
	self.buildExpr(node.X)
	self.emit(fmt.Sprintf("push_string %s", node.Sel.Name))
	self.emit("send_method :get_property")
}

func (self *IRBuilder) VisitIndexExpr(node *ast.IndexExpr) {
	self.buildExpr(node.X)
	self.buildExpr(node.Index)
	self.emit("send_method :get_index")
}

func (self *IRBuilder) VisitSliceExpr(node *ast.SliceExpr) {
	self.buildExpr(node.X)

	if node.Low == nil {
		self.emit("push_nil")
	} else {
		self.buildExpr(node.Low)
	}

	if node.High == nil {
		self.emit("push_nil")
	} else {
		self.buildExpr(node.High)
	}

	self.emit("send_method :slice")
}

func (self *IRBuilder) VisitCallExpr(node *ast.CallExpr) {
	for _, arg := range node.Args {
		self.buildExpr(arg)
	}

	self.buildExpr(node.Fun)
	self.emit(fmt.Sprintf("send_stack %d", len(node.Args)))
}

func (self *IRBuilder) VisitUnaryExpr(node *ast.UnaryExpr) {
	self.buildExpr(node.X)
	if node.Op == token.NOT {
		self.emit("send_method :__not__")
	} else if node.Op == token.SUB {
		self.emit("send_method :__minus__")
	}
}

func (self *IRBuilder) VisitBinaryExpr(node *ast.BinaryExpr) {
	self.buildExpr(node.Y)
	self.buildExpr(node.X)

	self.emit(fmt.Sprintf("send_stack :%s %d", OpFuncs[node.Op], 1))
}

func (self *IRBuilder) VisitArrayExpr(node *ast.ArrayExpr) {
	for _, elem := range node.Elems {
		self.buildExpr(elem)
	}
	self.emit(fmt.Sprintf("new_array %d", len(node.Elems)))
}

func (self *IRBuilder) VisitSetExpr(node *ast.SetExpr) {
	for _, elem := range node.Elems {
		self.buildExpr(elem)
	}
	self.emit(fmt.Sprintf("new_set %d", len(node.Elems)))
}

func (self *IRBuilder) VisitDictExpr(node *ast.DictExpr) {
	for _, field := range node.Fields {
		self.buildExpr(field.Name)
		self.buildExpr(field.Value)
	}
	self.emit(fmt.Sprintf("new_dict %d", len(node.Fields)))
}

func (self *IRBuilder) VisitFuncDeclExpr(node *ast.FuncDeclExpr) {
	n := self.PushClosure()
	for _, arg := range node.Args {
		self.C.AddArg(arg.Name)
	}
	node.Body.Accept(self)
	self.PopClosure()

	self.emit(fmt.Sprintf("push_closure %d", n))
	if node.Name != nil {
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
		self.emit("send_method :__inc__")
	} else if node.Tok == token.DEC {
		self.emit("send_method :__dec__")
	}
}

func (self *IRBuilder) VisitAssignStmt(node *ast.AssignStmt) {
	if node.Tok == token.ASSIGN {
		for i := 0; i < len(node.Rhs); i++ {
			self.buildExpr(node.Rhs[i])
		}
	} else {
		for i := 0; i < len(node.Lhs); i++ {
			self.buildExpr(node.Rhs[i])
		}
	}
}

func (self *IRBuilder) VisitGoStmt(node *ast.GoStmt) {
	self.emit("go_prepare")
	node.Call.Accept(self)
	self.emit("go_run")
}

func (self *IRBuilder) VisitReturnStmt(node *ast.ReturnStmt) {
	for _, res := range node.Results {
		self.buildExpr(res)
	}

	self.emit(fmt.Sprintf("raise_return %d", len(node.Results)))
}

func (self *IRBuilder) VisitBranchStmt(node *ast.BranchStmt) {
	if node.Tok == token.BREAK {
		self.emit("raise_break")
	}

	if node.Tok == token.CONTINUE {
		self.emit("raise_continue")
	}
}

func (self *IRBuilder) VisitBlockStmt(node *ast.BlockStmt) {
	self.emit("push_scope")
	for _, stmt := range node.List {
		stmt.Accept(self)
	}
	self.emit("pop_scope")
}

func (self *IRBuilder) VisitIfStmt(node *ast.IfStmt) {
	self.buildExpr(node.Cond)
	if node.Else != nil {
		self.emit("jump_if_false label_false")
	} else {
		self.emit("jump_if_false label_end")
	}
	node.Body.Accept(self)
	if node.Else != nil {
		self.emit("jump label_end")
		self.emit("label_false:")
		node.Else.Accept(self)
	}
	self.emit("label_end:")
}

func (self *IRBuilder) VisitCaseClause(node *ast.CaseClause) {
	// default
	if node.List != nil {
		for _, e := range node.List {
			_, ok := e.(*ast.BasicLit)
			self.buildExpr(e)
			if ok {
				self.emit("send_method :__eql__")
			}
			self.emit("jump_if_false label_out")
		}
	}

	for _, s := range node.Body {
		s.Accept(self)
	}

	self.emit("label_out:")
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

	self.emit("label_cond:")
	self.buildExpr(node.Cond)
	self.emit("jump_if_false label_out")
	node.Body.Accept(self)
	node.Post.Accept(self)
	self.emit("jump label_cond")
	self.emit("label_out:")
}

func (self *IRBuilder) VisitRangeStmt(node *ast.RangeStmt) {
	self.buildExpr(node.X)
}

func (self *IRBuilder) VisitImportStmt(node *ast.ImportStmt) {
	if len(node.Modules) == 1 {
		modname := node.Modules[0]
		modname = strings.Trim(modname, "\" ")
		self.emit("push_string " + modname)
		self.emit("import")
	}
}
