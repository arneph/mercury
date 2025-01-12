package text

import (
	"fmt"
	errors "go/scanner"
	positions "go/token"
	"math"
	"strconv"

	"github.com/arneph/mercury/logic"
	"github.com/arneph/mercury/logic/text/ast"
	"github.com/arneph/mercury/logic/text/parse"
	"github.com/arneph/mercury/logic/text/tokens"
)

func BuildFromFile(posFile *positions.File, src []byte) (*logic.System, errors.ErrorList) {
	astFile, errs := parse.ParseFile(posFile, src)
	if errs.Len() > 0 {
		return nil, errs
	}
	b := &builder{
		posFile: posFile,
		astFile: astFile,
		system:  logic.NewSystem(),
	}
	for _, astFileNode := range astFile.Nodes {
		b.buildFileNodeDeclaration(astFileNode)
	}
	for _, astFileNode := range astFile.Nodes {
		b.buildFileNodeDefinition(astFileNode)
	}
	return b.system, b.errs
}

type builder struct {
	posFile *positions.File
	astFile *ast.File
	errs    errors.ErrorList
	system  *logic.System
}

func (b *builder) buildFileNodeDeclaration(astFileNode ast.FileNode) {
	switch astFileNode := astFileNode.(type) {
	case *ast.Component:
		component := b.buildComponentDeclaration(astFileNode)
		if component != nil {
			b.system.Components[component.Name()] = component
		}
	case *ast.Test:
		break
	default:
		b.errs.Add(b.posFile.Position(astFileNode.Pos()), fmt.Sprintf("unexpected ast.FileNode: %v", astFileNode))
	}
}

func (b *builder) buildFileNodeDefinition(astFileNode ast.FileNode) {
	switch astFileNode := astFileNode.(type) {
	case *ast.Component:
		component := b.system.Components[astFileNode.Name.Name]
		b.buildComponentInstances(astFileNode, component)
	case *ast.Test:
		test := b.buildTest(astFileNode)
		if test != nil {
			b.system.Tests[test.Name()] = test
		}
	default:
		b.errs.Add(b.posFile.Position(astFileNode.Pos()), fmt.Sprintf("unexpected ast.FileNode: %v", astFileNode))
	}
}

func (b *builder) buildComponentDeclaration(astComponent *ast.Component) *logic.Component {
	name := astComponent.Name.Name
	if _, ok := b.system.Components[name]; ok {
		b.errs.Add(b.posFile.Position(astComponent.Name.Pos()), fmt.Sprintf("redefinition of component: %s", name))
		return nil
	}
	cb := b.newComponentBuilder()
	var inputs, outputs []*logic.Bus
	for _, astBus := range astComponent.Inputs.Defintions {
		input := cb.addBus(astBus)
		if input != nil {
			inputs = append(inputs, input)
		}
	}
	for _, astBus := range astComponent.Outputs.Defintions {
		output := cb.addBus(astBus)
		if output != nil {
			outputs = append(outputs, output)
		}
	}
	return logic.NewComponent(name, inputs, outputs)
}

func (b *builder) buildComponentInstances(astComponent *ast.Component, c *logic.Component) {
	cb := b.newComponentBuilderForComponent(c)
	for _, astEntry := range astComponent.Entries {
		switch astEntry := astEntry.(type) {
		case *ast.BusDefinitionEntry:
			cb.buildBusDefinitionEntry(astEntry)
		case *ast.ForLoop:
			cis := cb.buildForLoop(astEntry)
			c.Instances = append(c.Instances, cis...)
		case *ast.ComponentInstance:
			ci := cb.buildComponentInstance(astEntry)
			if ci != nil {
				c.Instances = append(c.Instances, ci)
			}
		default:
			b.errs.Add(b.posFile.Position(astEntry.Pos()), fmt.Sprintf("Unexpected ast.ComponentEntry: %v", astEntry))
		}
	}
}

func (b *builder) newComponentBuilder() *componentBuilder {
	return &componentBuilder{
		builder: b,
		buses:   make(map[string]*logic.Bus),
		vars:    make(map[string]int),
	}
}

