package domain

import (
	"gorm.io/gorm"
	"time"
)

type AuditLog struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (auditLog *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	auditLog.CreatedAt = time.Now()
	return
}

func (auditLog *AuditLog) BeforeUpdate(tx *gorm.DB) (err error) {
	auditLog.UpdatedAt = time.Now()
	return
}
