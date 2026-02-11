package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Download 从 URL 下载视频文件
func Download(videoURL string) (string, error) {
	// 解析 URL
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("无效的 URL: %w", err)
	}

	// 生成文件名
	filename := generateFilename(parsedURL)
	outputPath := filepath.Join("temp", filename)

	// 确保 temp 目录存在
	if err := os.MkdirAll("temp", 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 检查是否是特殊平台（YouTube, Bilibili 等）
	if isSpecialPlatform(parsedURL.Host) {
		return downloadWithYtDlp(videoURL, outputPath)
	}

	// 普通 HTTP 下载
	return downloadHTTP(videoURL, outputPath)
}

// downloadHTTP 普通 HTTP 下载
func downloadHTTP(videoURL, outputPath string) (string, error) {
	resp, err := http.Get(videoURL)
	if err != nil {
		return "", fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP 状态码: %d", resp.StatusCode)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return outputPath, nil
}

// downloadWithYtDlp 使用 yt-dlp 下载（支持 YouTube, Bilibili 等）
func downloadWithYtDlp(videoURL, outputPath string) (string, error) {
	// TODO: 实现 yt-dlp 调用
	// 需要系统安装 yt-dlp: brew install yt-dlp
	return "", fmt.Errorf("yt-dlp 下载暂未实现，请使用直链")
}

// generateFilename 生成文件名
func generateFilename(u *url.URL) string {
	path := u.Path
	if path == "" || path == "/" {
		return "video.mp4"
	}
	
	filename := filepath.Base(path)
	if !hasVideoExtension(filename) {
		filename += ".mp4"
	}
	return filename
}

// isSpecialPlatform 检查是否是需要特殊处理的平台
func isSpecialPlatform(host string) bool {
	specialHosts := []string{
		"youtube.com", "youtu.be",
		"bilibili.com", "b23.tv",
		"vimeo.com",
		"twitter.com", "x.com",
	}
	
	for _, h := range specialHosts {
		if strings.Contains(host, h) {
			return true
		}
	}
	return false
}

// hasVideoExtension 检查是否有视频扩展名
func hasVideoExtension(filename string) bool {
	extensions := []string{".mp4", ".webm", ".mkv", ".avi", ".mov", ".flv"}
	lower := strings.ToLower(filename)
	for _, ext := range extensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}