func (b *builder) newComponentBuilderForComponent(c *logic.Component) *componentBuilder {
	return &componentBuilder{
		builder: b,
		buses:   c.Buses,
		vars:    make(map[string]int),
	}
}

type componentBuilder struct {
	*builder
	buses map[string]*logic.Bus
	vars  map[string]int
}

func (b *componentBuilder) addBus(astBus *ast.BusDefinition) *logic.Bus {
	name := astBus.Name.Name
	if _, ok := b.buses[name]; ok {
		b.errs.Add(b.posFile.Position(astBus.Name.Pos()), fmt.Sprintf("redefinition of bus with name: %s", name))
		return nil
	} else if _, ok := b.vars[name]; ok {
		b.errs.Add(b.posFile.Position(astBus.Name.Pos()), fmt.Sprintf("redefinition of variable with name: %s", name))
		return nil
	}
	wireCount := 1
	if astBus.WireCount != nil {
		var ok bool
		wireCount, ok = b.evalInt(astBus.WireCount)
		if !ok {
			return nil
		}
	}
	bus := logic.NewBus(name, wireCount)
	b.buses[name] = bus
	return bus
}

type busReferenceMode int

const (
	DEFINITION_NOT_ALLOWED busReferenceMode = iota
	DEFINITION_ALLOWED
)

func (b *componentBuilder) buildBusReferenceList(astBusReferneceList *ast.BusReferenceList, mode busReferenceMode) []logic.BusWire {
	var wires []logic.BusWire
	for _, astBusReference := range astBusReferneceList.References {
		ws := b.buildBusReference(astBusReference, mode)
		if ws == nil {
			return nil
		}
		wires = append(wires, ws...)
	}
	return wires
}

func (b *componentBuilder) buildBusReference(astBusReference *ast.BusReference, mode busReferenceMode) []logic.BusWire {
	name := astBusReference.Name.Name
	bus, ok := b.buses[name]
	if !ok && mode != DEFINITION_ALLOWED {
		b.errs.Add(b.posFile.Position(astBusReference.Pos()), fmt.Sprintf("bus is undefined: %s", name))
		return nil
	} else if _, ok2 := b.vars[name]; !ok && ok2 {
		b.errs.Add(b.posFile.Position(astBusReference.Pos()), fmt.Sprintf("redefinition of variable with name: %s", name))
		return nil
	} else if !ok {
		if astBusReference.WireIndex != nil {
			b.errs.Add(b.posFile.Position(astBusReference.Pos()), fmt.Sprintf("cannot define bus with wire index: %s", name))
			return nil
		}
		bus = logic.NewBus(name, 1)
		b.buses[name] = bus
	}
	if astBusReference.WireIndex == nil {
		wires := make([]logic.BusWire, 0, bus.Wires())
		for i := 0; i < bus.Wires(); i++ {
			wires = append(wires, logic.BusWire{
				Bus:       bus,
				WireIndex: logic.WireIndex(i),
			})
		}
		return wires
	} else {
		index, ok := b.evalExpr(astBusReference.WireIndex)
		if !ok {
			return nil
		}
		return []logic.BusWire{
			{
				Bus:       bus,
				WireIndex: logic.WireIndex(index),
			},
		}
	}
}

func (b *componentBuilder) evalExpr(astExpr ast.Expr) (int, bool) {
	switch astExpr := astExpr.(type) {
	case *ast.UnaryExpr:
		return b.evalUnaryExpr(astExpr)
	case *ast.BinaryExpr:
		return b.evalBinaryExpr(astExpr)
	case *ast.Identifier:
		return b.evalIdentifier(astExpr)
	case *ast.Number:
		return b.evalInt(astExpr)
	default:
		b.errs.Add(b.posFile.Position(astExpr.Pos()), fmt.Sprintf("unexpected ast.Expr: %v", astExpr))
		return 0, false
	}
}

