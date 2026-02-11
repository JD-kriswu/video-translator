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

	"github.com/JD-kriswu/video-translator/internal/config"
)

// TranscriptionResponse Whisper API 响应
type TranscriptionResponse struct {
	Text string `json:"text"`
}

// Transcribe 将音频转录为文字
func Transcribe(audioPath string, cfg *config.TranscriberConfig) (string, error) {
	// 如果配置为空，使用环境变量
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	if cfg.APIKey == "" {
		return "", fmt.Errorf("请在配置文件中设置 api_key 或设置 OPENAI_API_KEY 环境变量")
	}

	return transcribeWithWhisper(audioPath, cfg)
}

// transcribeWithWhisper 调用 Whisper API
func transcribeWithWhisper(audioPath string, cfg *config.TranscriberConfig) (string, error) {
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
	if err := writer.WriteField("model", cfg.Model); err != nil {
		return "", fmt.Errorf("写入模型参数失败: %w", err)
	}

	writer.Close()

	// 发送请求
	url := cfg.BaseURL + "/audio/transcriptions"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
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
