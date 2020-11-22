package telegram

type Message struct {
	text   string
	chatId int64
}

func (m Message) Text() string {
	return m.text
}

func (m Message) ChatId() int64 {
	return m.chatId
}
