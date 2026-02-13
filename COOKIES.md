# Cookies 获取指南

为了提高抖音视频下载的成功率，建议配置真实的浏览器 cookies。

## 方法一：使用浏览器插件（推荐）

### Chrome/Edge 用户

1. 安装插件 [Get cookies.txt LOCALLY](https://chrome.google.com/webstore/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc)

2. 访问 [www.douyin.com](https://www.douyin.com) 并登录（未登录也可以，但登录后成功率更高）

3. 点击插件图标，选择 "Export" 导出 cookies

4. 将导出的文件保存为项目根目录下的 `cookies.txt`

### Firefox 用户

1. 安装插件 [cookies.txt](https://addons.mozilla.org/en-US/firefox/addon/cookies-txt/)

2. 访问 [www.douyin.com](https://www.douyin.com) 并登录

3. 点击插件图标导出 cookies

4. 保存为 `cookies.txt`

## 方法二：手动获取（Chrome DevTools）

1. 打开 Chrome 浏览器，访问 [www.douyin.com](https://www.douyin.com)

2. 按 `F12` 打开开发者工具

3. 切换到 "Application" 标签页

4. 左侧展开 "Cookies" > "https://www.douyin.com"

5. 创建 `cookies.txt` 文件，按照 Netscape Cookie 格式填写：

```
# Netscape HTTP Cookie File
.douyin.com	TRUE	/	FALSE	0	ttwid	YOUR_TTWID_VALUE
.douyin.com	TRUE	/	FALSE	0	msToken	YOUR_MSTOKEN_VALUE
```

关键 cookies 字段：
- `ttwid`: 设备标识
- `msToken`: 签名令牌
- `sessionid`: 会话ID（登录后才有）

## 使用方法

### 方法 1：放在项目根目录（默认）

将 `cookies.txt` 文件放在项目根目录，程序会自动识别。

```
video-translator/
├── cookies.txt          # <- 放这里
├── config.json
└── ...
```

### 方法 2：指定 cookies 文件路径

在代码中设置：

```go
import "video-translator/internal/downloader"

// 设置自定义 cookies 文件路径
downloader.SetCookiesFile("/path/to/your/cookies.txt")

// 然后正常下载
path, err := downloader.Download("https://v.douyin.com/xxx")
```

## 注意事项

1. **cookies 有效期**：cookies 通常在 1-7 天后过期，需要定期更新

2. **不要分享 cookies**：cookies 包含你的登录信息，不要分享给他人

3. **隐私安全**：
   - 添加 `cookies.txt` 到 `.gitignore`
   - 不要将 cookies 提交到代码仓库

4. **成功率提升**：
   - 未登录：成功率约 60-70%
   - 已登录：成功率约 80-90%

## 验证 cookies 是否生效

运行下载命令时，如果看到：

```
使用 cookies 文件: cookies.txt
✓ 使用 yt-dlp 下载成功
```

说明 cookies 已生效。

## 故障排查

### 问题：下载失败，提示 "HTTP Error 403"

**解决方案**：
1. 检查 cookies 是否过期，重新导出
2. 确保 cookies.txt 格式正确
3. 尝试登录抖音账号后重新导出

### 问题：无法找到 cookies.txt

**解决方案**：
1. 确认文件放在项目根目录
2. 检查文件名是否正确（区分大小写）
3. 使用绝对路径：`downloader.SetCookiesFile("/full/path/to/cookies.txt")`

## 高级技巧

### 自动更新 cookies

可以编写脚本定期从浏览器导出 cookies：

```bash
#!/bin/bash
# update-cookies.sh

# 从 Chrome 浏览器配置目录复制 cookies
# macOS Chrome cookies 位置
CHROME_COOKIES="$HOME/Library/Application Support/Google/Chrome/Default/Cookies"

# 使用 sqlite3 导出 cookies（需要安装 sqlite3）
sqlite3 "$CHROME_COOKIES" "SELECT host_key, TRUE, path, is_secure, expires_utc, name, value FROM cookies WHERE host_key like '%douyin.com%'" > cookies.txt
```

### 使用代理 + cookies（终极方案）

如果单独使用 cookies 仍失败，可以配合代理使用：

```bash
# 在 yt-dlp 命令中添加代理
yt-dlp --cookies cookies.txt --proxy "socks5://127.0.0.1:1080" "VIDEO_URL"
```

---

**提示**：如果你不想使用 cookies，程序会尝试其他下载方式（Playwright、API），只是成功率会降低。
