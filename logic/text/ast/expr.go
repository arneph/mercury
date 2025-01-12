package ast

import (
	positions "go/token"

	"github.com/arneph/mercury/logic/text/tokens"
)

type Expr interface {
	Node
	expr()
}

type UnaryExpr struct {
	Operator      tokens.Token
	OperatorStart positions.Pos
	Operand       Expr
}

func (u *UnaryExpr) Pos() positions.Pos {
	return u.OperatorStart
}

func (u *UnaryExpr) End() positions.Pos {
	return u.Operand.End()
}

func (u *UnaryExpr) expr() {}

type BinaryExpr struct {
	Operator      tokens.Token
	OperatorStart positions.Pos
	LhsOperand    Expr
	RhsOperand    Expr
}

func (b *BinaryExpr) Pos() positions.Pos {
	return b.LhsOperand.Pos()
}

func (b *BinaryExpr) End() positions.Pos {
	return b.RhsOperand.End()
}

func (b *BinaryExpr) expr() {}
