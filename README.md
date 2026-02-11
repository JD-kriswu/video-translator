# Video Translator

将视频链接转化为长文本，并翻译为中文的工具。

## 功能

- 🎬 从 URL 下载视频文件
- 🎤 提取音频并转录为文本（语音识别）
- 🌐 将文本翻译为中文
- 📄 输出完整的中文文本文档

## 技术栈

- Go 1.21+
- FFmpeg - 音视频处理
- Whisper API / OpenAI - 语音识别
- 翻译 API（OpenAI / DeepL / 其他）

## 安装

```bash
# 克隆项目
git clone <repo-url>
cd video-translator

# 安装依赖
go mod tidy

# 构建
go build -o video-translator ./cmd/main.go
```

## 依赖

确保系统已安装 FFmpeg：

```bash
# macOS
brew install ffmpeg

# Ubuntu/Debian
sudo apt install ffmpeg
```

## 使用

```bash
./video-translator -url <video-url> -output output.txt
```

## 项目结构

```
video-translator/
├── cmd/
│   └── main.go           # 入口文件
├── internal/
│   ├── downloader/       # 视频下载模块
│   │   └── downloader.go
│   ├── extractor/        # 音频提取模块
│   │   └── extractor.go
│   ├── transcriber/      # 语音转文字模块
│   │   └── transcriber.go
│   └── translator/       # 翻译模块
│       └── translator.go
├── pkg/
│   └── utils/            # 工具函数
├── output/               # 输出目录
├── go.mod
├── go.sum
└── README.md
```

## TODO

- [ ] 视频下载功能（支持 YouTube、Bilibili 等）
- [ ] FFmpeg 音频提取
- [ ] Whisper API 集成
- [ ] 翻译功能
- [ ] CLI 参数支持（cobra）
- [ ] 进度显示
- [ ] 并发处理长视频

## License

MIT
