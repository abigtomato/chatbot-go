package store

type Store interface {
	SetChatContext(userId string, content string)
	GetChatContext(userId string) string
}
