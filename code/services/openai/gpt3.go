package openai

import (
	"errors"
	"math/rand"
	"time"
)

const (
	maxTokens                = 2000
	temperature              = 1.0
	temperature_min          = 0.3
	temperature_max          = 1.5
	using_random_temperature = true
	engine                   = "gpt-3.5-turbo"
)

type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatGPTChoiceItem    `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}
type ChatGPTChoiceItem struct {
	Message      Messages `json:"message"`
	Index        int      `json:"index"`
	FinishReason string   `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string     `json:"model"`
	Messages         []Messages `json:"messages"`
	MaxTokens        int        `json:"max_tokens"`
	Temperature      float32    `json:"temperature"`
	TopP             int        `json:"top_p"`
	FrequencyPenalty int        `json:"frequency_penalty"`
	PresencePenalty  int        `json:"presence_penalty"`
}

// 生成步长为0.1的随机浮点数temperature
func randomTemperature() float32 {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子
	step := 1.0
	randomFloat := float64(temperature_min) + float64(rand.Float32()*(temperature_max-temperature_min))/step*step // 生成随机浮点数
	result := float32(int(randomFloat*10)) / 10                                                                   // 转换为最小单位为0.1的浮点数
	return result
}

func (gpt *ChatGPT) Completions(msg []Messages) (resp Messages,
	err error) {
	requestBody := ChatGPTRequestBody{
		Model:            engine,
		Messages:         msg,
		MaxTokens:        maxTokens,
		Temperature:      now_temp,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	gptResponseBody := &ChatGPTResponseBody{}
	url := gpt.FullUrl("chat/completions")
	//fmt.Println(url)
	if url == "" {
		return resp, errors.New("无法获取openai请求地址")
	}
	err = gpt.sendRequestWithBodyType(url, "POST", jsonBody, requestBody, gptResponseBody)
	if err == nil && len(gptResponseBody.Choices) > 0 {
		resp = gptResponseBody.Choices[0].Message
	} else {
		resp = Messages{}
		err = errors.New("openai 请求失败")
	}
	return resp, err
}
