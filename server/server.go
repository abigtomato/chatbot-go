package server

import (
	"github.com/abigtomato/chatbot-go/openai/chatgpt"
	"github.com/abigtomato/chatbot-go/openai/dalle"
	"github.com/abigtomato/chatbot-go/server/dingtalk"
	"github.com/abigtomato/chatbot-go/server/wechat"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RequestBody struct {
	Message string `json:"message"`
}

func Run() {
	engine := gin.Default()
	engine.Use(cors.Default())

	engine.Any("/wechat", func(context *gin.Context) {
		wechat.Handle(context.Writer, context.Request)
	})

	engine.Any("/dingtalk", func(context *gin.Context) {
		dingtalk.Handle(context.Writer, context.Request)
	})

	engine.GET("/images/:filename", func(context *gin.Context) {
		filename := context.Param("filename")
		root := "./data/images/"
		context.File(filepath.Join(root, filename))
	})

	engine.POST("/completions", func(context *gin.Context) {
		var requestBody RequestBody
		if err := context.BindJSON(&requestBody); err != nil {
			logrus.Errorln("bind json error: ", err)
			context.AbortWithStatusJSON(http.StatusBadRequest, "chatGPT睡着了，再试试呗～")
			return
		}

		message := requestBody.Message
		if message == "" {
			logrus.Errorln("message is empty")
			context.AbortWithStatusJSON(http.StatusBadRequest, "chatGPT睡着了，再试试呗～")
			return
		}
		logrus.Infoln("发消息了: ", message)

		if strings.HasPrefix(strings.TrimSpace(message), "/image") {
			path, err := dalle.Generations(message, dalle.CreateImageResponseFormatB64JSON)
			if err != nil {
				logrus.Errorln("dall-e generations error: ", err)
				context.AbortWithStatusJSON(http.StatusBadRequest, "chatGPT睡着了，再试试呗～")
				return
			}
			logrus.Infoln("生成了一张图片: ", path)
			context.JSON(http.StatusOK, path)
			return
		}

		reply, err := chatgpt.Completions(message)
		if err != nil {
			logrus.Errorln("chatGPT completions error: ", err)
			context.AbortWithStatusJSON(http.StatusInternalServerError, "completions error")
			return
		}
		logrus.Infoln("和chatGPT对话得到的回复: ", reply)
		context.JSON(http.StatusOK, reply)
	})

	server := &http.Server{
		Addr:           ":8090",
		Handler:        engine,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logrus.Infoln("start server", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		logrus.Errorln("start server error: ", err)
		return
	}
}
