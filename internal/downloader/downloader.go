package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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

	// 确保 temp 目录存在
	if err := os.MkdirAll("temp", 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 抖音专用下载器（优先）
	if isDouyinURL(parsedURL.Host) {
		path, err := downloadDouyin(videoURL)
		if err == nil {
			return path, nil
		}
		// 抖音下载失败，回退到 yt-dlp
		fmt.Printf("抖音专用下载失败: %v，尝试 yt-dlp\n", err)
	}

	// 检查是否是特殊平台（YouTube, Bilibili, 抖音等）
	if isSpecialPlatform(parsedURL.Host) {
		return downloadWithYtDlp(videoURL)
	}

	// 生成文件名
	filename := generateFilename(parsedURL)
	outputPath := filepath.Join("temp", filename)

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

// downloadWithYtDlp 使用 yt-dlp 下载（支持 YouTube, Bilibili, 抖音等）
func downloadWithYtDlp(videoURL string) (string, error) {
	// 检查 yt-dlp 是否可用
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return "", fmt.Errorf("yt-dlp 未安装，请运行: pip install yt-dlp")
	}

	// 输出文件模板
	outputTemplate := "temp/%(id)s.%(ext)s"

	// 构建 yt-dlp 命令
	cmd := exec.Command("yt-dlp",
		"-f", "best[ext=mp4]/best",  // 优先下载 mp4 格式
		"-o", outputTemplate,        // 输出文件名模板
		"--no-playlist",             // 不下载播放列表
		"--no-warnings",             // 不显示警告
		"--print", "after_move:filepath", // 打印最终文件路径
		videoURL,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("yt-dlp 下载失败: %w\n输出: %s", err, string(output))
	}

	// 解析输出获取文件路径
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("yt-dlp 未返回文件路径")
	}

	// 最后一行是文件路径
	filePath := strings.TrimSpace(lines[len(lines)-1])
	if filePath == "" || !fileExists(filePath) {
		// 尝试在 temp 目录查找最新文件
		files, _ := filepath.Glob("temp/*.*")
		if len(files) > 0 {
			filePath = files[len(files)-1]
		} else {
			return "", fmt.Errorf("下载完成但找不到文件")
		}
	}

	return filePath, nil
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
		"douyin.com", "iesdouyin.com",  // 抖音
		"tiktok.com",                    // TikTok
		"kuaishou.com",                  // 快手
		"weibo.com",                     // 微博
		"xiaohongshu.com",               // 小红书
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
