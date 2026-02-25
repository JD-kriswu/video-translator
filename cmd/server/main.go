package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JD-kriswu/video-translator/internal/auth"
	"github.com/JD-kriswu/video-translator/internal/config"
	"github.com/JD-kriswu/video-translator/internal/database"
	"github.com/JD-kriswu/video-translator/internal/model"
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

	// 初始化 MySQL
	if err := database.InitMySQL(&cfg.MySQL); err != nil {
		log.Fatalf("初始化 MySQL 失败: %v", err)
	}
	log.Println("✅ MySQL 连接成功")

	// 初始化 Redis
	if err := database.InitRedis(&cfg.Redis); err != nil {
		log.Fatalf("初始化 Redis 失败: %v", err)
	}
	log.Println("✅ Redis 连接成功")

	// 自动迁移数据库
	if err := model.AutoMigrate(); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	if err := model.AutoMigrateTranslation(); err != nil {
		log.Fatalf("翻译表迁移失败: %v", err)
	}
	log.Println("✅ 数据库迁移完成")

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
		// 认证相关 (无需登录)
		api.POST("/register", auth.Register)
		api.POST("/login", auth.Login)
		api.POST("/logout", auth.Logout)
		api.GET("/me", auth.GetCurrentUser)

		// 需要登录的接口
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			// 创建翻译任务
			protected.POST("/translate", server.CreateTask(cfg))

			// 查询任务状态
			protected.GET("/task/:id", server.GetTaskStatus)

			// 继续执行任务
			protected.POST("/task/:id/resume", server.ResumeTask(cfg))

			// 获取任务结果
			protected.GET("/task/:id/result", server.GetTaskResult)

			// 用户翻译列表
			protected.GET("/user/translations", server.GetUserTranslations)

			// 翻译记录详情
			protected.GET("/translations/:id", server.GetTranslationByID)
		}
	}

	// 启动服务
	port := ":80"
	log.Printf("🚀 服务启动: http://localhost%s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("启动失败: %v", err)
	}
}
