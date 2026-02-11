package server

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/JD-kriswu/video-translator/internal/config"
	"github.com/JD-kriswu/video-translator/internal/downloader"
	"github.com/JD-kriswu/video-translator/internal/extractor"
	"github.com/JD-kriswu/video-translator/internal/transcriber"
	"github.com/JD-kriswu/video-translator/internal/translator"
)

// Task 翻译任务
type Task struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Lang     string `json:"lang"`
	Status   string `json:"status"` // pending, downloading, extracting, transcribing, translating, done, error
	Progress int    `json:"progress"`
	Result   string `json:"result,omitempty"`
	Error    string `json:"error,omitempty"`
}

// 任务存储
var (
	tasks   = make(map[string]*Task)
	tasksMu sync.RWMutex
)

// TranslateRequest 翻译请求
type TranslateRequest struct {
	URL  string `json:"url" binding:"required"`
	Lang string `json:"lang"`
}

// CreateTask 创建翻译任务
func CreateTask(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TranslateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请提供视频 URL"})
			return
		}

		// 默认语言
		if req.Lang == "" {
			req.Lang = "zh"
		}

		// 创建任务
		task := &Task{
			ID:       uuid.New().String(),
			URL:      req.URL,
			Lang:     req.Lang,
			Status:   "pending",
			Progress: 0,
		}

		tasksMu.Lock()
		tasks[task.ID] = task
		tasksMu.Unlock()

		// 异步处理
		go processTask(task, cfg)

		c.JSON(http.StatusOK, gin.H{
			"id":      task.ID,
			"message": "任务已创建",
		})
	}
}

// GetTaskStatus 获取任务状态
func GetTaskStatus(c *gin.Context) {
	id := c.Param("id")

	tasksMu.RLock()
	task, exists := tasks[id]
	tasksMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// GetTaskResult 获取任务结果
func GetTaskResult(c *gin.Context) {
	id := c.Param("id")

	tasksMu.RLock()
	task, exists := tasks[id]
	tasksMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if task.Status != "done" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务未完成", "status": task.Status})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": task.Result,
	})
}

// processTask 处理翻译任务
func processTask(task *Task, cfg *config.Config) {
	updateTask := func(status string, progress int) {
		tasksMu.Lock()
		task.Status = status
		task.Progress = progress
		tasksMu.Unlock()
	}

	setError := func(err error) {
		tasksMu.Lock()
		task.Status = "error"
		task.Error = err.Error()
		tasksMu.Unlock()
	}

	// 1. 下载视频
	updateTask("downloading", 10)
	videoPath, err := downloader.Download(task.URL)
	if err != nil {
		setError(err)
		return
	}
	updateTask("downloading", 30)

	// 2. 提取音频
	updateTask("extracting", 40)
	audioPath, err := extractor.ExtractAudio(videoPath)
	if err != nil {
		setError(err)
		return
	}
	updateTask("extracting", 50)

	// 3. 语音转文字
	updateTask("transcribing", 60)
	text, err := transcriber.Transcribe(audioPath, &cfg.Transcriber)
	if err != nil {
		setError(err)
		return
	}
	updateTask("transcribing", 80)

	// 4. 翻译
	updateTask("translating", 85)
	translated, err := translator.Translate(text, task.Lang, &cfg.Translator)
	if err != nil {
		setError(err)
		return
	}

	// 完成
	tasksMu.Lock()
	task.Status = "done"
	task.Progress = 100
	task.Result = translated
	tasksMu.Unlock()
}
