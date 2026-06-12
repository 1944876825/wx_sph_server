package ai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"wx_video_help/config"
)

// ChatRequest GLM聊天请求
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

// ChatMessage 消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse GLM响应
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

// GenerateReply 根据用户消息和系统提示词生成AI回复
func GenerateReply(userMsg string, systemPrompt string) (string, error) {
	cfg := config.Conf.AI
	if !cfg.Enabled || cfg.APIKey == "" {
		return "", fmt.Errorf("AI 未启用或 API Key 未配置")
	}

	prompt := cfg.DefaultPrompt
	if systemPrompt != "" {
		prompt = systemPrompt
	}

	reqBody := ChatRequest{
		Model: cfg.Model,
		Messages: []ChatMessage{
			{Role: "system", Content: prompt},
			{Role: "user", Content: userMsg},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", cfg.BaseURL+"/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}

	if chatResp.Error.Message != "" {
		return "", fmt.Errorf("GLM API 错误: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("GLM API 返回空结果")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}
