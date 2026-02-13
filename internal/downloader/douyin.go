package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DouyinConfig 抖音下载配置
type DouyinConfig struct {
	Cookies   string
	UserAgent string
	WebID     string
	UIFID     string
	MsToken   string
	VerifyFp  string
}

// DefaultDouyinConfig 默认配置
var DefaultDouyinConfig = DouyinConfig{
	UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36",
	WebID:     "7585838252910429748",
	UIFID:     "5bdad390e71fd6e6e69e3cafe6018169c2447c8bc0b8484cc0f203a274f99fdb3a3d54e5ec3c0ad4f24c844b1f57f7f057502e41ea7778d2c38e21ad0cc807ef1393c5af3af876ae08a55933add6821eb5263820825b9dbacb05b06c93e50ef0270edfee1388d84722b4574173badfd7ba0e980d9e9cd782014fea89d0588af4bb3f4e5dd6858a997bbdd240a2686af6c7bd242f175368846485c0513fa5fc30",
	VerifyFp:  "verify_mjdz6bpy_mCAPmgkE_ol9G_4wZJ_84HZ_6g7NgZH1SD1l",
	MsToken:   "fiGH98aAiAEEXbDjlKmbasHnDIYAiL_Yon80Pqb0ivKC2SuUwjJrqziHK2n99qEJQdYZPY0DsuFAu2WoH6GcCGOiVxArHR1OSTIy0B_kQ2froVRR0OV7SVxpOQb781raOZDUJL8JrIVdkUIIhPQNaFA7RGFpV7vHPSZ2lEp7hgTR",
	Cookies:   "ttwid=1%7CJ1CYy68DAYeni0clNuH6aEHdPKUVMcjgLwxZQA3YpNI%7C1770865593%7C4210c2e5ec2844c9ce288365b21276a7d93cbb225d1982a4929fccd294e6fd44; UIFID=5bdad390e71fd6e6e69e3cafe6018169c2447c8bc0b8484cc0f203a274f99fdb3a3d54e5ec3c0ad4f24c844b1f57f7f057502e41ea7778d2c38e21ad0cc807ef1393c5af3af876ae08a55933add6821eb5263820825b9dbacb05b06c93e50ef0270edfee1388d84722b4574173badfd7ba0e980d9e9cd782014fea89d0588af4bb3f4e5dd6858a997bbdd240a2686af6c7bd242f175368846485c0513fa5fc30; passport_csrf_token=ad2ac00ed1b94857c8fba95a1135b416; passport_csrf_token_default=ad2ac00ed1b94857c8fba95a1135b416; odin_tt=4d4095511d996eab818dc07898de614bd72794b9c84cbf4541e1e6b393fd5aad4c0f6927971510a7aa3a20694a3b2b5f08aa41589a9a160be8ff34ce6ce2d7885176baae847a479d3a2cf3c3c1b0a47c; bd_ticket_guard_client_web_domain=2; volume_info=%7B%22isMute%22%3Atrue%2C%22isUserMute%22%3Atrue%2C%22volume%22%3A0.5%7D; strategyABtestKey=%221770864424.07%22; is_dash_user=1; IsDouyinActive=true; home_can_add_dy_2_desktop=%221%22; stream_recommend_feed_params=%22%7B%5C%22cookie_enabled%5C%22%3Atrue%2C%5C%22screen_width%5C%22%3A2560%2C%5C%22screen_height%5C%22%3A1440%2C%5C%22browser_online%5C%22%3Atrue%2C%5C%22cpu_core_num%5C%22%3A16%2C%5C%22device_memory%5C%22%3A8%2C%5C%22downlink%5C%22%3A10%2C%5C%22effective_type%5C%22%3A%5C%224g%5C%22%2C%5C%22round_trip_time%5C%22%3A100%7D%22",
}

// douyinClient 抖音专用 HTTP 客户端
var douyinClient = &http.Client{
	Timeout: 30 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 跟随重定向，但保留 headers
		if len(via) >= 10 {
			return fmt.Errorf("too many redirects")
		}
		return nil
	},
}

