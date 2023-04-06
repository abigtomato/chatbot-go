package dalle

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/abigtomato/chatbot-go/config"
	"github.com/abigtomato/chatbot-go/tools"
	"github.com/sirupsen/logrus"
	"image/png"
	"os"
	"strings"
	"time"
)

const (
	CreateImageResponseFormatURL     = "url"
	CreateImageResponseFormatB64JSON = "b64_json"
)

type DERequest struct {
	Prompt         string `json:"prompt"`
	N              int    `json:"n"`
	Size           string `json:"size"`
	ResponseFormat string `json:"response_format"`
	User           string `json:"user"`
}

type DEResponse struct {
	Created int64    `json:"created"`
	Data    []DEData `json:"data"`
	Error   DEError  `json:"error"`
}

type DEData struct {
	Url     string `json:"url"`
	B64Json string `json:"b64_json"`
}

type DEError struct {
	Message string `json:"message"`
}

func Generations(content, format string) (string, error) {
	return GenerationsWithConfig(content, format, config.LoadConfig())
}

func GenerationsWithConfig(content, format string, config *config.Configuration) (string, error) {
	deReq, err := json.Marshal(DERequest{
		Prompt:         content,
		N:              config.DallE.ImageN,
		Size:           config.DallE.ImageSize,
		ResponseFormat: format,
	})
	if err != nil {
		return "", err
	}
	logrus.Infoln("request dall-E: ", string(deReq))

	respBody, err := tools.SOpenAI(config.DallE.BaseUrl, deReq)
	if err != nil {
		return "", err
	}

	deResp := &DEResponse{}
	if err := json.Unmarshal(respBody, &deResp); err != nil {
		return "", err
	}

	switch format {
	case CreateImageResponseFormatURL:
		var out []string
		for _, v := range deResp.Data {
			out = append(out, v.Url)
		}
		return strings.Join(out, ","), nil
	case CreateImageResponseFormatB64JSON:
		imgBytes, err := base64.StdEncoding.DecodeString(deResp.Data[0].B64Json)
		if err != nil {
			return "", err
		}

		imgData, err := png.Decode(bytes.NewReader(imgBytes))
		if err != nil {
			return "", err
		}

		if err := os.MkdirAll("data/images", 0755); err != nil {
			return "", err
		}

		imageName := time.Now().Format("20060102-150405") + ".png"
		file, err := os.Create("data/images/" + imageName)
		if err != nil {
			return "", err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		if err := png.Encode(file, imgData); err != nil {
			return "", err
		}
		return config.Common.ServiceUrl + "/images/" + imageName, nil
	default:
		return "", nil
	}
}
