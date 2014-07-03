package ast

import (
	"github.com/jxwr/doubi/token"
)

type Node interface {
	//	Pos() int
	//	End() int
	Accept(Visitor)
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}

// Expression
type Ident struct {
	NamePos token.Pos
	Name    string
}

type BasicLit struct {
	ValuePos token.Pos
	Kind     token.Token
	Value    string
}

type ParenExpr struct {
	Lparen token.Pos
	X      Expr
	Rparen token.Pos
}

type SelectorExpr struct {
	X   Expr
	Sel *Ident
}

type IndexExpr struct {
	X      Expr
	Lbrack token.Pos
	Index  Expr
	Rbrack token.Pos
}

type SliceExpr struct {
	X      Expr
	Lbrack token.Pos
	Low    Expr
	High   Expr
	Rbrack token.Pos
}

type CallExpr struct {
	Fun    Expr
	Lparen token.Pos
	Args   []Expr
	Rparen token.Pos
}

type UnaryExpr struct {
	OpPos token.Pos
	Op    token.Token
	X     Expr
}

type BinaryExpr struct {
	X     Expr
	OpPos token.Pos
	Op    token.Token
	Y     Expr
}

type ArrayExpr struct {
	Lbrack token.Pos
	Elems  []Expr
	Rbrack token.Pos
}

type SetExpr struct {
	Lbrack token.Pos
	Elems  []Expr
	Rbrack token.Pos
}

type Field struct {
	Name     Expr
	ColonPos token.Pos
	Value    Expr
}

type DictExpr struct {
	Lbrace token.Pos
	Fields []*Field
	Rbrace token.Pos
}

type FuncDeclExpr struct {
	Func     token.Pos
	Recv     *Ident
	RecvType *Ident
	Name     *Ident
	Args     []*Ident
	Body     *BlockStmt
}

func (n Ident) exprNode()        {}
func (n BasicLit) exprNode()     {}
func (n ParenExpr) exprNode()    {}
func (n SelectorExpr) exprNode() {}
func (n IndexExpr) exprNode()    {}
func (n SliceExpr) exprNode()    {}
func (n CallExpr) exprNode()     {}
func (n UnaryExpr) exprNode()    {}
func (n BinaryExpr) exprNode()   {}
func (n ArrayExpr) exprNode()    {}
func (n SetExpr) exprNode()      {}
func (n DictExpr) exprNode()     {}
func (n FuncDeclExpr) exprNode() {}

func (n *Ident) Accept(v Visitor) {
	v.VisitIdent(n)
}

func (n *BasicLit) Accept(v Visitor) {
	v.VisitBasicLit(n)
}

func (n *ParenExpr) Accept(v Visitor) {
	v.VisitParenExpr(n)
}

func (n *SelectorExpr) Accept(v Visitor) {
	v.VisitSelectorExpr(n)
}

func (n *IndexExpr) Accept(v Visitor) {
	v.VisitIndexExpr(n)
}

func (n *SliceExpr) Accept(v Visitor) {
	v.VisitSliceExpr(n)
}

func (n *CallExpr) Accept(v Visitor) {
	v.VisitCallExpr(n)
}

func (n *UnaryExpr) Accept(v Visitor) {
	v.VisitUnaryExpr(n)
}

func (n *BinaryExpr) Accept(v Visitor) {
	v.VisitBinaryExpr(n)
}

func (n *ArrayExpr) Accept(v Visitor) {
	v.VisitArrayExpr(n)
}

func (n *SetExpr) Accept(v Visitor) {
	v.VisitSetExpr(n)
}

func (n *DictExpr) Accept(v Visitor) {
	v.VisitDictExpr(n)
}

func (n *FuncDeclExpr) Accept(v Visitor) {
	v.VisitFuncDeclExpr(n)
}

/// Stmts

type ExprStmt struct {
	X Expr
}

type SendStmt struct {
	Chan  Expr
	Arrow token.Pos
	Value Expr
}

type IncDecStmt struct {
	X      Expr
	TokPos token.Pos
	Tok    token.Token
}

type AssignStmt struct {
	Lhs    []Expr
	TokPos token.Pos
	Tok    token.Token
	Rhs    []Expr
}

type GoStmt struct {
	Go   token.Pos
	Call *CallExpr
}

type ReturnStmt struct {
	Return  token.Pos
	Results []Expr
}

type BranchStmt struct {
	TokPos token.Pos
	Tok    token.Token
}

type BlockStmt struct {
	Lbrace token.Pos
	List   []Stmt
	Rbrack token.Pos
}

type IfStmt struct {
	If   token.Pos
	Cond Expr
	Body *BlockStmt
	Else Stmt
}

type CaseClause struct {
	Case  token.Pos
	List  []Expr
	Colon token.Pos
	Body  []Stmt
}

type SwitchStmt struct {
	Switch token.Pos
	Init   Stmt
	Body   *BlockStmt
}

type SelectStmt struct {
	Select token.Pos
	Body   *BlockStmt
}

type ForStmt struct {
	For  token.Pos
	Init Stmt
	Cond Expr
	Post Stmt
	Body *BlockStmt
}

type RangeStmt struct {
	For      token.Pos
	KeyValue []Expr
	X        Expr
	Body     *BlockStmt
}

func (ExprStmt) stmtNode()   {}
func (SendStmt) stmtNode()   {}
func (IncDecStmt) stmtNode() {}
func (AssignStmt) stmtNode() {}
func (GoStmt) stmtNode()     {}
func (ReturnStmt) stmtNode() {}
func (BranchStmt) stmtNode() {}
func (BlockStmt) stmtNode()  {}
func (IfStmt) stmtNode()     {}
func (CaseClause) stmtNode() {}
func (SwitchStmt) stmtNode() {}
func (SelectStmt) stmtNode() {}
func (ForStmt) stmtNode()    {}
func (RangeStmt) stmtNode()  {}

func (n *ExprStmt) Accept(v Visitor) {
	v.VisitExprStmt(n)
}

func (n *SendStmt) Accept(v Visitor) {
	v.VisitSendStmt(n)
}

func (n *IncDecStmt) Accept(v Visitor) {
	v.VisitIncDecStmt(n)
}

func (n *AssignStmt) Accept(v Visitor) {
	v.VisitAssignStmt(n)
}

func (n *GoStmt) Accept(v Visitor) {
	v.VisitGoStmt(n)
}

func (n *ReturnStmt) Accept(v Visitor) {
	v.VisitReturnStmt(n)
}

func (n *BranchStmt) Accept(v Visitor) {
	v.VisitBranchStmt(n)
}

func (n *BlockStmt) Accept(v Visitor) {
	v.VisitBlockStmt(n)
}

func (n *IfStmt) Accept(v Visitor) {
	v.VisitIfStmt(n)
}

func (n *CaseClause) Accept(v Visitor) {
	v.VisitCaseClause(n)
}

func (n *SwitchStmt) Accept(v Visitor) {
	v.VisitSwitchStmt(n)
}

func (n *SelectStmt) Accept(v Visitor) {
	v.VisitSelectStmt(n)
}

func (n *ForStmt) Accept(v Visitor) {
	v.VisitForStmt(n)
}

func (n *RangeStmt) Accept(v Visitor) {
	v.VisitRangeStmt(n)
}
