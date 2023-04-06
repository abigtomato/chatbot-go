package dingtalk

import (
	"encoding/json"
	"fmt"
	"github.com/abigtomato/chatbot-go/openai/chatgpt"
	"github.com/abigtomato/chatbot-go/openai/dalle"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Errorln("read dingTalk request error: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(data) == 0 {
		logrus.Warning("dingTalk callback param is empty")
		return
	}

	receiveMsg := &ReceiveMsg{}
	if err := json.Unmarshal(data, receiveMsg); err != nil {
		logrus.Errorln("json unmarshal error: ", err)
		return
	}

	if receiveMsg.Text.Content == "" || receiveMsg.ChatBotUserID == "" {
		logrus.Errorln("receive message content is empty")
		return
	}

	logrus.Infoln("receive dingTalk message: ", receiveMsg.Text.Content)

	content := strings.TrimSpace(receiveMsg.Text.Content)
	switch {
	case strings.HasPrefix(content, "/画画"):
		content = strings.TrimSpace(strings.TrimPrefix(content, "/画画"))
		path, err := dalle.Generations(content, dalle.CreateImageResponseFormatB64JSON)
		if err != nil {
			logrus.Errorln("dall-e generations error: ", err)
			return
		}
		if err := receiveMsg.Reply(MARKDOWN, fmt.Sprintf("![](%s)", path)); err != nil {
			logrus.Errorln("receiveMsg reply error: ", err)
			return
		}
		return
	case strings.HasPrefix(content, "/帮我画"):
		content = strings.TrimSpace(strings.TrimPrefix(content, "/帮我画"))
		content = fmt.Sprintf("请根据“%s”生成dall-e的prompt，直接回复prompt，一条", content)
		reply, err := chatgpt.Completions(content)
		if err != nil {
			logrus.Errorln("chatGPT completions error: ", err)
			return
		}
		path, err := dalle.Generations(reply, dalle.CreateImageResponseFormatB64JSON)
		if err != nil {
			logrus.Errorln("dall-e generations error: ", err)
			return
		}
		if err := receiveMsg.Reply(MARKDOWN, fmt.Sprintf("![](%s)", path)); err != nil {
			logrus.Errorln("receiveMsg reply error: ", err)
			return
		}
		return
	default:
		reply, err := chatgpt.Completions(content)
		if err != nil {
			logrus.Errorln("chatGPT completions error: ", err)
			return
		}
		if err := receiveMsg.Reply(TEXT, reply); err != nil {
			logrus.Errorln("receiveMsg reply error: ", err)
			return
		}
		return
	}
}
