package ast

import positions "go/token"

type Identifier struct {
	Name  string
	Start positions.Pos
}

func (i *Identifier) Pos() positions.Pos {
	return i.Start
}

func (i *Identifier) End() positions.Pos {
	return i.Start + positions.Pos(len(i.Name))
}

func (i *Identifier) expr() {}

type Number struct {
	Value string
	Start positions.Pos
}

func (n *Number) Pos() positions.Pos {
	return n.Start
}

func (n *Number) End() positions.Pos {
	return n.Start + positions.Pos(len(n.Value))
}

func (n *Number) expr() {}
