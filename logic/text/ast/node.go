package ast

import positions "go/token"

type Node interface {
	Pos() positions.Pos
	End() positions.Pos
}
