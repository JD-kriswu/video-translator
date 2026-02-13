package transcriber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/JD-kriswu/video-translator/internal/config"
)

// TranscriptionResponse Whisper API 响应
type TranscriptionResponse struct {
	Text string `json:"text"`
}

// Transcribe 将音频转录为文字
func Transcribe(audioPath string, cfg *config.TranscriberConfig) (string, error) {
	switch cfg.Provider {
	case "local", "whisper":
		return transcribeWithLocalWhisper(audioPath, cfg)
	case "openai", "api":
		return transcribeWithAPI(audioPath, cfg)
	default:
		// 默认使用本地 whisper
		return transcribeWithLocalWhisper(audioPath, cfg)
	}
}

// transcribeWithLocalWhisper 使用本地 whisper 转录
func transcribeWithLocalWhisper(audioPath string, cfg *config.TranscriberConfig) (string, error) {
	// 检查 whisper 是否安装
	whisperPath, err := findWhisperExecutable()
	if err != nil {
		return "", fmt.Errorf("whisper 未安装: %w\n请运行: pip install openai-whisper", err)
	}

	// 确定模型
	model := cfg.Model
	if model == "" || model == "whisper-1" {
		model = "base" // 默认使用 base 模型，平衡速度和质量
	}

	// 输出文件路径
	outputDir := filepath.Dir(audioPath)
	baseName := strings.TrimSuffix(filepath.Base(audioPath), filepath.Ext(audioPath))

	// 构建命令
	// whisper audio.mp3 --model base --output_format txt --output_dir ./
	cmd := exec.Command(whisperPath,
		audioPath,
		"--model", model,
		"--output_format", "txt",
		"--output_dir", outputDir,
	)

	// 执行命令
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper 执行失败: %w\n输出: %s", err, string(output))
	}

	// 读取输出文件
	txtPath := filepath.Join(outputDir, baseName+".txt")
	text, err := os.ReadFile(txtPath)
	if err != nil {
		return "", fmt.Errorf("读取转录结果失败: %w", err)
	}

	// 清理临时文件
	os.Remove(txtPath)

	return strings.TrimSpace(string(text)), nil
}

// transcribeWithAPI 使用 API 转录
func transcribeWithAPI(audioPath string, cfg *config.TranscriberConfig) (string, error) {
	// 如果配置为空，使用环境变量
	if cfg.APIKey == "" {
		cfg.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	if cfg.APIKey == "" {
		return "", fmt.Errorf("请在配置文件中设置 api_key 或设置 OPENAI_API_KEY 环境变量")
	}

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
	model := cfg.Model
	if model == "" {
		model = "whisper-1"
	}
	if err := writer.WriteField("model", model); err != nil {
		return "", fmt.Errorf("写入模型参数失败: %w", err)
	}

	writer.Close()

	// 发送请求
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	url := baseURL + "/audio/transcriptions"
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

// findWhisperExecutable 查找 whisper 可执行文件
func findWhisperExecutable() (string, error) {
	// 尝试常见的命令名
	candidates := []string{"whisper", "whisper.exe"}

	for _, name := range candidates {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("whisper 命令未找到")
}

// GetAvailableModels 获取可用的本地模型列表
func GetAvailableModels() []string {
	return []string{
		"tiny",    // 最快，质量一般
		"base",    // 平衡速度和质量（推荐）
		"small",   // 较好质量
		"medium",  // 高质量
		"large",   // 最高质量，需要较多资源
		"large-v2",
		"large-v3",
	}
}
