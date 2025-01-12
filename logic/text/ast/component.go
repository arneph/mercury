package ast

import (
	positions "go/token"
)

type Component struct {
	Component    positions.Pos
	Name         *Identifier
	InputLParen  positions.Pos
	Inputs       *BusDefinitionList
	InputRParen  positions.Pos
	OutputLParen positions.Pos
	Outputs      *BusDefinitionList
	OutputRParen positions.Pos
	LBrace       positions.Pos
	Entries      []ComponentEntry
	RBRace       positions.Pos
}

func (c *Component) Pos() positions.Pos {
	return c.Component
}

func (c *Component) End() positions.Pos {
	return c.RBRace + 1
}

func (c *Component) fileNode() {}

type ComponentEntry interface {
	Node
	componentEntry()
}

type BusDefinitionEntry struct {
	Define      positions.Pos
	Definitions *BusDefinitionList
}

func (d *BusDefinitionEntry) Pos() positions.Pos {
	return d.Define
}

func (d *BusDefinitionEntry) End() positions.Pos {
	return d.Definitions.End()
}

func (d *BusDefinitionEntry) componentEntry() {}

type ForLoop struct {
	For      positions.Pos
	Variable *Identifier
	From     positions.Pos
	First    Expr
	To       positions.Pos
	Last     Expr
	LBrace   positions.Pos
	Entries  []ComponentEntry
	RBrace   positions.Pos
}

func (l *ForLoop) Pos() positions.Pos {
	return l.For
}

func (l *ForLoop) End() positions.Pos {
	return l.RBrace + 1
}

func (l *ForLoop) componentEntry() {}

type Instance interface {
	ComponentEntry
	instance()
}

type ConstantsInstace struct {
	Outputs   *BusReferenceList
	Colon     positions.Pos
	Constants *Constants
}

func (c *ConstantsInstace) Pos() positions.Pos {
	return c.Outputs.Pos()
}

func (c *ConstantsInstace) End() positions.Pos {
	return c.Constants.End()
}

func (c *ConstantsInstace) componentEntry() {}
func (c *ConstantsInstace) instance()       {}

type ComponentInstance struct {
	Outputs        *BusReferenceList
	Colon          positions.Pos
	DefinitionName *Identifier
	LParen         positions.Pos
	Inputs         *BusReferenceList
	RParen         positions.Pos
}

func (c *ComponentInstance) Pos() positions.Pos {
	return c.Outputs.Pos()
}

func (c *ComponentInstance) End() positions.Pos {
	return c.RParen
}

func (c *ComponentInstance) componentEntry() {}
func (c *ComponentInstance) instance()       {}

type BusDefinitionList struct {
	Defintions []*BusDefinition
}

func (l *BusDefinitionList) Pos() positions.Pos {
	return l.Defintions[0].Pos()
}

func (l *BusDefinitionList) End() positions.Pos {
	return l.Defintions[len(l.Defintions)-1].End()
}

type BusDefinition struct {
	Name      *Identifier
	LBrack    positions.Pos
	WireCount *Number
	RBrack    positions.Pos
}

func (d *BusDefinition) Pos() positions.Pos {
	return d.Name.Pos()
}

func (d *BusDefinition) End() positions.Pos {
	if d.WireCount == nil {
		return d.Name.End()
	} else {
		return d.RBrack + 1
	}
}
