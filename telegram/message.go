package telegram

type Message struct {
	text string
}

func (m Message) Text() string {
	return m.text
}