func (b *componentBuilder) evalUnaryExpr(astUnaryExpr *ast.UnaryExpr) (int, bool) {
	op, ok := b.evalExpr(astUnaryExpr.Operand)
	if !ok {
		return 0, false
	}
	switch astUnaryExpr.Operator {
	case tokens.ADD:
		return +op, true
	case tokens.SUB:
		return -op, true
	default:
		b.errs.Add(b.posFile.Position(astUnaryExpr.OperatorStart), fmt.Sprintf("unkown unary operator: %d", astUnaryExpr.Operator))
		return 0, false
	}
}

func (b *componentBuilder) evalBinaryExpr(astBinaryExpr *ast.BinaryExpr) (int, bool) {
	lhs, ok := b.evalExpr(astBinaryExpr.LhsOperand)
	if !ok {
		return 0, false
	}
	rhs, ok := b.evalExpr(astBinaryExpr.RhsOperand)
	if !ok {
		return 0, false
	}
	switch astBinaryExpr.Operator {
	case tokens.ADD:
		return lhs + rhs, true
	case tokens.SUB:
		return lhs - rhs, true
	case tokens.MUL:
		return lhs * rhs, true
	case tokens.QUO:
		return lhs / rhs, true
	case tokens.REM:
		return lhs % rhs, true
	default:
		b.errs.Add(b.posFile.Position(astBinaryExpr.OperatorStart), fmt.Sprintf("unkown binary operator: %d", astBinaryExpr.Operator))
		return 0, false
	}
}

func (b *componentBuilder) evalIdentifier(astIdentifier *ast.Identifier) (int, bool) {
	i, ok := b.vars[astIdentifier.Name]
	if !ok {
		b.errs.Add(b.posFile.Position(astIdentifier.Pos()), fmt.Sprintf("variable is undefined: %s", astIdentifier.Name))
		return 0, false
	}
	return i, true
}

func (b *builder) evalInt(astNumber *ast.Number) (int, bool) {
	i, err := strconv.ParseUint(astNumber.Value, 0, 64)
	if err != nil {
		b.errs.Add(b.posFile.Position(astNumber.Pos()), fmt.Sprintf("could not convert %q to int: %v", astNumber.Value, err))
		return 0, false
	}
	return int(i), true
}

func (b *componentBuilder) buildBusDefinitionEntry(astBusDefinitionEntry *ast.BusDefinitionEntry) {
	for _, astBus := range astBusDefinitionEntry.Definitions.Defintions {
		b.addBus(astBus)
	}
}

func (b *componentBuilder) buildForLoop(astForLoop *ast.ForLoop) []*logic.Instance {
	varName := astForLoop.Variable.Name
	if _, ok := b.buses[varName]; ok {
		b.errs.Add(b.posFile.Position(astForLoop.Variable.Pos()), fmt.Sprintf("redefinition of bus with name: %s", varName))
		return nil
	} else if _, ok := b.vars[varName]; ok {
		b.errs.Add(b.posFile.Position(astForLoop.Variable.Pos()), fmt.Sprintf("redefinition of variable with name: %s", varName))
		return nil
	}
	first, ok := b.evalExpr(astForLoop.First)
	if !ok {
		return nil
	}
	last, ok := b.evalExpr(astForLoop.Last)
	if !ok {
		return nil
	}
	if first > last {
		b.errs.Add(b.posFile.Position(astForLoop.Pos()), fmt.Sprintf("first value is larger than last value: %d > %d", first, last))
		return nil
	}
	instances := make([]*logic.Instance, 0, (last-first+1)*len(astForLoop.Entries))
	for i := first; i <= last; i++ {
		b.vars[varName] = i
		for _, astEntry := range astForLoop.Entries {
			switch astEntry := astEntry.(type) {
			case *ast.BusDefinitionEntry:
				b.errs.Add(b.posFile.Position(astEntry.Pos()), "bus definition now allowed in loop")
			case *ast.ForLoop:
				cis := b.buildForLoop(astEntry)
				instances = append(instances, cis...)
			case *ast.ComponentInstance:
				ci := b.buildComponentInstance(astEntry)
				if ci != nil {
					instances = append(instances, ci)
				}
			default:
				b.errs.Add(b.posFile.Position(astEntry.Pos()), fmt.Sprintf("Unexpected ast.ComponentEntry: %v", astEntry))
			}
		}
	}
	delete(b.vars, varName)
	return instances
}

