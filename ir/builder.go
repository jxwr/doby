package ir

import (
	"fmt"

	"github.com/jxwr/doubi/ast"
	"github.com/jxwr/doubi/ir/instr"
	"github.com/jxwr/doubi/token"
)

type IRBuilder struct {
	cs []instr.Instr
}

func NewIRBuilder() *IRBuilder {
	irb := &IRBuilder{[]instr.Instr{}}
	return irb
}

func (self *IRBuilder) AddInstr(instr instr.Instr) {
	self.cs = append(self.cs, instr)
}

func (self *IRBuilder) BuildExpr(node ast.Expr) {
	switch n := node.(type) {
	case *ast.Ident:
		self.BuildIdent(n)
	case *ast.BasicLit:
		self.BuildBasicLit(n)
	case *ast.BinaryExpr:
		self.BuildBinary(n)
	default:
		fmt.Printf("%#v\n", node)
	}
}

func (self *IRBuilder) BuildStmt(node ast.Stmt) {
	switch n := node.(type) {
	case *ast.ExprStmt:
		self.BuildExpr(n.X)
	default:
		fmt.Printf("%#v\n", node)
	}
}

func (self *IRBuilder) BuildIdent(node *ast.Ident) {
	fmt.Println("push ident", node.Name)
}

func (self *IRBuilder) BuildBasicLit(node *ast.BasicLit) {
	fmt.Println("push ident", node.Value)
}

func (self *IRBuilder) BuildBinary(node *ast.BinaryExpr) {
	self.BuildExpr(node.X)
	self.BuildExpr(node.Y)
	self.AddInstr(instr.MkSendMethod(OpFuncs[node.Op], 1))
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
