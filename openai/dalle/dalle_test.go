package dalle

import (
	"github.com/abigtomato/chatbot-go/config"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestDallE_GenerationsWithConfig(t *testing.T) {
	content := "晚霞"
	format := CreateImageResponseFormatB64JSON
	reply, err := GenerationsWithConfig(content, format, config.LoadConfigWithPath("../../config.dev.yaml"))
	if err != nil {
		logrus.Errorln(err)
		return
	}
	logrus.Infoln(reply)
}
