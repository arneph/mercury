package ast

import positions "go/token"

type BusReferenceList struct {
	References []*BusReference
}

func (l *BusReferenceList) Pos() positions.Pos {
	return l.References[0].Pos()
}

func (l *BusReferenceList) End() positions.Pos {
	return l.References[len(l.References)-1].End()
}

type BusReference struct {
	Name      *Identifier
	LBrack    positions.Pos
	WireIndex Expr
	RBrack    positions.Pos
}

func (r *BusReference) Pos() positions.Pos {
	return r.Name.Pos()
}

func (r *BusReference) End() positions.Pos {
	if r.WireIndex == nil {
		return r.Name.End()
	} else {
		return r.RBrack + 1
	}
}
