package model

// ChatContext 聊天上下文
type ChatContext struct {
	Id      string // 本次聊天的ID
	WeNick  string // 微信昵称
	Message string // 原消息
	Prompt  string // 提示语
	Content string // 上下文内容
}
