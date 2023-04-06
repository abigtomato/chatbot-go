package chatgpt

import (
	"encoding/json"
	"github.com/abigtomato/chatbot-go/config"
	"github.com/abigtomato/chatbot-go/tools"
	"github.com/sirupsen/logrus"
)

const DefaultUser = "user"

type GPTRequest struct {
	Model            string     `json:"model"`             // 模型
	Messages         *[]Message `json:"messages"`          // 消息
	MaxTokens        uint       `json:"max_tokens"`        // 生成结果时的最大单词数，不能超过模型的上下文长度
	Temperature      float64    `json:"temperature"`       // 随机因子，控制结果的随机性，如果希望结果更有创意可以尝试 0.9，或者希望有固定结果可以尝试 0.0
	TopP             int        `json:"top_p"`             // 随机因子2，一个可用于代替 temperature 的参数，对应机器学习中 nucleus sampling（核采样），如果设置 0.1 意味着只考虑构成前 10% 概率质量的 tokens
	FrequencyPenalty int        `json:"frequency_penalty"` // 重复度惩罚因子，是 -2.0 ~ 2.0 之间的数字，正值会根据新 tokens 在文本中的现有频率对其进行惩罚，从而降低模型逐字重复同一行的可能性
	PresencePenalty  int        `json:"presence_penalty"`  // 控制主题的重复度，是 -2.0 ~ 2.0 之间的数字，正值会根据到目前为止是否出现在文本中来惩罚新 tokens，从而增加模型谈论新主题的可能性
}

type Message struct {
	Role    string `json:"role"`    // 角色
	Content string `json:"content"` // 内容
}

type GPTResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int            `json:"created"`
	Model   string         `json:"model"`
	Usage   map[string]any `json:"usage"`
	Choices []ChoiceItem   `json:"choices"`
}

type ChoiceItem struct {
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
}

func Completions(content string) (string, error) {
	return CompletionsWithConfig(content, config.LoadConfig())
}

func CompletionsWithConfig(content string, config *config.Configuration) (string, error) {
	gptReq, err := json.Marshal(&GPTRequest{
		Model: config.ChatGPT.Model,
		Messages: &[]Message{
			{
				Role:    DefaultUser,
				Content: content,
			},
		},
		Temperature:      config.ChatGPT.Temperature,
		MaxTokens:        config.ChatGPT.MaxTokens,
		TopP:             config.ChatGPT.TopP,
		FrequencyPenalty: config.ChatGPT.FrequencyPenalty,
		PresencePenalty:  config.ChatGPT.PresencePenalty,
	})
	if err != nil {
		return "", err
	}
	logrus.Infoln("request chatGPT: ", string(gptReq))

	respBody, err := tools.SOpenAI(config.ChatGPT.BaseUrl, gptReq)
	if err != nil {
		logrus.Errorln("http send error: ", err)
		return "", err
	}
	logrus.Infoln("chatGPT response: ", string(respBody))

	gptResp := &GPTResponse{}
	if err := json.Unmarshal(respBody, gptResp); err != nil {
		return "", err
	}
	return gptResp.Choices[0].Message.Content, nil
}
