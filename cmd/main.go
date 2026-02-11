package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/JD-kriswu/video-translator/internal/config"
	"github.com/JD-kriswu/video-translator/internal/downloader"
	"github.com/JD-kriswu/video-translator/internal/extractor"
	"github.com/JD-kriswu/video-translator/internal/transcriber"
	"github.com/JD-kriswu/video-translator/internal/translator"
)

func main() {
	// 命令行参数
	url := flag.String("url", "", "视频 URL")
	output := flag.String("output", "output/result.txt", "输出文件路径")
	lang := flag.String("lang", "zh", "目标语言 (默认: zh)")
	configPath := flag.String("config", "", "配置文件路径 (默认: config.json)")
	keepFiles := flag.Bool("keep", false, "保留中间文件")
	initConfig := flag.Bool("init", false, "生成示例配置文件")
	flag.Parse()

	// 生成示例配置文件
	if *initConfig {
		path := *configPath
		if path == "" {
			path = "config.json"
		}
		if err := config.CreateExample(path); err != nil {
			log.Fatalf("生成配置文件失败: %v", err)
		}
		fmt.Printf("✅ 示例配置文件已生成: %s\n", path)
		fmt.Println("请编辑配置文件，填入你的 API Key")
		return
	}

	if *url == "" {
		fmt.Println("Usage: video-translator -url <video-url> [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\n示例:")
		fmt.Println("  video-translator -init                    # 生成示例配置文件")
		fmt.Println("  video-translator -url https://xxx.mp4     # 翻译视频")
		os.Exit(1)
	}

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 确保输出目录存在
	outputDir := filepath.Dir(*output)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("创建输出目录失败: %v", err)
	}

	fmt.Println("🎬 Video Translator")
	fmt.Println("==================")
	fmt.Printf("转录模型: %s (%s)\n", cfg.Transcriber.Model, cfg.Transcriber.BaseURL)
	fmt.Printf("翻译模型: %s (%s)\n", cfg.Translator.Model, cfg.Translator.BaseURL)

	// 1. 下载视频
	fmt.Println("\n📥 正在下载视频...")
	videoPath, err := downloader.Download(*url)
	if err != nil {
		log.Fatalf("下载失败: %v", err)
	}
	fmt.Printf("✅ 视频已下载: %s\n", videoPath)

	// 2. 提取音频
	fmt.Println("\n🎵 正在提取音频...")
	audioPath, err := extractor.ExtractAudio(videoPath)
	if err != nil {
		log.Fatalf("音频提取失败: %v", err)
	}
	fmt.Printf("✅ 音频已提取: %s\n", audioPath)

	// 3. 语音转文字
	fmt.Println("\n📝 正在转录音频...")
	text, err := transcriber.Transcribe(audioPath, &cfg.Transcriber)
	if err != nil {
		log.Fatalf("转录失败: %v", err)
	}
	fmt.Println("✅ 转录完成")

	// 4. 翻译
	fmt.Println("\n🌐 正在翻译...")
	translated, err := translator.Translate(text, *lang, &cfg.Translator)
	if err != nil {
		log.Fatalf("翻译失败: %v", err)
	}
	fmt.Println("✅ 翻译完成")

	// 5. 保存结果
	if err := os.WriteFile(*output, []byte(translated), 0644); err != nil {
		log.Fatalf("保存失败: %v", err)
	}
	fmt.Printf("\n📄 结果已保存: %s\n", *output)

	// 清理临时文件
	if !*keepFiles {
		os.Remove(videoPath)
		os.Remove(audioPath)
	}

	fmt.Println("\n🎉 完成!")
}
