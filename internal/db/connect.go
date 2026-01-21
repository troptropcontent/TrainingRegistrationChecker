package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DATABASE_FILE = "database.db"

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type Log struct {
	gorm.Model
	Level   LogLevel `gorm:"index;default:info"`
	Message string
}

func MustConnect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(DATABASE_FILE), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&Log{})

	return db
}
