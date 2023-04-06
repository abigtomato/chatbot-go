package tools

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/abigtomato/chatbot-go/config"
	"io"
	"net/http"
	"time"
)

func SOpenAI(baseUrl string, reqData []byte) ([]byte, error) {
	httpReq, err := http.NewRequest(http.MethodPost, baseUrl, bytes.NewBuffer(reqData))
	if err != nil {
		return nil, errors.New("new http request error: " + err.Error())
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Content-Type", "application/json;charset=utf-8")
	httpReq.Header.Set("Authorization", "Bearer "+config.C().ChatGPT.ApiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.New("send http error: " + err.Error())
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(httpResp.Body)

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, errors.New(fmt.Sprintf("response error, code=%d, details=%v\n", httpResp.StatusCode, string(body)))
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.New("read http response error: " + err.Error())
	}
	return body, err
}
