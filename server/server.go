package server

import (
	"github.com/abigtomato/chatbot-go/server/dingtalk"
	"github.com/abigtomato/chatbot-go/server/wechat"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run() {
	engine := gin.Default()

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

	server := &http.Server{
		Addr:           ":8090",
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logrus.Infoln("start server", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		logrus.Errorln("start server error: ", err)
		return
	}
}
