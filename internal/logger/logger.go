package logger

import (
	"fmt"
	"log"

	"github.com/troptropcontent/simon_bsc_checker/internal/db"
	"gorm.io/gorm"
)

type Logger struct {
	DB *gorm.DB
}

func NewLogger(database *gorm.DB) *Logger {
	return &Logger{
		DB: database,
	}
}

func (l *Logger) log(level db.LogLevel, format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	log.Printf("[%s] %s", level, message)
	if err := l.DB.Create(&db.Log{Level: level, Message: message}).Error; err != nil {
		log.Printf("[logger] failed to persist log: %v", err)
	}
}

func (l *Logger) Info(format string, v ...any) {
	l.log(db.LogLevelInfo, format, v...)
}

func (l *Logger) Warn(format string, v ...any) {
	l.log(db.LogLevelWarn, format, v...)
}

func (l *Logger) Error(format string, v ...any) {
	l.log(db.LogLevelError, format, v...)
}
