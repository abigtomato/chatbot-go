package chatgpt

import (
	"github.com/abigtomato/chatbot-go/config"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestChatGPT_CompletionsWithConfig(t *testing.T) {
	content := "你好，自我介绍一下"
	reply, err := CompletionsWithConfig(content, config.LoadConfigWithPath("../../config.dev.yaml"))
	if err != nil {
		logrus.Errorln(err)
		return
	}
	logrus.Infoln(reply)
}
