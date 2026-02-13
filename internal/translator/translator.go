package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/JD-kriswu/video-translator/internal/config"
)

// ChatRequest OpenAI Chat API 请求
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message 消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse OpenAI Chat API 响应
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Translate 翻译文本
func Translate(text, targetLang string, cfg *config.TranslatorConfig) (string, error) {
	// 如果配置为空，使用环境变量
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	if cfg.APIKey == "" {
		return "", fmt.Errorf("请在配置文件中设置 api_key 或设置 OPENAI_API_KEY 环境变量")
	}

	return translateWithOpenAI(text, targetLang, cfg)
}

// translateWithOpenAI 使用 OpenAI API 翻译
func translateWithOpenAI(text, targetLang string, cfg *config.TranslatorConfig) (string, error) {
	langName := getLanguageName(targetLang)

	prompt := fmt.Sprintf(`请将以下文本翻译成%s。要求：
1. 保持原文的语气和风格
2. 专业术语翻译准确
3. 语句通顺自然
4. 只输出翻译结果，不要添加任何解释

原文：
%s`, langName, text)

	reqBody := ChatRequest{
		Model: cfg.Model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	url := cfg.BaseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 错误 (%d): %s", resp.StatusCode, string(body))
	}

	var result ChatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("API 返回空结果")
	}

	content := result.Choices[0].Message.Content
	
	// 过滤掉 <think>...</think> 标签
	content = removeThinkTags(content)
	
	return strings.TrimSpace(content), nil
}

// removeThinkTags 移除 <think>...</think> 标签
func removeThinkTags(text string) string {
	// 简单的正则替换
	for {
		start := strings.Index(text, "<think>")
		if start == -1 {
			break
		}
		end := strings.Index(text, "</think>")
		if end == -1 {
			break
		}
		text = text[:start] + text[end+8:]
	}
	return text
}

// getLanguageName 获取语言名称
func getLanguageName(code string) string {
	languages := map[string]string{
		"zh": "中文",
		"en": "英文",
		"ja": "日文",
		"ko": "韩文",
		"fr": "法文",
		"de": "德文",
		"es": "西班牙文",
		"ru": "俄文",
		"pt": "葡萄牙文",
		"it": "意大利文",
	}

	if name, ok := languages[code]; ok {
		return name
	}
	return code
}