// downloadDouyin 抖音视频下载
func downloadDouyin(videoURL string) (string, error) {
	// 方案1: 使用增强版 yt-dlp（推荐，成功率最高）
	path, err := downloadWithYtDlpEnhanced(videoURL, "")
	if err == nil {
		fmt.Println("✓ 使用 yt-dlp 下载成功")
		return path, nil
	}
	fmt.Printf("yt-dlp 下载失败: %v，尝试 Playwright 方式\n", err)

	// 方案2: Playwright 方案（备用）
	path, err = downloadDouyinWithPlaywright(videoURL)
	if err == nil {
		fmt.Println("✓ 使用 Playwright 下载成功")
		return path, nil
	}
	fmt.Printf("Playwright 下载失败: %v，尝试 API 方式\n", err)

	// 方案3: API 方式（最后备选）
	// 1. 解析视频 ID
	awemeID, err := extractDouyinAwemeID(videoURL)
	if err != nil {
		return "", fmt.Errorf("解析抖音视频ID失败: %w", err)
	}

	// 2. 获取视频详情
	videoInfo, err := getDouyinVideoInfo(awemeID)
	if err != nil {
		return "", fmt.Errorf("获取视频信息失败: %w", err)
	}

	// 3. 下载视频
	outputPath := filepath.Join("temp", fmt.Sprintf("%s.mp4", awemeID))
	if err := downloadDouyinVideo(videoInfo.PlayURL, outputPath); err != nil {
		return "", fmt.Errorf("下载视频失败: %w", err)
	}

	fmt.Println("✓ 使用 API 下载成功")
	return outputPath, nil
}

// downloadDouyinWithPlaywright 使用 Playwright 下载抖音视频
func downloadDouyinWithPlaywright(videoURL string) (string, error) {
	// 生成输出文件名
	outputPath := filepath.Join("temp", fmt.Sprintf("douyin_%d.mp4", time.Now().UnixNano()))

	// 调用 Node.js 脚本
	cmd := exec.Command("node", "scripts/douyin-download.js", videoURL, outputPath)
	cmd.Dir = "/data/code/video-translator"

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Playwright 脚本执行失败: %w\n输出: %s", err, string(output))
	}

	// 解析输出，最后一行是文件路径
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("脚本未返回文件路径")
	}

	filePath := strings.TrimSpace(lines[len(lines)-1])
	if !fileExists(filePath) {
		// 尝试相对路径
		filePath = filepath.Join("/data/code/video-translator", filePath)
		if !fileExists(filePath) {
			return "", fmt.Errorf("下载完成但找不到文件: %s", filePath)
		}
	}

	return filePath, nil
}

// DouyinVideoInfo 视频信息
type DouyinVideoInfo struct {
	AwemeID  string
	Title    string
	PlayURL  string
	CoverURL string
	Author   string
}

// extractDouyinAwemeID 从 URL 提取视频 ID
func extractDouyinAwemeID(videoURL string) (string, error) {
	// 处理短链接重定向
	if strings.Contains(videoURL, "v.douyin.com") || strings.Contains(videoURL, "iesdouyin.com") {
		realURL, err := resolveDouyinShortURL(videoURL)
		if err != nil {
			return "", err
		}
		videoURL = realURL
	}

	// 匹配视频 ID
	// 格式: /video/7123456789012345678
	re := regexp.MustCompile(`/video/(\d+)`)
	matches := re.FindStringSubmatch(videoURL)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	// 格式: /note/7123456789012345678 (图文)
	re = regexp.MustCompile(`/note/(\d+)`)
	matches = re.FindStringSubmatch(videoURL)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	// 格式: modal_id=7123456789012345678
	re = regexp.MustCompile(`modal_id=(\d+)`)
	matches = re.FindStringSubmatch(videoURL)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	return "", fmt.Errorf("无法从 URL 提取视频 ID: %s", videoURL)
}

// resolveDouyinShortURL 解析短链接
func resolveDouyinShortURL(shortURL string) (string, error) {
	req, err := http.NewRequest("GET", shortURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", DefaultDouyinConfig.UserAgent)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // 不自动跟随重定向
		},
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location != "" {
		return location, nil
	}

	return shortURL, nil
}

