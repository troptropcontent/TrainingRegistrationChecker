package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Notifier constants and types
const BASE_NTFY_URL = "https://ntfy.sh/"

type Channel string

const (
	Success Channel = "NTFY_SUCCESS_CHANNEL"
	Logs    Channel = "NTFY_LOGS_CHANNEL"
)

func notify(message string, channel Channel) error {
	log.Printf("[%s] %s", channel, message)

	ntfyChannel := os.Getenv(string(channel))
	if ntfyChannel == "" {
		return fmt.Errorf("environment variable %s is not set", channel)
	}

	ntfyURL := BASE_NTFY_URL + ntfyChannel
	_, err := http.Post(ntfyURL, "text/plain", strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("unable to send notification: %v", err)
	}

	return nil
}

// Registration checker
const TRAINING_URL_ENV = "TRAINING_URL_ENV"

func isRegistrationOpen() (bool, error) {
	c := colly.NewCollector()

	var isOpen bool
	var fetchError error

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		isOpen = !strings.Contains(e.Text, "Candidature fermée.")
	})

	c.OnError(func(r *colly.Response, err error) {
		fetchError = err
	})

	trainingUrl := os.Getenv(TRAINING_URL_ENV)
	if trainingUrl == "" {
		fetchError = fmt.Errorf("environment variable %s is not set", TRAINING_URL_ENV)
	}

	c.Visit(trainingUrl)

	return isOpen, fetchError
}

func main() {
	log.Printf("Spining checker off ...")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		notify("Checking registration status...", Logs)

		open, err := isRegistrationOpen()
		if err != nil {
			notify(fmt.Sprintf("Error checking registration: %v", err), Logs)
			continue
		}

		if open {
			notify("Registration is open!", Success)
			log.Printf("Turning checker down ...")
			return
		}

		notify("Registration is closed!", Logs)
	}
}
