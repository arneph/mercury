package logic

import (
	"fmt"
	"strconv"
	"strings"
)

type Instance struct {
	Definition Definition
	Inputs     []BusWire
	Outputs    []BusWire
}

func (inst *Instance) String() string {
	var sb strings.Builder
	for i, output := range inst.Outputs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(output.Bus.Name)
		if output.Bus.Wires() > 1 {
			sb.WriteString("[")
			sb.WriteString(strconv.Itoa(int(output.WireIndex)))
			sb.WriteString("]")
		}
	}
	sb.WriteString(": ")
	switch def := inst.Definition.(type) {
	case *Constants:
		sb.WriteString(def.String())
	case NandGate, *Component:
		sb.WriteString(def.Name())
		sb.WriteString("(")
		for i, input := range inst.Inputs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(input.Bus.Name)
			if input.Bus.Wires() > 1 {
				sb.WriteString("[")
				sb.WriteString(strconv.Itoa(int(input.WireIndex)))
				sb.WriteString("]")
			}
		}
		sb.WriteString(")")
	default:
		panic(fmt.Errorf("unexpected Definition: %t", def))
	}
	return sb.String()
}
