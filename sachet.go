package sachet

type Provider interface {
	Send(message Message) error
}

type Message struct {
	To   []string
	From string
	Text string
	Type string
}
