package extractor

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ExtractAudio 从视频文件提取音频
func ExtractAudio(videoPath string) (string, error) {
	// 检查 ffmpeg 是否可用
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("ffmpeg 未安装，请运行: brew install ffmpeg")
	}

	// 生成输出路径
	baseName := strings.TrimSuffix(filepath.Base(videoPath), filepath.Ext(videoPath))
	audioPath := filepath.Join(filepath.Dir(videoPath), baseName+".mp3")

	// 构建 ffmpeg 命令
	// -i: 输入文件
	// -vn: 不处理视频
	// -acodec mp3: 使用 mp3 编码
	// -ar 16000: 采样率 16kHz（Whisper 推荐）
	// -ac 1: 单声道
	// -y: 覆盖已存在的文件
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vn",
		"-acodec", "mp3",
		"-ar", "16000",
		"-ac", "1",
		"-y",
		audioPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg 执行失败: %w\n输出: %s", err, string(output))
	}

	return audioPath, nil
}

// GetAudioDuration 获取音频时长（秒）
func GetAudioDuration(audioPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		audioPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("获取音频时长失败: %w", err)
	}

	var duration float64
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	if err != nil {
		return 0, fmt.Errorf("解析时长失败: %w", err)
	}

	return duration, nil
}
