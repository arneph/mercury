package ast

import positions "go/token"

type Test struct {
	Test    positions.Pos
	Name    *Identifier
	LBrace  positions.Pos
	Entries []TestEntry
	RBrace  positions.Pos
}

func (d *Test) Pos() positions.Pos {
	return d.Test
}

func (d *Test) End() positions.Pos {
	return d.RBrace + 1
}

func (d *Test) fileNode() {}

type TestEntry interface {
	Node
	testEntry()
}

type ComponentDecl struct {
	Component     positions.Pos
	Colon         positions.Pos
	ComponentName *Identifier
}

func (d *ComponentDecl) Pos() positions.Pos {
	return d.Component
}

func (d *ComponentDecl) End() positions.Pos {
	return d.ComponentName.End()
}

func (d *ComponentDecl) testEntry() {}

type SetInstr struct {
	Set       positions.Pos
	Inputs    *BusReferenceList
	Colon     positions.Pos
	Constants *Constants
}

func (s *SetInstr) Pos() positions.Pos {
	return s.Set
}

func (s *SetInstr) End() positions.Pos {
	return s.Constants.End()
}

func (s *SetInstr) testEntry() {}

type Assertion struct {
	Assert    positions.Pos
	Outputs   *BusReferenceList
	Is        positions.Pos
	Constants *Constants
}

func (a *Assertion) Pos() positions.Pos {
	return a.Assert
}

func (a *Assertion) End() positions.Pos {
	return a.Constants.End()
}

func (a *Assertion) testEntry() {}

type Expectation struct {
	Expect    positions.Pos
	Outputs   *BusReferenceList
	Is        positions.Pos
	Constants *Constants
}

func (e *Expectation) Pos() positions.Pos {
	return e.Expect
}

func (e *Expectation) End() positions.Pos {
	return e.Constants.End()
}

func (e *Expectation) testEntry() {}
