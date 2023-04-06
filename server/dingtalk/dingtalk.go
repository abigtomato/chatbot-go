package dingtalk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ReceiveMsg struct {
	ConversationID string `json:"conversationId"`
	AtUsers        []struct {
		DingTalkID string `json:"dingtalkId"`
	} `json:"atUsers"`
	ChatBotUserID             string `json:"chatbotUserId"`
	MsgID                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SenderStaffId             string `json:"senderStaffId"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	ConversationType          string `json:"conversationType"`
	SenderID                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      Text   `json:"text"`
	RobotCode                 string `json:"robotCode"`
	MsgType                   string `json:"msgtype"`
}

const (
	TEXT     string = "text"
	MARKDOWN string = "markdown"
)

type TextMessage struct {
	MsgType string `json:"msgtype"`
	At      *At    `json:"at"`
	Text    *Text  `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

type MarkDownMessage struct {
	MsgType  string    `json:"msgtype"`
	At       *At       `json:"at"`
	MarkDown *MarkDown `json:"markdown"`
}

type MarkDown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtUserIds []string `json:"atUserIds"`
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

func (r *ReceiveMsg) GetSenderIdentifier() (uid string) {
	uid = r.SenderStaffId
	if uid == "" {
		uid = r.SenderNick
	}
	return
}

func (r *ReceiveMsg) Reply(msgType, msg string) error {
	atUser := r.SenderStaffId
	if atUser == "" {
		msg = fmt.Sprintf("%s\n\n@%s", msg, r.SenderNick)
	}

	var msgTmp any
	switch msgType {
	case TEXT:
		msgTmp = &TextMessage{Text: &Text{Content: msg}, MsgType: TEXT, At: &At{AtUserIds: []string{atUser}}}
	case MARKDOWN:
		msgTmp = &MarkDownMessage{MsgType: MARKDOWN, At: &At{AtUserIds: []string{atUser}}, MarkDown: &MarkDown{Title: "Markdown Type", Text: msg}}
	default:
		msgTmp = &TextMessage{Text: &Text{Content: msg}, MsgType: TEXT, At: &At{AtUserIds: []string{atUser}}}
	}

	data, err := json.Marshal(msgTmp)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, r.SessionWebhook, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return nil
}
