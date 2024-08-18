package database

import (
	"billing-engine/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormConnection(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
