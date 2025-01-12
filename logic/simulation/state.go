package simulation

import "github.com/arneph/mercury/logic"

type ComponentState struct {
	Component *logic.Component
	BusStates map[*logic.Bus]logic.Value
}

func NewComponentState(c *logic.Component) *ComponentState {
	busStates := make(map[*logic.Bus]logic.Value, len(c.Buses))
	for _, bus := range c.Buses {
		busStates[bus] = make(logic.Value, bus.Wires())
	}
	s := &ComponentState{
		Component: c,
		BusStates: busStates,
	}
	s.simulateUntilStable()
	return s
}

func (s *ComponentState) Inputs() logic.Constants {
	vals := make([]logic.Value, len(s.Component.InputBusNames))
	for i, name := range s.Component.InputBusNames {
		bus := s.Component.Buses[name]
		vals[i] = s.BusStates[bus]
	}
	return *logic.NewConstants(vals)
}

func (s *ComponentState) SetInputs(c *logic.Constants) {
	for i, name := range s.Component.InputBusNames {
		bus := s.Component.Buses[name]
		s.BusStates[bus] = c.Values()[i]
	}
	s.simulateUntilStable()
}

func (s *ComponentState) Outputs() *logic.Constants {
	vals := make([]logic.Value, len(s.Component.OutputBusNames))
	for i, name := range s.Component.OutputBusNames {
		bus := s.Component.Buses[name]
		vals[i] = s.BusStates[bus]
	}
	return logic.NewConstants(vals)
}

func (s *ComponentState) simulateUntilStable() {
	for {
		stable := true
		for _, instance := range s.Component.Instances {
			_ = instance.Definition.(logic.NandGate)
			aWire := instance.Inputs[0]
			aState := s.BusStates[aWire.Bus][aWire.WireIndex]
			bWire := instance.Inputs[1]
			bState := s.BusStates[bWire.Bus][bWire.WireIndex]
			rWire := instance.Outputs[0]
			oldRState := s.BusStates[rWire.Bus][rWire.WireIndex]
			newRState := !(aState && bState)
			s.BusStates[rWire.Bus][rWire.WireIndex] = newRState
			if oldRState != newRState {
				stable = false
			}
		}
		if stable {
			break
		}
	}
}
