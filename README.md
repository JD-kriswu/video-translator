# Video Translator

将视频链接转化为长文本，并翻译为目标语言的工具。支持 CLI 和 Web 两种使用方式。

## 功能

- 🎬 从 URL 下载视频文件
- 🎤 提取音频并转录为文本（本地 Whisper / API）
- 🌐 将文本翻译为目标语言
- 📄 输出完整的翻译文档
- 🖥️ Web 界面，支持在线使用

## 技术栈

- Go 1.21+ / Gin
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

# 构建 CLI 版本
go build -o video-translator ./cmd/main.go

# 构建 Web 服务版本
go build -o video-translator-server ./cmd/server/main.go
```

## 系统依赖

```bash
# macOS
brew install ffmpeg
pip install openai-whisper
pip install yt-dlp  # 视频下载工具

# Ubuntu/Debian
sudo apt install ffmpeg
pip install openai-whisper
pip install yt-dlp
```

## 提高抖音下载成功率（可选但推荐）

为了提高抖音视频下载成功率，建议配置浏览器 cookies：

1. 使用浏览器插件导出 cookies（推荐 "Get cookies.txt LOCALLY"）
2. 访问 [www.douyin.com](https://www.douyin.com) 并导出 cookies
3. 将 `cookies.txt` 保存到项目根目录

详细说明请参考 [COOKIES.md](COOKIES.md)

**成功率对比：**
- 不使用 cookies：~60%
- 使用 cookies（未登录）：~70%
- 使用 cookies（已登录）：~90%

## 配置

创建 `config.json`：

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

### Whisper 模型

| 模型 | 大小 | 速度 | 质量 |
|------|------|------|------|
| tiny | ~75MB | 最快 | 一般 |
| base | ~150MB | 快 | 较好（推荐） |
| small | ~500MB | 中等 | 好 |
| medium | ~1.5GB | 较慢 | 很好 |
| large | ~3GB | 慢 | 最佳 |

## 使用方式

### Web 服务（推荐）

```bash
# 启动服务
./video-translator-server

# 访问 http://localhost:8080
```

打开浏览器，输入视频链接，选择目标语言，点击翻译即可。

### CLI 命令行

```bash
# 基本用法
./video-translator -url "https://example.com/video.mp4"

# 指定目标语言
./video-translator -url "https://xxx.mp4" -lang en

# 指定输出文件
./video-translator -url "https://xxx.mp4" -output result.txt
```

## API 接口

### 创建翻译任务

```bash
POST /api/translate
Content-Type: application/json

{
  "url": "https://example.com/video.mp4",
  "lang": "zh"
}
```

### 查询任务状态

```bash
GET /api/task/:id
```

### 获取翻译结果

```bash
GET /api/task/:id/result
```

## 项目结构

```
video-translator/
├── cmd/
│   ├── main.go              # CLI 入口
│   └── server/
│       ├── main.go          # Web 服务入口
│       └── static/          # 前端静态文件
├── internal/
│   ├── config/              # 配置管理
│   ├── server/              # HTTP 处理器
│   ├── downloader/          # 视频下载
│   ├── extractor/           # 音频提取
│   ├── transcriber/         # 语音转文字
│   └── translator/          # 翻译
├── config.json              # 配置文件
└── README.md
```

## License

MIT
sk-XIE4889qlgwZ9Ort6RbdBR5idOihzup5EtzJjPXGGZmC2PUm