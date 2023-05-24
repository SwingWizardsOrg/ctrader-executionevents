package messagehandler

type Imessage interface {
	setName(name string)
	getName() string
}
