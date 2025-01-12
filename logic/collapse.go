package logic

import (
	"fmt"
	"strconv"
)

func (c *Component) Collapse(newName string) *Component {
	cl := &collapser{
		result: &Component{
			name:      newName,
			Buses:     make(map[string]*Bus),
			Instances: nil,
		},
		componentInstanceCounts: make(map[string]int),
		wireLookup:              make(map[collapsedInstanceWire]BusWire),
	}
	instanceName := cl.defineInstance(c.Name())
	cl.resultInstanceName = instanceName
	cl.collapseBuses(instanceName, c.Buses, make(map[string]struct{}))
	for _, childInstance := range c.Instances {
		inputWires, outputWires := cl.ioWiresForInstance(childInstance, instanceName)
		cl.collapseInstance(childInstance, inputWires, outputWires)
	}
	for _, name := range c.InputBusNames {
		bus := c.Buses[name]
		for i := 0; i < bus.Wires(); i++ {
			wire := cl.lookupBusWire(instanceName, BusWire{
				Bus:       bus,
				WireIndex: WireIndex(i),
			})
			cl.result.InputBusNames = append(cl.result.InputBusNames, wire.Bus.Name)
		}
	}
	for _, name := range c.OutputBusNames {
		bus := c.Buses[name]
		for i := 0; i < bus.Wires(); i++ {
			wire := cl.lookupBusWire(instanceName, BusWire{
				Bus:       bus,
				WireIndex: WireIndex(i),
			})
			cl.result.OutputBusNames = append(cl.result.OutputBusNames, wire.Bus.Name)
		}
	}
	return cl.result
}

type collapsedInstanceName string
type collapsedInstanceWire struct {
	instanceName collapsedInstanceName
	busName      string
	wireIndex    WireIndex
}

type collapser struct {
	result                  *Component
	resultInstanceName      collapsedInstanceName
	componentInstanceCounts map[string]int
	wireLookup              map[collapsedInstanceWire]BusWire
}

func (c *collapser) defineInstance(name string) collapsedInstanceName {
	c.componentInstanceCounts[name]++
	i := c.componentInstanceCounts[name]
	return collapsedInstanceName(fmt.Sprintf("%s_i%d", name, i))
}

func (c *collapser) lookupWire(wire collapsedInstanceWire) BusWire {
	busWire, ok := c.wireLookup[wire]
	if !ok {
		panic(fmt.Errorf("undefined collapsed instance wire: %v", wire))
	}
	return busWire
}

func (c *collapser) lookupBusWire(instanceName collapsedInstanceName, busWire BusWire) BusWire {
	return c.lookupWire(collapsedInstanceWire{
		instanceName: instanceName,
		busName:      busWire.Bus.Name,
		wireIndex:    busWire.WireIndex,
	})
}

func (c *collapser) rememberWire(wire collapsedInstanceWire, busWire BusWire) {
	if old, ok := c.wireLookup[wire]; ok {
		panic(fmt.Errorf("already defined collapsed instance wire: old %v, new %v", old, wire))
	}
	c.result.Buses[busWire.Bus.Name] = busWire.Bus
	c.wireLookup[wire] = busWire
}

func (c *collapser) collapseBus(instanceName collapsedInstanceName, bus *Bus) {
	if bus.Wires() == 1 {
		var newBusName string
		if instanceName == c.resultInstanceName {
			newBusName = bus.Name
		} else {
			newBusName = string(instanceName) + "_" + bus.Name
		}
		c.rememberWire(collapsedInstanceWire{
			instanceName: instanceName,
			busName:      bus.Name,
			wireIndex:    0,
		}, BusWire{
			Bus:       NewBus(newBusName, 1),
			WireIndex: 0,
		})
	} else {
		for i := 0; i < bus.Wires(); i++ {
			var newBusName string
			if instanceName == c.resultInstanceName {
				newBusName = bus.Name + strconv.Itoa(i)
			} else {
				newBusName = string(instanceName) + "_" + bus.Name + strconv.Itoa(i)
			}
			c.rememberWire(collapsedInstanceWire{
				instanceName: instanceName,
				busName:      bus.Name,
				wireIndex:    WireIndex(i),
			}, BusWire{
				Bus:       NewBus(newBusName, 1),
				WireIndex: 0,
			})
		}
	}
}

func (c *collapser) collapseBuses(instanceName collapsedInstanceName, buses map[string]*Bus, exclude map[string]struct{}) {
	for name, bus := range buses {
		if _, ok := exclude[name]; ok {
			continue
		}
		c.collapseBus(instanceName, bus)
	}
}

func (c *collapser) ioWiresForInstance(instance *Instance, parentInstanceName collapsedInstanceName) ([]BusWire, []BusWire) {
	inputWires := make([]BusWire, len(instance.Inputs))
	for i, input := range instance.Inputs {
		inputWires[i] = c.lookupBusWire(parentInstanceName, input)
	}
	outputWires := make([]BusWire, len(instance.Outputs))
	for i, output := range instance.Outputs {
		outputWires[i] = c.lookupBusWire(parentInstanceName, output)
	}
	return inputWires, outputWires
}

func (c *collapser) collapseInstance(instance *Instance, inputWires []BusWire, outputWires []BusWire) {
	switch def := instance.Definition.(type) {
	case *Constants:
		c.result.Instances = append(c.result.Instances, &Instance{
			Definition: NewConstants(def.values),
			Inputs:     nil,
			Outputs:    outputWires,
		})
	case NandGate:
		c.result.Instances = append(c.result.Instances, &Instance{
			Definition: Nand,
			Inputs:     inputWires,
			Outputs:    outputWires,
		})
	case *Component:
		instanceName := c.defineInstance(def.name)
		ioBusNames := make(map[string]struct{})
		inputWireIndex := 0
		for _, busName := range def.InputBusNames {
			ioBusNames[busName] = struct{}{}
			bus := def.Buses[busName]
			for i := 0; i < bus.Wires(); i++ {
				c.rememberWire(collapsedInstanceWire{
					instanceName: instanceName,
					busName:      busName,
					wireIndex:    WireIndex(i),
				}, inputWires[inputWireIndex])
				inputWireIndex++
			}
		}
		outputWireIndex := 0
		for _, busName := range def.OutputBusNames {
			ioBusNames[busName] = struct{}{}
			bus := def.Buses[busName]
			for i := 0; i < bus.Wires(); i++ {
				c.rememberWire(collapsedInstanceWire{
					instanceName: instanceName,
					busName:      busName,
					wireIndex:    WireIndex(i),
				}, outputWires[outputWireIndex])
				outputWireIndex++
			}
		}
		c.collapseBuses(instanceName, def.Buses, ioBusNames)
		for _, childInstance := range def.Instances {
			childInputWires, childOutputWires := c.ioWiresForInstance(childInstance, instanceName)
			c.collapseInstance(childInstance, childInputWires, childOutputWires)
		}
	default:
		panic(fmt.Errorf("unexpected logic.Definition: %t", def))
	}
}
