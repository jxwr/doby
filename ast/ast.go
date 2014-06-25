package ast

import (
	"github.com/jxwr/doubi/token"
)

type Node interface {
	//	Pos() int
	//	End() int
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
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
	Sel Ident
}

type IndexExpr struct {
	X      Expr
	Lbrack token.Pos
	Index  Expr
	Bbrack token.Pos
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

func (n Ident) exprNode()        {}
func (n BasicLit) exprNode()     {}
func (n ParenExpr) exprNode()    {}
func (n SelectorExpr) exprNode() {}
func (n IndexExpr) exprNode()    {}
func (n SliceExpr) exprNode()    {}
func (n CallExpr) exprNode()     {}
func (n UnaryExpr) exprNode()    {}
func (n BinaryExpr) exprNode()   {}