func (b *componentBuilder) buildComponentInstance(astComponentInstance *ast.ComponentInstance) *logic.Instance {
	componentName := astComponentInstance.DefinitionName.Name
	var def logic.Definition
	var expectedInputs, expectedOutputs int
	if componentName == "nand" {
		def = logic.Nand
		expectedInputs = 2
		expectedOutputs = 1
	} else {
		component, ok := b.system.Components[componentName]
		if !ok {
			b.errs.Add(b.posFile.Position(astComponentInstance.DefinitionName.Pos()), fmt.Sprintf("undefined component: %s", componentName))
			return nil
		}
		def = component
		expectedInputs = component.InputWires()
		expectedOutputs = component.OutputWires()
	}
	outputs := b.buildBusReferenceList(astComponentInstance.Outputs, DEFINITION_ALLOWED)
	if outputs == nil {
		return nil
	} else if len(outputs) != expectedOutputs {
		b.errs.Add(b.posFile.Position(astComponentInstance.Outputs.Pos()), fmt.Sprintf("wrong number of output wires: expected %d, got %d", expectedOutputs, len(outputs)))
		return nil
	}
	inputs := b.buildBusReferenceList(astComponentInstance.Inputs, DEFINITION_ALLOWED)
	if inputs == nil {
		return nil
	} else if len(inputs) != expectedInputs {
		b.errs.Add(b.posFile.Position(astComponentInstance.Inputs.Pos()), fmt.Sprintf("wrong number of input wires: expected %d, got %d", expectedInputs, len(inputs)))
		return nil
	}
	return &logic.Instance{
		Definition: def,
		Inputs:     inputs,
		Outputs:    outputs,
	}
}

func (b *builder) buildTest(astTest *ast.Test) *logic.Test {
	name := astTest.Name.Name
	if _, ok := b.system.Tests[name]; ok {
		b.errs.Add(b.posFile.Position(astTest.Name.Pos()), fmt.Sprintf("redefinition of test: %s", name))
		return nil
	}
	var component *logic.Component
	for _, astTestEntry := range astTest.Entries {
		switch astTestEntry := astTestEntry.(type) {
		case *ast.ComponentDecl:
			if component != nil {
				b.errs.Add(b.posFile.Position(astTestEntry.Pos()), "redeclaration of test component")
				return nil
			}
			componentName := astTestEntry.ComponentName.Name
			c, ok := b.system.Components[componentName]
			if !ok {
				b.errs.Add(b.posFile.Position(astTestEntry.ComponentName.Pos()), fmt.Sprintf("undefined component: %s", componentName))
				return nil
			}
			component = c
		default:
			break
		}
	}
	t := logic.NewTest(name, astTest.Pos(), component)
	tb := b.newTestBuilderForTest(t)
	for _, astTestEntry := range astTest.Entries {
		switch astTestEntry := astTestEntry.(type) {
		case *ast.ComponentDecl:
			break
		case *ast.SetInstr:
			s := tb.buildSetInstr(astTestEntry)
			if s != nil {
				t.Steps = append(t.Steps, s)
			}
		case *ast.Assertion:
			s := tb.buildAssertion(astTestEntry)
			if s != nil {
				t.Steps = append(t.Steps, s)
			}
		case *ast.Expectation:
			s := tb.buildExpectation(astTestEntry)
			if s != nil {
				t.Steps = append(t.Steps, s)
			}
		default:
			b.errs.Add(b.posFile.Position(astTestEntry.Pos()), fmt.Sprintf("unexpected ast.TestEntry: %v", astTestEntry))
		}
	}
	return t
}

func (b *builder) newTestBuilderForTest(t *logic.Test) *testBuilder {
	return &testBuilder{
		builder: b,
		test:    t,
	}
}

type testBuilder struct {
	*builder
	test *logic.Test
}

