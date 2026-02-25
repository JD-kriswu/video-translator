package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/JD-kriswu/video-translator/internal/database"
)

// Translation 翻译记录
type Translation struct {
	ID         string         `gorm:"primarykey;size:36" json:"id"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	URL        string         `gorm:"size:500;not null" json:"url"`
	Lang       string         `gorm:"size:10" json:"lang"`
	Status     string         `gorm:"size:20;not null;default:pending" json:"status"`
	Progress   int            `gorm:"default:0" json:"progress"`
	VideoPath  string         `gorm:"size:500" json:"video_path,omitempty"`
	AudioPath  string         `gorm:"size:500" json:"audio_path,omitempty"`
	Transcript string         `gorm:"type:text" json:"transcript,omitempty"`
	Result     string         `gorm:"type:text" json:"result,omitempty"`
	Error      string         `gorm:"type:text" json:"error,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AutoMigrateTranslation 自动迁移翻译表
func AutoMigrateTranslation() error {
	return database.GetDB().AutoMigrate(&Translation{})
}

// CreateTranslation 创建翻译记录
func CreateTranslation(userID uint, id, url, lang string) (*Translation, error) {
	translation := &Translation{
		ID:     id,
		UserID: userID,
		URL:    url,
		Lang:   lang,
		Status: "pending",
	}

	result := database.GetDB().Create(translation)
	return translation, result.Error
}

// UpdateTranslation 更新翻译记录
func UpdateTranslation(translation *Translation) error {
	return database.GetDB().Save(translation).Error
}

// GetTranslationByID 根据 ID 获取翻译记录
func GetTranslationByID(id string) (*Translation, error) {
	var translation Translation
	result := database.GetDB().Where("id = ?", id).First(&translation)
	return &translation, result.Error
}

// GetUserTranslations 获取用户的翻译列表
func GetUserTranslations(userID uint, limit, offset int) ([]Translation, int64, error) {
	var translations []Translation
	var total int64

	db := database.GetDB().Model(&Translation{}).Where("user_id = ?", userID)

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	result := db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&translations)

	return translations, total, result.Error
}
