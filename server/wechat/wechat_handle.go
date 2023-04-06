package wechat

import (
	"fmt"
	"github.com/abigtomato/chatbot-go/config"
	"github.com/abigtomato/chatbot-go/openai/chatgpt"
	"github.com/abigtomato/chatbot-go/openai/dalle"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	DrawPrefix     = "画画"
	HelpDrawPrefix = "帮我画画"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	oa := wechat.NewWechat().GetOfficialAccount(&offConfig.Config{
		AppID:          config.C().Wechat.AppId,
		AppSecret:      config.C().Wechat.AppSecret,
		Token:          config.C().Wechat.Token,
		EncodingAESKey: config.C().Wechat.EncodingAESKey,
		Cache:          cache.NewMemory(),
	})
	server := oa.GetServer(r, w)

	// TODO 微信公众号请求超时5秒，会重试3次，需要解决
	server.SetMessageHandler(func(mixMessage *message.MixMessage) *message.Reply {
		answer := strings.TrimSpace(mixMessage.Content)
		logrus.Infoln("receive user message: ", answer)

		if strings.HasPrefix(answer, DrawPrefix) {
			answer = strings.TrimSpace(strings.TrimPrefix(answer, DrawPrefix))
			filename, err := dalle.Generations(answer, dalle.CreateImageResponseFormatB64JSON)
			if err != nil {
				return textMessage(fmt.Sprintf("dalle generations error: %v\n", err))
			}
			media, err := oa.GetMaterial().MediaUpload(material.MediaTypeImage, filename)
			if err != nil {
				return textMessage(fmt.Sprintf("media upload eror: %v\n", err))
			}
			return imageMessage(media.MediaID)
		}

		if strings.HasPrefix(answer, HelpDrawPrefix) {
			// TODO 将用户描述提交给chatGPT，让其代理生成dall-E的prompt
		}

		reply, err := chatgpt.Completions(answer)
		if err != nil {
			return textMessage(fmt.Sprintf("chatgpt completions error: %v\n", err))
		}
		return textMessage(reply)
	})

	if err := server.Serve(); err != nil {
		logrus.Errorln("message serve error: ", err)
		return
	}

	if err := server.Send(); err != nil {
		logrus.Errorln("message send error: ", err)
		return
	}
}

func textMessage(content string) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeText,
		MsgData: message.NewText(content),
	}
}

func imageMessage(mediaID string) *message.Reply {
	return &message.Reply{
		MsgType: message.MsgTypeImage,
		MsgData: message.NewImage(mediaID),
	}
}
