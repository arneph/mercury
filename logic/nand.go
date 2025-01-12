package logic

type NandGate struct{}

var Nand = NandGate{}

func (g NandGate) Name() string {
	return "nand"
}

func (g NandGate) InputNames() []string {
	return []string{"a", "b"}
}

func (g NandGate) OutputNames() []string {
	return []string{"r"}
}
