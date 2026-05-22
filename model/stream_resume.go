package model

import (
	"errors"
	"os"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const (
	ResumeStatusRunning = "RUNNING"
	ResumeStatusDone    = "DONE"
	ResumeStatusExpired = "EXPIRED"
	StreamResumeTTL     = 30 * time.Minute
)

type StreamResumeRecord struct {
	ID        string `gorm:"primaryKey;type:varchar(64)"`
	Status    string `gorm:"type:varchar(16);index"`
	Payload   string `gorm:"type:text"`
	CreatedAt int64  `gorm:"bigint;index"`
	UpdatedAt int64  `gorm:"bigint"`
	ExpiresAt int64  `gorm:"bigint;index"`
}

var RESUME_DB *gorm.DB

func ensureResumeDB() error {
	if RESUME_DB != nil {
		return nil
	}
	return InitResumeDB()
}

func InitResumeDB() error {
	if os.Getenv("RESUME_SQL_DSN") == "" {
		db, err := gorm.Open(sqlite.Open("resume.db"), &gorm.Config{PrepareStmt: true})
		if err != nil {
			return err
		}
		RESUME_DB = db
	} else {
		db, err := chooseDB("RESUME_SQL_DSN", true)
		if err != nil {
			return err
		}
		RESUME_DB = db
	}
	return RESUME_DB.AutoMigrate(&StreamResumeRecord{})
}

func CreateRunningResumeRecord(id string) error {
	if !common.IsMasterNode {
		return nil
	}
	if err := ensureResumeDB(); err != nil {
		return err
	}
	now := time.Now().Unix()
	rec := &StreamResumeRecord{ID: id, Status: ResumeStatusRunning, CreatedAt: now, UpdatedAt: now, ExpiresAt: now + int64(StreamResumeTTL.Seconds())}
	return RESUME_DB.Create(rec).Error
}

func CompleteResumeRecord(id string, payload string) {
	if !common.IsMasterNode {
		return
	}
	gopool.Go(func() {
		if err := ensureResumeDB(); err != nil {
			return
		}
		now := time.Now().Unix()
		_ = RESUME_DB.Model(&StreamResumeRecord{}).Where("id = ?", id).Updates(map[string]any{
			"status":     ResumeStatusDone,
			"payload":    payload,
			"updated_at": now,
			"expires_at": now + int64(StreamResumeTTL.Seconds()),
		}).Error
	})
}

func GetResumeRecord(id string) (*StreamResumeRecord, error) {
	if err := ensureResumeDB(); err != nil {
		return nil, err
	}
	var rec StreamResumeRecord
	if err := RESUME_DB.Where("id = ?", id).First(&rec).Error; err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	if rec.ExpiresAt <= now {
		rec.Status = ResumeStatusExpired
	}
	return &rec, nil
}

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
