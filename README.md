# Video Translator

将视频链接转化为长文本，并翻译为中文的工具。

## 功能

- 🎬 从 URL 下载视频文件
- 🎤 提取音频并转录为文本（本地 Whisper / API）
- 🌐 将文本翻译为中文
- 📄 输出完整的中文文本文档

## 技术栈

- Go 1.21+
- FFmpeg - 音视频处理
- OpenAI Whisper（本地） - 语音识别
- OpenAI API - 翻译

## 安装

```bash
# 克隆项目
git clone https://github.com/JD-kriswu/video-translator.git
cd video-translator

# 安装依赖
go mod tidy

# 构建
go build -o video-translator ./cmd/main.go
```

## 系统依赖

```bash
# macOS
brew install ffmpeg
pip install openai-whisper

# Ubuntu/Debian
sudo apt install ffmpeg
pip install openai-whisper
```

### Whisper 模型说明

首次运行时会自动下载模型，可选模型：

| 模型 | 大小 | 速度 | 质量 |
|------|------|------|------|
| tiny | ~75MB | 最快 | 一般 |
| base | ~150MB | 快 | 较好（推荐） |
| small | ~500MB | 中等 | 好 |
| medium | ~1.5GB | 较慢 | 很好 |
| large | ~3GB | 慢 | 最佳 |

## 配置

生成配置文件：

```bash
./video-translator -init
```

编辑 `config.json`：

```json
{
  "transcriber": {
    "provider": "local",
    "model": "base"
  },
  "translator": {
    "provider": "openai",
    "base_url": "https://api.openai.com/v1",
    "api_key": "your-api-key",
    "model": "gpt-4o-mini"
  }
}
```

### 转录配置

- `provider`: `local`（本地 Whisper）或 `openai`（API 调用）
- `model`: 本地模式可选 tiny/base/small/medium/large

### 翻译配置

- `base_url`: API 地址（支持兼容 OpenAI 的服务）
- `api_key`: API 密钥
- `model`: 模型名称

## 使用

```bash
# 基本用法
./video-translator -url "https://example.com/video.mp4"

# 指定输出文件
./video-translator -url "https://xxx.mp4" -output result.txt

# 指定目标语言
./video-translator -url "https://xxx.mp4" -lang en

# 保留中间文件（视频、音频）
./video-translator -url "https://xxx.mp4" -keep

# 使用指定配置文件
./video-translator -url "https://xxx.mp4" -config my-config.json
```

## 项目结构

```
video-translator/
├── cmd/
│   └── main.go              # 入口文件
├── internal/
│   ├── config/              # 配置管理
│   ├── downloader/          # 视频下载
│   ├── extractor/           # 音频提取
│   ├── transcriber/         # 语音转文字
│   └── translator/          # 翻译
├── config.example.json      # 示例配置
├── go.mod
└── README.md
```

## License

MIT
