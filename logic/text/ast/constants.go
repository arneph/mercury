package ast

import positions "go/token"

type Constants struct {
	Values []*Number
}

func (c *Constants) Pos() positions.Pos {
	return c.Values[0].Pos()
}

func (c *Constants) End() positions.Pos {
	return c.Values[len(c.Values)-1].End()
}
