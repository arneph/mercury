package logic

type Definition interface {
	Name() string
	InputNames() []string
	OutputNames() []string
}