// getDouyinVideoInfo 获取视频详情
func getDouyinVideoInfo(awemeID string) (*DouyinVideoInfo, error) {
	// 构建 API URL
	apiURL := buildDouyinAPIURL(awemeID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	setDouyinHeaders(req)

	resp, err := douyinClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 请求失败: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseDouyinResponse(body, awemeID)
}

// buildDouyinAPIURL 构建 API URL
func buildDouyinAPIURL(awemeID string) string {
	params := url.Values{}
	params.Set("device_platform", "webapp")
	params.Set("aid", "6383")
	params.Set("channel", "channel_pc_web")
	params.Set("aweme_id", awemeID)
	params.Set("pc_client_type", "1")
	params.Set("pc_libra_divert", "Windows")
	params.Set("version_code", "170400")
	params.Set("version_name", "17.4.0")
	params.Set("cookie_enabled", "true")
	params.Set("screen_width", "2560")
	params.Set("screen_height", "1440")
	params.Set("browser_language", "zh-CN")
	params.Set("browser_platform", "Win32")
	params.Set("browser_name", "Chrome")
	params.Set("browser_version", "144.0.0.0")
	params.Set("browser_online", "true")
	params.Set("engine_name", "Blink")
	params.Set("engine_version", "144.0.0.0")
	params.Set("os_name", "Windows")
	params.Set("os_version", "10")
	params.Set("cpu_core_num", "16")
	params.Set("device_memory", "8")
	params.Set("platform", "PC")
	params.Set("downlink", "10")
	params.Set("effective_type", "4g")
	params.Set("round_trip_time", "100")
	params.Set("webid", DefaultDouyinConfig.WebID)
	params.Set("uifid", DefaultDouyinConfig.UIFID)
	params.Set("verifyFp", DefaultDouyinConfig.VerifyFp)
	params.Set("fp", DefaultDouyinConfig.VerifyFp)
	params.Set("msToken", DefaultDouyinConfig.MsToken)

	return "https://www.douyin.com/aweme/v1/web/aweme/detail/?" + params.Encode()
}

// setDouyinHeaders 设置请求头
func setDouyinHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("User-Agent", DefaultDouyinConfig.UserAgent)
	req.Header.Set("Referer", "https://www.douyin.com/")
	req.Header.Set("Origin", "https://www.douyin.com")

	// 设置 cookies
	if DefaultDouyinConfig.Cookies != "" {
		req.Header.Set("Cookie", DefaultDouyinConfig.Cookies)
	}
}

// parseDouyinResponse 解析 API 响应
func parseDouyinResponse(body []byte, awemeID string) (*DouyinVideoInfo, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	// 检查状态码
	if statusCode, ok := result["status_code"].(float64); ok && statusCode != 0 {
		return nil, fmt.Errorf("API 返回错误: %v", result)
	}

	// 提取视频信息
	awemeDetail, ok := result["aweme_detail"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("无法获取视频详情，响应: %s", string(body))
	}

	info := &DouyinVideoInfo{
		AwemeID: awemeID,
	}

	// 标题
	if desc, ok := awemeDetail["desc"].(string); ok {
		info.Title = desc
	}

	// 作者
	if author, ok := awemeDetail["author"].(map[string]interface{}); ok {
		if nickname, ok := author["nickname"].(string); ok {
			info.Author = nickname
		}
	}

	// 视频播放地址
	if video, ok := awemeDetail["video"].(map[string]interface{}); ok {
		// 尝试获取无水印地址
		if playAddr, ok := video["play_addr"].(map[string]interface{}); ok {
			if urlList, ok := playAddr["url_list"].([]interface{}); ok && len(urlList) > 0 {
				if playURL, ok := urlList[0].(string); ok {
					// 替换为无水印地址
					info.PlayURL = strings.Replace(playURL, "playwm", "play", 1)
				}
			}
		}

		// 备用: bit_rate 中的高清地址
		if info.PlayURL == "" {
			if bitRate, ok := video["bit_rate"].([]interface{}); ok && len(bitRate) > 0 {
				if br, ok := bitRate[0].(map[string]interface{}); ok {
					if playAddr, ok := br["play_addr"].(map[string]interface{}); ok {
						if urlList, ok := playAddr["url_list"].([]interface{}); ok && len(urlList) > 0 {
							if playURL, ok := urlList[0].(string); ok {
								info.PlayURL = playURL
							}
						}
					}
				}
			}
		}
	}

	if info.PlayURL == "" {
		return nil, fmt.Errorf("无法获取视频播放地址")
	}

	return info, nil
}

// downloadDouyinVideo 下载视频文件
func downloadDouyinVideo(playURL, outputPath string) error {
	req, err := http.NewRequest("GET", playURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", DefaultDouyinConfig.UserAgent)
	req.Header.Set("Referer", "https://www.douyin.com/")

	resp, err := douyinClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// isDouyinURL 检查是否是抖音链接
func isDouyinURL(host string) bool {
	douyinHosts := []string{
		"douyin.com",
		"iesdouyin.com",
		"v.douyin.com",
	}
	for _, h := range douyinHosts {
		if strings.Contains(host, h) {
			return true
		}
	}
	return false
}

// SetDouyinCookies 设置抖音 cookies（外部调用）
func SetDouyinCookies(cookies string) {
	DefaultDouyinConfig.Cookies = cookies
}
