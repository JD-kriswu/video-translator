package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JD-kriswu/video-translator/internal/config"
	"github.com/JD-kriswu/video-translator/internal/server"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 静态文件
	staticFS, _ := fs.Sub(staticFiles, "static")
	r.StaticFS("/static", http.FS(staticFS))

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static/index.html")
	})

	// API 路由
	api := r.Group("/api")
	{
		// 创建翻译任务
		api.POST("/translate", server.CreateTask(cfg))

		// 查询任务状态
		api.GET("/task/:id", server.GetTaskStatus)

		// 获取任务结果
		api.GET("/task/:id/result", server.GetTaskResult)
	}

	// 启动服务
	port := ":8080"
	log.Printf("🚀 服务启动: http://localhost%s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
