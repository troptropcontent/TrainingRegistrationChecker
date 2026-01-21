package main

import (
	"log"

	"github.com/troptropcontent/simon_bsc_checker/internal/checker"
	"github.com/troptropcontent/simon_bsc_checker/internal/db"
	"github.com/troptropcontent/simon_bsc_checker/internal/logger"
	"github.com/troptropcontent/simon_bsc_checker/internal/web"
)

// Notifier constants
const SERVER_PORT = 3000

func main() {
	log.Printf("Spinning server off...")

	// Database
	db := db.MustConnect()
	// Logger
	logger := logger.NewLogger(db)
	// Checker
	checker := checker.MustNewChecker(&checker.NewCheckerConfig{Logger: logger})

	go checker.Start()
	// Launch Server
	server := web.NewServer(db, SERVER_PORT)
	server.Initialize()
}