func (b *testBuilder) buildSetInstr(astSetInstr *ast.SetInstr) *logic.SetInputs {
	expectedBuses := len(b.test.Component.InputBusNames)
	actualBuses := len(astSetInstr.Inputs.References)
	if actualBuses > expectedBuses {
		b.errs.Add(b.posFile.Position(astSetInstr.Inputs.Pos()), fmt.Sprintf("too many input buses: expected %d, got %d", expectedBuses, actualBuses))
		return nil
	} else if actualBuses < expectedBuses {
		b.errs.Add(b.posFile.Position(astSetInstr.Inputs.Pos()), fmt.Sprintf("too few input buses: expected %d, got %d", expectedBuses, actualBuses))
		return nil
	}
	actualConstants := len(astSetInstr.Constants.Values)
	if actualConstants > expectedBuses {
		b.errs.Add(b.posFile.Position(astSetInstr.Constants.Pos()), fmt.Sprintf("too many input values: expected %d, got %d", expectedBuses, actualConstants))
		return nil
	} else if actualConstants < expectedBuses {
		b.errs.Add(b.posFile.Position(astSetInstr.Constants.Pos()), fmt.Sprintf("too few input values: expected %d, got %d", expectedBuses, actualConstants))
		return nil
	}
	var cs []logic.Value
	for i, input := range astSetInstr.Inputs.References {
		if input.WireIndex != nil {
			b.errs.Add(b.posFile.Position(input.LBrack), "input buses in tests must be set in full")
		}
		actualName := input.Name.Name
		expectedName := b.test.Component.InputBusNames[i]
		if actualName != expectedName {
			b.errs.Add(b.posFile.Position(input.Name.Pos()), fmt.Sprintf("incorrect input bus name: expected %s, got %s", expectedName, actualName))
		}
		value, ok := b.evalInt(astSetInstr.Constants.Values[i])
		if !ok {
			return nil
		}
		bus := b.test.Component.Buses[expectedName]
		expectedWires := bus.Wires()
		minActualWires := 1
		if value > 0 {
			minActualWires = 1 + int(math.Ceil(math.Log2(float64(value))))
		}
		if minActualWires > expectedWires {
			b.errs.Add(b.posFile.Position(astSetInstr.Constants.Values[i].Pos()), fmt.Sprintf("too many input wires: expected %d, got at least %d", expectedWires, minActualWires))
			return nil
		}
		c := make(logic.Value, expectedWires)
		for i := 0; i < expectedWires; i++ {
			c[i] = (value>>i)%2 == 1
		}
		cs = append(cs, c)
	}
	return logic.NewSetInputsStep(astSetInstr.Pos(), logic.NewConstants(cs))
}

