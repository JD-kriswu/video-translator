# Video Translator

将视频链接转化为长文本，并翻译为中文的工具。

## 功能

- 🎬 从 URL 下载视频文件
- 🎤 提取音频并转录为文本（语音识别）
- 🌐 将文本翻译为中文
- 📄 输出完整的中文文本文档

## 技术栈

- Python 3.10+
- FFmpeg - 音视频处理
- Whisper / 其他 ASR - 语音识别
- 翻译 API（待定）

## 安装

```bash
# 克隆项目
git clone <repo-url>
cd video-translator

# 创建虚拟环境
python -m venv venv
source venv/bin/activate  # macOS/Linux
# venv\Scripts\activate  # Windows

# 安装依赖
pip install -r requirements.txt
```

## 使用

```bash
python main.py <video-url>
```

## 项目结构

```
video-translator/
├── main.py           # 入口文件
├── downloader.py     # 视频下载模块
├── transcriber.py    # 语音转文字模块
├── translator.py     # 翻译模块
├── requirements.txt  # 依赖列表
└── output/           # 输出目录
```

## TODO

- [ ] 视频下载功能
- [ ] 音频提取
- [ ] 语音识别集成
- [ ] 翻译功能
- [ ] CLI 参数支持
- [ ] 进度显示

## License

MIT
