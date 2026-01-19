package main

import (
	"fmt"
	"log"
	"time"

	"github.com/troptropcontent/simon_bsc_checker/internal/notifier"
	"github.com/troptropcontent/simon_bsc_checker/internal/registration"
)

func main() {
	log.Printf("Spining checker off ...")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		notifier.Notify("Checking registration status...", notifier.Logs)

		open, err := registration.IsOpen()
		if err != nil {
			notifier.Notify(fmt.Sprintf("Error checking registration: %v", err), notifier.Logs)

			continue
		}

		if open {
			notifier.Notify("Registration is open!", notifier.Success)

			log.Printf("Turning checker down ...")
			return
		}

		notifier.Notify("Registration is closed!", notifier.Logs)
	}
}
