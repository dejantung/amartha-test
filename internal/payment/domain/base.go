package domain

import (
	"gorm.io/gorm"
	"time"
)

type Base struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	base.CreatedAt = time.Now()
	base.UpdatedAt = base.CreatedAt
	return
}

func (base *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	base.UpdatedAt = time.Now()
	return
}
