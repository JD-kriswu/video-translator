package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 应用配置
type Config struct {
	Transcriber TranscriberConfig `json:"transcriber"`
	Translator  TranslatorConfig  `json:"translator"`
}

// TranscriberConfig 转录配置
type TranscriberConfig struct {
	Provider string `json:"provider"` // openai, azure, local
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
}

// TranslatorConfig 翻译配置
type TranslatorConfig struct {
	Provider string `json:"provider"` // openai, azure, deepl
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
	Model    string `json:"model"`
}

// DefaultConfig 默认配置
var DefaultConfig = Config{
	Transcriber: TranscriberConfig{
		Provider: "openai",
		BaseURL:  "https://api.openai.com/v1",
		Model:    "whisper-1",
	},
	Translator: TranslatorConfig{
		Provider: "openai",
		BaseURL:  "https://api.openai.com/v1",
		Model:    "gpt-4o-mini",
	},
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	// 如果没有指定路径，使用默认路径
	if path == "" {
		path = getDefaultConfigPath()
	}

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 文件不存在，返回默认配置
		cfg := DefaultConfig
		return &cfg, nil
	}

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 合并默认值
	mergeDefaults(&cfg)

	return &cfg, nil
}

// Save 保存配置文件
func Save(cfg *Config, path string) error {
	if path == "" {
		path = getDefaultConfigPath()
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 序列化配置
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// CreateExample 创建示例配置文件
func CreateExample(path string) error {
	example := Config{
		Transcriber: TranscriberConfig{
			Provider: "openai",
			BaseURL:  "https://api.openai.com/v1",
			APIKey:   "your-api-key-here",
			Model:    "whisper-1",
		},
		Translator: TranslatorConfig{
			Provider: "openai",
			BaseURL:  "https://api.openai.com/v1",
			APIKey:   "your-api-key-here",
			Model:    "gpt-4o-mini",
		},
	}

	return Save(&example, path)
}

// getDefaultConfigPath 获取默认配置文件路径
func getDefaultConfigPath() string {
	// 优先使用当前目录
	if _, err := os.Stat("config.json"); err == nil {
		return "config.json"
	}

	// 其次使用用户目录
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.json"
	}

	return filepath.Join(home, ".config", "video-translator", "config.json")
}

// mergeDefaults 合并默认值
func mergeDefaults(cfg *Config) {
	if cfg.Transcriber.Provider == "" {
		cfg.Transcriber.Provider = DefaultConfig.Transcriber.Provider
	}
	if cfg.Transcriber.BaseURL == "" {
		cfg.Transcriber.BaseURL = DefaultConfig.Transcriber.BaseURL
	}
	if cfg.Transcriber.Model == "" {
		cfg.Transcriber.Model = DefaultConfig.Transcriber.Model
	}

	if cfg.Translator.Provider == "" {
		cfg.Translator.Provider = DefaultConfig.Translator.Provider
	}
	if cfg.Translator.BaseURL == "" {
		cfg.Translator.BaseURL = DefaultConfig.Translator.BaseURL
	}
	if cfg.Translator.Model == "" {
		cfg.Translator.Model = DefaultConfig.Translator.Model
	}
}
