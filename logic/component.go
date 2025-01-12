package logic

import (
	"fmt"
	"strconv"
	"strings"
)

type Component struct {
	name           string
	Buses          map[string]*Bus
	InputBusNames  []string
	OutputBusNames []string
	Instances      []*Instance
}

func NewComponent(name string, inputs, outputs []*Bus) *Component {
	buses := make(map[string]*Bus)
	inputBusNames := make([]string, 0, len(inputs))
	for _, input := range inputs {
		if _, ok := buses[input.Name]; ok {
			panic(fmt.Errorf("bus name repeats: %q", input.Name))
		}
		buses[input.Name] = input
		inputBusNames = append(inputBusNames, input.Name)
	}
	outputBusNames := make([]string, 0, len(outputs))
	for _, output := range outputs {
		if _, ok := buses[output.Name]; ok {
			panic(fmt.Errorf("bus name repeats: %q", output.Name))
		}
		buses[output.Name] = output
		outputBusNames = append(outputBusNames, output.Name)
	}
	return &Component{
		name:           name,
		Buses:          buses,
		InputBusNames:  inputBusNames,
		OutputBusNames: outputBusNames,
	}
}

func (c *Component) Name() string {
	return c.name
}

func (c *Component) InputNames() []string {
	return c.InputBusNames
}

func (c *Component) InputWires() int {
	wires := 0
	for _, name := range c.InputBusNames {
		bus := c.Buses[name]
		wires += bus.Wires()
	}
	return wires
}

func (c *Component) OutputNames() []string {
	return c.OutputBusNames
}

func (c *Component) OutputWires() int {
	wires := 0
	for _, name := range c.OutputBusNames {
		bus := c.Buses[name]
		wires += bus.Wires()
	}
	return wires
}

func (c *Component) String() string {
	var sb strings.Builder
	sb.WriteString("component ")
	sb.WriteString(c.Name())
	sb.WriteString("(")
	ioBusNames := make(map[string]struct{})
	for i, name := range c.InputBusNames {
		ioBusNames[name] = struct{}{}
		if i > 0 {
			sb.WriteString(", ")
		}
		bus := c.Buses[name]
		sb.WriteString(name)
		if bus.Wires() > 1 {
			sb.WriteString("[")
			sb.WriteString(strconv.Itoa(bus.Wires()))
			sb.WriteString("]")
		}
	}
	sb.WriteString(")(")
	for i, name := range c.OutputBusNames {
		ioBusNames[name] = struct{}{}
		if i > 0 {
			sb.WriteString(", ")
		}
		bus := c.Buses[name]
		sb.WriteString(name)
		if bus.Wires() > 1 {
			sb.WriteString("[")
			sb.WriteString(strconv.Itoa(bus.Wires()))
			sb.WriteString("]")
		}
	}
	sb.WriteString(") {\n")
	for name, bus := range c.Buses {
		if bus.Wires() == 1 {
			continue
		} else if _, ok := ioBusNames[name]; ok {
			continue
		}
		sb.WriteString(indent)
		sb.WriteString("define ")
		sb.WriteString(name)
		if bus.Wires() > 1 {
			sb.WriteString("[")
			sb.WriteString(strconv.Itoa(bus.Wires()))
			sb.WriteString("]")
		}
		sb.WriteString("\n")
	}
	for _, instance := range c.Instances {
		sb.WriteString(indent)
		sb.WriteString(instance.String())
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}
