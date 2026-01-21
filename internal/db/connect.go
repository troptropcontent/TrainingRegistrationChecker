package db

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const DEFAULT_DATABASE_FILE = "database.db"

func getDatabasePath() string {
	if path := os.Getenv("DATABASE_PATH"); path != "" {
		return path
	}
	return DEFAULT_DATABASE_FILE
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

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
	dbPath := getDatabasePath()

	if err := ensureDir(dbPath); err != nil {
		log.Fatalf("failed to create database directory: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	db.AutoMigrate(&Log{})

	return db
}
