package model

import (
	"errors"
	"os"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type ConversationLog struct {
	Id         int    `json:"id" gorm:"primaryKey"`
	CreatedAt  int64  `json:"created_at" gorm:"bigint;index:idx_conv_created_at"`
	UserId     int    `json:"user_id" gorm:"index:idx_conv_user_token"`
	Username   string `json:"username" gorm:"index;default:''"`
	TokenId    int    `json:"token_id" gorm:"index:idx_conv_user_token"`
	TokenName  string `json:"token_name" gorm:"index;default:''"`
	ModelName  string `json:"model_name" gorm:"index;default:''"`
	RequestId  string `json:"request_id" gorm:"type:varchar(64);index;default:''"`
	PromptText string `json:"prompt_text" gorm:"type:text"`
	ReplyText  string `json:"reply_text" gorm:"type:text"`
}

var CONV_LOG_DB = LOG_DB

func ensureConversationLogDB() error {
	if CONV_LOG_DB != nil {
		return nil
	}
	return InitConversationLogDB()
}

func InitConversationLogDB() error {
	if os.Getenv("CONV_LOG_SQL_DSN") == "" {
		db, err := gorm.Open(sqlite.Open("oneapi-conversation-logs.db"), &gorm.Config{PrepareStmt: true})
		if err != nil {
			return err
		}
		CONV_LOG_DB = db
	} else {
		db, err := chooseDB("CONV_LOG_SQL_DSN", true)
		if err != nil {
			return err
		}
		CONV_LOG_DB = db
	}
	return CONV_LOG_DB.AutoMigrate(&ConversationLog{})
}

func RecordConversationLogAsync(c ConversationLog) {
	if !common.IsMasterNode || !setting.ConversationLogEnabled {
		return
	}
	if c.PromptText == "" && c.ReplyText == "" {
		return
	}
	if c.CreatedAt == 0 {
		c.CreatedAt = time.Now().Unix()
	}
	gopool.Go(func() {
		if err := ensureConversationLogDB(); err != nil {
			return
		}
		_ = CONV_LOG_DB.Create(&c).Error
	})
}

func GetConversationLogs(offset int, limit int, username string, tokenName string, modelName string) ([]*ConversationLog, int64, error) {
	if err := ensureConversationLogDB(); err != nil {
		return nil, 0, errors.New("conversation log db not initialized: " + err.Error())
	}
	tx := CONV_LOG_DB.Model(&ConversationLog{})
	if username != "" {
		tx = tx.Where("username = ?", username)
	}
	if tokenName != "" {
		tx = tx.Where("token_name = ?", tokenName)
	}
	if modelName != "" {
		tx = tx.Where("model_name = ?", modelName)
	}
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	items := make([]*ConversationLog, 0, limit)
	err := tx.Order("id desc").Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}
