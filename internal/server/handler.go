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
	ID         string `json:"id"`
	URL        string `json:"url"`
	Lang       string `json:"lang"`
	Status     string `json:"status"` // pending, downloading, extracting, transcribing, translating, done, error
	Progress   int    `json:"progress"`
	VideoPath  string `json:"video_path,omitempty"`
	AudioPath  string `json:"audio_path,omitempty"`
	Transcript string `json:"transcript,omitempty"` // 原始转录文本
	Result     string `json:"result,omitempty"`     // 翻译结果
	Error      string `json:"error,omitempty"`
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

// ResumeTask 继续执行任务
func ResumeTask(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		tasksMu.RLock()
		task, exists := tasks[id]
		tasksMu.RUnlock()

		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
			return
		}

		if task.Status != "error" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "任务状态不是错误，无需继续"})
			return
		}

		// 重置错误状态，继续执行
		tasksMu.Lock()
		task.Status = "pending"
		task.Error = ""
		tasksMu.Unlock()

		go processTask(task, cfg)

		c.JSON(http.StatusOK, gin.H{"message": "任务已继续执行"})
	}
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

	// 1. 下载视频（如果还没下载）
	if task.VideoPath == "" {
		updateTask("downloading", 10)
		videoPath, err := downloader.Download(task.URL)
		if err != nil {
			setError(err)
			return
		}
		tasksMu.Lock()
		task.VideoPath = videoPath
		tasksMu.Unlock()
	}
	updateTask("downloading", 30)

	// 2. 提取音频（如果还没提取）
	if task.AudioPath == "" {
		updateTask("extracting", 40)
		audioPath, err := extractor.ExtractAudio(task.VideoPath)
		if err != nil {
			setError(err)
			return
		}
		tasksMu.Lock()
		task.AudioPath = audioPath
		tasksMu.Unlock()
	}
	updateTask("extracting", 50)

	// 3. 语音转文字（如果还没转录）
	if task.Transcript == "" {
		updateTask("transcribing", 60)
		text, err := transcriber.Transcribe(task.AudioPath, &cfg.Transcriber)
		if err != nil {
			setError(err)
			return
		}
		tasksMu.Lock()
		task.Transcript = text
		tasksMu.Unlock()
	}
	updateTask("transcribing", 80)

	// 4. 翻译（如果还没翻译）
	if task.Result == "" {
		updateTask("translating", 85)
		translated, err := translator.Translate(task.Transcript, task.Lang, &cfg.Translator)
		if err != nil {
			setError(err)
			return
		}
		tasksMu.Lock()
		task.Result = translated
		tasksMu.Unlock()
	}

	// 完成
	tasksMu.Lock()
	task.Status = "done"
	task.Progress = 100
	tasksMu.Unlock()
}
