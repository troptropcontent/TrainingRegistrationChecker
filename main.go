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

// Notifier constants
const BASE_NTFY_URL = "https://ntfy.sh/"
const NTFY_CHANNEL_ENV = "NTFY_CHANNEL"

func notify(message string) error {
	ntfyChannel := os.Getenv(NTFY_CHANNEL_ENV)
	if ntfyChannel == "" {
		return fmt.Errorf("environment variable %s is not set", NTFY_CHANNEL_ENV)
	}

	ntfyURL := BASE_NTFY_URL + ntfyChannel
	resp, err := http.Post(ntfyURL, "text/plain", strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("unable to send notification: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("ntfy notification sent (status: %d)", resp.StatusCode)
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
		log.Printf("Checking registration status...")

		open, err := isRegistrationOpen()
		if err != nil {
			log.Printf("Error checking registration: %v", err)
			continue
		}

		if open {
			notify("Registration is open!")
			// log.Printf("Turning checker down ...")
			// return
		}

		log.Printf("Registration is closed.")
	}
}
