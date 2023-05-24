package messagehandler

type Message struct {
	name string
}

func (m *Message) setName(name string) {
	m.name = name
}

func (m *Message) getName() string {
	return m.name
}
