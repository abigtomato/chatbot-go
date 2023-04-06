package config

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	DefaultPath = "config.dev.yaml"

	DefaultChatGPTUrl       = "https://api.openai.com/v1/chat/completions"
	DefaultModel            = "gpt-3.5-turbo"
	DefaultMaxTokens        = 3072
	DefaultTemperature      = 0.9
	DefaultTopP             = 1
	DefaultFrequencyPenalty = 0
	DefaultPresencePenalty  = 0

	DefaultDallEUrl  = "https://api.openai.com/v1/images/generations"
	DefaultImageN    = 1
	DefaultImageSize = "512x512"
)

var config *Configuration
var once sync.Once

type Configuration struct {
	Common struct {
		ServiceUrl string `yaml:"serviceUrl"` // 服务地址
	} `yaml:"common"`
	Wechat struct {
		AppId          string `yaml:"appId"`          // appId
		AppSecret      string `yaml:"appSecret"`      // appSecret
		Token          string `yaml:"token"`          // token
		EncodingAESKey string `yaml:"encodingAESKey"` // 加密编码
	} `yaml:"wechat"`
	ChatGPT struct {
		ApiKey           string  `yaml:"apiKey"`           // apiKey
		BaseUrl          string  `yaml:"baseUrl"`          // 地址
		Model            string  `yaml:"model"`            // 模型，默认GPT-3.5
		MaxTokens        uint    `yaml:"maxTokens"`        // 生成结果时的最大单词数，GPT3.5是4000
		Temperature      float64 `yaml:"temperature"`      // 随机因子，控制结果的随机性，0～0.9
		TopP             int     `yaml:"topP"`             // 随机因子，0～0.9之间
		FrequencyPenalty int     `yaml:"frequencyPenalty"` // 重复度惩罚因子，是-2.0～2.0之间的数字
		PresencePenalty  int     `yaml:"presencePenalty"`  // 控制主题的重复度，是-2.0～2.0之间的数字
	} `yaml:"chatGPT"`
	DallE struct {
		BaseUrl   string `yaml:"baseUrl"`   // 地址
		ImageN    int    `yaml:"imageN"`    // 图片生成数量
		ImageSize string `yaml:"imageSize"` // 图片生成大小
	} `yaml:"dallE"`
}

func C() *Configuration {
	return LoadConfig()
}

func LoadConfig() *Configuration {
	return LoadConfigWithPath("")
}

func LoadConfigWithPath(path string) *Configuration {
	once.Do(func() {
		config = &Configuration{
			ChatGPT: struct {
				ApiKey           string  `yaml:"apiKey"`
				BaseUrl          string  `yaml:"baseUrl"`
				Model            string  `yaml:"model"`
				MaxTokens        uint    `yaml:"maxTokens"`
				Temperature      float64 `yaml:"temperature"`
				TopP             int     `yaml:"topP"`
				FrequencyPenalty int     `yaml:"frequencyPenalty"`
				PresencePenalty  int     `yaml:"presencePenalty"`
			}{
				BaseUrl:          DefaultChatGPTUrl,
				Model:            DefaultModel,
				MaxTokens:        DefaultMaxTokens,
				Temperature:      DefaultTemperature,
				TopP:             DefaultTopP,
				FrequencyPenalty: DefaultFrequencyPenalty,
				PresencePenalty:  DefaultPresencePenalty,
			},
			DallE: struct {
				BaseUrl   string `yaml:"baseUrl"`
				ImageN    int    `yaml:"imageN"`
				ImageSize string `yaml:"imageSize"`
			}{
				BaseUrl:   DefaultDallEUrl,
				ImageN:    DefaultImageN,
				ImageSize: DefaultImageSize,
			},
		}

		if path == "" {
			path = DefaultPath
		}

		file, err := os.Open(path)
		if err != nil {
			logrus.Errorln("open config error: ", err)
			return
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		if err := yaml.NewDecoder(file).Decode(&config); err != nil {
			logrus.Errorln("decode config error: ", err)
			return
		}
	})
	return config
}
