package logic

type System struct {
	Components map[string]*Component
	Tests      map[string]*Test
}

func NewSystem() *System {
	return &System{
		Components: make(map[string]*Component),
		Tests:      make(map[string]*Test),
	}
}

func (s *System) AddComponent(c *Component) {
	s.Components[c.name] = c
}

func (s *System) AddTest(t *Test) {
	s.Tests[t.name] = t
}
