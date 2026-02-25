## Context Structure

This project uses L0/L1/L2 context optimization.

- `.context/index.md` - L0: File index (one-liner per file)
- `.context/*/_overview.md` - L1: Structure summaries (signatures only)
- Source files - L2: Full code (load on demand)

When exploring the codebase:
1. Start with `.context/index.md` to find relevant files
2. Read `.context/<dir>/_overview.md` for structure
3. Only load full source when you need implementation details

## Project Structure

This is a video translation tool that:
- Downloads videos from various platforms (YouTube, Bilibili, Douyin, Xiaohongshu, etc.)
- Extracts audio from videos
- Transcribes audio to text using Whisper (local or API)
- Translates text to target language using OpenAI API

Key directories:
- `cmd/` - CLI and Web server entry points
- `internal/downloader/` - Video download logic (supports multiple platforms)
- `internal/extractor/` - Audio extraction using FFmpeg
- `internal/transcriber/` - Speech-to-text transcription
- `internal/translator/` - Text translation
- `internal/server/` - HTTP API handlers
- `internal/auth/` - User authentication
- `internal/database/` - MySQL and Redis clients

## Performance Tips

- For general exploration, read `.context/index.md` first (< 100 tokens)
- For understanding module structure, read relevant `_overview.md` (~ 2k tokens per module)
- Only read full source files when you need implementation details
