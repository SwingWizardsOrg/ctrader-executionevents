package messagehandler

type Trader struct {
	name string
}

func (m *Trader) setName(name string) {
	m.name = name
}

func (m *Trader) getName() string {
	return m.name
}
