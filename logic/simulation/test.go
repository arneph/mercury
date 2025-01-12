package simulation

import (
	"fmt"

	"github.com/arneph/mercury/logic"

	errors "go/scanner"
	positions "go/token"
)

func RunTest(test *logic.Test, posFile *positions.File) (errs errors.ErrorList) {
	c := test.Component.Collapse(test.Component.Name())
	state := NewComponentState(c)
	for _, step := range test.Steps {
		switch step := step.(type) {
		case *logic.SetInputs:
			state.SetInputs(collapseInput(step.Inputs))
		case *logic.CheckOutputs:
			expected := step.Outputs
			actual := groupOutput(state.Outputs(), test.Component)
			match := true
			for i, as := range actual.Values() {
				es := expected.Values()[i]
				for j, a := range as {
					e := es[j]
					if a != e {
						match = false
					}
				}
			}
			if match {
				continue
			}
			errs.Add(posFile.Position(step.Pos()), fmt.Sprintf("%v failed: expected %v, got %v", step.Kind, expected, actual))
			if step.Kind == logic.ASSERT {
				return
			}
		default:
			panic(fmt.Errorf("unexpected logic.TestStep: %t", step))
		}
	}
	return
}

func collapseInput(input *logic.Constants) *logic.Constants {
	var vals []logic.Value
	for _, val := range input.Values() {
		for _, v := range val {
			vals = append(vals, logic.Value{v})
		}
	}
	return logic.NewConstants(vals)
}

func groupOutput(output *logic.Constants, c *logic.Component) *logic.Constants {
	vals := make([]logic.Value, len(c.OutputBusNames))
	outputIndex := 0
	for i, name := range c.OutputBusNames {
		bus := c.Buses[name]
		val := make(logic.Value, bus.Wires())
		for j := 0; j < bus.Wires(); j++ {
			val[j] = output.Values()[outputIndex][0]
			outputIndex++
		}
		vals[i] = val
	}
	return logic.NewConstants(vals)
}
