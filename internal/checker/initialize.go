package checker

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/troptropcontent/simon_bsc_checker/internal/logger"
)

const NTFY_URL_ENV = "NTFY_URL"
const TRAINING_PAGE_URL_ENV = "TRAINING_PAGE_URL"
const DEFAULT_TICKER_INTERVAL = 1 * time.Minute

type Checker struct {
	NtfyUrl         string
	TrainingPageUrl string
	TickerInterval  time.Duration
	Logger          *logger.Logger
}

func (checker *Checker) notify(message string) error {
	resp, err := http.Post(checker.NtfyUrl, "text/plain", strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("unable to send notification: %v", err)
	}
	defer resp.Body.Close()

	checker.Logger.Info("ntfy notification sent (status: %d)", resp.StatusCode)

	return nil
}

func (checker *Checker) isRegistrationOpen() (bool, error) {
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

	c.Visit(checker.TrainingPageUrl)

	return isOpen, fetchError
}

type NewCheckerConfig struct {
	NtfyUrl         string
	TrainingPageUrl string
	TickerInterval  time.Duration
	Logger          *logger.Logger
}

func NewChecker(c *NewCheckerConfig) (*Checker, error) {
	ntfyUrl := c.NtfyUrl
	if ntfyUrl == "" {
		ntfyUrl = os.Getenv(NTFY_URL_ENV)
	}
	if ntfyUrl == "" {
		return nil, fmt.Errorf("could not determine the ntfy url, it must either be set through the configs or with the %v environement variable", NTFY_URL_ENV)
	}

	trainingPageUrl := c.TrainingPageUrl
	if trainingPageUrl == "" {
		trainingPageUrl = os.Getenv(TRAINING_PAGE_URL_ENV)
	}
	if trainingPageUrl == "" {
		return nil, fmt.Errorf("could not determine the training page URL, it must either be set through the configs or with the %v environement variable", TRAINING_PAGE_URL_ENV)
	}

	tickerInterval := c.TickerInterval
	if tickerInterval == 0 {
		tickerInterval = DEFAULT_TICKER_INTERVAL
	}

	return &Checker{
		NtfyUrl:         ntfyUrl,
		TrainingPageUrl: trainingPageUrl,
		TickerInterval:  tickerInterval,
		Logger:          c.Logger,
	}, nil
}

func MustNewChecker(c *NewCheckerConfig) *Checker {
	checker, err := NewChecker(c)
	if err != nil {
		log.Fatalf("failed to initialize checker : %v", err)
	}
	return checker
}

func (checker *Checker) Start() {
	ticker := time.NewTicker(checker.TickerInterval)
	defer ticker.Stop()

	for range ticker.C {
		checker.Logger.Info("Checking registration status...")

		open, err := checker.isRegistrationOpen()
		if err != nil {
			checker.Logger.Error("Error checking registration: %v", err)
			continue
		}

		if open {
			checker.Logger.Info("Registration is open!")
			if err := checker.notify("Registration is open!"); err != nil {
				checker.Logger.Error("Failed to send notification: %v", err)
			}
			continue
		}

		checker.Logger.Info("Registration is closed.")
	}
}
