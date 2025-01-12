package logic

type Source interface{}
type Sink interface{}

type Bus struct {
	Name    string
	sources [][]Source
	sinks   [][]Sink
}

type WireIndex int

type BusWire struct {
	Bus       *Bus
	WireIndex WireIndex
}

func NewBus(name string, wireCount int) *Bus {
	return &Bus{
		Name:    name,
		sources: make([][]Source, wireCount),
		sinks:   make([][]Sink, wireCount),
	}
}

func (b Bus) Wires() int {
	return len(b.sources)
}

func (b Bus) Sources() [][]Source {
	return b.sources
}

func (b Bus) SourcesForWire(i WireIndex) []Source {
	return b.sources[i]
}

func (b Bus) AddSourceForWire(i WireIndex, source Source) {
	b.sources[i] = append(b.sources[i], source)
}

func (b Bus) Sinks() [][]Sink {
	return b.sinks
}

func (b Bus) SinksForWire(i WireIndex) []Sink {
	return b.sinks[i]
}

func (b Bus) AddSinkForWire(i WireIndex, sink Sink) {
	b.sinks[i] = append(b.sinks[i], sink)
}
