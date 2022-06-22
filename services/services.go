package servers

type Service interface {
	OnServiceMessage()
	OnUserMessage()
}
