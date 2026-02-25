package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/JD-kriswu/video-translator/internal/model"
)

// GetUserTranslations 获取当前用户的翻译列表
func GetUserTranslations(c *gin.Context) {
	// 从上下文获取用户ID（由 AuthMiddleware 设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 查询数据库
	translations, total, err := model.GetUserTranslations(userID.(uint), pageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": translations,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetTranslationByID 根据ID获取翻译记录
func GetTranslationByID(c *gin.Context) {
	id := c.Param("id")

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 查询翻译记录
	translation, err := model.GetTranslationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "翻译记录不存在"})
		return
	}

	// 验证所有权
	if translation.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此记录"})
		return
	}

	c.JSON(http.StatusOK, translation)
}
