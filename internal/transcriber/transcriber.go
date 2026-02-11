package transcriber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Config 转录配置
type Config struct {
	APIKey   string
	BaseURL  string
	Model    string
}

// DefaultConfig 默认配置（使用 OpenAI Whisper API）
var DefaultConfig = Config{
	BaseURL: "https://api.openai.com/v1",
	Model:   "whisper-1",
}

// TranscriptionResponse Whisper API 响应
type TranscriptionResponse struct {
	Text string `json:"text"`
}

// Transcribe 将音频转录为文字
func Transcribe(audioPath string) (string, error) {
	config := DefaultConfig
	
	// 从环境变量获取 API Key
	config.APIKey = os.Getenv("OPENAI_API_KEY")
	if config.APIKey == "" {
		return "", fmt.Errorf("请设置 OPENAI_API_KEY 环境变量")
	}

	return transcribeWithWhisper(audioPath, config)
}

// TranscribeWithConfig 使用自定义配置转录
func TranscribeWithConfig(audioPath string, config Config) (string, error) {
	return transcribeWithWhisper(audioPath, config)
}

// transcribeWithWhisper 调用 Whisper API
func transcribeWithWhisper(audioPath string, config Config) (string, error) {
	// 打开音频文件
	file, err := os.Open(audioPath)
	if err != nil {
		return "", fmt.Errorf("打开音频文件失败: %w", err)
	}
	defer file.Close()

	// 构建 multipart 请求
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件
	part, err := writer.CreateFormFile("file", filepath.Base(audioPath))
	if err != nil {
		return "", fmt.Errorf("创建表单失败: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("写入文件数据失败: %w", err)
	}

	// 添加模型参数
	if err := writer.WriteField("model", config.Model); err != nil {
		return "", fmt.Errorf("写入模型参数失败: %w", err)
	}

	writer.Close()

	// 发送请求
	url := config.BaseURL + "/audio/transcriptions"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 错误 (%d): %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result TranscriptionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	return result.Text, nil
}
