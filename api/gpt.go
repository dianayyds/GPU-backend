package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	nurl "net/url"
	"time"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}
type Choice struct {
	Index        int              `json:"index"`
	Message      Message          `json:"message"`
	Logprobs     *json.RawMessage `json:"logprobs"` // 使用指针允许空值
	FinishReason string           `json:"finish_reason"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func GptHandler(g *gin.Context) {
	start := time.Now() // 获取开始时间
	text := g.Query("text")
	//// 将替换为你的OpenAI API密钥
	apiKey := "YOUR_API_KEY"
	url := "https://api.openai.com/v1/chat/completions"
	// 构建请求体
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo-0125", 
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": text,
					},
				},
			},
		},
	})
	proxyUrl, _ := nurl.Parse("http://127.0.0.1:59527")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	client := &http.Client{
		Transport: transport,
	}
	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		seelog.Error(err)
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		seelog.Error(err)
	}
	defer resp.Body.Close()
	// 读取响应体
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		seelog.Error(err)
	}
	var response ApiResponse
	json.Unmarshal([]byte(responseBody), &response)
	elapsed := time.Since(start)
	// 打印消息内容
	g.JSON(200, gin.H{
		"code":   0,
		"answer": response.Choices[0].Message.Content,
		"time":   elapsed,
	})

}
