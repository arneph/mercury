package logic

import (
	"fmt"
	positions "go/token"
)

type Test struct {
	name      string
	pos       positions.Pos
	Component *Component
	Steps     []TestStep
}

func NewTest(name string, pos positions.Pos, component *Component) *Test {
	return &Test{
		name:      name,
		pos:       pos,
		Component: component,
	}
}

type TestStep interface {
	Pos() positions.Pos
}

func (t *Test) Name() string {
	return t.name
}

type SetInputs struct {
	pos    positions.Pos
	Inputs *Constants
}

func NewSetInputsStep(pos positions.Pos, inputs *Constants) *SetInputs {
	return &SetInputs{
		pos:    pos,
		Inputs: inputs,
	}
}

func (s *SetInputs) Pos() positions.Pos {
	return s.pos
}

type CheckKind int

const (
	ASSERT CheckKind = iota
	EXPECT
)

func (k CheckKind) String() string {
	switch k {
	case ASSERT:
		return "assert"
	case EXPECT:
		return "expect"
	default:
		panic(fmt.Errorf("unexpected logic.CheckKind: %d", k))
	}
}

type CheckOutputs struct {
	pos     positions.Pos
	Kind    CheckKind
	Outputs *Constants
}

func NewCheckOutputsStep(pos positions.Pos, kind CheckKind, outputs *Constants) *CheckOutputs {
	return &CheckOutputs{
		pos:     pos,
		Kind:    kind,
		Outputs: outputs,
	}
}

func (c *CheckOutputs) Pos() positions.Pos {
	return c.pos
}