func (b *testBuilder) buildAssertion(astAssertion *ast.Assertion) *logic.CheckOutputs {
	expectedBuses := len(b.test.Component.OutputBusNames)
	actualBuses := len(astAssertion.Outputs.References)
	if actualBuses > expectedBuses {
		b.errs.Add(b.posFile.Position(astAssertion.Outputs.Pos()), fmt.Sprintf("too many output buses: expected %d, got %d", expectedBuses, actualBuses))
		return nil
	} else if actualBuses < expectedBuses {
		b.errs.Add(b.posFile.Position(astAssertion.Outputs.Pos()), fmt.Sprintf("too few output buses: expected %d, got %d", expectedBuses, actualBuses))
		return nil
	}
	actualConstants := len(astAssertion.Constants.Values)
	if actualConstants > expectedBuses {
		b.errs.Add(b.posFile.Position(astAssertion.Constants.Pos()), fmt.Sprintf("too many output values: expected %d, got %d", expectedBuses, actualConstants))
		return nil
	} else if actualConstants < expectedBuses {
		b.errs.Add(b.posFile.Position(astAssertion.Constants.Pos()), fmt.Sprintf("too few output values: expected %d, got %d", expectedBuses, actualConstants))
		return nil
	}
	var cs []logic.Value
	for i, output := range astAssertion.Outputs.References {
		if output.WireIndex != nil {
			b.errs.Add(b.posFile.Position(output.LBrack), "output buses in tests must be set in full")
		}
		actualName := output.Name.Name
		expectedName := b.test.Component.OutputBusNames[i]
		if actualName != expectedName {
			b.errs.Add(b.posFile.Position(output.Name.Pos()), fmt.Sprintf("incorrect output bus name: expected %s, got %s", expectedName, actualName))
		}
		value, ok := b.evalInt(astAssertion.Constants.Values[i])
		if !ok {
			return nil
		}
		bus := b.test.Component.Buses[expectedName]
		expectedWires := bus.Wires()
		minActualWires := 1
		if value > 0 {
			minActualWires = 1 + int(math.Ceil(math.Log2(float64(value))))
		}
		if minActualWires > expectedWires {
			b.errs.Add(b.posFile.Position(astAssertion.Constants.Values[i].Pos()), fmt.Sprintf("too many output wires: expected %d, got at least %d", expectedWires, minActualWires))
			return nil
		}
		c := make(logic.Value, expectedWires)
		for i := 0; i < expectedWires; i++ {
			c[i] = (value>>i)%2 == 1
		}
		cs = append(cs, c)
	}
	return logic.NewCheckOutputsStep(astAssertion.Pos(), logic.ASSERT, logic.NewConstants(cs))
}

func (b *testBuilder) buildExpectation(astExpectation *ast.Expectation) *logic.CheckOutputs {
	expectedBuses := len(b.test.Component.OutputBusNames)
	actualBuses := len(astExpectation.Outputs.References)
	if actualBuses > expectedBuses {
		b.errs.Add(b.posFile.Position(astExpectation.Outputs.Pos()), fmt.Sprintf("too many output buses: expected %d, got %d", expectedBuses, actualBuses))
		return nil
	} else if actualBuses < expectedBuses {
		b.errs.Add(b.posFile.Position(astExpectation.Outputs.Pos()), fmt.Sprintf("too few output buses: expected %d, got %d", expectedBuses, actualBuses))
		return nil
	}
	actualConstants := len(astExpectation.Constants.Values)
	if actualConstants > expectedBuses {
		b.errs.Add(b.posFile.Position(astExpectation.Constants.Pos()), fmt.Sprintf("too many output values: expected %d, got %d", expectedBuses, actualConstants))
		return nil
	} else if actualConstants < expectedBuses {
		b.errs.Add(b.posFile.Position(astExpectation.Constants.Pos()), fmt.Sprintf("too few output values: expected %d, got %d", expectedBuses, actualConstants))
		return nil
	}
	var cs []logic.Value
	for i, output := range astExpectation.Outputs.References {
		if output.WireIndex != nil {
			b.errs.Add(b.posFile.Position(output.LBrack), "output buses in tests must be set in full")
		}
		actualName := output.Name.Name
		expectedName := b.test.Component.OutputBusNames[i]
		if actualName != expectedName {
			b.errs.Add(b.posFile.Position(output.Name.Pos()), fmt.Sprintf("incorrect output bus name: expected %s, got %s", expectedName, actualName))
		}
		value, ok := b.evalInt(astExpectation.Constants.Values[i])
		if !ok {
			return nil
		}
		bus := b.test.Component.Buses[expectedName]
		expectedWires := bus.Wires()
		minActualWires := 1
		if value > 0 {
			minActualWires = 1 + int(math.Ceil(math.Log2(float64(value))))
		}
		if minActualWires > expectedWires {
			b.errs.Add(b.posFile.Position(astExpectation.Constants.Values[i].Pos()), fmt.Sprintf("too many output wires: expected %d, got at least %d", expectedWires, minActualWires))
			return nil
		}
		c := make(logic.Value, expectedWires)
		for i := 0; i < expectedWires; i++ {
			c[i] = (value>>i)%2 == 1
		}
		cs = append(cs, c)
	}
	return logic.NewCheckOutputsStep(astExpectation.Pos(), logic.EXPECT, logic.NewConstants(cs))
}
